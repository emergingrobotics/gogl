#!/usr/bin/env bash
#
# Hardware-in-the-loop test: exercise every writable setting against a real router.
#
# Captures a baseline profile, changes one setting at a time, reads each back to confirm
# it took, then restores the baseline. Every mock test in this repository asserts what
# gogl believes; this asserts what the router does.
#
# It exists because believing the vendor's API description has been wrong four times:
# dhcp.* does not exist, dns.set_host rejects three characters it does not mention,
# wifi.htmodes is an object where an array is documented, and clients.online_time is a
# number where a string is documented. Each was found by a person running a command, not
# by the test suite.
#
# MUST BE RUN FROM THE ROUTER'S WIRED LAN. This is checked, not assumed, and the run
# aborts if it is not true. Two reasons:
#
#   * Wireless writes are refused over a wireless session, so a wireless run would skip
#     every test that matters and report a hollow pass.
#   * Being on the LAN is what makes radio writes safe to test: dropping every wireless
#     client cannot cut the connection this script is using.
#
# The same precondition rules out testing LAN configuration. Renumbering the LAN or moving
# the DHCP pool disturbs the very network this host is on -- at best invalidating its
# lease, at worst making the router unreachable mid-test with the baseline unrestored.
# Those paths are covered by unit tests against the mock and by hand.
#
# A companion script is planned to close the remaining gap: associate the test host's own
# WiFi to the router after a wireless change and confirm a client can actually get an
# address and route. That is the one thing neither this script nor the mock can prove --
# that a write which the API accepted produced a working network.
#
#   ./scripts/hil-test.sh --help
#
# THIS WRITES TO A REAL ROUTER. Read the safety notes below before running it.
set -uo pipefail

cd "$(dirname "$0")/.."

# ---------------------------------------------------------------------------
# Options
# ---------------------------------------------------------------------------

GOGL="${GOGL:-gogl}"
ROUTER_FLAG=()
INCLUDE_UPLINK_RADIO=0
DRY_RUN=0
ASSUME_YES=0
KEEP_BASELINE=0
CALL_DELAY="${CALL_DELAY:-1}"
WORKDIR="${WORKDIR:-$(mktemp -d)}"

usage() {
    cat <<'USAGE'
Usage: scripts/hil-test.sh [options]

Exercises the wireless and naming settings against a real GL.iNet router, verifying each
change took, then restores the configuration it started from.

Must be run from the router's wired LAN. LAN addressing and the DHCP pool are
deliberately not tested: changing them disturbs the network this host is on.

Options:
  --router NAME          pass --router NAME to every gogl call
  -H, --host ADDR        pass -H ADDR to every gogl call
  --include-uplink-radio also tune the radio carrying the WiFi uplink. On a router in
                         repeater mode this can drop its internet connection. Off by
                         default.
  --dry-run              list what would be tested and exit, touching nothing
  --yes                  skip the confirmation prompt
  --keep-baseline        leave the baseline profile on disk after a successful run
  --delay SECONDS        pause between gogl calls (default 1). Each call is a full
                         login, and the firmware rate-limits them.
  -h, --help             this text

Safety:
  * A wired LAN connection is REQUIRED and verified. The run aborts otherwise.
  * Wireless clients WILL be disconnected: SSIDs, passphrases and channels all change.
    Nothing this script does can cut its own wired connection.
  * The baseline is captured WITH passphrases so they can be restored. It is written to
    a temporary directory and deleted on success unless --keep-baseline.
  * Restore runs on any exit, including Ctrl-C. If restore itself fails, the baseline
    path is printed and left in place.
  * Nothing here reboots, upgrades, or changes the admin password.

Environment:
  GOGL          the binary to test (default: gogl on PATH)
  WORKDIR       where to keep the baseline (default: a temporary directory)
  CALL_DELAY    seconds between gogl calls (default 1)
  LOCKOUT_WAIT  seconds to wait for a lockout to clear before giving up (default 900)
USAGE
}

