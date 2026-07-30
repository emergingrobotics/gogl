package conn

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

// ReadSecret obtains a secret without putting it on the command line.
//
// This exists to fix a real defect. `goglnet --set-key 'passphrase'` placed a WiFi
// passphrase in argv, where it is visible to any other user via ps and is recorded in
// shell history -- contradicting the reason the router password was never given a flag
// in the first place. The same tool enforced one rule and broke it one flag over.
//
// Resolution order: the named command, then a prompt on the terminal. There is
// deliberately no path that accepts the value as an argument.
func ReadSecret(prompt, command string) (string, error) {
	if command != "" {
		return runSecretCommand(command)
	}
	return promptSecret(prompt)
}

// promptSecret reads from the terminal with echo disabled.
//
// Reads from /dev/tty rather than stdin so that a secret can still be typed while
// stdin carries data -- `gogl lan reservations import - ` piping a host file, for
// instance.
func promptSecret(prompt string) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("no terminal to prompt on: %w", err)
	}
	defer tty.Close()

	fd := int(tty.Fd())
	if !term.IsTerminal(fd) {
		return "", errors.New("not a terminal; use the --*-command form instead")
	}

	fmt.Fprint(tty, prompt)
	raw, err := term.ReadPassword(fd)
	fmt.Fprintln(tty)
	if err != nil {
		return "", fmt.Errorf("reading the secret: %w", err)
	}

	secret := strings.TrimRight(string(raw), "\r\n")
	if secret == "" {
		return "", errors.New("empty secret")
	}
	return secret, nil
}

// runSecretCommand executes command and returns its first line.
//
// First line only, because `pass show` prints the secret first and metadata after.
func runSecretCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("secret command failed (%s): %w", command, err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	if !scanner.Scan() {
		return "", fmt.Errorf("secret command produced no output: %s", command)
	}
	secret := strings.TrimRight(scanner.Text(), "\r")
	if secret == "" {
		return "", fmt.Errorf("secret command produced an empty first line: %s", command)
	}
	return secret, nil
}
