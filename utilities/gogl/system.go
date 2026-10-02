package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSystemCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "Device identity",

		// Runnable + Args so an unknown subcommand here is a usage error (exit 2) rather
		// than cobra's silent help-with-exit-0 for a non-runnable parent.
		Args: wrapArgsError(unknownSubcommandArgs),
		RunE: showHelp,
	}
	cmd.AddCommand(newSystemInfoCommand())
	return cmd
}

func newSystemInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Report the model, firmware version and uptime",
		Args:  wrapArgsError(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := connect()
			if err != nil {
				return err
			}
			defer client.Close()

			info, err := client.System().Info(cmd.Context())
			if err != nil {
				return explain(err)
			}
			out := cmd.OutOrStdout()
			if asJSON() {
				return writeJSON(out, info)
			}
			fmt.Fprintf(out, "MODEL      %s\nFIRMWARE   %s\n", info.Model, info.Firmware)
			if info.MAC != "" {
				fmt.Fprintf(out, "MAC        %s\n", info.MAC)
			}
			fmt.Fprintf(out, "ENDPOINT   %s\n", client.Endpoint())
			return nil
		},
	}
}
