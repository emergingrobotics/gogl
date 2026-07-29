package main

import (
	"strings"
	"testing"
)

const wellFormed = `# gofips fixed IP assignments
# exported from UDM at 192.168.4.1

host myserver {
    hardware ethernet aa:bb:cc:dd:ee:01;
    fixed-address 192.168.4.10;
}

host printer {
    hardware ethernet AA:BB:CC:DD:EE:02;
    fixed-address 192.168.4.11;
}
`

func TestParseHosts(t *testing.T) {
	got, errs := ParseHosts(strings.NewReader(wellFormed))
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(got) != 2 {
		t.Fatalf("parsed %d declarations, want 2", len(got))
	}

	if got[0].Reservation.Name != "myserver" {
		t.Errorf("first name = %q, want myserver", got[0].Reservation.Name)
	}
	// MACs normalize to lowercase on input regardless of how they were written.
	if got[1].Reservation.MAC != "aa:bb:cc:dd:ee:02" {
		t.Errorf("second MAC = %q, want lowercase", got[1].Reservation.MAC)
	}
	if got[0].Reservation.IP != "192.168.4.10" {
		t.Errorf("first IP = %q", got[0].Reservation.IP)
	}
	if got[0].Line != 4 {
		t.Errorf("first declaration Line = %d, want 4", got[0].Line)
	}
}

// Whitespace is flexible on input: tabs, no indentation, and the whole block on
// one line all parse.
func TestParseHostsToleratesWhitespace(t *testing.T) {
	inputs := []string{
		"host a {\nhardware ethernet aa:bb:cc:dd:ee:01;\nfixed-address 10.0.0.1;\n}\n",
		"host a {\n\t\thardware ethernet aa:bb:cc:dd:ee:01;\n\t\tfixed-address 10.0.0.1;\n}\n",
		"  host a  {  hardware ethernet aa:bb:cc:dd:ee:01;  fixed-address 10.0.0.1;  }  \n",
		"host a{hardware ethernet aa:bb:cc:dd:ee:01;fixed-address 10.0.0.1;}",
	}
	for i, in := range inputs {
		got, errs := ParseHosts(strings.NewReader(in))
		if len(errs) != 0 {
			t.Errorf("input %d: errors %v", i, errs)
			continue
		}
		if len(got) != 1 || got[0].Reservation.Name != "a" {
			t.Errorf("input %d: got %v", i, got)
		}
	}
}

// Non-host directives are ignored, so a real dhcpd.conf can be fed in directly.
func TestParseHostsIgnoresOtherDirectives(t *testing.T) {
	const input = `
option domain-name "example.org";
default-lease-time 600;

subnet 192.168.4.0 netmask 255.255.255.0 {
    range 192.168.4.100 192.168.4.200;
    option routers 192.168.4.1;
}

host myserver {
    hardware ethernet aa:bb:cc:dd:ee:01;
    fixed-address 192.168.4.10;
}
`
	got, errs := ParseHosts(strings.NewReader(input))
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	if len(got) != 1 || got[0].Reservation.Name != "myserver" {
		t.Errorf("got %v, want just myserver", got)
	}
}

// Unknown statements inside a host block are ignored rather than fatal.
func TestParseHostsIgnoresUnknownStatementsInBlock(t *testing.T) {
	const input = `host a {
    hardware ethernet aa:bb:cc:dd:ee:01;
    fixed-address 10.0.0.1;
    option host-name "whatever";
}
`
	got, errs := ParseHosts(strings.NewReader(input))
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	if len(got) != 1 {
		t.Errorf("parsed %d declarations, want 1", len(got))
	}
}

// Errors carry line numbers, since the point is to find the problem in a file.
func TestParseHostsReportsLineNumbers(t *testing.T) {
	const input = `host good {
    hardware ethernet aa:bb:cc:dd:ee:01;
    fixed-address 192.168.4.10;
}

host bad {
    hardware ethernet not-a-mac;
    fixed-address 192.168.4.11;
}
`
	got, errs := ParseHosts(strings.NewReader(input))
	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Error(), "line 7") {
		t.Errorf("error %q does not name line 7", errs[0])
	}
	// The good block must still parse: one bad block does not poison the file.
	if len(got) != 1 || got[0].Reservation.Name != "good" {
		t.Errorf("got %v, want the good declaration", got)
	}
}

