// Command probe diagnoses GL.iNet firmware 4.x authentication and records
// verbatim JSON-RPC exchanges, so that group and method names can be captured
// rather than guessed.
//
// It deliberately reimplements the login rather than calling src/auth, so that a
// disagreement between this program and the library isolates where the fault is.
// The digest variants below are the plausible readings of the scheme; the router
// decides which is right.
//
// Usage:
//
//	go run ./discovery -H 192.168.8.1                       # diagnose login
//	go run ./discovery -H 192.168.8.1 -group system -method get_status
//
// Output may contain a live session id. Review before pasting anywhere public.
package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/md5_crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
)

func main() {
	host := flag.String("H", os.Getenv("GL_ROUTER_IP"), "router address")
	port := flag.Int("p", 80, "router port")
	user := flag.String("u", envOr("GL_USERNAME", "root"), "username")
	group := flag.String("group", "", "functional group to probe (omit to diagnose login only)")
	method := flag.String("method", "", "method to probe")
	argsJSON := flag.String("args", "{}", "JSON object of arguments")
	batch := flag.String("batch", "", "file of 'group method' lines, or - for stdin; one login for all of them")
	full := flag.Bool("full", false, "with -batch, print each full result rather than a one-line summary")
	variants := flag.Bool("variants", false,
		"sweep every digest variant instead of trusting hash-method; each wrong guess "+
			"counts against the router's login-failure limit")
	flag.Parse()

	if *host == "" {
		fmt.Fprintln(os.Stderr, "no host: pass -H or set GL_ROUTER_IP")
		os.Exit(1)
	}
	password := os.Getenv("GL_PASSWORD")
	if password == "" {
		fmt.Fprintln(os.Stderr, "GL_PASSWORD is not set")
		os.Exit(1)
	}

	url := fmt.Sprintf("http://%s:%d/rpc", *host, *port)
	fmt.Printf("endpoint: %s\nusername: %s\n\n", url, *user)

	sid, err := diagnose(url, *user, password, *variants)
	if err != nil {
		fmt.Fprintln(os.Stderr, "\nFAILED:", err)
		os.Exit(1)
	}

	if *batch != "" {
		if err := runBatch(url, sid, *batch, *full); err != nil {
			fmt.Fprintln(os.Stderr, "batch:", err)
			os.Exit(1)
		}
		return
	}

	if *group == "" {
		fmt.Printf("\nLogin succeeded. Probe an endpoint with:\n"+
			"  go run ./discovery -H %s -group <group> -method <method>\n"+
			"  go run ./discovery -H %s -batch candidates.txt\n", *host, *host)
		return
	}

	var args any
	if err := json.Unmarshal([]byte(*argsJSON), &args); err != nil {
		fmt.Fprintln(os.Stderr, "parse -args:", err)
		os.Exit(1)
	}

	raw, err := post(url, map[string]any{
		"jsonrpc": "2.0", "id": 99, "method": "call",
		"params": []any{sid, *group, *method, args},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "call:", err)
		os.Exit(1)
	}
	fmt.Printf("\n=== %s.%s ===\n%s\n", *group, *method, indent(raw))
}

// challengeResult is decoded loosely: the point is to see what the router really
// sends, including any field this project does not yet know about.
type challengeResult struct {
	Alg        json.RawMessage `json:"alg"`
	Salt       string          `json:"salt"`
	Nonce      string          `json:"nonce"`
	HashMethod string          `json:"hash-method"`
}

