package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

const (
	noAddress     = "-"
	stateReserved = "reserved"
	stateDynamic  = "dynamic"
)

func FormatText(w io.Writer, entries []Entry, showReserved bool) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	for _, e := range entries {
		ip := e.IP
		if ip == "" {
			ip = noAddress
		}
		if showReserved {
			state := stateDynamic
			if e.Reserved {
				state = stateReserved
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", e.MAC, ip, e.Name, e.Manufacturer, state)
			continue
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", e.MAC, ip, e.Name, e.Manufacturer)
	}

	return tw.Flush()
}

func FormatJSON(w io.Writer, entries []Entry) error {
	// A nil slice must marshal as [] rather than null, so a consumer piping this
	// into jq never has to special-case the empty result.
	if entries == nil {
		entries = []Entry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