func TestParseHostsRejects(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"missing mac", "host a {\n fixed-address 10.0.0.1;\n}\n", "hardware ethernet"},
		{"missing address", "host a {\n hardware ethernet aa:bb:cc:dd:ee:01;\n}\n", "fixed-address"},
		{"missing semicolon", "host a {\n hardware ethernet aa:bb:cc:dd:ee:01\n fixed-address 10.0.0.1;\n}\n", "semicolon"},
		{"trailing statement no semicolon", "host a {\n hardware ethernet aa:bb:cc:dd:ee:01;\n fixed-address 10.0.0.1\n}\n", "semicolon"},
		{"unclosed block", "host a {\n hardware ethernet aa:bb:cc:dd:ee:01;\n fixed-address 10.0.0.1;\n", "unclosed"},
		{"bad ip", "host a {\n hardware ethernet aa:bb:cc:dd:ee:01;\n fixed-address 999.1.1.1;\n}\n", "IPv4"},
		{"ipv6 address", "host a {\n hardware ethernet aa:bb:cc:dd:ee:01;\n fixed-address fe80::1;\n}\n", "IPv4"},
		{"no hostname", "host {\n hardware ethernet aa:bb:cc:dd:ee:01;\n fixed-address 10.0.0.1;\n}\n", "hostname"},
		{"missing brace", "host a\n hardware ethernet aa:bb:cc:dd:ee:01;\n", "brace"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := ParseHosts(strings.NewReader(tt.input))
			if len(errs) == 0 {
				t.Fatalf("parsed without error, want one mentioning %q", tt.want)
			}
			if !strings.Contains(errs[0].Error(), tt.want) {
				t.Errorf("error %q does not mention %q", errs[0], tt.want)
			}
		})
	}
}

// The one place a gofips file may not import unchanged: an underscore is legal on
// UniFi but is not a legal DNS label character, so it is rejected rather than
// silently rewritten.
func TestParseHostsRejectsUnderscoreName(t *testing.T) {
	const input = "host my_server {\n hardware ethernet aa:bb:cc:dd:ee:01;\n fixed-address 10.0.0.1;\n}\n"
	_, errs := ParseHosts(strings.NewReader(input))
	if len(errs) == 0 {
		t.Fatal("accepted an underscore in a hostname")
	}
	if !strings.Contains(errs[0].Error(), "_") {
		t.Errorf("error %q does not name the offending character", errs[0])
	}
}

// All errors are collected, not just the first, so one run finds every problem.
func TestParseHostsCollectsAllErrors(t *testing.T) {
	const input = `host bad1 {
    hardware ethernet nope;
    fixed-address 10.0.0.1;
}
host bad2 {
    hardware ethernet aa:bb:cc:dd:ee:02;
    fixed-address 999.9.9.9;
}
host bad3 {
    hardware ethernet aa:bb:cc:dd:ee:03;
    fixed-address 10.0.0.3;
`
	_, errs := ParseHosts(strings.NewReader(input))
	if len(errs) < 3 {
		t.Errorf("got %d errors, want at least 3: %v", len(errs), errs)
	}
}

func TestParseHostsEmpty(t *testing.T) {
	for _, input := range []string{"", "# only a comment\n", "\n\n\n"} {
		got, errs := ParseHosts(strings.NewReader(input))
		if len(errs) != 0 {
			t.Errorf("input %q: errors %v", input, errs)
		}
		if len(got) != 0 {
			t.Errorf("input %q: parsed %d declarations, want 0", input, len(got))
		}
	}
}

// A comment on the same line as a statement must not break parsing.
func TestParseHostsInlineComment(t *testing.T) {
	const input = `host a {
    hardware ethernet aa:bb:cc:dd:ee:01;  # the NAS
    fixed-address 10.0.0.1;
}
`
	got, errs := ParseHosts(strings.NewReader(input))
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	if len(got) != 1 {
		t.Errorf("parsed %d declarations, want 1", len(got))
	}
}

// A fragment with no trailing newline is what --add passes in.
func TestParseHostsFragmentWithoutTrailingNewline(t *testing.T) {
	const input = "host a {\n    hardware ethernet aa:bb:cc:dd:ee:01;\n    fixed-address 10.0.0.1;\n}"
	got, errs := ParseHosts(strings.NewReader(input))
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	if len(got) != 1 {
		t.Errorf("parsed %d declarations, want 1", len(got))
	}
}
