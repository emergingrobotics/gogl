package conn

import (
	"fmt"
	"io"
	"runtime/debug"
)

// Version is stamped at build time with -ldflags "-X ...conn.Version=v1.2.3".
//
// It exists because two separate confusions this week came from a stale binary in
// ~/bin shadowing a fresh one in ./bin, each diagnosed by comparing `strings` output
// against source. One command should answer "which build am I running".
var Version = ""

// PrintVersion writes the build identity.
//
// Falls back to the module's VCS stamp, which the Go toolchain embeds automatically for
// builds from a repository. That covers `go install` and `go build` without -ldflags,
// which is how most people will get this.
func PrintVersion(w io.Writer, program string) {
	version := Version
	revision, modified := "", false

	if info, ok := debug.ReadBuildInfo(); ok {
		if version == "" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.modified":
				modified = setting.Value == "true"
			}
		}
	}
	if version == "" {
		version = "devel"
	}

	fmt.Fprintf(w, "%s %s", program, version)
	if revision != "" {
		short := revision
		if len(short) > 12 {
			short = short[:12]
		}
		fmt.Fprintf(w, " (%s", short)
		if modified {
			fmt.Fprint(w, ", dirty")
		}
		fmt.Fprint(w, ")")
	}
	fmt.Fprintln(w)
}
