package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func write(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}

const example = `
default = "player-test"
output  = "json"

[routers.player-test]
host             = "192.168.8.1"
domain           = "herlein.me"
password_command = "printf secret"

[routers.lab]
host = "192.168.4.1"
port = 8080
`

func TestLoad(t *testing.T) {
	f, err := LoadFrom(write(t, example))
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}

	if f.Default != "player-test" || f.Output != "json" {
		t.Errorf("top level = %+v", f)
	}
	if got := f.Names(); !reflect.DeepEqual(got, []string{"lab", "player-test"}) {
		t.Errorf("Names = %v, want sorted", got)
	}

	r := f.Routers["lab"]
	if r.Host != "192.168.4.1" || r.Port != 8080 {
		t.Errorf("lab = %+v", r)
	}
	// The map key must reach the struct, or errors cannot name the router.
	if r.Name() != "lab" {
		t.Errorf("Name() = %q", r.Name())
	}
}

// Requiring a config file to run one command would be hostile: the tool works from
// flags and the environment alone.
func TestLoadMissingFileIsNotAnError(t *testing.T) {
	f, err := LoadFrom(filepath.Join(t.TempDir(), "absent.toml"))
	if err != nil {
		t.Fatalf("a missing file was an error: %v", err)
	}
	if len(f.Routers) != 0 {
		t.Errorf("routers = %+v, want empty", f.Routers)
	}
	r, err := f.Resolve("")
	if err != nil || r != nil {
		t.Errorf("Resolve on an empty file = %v, %v; want nil, nil", r, err)
	}
}

func TestResolve(t *testing.T) {
	f, err := LoadFrom(write(t, example))
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}

	byDefault, err := f.Resolve("")
	if err != nil || byDefault.Name() != "player-test" {
		t.Errorf("default = %v, %v", byDefault, err)
	}

	byName, err := f.Resolve("lab")
	if err != nil || byName.Name() != "lab" {
		t.Errorf("named = %v, %v", byName, err)
	}
}

// A typo in --router is the likeliest mistake, so the error has to list the real names.
func TestResolveUnknownNamesTheRealOnes(t *testing.T) {
	f, err := LoadFrom(write(t, example))
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}

	_, err = f.Resolve("playertest")
	if err == nil {
		t.Fatal("an unknown router was accepted")
	}
	for _, want := range []string{"player-test", "lab"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not list %q: %v", want, err)
		}
	}
}

// A single-router file should not have to name itself twice.
func TestResolveSingleRouterNeedsNoDefault(t *testing.T) {
	f, err := LoadFrom(write(t, "[routers.only]\nhost = \"192.168.8.1\"\n"))
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}
	r, err := f.Resolve("")
	if err != nil || r == nil || r.Name() != "only" {
		t.Errorf("Resolve = %v, %v", r, err)
	}
}

// Several routers and no default is ambiguous, and guessing would eventually write to
// the wrong device.
func TestResolveAmbiguousWithoutDefault(t *testing.T) {
	f, err := LoadFrom(write(t,
		"[routers.a]\nhost = \"192.168.8.1\"\n[routers.b]\nhost = \"192.168.4.1\"\n"))
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}
	_, err = f.Resolve("")
	if err == nil {
		t.Fatal("an ambiguous file resolved to something")
	}
	if !strings.Contains(err.Error(), "--router") {
		t.Errorf("error does not say how to disambiguate: %v", err)
	}
}

