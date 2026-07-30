package netcfg

import (
	"errors"
	"strconv"
	"strings"
)

// Go's flag package cannot distinguish "--set-hidden=false" from "--set-hidden not
// given", because both leave a *bool at false. That difference is the whole contract
// of a partial update: an absent field must be left alone, and a false one must be
// written. These two types carry the distinction.
//
// The alternative, a pair of flags per field (--set-hidden / --set-visible), doubles
// the surface and still cannot express "leave it alone" when both are absent -- it
// just hides the problem behind more flags.

// optionalBool is a bool flag that remembers whether it was given.
type optionalBool struct {
	value bool
	set   bool
}

func (b *optionalBool) String() string {
	if !b.set {
		return ""
	}
	return strconv.FormatBool(b.value)
}

func (b *optionalBool) Set(raw string) error {
	// Accept the spellings an operator will reach for, rather than only Go's.
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "true", "yes", "on", "1":
		b.value, b.set = true, true
	case "false", "no", "off", "0":
		b.value, b.set = false, true
	default:
		return errors.New(`must be true or false`)
	}
	return nil
}

// IsBoolFlag is false deliberately: it forces "--set-hidden=true" rather than
// allowing a bare "--set-hidden", so the intent is always written down. A bare flag
// that means true is exactly how "--set-enabled" ends up disabling something.
func (b *optionalBool) IsBoolFlag() bool { return false }

// Get returns the value and whether it was set.
func (b *optionalBool) Get() (bool, bool) { return b.value, b.set }

// Ptr returns a pointer to the value, or nil if the flag was not given, which is the
// shape the services layer takes for a partial update.
func (b *optionalBool) Ptr() *bool {
	if !b.set {
		return nil
	}
	v := b.value
	return &v
}

// optionalInt is an int flag that remembers whether it was given.
//
// Needed because 0 is a meaningful value here: channel 0 means "pick one
// automatically", so a zero-versus-unset sentinel is unavailable.
type optionalInt struct {
	value int
	set   bool
}

func (i *optionalInt) String() string {
	if !i.set {
		return ""
	}
	return strconv.Itoa(i.value)
}

func (i *optionalInt) Set(raw string) error {
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return errors.New("must be a number")
	}
	i.value, i.set = v, true
	return nil
}

// Ptr returns a pointer to the value, or nil if the flag was not given.
func (i *optionalInt) Ptr() *int {
	if !i.set {
		return nil
	}
	v := i.value
	return &v
}

// optionalString is a string flag that remembers whether it was given, so that
// setting a value to the empty string stays distinguishable from not setting it.
type optionalString struct {
	value string
	set   bool
}

func (s *optionalString) String() string { return s.value }

func (s *optionalString) Set(raw string) error {
	s.value, s.set = raw, true
	return nil
}

// Ptr returns a pointer to the value, or nil if the flag was not given.
func (s *optionalString) Ptr() *string {
	if !s.set {
		return nil
	}
	v := s.value
	return &v
}
