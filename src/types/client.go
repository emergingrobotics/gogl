package types

import (
	"strconv"
	"strings"
	"time"
)

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

	// OnlineTime is when the client came online, verbatim as the firmware reports it.
	//
	// A string because that is what the API description calls it, and the format is
	// uncaptured: it may be a unix timestamp, or seconds elapsed, or a formatted date.
	// GL.iNet is not consistent about this -- network.get_dhcp_leases calls its field
	// "expires" and reports seconds remaining, not a timestamp -- so the value is kept
	// unparsed and rendered by SinceOnline, which says what it does not know.
	OnlineTime string `json:"online_time,omitempty"`

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

// SinceOnline renders OnlineTime for a human.
//
// The format is uncaptured, so this recognises the two plausible integer readings and
// otherwise returns the value untouched. Printing a guess as though it were known is how
// a lease "expires" field got rendered as a 1970 date earlier in this project.
//
// The second return reports whether the value was understood, so a caller can choose to
// show nothing rather than something meaningless.
func (c *Client) SinceOnline(now time.Time) (string, bool) {
	raw := strings.TrimSpace(c.OnlineTime)
	if raw == "" {
		return "", false
	}

	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		// Not an integer: possibly a formatted date. Pass it through rather than
		// discard information we cannot classify.
		return raw, true
	}

	switch {
	case n <= 0:
		return "", false
	case n >= unixTimestampFloor:
		// Large enough to only make sense as seconds since the epoch.
		return now.Sub(time.Unix(n, 0)).Truncate(time.Second).String(), true
	default:
		// Small enough to be an elapsed count rather than an absolute time.
		return (time.Duration(n) * time.Second).String(), true
	}
}

// unixTimestampFloor is 2001-09-09, the point past which an integer is far more likely
// to be a unix timestamp than a count of elapsed seconds. A device claiming 31 years of
// uptime is not the reading to prefer.
const unixTimestampFloor = 1_000_000_000
