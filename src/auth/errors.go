package auth

import "errors"

// Sentinels live here rather than in the root package because auth is a leaf.
// The root package re-exports them.
var (
	// ErrUnsupportedAlgorithm means the challenge named a crypt algorithm this
	// package does not implement. Never falls back to a weaker algorithm: a
	// silent downgrade is worse than a failed login.
	ErrUnsupportedAlgorithm = errors.New("gogl: unsupported crypt algorithm")

	// ErrUnsupportedHashMethod means the challenge named a login digest this
	// package does not implement. Never falls back to another hash: a wrong digest
	// is indistinguishable from a wrong password.
	ErrUnsupportedHashMethod = errors.New("gogl: unsupported login hash method")

	// ErrNonceExpired means the login nonce died before it was used. Retriable.
	ErrNonceExpired = errors.New("gogl: challenge nonce expired")

	// ErrUnauthorized means the router rejected the credentials.
	ErrUnauthorized = errors.New("gogl: unauthorized")

	// ErrLoginRateLimited means the firmware's brute-force protection has locked
	// the account. Observed on an SFT1200 after roughly a dozen failed logins; the
	// lockout is around ten minutes and no amount of retrying shortens it.
	ErrLoginRateLimited = errors.New("gogl: too many failed logins, router is rate limiting")
)