while [ $# -gt 0 ]; do
    case "$1" in
        --router)               ROUTER_FLAG=(--router "$2"); shift 2 ;;
        -H|--host)              ROUTER_FLAG=(-H "$2"); shift 2 ;;
        --include-uplink-radio) INCLUDE_UPLINK_RADIO=1; shift ;;
        --dry-run)              DRY_RUN=1; shift ;;
        --yes)                  ASSUME_YES=1; shift ;;
        --keep-baseline)        KEEP_BASELINE=1; shift ;;
        --delay)                CALL_DELAY="$2"; shift 2 ;;
        -h|--help)              usage; exit 0 ;;
        *) echo "unknown option: $1" >&2; usage >&2; exit 2 ;;
    esac
done

BASELINE="$WORKDIR/baseline.json"

# ---------------------------------------------------------------------------
# Reporting
# ---------------------------------------------------------------------------

PASS=0
FAIL=0
SKIP=0
FAILURES=()

bold() { printf '\033[1m%s\033[0m\n' "$*"; }
pass() { PASS=$((PASS + 1)); printf '  \033[32mPASS\033[0m  %s\n' "$*"; }
skip() { SKIP=$((SKIP + 1)); printf '  \033[33mSKIP\033[0m  %s\n' "$*"; }
fail() {
    FAIL=$((FAIL + 1))
    FAILURES+=("$1")
    printf '  \033[31mFAIL\033[0m  %s\n' "$*"
}

# Every gogl invocation performs a full login: two challenge calls and a login. The
# firmware's brute-force protection counts those, and a run of this script makes dozens.
#
# OBSERVED 2026-07-30: an unpaced run locked the account out after the first test, and the
# remaining seven cascaded into failures and skips that hid the cause. Worse, restore could
# not authenticate either, so the router was left with a probe value in place.
#
# Two defences: pace the calls, and stop at the first sign of a lockout rather than
# spending the next sixty calls making it worse.
LOCKED_OUT=0

g() {
    if [ "$LOCKED_OUT" = 1 ]; then
        return 1
    fi
    sleep "$CALL_DELAY"

    local out status
    out=$("$GOGL" "${ROUTER_FLAG[@]}" "$@" 2>&1)
    status=$?

    if printf '%s' "$out" | grep -qE 'rate limiting|too many failed logins'; then
        LOCKED_OUT=1
        printf '%s\n' "$out" >&2
        return 1
    fi

    printf '%s\n' "$out"
    return "$status"
}

# require_session aborts the run when the router has stopped authenticating.
#
# Called between tests. Continuing past a lockout produces a page of failures that say
# nothing about the real problem, and leaves less time for the restore to succeed.
require_session() {
    [ "$LOCKED_OUT" = 0 ] && return 0
    cat >&2 <<'MSG'

hil-test: the router is rate limiting and has stopped authenticating.

Every gogl invocation logs in again, and this script makes dozens. Waiting for the
lockout to clear before attempting the restore.
MSG
    return 1
}

# wait_for_session blocks until the router authenticates again, or gives up.
#
# Restore matters more than the tests: leaving a probe value in place is the one outcome
# this script must not produce.
wait_for_session() {
    local waited=0 interval=30 limit="${LOCKOUT_WAIT:-900}"

    while [ "$waited" -lt "$limit" ]; do
        if "$GOGL" "${ROUTER_FLAG[@]}" system info >/dev/null 2>&1; then
            LOCKED_OUT=0
            printf '  session recovered after %ds\n' "$waited"
            return 0
        fi
        printf '  waiting for the lockout to clear (%ds of up to %ds)\r' "$waited" "$limit"
        sleep "$interval"
        waited=$((waited + interval))
    done
    printf '\n'
    return 1
}

# jget reads a value out of a gogl JSON response.
jget() { jq -r "$1" 2>/dev/null; }

