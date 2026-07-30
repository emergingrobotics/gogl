#!/usr/bin/env bash
#
# Verify README.md against the built binaries.
#
# Exists because documentation drifted twice in ways review did not catch: goglnet
# grew --dry-run and --show-key with no table entry, and a README Go snippet kept
# calling Network().Set with two arguments after it took three. Both are mechanically
# checkable, so they should be checked mechanically.
# No pipefail. `printf | grep -q` is the natural spelling for these checks, and grep -q
# exits on its first match -- which closes the pipe, kills printf with SIGPIPE, and makes
# pipefail report the whole pipeline as failed even though the match succeeded. That
# produced false "undocumented flag" reports that varied with where in the file the match
# happened to be. The checks below use bash pattern matching instead of pipes, so there is
# nothing for pipefail to get wrong.
set -u

cd "$(dirname "$0")/.."

fail=0
note() { printf '%s\n' "$*" >&2; fail=1; }

# --- 1. README is narrative, not a flag reference. ----------------------------------
#
# This used to check README section by section against each binary's --help. That made sense
# when README *was* the reference; docs/gogl-guide.md is now, and README delegates to it.
# Checking README for exhaustive flag coverage would only pressure it back into duplicating
# the guide, which is how the two drift apart.
#
# README is still checked for broken links and anchors below, and the guide is checked for
# flag coverage next.
if ! grep -q 'docs/gogl-guide.md' README.md; then
    note "check-docs: README.md does not link to docs/gogl-guide.md"
fi

# --- 1b. The guide documents every flag the binary has. ----------------------------
#
# Checked separately from README because the guide is the complete reference: a flag missing
# there is a documentation bug, whereas README is allowed to be selective.
# Flags documented once under the guide's global-flags table rather than per command.
shared="host port https secure router output version json"

if [ -x bin/gogl ] && [ -f docs/gogl-guide.md ]; then
    guide=$(cat docs/gogl-guide.md)
    for leaf in "lan show" "lan set" "lan leases" \
                "lan reservations list" "lan reservations export" "lan reservations import" \
                "lan reservations add" "lan reservations rm" "lan reservations clear" \
                "lan dns show" "lan dns set" "lan dns add" "lan dns rm" "lan dns clear" \
                "radio list" "radio show" "radio set" \
                "wifi list" "wifi show" "wifi set" \
                "clients list" "clients vendor" \
                "profile export" "profile import" "system info" \
                "config show" "config routers" "config init"; do
        # shellcheck disable=SC2086
        body=$(bin/gogl $leaf --help 2>&1 | awk '/^Flags:/{f=1;next} /^Global Flags:/{f=0} f')
        for flag in $(printf '%s' "$body" | grep -oE '\-\-[a-z][a-z-]*' | sed 's/^--//' | sort -u); do
            case " $shared help " in *" $flag "*) continue ;; esac
            if [[ $guide != *'`--'"$flag"* ]]; then
                note "check-docs: docs/gogl-guide.md does not document \`gogl $leaf --$flag\`"
            fi
        done
    done
fi

# --- 2. In-page anchor links resolve. ----------------------------------------------
anchors=$(grep -oE '^#{1,6} .*' README.md \
    | sed -e 's/^#* //' -e 's/[`*_]//g' -e 's/[^[:alnum:] -]//g' \
    | tr '[:upper:] ' '[:lower:]-' | sort -u)

for link in $(grep -oE '\]\(#[^)]+\)' README.md | sed -e 's/](#//' -e 's/)//' | sort -u); do
    if [[ $'\n'"$anchors"$'\n' != *$'\n'"$link"$'\n'* ]]; then
        note "check-docs: broken anchor #$link"
    fi
done

# --- 3. Relative file links resolve. ----------------------------------------------
for doc in README.md VISION.md TODO.md docs/*.md; do
    [ -f "$doc" ] || continue
    dir=$(dirname "$doc")
    for target in $(grep -oE '\]\([^)#][^)]*\)' "$doc" | sed -e 's/](//' -e 's/)//' \
                    | grep -vE '^(https?|mailto):' | sed 's/#.*//' | sort -u); do
        [ -n "$target" ] || continue
        [ -e "$dir/$target" ] || note "check-docs: $doc links to missing $target"
    done
done

if [ "$fail" -eq 0 ]; then
    echo "check-docs: README and docs agree with the binaries."
fi
exit "$fail"
