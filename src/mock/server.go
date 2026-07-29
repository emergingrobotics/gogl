// Package mock provides an in-process JSON-RPC server that behaves like a
// GL.iNet router, including its authentication quirks.
//
// Tests run against it rather than hardware, which is why gogl forbids SSH and
// UCI: everything the module does must be reachable this way. The mock
// deliberately reproduces the protocol's hostile behaviors -- nonce expiry,
// session expiry, single-use nonces -- because a mock that only serves
// happy-path payloads cannot catch the bugs a real router would.
package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/emergingrobotics/gogl/src/auth"
	"github.com/emergingrobotics/gogl/src/types"
)

// Defaults mirror the real device. Tests shorten the TTLs to drive expiry
// deterministically rather than by sleeping for tens of seconds.
const (
	DefaultUsername   = "root"
	DefaultNonceTTL   = 1 * time.Second
	DefaultSessionTTL = 35 * time.Second

	// defaultSalt is stable across challenges, as on the real device, so that a
	// client can cache the expensive cipher.
	defaultSalt = "abcdefgh"

	// FactoryHostFile is the host-file content a GL-SFT1200 ships, captured
	// verbatim. Tests that write host entries must leave this intact: it is the
	// router's own loopback and IPv6 resolution.
	FactoryHostFile = "127.0.0.1 localhost\n\n::1     localhost ip6-localhost ip6-loopback\n" +
		"ff02::1 ip6-allnodes\nff02::2 ip6-allrouters\n"

	// JSON-RPC error codes the mock emits.
	CodeAccessDenied     = -32000
	CodeNotFound         = -32001
	CodeLoginRateLimited = -32003
	CodeBadRequest       = -32600
)

// Options configures a Server.
type Options struct {
	Username string
	Password string

	// Alg selects the crypt algorithm advertised by challenge, so tests can
	// exercise the MD5, SHA-256 and SHA-512 paths. Defaults to SHA-512.
	Alg int

	// AlgAsString makes challenge emit alg as a JSON string rather than a
	// number, reproducing observed firmware variation.
	AlgAsString bool

	// HashMethod is the login digest the mock advertises and requires. Empty
	// means the field is omitted from the challenge entirely, reproducing
	// firmware that predates it and expects MD5.
	HashMethod string

	// MaxLoginFailures reproduces the firmware's brute-force lockout. Zero
	// disables it. Once exceeded, challenge and login both return
	// CodeLoginRateLimited with a wait, as real hardware does.
	MaxLoginFailures int

	// LockoutWait is the seconds reported in the rate-limit error's data.wait.
	LockoutWait int

	NonceTTL   time.Duration
	SessionTTL time.Duration
}

// RPCFault is an injected error, returned once for the next matching call.
type RPCFault struct {
	Code    int
	Message string
}

type nonceRecord struct {
	issued time.Time
	used   bool
}

// Server is an httptest-backed /rpc endpoint.
type Server struct {
	t    *testing.T
	http *httptest.Server
	opts Options

	mu       sync.Mutex
	nonces   map[string]*nonceRecord
	nonceSeq int
	sid      string
	sidSeen  time.Time

	loginCount    int
	aliveCount    int
	loginFailures int

	// state holds per-group fixture payloads, keyed "group.method".
	state map[string]json.RawMessage
	// failures holds one-shot injected errors, keyed "group.method".
	failures map[string]*RPCFault
}

// NewServer starts a mock router and registers cleanup with t.
func NewServer(t *testing.T, opts Options) *Server {
	t.Helper()

	if opts.Username == "" {
		opts.Username = DefaultUsername
	}
	if opts.Alg == 0 {
		opts.Alg = auth.AlgSHA512
	}
	if opts.NonceTTL == 0 {
		opts.NonceTTL = DefaultNonceTTL
	}
	if opts.SessionTTL == 0 {
		opts.SessionTTL = DefaultSessionTTL
	}

	s := &Server{
		t:        t,
		opts:     opts,
		nonces:   make(map[string]*nonceRecord),
		state:    make(map[string]json.RawMessage),
		failures: make(map[string]*RPCFault),
	}
	s.http = httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(s.http.Close)

	// Every real router ships a host file, so the mock does too. It carries no gogl
	// block, which is the factory state: no domain configured.
	s.SetHostFile(FactoryHostFile)
	return s
}

// URL returns the JSON-RPC endpoint.
func (s *Server) URL() string { return s.http.URL + "/rpc" }

// Username returns the account the mock accepts.
func (s *Server) Username() string { return s.opts.Username }

// Password returns the credential the mock accepts.
func (s *Server) Password() string { return s.opts.Password }

// Close stops the server. Registered with t.Cleanup already; call it only to
// stop early.
func (s *Server) Close() { s.http.Close() }

// LoginCount reports how many successful logins the mock has served. Tests use
// it to prove that a burst of concurrent calls produced exactly one login.
func (s *Server) LoginCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loginCount
}

// AliveCount reports how many alive calls the mock has served, so a test can
// prove the keepalive stopped after Close.
func (s *Server) AliveCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.aliveCount
}

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, nil, CodeBadRequest, "malformed request")
		return
	}

	switch req.Method {
	case "challenge":
		s.handleChallenge(w, &req)
	case "login":
		s.handleLogin(w, &req)
	case "alive":
		s.handleAlive(w, &req)
	case "call":
		s.handleCall(w, &req)
	default:
		s.writeError(w, req.ID, CodeBadRequest, "unknown method "+req.Method)
	}
}

func (s *Server) handleChallenge(w http.ResponseWriter, req *request) {
	var params struct {
		Username string `json:"username"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.writeError(w, req.ID, CodeBadRequest, "bad challenge params")
		return
	}
	if params.Username != s.opts.Username {
		s.writeError(w, req.ID, CodeAccessDenied, "unknown user")
		return
	}
	if s.rateLimited() {
		s.writeRateLimit(w, req.ID)
		return
	}

	s.mu.Lock()
	s.nonceSeq++
	nonce := fmt.Sprintf("nonce-%d", s.nonceSeq)
	s.nonces[nonce] = &nonceRecord{issued: time.Now()}
	s.mu.Unlock()

	// The salt is stable across challenges so a client can cache the expensive
	// cipher; only the nonce rotates.
	result := map[string]any{"salt": defaultSalt, "nonce": nonce}
	// Omitted entirely when empty: that is how older firmware behaves, and a
	// client must read its absence as MD5.
	if s.opts.HashMethod != "" {
		result["hash-method"] = s.opts.HashMethod
	}
	if s.opts.AlgAsString {
		result["alg"] = strconv.Itoa(s.opts.Alg)
	} else {
		result["alg"] = s.opts.Alg
	}
	s.writeResult(w, req.ID, result)
}

func (s *Server) handleLogin(w http.ResponseWriter, req *request) {
	var params struct {
		Username string `json:"username"`
		Hash     string `json:"hash"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.writeError(w, req.ID, CodeBadRequest, "bad login params")
		return
	}

	cipher, err := auth.Cipher(s.opts.Password, defaultSalt, s.opts.Alg)
	if err != nil {
		// Errorf rather than Fatalf: this runs on the server's goroutine, where
		// Fatalf's runtime.Goexit would be a bug.
		s.t.Errorf("mock: cipher: %v", err)
		s.writeError(w, req.ID, CodeBadRequest, "mock cipher failure")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the nonce this digest was computed against. The real device tracks
	// the outstanding nonce; reproducing that is what makes reuse and expiry
	// observable in tests.
	var matched string
	for nonce := range s.nonces {
		expected, err := auth.LoginHash(params.Username, cipher, nonce, s.opts.HashMethod)
		if err != nil {
			s.t.Errorf("mock: login hash: %v", err)
			s.writeErrorLocked(w, req.ID, CodeBadRequest, "mock hash failure")
			return
		}
		if expected == params.Hash {
			matched = nonce
			break
		}
	}
	if matched == "" {
		s.loginFailures++
		if s.opts.MaxLoginFailures > 0 && s.loginFailures > s.opts.MaxLoginFailures {
			s.writeRateLimitLocked(w, req.ID)
			return
		}
		s.writeErrorLocked(w, req.ID, CodeAccessDenied, "invalid credentials")
		return
	}

	rec := s.nonces[matched]
	if rec.used {
		s.writeErrorLocked(w, req.ID, CodeAccessDenied, "nonce already used")
		return
	}
	if time.Since(rec.issued) > s.opts.NonceTTL {
		s.writeErrorLocked(w, req.ID, CodeAccessDenied, "nonce expired")
		return
	}
	rec.used = true

	s.nonceSeq++
	s.sid = fmt.Sprintf("sid-%d", s.nonceSeq)
	s.sidSeen = time.Now()
	s.loginCount++
	s.writeResultLocked(w, req.ID, map[string]any{"sid": s.sid, "username": params.Username})
}

func (s *Server) handleAlive(w http.ResponseWriter, req *request) {
	sid, ok := s.extractSID(req.Params)
	if !ok {
		s.writeError(w, req.ID, CodeBadRequest, "bad alive params")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.aliveCount++

	if !s.sessionValidLocked(sid) {
		s.writeErrorLocked(w, req.ID, CodeAccessDenied, "Access denied")
		return
	}
	s.sidSeen = time.Now()
	s.writeResultLocked(w, req.ID, map[string]any{"alive": true})
}

// extractSID accepts the sid as a bare params object or as the first element of
// a params array, both of which have been observed for alive.
func (s *Server) extractSID(params json.RawMessage) (string, bool) {
	var object struct {
		SID string `json:"sid"`
	}
	if err := json.Unmarshal(params, &object); err == nil {
		return object.SID, true
	}

	var array []json.RawMessage
	if err := json.Unmarshal(params, &array); err != nil || len(array) == 0 {
		return "", false
	}
	var sid string
	if err := json.Unmarshal(array[0], &sid); err != nil {
		return "", false
	}
	return sid, true
}

// sessionValidLocked reports whether sid is the current session and has not
// idled out. Callers must hold s.mu.
func (s *Server) sessionValidLocked(sid string) bool {
	if sid == "" || sid != s.sid {
		return false
	}
	return time.Since(s.sidSeen) <= s.opts.SessionTTL
}

// rateLimited reports whether the mock's lockout has tripped.
func (s *Server) rateLimited() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.opts.MaxLoginFailures > 0 && s.loginFailures > s.opts.MaxLoginFailures
}

func (s *Server) writeRateLimit(w http.ResponseWriter, id json.RawMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writeRateLimitLocked(w, id)
}

func (s *Server) writeRateLimitLocked(w http.ResponseWriter, id json.RawMessage) {
	wait := s.opts.LockoutWait
	if wait == 0 {
		wait = 600
	}
	s.encode(w, map[string]any{
		"jsonrpc": "2.0", "id": rawOrNull(id),
		"error": map[string]any{
			"code":    CodeLoginRateLimited,
			"message": "Login fail number over limit",
			"data":    map[string]any{"wait": wait},
		},
	})
}

func (s *Server) writeResult(w http.ResponseWriter, id json.RawMessage, result any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writeResultLocked(w, id, result)
}

func (s *Server) writeResultLocked(w http.ResponseWriter, id json.RawMessage, result any) {
	s.encode(w, map[string]any{"jsonrpc": "2.0", "id": rawOrNull(id), "result": result})
}

func (s *Server) writeError(w http.ResponseWriter, id json.RawMessage, code int, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writeErrorLocked(w, id, code, message)
}

func (s *Server) writeErrorLocked(w http.ResponseWriter, id json.RawMessage, code int, message string) {
	s.encode(w, map[string]any{
		"jsonrpc": "2.0", "id": rawOrNull(id),
		"error": map[string]any{"code": code, "message": message},
	})
}

func (s *Server) encode(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		s.t.Errorf("mock: encode response: %v", err)
	}
}

func rawOrNull(id json.RawMessage) any {
	if len(id) == 0 {
		return nil
	}
	return id
}

// HostFileWith builds host-file content: the factory boilerplate, plus gogl's
// managed block carrying domain and entries.
//
// It exists so tests never hand-assemble the marker line. They used to, and when
// hardware forced the marker format to change, every one of those literals became a
// silent lie: the tests still passed while asserting against a format no router
// would accept. Building through types.HostFile means the format has exactly one
// definition.
//
// Each entry is "IP name [name...]", as it appears in the file.
func HostFileWith(domain string, entries ...string) string {
	f := types.ParseHostFile(FactoryHostFile)
	f.Domain = domain
	for _, e := range entries {
		fields := strings.Fields(e)
		if len(fields) < 2 {
			continue
		}
		f.Entries = append(f.Entries, types.HostEntry{IP: fields[0], Names: fields[1:]})
	}
	return f.Render()
}