func diagnose(url, user, password string, tryVariants bool) (string, error) {
	fmt.Println("=== step 1: challenge ===")
	raw, err := post(url, map[string]any{
		"jsonrpc": "2.0", "id": 1, "method": "challenge",
		"params": map[string]any{"username": user},
	})
	if err != nil {
		return "", err
	}
	fmt.Println(indent(raw))

	var envelope struct {
		Result challengeResult `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return "", fmt.Errorf("decode challenge: %w", err)
	}
	if envelope.Error != nil {
		return "", fmt.Errorf("challenge rejected: %s (code %d)", envelope.Error.Message, envelope.Error.Code)
	}

	first := envelope.Result
	alg := strings.Trim(string(first.Alg), `"`)
	fmt.Printf("\nalg=%s  salt=%q  nonce=%q\n", alg, first.Salt, first.Nonce)
	if first.Salt == "" {
		return "", fmt.Errorf("router returned an empty salt")
	}

	// The salt is where undocumented behavior hides. A shadow entry may carry a
	// non-default rounds= prefix, and crypt truncates a SHA-512 salt at 16 bytes,
	// so an unusual length or an embedded $ changes the answer.
	fmt.Printf("salt length: %d bytes\n", len(first.Salt))
	if strings.Contains(first.Salt, "rounds=") {
		fmt.Println("NOTE: salt carries a rounds= prefix -- non-default cost, this matters")
	}
	if strings.Contains(first.Salt, "$") {
		fmt.Println("NOTE: salt contains a '$' -- it may already be a full crypt setting")
	}
	if len(first.Salt) > 16 {
		fmt.Println("NOTE: salt is longer than 16 bytes, which SHA-512 crypt truncates")
	}

	cipher, err := unixCrypt(password, first.Salt, alg)
	if err != nil {
		return "", err
	}
	fmt.Printf("\ncrypt output: %s\n", cipher)

	tail := cipher
	if i := strings.LastIndex(cipher, "$"); i >= 0 {
		tail = cipher[i+1:]
	}

	// If the router handed back a full setting rather than a bare salt, crypt it
	// verbatim as well.
	verbatim := cipher
	if strings.HasPrefix(first.Salt, "$") {
		if v, err := crypterFor(alg); err == nil {
			if out, err := v.Generate([]byte(password), []byte(first.Salt)); err == nil {
				verbatim = out
			}
		}
	}

	// Truncated salt, in case the stored hash used a shorter one than the
	// challenge reports.
	shortSalt := first.Salt
	if len(shortSalt) > 8 {
		shortSalt = shortSalt[:8]
	}
	shortCipher, _ := unixCrypt(password, shortSalt, alg)

	// Each variant is a plausible reading of the scheme. The router arbitrates.
	variants := []struct {
		name   string
		digest func(nonce string) string
	}{
		{"md5(user:fullcrypt:nonce)  [what gogl currently sends]",
			func(n string) string { return md5hex(user + ":" + cipher + ":" + n) }},
		{"md5(user:crypttail:nonce)",
			func(n string) string { return md5hex(user + ":" + tail + ":" + n) }},
		{"md5(fullcrypt:nonce)",
			func(n string) string { return md5hex(cipher + ":" + n) }},
		{"md5(user:password:nonce)",
			func(n string) string { return md5hex(user + ":" + password + ":" + n) }},
		{"sha256(user:fullcrypt:nonce)",
			func(n string) string { return sha256hex(user + ":" + cipher + ":" + n) }},
		{"MD5(user:fullcrypt:nonce) uppercase hex",
			func(n string) string { return strings.ToUpper(md5hex(user + ":" + cipher + ":" + n)) }},
		{"md5(user:fullcrypt:nonce) with salt crypted verbatim",
			func(n string) string { return md5hex(user + ":" + verbatim + ":" + n) }},
		{"md5(user:fullcrypt:nonce) with salt truncated to 8 bytes",
			func(n string) string { return md5hex(user + ":" + shortCipher + ":" + n) }},
		{"md5(nonce:fullcrypt:user)  [reversed order]",
			func(n string) string { return md5hex(n + ":" + cipher + ":" + user) }},
	}

	// Honor the router's own answer first. Sweeping variants costs a failed login
	// per wrong guess, and the firmware locks the account out after a handful of
	// failures -- ten minutes of lockout is a steep price for a guess the router
	// already told us how to avoid.
	if !tryVariants {
		if first.HashMethod == "" {
			fmt.Println("\nNo hash-method in the challenge; assuming md5 (pre-4.8 firmware).")
		} else {
			fmt.Printf("\nRouter advertises hash-method=%q; using it.\n", first.HashMethod)
		}

		sid, err := attemptLogin(url, user, cipher, first.HashMethod)
		if err == nil {
			fmt.Println("\n=== login succeeded ===")
			return sid, nil
		}
		return "", fmt.Errorf("login with the advertised hash-method failed: %w\n"+
			"  re-run with -variants to sweep the alternatives, but note that each\n"+
			"  wrong guess counts against the router's login-failure limit", err)
	}

	fmt.Println("\n=== step 2: trying each digest variant ===")
	fmt.Println("(each wrong guess counts against the router's login-failure limit,")
	fmt.Println(" which locks the account for ~10 minutes once tripped)")

	for _, v := range variants {
		// A fresh challenge per attempt: the nonce is single-use, so reusing one
		// would fail every variant after the first for the wrong reason.
		fresh, err := freshNonce(url, user)
		if err != nil {
			return "", err
		}

		body, err := post(url, map[string]any{
			"jsonrpc": "2.0", "id": 2, "method": "login",
			"params": map[string]any{"username": user, "hash": v.digest(fresh)},
		})
		if err != nil {
			return "", err
		}

		var loginEnvelope struct {
			Result struct {
				SID string `json:"sid"`
			} `json:"result"`
			Error *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &loginEnvelope); err != nil {
			return "", fmt.Errorf("decode login: %w", err)
		}

		switch {
		case loginEnvelope.Result.SID != "":
			fmt.Printf("  OK    %s\n", v.name)
			fmt.Printf("\n=== step 3: this is the scheme your firmware uses ===\n%s\n", v.name)
			return loginEnvelope.Result.SID, nil
		case loginEnvelope.Error != nil:
			fmt.Printf("  no    %s  (%s)\n", v.name, loginEnvelope.Error.Message)
		default:
			fmt.Printf("  no    %s  (no sid, no error: %s)\n", v.name, string(body))
		}
	}

	return "", fmt.Errorf("no digest variant was accepted; the password may simply be wrong, " +
		"or this firmware uses a scheme not listed above")
}

// attemptLogin performs one login with the given hash method, taking a fresh
// nonce because a nonce is single-use and short-lived.
func attemptLogin(url, user, cipher, hashMethod string) (string, error) {
	nonce, err := freshNonce(url, user)
	if err != nil {
		return "", err
	}

	var digest string
	switch hashMethod {
	case "", "md5":
		digest = md5hex(user + ":" + cipher + ":" + nonce)
	case "sha256":
		digest = sha256hex(user + ":" + cipher + ":" + nonce)
	default:
		return "", fmt.Errorf("unsupported hash-method %q -- record this; it is new", hashMethod)
	}

	body, err := post(url, map[string]any{
		"jsonrpc": "2.0", "id": 2, "method": "login",
		"params": map[string]any{"username": user, "hash": digest},
	})
	if err != nil {
		return "", err
	}

	var envelope struct {
		Result struct {
			SID string `json:"sid"`
		} `json:"result"`
		Error *struct {
			Code    int             `json:"code"`
			Message string          `json:"message"`
			Data    json.RawMessage `json:"data"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return "", err
	}
	if envelope.Error != nil {
		if len(envelope.Error.Data) > 0 {
			return "", fmt.Errorf("%s (code %d, data %s)",
				envelope.Error.Message, envelope.Error.Code, envelope.Error.Data)
		}
		return "", fmt.Errorf("%s (code %d)", envelope.Error.Message, envelope.Error.Code)
	}
	if envelope.Result.SID == "" {
		return "", fmt.Errorf("no sid in response: %s", body)
	}
	return envelope.Result.SID, nil
}

func freshNonce(url, user string) (string, error) {
	raw, err := post(url, map[string]any{
		"jsonrpc": "2.0", "id": 1, "method": "challenge",
		"params": map[string]any{"username": user},
	})
	if err != nil {
		return "", err
	}
	var envelope struct {
		Result challengeResult `json:"result"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return "", err
	}
	return envelope.Result.Nonce, nil
}

func crypterFor(alg string) (crypt.Crypter, error) {
	switch alg {
	case "1":
		return crypt.MD5.New(), nil
	case "5":
		return crypt.SHA256.New(), nil
	case "6":
		return crypt.SHA512.New(), nil
	default:
		return nil, fmt.Errorf("unsupported alg %q -- this is new information; record it", alg)
	}
}

func unixCrypt(password, salt, alg string) (string, error) {
	c, err := crypterFor(alg)
	if err != nil {
		return "", err
	}
	return c.Generate([]byte(password), []byte("$"+alg+"$"+salt))
}

func md5hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func sha256hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func post(url string, body any) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("HTTP %s: %s", resp.Status, strings.TrimSpace(string(out)))
	}
	return out, nil
}

