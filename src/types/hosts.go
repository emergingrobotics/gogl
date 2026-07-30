package types

import (
	"fmt"
	"net"
	"strings"
)

// The router's host file is a plain hosts(5) file served through dns.get_host and
// dns.set_host. dnsmasq answers from it, so this is where gogl creates real DNS
// records -- a reservation does not, see Reservation.
//
// gogl owns a delimited block inside that file and never touches anything outside
// it. The file also holds the loopback and IPv6 boilerplate the firmware ships,
// and clobbering that would be rude at best.
const (
	// BeginMarker opens gogl's block. The domain rides in this line so it is
	// persisted on the device and can be read back: the firmware exposes no
	// setting for a DNS domain, and /ubus is not routed, so there is nowhere else
	// to put it.
	BeginMarker = "# BEGIN gogl managed hosts"

	// EndMarker closes the block.
	EndMarker = "# END gogl managed hosts"

	// domainKey introduces the domain in the marker line, as a bare word.
	//
	// It was "domain:" wrapped in parentheses until hardware said otherwise:
	// dns.set_host rejects content containing "(", ")" or "=" with -32602 Invalid
	// params, however well-formed the file is otherwise. Colons, dots, hyphens and
	// long lines are all fine. This exact marker shape is the one verified accepted
	// with an entry present.
	domainKey = "domain"
)

// rejectedByFirmware are the characters dns.set_host refuses anywhere in the file.
//
// Isolated one character at a time against a GL-SFT1200 on 4.3.28. The firmware
// gives no indication which character offended it, so gogl checks before writing
// rather than passing -32602 back to the caller.
const rejectedByFirmware = `()=`

// ValidateContent reports whether content is writable through dns.set_host.
//
// Callers do not need this: HostsService.Put runs it. It is exported because a
// caller assembling a file by hand deserves the same error the library gets,
// naming the character and the line rather than the firmware's bare "Invalid
// params".
func ValidateContent(content string) error {
	for lineNumber, line := range strings.Split(content, "\n") {
		if idx := strings.IndexAny(line, rejectedByFirmware); idx >= 0 {
			return fmt.Errorf("%w: line %d contains %q, which dns.set_host rejects: %s",
				ErrUnwritableContent, lineNumber+1, line[idx], line)
		}
	}
	return nil
}

// HostEntry is one line of the managed block: an address and the names it answers.
type HostEntry struct {
	IP    string   `json:"ip"`
	Names []string `json:"names"`
}

// HostFile is a parsed host file, split into the parts gogl owns and the parts it
// must preserve.
type HostFile struct {
	// Before and After are the unmanaged content surrounding gogl's block,
	// preserved verbatim.
	Before string
	After  string

	// Domain is the suffix gogl appends when generating FQDNs. Empty means it has
	// never been configured, which several operations refuse to proceed without.
	Domain string

	// Entries are the managed lines.
	Entries []HostEntry
}

// ParseHostFile splits raw host-file content around gogl's managed block.
//
// A file with no block parses as all-unmanaged with no domain, which is the
// factory state and is how "domain not configured" is detected.
func ParseHostFile(raw string) *HostFile {
	f := &HostFile{}

	begin := strings.Index(raw, BeginMarker)
	if begin < 0 {
		f.Before = raw
		return f
	}

	f.Before = raw[:begin]
	rest := raw[begin:]

	// The marker line carries the domain.
	newline := strings.Index(rest, "\n")
	if newline < 0 {
		// A truncated marker line: treat the whole thing as ours and empty.
		f.Domain = parseDomain(rest)
		return f
	}
	f.Domain = parseDomain(rest[:newline])
	rest = rest[newline+1:]

	end := strings.Index(rest, EndMarker)
	if end < 0 {
		// No closing marker. Everything after the opener is ours; better to adopt
		// it than to leave a half-block that the next write would nest inside.
		f.Entries = parseEntries(rest)
		return f
	}

	f.Entries = parseEntries(rest[:end])

	after := rest[end:]
	if newline := strings.Index(after, "\n"); newline >= 0 {
		f.After = after[newline+1:]
	}
	return f
}

