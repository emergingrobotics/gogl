// Command goglnet reports a GL.iNet router's LAN address, DHCP pool, and DNS
// settings. Read-only: it never writes to the router.
//
// It is the companion to goglps in the way gofinet is the companion to gofips.
// The pool boundaries define which addresses the router hands out dynamically;
// everything else in the subnet is available for a static reservation.
//
// Note the DNS line reports the resolvers the router advertises to clients, not
// any per-host names: a reservation does not create a DNS record on this firmware.
// Use goglps --domain and the host file for names.
//
// With the --set-* flags it also writes the LAN address and DHCP pool. That is
// refused while any reservation exists, and it drops the session, since the router
// moves to a new address mid-call. --dry-run shows the change and the refusal
// without performing the write.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/emergingrobotics/gogl/src/types"
	"github.com/emergingrobotics/gogl/utilities/internal/conn"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "goglnet:", err)
		os.Exit(1)
	}
}

func run() error {
	var flags conn.Flags
	fs := flag.NewFlagSet("goglnet", flag.ExitOnError)
	flags.Register(fs)

	asJSON := fs.Bool("json", false, "output JSON instead of text")
	fs.BoolVar(asJSON, "j", false, "output JSON instead of text (shorthand)")

	// Write flags. All four are required together: a partial network is not a
	// thing the firmware can apply, and guessing the rest would be worse than
	// asking.
	setIP := fs.String("set-ip", "", "new LAN address, e.g. 192.168.2.1 (requires --set-mask, --set-start, --set-end)")
	setMask := fs.String("set-mask", "", "new netmask")
	setStart := fs.String("set-start", "", "new DHCP pool start")
	setEnd := fs.String("set-end", "", "new DHCP pool end")
	setIface := fs.String("set-interface", types.InterfaceLAN,
		"which interface to write: lan or guest")
	dryRun := fs.Bool("dry-run", false, "with --set-*, show the change without applying it")

	if err := conn.Parse(fs, os.Args[1:]); err != nil {
		return err
	}

	writing := *setIP != "" || *setMask != "" || *setStart != "" || *setEnd != ""
	if writing {
		if *setIP == "" || *setMask == "" || *setStart == "" || *setEnd == "" {
			return errors.New("--set-ip, --set-mask, --set-start and --set-end are all required together")
		}
	} else if *dryRun {
		// Silently reporting the network would suggest the dry run had something
		// to say about a write that was never requested.
		return errors.New("--dry-run applies to --set-*; there is nothing to preview")
	}

	client, err := flags.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()

	if writing {
		return flags.Explain(runSetNetwork(ctx, client, &types.Network{
			Interface: *setIface,
			LANIP:     *setIP,
			Netmask:   *setMask,
			DHCPStart: *setStart,
			DHCPStop:  *setEnd,
		}, *dryRun))
	}

	report, err := buildReport(ctx, client.Network(), client.System(), client.Reservations())
	if err != nil {
		return flags.Explain(err)
	}

	if *asJSON {
		return formatJSON(os.Stdout, report)
	}
	return formatText(os.Stdout, report)
}
