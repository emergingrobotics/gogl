// Package config reads gogl's TOML configuration and resolves how to reach a router.
//
// The file holds everything except secrets. A password never appears in it: it comes
// from the environment, from a command the file names, or from a prompt. That split is
// the same one git makes with credential.helper and restic with --password-command,
// and it exists because a secret in a config file is a secret in every backup of that
// file.
package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// Environment variables, unchanged from the four-utility era so existing scripts keep
// working.
const (
	EnvPassword = "GL_PASSWORD"
	EnvUsername = "GL_USERNAME"
	EnvHost     = "GL_ROUTER_IP"
)

// DefaultPort is 80, which is what GL.iNet firmware serves. Not 443: a habit carried
// over from other vendors' controllers would fail here.
const DefaultPort = 80

// File is a parsed configuration file.
type File struct {
	// Default names the router used when --router is not given. Optional when the
	// file defines exactly one router.
	Default string `toml:"default"`

	// Output is "text" or "json".
	Output string `toml:"output"`

	Routers map[string]*Router `toml:"routers"`

	// path records where this was read from, for error messages and `config show`.
	path string
}

// Router is one device's connection settings.
type Router struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`

	// HTTPS and Insecure control TLS. Both default false, because these devices
	// serve plain HTTP on port 80.
	HTTPS    bool `toml:"https"`
	Insecure bool `toml:"insecure"`

	// Domain is the DNS suffix this router should use. Convenience only: it is
	// applied by an explicit command, never implicitly.
	Domain string `toml:"domain"`

	// PasswordCommand prints the router's password on stdout.
	//
	// It runs through a shell, so it can be a pipeline. That is deliberate -- "pass
	// show x" and "gpg -d x.gpg | head -1" are both things people have -- and it is
	// safe in a way that interpolating user data into a command is not: this string
	// comes from a file the operator wrote, not from an argument or a device.
	PasswordCommand string `toml:"password_command"`

	// name is the key this router was defined under.
	name string
}

// Name returns the key the router was defined under.
func (r *Router) Name() string { return r.name }

// Path returns the configuration file's location, honoring XDG_CONFIG_HOME.
func Path() string {
	if explicit := os.Getenv("GOGL_CONFIG"); explicit != "" {
		return explicit
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			// No home directory: return a relative path rather than failing here, so
			// the caller reports "no config" rather than a confusing lookup error.
			return filepath.Join(".config", "gogl", "config.toml")
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "gogl", "config.toml")
}

// CacheDir returns where re-downloadable data lives, honoring XDG_CACHE_HOME.
//
// Cache rather than data: the IEEE OUI database is fetched from the internet and can
// be thrown away at any time. It previously lived under XDG_DATA_HOME, which is for
// things a user would miss.
func CacheDir() (string, error) {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("locating a cache directory: %w", err)
		}
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(base, "gogl"), nil
}

// Load reads the configuration file.
//
// A missing file is not an error: gogl works entirely from flags and the environment,
// and demanding a config file to run one command would be hostile.
func Load() (*File, error) {
	return LoadFrom(Path())
}

// LoadFrom reads a specific path.
func LoadFrom(path string) (*File, error) {
	f := &File{Routers: map[string]*Router{}, path: path}

	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return f, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	if err := toml.Unmarshal(raw, f); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if f.Routers == nil {
		f.Routers = map[string]*Router{}
	}
	for name, r := range f.Routers {
		r.name = name
	}
	f.path = path

	if err := f.validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return f, nil
}

// validate rejects a file that would behave surprisingly, rather than letting the
// surprise happen at connection time.
func (f *File) validate() error {
	if f.Output != "" && f.Output != "text" && f.Output != "json" {
		return fmt.Errorf("output = %q, want \"text\" or \"json\"", f.Output)
	}
	if f.Default != "" {
		if _, ok := f.Routers[f.Default]; !ok {
			return fmt.Errorf("default = %q but no [routers.%s] section", f.Default, f.Default)
		}
	}
	for name, r := range f.Routers {
		if r.Host == "" {
			return fmt.Errorf("routers.%s has no host", name)
		}
	}
	return nil
}

// Path returns where this file was read from.
func (f *File) Path() string { return f.path }

// Names returns the configured router names, sorted.
func (f *File) Names() []string {
	names := make([]string, 0, len(f.Routers))
	for name := range f.Routers {
		names = append(names, name)
	}
	// Sorted so `config routers` output is stable and diffable.
	sort.Strings(names)
	return names
}

// Resolve picks a router by name, falling back to the default.
//
// With no name, no default, and exactly one router defined, that one is used: a
// single-router file should not need to name itself twice.
func (f *File) Resolve(name string) (*Router, error) {
	if name != "" {
		r, ok := f.Routers[name]
		if !ok {
			return nil, fmt.Errorf("no router named %q in %s (have: %s)",
				name, f.path, strings.Join(f.Names(), ", "))
		}
		return r, nil
	}

	if f.Default != "" {
		return f.Routers[f.Default], nil
	}
	if len(f.Routers) == 1 {
		for _, r := range f.Routers {
			return r, nil
		}
	}
	if len(f.Routers) > 1 {
		return nil, fmt.Errorf("several routers in %s and no default; pass --router (have: %s)",
			f.path, strings.Join(f.Names(), ", "))
	}
	return nil, nil
}

// Password resolves the router's password.
//
// Order, highest first: the environment, then the configured command, then a prompt.
// There is deliberately no --password flag: a secret on argv is visible in ps and
// lands in shell history.
//
// prompt is called only when nothing else supplies a password, and may be nil in
// non-interactive contexts.
func (r *Router) Password(prompt func() (string, error)) (string, error) {
	if v := os.Getenv(EnvPassword); v != "" {
		return v, nil
	}
	if r != nil && r.PasswordCommand != "" {
		out, err := runPasswordCommand(r.PasswordCommand)
		if err != nil {
			return "", err
		}
		if out != "" {
			return out, nil
		}
		return "", fmt.Errorf("password_command for %q produced no output: %s",
			r.name, r.PasswordCommand)
	}
	if prompt != nil {
		return prompt()
	}
	return "", fmt.Errorf("no password: set %s, configure password_command, or run interactively",
		EnvPassword)
}

// runPasswordCommand executes the command and returns its first line.
//
// First line rather than all output because `pass show` prints the password first and
// metadata after, which is the common case this has to work with.
func runPasswordCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	// Inherit stderr so a failing helper can explain itself -- a gpg passphrase
	// prompt, for instance, has nowhere else to go.
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("password_command failed (%s): %w", command, err)
	}

	line := string(out)
	if idx := strings.IndexByte(line, '\n'); idx >= 0 {
		line = line[:idx]
	}
	return strings.TrimRight(line, "\r"), nil
}