// parseDomain pulls the domain out of a begin-marker line.
//
// The line reads "# BEGIN gogl managed hosts domain example.test", so the domain is
// the token after the "domain" keyword. Fields rather than an index scan, because
// "domain" also appears inside the marker text of no other line gogl writes, and
// tokenizing makes trailing whitespace a non-issue.
func parseDomain(line string) string {
	fields := strings.Fields(line)
	for i, f := range fields {
		if f == domainKey && i+1 < len(fields) {
			return fields[i+1]
		}
	}
	return ""
}

func parseEntries(block string) []HostEntry {
	var entries []HostEntry
	for _, line := range strings.Split(block, "\n") {
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = line[:idx]
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		entries = append(entries, HostEntry{IP: fields[0], Names: fields[1:]})
	}
	return entries
}

// Render produces the full host-file content: unmanaged parts verbatim, with
// gogl's block rebuilt between them.
func (f *HostFile) Render() string {
	var b strings.Builder

	b.WriteString(f.Before)
	if f.Before != "" && !strings.HasSuffix(f.Before, "\n") {
		b.WriteString("\n")
	}

	// A block with neither a domain nor entries says nothing. Writing one anyway
	// meant that clearing an already-empty router still pushed a file to the
	// device, which is how a no-op came to fail.
	if f.Domain != "" || len(f.Entries) > 0 {
		b.WriteString(BeginMarker)
		if f.Domain != "" {
			fmt.Fprintf(&b, " %s %s", domainKey, f.Domain)
		}
		b.WriteString("\n")

		for _, e := range f.Entries {
			fmt.Fprintf(&b, "%s %s\n", e.IP, strings.Join(e.Names, " "))
		}
		b.WriteString(EndMarker + "\n")
	}

	b.WriteString(f.After)
	return b.String()
}

// FQDN returns the fully-qualified form of name under the configured domain, or
// just the name when no domain is set.
func (f *HostFile) FQDN(name string) string {
	if f.Domain == "" || strings.Contains(name, ".") {
		return name
	}
	return name + "." + f.Domain
}

// bothForms returns the bare and qualified spellings of name.
//
// A caller may pass either, and means the same host regardless. Matching only the
// spelling given would strip half an entry -- removing "nas.lab.example" would
// leave a bare "nas" still answering.
func (f *HostFile) bothForms(name string) (bare, qualified string) {
	bare = name
	if f.Domain != "" {
		bare = strings.TrimSuffix(name, "."+f.Domain)
	}
	return bare, f.FQDN(bare)
}

// Set adds or replaces the entry for name, pointing it at ip.
//
// The entry carries both the bare name and its FQDN, so either resolves. Any other
// entry claiming the same name is dropped: two answers for one name is not a state
// worth preserving.
func (f *HostFile) Set(name, ip string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if parsed := net.ParseIP(ip); parsed == nil || parsed.To4() == nil {
		return fmt.Errorf("%w: %q", ErrInvalidIP, ip)
	}

	names := []string{name}
	if fqdn := f.FQDN(name); fqdn != name {
		names = append(names, fqdn)
	}

	f.Remove(name)
	f.Entries = append(f.Entries, HostEntry{IP: ip, Names: names})
	return nil
}

// Remove drops every entry answering name, in bare or qualified form. It reports
// whether anything was removed.
func (f *HostFile) Remove(name string) bool {
	bare, qualified := f.bothForms(name)

	kept := make([]HostEntry, 0, len(f.Entries))
	removed := false
	for _, e := range f.Entries {
		survivors := make([]string, 0, len(e.Names))
		for _, n := range e.Names {
			if strings.EqualFold(n, bare) || strings.EqualFold(n, qualified) {
				removed = true
				continue
			}
			survivors = append(survivors, n)
		}
		// An entry with every name stripped has nothing left to answer.
		if len(survivors) > 0 {
			kept = append(kept, HostEntry{IP: e.IP, Names: survivors})
		}
	}
	f.Entries = kept
	return removed
}

// Lookup returns the address for name, in bare or qualified form.
func (f *HostFile) Lookup(name string) (string, bool) {
	bare, qualified := f.bothForms(name)
	for _, e := range f.Entries {
		for _, n := range e.Names {
			if strings.EqualFold(n, bare) || strings.EqualFold(n, qualified) {
				return e.IP, true
			}
		}
	}
	return "", false
}

// Clear removes every managed entry, leaving the block and its domain in place.
func (f *HostFile) Clear() { f.Entries = nil }
