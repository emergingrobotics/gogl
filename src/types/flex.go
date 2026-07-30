package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexString is a string that also decodes from a JSON number.
//
// It exists because GL.iNet's API description is not a reliable guide to its wire
// format. `clients.online_time` is documented as a string; firmware 4.3.28 sends a
// number. Typing it from the description made every `clients.get_list` decode fail.
//
// This is the fourth time the description has been wrong in a way only hardware could
// reveal: `dhcp.*` does not exist, `dns.set_host` rejects three characters it does not
// mention, `wifi.htmodes` is an object where an array is documented, and now this. The
// project's original decision not to port gofi's FlexInt/FlexBool said to add such types
// only when a fixture proved the need. A fixture has now proved it twice.
//
// The underlying value is kept as text so a caller can decide how to read it: a numeric
// timestamp, a numeric duration, or a formatted date are all possible, and guessing which
// is how a lease "expires" field nearly got rendered as a 1970 date.
type FlexString string

// UnmarshalJSON accepts a JSON string, number, or null.
func (f *FlexString) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	switch {
	case len(b) == 0, bytes.Equal(b, []byte("null")):
		*f = ""
		return nil
	case b[0] == '"':
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*f = FlexString(s)
		return nil
	default:
		// A number, or anything else scalar: keep the literal text. json.Number would
		// reject a bool or a bare token, and this type's job is to not fail.
		if !isJSONScalar(b) {
			return fmt.Errorf("cannot decode %s into a string-or-number field", b)
		}
		*f = FlexString(b)
		return nil
	}
}

// MarshalJSON always emits a string, so gogl's own output has one stable shape
// regardless of what the firmware sent.
func (f FlexString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}

// String returns the value as text.
func (f FlexString) String() string { return string(f) }

// Int64 parses the value as an integer, reporting whether it is one.
func (f FlexString) Int64() (int64, bool) {
	n, err := strconv.ParseInt(string(f), 10, 64)
	return n, err == nil
}

// isJSONScalar reports whether b looks like a number, bool, or other bare token rather
// than an object or array.
func isJSONScalar(b []byte) bool {
	return b[0] != '{' && b[0] != '['
}