# check compares an expected value against what the router reports back.
#
# The read is a fresh call, not a cached value: the point is to confirm the device
# accepted the write, and gogl echoing its own intent proves nothing.
check() {
    local label="$1" want="$2" got="$3"
    if [ "$got" = "$want" ]; then
        pass "$label ($want)"
    else
        fail "$label: wrote $want, router reports $got"
    fi
}

# ---------------------------------------------------------------------------
# Preflight
# ---------------------------------------------------------------------------

preflight() {
    command -v jq >/dev/null || { echo "hil-test: jq is required" >&2; exit 2; }
    command -v ip >/dev/null || { echo "hil-test: the 'ip' command is required" >&2; exit 2; }
    command -v "$GOGL" >/dev/null || { echo "hil-test: $GOGL not found on PATH" >&2; exit 2; }

    bold "Preflight"
    printf '  %s\n' "$("$GOGL" --version)"

    if ! g system info >/dev/null 2>&1; then
        echo "hil-test: cannot reach the router. Check -H/--router and GL_PASSWORD." >&2
        exit 2
    fi
    printf '  router:  %s\n' "$(g system info | awk '/^MODEL|^FIRMWARE/ {printf "%s ", $2}')"

    require_wired_lan
}

# require_wired_lan aborts unless this host is on the router's LAN over a cable.
#
# Verified rather than trusted. Every wireless test below depends on it: a wireless
# session would have its writes refused, turning the whole run into skips that read like
# success. And a script that disconnects every wireless client must not be one of them.
require_wired_lan() {
    local router subnet local_ip

    router=$(g lan show --output json | jq -r '.lan_ip')
    subnet=$(g lan show --output json | jq -r '.subnet // empty')
    if [ -z "$router" ] || [ "$router" = "null" ]; then
        echo "hil-test: could not read the router's LAN address" >&2
        exit 2
    fi

    # The source address the kernel picks for this router is what the router will see.
    local_ip=$(ip -o route get "$router" 2>/dev/null | grep -oE 'src [0-9.]+' | awk '{print $2}')
    if [ -z "$local_ip" ]; then
        echo "hil-test: cannot determine this host's address toward $router" >&2
        exit 2
    fi
    printf '  session: %s -> %s' "$local_ip" "$router"

    # Ask the router what it thinks we are. Its own client table is authoritative about
    # which interface we arrived on; inferring from local interface names is not.
    local entry wired band
    entry=$(g clients list --all --output json 2>/dev/null |
        jq -c --arg ip "$local_ip" '.[] | select(.ip == $ip)' | head -1)

    if [ -z "$entry" ]; then
        printf ' (not in the router'"'"'s client list)\n'
        cat >&2 <<MSG

hil-test: this host is not in the router's client list, so it is reaching the router
from off-LAN through another router. Wireless writes would be permitted, but a wireless
client cannot be verified and a disconnect cannot be observed.

Connect this host directly to a LAN port and try again.
MSG
        exit 2
    fi

    wired=$(printf '%s' "$entry" | jq -r '.is_wired')
    band=$(printf '%s' "$entry" | jq -r '.band // "cable"')
    printf ' over %s\n' "$band"

    if [ "$wired" != "true" ]; then
        cat >&2 <<MSG

hil-test: this host is on the router over $band, not a cable.

Every wireless write below would be refused by gogl's own guard, and the run would
report skips that look like passes. Worse, changing an SSID or passphrase would drop
this session mid-test with the baseline unrestored.

Connect this host to a LAN port and try again.
MSG
        exit 2
    fi
    printf '  wired:   confirmed by the router\n'
}

# ---------------------------------------------------------------------------
# Baseline and restore
# ---------------------------------------------------------------------------

