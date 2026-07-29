// Command goglps manages static IP reservations and their DNS names on a GL.iNet
// travel router, using ISC DHCP host declaration format.
//
// A reservation pins a MAC to an IP. It does NOT create a DNS record -- that was
// tested against a GL-SFT1200 on 4.3.28 and is false. The router's DNS answers
// from DHCP lease hostnames, which clients announce themselves, so the name in a
// reservation is a label: it identifies the entry here and in the admin panel.
//
// So this tool reproduces a network's addresses, not its names. gofips's
// --keep-dns still has no analogue, but for a duller reason than originally
// assumed: there is no DNS record to keep.
//
// The file format is identical to gofips's, so a file exported from a UniFi
// controller imports here without conversion. The one exception is hostnames
// containing an underscore, which is legal on UniFi but not a legal DNS label and
// is rejected rather than silently rewritten.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/emergingrobotics/gogl/utilities/internal/conn"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "goglps:", err)
		os.Exit(1)
	}
}

// modeFlags holds the mode selection and its modifiers.
type modeFlags struct {
	get   bool
	set   bool
	add   bool
	del   bool
	clear bool

	// domain sets the DNS domain, which reservation writes require.
	domain string

	name string
	mac  string
	ip   string

	force  bool
	prune  bool
	dryRun bool
}

func run() error {
	var (
		flags conn.Flags
		modes modeFlags
	)
	fs := flag.NewFlagSet("goglps", flag.ExitOnError)
	flags.Register(fs)

	fs.BoolVar(&modes.get, "get", false, "export all reservations in ISC DHCP format")
	fs.BoolVar(&modes.get, "g", false, "export all reservations (shorthand)")
	fs.BoolVar(&modes.set, "set", false, "import host declarations from a file or stdin")
	fs.BoolVar(&modes.set, "s", false, "import host declarations (shorthand)")
	fs.BoolVar(&modes.add, "add", false, "add a single host from a declaration fragment")
	fs.BoolVar(&modes.add, "a", false, "add a single host (shorthand)")
	fs.BoolVar(&modes.del, "del", false, "delete a host by --name, --mac, or --ip")
	fs.BoolVar(&modes.del, "d", false, "delete a host (shorthand)")
	fs.BoolVar(&modes.clear, "clear", false, "delete ALL reservations and DNS names")
	fs.StringVar(&modes.domain, "domain", "", "set the DNS domain (required before any reservation write)")

	fs.StringVar(&modes.name, "name", "", "hostname, with --del")
	fs.StringVar(&modes.name, "n", "", "hostname, with --del (shorthand)")
	fs.StringVar(&modes.mac, "mac", "", "MAC address, with --del")
	fs.StringVar(&modes.mac, "m", "", "MAC address, with --del (shorthand)")
	fs.StringVar(&modes.ip, "ip", "", "IP address, with --del")
	fs.StringVar(&modes.ip, "i", "", "IP address, with --del (shorthand)")

	fs.BoolVar(&modes.force, "force", false, "proceed past conflicts")
	fs.BoolVar(&modes.force, "f", false, "proceed past conflicts (shorthand)")
	fs.BoolVar(&modes.prune, "prune", false, "with --set, delete reservations and names absent from the file")
	fs.BoolVar(&modes.prune, "P", false, "with --set, delete reservations and names absent from the file (shorthand)")
	fs.BoolVar(&modes.dryRun, "dry-run", false, "show what would change without changing it")

	if err := conn.Parse(fs, os.Args[1:]); err != nil {
		return err
	}

	if err := checkModes(modes); err != nil {
		fs.Usage()
		return err
	}

	client, err := flags.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	date := time.Now().Format(time.DateOnly)

	switch {
	case modes.domain != "":
		return flags.Explain(runSetDomain(ctx, client, modes.domain))
	case modes.get:
		return flags.Explain(runGet(ctx, os.Stdout, client.Reservations(), client.Network(), date))
	case modes.set:
		return flags.Explain(runSet(ctx, client, fs.Arg(0), modes))
	case modes.add:
		return flags.Explain(runAdd(ctx, client, fs.Arg(0), modes))
	case modes.clear:
		return flags.Explain(runClear(ctx, client, modes))
	default:
		return flags.Explain(runDel(ctx, client, modes))
	}
}

// checkModes enforces exactly one mode flag.
func checkModes(modes modeFlags) error {
	selected := 0
	for _, on := range []bool{modes.get, modes.set, modes.add, modes.del, modes.clear, modes.domain != ""} {
		if on {
			selected++
		}
	}
	if selected != 1 {
		return errors.New("exactly one of --get, --set, --add, --del, --clear, --domain is required")
	}
	return nil
}
