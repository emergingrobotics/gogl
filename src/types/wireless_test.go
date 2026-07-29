package types

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// htmodes is the field the vendored API description gets wrong: it documents a "List
// of supported bandwidths" and the device sends an object keyed by hardware mode,
// mixing numbers with a bool. Typing it from the description made every read of
// wifi.get_config fail, so the real payload is pinned here verbatim.

const captured2G = `{"11b/g/n": 40, "11g/n": 40, "11n": 40, "auto": true}`
const captured5G = `{"11a/n/ac": 80, "11ac": 80, "11n/ac": 80, "auto": false}`

func TestHTModesUnmarshalCaptured(t *testing.T) {
	var got HTModes
	if err := json.Unmarshal([]byte(captured2G), &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !got.Auto {
		t.Error("auto is false, want true")
	}
	want := map[string]int{"11b/g/n": 40, "11g/n": 40, "11n": 40}
	if !reflect.DeepEqual(got.MaxWidthMHz, want) {
		t.Errorf("MaxWidthMHz = %v, want %v", got.MaxWidthMHz, want)
	}

	var fiveGHz HTModes
	if err := json.Unmarshal([]byte(captured5G), &fiveGHz); err != nil {
		t.Fatalf("Unmarshal 5G: %v", err)
	}
	if fiveGHz.Auto {
		t.Error("5G auto is true, want false")
	}
	if fiveGHz.MaxWidthMHz["11ac"] != 80 {
		t.Errorf("11ac max = %d, want 80", fiveGHz.MaxWidthMHz["11ac"])
	}
}

// The mock re-marshals the stored config after every write, so a lossy round trip
// would quietly corrupt the fixture between one write and the next.
func TestHTModesRoundTrip(t *testing.T) {
	for _, captured := range []string{captured2G, captured5G} {
		var first HTModes
		if err := json.Unmarshal([]byte(captured), &first); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		encoded, err := json.Marshal(first)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		var second HTModes
		if err := json.Unmarshal(encoded, &second); err != nil {
			t.Fatalf("Unmarshal round trip: %v", err)
		}
		if !reflect.DeepEqual(first, second) {
			t.Errorf("round trip changed the value:\n%+v\n%+v", first, second)
		}
	}
}

func TestHTModesRejectsTheDocumentedShape(t *testing.T) {
	// What the vendored description claims, and what gogl wrongly expected.
	var got HTModes
	err := json.Unmarshal([]byte(`["HT20","HT40"]`), &got)
	if err == nil {
		t.Fatal("an array was accepted")
	}
	if !strings.Contains(err.Error(), "object") {
		t.Errorf("error does not say what was expected: %v", err)
	}
}

func TestHTModeOptions(t *testing.T) {
	twoGHz := WirelessRadio{
		Band: Band2G, HWMode: "11g/n",
		HTModes: HTModes{Auto: true, MaxWidthMHz: map[string]int{"11g/n": 40, "11n": 40}},
	}
	// auto first because it is what the radio ships with, then the widths up to the
	// maximum. 80 must not appear: this radio caps at 40.
	if got := twoGHz.HTModeOptions(); !reflect.DeepEqual(got, []string{"auto", "20", "40"}) {
		t.Errorf("2.4GHz options = %v", got)
	}

	fiveGHz := WirelessRadio{
		Band: Band5G, HWMode: "11a/n/ac",
		HTModes: HTModes{Auto: false, MaxWidthMHz: map[string]int{"11a/n/ac": 80}},
	}
	if got := fiveGHz.HTModeOptions(); !reflect.DeepEqual(got, []string{"20", "40", "80"}) {
		t.Errorf("5GHz options = %v", got)
	}
}

// Mid-hwmode-change, the current hwmode may not be a key in the map. Falling back to
// the widest the radio reports anywhere beats returning nothing, which would refuse
// every bandwidth.
func TestHTModeOptionsWithUnknownHWMode(t *testing.T) {
	r := WirelessRadio{
		HWMode:  "11ax",
		HTModes: HTModes{MaxWidthMHz: map[string]int{"11ac": 80}},
	}
	if got := r.HTModeOptions(); !reflect.DeepEqual(got, []string{"20", "40", "80"}) {
		t.Errorf("options = %v, want the widest the radio reports", got)
	}
}

func TestRadioChangesValidatesHTModeAgainstOptions(t *testing.T) {
	r := &WirelessRadio{
		Band: Band2G, HWMode: "11g/n",
		HTModes: HTModes{Auto: true, MaxWidthMHz: map[string]int{"11g/n": 40}},
	}

	for _, ok := range []string{"auto", "20", "40"} {
		c := RadioChanges{HTMode: &ok}
		if err := c.Validate(r); err != nil {
			t.Errorf("htmode %q rejected: %v", ok, err)
		}
	}

	// 80 is beyond this radio, and "VHT80" is the form the description implied but
	// the device never uses.
	for _, bad := range []string{"80", "VHT80", "HT20"} {
		c := RadioChanges{HTMode: &bad}
		err := c.Validate(r)
		if err == nil {
			t.Errorf("htmode %q accepted", bad)
			continue
		}
		if !strings.Contains(err.Error(), "auto") {
			t.Errorf("error for %q does not list the options: %v", bad, err)
		}
	}
}

// The captured encryption list includes WPA3 (sae) and excludes bare "psk", which the
// hand-written fixture wrongly had.
func TestInterfaceChangesValidatesEncryption(t *testing.T) {
	r := &WirelessRadio{
		Band:        Band5G,
		Encryptions: []string{"none", "psk2", "psk-mixed", "sae", "sae-mixed"},
	}

	for _, ok := range []string{"psk2", "sae", "sae-mixed", "none"} {
		c := InterfaceChanges{Encryption: &ok}
		if err := c.Validate(r); err != nil {
			t.Errorf("encryption %q rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"psk", "wep", "wpa3"} {
		c := InterfaceChanges{Encryption: &bad}
		if err := c.Validate(r); err == nil {
			t.Errorf("encryption %q accepted", bad)
		}
	}
}

func TestIsDFS(t *testing.T) {
	r := WirelessRadio{Channels: []WirelessChannel{
		{Channel: 36, DFS: false},
		{Channel: 52, DFS: true},
	}}
	if r.IsDFS(36) {
		t.Error("36 reported as DFS")
	}
	if !r.IsDFS(52) {
		t.Error("52 not reported as DFS")
	}
	// An unknown channel is not DFS: the warning is for channels the radio listed.
	if r.IsDFS(999) {
		t.Error("an unknown channel reported as DFS")
	}
}

func TestChannelNumbers(t *testing.T) {
	r := WirelessRadio{Channels: []WirelessChannel{
		{Channel: 1}, {Channel: 6}, {Channel: 11},
	}}
	if got := r.ChannelNumbers(); !reflect.DeepEqual(got, []int{1, 6, 11}) {
		t.Errorf("ChannelNumbers = %v", got)
	}
}

// The fixture the mock serves must decode through the same types the device's payload
// does. This is the assertion that would have caught the htmodes bug: it fails if the
// fixture and the types disagree, regardless of what either says.
func TestCapturedRadioDecodes(t *testing.T) {
	const captured = `{
      "band": "2G", "channel": 6, "device": "radio0",
      "encryptions": ["none","psk2","psk-mixed","sae","sae-mixed"],
      "htmode": "auto",
      "htmodes": {"11b/g/n": 40, "11g/n": 40, "11n": 40, "auto": true},
      "hwmode": "11g/n",
      "hwmodes": ["11n","11g/n","11b/g/n"],
      "channels": [{"channel":1,"dfs":false},{"channel":6,"dfs":false}],
      "txpower": "Max",
      "ifaces": [{"enabled":true,"encryption":"psk2","guest":false,"hidden":false,
                  "key":"mockpass","name":"default_radio0","ssid":"GL-SFT1200-c41"}]
    }`

	var got WirelessRadio
	if err := json.Unmarshal([]byte(captured), &got); err != nil {
		t.Fatalf("the captured payload does not decode: %v", err)
	}

	if got.HWMode != "11g/n" {
		t.Errorf("hwmode = %q", got.HWMode)
	}
	// Slash-joined combinations, not the bare "11b"/"11g"/"11n" the description implies.
	if !reflect.DeepEqual(got.HWModes, []string{"11n", "11g/n", "11b/g/n"}) {
		t.Errorf("hwmodes = %v", got.HWModes)
	}
	// htmode is a width string, not "HT20".
	if got.HTMode != "auto" {
		t.Errorf("htmode = %q, want auto", got.HTMode)
	}
	if got.Channel != 6 {
		t.Errorf("channel = %d", got.Channel)
	}
	if len(got.Ifaces) != 1 || got.Ifaces[0].SSID != "GL-SFT1200-c41" {
		t.Errorf("ifaces = %+v", got.Ifaces)
	}
}

// The channel list comes from wifi.get_config, which reports what GL.iNet's firmware
// offers rather than what the radio supports. On the captured SFT1200 the API offers
// nine 5GHz channels while the driver reports twenty-five, the difference being every
// DFS channel. So the refusal must not claim the radio is incapable.
func TestChannelRefusalBlamesTheFirmwareNotTheRadio(t *testing.T) {
	r := &WirelessRadio{
		Band: Band5G,
		Channels: []WirelessChannel{
			{Channel: 36}, {Channel: 40}, {Channel: 149},
		},
	}
	channel := 52 // a DFS channel the hardware supports and the API withholds
	c := RadioChanges{Channel: &channel}

	err := c.Validate(r)
	if err == nil {
		t.Fatal("a channel outside the API's list was accepted")
	}
	if !strings.Contains(err.Error(), "firmware") {
		t.Errorf("refusal does not attribute the limit to the firmware: %v", err)
	}
	for _, wrong := range []string{"not available", "cannot", "does not support"} {
		if strings.Contains(err.Error(), wrong) {
			t.Errorf("refusal implies a hardware limit with %q: %v", wrong, err)
		}
	}
	// It still has to say what is on offer.
	if !strings.Contains(err.Error(), "149") {
		t.Errorf("refusal does not list the offered channels: %v", err)
	}
}
