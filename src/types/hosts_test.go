package types

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

// The boilerplate a factory GL-SFT1200 ships in its host file, captured verbatim.
// Preserving this is the whole reason gogl works inside a delimited block: it is
// the router's own loopback and IPv6 resolution.
const factoryHostFile = "127.0.0.1 localhost\n\n::1     localhost ip6-localhost ip6-loopback\n" +
	"ff02::1 ip6-allnodes\nff02::2 ip6-allrouters\n"

func TestParseHostFileFactoryState(t *testing.T) {
	f := ParseHostFile(factoryHostFile)

	if f.Domain != "" {
		t.Errorf("Domain = %q, want empty on a factory router", f.Domain)
	}
	if len(f.Entries) != 0 {
		t.Errorf("Entries = %v, want none", f.Entries)
	}
	// Every byte must survive, or a write clobbers the router's own resolution.
	if f.Before != factoryHostFile {
		t.Errorf("Before = %q, want the file verbatim", f.Before)
	}
}

func TestRenderPreservesUnmanagedContent(t *testing.T) {
	f := ParseHostFile(factoryHostFile)
	f.Domain = "lab.example"
	if err := f.Set("nas", "192.168.8.13"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	out := f.Render()
	for _, line := range []string{
		"127.0.0.1 localhost",
		"::1     localhost ip6-localhost ip6-loopback",
		"ff02::1 ip6-allnodes",
		"ff02::2 ip6-allrouters",
	} {
		if !strings.Contains(out, line) {
			t.Errorf("render dropped unmanaged line %q:\n%s", line, out)
		}
	}
	if !strings.Contains(out, BeginMarker) || !strings.Contains(out, EndMarker) {
		t.Errorf("render omitted the block markers:\n%s", out)
	}
}

// The trailing newline matters: without it the first managed line concatenates
// onto the last unmanaged one, which silently turns ff02::2 into the answer for
// your hostname. That happened for real during discovery.
func TestRenderSeparatesBlockFromPrecedingContent(t *testing.T) {
	f := ParseHostFile("ff02::2 ip6-allrouters") // no trailing newline
	f.Domain = "lab.example"
	if err := f.Set("nas", "192.168.8.13"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	out := f.Render()
	if strings.Contains(out, "ip6-allrouters"+BeginMarker) ||
		strings.Contains(out, "ip6-allrouters#") {
		t.Errorf("block ran onto the previous line:\n%s", out)
	}
	if !strings.Contains(out, "ip6-allrouters\n") {
		t.Errorf("missing separator after unmanaged content:\n%s", out)
	}
}

func TestHostFileRoundTrip(t *testing.T) {
	f := ParseHostFile(factoryHostFile)
	f.Domain = "lab.example"
	for name, ip := range map[string]string{"nas": "192.168.8.13", "printer": "192.168.8.14"} {
		if err := f.Set(name, ip); err != nil {
			t.Fatalf("Set(%s): %v", name, err)
		}
	}

	again := ParseHostFile(f.Render())
	if again.Domain != "lab.example" {
		t.Errorf("Domain = %q after round trip", again.Domain)
	}
	if len(again.Entries) != 2 {
		t.Fatalf("Entries = %d after round trip, want 2", len(again.Entries))
	}
	if ip, ok := again.Lookup("nas"); !ok || ip != "192.168.8.13" {
		t.Errorf("Lookup(nas) = %q, %v", ip, ok)
	}
	if again.Before != factoryHostFile {
		t.Errorf("unmanaged content changed across a round trip")
	}
}

// Writing twice must not nest or duplicate the block.
func TestHostFileIsIdempotentAcrossWrites(t *testing.T) {
	f := ParseHostFile(factoryHostFile)
	f.Domain = "lab.example"
	if err := f.Set("nas", "192.168.8.13"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	content := f.Render()
	for i := 0; i < 3; i++ {
		content = ParseHostFile(content).Render()
	}
	if got := strings.Count(content, BeginMarker); got != 1 {
		t.Errorf("BeginMarker appears %d times after repeated writes, want 1:\n%s", got, content)
	}
	if got := strings.Count(content, "192.168.8.13"); got != 1 {
		t.Errorf("entry appears %d times, want 1", got)
	}
}

func TestHostFileSetBothBareAndFQDN(t *testing.T) {
	f := &HostFile{Domain: "lab.example"}
	if err := f.Set("nas", "192.168.8.13"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	if len(f.Entries) != 1 {
		t.Fatalf("Entries = %d, want 1", len(f.Entries))
	}
	names := f.Entries[0].Names
	if len(names) != 2 || names[0] != "nas" || names[1] != "nas.lab.example" {
		t.Errorf("Names = %v, want [nas nas.lab.example]", names)
	}
}

// An already-qualified name is not double-suffixed.
func TestHostFileFQDNDoesNotDoubleQualify(t *testing.T) {
	f := &HostFile{Domain: "lab.example"}
	if got := f.FQDN("nas.other.example"); got != "nas.other.example" {
		t.Errorf("FQDN = %q, want it unchanged", got)
	}
	if got := f.FQDN("nas"); got != "nas.lab.example" {
		t.Errorf("FQDN = %q, want nas.lab.example", got)
	}
}

func TestHostFileFQDNWithoutDomain(t *testing.T) {
	f := &HostFile{}
	if got := f.FQDN("nas"); got != "nas" {
		t.Errorf("FQDN = %q, want nas when no domain is set", got)
	}
}

// Re-setting a name replaces it rather than leaving two answers.
func TestHostFileSetReplaces(t *testing.T) {
	f := &HostFile{Domain: "lab.example"}
	for _, ip := range []string{"192.168.8.13", "192.168.8.20"} {
		if err := f.Set("nas", ip); err != nil {
			t.Fatalf("Set: %v", err)
		}
	}
	if len(f.Entries) != 1 {
		t.Fatalf("Entries = %d, want 1 after replacing", len(f.Entries))
	}
	if ip, _ := f.Lookup("nas"); ip != "192.168.8.20" {
		t.Errorf("Lookup = %q, want the replacement", ip)
	}
}

func TestHostFileSetValidates(t *testing.T) {
	f := &HostFile{Domain: "lab.example"}
	if err := f.Set("bad_name", "192.168.8.13"); !errors.Is(err, ErrInvalidName) {
		t.Errorf("error = %v, want ErrInvalidName", err)
	}
	if err := f.Set("nas", "999.1.1.1"); !errors.Is(err, ErrInvalidIP) {
		t.Errorf("error = %v, want ErrInvalidIP", err)
	}
	if err := f.Set("nas", "fe80::1"); !errors.Is(err, ErrInvalidIP) {
		t.Errorf("error = %v, want ErrInvalidIP for IPv6", err)
	}
}

func TestHostFileRemove(t *testing.T) {
	f := &HostFile{Domain: "lab.example"}
	if err := f.Set("nas", "192.168.8.13"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	if !f.Remove("nas") {
		t.Error("Remove reported nothing removed")
	}
	if len(f.Entries) != 0 {
		t.Errorf("Entries = %v after Remove", f.Entries)
	}
	if f.Remove("nas") {
		t.Error("Remove reported a removal on an absent name")
	}
}

// Removing by either the bare name or the FQDN must work.
func TestHostFileRemoveByFQDN(t *testing.T) {
	f := &HostFile{Domain: "lab.example"}
	if err := f.Set("nas", "192.168.8.13"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if !f.Remove("nas.lab.example") {
		t.Error("Remove by FQDN reported nothing removed")
	}
	if len(f.Entries) != 0 {
		t.Errorf("Entries = %v", f.Entries)
	}
}

// An unmanaged entry sharing an address must keep its other names.
func TestHostFileRemoveKeepsOtherNamesOnSameEntry(t *testing.T) {
	f := ParseHostFile(BeginMarker + " domain lab.example\n" +
		"192.168.8.13 nas nas.lab.example fileserver\n" + EndMarker + "\n")

	if !f.Remove("nas") {
		t.Fatal("Remove reported nothing removed")
	}
	if len(f.Entries) != 1 {
		t.Fatalf("Entries = %v, want the entry kept for its remaining name", f.Entries)
	}
	if got := f.Entries[0].Names; len(got) != 1 || got[0] != "fileserver" {
		t.Errorf("Names = %v, want [fileserver]", got)
	}
}

func TestHostFileClear(t *testing.T) {
	f := ParseHostFile(factoryHostFile)
	f.Domain = "lab.example"
	if err := f.Set("nas", "192.168.8.13"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	f.Clear()
	if len(f.Entries) != 0 {
		t.Errorf("Entries = %v after Clear", f.Entries)
	}
	// Clear drops entries, not the domain or the surrounding file.
	if f.Domain != "lab.example" {
		t.Errorf("Clear discarded the domain")
	}
	if !strings.Contains(f.Render(), "127.0.0.1 localhost") {
		t.Error("Clear discarded unmanaged content")
	}
}

func TestParseHostFileComments(t *testing.T) {
	f := ParseHostFile(BeginMarker + " domain lab.example\n" +
		"# a comment inside the block\n" +
		"192.168.8.13 nas   # trailing comment\n" +
		"\n" +
		EndMarker + "\n")

	if len(f.Entries) != 1 {
		t.Fatalf("Entries = %v, want one", f.Entries)
	}
	if got := f.Entries[0].Names; len(got) != 1 || got[0] != "nas" {
		t.Errorf("Names = %v; a trailing comment leaked in", got)
	}
}

// A block with no closing marker is adopted rather than left to nest.
func TestParseHostFileUnterminatedBlock(t *testing.T) {
	f := ParseHostFile("127.0.0.1 localhost\n" + BeginMarker + " domain lab.example\n" +
		"192.168.8.13 nas\n")

	if f.Domain != "lab.example" {
		t.Errorf("Domain = %q", f.Domain)
	}
	if len(f.Entries) != 1 {
		t.Errorf("Entries = %v, want the orphaned entry adopted", f.Entries)
	}
	if got := strings.Count(f.Render(), BeginMarker); got != 1 {
		t.Errorf("BeginMarker appears %d times after re-render, want 1", got)
	}
}

func TestParseHostFilePreservesContentAfterBlock(t *testing.T) {
	const trailing = "10.0.0.1 something-else\n"
	f := ParseHostFile(BeginMarker + " domain lab.example\n" +
		"192.168.8.13 nas\n" + EndMarker + "\n" + trailing)

	if f.After != trailing {
		t.Errorf("After = %q, want %q", f.After, trailing)
	}
	if !strings.Contains(f.Render(), trailing) {
		t.Error("render dropped content following the block")
	}
}

func TestParseDomainVariants(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{BeginMarker + " domain lab.example", "lab.example"},
		{BeginMarker + " domain lan", "lan"},
		{BeginMarker, ""},
		{BeginMarker + " domain", ""},
	}
	for _, tt := range tests {
		if got := parseDomain(tt.line); got != tt.want {
			t.Errorf("parseDomain(%q) = %q, want %q", tt.line, got, tt.want)
		}
	}
}

// Lookup must accept either spelling, for the same reason Remove must.
func TestHostFileLookupByEitherForm(t *testing.T) {
	f := &HostFile{Domain: "lab.example"}
	if err := f.Set("nas", "192.168.8.13"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	for _, name := range []string{"nas", "nas.lab.example", "NAS"} {
		ip, ok := f.Lookup(name)
		if !ok || ip != "192.168.8.13" {
			t.Errorf("Lookup(%q) = %q, %v", name, ip, ok)
		}
	}
	if _, ok := f.Lookup("other"); ok {
		t.Error("Lookup found a name that was never set")
	}
}

// Setting by the FQDN and then by the bare name is one host, not two.
func TestHostFileSetByEitherFormIsOneEntry(t *testing.T) {
	f := &HostFile{Domain: "lab.example"}
	if err := f.Set("nas.lab.example", "192.168.8.13"); err != nil {
		t.Fatalf("Set FQDN: %v", err)
	}
	if err := f.Set("nas", "192.168.8.20"); err != nil {
		t.Fatalf("Set bare: %v", err)
	}
	if len(f.Entries) != 1 {
		t.Errorf("Entries = %v, want one host", f.Entries)
	}
	if ip, _ := f.Lookup("nas"); ip != "192.168.8.20" {
		t.Errorf("Lookup = %q, want the replacement", ip)
	}
}

// The firmware rejects "(", ")" and "=" anywhere in the host file, with -32602 and
// no indication which character offended it. These tests exist because gogl shipped
// a marker line containing parentheses: the mock accepted it, every test passed, and
// the device refused every write to the host file.

func TestValidateContentRejectsWhatTheFirmwareRejects(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"open paren", "127.0.0.1 localhost\n# a (comment\n"},
		{"close paren", "# a comment)\n"},
		{"equals", "# domain=example.test\n"},
		{"the old marker format", BeginMarker + " (domain: lab.example)\n" + EndMarker + "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContent(tt.content)
			if !errors.Is(err, ErrUnwritableContent) {
				t.Fatalf("error = %v, want ErrUnwritableContent", err)
			}
			// The firmware will not say which character it disliked, so gogl must.
			if !strings.Contains(err.Error(), "line") {
				t.Errorf("error does not locate the problem: %v", err)
			}
		})
	}
}

func TestValidateContentAcceptsWhatTheFirmwareAccepts(t *testing.T) {
	// Every one of these was confirmed accepted by a GL-SFT1200 on 4.3.28.
	accepted := []string{
		"127.0.0.1 localhost\n\n::1     localhost ip6-localhost ip6-loopback\n",
		"# a plain comment\n",
		"# a comment with a colon: here\n",
		"# a comment with a dot example.test\n",
		"# a deliberately long comment with no punctuation at all just plain words here\n",
		BeginMarker + " domain example.test\n192.168.2.253 p p.example.test\n" + EndMarker + "\n",
	}
	for _, content := range accepted {
		if err := ValidateContent(content); err != nil {
			t.Errorf("ValidateContent(%q) = %v, want nil", content, err)
		}
	}
}

// The guarantee that matters: nothing Render produces can be unwritable. Validating
// in Put is a safety net; this is the property that makes the net unnecessary.
func TestRenderNeverProducesUnwritableContent(t *testing.T) {
	domains := []string{"", "lan", "lab.example", "herlein.me", "a-b.c-d.test"}
	names := []string{"nas", "pi-2", "a.b.c", "host.other.test"}

	for _, domain := range domains {
		f := ParseHostFile("127.0.0.1 localhost\n")
		f.Domain = domain
		for i, n := range names {
			// Set requires a domain to qualify against; append directly so the
			// no-domain case is covered too.
			f.Entries = append(f.Entries, HostEntry{
				IP:    "192.168.8." + strconv.Itoa(10+i),
				Names: []string{n, f.FQDN(n)},
			})
		}
		if err := ValidateContent(f.Render()); err != nil {
			t.Errorf("Render with domain %q produced unwritable content: %v\n%s",
				domain, err, f.Render())
		}
	}
}

// An empty block says nothing, and writing one is what made clearing an empty
// router fail.
func TestRenderOmitsAnEmptyBlock(t *testing.T) {
	f := ParseHostFile("127.0.0.1 localhost\n")

	got := f.Render()
	if strings.Contains(got, BeginMarker) {
		t.Errorf("Render emitted a block for an empty, domainless file:\n%s", got)
	}
	if got != "127.0.0.1 localhost\n" {
		t.Errorf("Render = %q, want the input unchanged", got)
	}
}

// A domain alone justifies a block: it is the thing being persisted.
func TestRenderKeepsABlockForADomainWithNoEntries(t *testing.T) {
	f := ParseHostFile("127.0.0.1 localhost\n")
	f.Domain = "lab.example"

	got := f.Render()
	if !strings.Contains(got, BeginMarker+" domain lab.example") {
		t.Errorf("domain not persisted:\n%s", got)
	}
	if ParseHostFile(got).Domain != "lab.example" {
		t.Errorf("domain did not survive a round trip:\n%s", got)
	}
}
