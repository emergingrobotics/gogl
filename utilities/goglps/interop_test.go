package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// These fixtures replicate gofips's output byte for byte, transcribed from
// gofi's utilities/gofips/format.go rather than invented: the header lines, the
// blank line before each block, the four-space indent, and the closing brace on
// its own line.
//
// This file is the interoperability contract between gofi and gogl, and the only
// thing that would catch gofips changing its output format. If it starts failing,
// check gofips's formatter before changing anything here.
const (
	exportFixture = "testdata/gofips-export.hosts"
	emptyFixture  = "testdata/gofips-empty.hosts"
)

func TestParsesRealGofipsExport(t *testing.T) {
	data, err := os.ReadFile(exportFixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	got, errs := ParseHosts(bytes.NewReader(data))
	if len(errs) != 0 {
		t.Fatalf("gofips export did not parse: %v", errs)
	}
	if len(got) != 3 {
		t.Fatalf("parsed %d declarations, want 3", len(got))
	}

	want := map[string]string{
		"myserver": "192.168.4.10",
		"printer":  "192.168.4.11",
		"nas":      "192.168.4.13",
	}
	for _, d := range got {
		wantIP, ok := want[d.Reservation.Name]
		if !ok {
			t.Errorf("unexpected host %q", d.Reservation.Name)
			continue
		}
		if d.Reservation.IP != wantIP {
			t.Errorf("host %q has IP %q, want %q", d.Reservation.Name, d.Reservation.IP, wantIP)
		}
		delete(want, d.Reservation.Name)
	}
	for name := range want {
		t.Errorf("host %q was not parsed", name)
	}
}

// gofips's empty-table output is entirely commented, so it must import as zero
// reservations rather than erroring or inventing an "example" host.
func TestParsesEmptyGofipsExport(t *testing.T) {
	data, err := os.ReadFile(emptyFixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	got, errs := ParseHosts(bytes.NewReader(data))
	if len(errs) != 0 {
		t.Fatalf("empty gofips export did not parse: %v", errs)
	}
	if len(got) != 0 {
		t.Errorf("parsed %d declarations from an empty export, want 0: %v", len(got), got)
	}
}

// Our output must be re-readable by the same rules, so a file can move back and
// forth between the two tools.
func TestOutputIsGofipsCompatible(t *testing.T) {
	data, err := os.ReadFile(exportFixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	parsed, errs := ParseHosts(bytes.NewReader(data))
	if len(errs) != 0 {
		t.Fatalf("parse: %v", errs)
	}

	var buf bytes.Buffer
	if err := FormatHosts(&buf, reservationsOf(parsed), Header{Host: "192.168.4.1", Date: "2026-07-27"}); err != nil {
		t.Fatalf("FormatHosts: %v", err)
	}

	reparsed, errs := ParseHosts(bytes.NewReader(buf.Bytes()))
	if len(errs) != 0 {
		t.Fatalf("our own output did not re-parse: %v", errs)
	}
	if len(reparsed) != len(parsed) {
		t.Errorf("re-parsed %d declarations, want %d", len(reparsed), len(parsed))
	}
}

// The structural details gofips emits must all be tolerated, one at a time, so a
// future divergence points at the specific thing that changed.
func TestGofipsStructuralDetails(t *testing.T) {
	data, err := os.ReadFile(exportFixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	text := string(data)

	if !strings.HasPrefix(text, "# gofips fixed IP assignments") {
		t.Error("fixture no longer carries gofips's header; re-check gofips/format.go")
	}
	if !strings.Contains(text, "\n\nhost ") {
		t.Error("fixture no longer has a blank line before each block")
	}
	if !strings.Contains(text, "\n    hardware ethernet ") {
		t.Error("fixture no longer uses a four-space indent")
	}
	if !strings.Contains(text, ";\n}\n") {
		t.Error("fixture no longer closes blocks on their own line")
	}
}

// A file exported from UniFi is typically in a different subnet than the travel
// router ships with. That is not a parse problem, and must not be treated as one:
// it surfaces later, in validateAgainstDevice.
func TestGofipsExportParsesRegardlessOfSubnet(t *testing.T) {
	data, err := os.ReadFile(exportFixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	parsed, errs := ParseHosts(bytes.NewReader(data))
	if len(errs) != 0 {
		t.Fatalf("parse: %v", errs)
	}

	// The fixture is 192.168.4.x while the router default is 192.168.8.x.
	_, deviceErrs := validateAgainstDevice(reservationsOf(parsed), testLAN())
	if len(deviceErrs) == 0 {
		t.Error("expected a subnet mismatch against the default router LAN")
	}
	if !strings.Contains(errsText(deviceErrs), "3 of 3") {
		t.Errorf("mismatch report should count all three entries:\n%s", errsText(deviceErrs))
	}
}
