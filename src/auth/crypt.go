// Package auth implements GL.iNet firmware 4.x challenge/response
// authentication. The password is never transmitted; only a digest derived from
// it, salted by the router and bound to a short-lived nonce.
package auth

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/GehirnInc/crypt"
	// Registering the three algorithms the firmware can name. Blank imports are
	// how this library exposes them.
	_ "github.com/GehirnInc/crypt/md5_crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
)

// Algorithm identifiers as returned in the challenge response's alg field. These
// select the unix crypt(3) variant used to derive the password cipher.
const (
	AlgMD5    = 1
	AlgSHA256 = 5
	AlgSHA512 = 6
)

// CodeLoginRateLimited is the JSON-RPC code the firmware returns once its
// brute-force protection trips. Observed on a GL-SFT1200 running 4.3.28, with a
// data.wait field giving the remaining lockout in seconds.
const CodeLoginRateLimited = -32003

// CodeAccessDenied is what the firmware returns for a refused credential -- and also
// for a challenge call refused by brute-force protection, which is why the two cases
// have to be told apart by which method was called.
const CodeAccessDenied = -32000

// Hash methods as returned in the challenge response's "hash-method" field. This
// is a separate choice from alg: it selects the digest applied to the login
// string, not the crypt used on the password.
//
// Firmware 4.3.28 on the GL-SFT1200 reports "sha256". Older firmware omits the
// field entirely and expects MD5, which is what every public client library
// implements -- so the field must be honored, and its absence must mean MD5.
const (
	HashMethodMD5    = "md5"
	HashMethodSHA256 = "sha256"
)

// Cipher derives the unix crypt(3) hash of password under salt, using the
// algorithm the router named. The result is the full crypt string including its
// "$alg$salt$" prefix, which is what the login digest is computed over.
//
// This is deliberately slow for SHA-256 and SHA-512 (5000 rounds). Callers must
// obtain a fresh nonce after calling this, because the nonce from the challenge
// that supplied the salt has very likely expired by the time it returns.
func Cipher(password, salt string, alg int) (string, error) {
	if salt == "" {
		return "", fmt.Errorf("gogl: crypt salt is empty")
	}

	var c crypt.Crypter
	switch alg {
	case AlgMD5:
		c = crypt.MD5.New()
	case AlgSHA256:
		c = crypt.SHA256.New()
	case AlgSHA512:
		c = crypt.SHA512.New()
	default:
		return "", fmt.Errorf("%w: %d", ErrUnsupportedAlgorithm, alg)
	}

	setting := fmt.Sprintf("$%d$%s", alg, salt)
	hashed, err := c.Generate([]byte(password), []byte(setting))
	if err != nil {
		return "", fmt.Errorf("gogl: crypt: %w", err)
	}
	return hashed, nil
}

// LoginHash computes the digest sent as the login method's hash parameter:
// H(username:cipher:nonce), where H is the hash the router named.
//
// The digest is not the password hash: it binds the already-strong crypt output
// to a single-use nonce. An empty hashMethod means MD5, which is what firmware
// that predates the "hash-method" field expects.
func LoginHash(username, cipher, nonce, hashMethod string) (string, error) {
	payload := []byte(username + ":" + cipher + ":" + nonce)

	switch hashMethod {
	case "", HashMethodMD5:
		sum := md5.Sum(payload)
		return hex.EncodeToString(sum[:]), nil
	case HashMethodSHA256:
		sum := sha256.Sum256(payload)
		return hex.EncodeToString(sum[:]), nil
	default:
		// Never fall back to a different hash: a wrong digest is indistinguishable
		// from a wrong password, which is a miserable failure to debug.
		return "", fmt.Errorf("%w: hash-method %q", ErrUnsupportedHashMethod, hashMethod)
	}
}
