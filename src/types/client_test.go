package types

import (
	"encoding/json"
	"testing"
)

func TestClientHostname(t *testing.T) {
	named := &Client{MAC: "aa:bb:cc:dd:ee:01", Name: "nas"}
	if got := named.Hostname(); got != "nas" {
		t.Errorf("Hostname() = %q, want nas", got)
	}

	nameless := &Client{MAC: "aa:bb:cc:dd:ee:02"}
	if got := nameless.Hostname(); got != UnknownHostname {
		t.Errorf("Hostname() = %q, want %q", got, UnknownHostname)
	}
}

// Wired/wireless comes from the iface string, not a boolean. Getting this wrong
// makes goglmac's --wired match nothing at all, which is exactly what an earlier
// guessed field set did against real hardware.
func TestClientIsWired(t *testing.T) {
	tests := []struct {
		iface string
		wired bool
		band  string
	}{
		{IfaceCable, true, ""},
		{Iface2GHz, false, "2.4G"},
		{Iface5GHz, false, "5G"},
		{"", false, ""},
	}
	for _, tt := range tests {
		c := &Client{Iface: tt.iface}
		if got := c.IsWired(); got != tt.wired {
			t.Errorf("Client{Iface:%q}.IsWired() = %v, want %v", tt.iface, got, tt.wired)
		}
		if got := c.Band(); got != tt.band {
			t.Errorf("Client{Iface:%q}.Band() = %q, want %q", tt.iface, got, tt.band)
		}
	}
}

// Decoded verbatim from clients.get_list on a GL-SFT1200 running 4.3.28, with the
// long rate arrays trimmed. If the firmware's field names ever drift, this fails
// rather than silently reporting every client as wireless.
func TestClientDecodesRealPayload(t *testing.T) {
	const payload = `{
      "limit_tx": 0,
      "ip": "192.168.8.185",
      "limit_rx": 0,
      "total_tx_init": 0,
      "total_tx": 4096,
      "last_update_rate": 1750905221,
      "type": 0,
      "mac": "B4:0E:CF:2A:85:6F",
      "remote": false,
      "iface": "2.4G",
      "tx": 0,
      "online": true,
      "name": "Bouffalolab_bl606p-2a856f",
      "blocked": false,
      "total_rx_init": 0,
      "total_rx": 8192,
      "rx": 0
    }`

	var c Client
	if err := json.Unmarshal([]byte(payload), &c); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if c.MAC != "B4:0E:CF:2A:85:6F" {
		t.Errorf("MAC = %q", c.MAC)
	}
	if c.IP != "192.168.8.185" {
		t.Errorf("IP = %q", c.IP)
	}
	if c.Name != "Bouffalolab_bl606p-2a856f" {
		t.Errorf("Name = %q", c.Name)
	}
	if !c.Online {
		t.Error("Online = false")
	}
	if c.Iface != Iface2GHz {
		t.Errorf("Iface = %q, want %q", c.Iface, Iface2GHz)
	}
	if c.IsWired() {
		t.Error("a 2.4G client reported as wired")
	}
	if c.Band() != "2.4G" {
		t.Errorf("Band() = %q", c.Band())
	}
	// Cumulative totals, not the instantaneous rates.
	if c.RXBytes != 8192 || c.TXBytes != 4096 {
		t.Errorf("counters = rx %d, tx %d; want 8192, 4096 from total_rx/total_tx", c.RXBytes, c.TXBytes)
	}
}

func TestClientDecodesWiredPayload(t *testing.T) {
	const payload = `{"mac":"08:26:AE:35:D7:A6","ip":"192.168.8.241","name":"helios","online":true,"iface":"cable","blocked":false}`

	var c Client
	if err := json.Unmarshal([]byte(payload), &c); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !c.IsWired() {
		t.Error("a cable client is not reported as wired")
	}
	if c.Band() != "" {
		t.Errorf("Band() = %q for a wired client, want empty", c.Band())
	}
}

// Optional fields must be omitted when absent, so a consumer can tell "not
// reported" from "reported as zero".
func TestClientOmitsEmptyFields(t *testing.T) {
	b, err := json.Marshal(&Client{MAC: "aa:bb:cc:dd:ee:01", Online: true})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	for _, absent := range []string{"ip", "name", "iface", "blocked", "total_rx", "total_tx"} {
		if _, present := decoded[absent]; present {
			t.Errorf("field %q should be omitted when empty", absent)
		}
	}
	for _, present := range []string{"mac", "online"} {
		if _, ok := decoded[present]; !ok {
			t.Errorf("field %q should always be present", present)
		}
	}
}
