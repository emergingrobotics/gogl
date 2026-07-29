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
	row("AVAILABLE", strconv.Itoa(r.AvailableCount))

	return tw.Flush()
}

func formatJSON(w io.Writer, r *Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