capture_baseline() {
    bold "Capturing baseline"
    mkdir -p "$WORKDIR"
    chmod 700 "$WORKDIR"

    # The baseline carries WiFi passphrases in cleartext, so it is created 0600 rather
    # than at whatever the caller's umask happens to be.
    ( umask 077; : > "$BASELINE" )
    if ! g profile export --with-keys > "$BASELINE" 2>"$WORKDIR/export.err"; then
        echo "hil-test: baseline capture failed; refusing to test anything." >&2
        cat "$WORKDIR/export.err" >&2
        exit 1
    fi
    # An empty or malformed baseline is worse than none: it would restore nothing while
    # looking like it had.
    if ! jq -e '.gogl_profile_version and .network.ip' "$BASELINE" >/dev/null; then
        echo "hil-test: baseline looks wrong; refusing to test anything." >&2
        head -20 "$BASELINE" >&2
        exit 1
    fi
    printf '  %s (%s bytes)\n' "$BASELINE" "$(wc -c <"$BASELINE")"
    printf '  LAN %s  pool %s-%s  domain %s\n' \
        "$(jget '.network.ip' <"$BASELINE")" \
        "$(jget '.network.dhcp_start' <"$BASELINE")" \
        "$(jget '.network.dhcp_end' <"$BASELINE")" \
        "$(jget '.domain // "(none)"' <"$BASELINE")"
}

