package netcfg

import (
	"flag"
	"testing"
)

// These types exist for one reason: a partial update must distinguish "set it to
// false" from "leave it alone", and flag.Bool cannot. Every test here is about that
// distinction, because getting it wrong means --set-key silently resets --set-enabled.

func TestOptionalBoolUnsetIsNil(t *testing.T) {
	var b optionalBool
	if b.Ptr() != nil {
		t.Error("an unparsed optionalBool is not nil")
	}
}

func TestOptionalBoolFalseIsNotNil(t *testing.T) {
	var b optionalBool
	if err := b.Set("false"); err != nil {
		t.Fatalf("Set(false): %v", err)
	}

	got := b.Ptr()
	if got == nil {
		t.Fatal("--set-x=false produced nil, which would leave the field alone")
	}
	if *got {
		t.Error("--set-x=false produced true")
	}
}

func TestOptionalBoolSpellings(t *testing.T) {
	truthy := []string{"true", "TRUE", "yes", "on", "1", " true "}
	for _, raw := range truthy {
		var b optionalBool
		if err := b.Set(raw); err != nil {
			t.Errorf("Set(%q): %v", raw, err)
			continue
		}
		if v, _ := b.Get(); !v {
			t.Errorf("Set(%q) gave false", raw)
		}
	}

	falsy := []string{"false", "FALSE", "no", "off", "0"}
	for _, raw := range falsy {
		var b optionalBool
		if err := b.Set(raw); err != nil {
			t.Errorf("Set(%q): %v", raw, err)
			continue
		}
		if v, _ := b.Get(); v {
			t.Errorf("Set(%q) gave true", raw)
		}
	}

	for _, raw := range []string{"", "maybe", "2", "y"} {
		var b optionalBool
		if err := b.Set(raw); err == nil {
			t.Errorf("Set(%q) was accepted", raw)
		}
	}
}

// A bare --set-enabled must not be silently read as true. Requiring the value is what
// stops "--set-enabled --iface x" from disabling something by consuming "--iface" or
// defaulting to on.
func TestOptionalBoolRequiresAValue(t *testing.T) {
	var b optionalBool
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.SetOutput(discard{})
	fs.Var(&b, "set-enabled", "")

	if err := fs.Parse([]string{"--set-enabled"}); err == nil {
		t.Error("a bare --set-enabled was accepted")
	}
}

func TestOptionalIntZeroIsNotUnset(t *testing.T) {
	var i optionalInt
	if i.Ptr() != nil {
		t.Error("an unparsed optionalInt is not nil")
	}

	if err := i.Set("0"); err != nil {
		t.Fatalf("Set(0): %v", err)
	}
	got := i.Ptr()
	if got == nil {
		t.Fatal("--set-channel=0 produced nil; channel 0 means auto and must be writable")
	}
	if *got != 0 {
		t.Errorf("--set-channel=0 produced %d", *got)
	}
}

func TestOptionalIntRejectsNonNumbers(t *testing.T) {
	var i optionalInt
	for _, raw := range []string{"", "auto", "36MHz"} {
		if err := i.Set(raw); err == nil {
			t.Errorf("Set(%q) was accepted", raw)
		}
	}
}

func TestOptionalStringEmptyIsNotUnset(t *testing.T) {
	var s optionalString
	if s.Ptr() != nil {
		t.Error("an unparsed optionalString is not nil")
	}

	if err := s.Set(""); err != nil {
		t.Fatalf("Set(\"\"): %v", err)
	}
	got := s.Ptr()
	if got == nil {
		t.Fatal("--set-key= produced nil rather than an empty value")
	}
	if *got != "" {
		t.Errorf("--set-key= produced %q", *got)
	}
}

// Parsing through a real FlagSet, because that is where the distinction has to
// survive: the whole point is what reaches the services layer.
func TestOptionalFlagsThroughAFlagSet(t *testing.T) {
	var ssid, key optionalString
	var hidden, enabled optionalBool
	var channel optionalInt

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Var(&ssid, "set-ssid", "")
	fs.Var(&key, "set-key", "")
	fs.Var(&hidden, "set-hidden", "")
	fs.Var(&enabled, "set-enabled", "")
	fs.Var(&channel, "set-channel", "")

	if err := fs.Parse([]string{"--set-ssid", "player-test", "--set-hidden=false"}); err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if ssid.Ptr() == nil || *ssid.Ptr() != "player-test" {
		t.Errorf("ssid = %v", ssid.Ptr())
	}
	if hidden.Ptr() == nil || *hidden.Ptr() {
		t.Errorf("hidden = %v, want a non-nil false", hidden.Ptr())
	}
	// The flags nobody passed must stay nil, or this write would clobber them.
	if key.Ptr() != nil {
		t.Errorf("key was not passed but is %v", key.Ptr())
	}
	if enabled.Ptr() != nil {
		t.Errorf("enabled was not passed but is %v", enabled.Ptr())
	}
	if channel.Ptr() != nil {
		t.Errorf("channel was not passed but is %v", channel.Ptr())
	}
}

type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }
