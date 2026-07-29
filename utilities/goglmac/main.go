// Command goglmac lists clients connected to a GL.iNet router, with independent
// IEEE OUI manufacturer lookup. Read-only.
//
// Its practical role in the workflow is discovering the MAC addresses to put into
// a host file for goglps.
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
		fmt.Fprintln(os.Stderr, "goglmac:", err)
		os.Exit(1)
	}
}

func run() error {
	var flags conn.Flags
	fs := flag.NewFlagSet("goglmac", flag.ExitOnError)
	flags.Register(fs)

	wifi := fs.Bool("wifi", false, "list only WiFi clients")
	fs.BoolVar(wifi, "w", false, "list only WiFi clients (shorthand)")
	wired := fs.Bool("wired", false, "list only wired clients")
	fs.BoolVar(wired, "e", false, "list only wired clients (shorthand)")
	all := fs.Bool("all", false, "list all clients (default)")
	fs.BoolVar(all, "a", false, "list all clients (shorthand)")

	asJSON := fs.Bool("json", false, "output JSON instead of text")
	fs.BoolVar(asJSON, "j", false, "output JSON instead of text (shorthand)")
	showReserved := fs.Bool("reserved", false, "mark which clients have a reservation")
	fs.BoolVar(showReserved, "r", false, "mark which clients have a reservation (shorthand)")

	if err := conn.Parse(fs, os.Args[1:]); err != nil {
		return err
	}

	if *wifi && *wired {
		return errors.New("--wifi and --wired are mutually exclusive")
	}
	keep := filterAll
	switch {
	case *wifi:
		keep = filterWiFi
	case *wired:
		keep = filterWired
	}

	// Load the OUI database before connecting: a stale-cache warning belongs
	// before any device output, and a hard OUI failure should not leave a session
	// open.
	cacheDir, err := conn.OUICacheDir()
	if err != nil {
		return err
	}
	db, err := LoadOUI(cacheDir, nil)
	if err != nil {
		return err
	}

	client, err := flags.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	clients, err := client.Clients().List(ctx)
	if err != nil {
		return flags.Explain(err)
	}

	// Only fetch reservations when they will be shown: it is an extra round trip
	// against a small SoC.
	var reservations []types.Reservation
	if *showReserved {
		reservations, err = client.Reservations().List(ctx)
		if err != nil {
			return flags.Explain(err)
		}
	}

	entries := buildEntries(clients, reservations, db, keep)

	if *asJSON {
		return formatJSON(os.Stdout, entries)
	}
	return formatText(os.Stdout, entries, *showReserved)
}