// indent pretty-prints JSON, falling back to the raw bytes when the response is
// not JSON at all -- which is itself worth seeing.
func indent(raw []byte) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "  ", "  "); err != nil {
		return "  " + string(raw)
	}
	return "  " + buf.String()
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// runBatch probes many group/method pairs on one session. A full login costs 5000
// crypt rounds, so doing it per probe makes broad discovery painfully slow.
func runBatch(url, sid, source string, full bool) error {
	var input io.Reader = os.Stdin
	if source != "-" {
		f, err := os.Open(source)
		if err != nil {
			return err
		}
		defer f.Close()
		input = f
	}

	scanner := bufio.NewScanner(input)
	found := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("  SKIP  %q (want 'group method')\n", line)
			continue
		}
		group, method := fields[0], fields[1]

		raw, err := post(url, map[string]any{
			"jsonrpc": "2.0", "id": 99, "method": "call",
			"params": []any{sid, group, method, map[string]any{}},
		})
		if err != nil {
			fmt.Printf("  ERR   %s.%s: %v\n", group, method, err)
			continue
		}

		var envelope struct {
			Result json.RawMessage `json:"result"`
			Error  *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(raw, &envelope); err != nil {
			fmt.Printf("  ERR   %s.%s: undecodable: %s\n", group, method, raw)
			continue
		}

		if envelope.Error != nil {
			fmt.Printf("  no    %-28s %s (%d)\n", group+"."+method,
				envelope.Error.Message, envelope.Error.Code)
			continue
		}

		found++
		if full {
			fmt.Printf("  OK    %s\n%s\n", group+"."+method, indent(envelope.Result))
			continue
		}
		fmt.Printf("  OK    %-28s %s\n", group+"."+method, summarize(envelope.Result))
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("\n%d endpoint(s) responded.\n", found)
	return nil
}

// summarize renders a one-line shape of a result: the top-level keys, which is
// what identifies a useful endpoint at a glance.
func summarize(raw json.RawMessage) string {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil {
		flat := strings.Join(strings.Fields(string(raw)), " ")
		if len(flat) > 90 {
			flat = flat[:90] + "..."
		}
		return flat
	}
	keys := make([]string, 0, len(object))
	for k := range object {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return "keys: " + strings.Join(keys, ", ")
}
