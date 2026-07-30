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
	stateOnline   = "yes"
	stateOffline  = "no"
)

// FormatOptions selects which columns appear.
//
// A struct rather than a run of positional booleans: three bools at a call site is
// unreadable, and this list will grow again.
type FormatOptions struct {
	// ShowReserved adds a column marking which clients hold a reservation.
	ShowReserved bool

	// ShowOnline adds a column reporting connected state. Set when offline clients may
	// be present; redundant when the caller has already filtered to online only.
	ShowOnline bool
}

func FormatText(w io.Writer, entries []Entry, opts FormatOptions) error {
	// A header with no rows under it is noise, and with the online-only default an
	// empty result needs explaining rather than leaving the reader to wonder whether
	// the command worked. JSON output stays machine-clean: it emits [].
	if len(entries) == 0 {
		if opts.ShowOnline {
			_, err := fmt.Fprintln(w, "no clients")
			return err
		}
		_, err := fmt.Fprintln(w, "no clients online; pass --all to include offline ones")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// A header, because without one a row of MAC, address, hostname and manufacturer is
	// four unlabelled columns and the reader has to infer which is which.
	header := "MAC\tADDRESS\tHOSTNAME\tMANUFACTURER\tSINCE"
	if opts.ShowOnline {
		header += "\tONLINE"
	}
	if opts.ShowReserved {
		header += "\tSTATE"
	}
	fmt.Fprintln(tw, header)

	for _, e := range entries {
		ip := e.IP
		if ip == "" {
			ip = noAddress
		}
		name := e.Name
		if name == "" {
			name = noAddress
		}
		since := e.Since
		if since == "" {
			since = noAddress
		}

		row := fmt.Sprintf("%s\t%s\t%s\t%s\t%s", e.MAC, ip, name, e.Manufacturer, since)
		if opts.ShowOnline {
			state := stateOffline
			if e.Online {
				state = stateOnline
			}
			row += "\t" + state
		}
		if opts.ShowReserved {
			state := stateDynamic
			if e.Reserved {
				state = stateReserved
			}
			row += "\t" + state
		}
		fmt.Fprintln(tw, row)
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
