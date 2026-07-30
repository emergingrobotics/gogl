package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/emergingrobotics/gogl/utilities/internal/clients"
	"github.com/emergingrobotics/gogl/utilities/internal/conn"
)

func newClientsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clients",
		Short: "Stations connected to the router, with manufacturer lookup",
		Long: `Stations connected to the router.

Its own area rather than part of lan, because a station arrives over cable, 2.4GHz or
5GHz and the useful view is all of them together. Manufacturer lookup is done from the
IEEE OUI registry independently of the router.`,
	}
	cmd.AddCommand(newClientsListCommand(), newClientsVendorCommand())
	return cmd
}

func newClientsListCommand() *cobra.Command {
	var wifi, wired, reserved bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List connected stations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			keep, err := clients.FilterFor(wifi, wired)
			if err != nil {
				return fmt.Errorf("%w: %s", errUsage, err)
			}

			client, err := connect()
			if err != nil {
				return err
			}
			defer client.Close()

			return explain(clients.List(cmd.Context(), client, os.Stdout, clients.Options{
				Keep:         keep,
				ShowReserved: reserved,
				JSON:         asJSON(),
			}))
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&wifi, "wifi", "w", false, "only wireless stations")
	f.BoolVarP(&wired, "wired", "e", false, "only wired stations")
	f.BoolVarP(&reserved, "reserved", "r", false, "mark which stations hold a reservation")
	return cmd
}

func newClientsVendorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "vendor <mac>",
		Short: "Look up a MAC address's manufacturer, without contacting the router",
		Long: `Look up a MAC address's manufacturer.

Entirely offline: it reads the cached IEEE OUI registry and never opens a session. Useful
for identifying a device from a MAC someone sent you.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cacheDir, err := conn.OUICacheDir()
			if err != nil {
				return err
			}
			db, err := clients.LoadOUI(cacheDir, nil)
			if err != nil {
				return err
			}

			vendor := db.Lookup(args[0])
			if vendor == "" {
				// A locally-administered or randomized address has no registry entry,
				// and saying "unknown" would imply a failed lookup rather than an
				// address that by design identifies nobody.
				fmt.Printf("%s  no registered manufacturer (locally administered or randomized)\n",
					strings.ToLower(args[0]))
				return nil
			}
			if asJSON() {
				return writeJSON(os.Stdout, map[string]string{
					"mac": strings.ToLower(args[0]), "vendor": vendor,
				})
			}
			fmt.Printf("%s  %s\n", strings.ToLower(args[0]), vendor)
			return nil
		},
	}
}