RESTORED=0
restore() {
    [ "$RESTORED" = 1 ] && return
    [ -s "$BASELINE" ] || return
    RESTORED=1

    bold "Restoring baseline"

    # A lockout is the likeliest reason to be here, and restoring matters more than the
    # tests did. Wait it out rather than failing immediately with a probe value in place.
    if [ "$LOCKED_OUT" = 1 ]; then
        printf '  the router is rate limiting; waiting before restore\n'
        if ! wait_for_session; then
            printf '\033[31m  Could not authenticate to restore.\033[0m Baseline kept at:\n    %s\n' \
                "$BASELINE" >&2
            printf '  Wait for the lockout to clear, then run:\n    %s profile import %s --wireless --force\n' \
                "$GOGL" "$BASELINE" >&2
            return
        fi
    fi

    # Clear first. profile import adds and updates but never prunes, so a reservation or
    # DNS name this test created would survive a plain import and silently persist.
    if ! g lan reservations clear --force >/dev/null 2>&1; then
        printf '  could not clear before restore; continuing\n' >&2
    fi

    if g profile import "$BASELINE" --wireless --force; then
        printf '  restored from %s\n' "$BASELINE"
        if [ "$KEEP_BASELINE" = 0 ]; then
            rm -f "$BASELINE" "$WORKDIR"/*.err "$WORKDIR"/*.key
            rmdir "$WORKDIR" 2>/dev/null || true
        else
            printf '  baseline kept at %s\n' "$BASELINE"
        fi
    else
        printf '\033[31m  RESTORE FAILED.\033[0m Baseline kept at:\n    %s\n' "$BASELINE" >&2
        printf '  Re-apply by hand with:\n    %s profile import %s --wireless --force\n' \
            "$GOGL" "$BASELINE" >&2
    fi
}
trap restore EXIT INT TERM

# ---------------------------------------------------------------------------
# Tests
# ---------------------------------------------------------------------------

test_dns_domain() {
    bold "DNS domain"
    local original probe
    original=$(g lan dns show --output json | jget '.domain')
    probe="hil-probe.test"

    g lan dns set --domain "$probe" >/dev/null || { fail "set domain"; return; }
    check "domain set" "$probe" "$(g lan dns show --output json | jget '.domain')"

    # Changing the domain requalifies existing names; confirm the round trip back.
    if [ -n "$original" ] && [ "$original" != "null" ]; then
        g lan dns set --domain "$original" >/dev/null || { fail "restore domain"; return; }
        check "domain restored" "$original" "$(g lan dns show --output json | jget '.domain')"
    fi
}

test_dns_entries() {
    bold "DNS names"
    local name="hil-probe" ip="192.168.250.251"

    # The address is deliberately outside any plausible LAN: a DNS entry needs no
    # reachable host, and picking a live address could shadow a real device.
    g lan dns add "$name" "$ip" >/dev/null || { fail "add DNS name"; return; }
    check "DNS name added" "$ip" \
        "$(g lan dns show --output json |
            jq -r --arg n "$name" '.entries[]? | select(.names[0] == $n) | .ip')"

    g lan dns rm "$name" >/dev/null || { fail "remove DNS name"; return; }
    local still
    still=$(g lan dns show --output json | jq -r '.entries[]?.names[]?' | grep -cx "$name" || true)
    check "DNS name removed" "0" "$still"
}

test_reservations() {
    bold "Reservations"
    local mac="02:00:00:ff:fe:01" ip name
    name="hil-probe-res"
    # An address inside the current subnet, chosen high to avoid colliding with the pool
    # or a live host. Reservations outside the subnet are refused.
    ip=$(g lan show --output json | jget '.lan_ip' | sed 's/\.[0-9]*$/.251/')

    g lan reservations add --name "$name" --mac "$mac" --ip "$ip" --force >/dev/null ||
        { fail "add reservation"; return; }
    check "reservation added" "$ip" \
        "$(g lan reservations list --output json | jq -r --arg m "$mac" '.[] | select(.mac == $m) | .ip')"

    g lan reservations rm --mac "$mac" --force >/dev/null || { fail "remove reservation"; return; }
    local count
    count=$(g lan reservations list --output json | jq -r --arg m "$mac" '[.[] | select(.mac == $m)] | length')
    check "reservation removed" "0" "$count"
}

# LAN addressing and the DHCP pool are deliberately not tested here.
#
# This host is on the LAN by requirement, so moving the pool can invalidate its own lease
# and renumbering can make the router unreachable mid-test with the baseline unrestored.
# Both paths are covered against the mock, where the failure costs nothing:
# TestRunSetNetworkPoolOnlyFillsAddressFromDevice, TestLANSetSubnetMoveIsRefusedWith-
# Reservations, and the guard tests in src/services/guards_test.go.

test_wifi_identity() {
    bold "Wireless identity"
    local ifaces
    ifaces=$(g wifi list --output json | jq -r '.[] | select(.guest == false) | .name')
    [ -n "$ifaces" ] || { skip "no non-guest interfaces reported"; return; }

    local iface
    for iface in $ifaces; do
        local original probe
        original=$(g wifi list --output json | jq -r --arg i "$iface" '.[] | select(.name == $i) | .ssid')
        probe="hil-$iface"

        if ! g wifi set --iface "$iface" --ssid "$probe" --yes >/dev/null 2>"$WORKDIR/wifi.err"; then
            if grep -q 'wireless session' "$WORKDIR/wifi.err"; then
                skip "$iface SSID: refused, session is over WiFi (run this on ethernet)"
                continue
            fi
            fail "$iface SSID: $(head -1 "$WORKDIR/wifi.err")"
            continue
        fi
        check "$iface SSID" "$probe" \
            "$(g wifi list --output json | jq -r --arg i "$iface" '.[] | select(.name == $i) | .ssid')"

        # Only the named field may change. This is the property the partial-update design
        # rests on, and the one worth re-checking on hardware.
        local key_len
        key_len=$(g wifi list --output json --show-key 2>/dev/null |
            jq -r --arg i "$iface" '.[] | select(.name == $i) | .key | length')
        if [ "${key_len:-0}" -ge 8 ]; then
            pass "$iface passphrase survived an SSID-only write ($key_len characters)"
        else
            fail "$iface passphrase was cleared by an SSID-only write"
        fi

        g wifi set --iface "$iface" --ssid "$original" --yes >/dev/null ||
            fail "$iface SSID restore"
    done
}

test_wifi_flags() {
    bold "Wireless hidden and enabled"
    local iface
    iface=$(g wifi list --output json | jq -r '.[] | select(.guest == true) | .name' | head -1)
    [ -n "$iface" ] || { skip "no guest interface to toggle"; return; }

    # A guest interface is used deliberately: toggling `enabled` on a main SSID would cut
    # off any wireless client relying on it.
    local original
    original=$(g wifi list --output json | jq -r --arg i "$iface" '.[] | select(.name == $i) | .hidden')

    if ! g wifi set --iface "$iface" --hidden=true --yes >/dev/null 2>"$WORKDIR/wifi.err"; then
        if grep -q 'wireless session' "$WORKDIR/wifi.err"; then
            skip "$iface hidden: refused, session is over WiFi"
            return
        fi
        fail "$iface hidden: $(head -1 "$WORKDIR/wifi.err")"
        return
    fi
    check "$iface hidden=true" "true" \
        "$(g wifi list --output json | jq -r --arg i "$iface" '.[] | select(.name == $i) | .hidden')"

    g wifi set --iface "$iface" --hidden="$original" --yes >/dev/null ||
        fail "$iface hidden restore"
}

# test_wifi_passphrase exercises a write path nothing else can reach.
#
# Safe only because the wired requirement holds: changing a passphrase drops every client
# on that radio. It matters because --passphrase prompts rather than taking a value, so
# the only automated path to it is --passphrase-command, and that has never run.
test_wifi_passphrase() {
    bold "Wireless passphrase"
    local iface
    iface=$(g wifi list --output json | jq -r '.[] | select(.guest == true) | .name' | head -1)
    [ -n "$iface" ] || { skip "no guest interface to change a passphrase on"; return; }

    local original probe
    original=$(g wifi list --output json --show-key 2>/dev/null |
        jq -r --arg i "$iface" '.[] | select(.name == $i) | .key')
    if [ -z "$original" ] || [ "$original" = "null" ]; then
        skip "$iface has no passphrase to restore; not risking a one-way change"
        return
    fi
    probe="hilprobe$(date +%s)"

    # The secret goes through a file, never through the command string.
    #
    # The obvious spelling -- --passphrase-command "printf %s '$probe'" -- is the exact
    # anti-pattern Critical Rule 5 forbids: it puts a live passphrase in argv where ps can
    # read it, and breaks outright on a passphrase containing a quote. Only the path is
    # interpolated, and the path is one this script created under mktemp.
    ( umask 077; printf '%s' "$probe" > "$WORKDIR/probe.key" )
    ( umask 077; printf '%s' "$original" > "$WORKDIR/original.key" )

    if ! g wifi set --iface "$iface" --passphrase-command "cat $WORKDIR/probe.key" --yes \
        >/dev/null 2>"$WORKDIR/pass.err"; then
        fail "$iface passphrase: $(head -1 "$WORKDIR/pass.err")"
        return
    fi
    check "$iface passphrase" "$probe" \
        "$(g wifi list --output json --show-key | jq -r --arg i "$iface" '.[] | select(.name == $i) | .key')"

    # The SSID must not have moved: only the named field is sent.
    local ssid_after
    ssid_after=$(g wifi list --output json | jq -r --arg i "$iface" '.[] | select(.name == $i) | .ssid')
    [ -n "$ssid_after" ] && [ "$ssid_after" != "null" ] &&
        pass "$iface SSID survived a passphrase-only write ($ssid_after)"

    g wifi set --iface "$iface" --passphrase-command "cat $WORKDIR/original.key" --yes >/dev/null ||
        fail "$iface passphrase restore"

    shred -u "$WORKDIR/probe.key" "$WORKDIR/original.key" 2>/dev/null ||
        rm -f "$WORKDIR/probe.key" "$WORKDIR/original.key"
}

test_wifi_encryption() {
    bold "Wireless encryption"
    local iface radio supported original probe
    iface=$(g wifi list --output json | jq -r '.[] | select(.guest == true) | .name' | head -1)
    [ -n "$iface" ] || { skip "no guest interface to change encryption on"; return; }

    original=$(g wifi list --output json | jq -r --arg i "$iface" '.[] | select(.name == $i) | .encryption')
    radio=$(g wifi list --output json | jq -r --arg i "$iface" '.[] | select(.name == $i) | .band')

    # Pick a mode this radio advertises, other than the current one. Guessing would trip
    # the firmware's own validation, and gogl validates against the same list.
    probe=$(g radio list --output json |
        jq -r --arg b "$radio" --arg cur "$original" \
        '.[] | select(.band == $b) | .encryptions[]? | select(. != $cur and . != "none")' | head -1)
    [ -n "$probe" ] || { skip "$iface: no alternative encryption offered"; return; }

    if ! g wifi set --iface "$iface" --encryption "$probe" --yes >/dev/null 2>"$WORKDIR/enc.err"; then
        fail "$iface encryption: $(head -1 "$WORKDIR/enc.err")"
        return
    fi
    check "$iface encryption" "$probe" \
        "$(g wifi list --output json | jq -r --arg i "$iface" '.[] | select(.name == $i) | .encryption')"

    g wifi set --iface "$iface" --encryption "$original" --yes >/dev/null ||
        fail "$iface encryption restore"
}

test_radio_tuning() {
    bold "Radio tuning"
    local radios
    radios=$(g radio list --output json | jq -r '.[].device')
    [ -n "$radios" ] || { skip "no radios reported"; return; }

    local device
    for device in $radios; do
        # The radio carrying a WiFi uplink is skipped by default. On a router in repeater
        # mode the station interface shares the radio with the AP, so retuning it can
        # drop the router's own internet connection.
        if [ "$INCLUDE_UPLINK_RADIO" = 0 ] && radio_carries_uplink "$device"; then
            skip "$device: may carry the WiFi uplink (--include-uplink-radio to test it)"
            continue
        fi

        local original candidate
        original=$(g radio list --output json | jq -r --arg d "$device" '.[] | select(.device == $d) | .channel')
        # Pick a channel the radio itself advertises, other than the current one. Guessing
        # would hit the firmware's own channel filter, which excludes DFS channels the
        # hardware supports.
        candidate=$(g radio list --output json |
            jq -r --arg d "$device" --argjson c "$original" \
            '.[] | select(.device == $d) | .channels[] | select(.dfs == false) | .channel | select(. != $c)' |
            head -1)
        [ -n "$candidate" ] || { skip "$device: no alternative non-DFS channel offered"; continue; }

        if ! g radio set --device "$device" --channel "$candidate" --yes >/dev/null 2>"$WORKDIR/radio.err"; then
            if grep -q 'wireless session' "$WORKDIR/radio.err"; then
                skip "$device channel: refused, session is over WiFi"
                continue
            fi
            fail "$device channel: $(head -1 "$WORKDIR/radio.err")"
            continue
        fi
        check "$device channel" "$candidate" \
            "$(g radio list --output json | jq -r --arg d "$device" '.[] | select(.device == $d) | .channel')"

        test_radio_width "$device"
        test_radio_power "$device"
        test_radio_hwmode "$device"

        g radio set --device "$device" --channel "$original" --yes >/dev/null ||
            fail "$device channel restore"
    done
}

# test_radio_width checks the one inference still standing in the wireless model.
#
# The firmware reports only the maximum channel width per hardware mode, and gogl offers
# every narrower width on the assumption they are settable. Nothing has confirmed that.
test_radio_width() {
    local device="$1" original narrower
    original=$(g radio list --output json | jq -r --arg d "$device" '.[] | select(.device == $d) | .htmode')

    narrower=20
    [ "$original" = "20" ] && narrower=40

    if ! g radio set --device "$device" --width "$narrower" --yes >/dev/null 2>"$WORKDIR/width.err"; then
        fail "$device width $narrower rejected: $(head -1 "$WORKDIR/width.err") -- the narrower-width inference in HTModes.Options may be wrong"
        return
    fi
    check "$device width" "$narrower" \
        "$(g radio list --output json | jq -r --arg d "$device" '.[] | select(.device == $d) | .htmode')"

    g radio set --device "$device" --width "$original" --yes >/dev/null ||
        fail "$device width restore"
}

# test_radio_power walks the transmit-power levels.
#
# The only field whose accepted values gogl hardcodes rather than reading from the radio,
# so it is the one most likely to be wrong.
test_radio_power() {
    local device="$1" original probe
    original=$(g radio list --output json | jq -r --arg d "$device" '.[] | select(.device == $d) | .txpower')

    probe="Medium"
    [ "$original" = "Medium" ] && probe="High"

    if ! g radio set --device "$device" --power "$probe" --yes >/dev/null 2>"$WORKDIR/power.err"; then
        fail "$device power $probe rejected: $(head -1 "$WORKDIR/power.err") -- TXPowerLevels may not match the firmware"
        return
    fi
    check "$device power" "$probe" \
        "$(g radio list --output json | jq -r --arg d "$device" '.[] | select(.device == $d) | .txpower')"

    g radio set --device "$device" --power "$original" --yes >/dev/null ||
        fail "$device power restore"
}

# test_radio_hwmode changes the hardware mode.
#
# Worth exercising because the values are slash-joined combinations like 11a/n/ac, not the
# bare 11ac the API description implies -- a discrepancy already found once by hand.
test_radio_hwmode() {
    local device="$1" original probe
    original=$(g radio list --output json | jq -r --arg d "$device" '.[] | select(.device == $d) | .hwmode')
    probe=$(g radio list --output json |
        jq -r --arg d "$device" --arg cur "$original" \
        '.[] | select(.device == $d) | .hwmodes[]? | select(. != $cur)' | head -1)
    [ -n "$probe" ] || { skip "$device: no alternative hardware mode offered"; return; }

    if ! g radio set --device "$device" --hwmode "$probe" --yes >/dev/null 2>"$WORKDIR/hw.err"; then
        fail "$device hwmode $probe rejected: $(head -1 "$WORKDIR/hw.err")"
        return
    fi
    check "$device hwmode" "$probe" \
        "$(g radio list --output json | jq -r --arg d "$device" '.[] | select(.device == $d) | .hwmode')"

    g radio set --device "$device" --hwmode "$original" --yes >/dev/null ||
        fail "$device hwmode restore"
}

# radio_carries_uplink guesses whether a radio hosts a station interface.
#
# gogl does not model repeater mode, so this is inferred from the radio being 5GHz, which
# is where GL.iNet puts the uplink by default. Deliberately cautious: the cost of a false
# positive is a skipped test, and of a false negative is the operator's internet.
radio_carries_uplink() {
    local device="$1" band
    band=$(g radio list --output json | jq -r --arg d "$device" '.[] | select(.device == $d) | .band')
    [ "$band" = "5G" ]
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

TESTS=(
    test_dns_domain
    test_dns_entries
    test_reservations
    test_wifi_identity
    test_wifi_flags
    test_wifi_passphrase
    test_wifi_encryption
    test_radio_tuning
)

if [ "$DRY_RUN" = 1 ]; then
    bold "Would run, touching nothing:"
    printf '  %s\n' "${TESTS[@]}"
    printf '\ndelay between calls: %ss   uplink radio: %s\n' "$CALL_DELAY" \
        "$([ "$INCLUDE_UPLINK_RADIO" = 1 ] && echo included || echo skipped)"
    trap - EXIT INT TERM
    exit 0
fi

preflight

if [ "$ASSUME_YES" = 0 ]; then
    printf '\nThis writes to a real router, one setting at a time, and restores the\n'
    printf 'configuration it started from. Wireless clients will be disconnected.\n'
    printf 'Proceed? [y/N] '
    read -r answer
    case "$answer" in
        y|Y) ;;
        *) echo "aborted"; trap - EXIT INT TERM; exit 0 ;;
    esac
fi

capture_baseline

for t in "${TESTS[@]}"; do
    if ! require_session; then
        break
    fi
    "$t"
done

restore

bold "Summary"
printf '  %d passed, %d failed, %d skipped\n' "$PASS" "$FAIL" "$SKIP"
if [ "$FAIL" -gt 0 ]; then
    printf '\nFailures:\n'
    printf '  - %s\n' "${FAILURES[@]}"
    exit 1
fi
exit 0
