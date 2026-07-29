package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Bands the firmware reports. wifi.get_config returns one entry per radio.
const (
	Band2G = "2G"
	Band5G = "5G"
)

// SSID and passphrase limits, from wifi.set_config's own documentation. Enforced
// before the write so a rejection names the field rather than arriving as a bare
// firmware error code.
const (
	MaxSSIDLength = 32
	MinKeyLength  = 8
	MaxKeyLength  = 63
)

// WirelessInterface is one SSID on one radio.
//
// Field names mirror wifi.get_config verbatim. Note that Key is the passphrase in
// cleartext: the firmware returns it that way, and callers that print it should
// consider where the output lands.
type WirelessInterface struct {
	// Name is the interface identifier wifi.set_config requires as iface_name,
	// such as "default_radio0" or "guest5g". It is not the SSID and not a band.
	Name string `json:"name"`

	SSID       string `json:"ssid"`
	Key        string `json:"key"`
	Encryption string `json:"encryption"`

	Enabled bool `json:"enabled"`
	Hidden  bool `json:"hidden"`
	Guest   bool `json:"guest"`

	// Band is filled in from the enclosing radio, which is where the firmware
	// reports it. Carrying it here means a caller holding one interface knows which
	// radio it belongs to without walking back up.
	Band string `json:"-"`
}

// WirelessChannel is one selectable channel and whether it is subject to DFS.
//
// DFS channels require radar detection: the radio must vacate one if it sees a
// weather or military radar pulse, taking every client with it for the minutes it
// spends re-scanning. Fine for a fixed install, a poor choice for kit that has to
// come up reliably in an unfamiliar building.
type WirelessChannel struct {
	Channel int  `json:"channel"`
	DFS     bool `json:"dfs"`
}

// TXPower levels the firmware accepts.
var TXPowerLevels = []string{"Max", "High", "Medium", "Low"}

// AutoHTMode is the htmode value meaning "choose a width automatically".
const AutoHTMode = "auto"

// Channel widths in MHz that htmode can name, narrowest first.
var htModeWidths = []int{20, 40, 80, 160}

// HTModes is what wifi.get_config reports under "htmodes".
//
// CAPTURED from a GL-SFT1200 on 4.3.28, because the vendored API description is
// wrong about this field. It calls htmodes a "List of supported bandwidths", an
// array of strings. The device sends an object keyed by hardware mode whose values
// are the maximum channel width in MHz, plus an "auto" key whose value is a bool:
//
//	"htmodes": {"11b/g/n": 40, "11g/n": 40, "11n": 40, "auto": true}
//	"htmodes": {"11a/n/ac": 80, "11ac": 80, "11n/ac": 80, "auto": false}
//
// The current htmode is a width string, not the "HT20"/"VHT80" form the description
// implies: observed values are "auto" on the 2.4GHz radio and "20" on the 5GHz one.
type HTModes struct {
	// Auto reports whether "auto" is an available htmode.
	Auto bool

	// MaxWidthMHz maps a hardware mode to the widest channel it supports.
	MaxWidthMHz map[string]int
}

// UnmarshalJSON reads the mixed-type object the firmware sends.
func (h *HTModes) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("htmodes is not an object: %w", err)
	}

	h.MaxWidthMHz = make(map[string]int, len(raw))
	for key, value := range raw {
		if key == AutoHTMode {
			if err := json.Unmarshal(value, &h.Auto); err != nil {
				return fmt.Errorf("htmodes.auto is not a bool: %w", err)
			}
			continue
		}
		var width int
		if err := json.Unmarshal(value, &width); err != nil {
			return fmt.Errorf("htmodes[%q] is not a number: %w", key, err)
		}
		h.MaxWidthMHz[key] = width
	}
	return nil
}

// MarshalJSON writes the same shape back, so a round trip through this type is
// faithful. The mock relies on it.
func (h HTModes) MarshalJSON() ([]byte, error) {
	out := make(map[string]any, len(h.MaxWidthMHz)+1)
	for mode, width := range h.MaxWidthMHz {
		out[mode] = width
	}
	out[AutoHTMode] = h.Auto
	return json.Marshal(out)
}

