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
//
// With --set-ssid it writes one wireless interface's SSID. That is refused when the
// session arrives over WiFi, since applying it would sever the session with no
// address to reconnect at, and it requires a confirmation or --yes.
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

	// Write flags. Moving the address requires the mask and both pool bounds, since
	// the old pool cannot be valid in a new subnet. A pool-only change needs neither:
	// the router is not moving, so the address and mask are read from the device.
	setIP := fs.String("set-ip", "", "new LAN address, e.g. 192.168.2.1 (requires --set-mask, --set-start, --set-end)")
	setMask := fs.String("set-mask", "", "new netmask")
	setStart := fs.String("set-start", "", "new DHCP pool start")
	setEnd := fs.String("set-end", "", "new DHCP pool end")
	setIface := fs.String("set-interface", types.InterfaceLAN,
		"which interface to write: lan or guest")
	dryRun := fs.Bool("dry-run", false, "with --set-*, show the change without applying it")
	force := fs.Bool("force", false,
		"with a subnet change, proceed despite existing reservations; the firmware rewrites them")

	// Wireless. Reported always; written only when one of the --set-* wireless flags
	// is given. These use optional types rather than plain flags because a partial
	// update has to tell "set it to false" apart from "leave it alone", which
	// flag.Bool cannot express.
	iface := fs.String("iface", "", "wireless interface to write, e.g. default_radio0")
	device := fs.String("device", "", "radio to tune, e.g. radio0")
	showKey := fs.Bool("show-key", false, "print WiFi passphrases instead of masking them")
	yes := fs.Bool("yes", false, "with a wireless write, skip the confirmation prompt")

	var setSSID, setKey, setEncryption optionalString
	var setHidden, setEnabled optionalBool
	fs.Var(&setSSID, "set-ssid", "new SSID (requires --iface)")
	fs.Var(&setKey, "set-key", "new WPA passphrase, 8 to 63 characters (requires --iface)")
	fs.Var(&setEncryption, "set-encryption", "new encryption mode, e.g. psk2 (requires --iface)")
	fs.Var(&setHidden, "set-hidden", "hide or advertise the SSID: true or false (requires --iface)")
	fs.Var(&setEnabled, "set-enabled", "enable or disable the interface: true or false (requires --iface)")

	var setChannel optionalInt
	var setHTMode, setHWMode, setTXPower optionalString
	fs.Var(&setChannel, "set-channel", "new channel, or 0 for auto (requires --device)")
	fs.Var(&setHTMode, "set-htmode", "new bandwidth, e.g. HT20 or VHT80 (requires --device)")
	fs.Var(&setHWMode, "set-hwmode", "new hardware mode, e.g. 11ac (requires --device)")
	fs.Var(&setTXPower, "set-txpower", "transmit power: Max, High, Medium or Low (requires --device)")

	if err := conn.Parse(fs, os.Args[1:]); err != nil {
		return err
	}

	ifaceChanges := types.InterfaceChanges{
		SSID:       setSSID.Ptr(),
		Key:        setKey.Ptr(),
		Encryption: setEncryption.Ptr(),
		Hidden:     setHidden.Ptr(),
		Enabled:    setEnabled.Ptr(),
	}
	radioChanges := types.RadioChanges{
		Channel: setChannel.Ptr(),
		HTMode:  setHTMode.Ptr(),
		HWMode:  setHWMode.Ptr(),
		TXPower: setTXPower.Ptr(),
	}

	// The firmware scopes interface fields by iface_name and radio fields by device,
	// so a write with neither cannot be routed. Say which one is missing.
	if !ifaceChanges.Empty() && *iface == "" {
		return errors.New("interface changes require --iface; run goglnet with no flags to list the interfaces")
	}
	if !radioChanges.Empty() && *device == "" {
		return errors.New("radio changes require --device; run goglnet with no flags to list the radios")
	}

	movingAddress := *setIP != "" || *setMask != ""
	changingPool := *setStart != "" || *setEnd != ""
	writing := movingAddress || changingPool

	if movingAddress {
		// A pool from the old subnet cannot be valid in a new one, and guessing a new
		// one would be worse than asking.
		if *setIP == "" || *setMask == "" || *setStart == "" || *setEnd == "" {
			return errors.New("moving the LAN address requires --set-ip, --set-mask, --set-start and --set-end together")
		}
	} else if changingPool {
		if *setStart == "" || *setEnd == "" {
			return errors.New("--set-start and --set-end are required together")
		}
	} else if *dryRun && ifaceChanges.Empty() && radioChanges.Empty() {
		// Silently reporting the network would suggest the dry run had something
		// to say about a write that was never requested.
		return errors.New("--dry-run applies to --set-* and --set-ssid; there is nothing to preview")
	}

	client, err := flags.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()

	if !ifaceChanges.Empty() || !radioChanges.Empty() {
		return flags.Explain(runSetWireless(ctx, client, *iface, ifaceChanges,
			*device, radioChanges, wirelessModes{dryRun: *dryRun, yes: *yes}))
	}

	if writing {
		return flags.Explain(runSetNetwork(ctx, client, &types.Network{
			Interface: *setIface,
			LANIP:     *setIP,
			Netmask:   *setMask,
			DHCPStart: *setStart,
			DHCPStop:  *setEnd,
		}, networkModes{dryRun: *dryRun, force: *force}))
	}

	report, err := buildReport(ctx, client.Network(), client.System(), client.Reservations())
	if err != nil {
		return flags.Explain(err)
	}

	if *asJSON {
		return formatJSON(os.Stdout, report)
	}
	if err := formatText(os.Stdout, report); err != nil {
		return err
	}

	// Wireless is reported alongside the wired network rather than behind a flag:
	// "what is this network" includes what devices associate to.
	radios, err := client.Wireless().Radios(ctx)
	if err != nil {
		// A router that will not report wireless is still worth reporting the LAN
		// for, so this is a warning rather than a failure.
		fmt.Fprintln(os.Stderr, "warning: could not read wireless config:", err)
		return nil
	}
	fmt.Fprintln(os.Stdout)
	return formatWireless(os.Stdout, radios, *showKey)
}
