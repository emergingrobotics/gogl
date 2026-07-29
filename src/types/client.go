package types

// Interface values the router reports in a client's Iface field. Confirmed
// against a GL-SFT1200 on firmware 4.3.28: wired clients report "cable", wireless
// ones report the band they are associated on.
const (
	IfaceCable = "cable"
	Iface2GHz  = "2.4G"
	Iface5GHz  = "5G"
)

// Client is a station currently known to the router.
//
// The field set and JSON tags mirror what firmware 4.3.28 actually returns from
// clients.get_list. Notably the firmware reports the connection as a single
// `iface` string rather than a boolean, and cumulative counters as total_rx and
// total_tx; an earlier version of this type guessed `is_wired`, `band`, `rx_bytes`
// and `signal`, none of which exist, so wired/wireless filtering silently matched
// nothing.
type Client struct {
	MAC    string `json:"mac"`
	IP     string `json:"ip,omitempty"`
	Name   string `json:"name,omitempty"`
	Online bool   `json:"online"`

	// Iface is the connection the client is on: "cable", "2.4G" or "5G".
	Iface string `json:"iface,omitempty"`

	// Blocked reports whether the router is blocking this client.
	Blocked bool `json:"blocked,omitempty"`

	// RXBytes and TXBytes are cumulative totals. The firmware also reports
	// instantaneous rx/tx rates, which this type deliberately omits: a byte total
	// is meaningful in a report, a momentary rate is not.
	RXBytes uint64 `json:"total_rx,omitempty"`
	TXBytes uint64 `json:"total_tx,omitempty"`
}

// UnknownHostname is reported for a client the router names nothing.
const UnknownHostname = "unknown"

// Hostname returns the name to display for c. Utilities call this rather than
// reimplementing the fallback, so the three of them agree.
func (c *Client) Hostname() string {
	if c.Name != "" {
		return c.Name
	}
	return UnknownHostname
}

// IsWired reports whether the client is on the wired LAN.
func (c *Client) IsWired() bool {
	return c.Iface == IfaceCable
}

// Band returns the radio band a wireless client is associated on, or the empty
// string for a wired client.
func (c *Client) Band() string {
	if c.Iface == "" || c.IsWired() {
		return ""
	}
	return c.Iface
}
