package conn

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintVersionWithStamp(t *testing.T) {
	original := Version
	t.Cleanup(func() { Version = original })
	Version = "v0.9.1"

	var buf bytes.Buffer
	PrintVersion(&buf, "gogl")

	got := buf.String()
	if !strings.HasPrefix(got, "gogl v0.9.1") {
		t.Errorf("output = %q, want it to start with the program and version", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("output = %q, want a trailing newline", got)
	}
}

// An unstamped build must still say something useful rather than printing an empty
// version, which reads as a broken binary.
func TestPrintVersionWithoutStamp(t *testing.T) {
	original := Version
	t.Cleanup(func() { Version = original })
	Version = ""

	var buf bytes.Buffer
	PrintVersion(&buf, "gogl")

	got := buf.String()
	if !strings.HasPrefix(got, "gogl ") {
		t.Errorf("output = %q", got)
	}
	if strings.Contains(got, "gogl \n") {
		t.Errorf("output has an empty version: %q", got)
	}
}

// The point of the whole thing: two stale-binary confusions this week were diagnosed by
// comparing `strings` output to source. A revision in the output answers it directly.
func TestPrintVersionIncludesRevisionWhenAvailable(t *testing.T) {
	var buf bytes.Buffer
	PrintVersion(&buf, "gogl")

	// Under `go test` the binary carries no VCS stamp, so this asserts only that the
	// revision path does not corrupt the line when absent.
	got := buf.String()
	if strings.Contains(got, "()") {
		t.Errorf("empty revision parenthetical: %q", got)
	}
	if n := strings.Count(got, "\n"); n != 1 {
		t.Errorf("want exactly one line, got %d: %q", n, got)
	}
}
