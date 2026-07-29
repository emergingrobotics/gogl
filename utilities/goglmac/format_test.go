package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func testEntries() []Entry {
	return []Entry{
		{MAC: "00:1b:63:00:00:01", IP: "192.168.8.13", Name: "nas", Manufacturer: "Apple, Inc.", IsWired: true, Online: true, Reserved: true},
		{MAC: "02:aa:bb:cc:dd:ee", Name: "unknown", Manufacturer: randomizedManufacturer, Online: true},
	}
}

func TestFormatText(t *testing.T) {
	var buf bytes.Buffer
	if err := formatText(&buf, testEntries(), false); err != nil {
		t.Fatalf("formatText error: %v", err)
	}
	out := buf.String()

	for _, want := range []string{"00:1b:63:00:00:01", "192.168.8.13", "nas", "Apple, Inc.", randomizedManufacturer} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	// A client with no address shows a dash, not an empty column.
	if !strings.Contains(out, noAddress) {
		t.Errorf("address-less client not marked with %q:\n%s", noAddress, out)
	}
}

func TestFormatTextWithReservedColumn(t *testing.T) {
	var buf bytes.Buffer
	if err := formatText(&buf, testEntries(), true); err != nil {
		t.Fatalf("formatText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, stateReserved) || !strings.Contains(out, stateDynamic) {
		t.Errorf("reserved column missing:\n%s", out)
	}
}

func TestFormatTextWithoutReservedColumn(t *testing.T) {
	var buf bytes.Buffer
	if err := formatText(&buf, testEntries(), false); err != nil {
		t.Fatalf("formatText error: %v", err)
	}
	if strings.Contains(buf.String(), stateDynamic) {
		t.Errorf("reserved column present without -r:\n%s", buf.String())
	}
}

func TestFormatTextEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := formatText(&buf, nil, false); err != nil {
		t.Fatalf("formatText error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for no entries, got %q", buf.String())
	}
}

func TestFormatJSONIsArray(t *testing.T) {
	var buf bytes.Buffer
	if err := formatJSON(&buf, testEntries()); err != nil {
		t.Fatalf("formatJSON error: %v", err)
	}

	var decoded []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not a JSON array: %v", err)
	}
	if len(decoded) != 2 {
		t.Fatalf("decoded %d entries, want 2", len(decoded))
	}
	if decoded[0]["manufacturer"] != "Apple, Inc." {
		t.Errorf("manufacturer = %v", decoded[0]["manufacturer"])
	}
	// Zero-valued fields are omitted, so a consumer can tell "not reported" from
	// "reported as zero".
	if _, present := decoded[1]["ip"]; present {
		t.Error("empty ip should be omitted from JSON")
	}
	if _, present := decoded[1]["reserved"]; present {
		t.Error("false reserved should be omitted from JSON")
	}
}

// A nil slice must marshal as [] rather than null, so piping into jq never needs
// an empty-result special case.
func TestFormatJSONEmptyIsEmptyArray(t *testing.T) {
	var buf bytes.Buffer
	if err := formatJSON(&buf, nil); err != nil {
		t.Fatalf("formatJSON error: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "[]" {
		t.Errorf("empty output = %q, want []", got)
	}
}