// Rejecting a bad file at load time beats letting the surprise happen at connection
// time, when a write may already be in flight.
func TestLoadRejectsBadFiles(t *testing.T) {
	tests := map[string]string{
		"router with no host":  "[routers.a]\ndomain=\"x.test\"\n",
		"default names naught": "default=\"ghost\"\n[routers.a]\nhost=\"1.2.3.4\"\n",
		"bad output":           "output=\"yaml\"\n[routers.a]\nhost=\"1.2.3.4\"\n",
		"not TOML":             "{[[[\n",
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := LoadFrom(write(t, body)); err == nil {
				t.Errorf("accepted %s", name)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Password resolution
// ---------------------------------------------------------------------------

func TestPasswordPrefersEnvironment(t *testing.T) {
	t.Setenv(EnvPassword, "from-env")
	r := &Router{PasswordCommand: "printf from-command"}

	got, err := r.Password(nil)
	if err != nil {
		t.Fatalf("Password: %v", err)
	}
	if got != "from-env" {
		t.Errorf("password = %q, want the environment to win", got)
	}
}

func TestPasswordFromCommand(t *testing.T) {
	t.Setenv(EnvPassword, "")
	r := &Router{name: "x", PasswordCommand: "printf from-command"}

	got, err := r.Password(nil)
	if err != nil {
		t.Fatalf("Password: %v", err)
	}
	if got != "from-command" {
		t.Errorf("password = %q", got)
	}
}

// `pass show` prints the password first and metadata after, which is the common case
// this has to work with.
func TestPasswordCommandTakesTheFirstLineOnly(t *testing.T) {
	t.Setenv(EnvPassword, "")
	r := &Router{name: "x", PasswordCommand: "printf 'secret\\nurl: example.test\\n'"}

	got, err := r.Password(nil)
	if err != nil {
		t.Fatalf("Password: %v", err)
	}
	if got != "secret" {
		t.Errorf("password = %q, want just the first line", got)
	}
}

func TestPasswordCommandHandlesCRLF(t *testing.T) {
	t.Setenv(EnvPassword, "")
	r := &Router{name: "x", PasswordCommand: "printf 'secret\\r\\n'"}

	got, err := r.Password(nil)
	if err != nil {
		t.Fatalf("Password: %v", err)
	}
	if got != "secret" {
		t.Errorf("password = %q, want the carriage return stripped", got)
	}
}

// A helper that fails must not be mistaken for "no password configured", which would
// send the operator looking in the wrong place.
func TestPasswordCommandFailureIsReported(t *testing.T) {
	t.Setenv(EnvPassword, "")
	r := &Router{name: "x", PasswordCommand: "exit 3"}

	_, err := r.Password(nil)
	if err == nil {
		t.Fatal("a failing password_command was accepted")
	}
	if !strings.Contains(err.Error(), "password_command") {
		t.Errorf("error does not name the cause: %v", err)
	}
}

// Silent success with empty output is the dangerous case: it would look like an empty
// password and produce an authentication failure that blames the router.
func TestPasswordCommandEmptyOutputIsAnError(t *testing.T) {
	t.Setenv(EnvPassword, "")
	r := &Router{name: "x", PasswordCommand: "true"}

	_, err := r.Password(nil)
	if err == nil {
		t.Fatal("a password_command producing nothing was accepted")
	}
	if !strings.Contains(err.Error(), "no output") {
		t.Errorf("error does not say what happened: %v", err)
	}
}

func TestPasswordFallsBackToPrompt(t *testing.T) {
	t.Setenv(EnvPassword, "")
	r := &Router{name: "x"}

	got, err := r.Password(func() (string, error) { return "typed", nil })
	if err != nil {
		t.Fatalf("Password: %v", err)
	}
	if got != "typed" {
		t.Errorf("password = %q", got)
	}
}

// Non-interactive with nothing configured must explain all three options rather than
// just failing.
func TestPasswordWithNothingAvailable(t *testing.T) {
	t.Setenv(EnvPassword, "")
	r := &Router{name: "x"}

	_, err := r.Password(nil)
	if err == nil {
		t.Fatal("Password succeeded with nothing configured")
	}
	for _, want := range []string{EnvPassword, "password_command", "interactively"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not mention %q: %v", want, err)
		}
	}
}

// A nil Router still has to work: gogl runs from flags alone with no config file.
func TestPasswordOnNilRouter(t *testing.T) {
	t.Setenv(EnvPassword, "from-env")
	var r *Router

	got, err := r.Password(nil)
	if err != nil || got != "from-env" {
		t.Errorf("Password on a nil router = %q, %v", got, err)
	}
}

// ---------------------------------------------------------------------------
// Paths
// ---------------------------------------------------------------------------

func TestPathHonorsXDG(t *testing.T) {
	t.Setenv("GOGL_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "/xdg")

	if got, want := Path(), filepath.Join("/xdg", "gogl", "config.toml"); got != want {
		t.Errorf("Path = %q, want %q", got, want)
	}
}

func TestPathExplicitOverride(t *testing.T) {
	t.Setenv("GOGL_CONFIG", "/tmp/elsewhere.toml")
	t.Setenv("XDG_CONFIG_HOME", "/xdg")

	if got := Path(); got != "/tmp/elsewhere.toml" {
		t.Errorf("Path = %q, want the explicit override", got)
	}
}

func TestPathFallsBackToHome(t *testing.T) {
	t.Setenv("GOGL_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "/home/someone")

	if got, want := Path(), "/home/someone/.config/gogl/config.toml"; got != want {
		t.Errorf("Path = %q, want %q", got, want)
	}
}

// Cache, not data: the OUI database is re-downloadable and a user would not miss it.
func TestCacheDirHonorsXDG(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "/xdgcache")

	got, err := CacheDir()
	if err != nil {
		t.Fatalf("CacheDir: %v", err)
	}
	if want := filepath.Join("/xdgcache", "gogl"); got != want {
		t.Errorf("CacheDir = %q, want %q", got, want)
	}
}

func TestCacheDirFallsBackToHome(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("HOME", "/home/someone")

	got, err := CacheDir()
	if err != nil {
		t.Fatalf("CacheDir: %v", err)
	}
	if want := "/home/someone/.cache/gogl"; got != want {
		t.Errorf("CacheDir = %q, want %q", got, want)
	}
}

// The file is read with os.ReadFile, so a directory in its place is an error rather
// than a panic.
func TestLoadFromDirectory(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadFrom(dir)
	if err == nil {
		t.Error("a directory was accepted as a config file")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Error("a directory was reported as not existing")
	}
}
