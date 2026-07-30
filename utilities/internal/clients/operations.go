package clients

import (
	"net"
	"sort"
	"strings"

	"github.com/emergingrobotics/gogl/src/ipmath"
	"github.com/emergingrobotics/gogl/src/types"
)

// Entry is one row of goglmac output. JSON tags define the -j contract.
//
// The field set is narrower than gofimac's because the router reports less than a
// UniFi controller does. A field is added only when the firmware is observed
// returning it; nothing is invented for parity.
type Entry struct {
	MAC          string `json:"mac"`
	IP           string `json:"ip,omitempty"`
	Name         string `json:"hostname"`
	Manufacturer string `json:"manufacturer"`
	IsWired      bool   `json:"is_wired"`
	Online       bool   `json:"online"`
	Reserved     bool   `json:"reserved,omitempty"`
	Blocked      bool   `json:"blocked,omitempty"`
	RXBytes      uint64 `json:"rx_bytes,omitempty"`
	TXBytes      uint64 `json:"tx_bytes,omitempty"`
	Band         string `json:"band,omitempty"`
}

// Filter selects which clients to report.
type Filter func(types.Client) bool

func FilterAll(types.Client) bool     { return true }
func FilterWired(c types.Client) bool { return c.IsWired() }
func FilterWiFi(c types.Client) bool  { return !c.IsWired() }

// BuildEntries joins clients with the OUI database and the reservation table,
// sorted numerically by address with address-less clients last.
//
// The manufacturer always comes from our own OUI lookup, never from any value the
// router reports: on an OpenWrt 18.06 base the router's own table is likely to be
// years stale.
func BuildEntries(clients []types.Client, reservations []types.Reservation, db OUIDatabase, keep Filter) []Entry {
	reserved := make(map[string]bool, len(reservations))
	for _, r := range reservations {
		reserved[strings.ToLower(r.MAC)] = true
	}

	entries := make([]Entry, 0, len(clients))
	for i := range clients {
		c := clients[i]
		if !keep(c) {
			continue
		}
		mac := strings.ToLower(c.MAC)
		entries = append(entries, Entry{
			MAC:          mac,
			IP:           c.IP,
			Name:         c.Hostname(),
			Manufacturer: db.Lookup(c.MAC),
			IsWired:      c.IsWired(),
			Online:       c.Online,
			Reserved:     reserved[mac],
			Blocked:      c.Blocked,
			RXBytes:      c.RXBytes,
			TXBytes:      c.TXBytes,
			Band:         c.Band(),
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		left, right := net.ParseIP(entries[i].IP), net.ParseIP(entries[j].IP)
		// Clients without an address sort last: they are the least useful rows and
		// would otherwise all collide at position zero.
		switch {
		case left == nil && right == nil:
			return entries[i].MAC < entries[j].MAC
		case left == nil:
			return false
		case right == nil:
			return true
		}
		return ipmath.ToUint32(left) < ipmath.ToUint32(right)
	})

	return entries
}
