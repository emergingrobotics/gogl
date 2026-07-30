package types

import (
	"encoding/json"
	"testing"
)

// The API description calls online_time a string. Firmware 4.3.28 sends a number. Both
// have to decode, and this is the test that was missing: it exercises the wire format
// rather than a Go struct marshalled back out.
func TestFlexStringDecodesNumberAndString(t *testing.T) {
	tests := map[string]struct {
		json string
		want string
	}{
		"number":       {`{"online_time": 1784548800}`, "1784548800"},
		"string":       {`{"online_time": "1784548800"}`, "1784548800"},
		"zero":         {`{"online_time": 0}`, "0"},
		"empty string": {`{"online_time": ""}`, ""},
		"null":         {`{"online_time": null}`, ""},
		"absent":       {`{}`, ""},
		"float":        {`{"online_time": 1784548800.0}`, "1784548800.0"},
		"quoted date":  {`{"online_time": "2026-07-30 11:00:00"}`, "2026-07-30 11:00:00"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var got struct {
				OnlineTime FlexString `json:"online_time"`
			}
			if err := json.Unmarshal([]byte(tt.json), &got); err != nil {
				t.Fatalf("Unmarshal(%s): %v", tt.json, err)
			}
			if got.OnlineTime.String() != tt.want {
				t.Errorf("OnlineTime = %q, want %q", got.OnlineTime, tt.want)
			}
		})
	}
}

// An object or array is a genuine type error, not something to swallow: silently
// accepting one would hide a real API change behind an empty value.
func TestFlexStringRejectsStructuredValues(t *testing.T) {
	for _, body := range []string{`{"online_time": {}}`, `{"online_time": [1]}`} {
		var got struct {
			OnlineTime FlexString `json:"online_time"`
		}
		if err := json.Unmarshal([]byte(body), &got); err == nil {
			t.Errorf("Unmarshal(%s) succeeded, want a type error", body)
		}
	}
}

// gogl's own output has one shape regardless of what the firmware sent, so a consumer of
// `--output json` never has to handle both.
func TestFlexStringMarshalsAsAString(t *testing.T) {
	out, err := json.Marshal(struct {
		OnlineTime FlexString `json:"online_time"`
	}{OnlineTime: "1784548800"})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := `{"online_time":"1784548800"}`; string(out) != want {
		t.Errorf("Marshal = %s, want %s", out, want)
	}
}

func TestFlexStringInt64(t *testing.T) {
	if n, ok := FlexString("1784548800").Int64(); !ok || n != 1784548800 {
		t.Errorf("Int64 = %d, %t", n, ok)
	}
	if _, ok := FlexString("2026-07-30").Int64(); ok {
		t.Error("a date parsed as an integer")
	}
	if _, ok := FlexString("").Int64(); ok {
		t.Error("an empty value parsed as an integer")
	}
}

// The whole client payload must decode, not just the field in isolation. This is the
// assertion whose absence let the bug ship: every clients test built types.Client values
// in Go and never decoded a captured response.
func TestClientDecodesCapturedPayload(t *testing.T) {
	const captured = `{"mac":"10:51:07:1f:8d:1c","ip":"192.168.8.10","name":"europa",
	  "iface":"cable","online":true,"online_time":1784548800,"blocked":false,
	  "type":2,"remote":false,"tx":0,"rx":0,"total_tx":102400,"total_rx":204800}`

	var c Client
	if err := json.Unmarshal([]byte(captured), &c); err != nil {
		t.Fatalf("a captured client does not decode: %v", err)
	}

	if c.MAC != "10:51:07:1f:8d:1c" || c.IP != "192.168.8.10" {
		t.Errorf("client = %+v", c)
	}
	if !c.Online {
		t.Error("online did not decode")
	}
	if c.OnlineTime.String() != "1784548800" {
		t.Errorf("online_time = %q", c.OnlineTime)
	}
	if !c.IsWired() {
		t.Error("a cable client did not report as wired")
	}
	if c.RXBytes != 204800 || c.TXBytes != 102400 {
		t.Errorf("byte totals = %d/%d", c.RXBytes, c.TXBytes)
	}
}
