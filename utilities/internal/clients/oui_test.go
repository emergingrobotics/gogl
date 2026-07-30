package clients

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const sampleOUI = `
OUI/MA-L                                                    Organization
company_id                                                  Organization
                                                            Address

AC-DE-48   (hex)		Private
ACDE48     (base 16)		Private

00-1B-63   (hex)		Apple, Inc.
001B63     (base 16)		Apple, Inc.
				1 Infinite Loop
				Cupertino CA 95014
				US

D4-9A-20   (hex)		Dell Inc.
D49A20     (base 16)		Dell Inc.
`

func TestParseOUI(t *testing.T) {
	db, err := ParseOUI(strings.NewReader(sampleOUI))
	if err != nil {
		t.Fatalf("ParseOUI error: %v", err)
	}
	if len(db) != 3 {
		t.Fatalf("parsed %d entries, want 3: %v", len(db), db)
	}
	if got := db["00:1b:63"]; got != "Apple, Inc." {
		t.Errorf("db[00:1b:63] = %q, want %q", got, "Apple, Inc.")
	}
	if got := db["d4:9a:20"]; got != "Dell Inc." {
		t.Errorf("db[d4:9a:20] = %q, want %q", got, "Dell Inc.")
	}
}

func TestParseOUIEmpty(t *testing.T) {
	db, err := ParseOUI(strings.NewReader("nothing useful here\n"))
	if err != nil {
		t.Fatalf("ParseOUI error: %v", err)
	}
	if len(db) != 0 {
		t.Errorf("parsed %d entries from junk, want 0", len(db))
	}
}

func TestLookup(t *testing.T) {
	db, err := ParseOUI(strings.NewReader(sampleOUI))
	if err != nil {
		t.Fatalf("ParseOUI error: %v", err)
	}

	tests := []struct {
		mac  string
		want string
	}{
		{"00:1b:63:aa:bb:cc", "Apple, Inc."},
		{"00:1B:63:AA:BB:CC", "Apple, Inc."},
		{"d4:9a:20:11:22:33", "Dell Inc."},
		// fc has the locally-administered bit clear, so an absent entry is
		// genuinely unknown rather than randomized. ff:ff:ff would NOT work here:
		// 0xff has that bit set.
		{"fc:ff:ff:11:22:33", unknownManufacturer},
	}
	for _, tt := range tests {
		if got := db.Lookup(tt.mac); got != tt.want {
			t.Errorf("Lookup(%s) = %q, want %q", tt.mac, got, tt.want)
		}
	}
}

// A locally-administered MAC will never appear in the IEEE database, so
// reporting "unknown" would be misleading: it is a randomized address, and a
// poor choice to reserve an address for.
func TestLookupRandomized(t *testing.T) {
	db := OUIDatabase{}
	for _, mac := range []string{
		"02:11:22:33:44:55",
		"06:aa:bb:cc:dd:ee",
		"0a:00:00:00:00:01",
		"de:ad:be:ef:00:01",
	} {
		if got := db.Lookup(mac); got != randomizedManufacturer {
			t.Errorf("Lookup(%s) = %q, want %q", mac, got, randomizedManufacturer)
		}
	}
}

func TestLookupGloballyAdministeredIsNotRandomized(t *testing.T) {
	db := OUIDatabase{}
	// Second-least-significant bit of the first octet clear means globally
	// administered, so an absent entry is genuinely unknown.
	for _, mac := range []string{"00:1b:63:aa:bb:cc", "d4:9a:20:11:22:33", "fc:00:00:00:00:01"} {
		if got := db.Lookup(mac); got != unknownManufacturer {
			t.Errorf("Lookup(%s) = %q, want %q", mac, got, unknownManufacturer)
		}
	}
}

// A registered prefix wins over the randomized heuristic.
func TestLookupRegisteredBeatsLocalAdminBit(t *testing.T) {
	db := OUIDatabase{"02:11:22": "Weird Corp"}
	if got := db.Lookup("02:11:22:33:44:55"); got != "Weird Corp" {
		t.Errorf("Lookup = %q, want the registered name", got)
	}
}

