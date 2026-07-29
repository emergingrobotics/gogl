package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	gogl "github.com/emergingrobotics/gogl/src"

	"github.com/emergingrobotics/gogl/src/types"
)

// Report is the flattened view goglnet prints. JSON tags define the -j output
// contract.
type Report struct {
	Model    string `json:"model,omitempty"`
	Firmware string `json:"firmware,omitempty"`

	LANIP   string `json:"lan_ip"`
	Netmask string `json:"netmask"`
	Subnet  string `json:"subnet,omitempty"`

	DHCPEnabled bool            `json:"dhcp_enabled"`
	DHCPStart   string          `json:"dhcp_start,omitempty"`
	DHCPStop    string          `json:"dhcp_stop,omitempty"`
	PoolSize    int             `json:"pool_size,omitempty"`
	DHCPLease   types.LeaseTime `json:"dhcp_lease,omitempty"`

	Interface string   `json:"interface,omitempty"`
	Gateway   string   `json:"gateway,omitempty"`
	DNS       []string `json:"dns,omitempty"`

	ReservedCount  int `json:"reserved_count"`
	AvailableCount int `json:"available_count"`
}

// The three narrow interfaces below exist so buildReport can be tested with stubs
// instead of a live client.
type networkGetter interface {
	Get(context.Context) (*types.Network, error)
}

type systemInfoer interface {
	Info(context.Context) (*types.SystemInfo, error)
}

type reservationLister interface {
	List(context.Context) ([]types.Reservation, error)
}

func buildReport(ctx context.Context, nets networkGetter, sys systemInfoer, res reservationLister) (*Report, error) {
	network, err := nets.Get(ctx)
	if err != nil {
		return nil, err
	}
	reservations, err := res.List(ctx)
	if err != nil {
		return nil, err
	}

	report := &Report{
		LANIP:         network.LANIP,
		Netmask:       network.Netmask,
		DHCPEnabled:   bool(network.DHCPEnabled),
		DHCPStart:     network.DHCPStart,
		DHCPStop:      network.DHCPStop,
		PoolSize:      network.PoolSize(),
		DHCPLease:     network.DHCPLease,
		Interface:     network.Interface,
		Gateway:       network.Gateway,
		DNS:           network.DNS,
		ReservedCount: len(reservations),
	}

	// System info is a convenience: a router that will not report its model is
	// still worth reporting the network of.
	if info, err := sys.Info(ctx); err == nil {
		report.Model, report.Firmware = info.Model, info.Firmware
	}

	if subnet, err := network.Subnet(); err == nil {
		report.Subnet = subnet.String()
	}

	report.AvailableCount = countAvailable(network, reservations)
	return report, nil
}

// countAvailable returns host addresses that are neither pooled, nor reserved,
// nor the router's own address.
//
// Reservations already inside the pool are not subtracted twice, which would
// understate what is actually free.
func countAvailable(network *types.Network, reservations []types.Reservation) int {
	available := network.UsableHosts()
	if available == 0 {
		return 0
	}

	available -= network.PoolSize()

	// The router's own address is never assignable, but only subtract it if the
	// pool has not already accounted for it.
	routerIP := net.ParseIP(network.LANIP)
	if routerIP != nil {
		if pooled, err := network.InDHCPPool(routerIP); err == nil && !pooled {
			available--
		}
	}

	for _, r := range reservations {
		ip := net.ParseIP(r.IP)
		if ip == nil {
			continue
		}
		inside, err := network.Contains(ip)
		if err != nil || !inside {
			continue
		}
		if pooled, err := network.InDHCPPool(ip); err == nil && pooled {
			continue
		}
		if routerIP != nil && ip.Equal(routerIP) {
			continue
		}
		available--
	}

	if available < 0 {
		return 0
	}
	return available
}

// runSetNetwork writes a new address and pool to the router.
//
// Two things make this unlike every other write in gogl. It is refused while
// reservations exist, because renumbering underneath them leaves each one pointing
// outside the new subnet. And it moves the router to a different address, so the
// call cannot report success in the usual way: the connection dies as the change
// takes effect. dryRun runs every check and prints the same summary, but stops
// short of the write.
func runSetNetwork(ctx context.Context, client *gogl.Client, n *types.Network, dryRun bool) error {
	current, err := client.Network().Get(ctx)
	if err != nil {
		return err
	}

	// The same validation the write performs, so a dry run cannot approve what
	// the real call would reject.
	if err := n.ValidateForWrite(); err != nil {
		return err
	}
	subnet, err := n.Subnet()
	if err != nil {
		return err
	}

	// Check the blocking condition before announcing the change. Printing "the
	// router will move" and then refusing reads as a failed write rather than a
	// refused one.
	existing, err := client.Reservations().List(ctx)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return fmt.Errorf("%w: %d reservation(s) present\n"+
			"  renumbering now would leave every one of them outside %s, present but inert\n"+
			"  clear them first:  goglps --clear",
			types.ErrReservationsExist, len(existing), subnet)
	}

	fmt.Fprintf(os.Stderr, "interface %s: %s/%s -> %s (%s)\n",
		n.Interface, current.LANIP, current.Netmask, n.LANIP, subnet)
	fmt.Fprintf(os.Stderr, "pool: %s-%s -> %s-%s\n",
		current.DHCPStart, current.DHCPStop, n.DHCPStart, n.DHCPStop)
	if dryRun {
		fmt.Fprintln(os.Stderr, "dry run: nothing was changed")
		return nil
	}

	fmt.Fprintln(os.Stderr,
		"the router will move to the new address; this session will drop")

	if err := client.Network().Set(ctx, n); err != nil {
		// A dropped connection here is the expected outcome of success, not a
		// failure, and saying so saves the operator a diagnosis.
		if errors.Is(err, context.DeadlineExceeded) || isConnectionLost(err) {
			fmt.Fprintln(os.Stderr,
				"connection lost, which is expected: the router has moved.")
			fmt.Fprintf(os.Stderr, "reconnect with -H %s\n", n.LANIP)
			return nil
		}
		return err
	}

	fmt.Printf("network set: %s on %s, pool %s-%s\n",
		n.LANIP, n.Interface, n.DHCPStart, n.DHCPStop)
	fmt.Fprintf(os.Stderr, "reconnect with -H %s\n", n.LANIP)
	return nil
}

// isConnectionLost reports whether err looks like the router going away mid-call,
// which is the normal consequence of a successful re-address.
func isConnectionLost(err error) bool {
	if err == nil {
		return false
	}
	text := err.Error()
	for _, sign := range []string{
		"connection refused", "connection reset", "EOF",
		"no route to host", "network is unreachable", "Client.Timeout",
	} {
		if strings.Contains(text, sign) {
			return true
		}
	}
	return false
}
