# gogl v2: one binary

**Status:** phases 1 and 2 shipped. Written 2026-07-30, status revised 2026-07-30.
Phases 3 and 4 wait on a hardware capture pass — see [Phases](#phases).

This document specified the move from four binaries to one. That move is complete: the CLI
surface it describes is what `gogl` now does, and [`gogl-guide.md`](gogl-guide.md) is the
reference for it. The library under `src/` kept its shape throughout; `src/services`
interfaces and `src/types` are the seam this design leaned on.

---

## What changed and why

**One binary, `gogl <area> <action>`.** Four binaries meant four flag parsers, four help texts,
and four copies of the connection flags. Noun-verb is what git, docker, kubectl and gh converged
on, and it scales to the ~35 leaf commands this surface needs.

**No external file-format obligation.** An earlier constraint held the command names in step
with a sibling project's. It was dropped rather than allowed to decline a better command tree.
The ISC DHCP host-declaration format **stays**, purely on its own merits: it diffs, it reviews,
it lives in git, and it is readable without a tool.

**Configuration moves to TOML under XDG paths**, with named routers and a password command.

---

## Non-goals

- **GL.iNet firmware 4.x only.** A UCI-over-`/ubus` backend for generic OpenWrt was specified
  here and then dropped: it would have meant two implementations of every service, two credential
  paths, and two sets of field translations, for hardware this project does not target. See
  [Scope](../VISION.md#scope).
- **No SSH.** Per Critical Rule 5: the API is a capability boundary, and user data never enters a
  command string.
- **No root on the operator's machine.** Nothing this tool does needs it; it installs to
  `~/.local/bin`.
- **Not a full router backup.** `gogl profile` captures a network, not a device. `sysupgrade -b`
  over SSH remains the right tool for a byte-exact restore, and always will be.

---

## Command tree

**Endpoint** is whether the underlying call has been exercised against real hardware —
the scheduling input, not a wish list. **Shipped** is whether the command exists today.
[`gogl-guide.md`](gogl-guide.md) documents every shipped flag; this table is the plan's
own scorecard.

| Command | Endpoint | Shipped |
|---|---|---|
| `gogl lan show [--show-key]` | verified | yes |
| `gogl lan set [--guest] --ip --mask --pool-start --pool-end [--force] [--dry-run]` | verified | yes |
| `gogl lan leases` | verified | yes |
| `gogl lan reservations list` | verified | yes |
| `gogl lan reservations add --name --mac --ip` | verified | yes |
| `gogl lan reservations rm --name\|--mac\|--ip` | verified | yes |
| `gogl lan reservations clear` | verified | yes |
| `gogl lan reservations export` | verified | yes |
| `gogl lan reservations import <file> [--prune] [--force] [--dry-run]` | verified | yes |
| `gogl lan dns show` | verified | yes |
| `gogl lan dns set --domain <d>` | workaround | yes |
| `gogl lan dns add <name> <ip>` / `rm <name>` / `clear` | verified | yes |
| `gogl radio list` / `show --band 5` | verified | yes |
| `gogl radio set --band 5 --channel --width --hwmode --power` | verified | yes |
| `gogl wifi list` / `show --band 5 [--guest]` | verified | yes |
| `gogl wifi set --band 5 [--guest] --ssid --passphrase --encryption --hidden --enabled` | verified | yes |
| `gogl clients list [-w\|-e] [-r] [-a]` | verified | yes |
| `gogl clients vendor <mac>` | offline | yes |
| `gogl profile export [--with-keys]` / `import <file> [--wireless]` | verified | yes |
| `gogl system info` | verified | yes |
| `gogl config show` / `routers` / `init` | local | yes |
| `gogl completion bash\|zsh\|fish` | free | yes |
| `gogl system reboot` | **not captured** | no |
| `gogl access check` | verified | no |
| `gogl access passwd` | **not captured** | no |
| `gogl access users list\|add\|rm` | **not captured** | no |
| `gogl wan show` | **not captured** | no |
| `gogl wan set-mode dhcp\|static\|pppoe` | **not captured** | no |
| `gogl wan repeater scan\|connect` | **not captured** | no |
| `gogl wan tethering show` / `wan modem show` | **not captured** | no |

Two departures from the plan as written, both settled by implementation:

- **`--guest` is on `lan set`, not `lan show`.** The report already names every interface it
  finds, so a flag to select one had nothing to do.
- **`clients list` filters by link type, not by band** — `-w/--wifi` and `-e/--wired`, plus
  `-r/--reserved` and `-a/--all`. A band filter implied the client list carried a band per
  station; it does not.

`lan reservations` accepts `res` as an alias, since four words before an argument is a lot.

### Verb vocabulary, held strictly

| Verb | Means |
|---|---|
| `list` | many items |
| `show` | one thing, in detail |
| `set` | write fields on a thing |
| `add` / `rm` | collection membership |
| `clear` | empty a collection |
| `import` / `export` | file in, file out |
| `check` | verify without changing |

Banned: `get` (ambiguous with `show`), `del`/`delete` (pick `rm`), `create` (use `add`), `update`
(use `set`).

This is why `lan dns set --domain X` is a field write and `lan dns add nas 192.168.8.13` is a
member add. One verb, one meaning, everywhere.

### Band abstraction

`--band 2.4|5|6` resolves to a device at runtime by reading the band each radio reports. Not a
static map: newer models have a third radio, and nothing guarantees `radio0` is 2.4GHz across the
product line. If two radios report the same band, `--device` becomes required rather than gogl
guessing.

The same resolution serves `wifi`, where interface names are worse than device names:
`--band 5` finds `default_radio1`, `--band 5 --guest` finds `guest5g`.

### Guest

A `--guest` flag on `lan` and `wifi`, not an area. That follows the firmware's own model:
`lan.get_config_list` returns `lan` and `guest` as interface variants, and guest SSIDs are more
entries in `wifi`'s `ifaces` array. One concept, two facets, no duplicated verbs.

---

## One backend, deliberately

The service interfaces in `src/services` are semantic rather than wire-shaped, which made a
second backend architecturally cheap — and it was specified here, over UCI on `POST /ubus`, for
generic OpenWrt. It is dropped.

The reasoning is worth keeping, because it explains field names that otherwise look arbitrary.
The devices run OpenWrt 18.06, and every group `gogl` touches — `lan`, `dns`, `wifi`, `clients`,
`network`, `system` — wraps UCI `network`, `dhcp`, `wireless` and the `hosts(5)` file. So GL.iNet
is translating standard concepts into its own JSON, and the translation is lossy in places:

| Concept | GL.iNet API | UCI underneath |
|---|---|---|
| DHCP pool | `start` + `end`, both addresses | `start` + `limit`, an offset and a **count** |
| Channel width | `20`, `40`, `80`, `auto` | `htmode`: `HT20`, `VHT80`, `HE80` |
| dnsmasq domain | **absent** | `dhcp.@dnsmasq[0].domain` |
| Hardware mode | `11a/n/ac`, slash-joined | `hwmode` plus `htmode` |

That table is documentation, not a compatibility plan. It is why `VHT80` felt like the natural
value and was wrong — UCI vocabulary leaking into a GL.iNet API — and why the DHCP pool needs a
real translation rather than a rename. Reading it as "we could support both" is what this section
now exists to correct.

Two consequences follow, both simplifications:

**`src/types` holds GL.iNet's shapes without apology.** `IntBool` exists because the firmware
sends `enable: 1`; `HTModes` is its capability object with the narrower-widths inference on top.
An earlier version of this document listed those as vendor artifacts to be moved out of a neutral
layer. There is no neutral layer, so that work is cancelled.

**The domain-in-a-marker-line is the answer, not a workaround.** GL.iNet exposes no dnsmasq domain
setting, `POST /ubus` returns 404 on this hardware, and no better mechanism is coming. It stops
being described as temporary.

## Configuration

`${XDG_CONFIG_HOME:-~/.config}/gogl/config.toml`:

```toml
default = "player-test"
output  = "text"            # or "json"

[routers.player-test]
host             = "192.168.8.1"
domain           = "bench.test"
password_command = "pass show routers/player-test"

[routers.lab]
host = "192.168.4.1"
```

`gogl --router lab lan show`, or omit `--router` for the default.

**Password resolution, highest wins:** `GL_PASSWORD` env, `password_command` in TOML,
interactive prompt without echo. Never a `--password` flag: a secret on argv is visible in `ps`
and lands in shell history. As shipped there is no `--password-command` flag either — the
router password comes only from those three sources, and the flag was dropped as a fourth way
to say what the TOML already says. (WiFi passphrases *do* have `--passphrase-command`, because
they are per-write rather than per-router.)

`password_command` executes a command named in a file the user owns at `0600` — the same trust
model as `git`'s `credential.helper` and `restic`'s `--password-command`.

### Paths

| What | Where |
|---|---|
| Config | `${XDG_CONFIG_HOME:-~/.config}/gogl/config.toml` |
| OUI cache | `${XDG_CACHE_HOME:-~/.cache}/gogl/oui.txt` |
| Install | `~/.local/bin/gogl` |

Both moved as specified. The OUI database had been under `XDG_DATA_HOME`, which was wrong — it
is a re-downloadable cache, not user data. Install had targeted `~/bin`; `~/.local/bin` is on
`PATH` by default via systemd and Debian's `.profile`, and `make uninstall` still clears the
legacy location so a stale copy cannot shadow a fresh build.

### Global flags

`--router`, `-H/--host`, `-p/--port`, `--https`, `--insecure`, `--output text|json`, `--yes`.
Per-command where they mean something: `--dry-run`, `--force`, `--guest`, `--band`, `--device`.

**Exit codes:** `0` success, `1` error, `2` usage, `3` refused by a guard. The last one lets a
script distinguish "I was blocked" from "it broke", which matters for the domain and reservation
guards.

---

## Secrets on the command line

`--passphrase` with no value prompts without echo; `--passphrase-command` reads it from a
command. Never a bare `--passphrase <value>`.

Named `--passphrase` rather than `--password` so that `gogl access passwd` — the router's own
password — cannot be confused with a WiFi key.

This was a fix, not a new rule: the v1 wireless utility's `--set-key 'value'` put a WiFi
passphrase in argv and in shell history, contradicting the reason `GL_PASSWORD` was never given
a flag. Passing a value to `--passphrase` is now an error rather than being silently accepted.

---

## Implementation strategy

**The library does not change.** `src/services`, `src/transport`, `src/auth` and `src/mock` keep
their shapes.

**The four utilities became importable packages** rather than being merged into one. Merging four
`package main`s would have meant renaming ~15 colliding identifiers -- `formatText`, `formatJSON`,
`mockClient`, `isConnectionLost`, `stubNetwork`, `testLAN`, `lanFixture` and several `Test*` names
-- which is the sweeping-rename shape that has already corrupted files in this project once.
Separate packages keep every namespace intact and need no renames at all.

| Was (v1 binary) | Is | Exported entry points |
|---|---|---|
| reservations utility | `utilities/internal/reservations` | `Get`, `Set`, `Add`, `Del`, `Clear`, `SetDomain`, `ParseHosts`, `FormatHosts`, `Modes` |
| network utility | `utilities/internal/netcfg` | `Show`, `BuildReport`, `FormatText`, `FormatJSON`, `FormatWireless`, `SetNetwork`, `SetWireless`, `SetSSID`, `NetworkModes`, `WirelessModes` |
| clients utility | `utilities/internal/clients` | `List`, `BuildEntries`, `FormatText`, `FormatJSON`, `FilterFor`, `LoadOUI`, `ParseOUI` |
| profile utility | `utilities/internal/profile` | `Capture`, `Apply`, `ReadProfile`, `Profile`, `CaptureOptions`, `ApplyModes` |

Two compositions were extracted from the deleted `main.go` files rather than dissolved into the
command layer, so their ordering stays covered by the package's own tests: `netcfg.Show` (report
the LAN, then the radios, treating a wireless read failure as a warning) and `clients.List` (load
the OUI database before touching the device, then read clients and optionally reservations).

Two things were deleted rather than carried forward. `checkModes` and its test enforced "exactly
one mode selected", which a subcommand tree makes structurally impossible to violate. And
`reservations.Modes` lost its `Get`/`Set`/`Add`/`Del`/`Clear` booleans for the same reason.

`netcfg.optionalBool`, `optionalInt` and `optionalString` are **gone**, as anticipated: pflag's
`Changed()` answers the set-versus-unset question natively, so the command tree made them
redundant.

The move kept the whole test suite alive; it now stands at 549 test functions and `make test`
is clean under `-race`.

**The four binaries were retired in one commit** once the last area was ported, with a
`make uninstall` that also clears the stale copies in `~/bin`.

### Phases

1. **Shipped.** Command tree over verified endpoints: `lan`, `radio`, `wifi`, `clients`,
   `profile`, `system info`, `config`. The four utilities became importable packages with
   tests intact, and `utilities/gogl` is the cobra tree wiring them together.
2. **Shipped.** TOML config with named routers and `password_command`, XDG paths, install to
   `~/.local/bin`, `make uninstall`, `--version`, `make check-docs`.
3. **Not started.** Capture pass on the SFT1200 for `wan`, `access` and `system reboot`, with
   `discovery/shape`. See [`../TODO.md`](../TODO.md).
4. **Not started.** `wan` and `access`, written from those captures rather than from GL.iNet's
   descriptions.

Phases 1 and 2 needed no hardware, which is why they went first. Nothing in 3 or 4 gets written
before the capture exists, because building from the vendor description has been wrong three
times: `dhcp.*` does not exist, `dns.set_host` rejects `(`, `)` and `=`, and `htmodes` is an
object rather than an array.

The former phases 4 and 5 — a `backend/openwrt` implementation, and moving GL.iNet wire types out
of `src/types` — are **cancelled** with the portability goal.

## Open questions

- **`wan` shape.** Whether it stays one area, and how `repeater`, `tethering` and `modem` sit
  under it. Cannot be settled before the capture. No longer complicated by "how do
  vendor-specific commands fit a portable tool", since every command here is vendor-specific.

Settled by the GL.iNet-only decision: the name stays `gogl`, the logo stays, and there is no
convention needed for vendor-specific commands. Settled by implementation: `netcfg`'s
`optional*` helpers are gone.

## Related

- [`../VISION.md`](../VISION.md) — requirements, Critical Rule 5 as reworded
- [`DESIGN.md`](DESIGN.md) — the current architecture this builds on
- [`../GL_INET_4X_API_DOCUMENTATION.md`](../GL_INET_4X_API_DOCUMENTATION.md) — hardware-verified
  GL.iNet surface, and where it disagrees with its own documentation
- [`../TODO.md`](../TODO.md) — the capture pass this design waits on
