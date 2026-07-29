package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	disabledMarker = "(disabled)"
	emptyMarker    = "-"
)

func formatText(w io.Writer, r *Report) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	row := func(label, value string) {
		if value == "" {
			value = emptyMarker
		}
		fmt.Fprintf(tw, "%s\t%s\n", label, value)
	}

	row("MODEL", r.Model)
	row("FIRMWARE", r.Firmware)

	lan := r.Subnet
	if lan == "" {
		lan = r.LANIP
	}
	row("LAN", fmt.Sprintf("%s  (%s)", lan, r.Netmask))

	if r.DHCPEnabled {
		row("DHCP", "enabled")
		row("POOL", fmt.Sprintf("%s - %s  (%d addresses)", r.DHCPStart, r.DHCPStop, r.PoolSize))
		row("LEASE", r.DHCPLease.String())
	} else {
		row("DHCP", "disabled")
		row("POOL", disabledMarker)
		row("LEASE", disabledMarker)
	}

	row("INTERFACE", r.Interface)
	row("GATEWAY", r.Gateway)
	row("DNS", strings.Join(r.DNS, ","))
	row("RESERVED", strconv.Itoa(r.ReservedCount))
	if len(r.InPool) > 0 {
		// Spelled out on the RESERVED line's heels, because it is the explanation for
		// an AVAILABLE count that otherwise looks wrong: a reservation inside the pool
		// is neither dynamically assignable nor counted as available.
		row("  IN POOL", fmt.Sprintf("%d  (honored by dnsmasq, excluded from the pool)", len(r.InPool)))
	}
	row("AVAILABLE", strconv.Itoa(r.AvailableCount))

	if err := tw.Flush(); err != nil {
		return err
	}

	if len(r.InPool) > 0 {
		fmt.Fprintf(w, "\n%d reservation(s) fall inside the DHCP pool %s-%s:\n",
			len(r.InPool), r.DHCPStart, r.DHCPStop)
		for _, res := range r.InPool {
			name := res.Name
			if name == "" {
				name = emptyMarker
			}
			fmt.Fprintf(w, "  %-15s %-17s %s\n", res.IP, res.MAC, name)
		}
		fmt.Fprintln(w, "These work: dnsmasq honors a static bind inside the dynamic range and")
		fmt.Fprintln(w, "excludes that address from allocation. Move the pool with --set-start")
		fmt.Fprintln(w, "and --set-end if you would rather keep the ranges separate.")
	}

	return nil
}

func formatJSON(w io.Writer, r *Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