func TestLookupMalformed(t *testing.T) {
	db := OUIDatabase{}
	for _, mac := range []string{"", "zz", "00:1b", "not:a"} {
		if got := db.Lookup(mac); got != unknownManufacturer {
			t.Errorf("Lookup(%q) = %q, want %q", mac, got, unknownManufacturer)
		}
	}
}

// A fresh cache must be used without a download.
func TestLoadOUIUsesFreshCache(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ouiFilename), []byte(sampleOUI), 0o644); err != nil {
		t.Fatalf("seed cache: %v", err)
	}

	fetched := false
	db, err := LoadOUI(dir, func() (string, error) {
		fetched = true
		return "", errors.New("should not be called")
	})
	if err != nil {
		t.Fatalf("LoadOUI error: %v", err)
	}
	if fetched {
		t.Error("LoadOUI downloaded despite a fresh cache")
	}
	if db.Lookup("00:1b:63:aa:bb:cc") != "Apple, Inc." {
		t.Error("cache was not used")
	}
}

func TestLoadOUIDownloadsWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	db, err := LoadOUI(dir, func() (string, error) { return sampleOUI, nil })
	if err != nil {
		t.Fatalf("LoadOUI error: %v", err)
	}
	if db.Lookup("d4:9a:20:11:22:33") != "Dell Inc." {
		t.Error("downloaded data was not parsed")
	}
	// The download must be cached for next time.
	if _, err := os.Stat(filepath.Join(dir, ouiFilename)); err != nil {
		t.Errorf("download was not cached: %v", err)
	}
}

func TestLoadOUIRefreshesStaleCache(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ouiFilename)
	if err := os.WriteFile(path, []byte("AC-DE-48   (hex)\t\tStale Corp\n"), 0o644); err != nil {
		t.Fatalf("seed cache: %v", err)
	}
	stale := time.Now().Add(-(maxCacheAge + 24*time.Hour))
	if err := os.Chtimes(path, stale, stale); err != nil {
		t.Fatalf("age cache: %v", err)
	}

	db, err := LoadOUI(dir, func() (string, error) { return sampleOUI, nil })
	if err != nil {
		t.Fatalf("LoadOUI error: %v", err)
	}
	if db.Lookup("00:1b:63:aa:bb:cc") != "Apple, Inc." {
		t.Error("stale cache was not refreshed")
	}
}

// Download failure with a stale cache present is a warning, not a failure: a
// slightly out-of-date manufacturer name beats no output at all.
func TestLoadOUIFallsBackToStaleCache(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ouiFilename)
	if err := os.WriteFile(path, []byte(sampleOUI), 0o644); err != nil {
		t.Fatalf("seed cache: %v", err)
	}
	stale := time.Now().Add(-(maxCacheAge + 24*time.Hour))
	if err := os.Chtimes(path, stale, stale); err != nil {
		t.Fatalf("age cache: %v", err)
	}

	db, err := LoadOUI(dir, func() (string, error) { return "", errors.New("no network") })
	if err != nil {
		t.Fatalf("LoadOUI should fall back to a stale cache, got: %v", err)
	}
	if db.Lookup("00:1b:63:aa:bb:cc") != "Apple, Inc." {
		t.Error("stale cache was not used as a fallback")
	}
}

// No cache and no network is a hard failure, matching gofimac. Documented as a
// known limitation rather than worked around.
func TestLoadOUIFailsWithNoCacheAndNoNetwork(t *testing.T) {
	if _, err := LoadOUI(t.TempDir(), func() (string, error) {
		return "", errors.New("no network")
	}); err == nil {
		t.Error("LoadOUI succeeded with no cache and no network, want error")
	}
}

func TestLoadOUICreatesCacheDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "does", "not", "exist")
	if _, err := LoadOUI(dir, func() (string, error) { return sampleOUI, nil }); err != nil {
		t.Fatalf("LoadOUI error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ouiFilename)); err != nil {
		t.Errorf("cache directory was not created: %v", err)
	}
}