// Options returns the htmode values settable for hwmode, narrowest width first,
// with "auto" included when available.
//
// Derived rather than reported: the firmware says only what the widest width is, so
// the narrower ones are inferred as also valid. That inference is why an unsupported
// width is reported as a refusal naming these options rather than as a hard fact
// about the radio.
func (h *HTModes) Options(hwmode string) []string {
	max, ok := h.MaxWidthMHz[hwmode]
	if !ok {
		// An unrecognized hardware mode: fall back to the widest this radio reports
		// anywhere, so a caller mid-hwmode-change still gets useful options.
		for _, width := range h.MaxWidthMHz {
			if width > max {
				max = width
			}
		}
	}

	var out []string
	if h.Auto {
		out = append(out, AutoHTMode)
	}
	for _, width := range htModeWidths {
		if width <= max {
			out = append(out, strconv.Itoa(width))
		}
	}
	return out
}

// AutoChannel is what the firmware means by channel 0: pick one.
const AutoChannel = 0

// WirelessRadio is one physical radio, the interfaces on it, and what it supports.
//
// The supported-value lists are why radio writes can be validated at all: the
// firmware answers a bad channel with a bare error code, while the radio itself
// already said which channels exist.
type WirelessRadio struct {
	Band    string              `json:"band"`
	Device  string              `json:"device"`
	Channel int                 `json:"channel"`
	HTMode  string              `json:"htmode"`
	HWMode  string              `json:"hwmode"`
	TXPower string              `json:"txpower"`
	Ifaces  []WirelessInterface `json:"ifaces"`

	// HWModes are hardware modes such as "11g/n" or "11a/n/ac". Note these are
	// slash-joined combinations, not the bare "11b"/"11g"/"11n" the vendored
	// description implies.
	HWModes []string `json:"hwmodes"`

	HTModes     HTModes           `json:"htmodes"`
	Channels    []WirelessChannel `json:"channels"`
	Encryptions []string          `json:"encryptions"`
}

// HTModeOptions returns the htmode values settable on this radio as it is currently
// configured.
func (r *WirelessRadio) HTModeOptions() []string {
	return r.HTModes.Options(r.HWMode)
}

// ChannelNumbers returns the selectable channels.
func (r *WirelessRadio) ChannelNumbers() []int {
	out := make([]int, 0, len(r.Channels))
	for _, c := range r.Channels {
		out = append(out, c.Channel)
	}
	return out
}

// IsDFS reports whether channel requires radar detection on this radio.
func (r *WirelessRadio) IsDFS(channel int) bool {
	for _, c := range r.Channels {
		if c.Channel == channel {
			return c.DFS
		}
	}
	return false
}

// InterfaceChanges is a partial update to one wireless interface. A nil field is
// left alone, which is what wifi.set_config does with an absent one.
type InterfaceChanges struct {
	SSID       *string
	Key        *string
	Encryption *string
	Hidden     *bool
	Enabled    *bool
}

// Empty reports whether there is nothing to write.
func (c *InterfaceChanges) Empty() bool {
	return c.SSID == nil && c.Key == nil && c.Encryption == nil &&
		c.Hidden == nil && c.Enabled == nil
}

// Validate checks the requested changes against what radio supports.
func (c *InterfaceChanges) Validate(radio *WirelessRadio) error {
	if c.SSID != nil {
		if err := ValidateSSID(*c.SSID); err != nil {
			return err
		}
	}
	if c.Key != nil {
		if err := ValidateWirelessKey(*c.Key); err != nil {
			return err
		}
	}
	if c.Encryption != nil && radio != nil && len(radio.Encryptions) > 0 {
		if !contains(radio.Encryptions, *c.Encryption) {
			return fmt.Errorf("%w: encryption %q is not supported on %s (have: %s)",
				ErrInvalidInput, *c.Encryption, radio.Band, strings.Join(radio.Encryptions, ", "))
		}
	}
	return nil
}

// RadioChanges is a partial update to one radio's tuning.
type RadioChanges struct {
	Channel *int
	HTMode  *string
	HWMode  *string
	TXPower *string
}

// Empty reports whether there is nothing to write.
func (c *RadioChanges) Empty() bool {
	return c.Channel == nil && c.HTMode == nil && c.HWMode == nil && c.TXPower == nil
}

