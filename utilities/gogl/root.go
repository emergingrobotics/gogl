package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	gogl "github.com/emergingrobotics/gogl/src"
	"github.com/emergingrobotics/gogl/utilities/internal/config"
	"github.com/emergingrobotics/gogl/utilities/internal/conn"
)

// Global connection settings, shared by every command through persistent flags.
//
// Kept in one place rather than per-command because a router is not a per-command
// concern: `--router lab` should mean the same thing wherever it appears.
type globals struct {
	router   string
	host     string
	port     int
	https    bool
	insecure bool
	output   string

	// file is the loaded configuration, read once in PersistentPreRunE.
	file *config.File

	// flags carries the resolved connection settings. Populated from the config file
	// and then overridden by any flag the operator actually passed, which is why
	// resolution needs cobra's Changed() rather than comparing against defaults.
	flags conn.Flags
}

var opts globals

const (
	outputText = "text"
	outputJSON = "json"
)

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "gogl",
		Short: "Manage a GL.iNet travel router's network, DNS and wireless",
		Long: `gogl manages GL.iNet 4.x travel routers over their JSON-RPC API.

It exists to make a network reproducible: capture the addressing and names a site
depends on, then recreate them on a pocket router somewhere else.

Every write goes through the router's own API. There is no SSH and no shell, so the
blast radius of a mistake is bounded by what the API can express.`,

		// Errors are printed by main with a consistent prefix and exit code, and usage
		// on a runtime failure buries the message under a wall of flags.
		SilenceErrors: true,
		SilenceUsage:  true,

		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return loadConfig(cmd)
		},
	}

	flags := root.PersistentFlags()
	flags.StringVar(&opts.router, "router", "", "named router from the config file")
	flags.StringVarP(&opts.host, "host", "H", "", "router address")
	flags.IntVarP(&opts.port, "port", "p", config.DefaultPort, "router port")
	flags.BoolVar(&opts.https, "https", false, "use HTTPS instead of HTTP")
	flags.BoolVar(&opts.insecure, "insecure", true,
		"skip TLS certificate verification (default true: these devices ship self-signed certificates)")
	flags.StringVar(&opts.output, "output", outputText, "output format: text or json")

	root.AddCommand(
		newLANCommand(),
		newRadioCommand(),
		newWiFiCommand(),
		newClientsCommand(),
		newProfileCommand(),
		newSystemCommand(),
		newConfigCommand(),
	)
	return root
}

// loadConfig reads the configuration file and resolves the target router.
//
// Runs for every command including the offline ones. That is deliberate: a malformed
// config file should be reported by `gogl config show` rather than only surfacing when
// something tries to connect.
func loadConfig(cmd *cobra.Command) error {
	file, err := config.Load()
	if err != nil {
		return err
	}
	opts.file = file

	if opts.output != outputText && opts.output != outputJSON {
		return fmt.Errorf("--output %q: want %q or %q", opts.output, outputText, outputJSON)
	}
	// A file-level output preference applies only when the flag was not given.
	if !cmd.Flags().Changed("output") && file.Output != "" {
		opts.output = file.Output
	}

	router, err := file.Resolve(opts.router)
	if err != nil {
		return err
	}

	opts.flags = conn.Flags{
		Host:   opts.host,
		Port:   opts.port,
		HTTPS:  opts.https,
		Secure: !opts.insecure,
	}

	// The config file fills what the flags did not. Checked with Changed() rather than
	// against zero values, so that `--port 80` explicitly is distinguishable from not
	// passing --port at all.
	if router != nil {
		if !cmd.Flags().Changed("host") && opts.flags.Host == "" {
			opts.flags.Host = router.Host
		}
		if !cmd.Flags().Changed("port") && router.Port != 0 {
			opts.flags.Port = router.Port
		}
		if !cmd.Flags().Changed("https") && router.HTTPS {
			opts.flags.HTTPS = true
		}
		if !cmd.Flags().Changed("insecure") && router.Insecure {
			opts.flags.Secure = false
		}
		if router.Username != "" {
			opts.flags.Username = router.Username
		}
	}
	if opts.flags.Host == "" {
		opts.flags.Host = os.Getenv(config.EnvHost)
	}
	return nil
}

// connect opens a session to the resolved router.
//
// The password is resolved last and separately, because it is the one setting that
// never comes from the config file: environment, then a command the file names, then a
// prompt with echo off. There is no --password flag, since a secret on argv is visible
// in ps and recorded in shell history.
func connect() (*gogl.Client, error) {
	router, err := opts.file.Resolve(opts.router)
	if err != nil {
		return nil, err
	}

	password, err := router.Password(func() (string, error) {
		return conn.ReadSecret(fmt.Sprintf("password for %s: ", opts.flags.Host), "")
	})
	if err != nil {
		return nil, err
	}
	opts.flags.Password = password

	client, err := opts.flags.Connect()
	if err != nil {
		return nil, err
	}
	for _, w := range opts.flags.Warnings() {
		fmt.Fprintln(os.Stderr, "warning:", w)
	}
	return client, nil
}

// explain annotates an error with what gogl knows about the connection, which turns an
// opaque "access denied" into something actionable.
func explain(err error) error { return opts.flags.Explain(err) }

// asJSON reports whether output should be machine-readable.
func asJSON() bool { return opts.output == outputJSON }

// wantedRouterName is used by `config show` to report what would be connected to.
func wantedRouterName() string {
	if opts.router != "" {
		return opts.router
	}
	if opts.file != nil && opts.file.Default != "" {
		return opts.file.Default
	}
	return ""
}

// errRefused marks an error as a guard refusal rather than a failure, so main can exit
// with a distinct code and a script can tell "I was blocked" from "it broke".
var errRefused = errors.New("refused")