// IEEE fronts the OUI file with a bot Filter that answers Go's default
// "Go-http-client/1.1" with HTTP 418. That failure is indistinguishable from a
// broken network at the command line, so the User-Agent is pinned here.
func TestFetchOUISendsUserAgent(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("User-Agent")
		fmt.Fprint(w, sampleOUI)
	}))
	defer srv.Close()

	body, err := fetchOUI(srv.URL)
	if err != nil {
		t.Fatalf("fetchOUI error: %v", err)
	}
	if body == "" {
		t.Error("fetchOUI returned an empty body")
	}
	if got == "" {
		t.Fatal("no User-Agent was sent")
	}
	if strings.HasPrefix(got, "Go-http-client") {
		t.Errorf("User-Agent = %q, which IEEE answers with HTTP 418", got)
	}
	if !strings.Contains(got, "goglmac") {
		t.Errorf("User-Agent = %q, want it to identify this tool", got)
	}
}

func TestFetchOUIReportsHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer srv.Close()

	_, err := fetchOUI(srv.URL)
	if err == nil {
		t.Fatal("fetchOUI succeeded on a 418")
	}
	if !strings.Contains(err.Error(), "418") {
		t.Errorf("error %q does not report the status code", err)
	}
}

// The downloaded body must actually parse, not merely arrive.
func TestFetchOUIRoundTripsThroughParse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, sampleOUI)
	}))
	defer srv.Close()

	body, err := fetchOUI(srv.URL)
	if err != nil {
		t.Fatalf("fetchOUI error: %v", err)
	}
	db, err := ParseOUI(strings.NewReader(body))
	if err != nil {
		t.Fatalf("ParseOUI error: %v", err)
	}
	if db.Lookup("00:1b:63:aa:bb:cc") != "Apple, Inc." {
		t.Error("downloaded body did not parse into a usable database")
	}
}

// The live IEEE file uses CRLF line endings. The LF-only sampleOUI above would
// not have caught a regression here, so this pins the real format explicitly:
// a stray carriage return would end up inside every manufacturer name.
func TestParseOUIHandlesCRLF(t *testing.T) {
	const crlf = "00-1B-63   (hex)\t\tApple, Inc.\r\n" +
		"B8-27-EB   (hex)\t\tRaspberry Pi Foundation\r\n" +
		"00-50-56   (hex)\t\tVMware, Inc.\r\n"

	db, err := ParseOUI(strings.NewReader(crlf))
	if err != nil {
		t.Fatalf("ParseOUI error: %v", err)
	}

	want := map[string]string{
		"00:1b:63:aa:bb:cc": "Apple, Inc.",
		"b8:27:eb:11:22:33": "Raspberry Pi Foundation",
		"00:50:56:11:22:33": "VMware, Inc.",
	}
	for mac, organization := range want {
		got := db.Lookup(mac)
		if got != organization {
			t.Errorf("Lookup(%s) = %q, want %q", mac, got, organization)
		}
		if strings.ContainsAny(got, "\r\n") {
			t.Errorf("Lookup(%s) = %q, which still carries a line ending", mac, got)
		}
	}
}

// Real-world prefixes verified against the live IEEE database on 2026-07-27.
// These are stable registrations, and they double as a sanity check that the
// parser's prefix normalisation matches what the file actually contains.
func TestParseOUIRealWorldPrefixes(t *testing.T) {
	const real = "00-1B-63   (hex)\t\tApple, Inc.\r\n" +
		"D4-9A-20   (hex)\t\tApple, Inc.\r\n" +
		"B8-27-EB   (hex)\t\tRaspberry Pi Foundation\r\n" +
		"DC-A6-32   (hex)\t\tRaspberry Pi Trading Ltd\r\n"

	db, err := ParseOUI(strings.NewReader(real))
	if err != nil {
		t.Fatalf("ParseOUI error: %v", err)
	}
	// A Raspberry Pi is the likeliest thing an emergingrobotics kit pins an
	// address for, so it is worth asserting by name.
	if got := db.Lookup("b8:27:eb:00:11:22"); got != "Raspberry Pi Foundation" {
		t.Errorf("Lookup of a Raspberry Pi MAC = %q", got)
	}
	if got := db.Lookup("dc:a6:32:00:11:22"); got != "Raspberry Pi Trading Ltd" {
		t.Errorf("Lookup of a newer Raspberry Pi MAC = %q", got)
	}
}
