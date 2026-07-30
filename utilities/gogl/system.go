package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newSystemCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "Device identity",
	}
	cmd.AddCommand(newSystemInfoCommand())
	return cmd
}

func newSystemInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Report the model, firmware version and uptime",
		Args:  cobra.NoArgs,
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
			if asJSON() {
				return writeJSON(os.Stdout, info)
			}
			fmt.Printf("MODEL      %s\nFIRMWARE   %s\n", info.Model, info.Firmware)
			if info.MAC != "" {
				fmt.Printf("MAC        %s\n", info.MAC)
			}
			fmt.Printf("ENDPOINT   %s\n", client.Endpoint())
			return nil
		},
	}
}
