#!/usr/bin/env bash
#
# Verify README.md against the built binaries.
#
# Exists because documentation drifted twice in ways review did not catch: goglnet
# grew --dry-run and --show-key with no table entry, and a README Go snippet kept
# calling Network().Set with two arguments after it took three. Both are mechanically
# checkable, so they should be checked mechanically.
set -uo pipefail

cd "$(dirname "$0")/.."

fail=0
note() { printf '%s\n' "$*" >&2; fail=1; }

# --- 1. Every flag in --help is documented, and every documented flag exists. -------
#
# Connection flags are documented once in Configuration rather than per tool, so they
# are exempt from the per-section check.
# Flags every tool shares, documented once under "Common flags" and "Connection flags"
# rather than repeated in each tool's section.
shared="host port https secure router output version json"

if ! compgen -G 'bin/*' >/dev/null; then
    echo "check-docs: no binaries in bin/; skipping the flag check." >&2
fi

for tool in bin/*; do
    name=$(basename "$tool")
    [ -x "$tool" ] || continue
    [ -f "$tool" ] || continue

    section=$(awk -v t="### \`$name\`" '
        index($0, t) == 1 { inside = 1; next }
        inside && /^---$/  { exit }
        inside             { print }
    ' README.md)

    if [ -z "$section" ]; then
        note "check-docs: README.md has no '### \`$name\`' section"
        continue
    fi

    # Strip only the leading whitespace and dashes. Stripping every dash turns
    # --dry-run into dryrun and compares nothing useful.
    real=$("$tool" --help 2>&1 \
        | grep -oE '^[[:space:]]+-{1,2}[a-z][a-z-]*' \
        | sed -E 's/^[[:space:]]+-{1,2}//' | sort -u)

    # Prose contains wildcards like --set-* and cross-references to other tools'
    # flags. A trailing hyphen marks the former; drop those rather than reporting a
    # flag named "set-".
    documented=$(printf '%s' "$section" | grep -oE '`--[a-z][a-z-]*' \
        | sed 's/`--//' | grep -vE -- '-$' | sort -u)

    for flag in $real; do
        case " $shared " in *" $flag "*) continue ;; esac
        [ ${#flag} -lt 2 ] && continue
        if ! printf '%s\n' "$documented" | grep -qx "$flag"; then
            note "check-docs: $name --$flag exists but is not documented"
        fi
    done

    for flag in $documented; do
        case " $shared " in *" $flag "*) continue ;; esac
        if ! printf '%s\n' "$real" | grep -qx "$flag"; then
            note "check-docs: $name --$flag is documented but does not exist"
        fi
    done
done

# --- 2. In-page anchor links resolve. ----------------------------------------------
anchors=$(grep -oE '^#{1,6} .*' README.md \
    | sed -e 's/^#* //' -e 's/[`*_]//g' -e 's/[^[:alnum:] -]//g' \
    | tr '[:upper:] ' '[:lower:]-' | sort -u)

for link in $(grep -oE '\]\(#[^)]+\)' README.md | sed -e 's/](#//' -e 's/)//' | sort -u); do
    printf '%s\n' "$anchors" | grep -qx "$link" || note "check-docs: broken anchor #$link"
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
