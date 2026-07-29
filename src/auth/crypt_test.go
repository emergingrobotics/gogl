package auth

import (
	"errors"
	"strings"
	"testing"
)

// Vectors generated with:
//
//	openssl passwd -1 -salt abcdefgh testpassword
//	openssl passwd -5 -salt abcdefgh testpassword
//	openssl passwd -6 -salt abcdefgh testpassword
//
// These are authoritative. If an implementation disagrees, the implementation is
// wrong.
const (
	testPassword = "testpassword"
	testSalt     = "abcdefgh"

	wantMD5    = "$1$abcdefgh$H6JMmWFBXCyBkxzBuU/es0"
	wantSHA256 = "$5$abcdefgh$O0RDERJFpTqZJIJKvF.ES67YlwQkXIZRUnti0faDht5"
	wantSHA512 = "$6$abcdefgh$.ofHZDk5EnkwHnbcCRFECyA9NAXafNK89M2N49HOc2iXEMuAVgw2VQrHEjAL6PQe8YtZ8W02Ai/xrAzwN5LIK1"
)

func TestCipher(t *testing.T) {
	tests := []struct {
		name string
		alg  int
		want string
	}{
		{"MD5", AlgMD5, wantMD5},
		{"SHA256", AlgSHA256, wantSHA256},
		{"SHA512", AlgSHA512, wantSHA512},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Cipher(testPassword, testSalt, tt.alg)
			if err != nil {
				t.Fatalf("Cipher error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Cipher(alg=%d) =\n  %q\nwant\n  %q", tt.alg, got, tt.want)
			}
		})
	}
}

// An unrecognized algorithm must fail rather than fall back to a weaker one.
func TestCipherRejectsUnknownAlgorithm(t *testing.T) {
	for _, alg := range []int{0, 2, 3, 4, 7, 99, -1} {
		if _, err := Cipher(testPassword, testSalt, alg); !errors.Is(err, ErrUnsupportedAlgorithm) {
			t.Errorf("Cipher(alg=%d) error = %v, want ErrUnsupportedAlgorithm", alg, err)
		}
	}
}

func TestCipherRejectsEmptySalt(t *testing.T) {
	if _, err := Cipher(testPassword, "", AlgSHA512); err == nil {
		t.Error("Cipher with empty salt succeeded, want error")
	}
}

// A different password must yield a different cipher under the same salt, or the
// hash is not actually a function of the password.
func TestCipherVariesWithPassword(t *testing.T) {
	a, err := Cipher("one", testSalt, AlgSHA512)
	if err != nil {
		t.Fatalf("Cipher: %v", err)
	}
	b, err := Cipher("two", testSalt, AlgSHA512)
	if err != nil {
		t.Fatalf("Cipher: %v", err)
	}
	if a == b {
		t.Error("Cipher produced the same hash for different passwords")
	}
}

// Vectors generated with:
//
//	printf 'root:%s:testnonce' "$SHA512_CIPHER" | md5sum
//	printf 'root:%s:testnonce' "$SHA512_CIPHER" | sha256sum
const (
	wantDigestMD5    = "c31c85d0648225cc107b5dfc0a410060"
	wantDigestSHA256 = "4f16857a07ffa841fd9227429325cdb28bc7c080e506c20835aa92fc15fa3497"
)

func TestLoginHash(t *testing.T) {
	tests := []struct {
		name       string
		hashMethod string
		want       string
	}{
		// An absent hash-method means MD5: that is how firmware predating the
		// field behaves, and it is what every public client library implements.
		{"omitted defaults to MD5", "", wantDigestMD5},
		{"explicit md5", HashMethodMD5, wantDigestMD5},
		{"sha256", HashMethodSHA256, wantDigestSHA256},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoginHash("root", wantSHA512, "testnonce", tt.hashMethod)
			if err != nil {
				t.Fatalf("LoginHash error: %v", err)
			}
			if got != tt.want {
				t.Errorf("LoginHash(%q) = %q, want %q", tt.hashMethod, got, tt.want)
			}
		})
	}
}

// The two methods must not agree, or the selection is not doing anything.
func TestLoginHashMethodsDiffer(t *testing.T) {
	if wantDigestMD5 == wantDigestSHA256 {
		t.Fatal("the two vectors are identical; one of them is wrong")
	}
}

// An unrecognized method must fail rather than silently pick one. A wrong digest
// is indistinguishable from a wrong password, which is a miserable failure to
// debug -- it cost a full debugging session on real hardware to find.
func TestLoginHashRejectsUnknownMethod(t *testing.T) {
	for _, method := range []string{"sha512", "bcrypt", "none", "MD5", "SHA256"} {
		if _, err := LoginHash("root", wantSHA512, "n", method); !errors.Is(err, ErrUnsupportedHashMethod) {
			t.Errorf("LoginHash(%q) error = %v, want ErrUnsupportedHashMethod", method, err)
		}
	}
}

// A different nonce must yield a different digest, or replay protection is
// absent.
func TestLoginHashVariesWithNonce(t *testing.T) {
	for _, method := range []string{HashMethodMD5, HashMethodSHA256} {
		a, err := LoginHash("root", wantSHA512, "nonce-one", method)
		if err != nil {
			t.Fatalf("LoginHash: %v", err)
		}
		b, err := LoginHash("root", wantSHA512, "nonce-two", method)
		if err != nil {
			t.Fatalf("LoginHash: %v", err)
		}
		if a == b {
			t.Errorf("%s produced the same digest for different nonces", method)
		}
	}
}

func TestLoginHashVariesWithUsername(t *testing.T) {
	a, err := LoginHash("root", wantSHA512, "n", HashMethodSHA256)
	if err != nil {
		t.Fatalf("LoginHash: %v", err)
	}
	b, err := LoginHash("admin", wantSHA512, "n", HashMethodSHA256)
	if err != nil {
		t.Fatalf("LoginHash: %v", err)
	}
	if a == b {
		t.Error("LoginHash produced the same digest for different usernames")
	}
}

// Captured from a live GL-SFT1200 on firmware 4.3.28: alg 5 (SHA-256 crypt) with a
// 16-byte salt, and hash-method "sha256". The salt is real; the password is this
// package's test password, so no credential derived from a real one is committed.
func TestRealSFT1200ChallengeShape(t *testing.T) {
	const observedSalt = "j2K6qjQJi8fCtAzO"

	cipher, err := Cipher(testPassword, observedSalt, AlgSHA256)
	if err != nil {
		t.Fatalf("Cipher: %v", err)
	}
	if !strings.HasPrefix(cipher, "$5$"+observedSalt+"$") {
		t.Errorf("cipher = %q, want a $5$ crypt over the observed salt", cipher)
	}

	// The firmware pairs alg 5 with hash-method sha256; both must be honored
	// together, and the result must not equal what the MD5 path would produce.
	viaSHA256, err := LoginHash("root", cipher, "nonce", HashMethodSHA256)
	if err != nil {
		t.Fatalf("LoginHash: %v", err)
	}
	viaMD5, err := LoginHash("root", cipher, "nonce", HashMethodMD5)
	if err != nil {
		t.Fatalf("LoginHash: %v", err)
	}
	if viaSHA256 == viaMD5 {
		t.Error("sha256 and md5 digests agree; the selection is broken")
	}
	if len(viaSHA256) != 64 {
		t.Errorf("sha256 digest is %d hex chars, want 64", len(viaSHA256))
	}
	if len(viaMD5) != 32 {
		t.Errorf("md5 digest is %d hex chars, want 32", len(viaMD5))
	}
}
