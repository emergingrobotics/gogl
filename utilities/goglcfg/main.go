// Command goglcfg captures a GL.iNet router's reproducible configuration to a JSON
// profile, and applies one back.
//
// It is the fourth utility, with no gofi counterpart. The other three mirror gofi so
// that knowing one set means knowing the other; this one exists because reproducing a
// network spans all three of them and belongs in none.
//
// A profile is not a router image. It carries what defines a network -- the LAN and
// pool, reservations, DNS names and domain, wireless identity and tuning -- and omits
// everything identifying a particular unit. That is what makes it usable on a second
// router, which is the point.
//
// Every section comes from an endpoint verified against hardware. Nothing here is
// built on GL.iNet's API description alone, because doing that has been wrong three
// times in this project.
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
		fmt.Fprintln(os.Stderr, "goglcfg:", err)
		os.Exit(1)
	}
}

func run() error {
	var flags conn.Flags
	fs := flag.NewFlagSet("goglcfg", flag.ExitOnError)
	flags.Register(fs)

	get := fs.Bool("get", false, "capture a profile to stdout")
	fs.BoolVar(get, "g", false, "capture a profile to stdout (shorthand)")
	set := fs.Bool("set", false, "apply a profile from a file, or stdin")
	fs.BoolVar(set, "s", false, "apply a profile (shorthand)")

	withKeys := fs.Bool("with-keys", false,
		"with --get, include WiFi passphrases in cleartext")
	wireless := fs.Bool("wireless", false,
		"with --set, apply the wireless sections too; requires a wired session")
	dryRun := fs.Bool("dry-run", false, "with --set, show what would change without changing it")
	force := fs.Bool("force", false,
		"with --set, allow a subnet move while reservations exist")

	if err := conn.Parse(fs, os.Args[1:]); err != nil {
		return err
	}

	if *get == *set {
		fs.Usage()
		return errors.New("exactly one of --get or --set is required")
	}

	client, err := flags.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()

	if *get {
		profile, err := Capture(ctx, client, CaptureOptions{
			WithKeys: *withKeys,
			Host:     flags.Host,
			Captured: time.Now().Format(time.RFC3339),
		}, os.Stderr)
		if err != nil {
			return flags.Explain(err)
		}
		if !*withKeys && len(profile.Wireless) > 0 {
			fmt.Fprintln(os.Stderr,
				"note: WiFi passphrases omitted. Applying this profile leaves the target's\n"+
					"      existing passphrases alone. Use --with-keys to include them.")
		}
		return profile.Write(os.Stdout)
	}

	input := os.Stdin
	if path := fs.Arg(0); path != "" {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open %s: %w", path, err)
		}
		defer f.Close()
		input = f
	}

	profile, err := ReadProfile(input)
	if err != nil {
		return err
	}

	return flags.Explain(Apply(ctx, client, profile, applyModes{
		dryRun:   *dryRun,
		force:    *force,
		wireless: *wireless,
	}, os.Stderr))
}
