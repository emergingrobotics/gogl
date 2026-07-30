package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Authenticator performs the challenge/response login. It caches the derived
// cipher, which depends only on password and salt, so that re-authentication
// does not repeat the deliberately slow crypt step.
type Authenticator struct {
	httpClient *http.Client
	url        string
	username   string
	password   string

	mu     sync.Mutex
	cipher string
	salt   string
}

func NewAuthenticator(httpClient *http.Client, url, username, password string) *Authenticator {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Authenticator{httpClient: httpClient, url: url, username: username, password: password}
}

type challengeResult struct {
	Alg   json.RawMessage `json:"alg"`
	Salt  string          `json:"salt"`
	Nonce string          `json:"nonce"`

	// HashMethod selects the digest applied to the login string. Firmware that
	// predates the field omits it and expects MD5; 4.3.28 on the SFT1200 reports
	// "sha256". It is a separate choice from Alg, which selects the crypt.
	HashMethod string `json:"hash-method"`
}

// Login returns a fresh session id.
//
// It issues challenge twice on purpose. The first supplies the salt, from which
// the cipher is derived at a cost of 5000 hash rounds; by the time that
// completes, the first challenge's nonce has very likely expired, since the
// router keeps it alive for only about a second. The second challenge supplies a
// nonce that is fresh at the moment the cheap MD5 digest is computed over it.
//
// A single-challenge implementation races against its own crypt cost and fails
// intermittently, which is far worse than failing consistently.
func (a *Authenticator) Login(ctx context.Context) (string, error) {
	first, err := a.challenge(ctx)
	if err != nil {
		return "", err
	}

	alg, err := decodeAlg(first.Alg)
	if err != nil {
		return "", err
	}

	cipher, err := a.cipherFor(first.Salt, alg)
	if err != nil {
		return "", err
	}

	second, err := a.challenge(ctx)
	if err != nil {
		return "", err
	}
	// A salt change would invalidate the cached cipher. It should not happen, but
	// signing with a stale salt would surface as a wrong password, which is a
	// miserable thing to debug.
	if second.Salt != first.Salt {
		cipher, err = a.cipherFor(second.Salt, alg)
		if err != nil {
			return "", err
		}
	}

	hash, err := LoginHash(a.username, cipher, second.Nonce, second.HashMethod)
	if err != nil {
		return "", err
	}

	var result struct {
		SID string `json:"sid"`
	}
	params := map[string]any{
		"username": a.username,
		"hash":     hash,
	}
	if err := a.post(ctx, "login", params, &result); err != nil {
		return "", fmt.Errorf("%w: %w", ErrUnauthorized, err)
	}
	if result.SID == "" {
		return "", fmt.Errorf("%w: router returned no session id", ErrUnauthorized)
	}
	return result.SID, nil
}

// cipherFor returns the cached cipher, deriving it if the salt has changed.
func (a *Authenticator) cipherFor(salt string, alg int) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cipher != "" && a.salt == salt {
		return a.cipher, nil
	}
	cipher, err := Cipher(a.password, salt, alg)
	if err != nil {
		return "", err
	}
	a.cipher, a.salt = cipher, salt
	return cipher, nil
}

func (a *Authenticator) challenge(ctx context.Context) (challengeResult, error) {
	var result challengeResult
	err := a.post(ctx, "challenge", map[string]any{"username": a.username}, &result)
	return result, err
}

// decodeAlg tolerates alg arriving as a JSON number or a JSON string, both of
// which have been observed from firmware in the wild.
func decodeAlg(raw json.RawMessage) (int, error) {
	var asNumber int
	if err := json.Unmarshal(raw, &asNumber); err == nil {
		return asNumber, nil
	}
	var asString string
	if err := json.Unmarshal(raw, &asString); err != nil {
		return 0, fmt.Errorf("%w: alg is neither number nor string: %s", ErrUnsupportedAlgorithm, raw)
	}
	n, err := strconv.Atoi(asString)
	if err != nil {
		return 0, fmt.Errorf("%w: alg %q is not numeric", ErrUnsupportedAlgorithm, asString)
	}
	return n, nil
}

func (a *Authenticator) post(ctx context.Context, method string, params, out any) error {
	body, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	})
	if err != nil {
		return fmt.Errorf("gogl: marshal %s: %w", method, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gogl: build %s request: %w", method, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gogl: %s: %w", method, err)
	}
	defer resp.Body.Close()

	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				Wait int `json:"wait"`
			} `json:"data"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("gogl: decode %s response: %w", method, err)
	}
	if envelope.Error != nil {
		// The lockout deserves its own sentinel: retrying makes it worse, and the
		// generic message reads as a wrong password, which sends you looking in
		// entirely the wrong place.
		if envelope.Error.Code == CodeLoginRateLimited {
			if wait := envelope.Error.Data.Wait; wait > 0 {
				return fmt.Errorf("%w: wait %d seconds (%s) before trying again",
					ErrLoginRateLimited, wait, (time.Duration(wait) * time.Second).String())
			}
			return fmt.Errorf("%w: %s", ErrLoginRateLimited, envelope.Error.Message)
		}

		// An access denial on challenge is also rate limiting, reported under a
		// different code.
		//
		// challenge carries only a username: no password, no digest, no session. A
		// denial there cannot mean a wrong password, so reporting it as "Access denied"
		// sends the operator to check a credential that was never sent.
		//
		// OBSERVED 2026-07-30: a script making ~70 separate gogl invocations, each
		// performing a full login, drew -32000 from challenge on every call after the
		// first few, with no wait value and no -32003.
		if method == "challenge" && envelope.Error.Code == CodeAccessDenied {
			return fmt.Errorf("%w: the challenge call was denied, and it carries no "+
				"password -- the router is rate limiting after repeated logins. Wait a "+
				"few minutes. Note that each gogl invocation logs in again, so a script "+
				"making many separate calls can trip this: %s",
				ErrLoginRateLimited, envelope.Error.Message)
		}
		return fmt.Errorf("gogl: %s failed: %s (code %d)", method, envelope.Error.Message, envelope.Error.Code)
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(envelope.Result, out); err != nil {
		return fmt.Errorf("gogl: decode %s result: %w", method, err)
	}
	return nil
}