// Validate checks the requested tuning against what radio advertises.
//
// The radio reports its own supported channels, bandwidths and hardware modes, so a
// bad value is caught here with the valid ones named, rather than arriving as an
// error code that says only "Invalid params".
func (c *RadioChanges) Validate(radio *WirelessRadio) error {
	if radio == nil {
		return fmt.Errorf("%w: no radio to validate against", ErrInvalidInput)
	}

	if c.Channel != nil && *c.Channel != AutoChannel && len(radio.Channels) > 0 {
		if !containsInt(radio.ChannelNumbers(), *c.Channel) {
			// "the firmware does not offer" rather than "the radio cannot": the list
			// comes from wifi.get_config, which reports GL.iNet's policy and not the
			// driver's capability. On the captured SFT1200 the API offers nine 5GHz
			// channels while iw phy reports twenty-five, the difference being every
			// DFS channel. Blaming the hardware would send someone looking for a
			// hardware answer to a firmware question.
			return fmt.Errorf("%w: the firmware does not offer channel %d on %s (offers: %s, or 0 for auto)",
				ErrInvalidInput, *c.Channel, radio.Band, joinInts(radio.ChannelNumbers()))
		}
	}
	if c.HTMode != nil {
		options := radio.HTModeOptions()
		if len(options) > 0 && !contains(options, *c.HTMode) {
			return fmt.Errorf("%w: bandwidth %q is not settable on %s with hardware mode %s (have: %s)",
				ErrInvalidInput, *c.HTMode, radio.Band, radio.HWMode, strings.Join(options, ", "))
		}
	}
	if c.HWMode != nil && len(radio.HWModes) > 0 && !contains(radio.HWModes, *c.HWMode) {
		return fmt.Errorf("%w: hardware mode %q is not supported on %s (have: %s)",
			ErrInvalidInput, *c.HWMode, radio.Band, strings.Join(radio.HWModes, ", "))
	}
	if c.TXPower != nil && !contains(TXPowerLevels, *c.TXPower) {
		return fmt.Errorf("%w: transmit power %q is not a level (have: %s)",
			ErrInvalidInput, *c.TXPower, strings.Join(TXPowerLevels, ", "))
	}
	return nil
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

func containsInt(haystack []int, needle int) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

func joinInts(values []int) string {
	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, strconv.Itoa(v))
	}
	return strings.Join(parts, ", ")
}

// MaskedKey returns the passphrase's length rather than its value, for output that
// should confirm a key is set without disclosing it.
func (w *WirelessInterface) MaskedKey() string {
	if w.Key == "" {
		return "(none)"
	}
	return fmt.Sprintf("(%d characters, --show-key to reveal)", utf8.RuneCountInString(w.Key))
}

// Describe names the interface the way an operator would recognize it.
func (w *WirelessInterface) Describe() string {
	kind := "main"
	if w.Guest {
		kind = "guest"
	}
	return fmt.Sprintf("%s %s (%s)", w.Band, kind, w.Name)
}

// ValidateSSID reports whether ssid is writable.
//
// Rejected rather than trimmed or escaped: an SSID that is silently altered is one
// the operator's devices will not find, and discovering that means walking to the
// hardware.
func ValidateSSID(ssid string) error {
	if ssid == "" {
		return fmt.Errorf("%w: an SSID cannot be empty", ErrInvalidInput)
	}
	if n := utf8.RuneCountInString(ssid); n > MaxSSIDLength {
		return fmt.Errorf("%w: SSID is %d characters, the limit is %d",
			ErrInvalidInput, n, MaxSSIDLength)
	}
	if !utf8.ValidString(ssid) {
		return fmt.Errorf("%w: SSID is not valid UTF-8", ErrInvalidInput)
	}
	// Leading or trailing whitespace in an SSID is invisible on a client's network
	// list and a genuine source of "it says the name is right but it will not
	// connect".
	if strings.TrimSpace(ssid) != ssid {
		return fmt.Errorf("%w: SSID has leading or trailing whitespace, which is invisible to clients",
			ErrInvalidInput)
	}
	for _, r := range ssid {
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("%w: SSID contains a control character", ErrInvalidInput)
		}
	}
	return nil
}

// ValidateWirelessKey reports whether key is writable as a WPA passphrase.
func ValidateWirelessKey(key string) error {
	if n := len(key); n < MinKeyLength || n > MaxKeyLength {
		return fmt.Errorf("%w: passphrase is %d characters, must be %d to %d",
			ErrInvalidInput, n, MinKeyLength, MaxKeyLength)
	}
	return nil
}
