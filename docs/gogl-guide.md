# gogl: the complete command reference

`gogl` manages GL.iNet 4.x travel routers over the firmware's JSON-RPC API.

Its purpose is making a network reproducible: capture the addressing and names a site
depends on, then recreate them on a pocket router somewhere else. Everything goes through
the router's own API — no SSH, no shell — so the blast radius of a mistake is bounded by
what the API can express.

**Scope:** GL.iNet firmware 4.x only. Verified against a GL-SFT1200 (Opal) on 4.3.28. Not
a portable OpenWrt tool and will not become one; see [Scope](../VISION.md#scope).

---

## Contents

- [Shape of the command line](#shape-of-the-command-line)
- [Global flags](#global-flags)
- [Configuration file](#configuration-file)
- [Passwords and secrets](#passwords-and-secrets)
- [Exit codes](#exit-codes)
- [`gogl lan`](#gogl-lan) — address, pool, leases
- [`gogl lan reservations`](#gogl-lan-reservations) — static bindings
- [`gogl lan dns`](#gogl-lan-dns) — domain and names
- [`gogl radio`](#gogl-radio) — channel, width, power
- [`gogl wifi`](#gogl-wifi) — SSID, passphrase, encryption
- [`gogl clients`](#gogl-clients) — connected stations
- [`gogl profile`](#gogl-profile) — whole-network capture and apply
- [`gogl system`](#gogl-system) — device identity
- [`gogl config`](#gogl-config) — gogl's own settings
- [`gogl completion`](#gogl-completion) — shell completion
- [The three guards](#the-three-guards)
- [Two things the firmware does that will surprise you](#two-things-the-firmware-does-that-will-surprise-you)
- [Recipes](#recipes)
- [Troubleshooting](#troubleshooting)

---

## Shape of the command line

```
gogl [global flags] <area> [subarea] <action> [flags] [arguments]
```

Nine areas. Eight act on the router; `config` acts on your machine.

| Area | What it covers |
|---|---|
| `lan` | LAN address, DHCP pool, leases, and the `reservations` and `dns` subareas |
| `radio` | Per-radio tuning: channel, width, hardware mode, transmit power |
| `wifi` | Per-interface identity: SSID, passphrase, encryption, hidden, enabled |
| `clients` | Connected stations, with IEEE OUI manufacturer lookup |
| `profile` | A whole network captured to JSON, and applied back |
| `system` | Model, firmware, endpoint |
| `config` | gogl's own configuration file |
| `completion` | Shell completion scripts |

### Verb vocabulary

Held strictly, so a verb means one thing everywhere.

| Verb | Means |
|---|---|
| `list` | many items |
| `show` | one thing, in detail |
| `set` | write fields on a thing |
| `add` / `rm` | collection membership |
| `clear` | empty a collection |
| `import` / `export` | file in, file out |

Not used: `get` (ambiguous with `show`), `delete` (it is `rm`), `create` (it is `add`),
`update` (it is `set`).

This is why `lan dns set --domain X` is a field write while `lan dns add nas 192.168.8.13`
is a member add.

---

## Global flags

Accepted by every command.

| Flag | Default | Meaning |
|---|---|---|
| `--router NAME` | the config file's `default` | Which configured router to use |
| `-H`, `--host ADDR` | from config, then `GL_ROUTER_IP` | Router address, overriding the config |
| `-p`, `--port N` | `80` | Router port. Not 443: this firmware serves plain HTTP |
| `--https` | off | Use TLS |
| `--insecure` | **on** | Skip certificate verification |
| `--output text\|json` | `text`, or the config's `output` | Output format |
| `-v`, `--version` | | Print version, revision, and whether the tree was dirty |
| `-h`, `--help` | | Help for any command or subcommand |

`--insecure` defaults to *on* because these devices ship self-signed certificates and a
CLI that cannot reach one out of the box is useless. The Go library defaults the other
way: `gogl.Config` verifies TLS at its zero value, because a library must not be insecure
by default.

`--version` exists because a stale binary shadowing a fresh one is hard to spot:

```console
$ gogl --version
gogl v0.9.1 (491fb2f11996, dirty)
```

---

## Configuration file

`${XDG_CONFIG_HOME:-~/.config}/gogl/config.toml`. Override the whole path with
`GOGL_CONFIG`.

A missing file is not an error — gogl works from flags and the environment alone.

```toml
default = "home"
output  = "text"          # or "json"

[routers.home]
host             = "192.168.8.1"
domain           = "lab.example"
password_command = "pass show routers/home"

[routers.travel]
host = "192.168.8.1"
port = 80
```

| Key | Scope | Meaning |
|---|---|---|
| `default` | top level | Router used when `--router` is absent. Optional with one router defined |
| `output` | top level | `text` or `json` |
| `host` | per router | **Required.** Address |
| `port` | per router | Defaults to 80 |
| `username` | per router | Defaults to `root`, the only account standard firmware has |
| `https`, `insecure` | per router | TLS behaviour |
| `domain` | per router | A note to yourself. Applied only by `lan dns set`, never implicitly |
| `password_command` | per router | A command printing the password on its first line |

`gogl config init` writes a commented starting point.

Precedence, highest first: **command-line flag → config file → environment variable.** A
flag counts as given only if you actually typed it, so `--port 80` is distinguishable from
omitting `--port`.

Related paths:

| What | Where |
|---|---|
| Config | `${XDG_CONFIG_HOME:-~/.config}/gogl/config.toml` |
| OUI cache | `${XDG_CACHE_HOME:-~/.cache}/gogl/` |
| Install | `~/.local/bin/gogl` |

---

## Passwords and secrets

**There is no `--password` flag, and there never will be.** A secret on the command line
is visible to every user through `ps` and is recorded in your shell history.

Router password resolution, highest first:

1. `GL_PASSWORD` in the environment
2. `password_command` in the config file
3. An interactive prompt, echo off, read from `/dev/tty`

```bash
read -rsp 'router password: ' GL_PASSWORD; export GL_PASSWORD    # for a session
```

`password_command` follows the same pattern as `git`'s `credential.helper` and `restic`'s
`--password-command`. Only the first line of output is used, because `pass show` prints the
password first and metadata after.

**WiFi passphrases follow the same rule.** `gogl wifi set --passphrase` takes no value; it
prompts. For scripts, `--passphrase-command` reads it from a command's output. Passing a
value is rejected rather than silently ignored, so nothing lands in your history by
accident.

Two things to know about what the API exposes:

- `wifi.get_config` returns passphrases in **cleartext** over plain HTTP. Reading them
  needs LAN access and the admin password — the same bar as opening the admin panel — so
  gogl masks them in output unless you pass `--show-key`. That keeps a key out of your
  scrollback; it is not an access control.
- `gogl profile export` omits passphrases by default. See [`gogl profile`](#gogl-profile).

---

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | Failure: the router was unreachable, a write was rejected, something broke |
| `2` | Usage: the command was invoked wrongly |
| `3` | **Refused by a guard:** the request was well-formed but the state was wrong |

Code 3 exists so a script can tell "I was blocked" from "it broke". It covers
`ErrDomainNotSet`, `ErrReservationsExist`, and `ErrWirelessSession` — the three
[guards](#the-three-guards).

---

## `gogl lan`

The wired network.

### `gogl lan show`

Reports the LAN, the pool, reservation counts, and every radio.

| Flag | Meaning |
|---|---|
| `--show-key` | Print WiFi passphrases instead of masking them |

```console
$ gogl lan show
MODEL      sft1200
FIRMWARE   4.3.28
LAN        192.168.8.0/24  (255.255.255.0)
DHCP       enabled
POOL       192.168.8.100 - 192.168.8.149  (50 addresses)
LEASE      12h
INTERFACE  lan
GATEWAY    -
DNS        -
RESERVED   27
  IN POOL  20  (honored by dnsmasq, excluded from the pool)
AVAILABLE  96
```

Wireless is reported alongside the wired network, because "what is this network" includes
what devices associate to. A router that will not report wireless still gets its LAN
reported — the wireless read is a warning, not a failure.

`IN POOL` appears only when reservations fall inside the DHCP range. Nothing is broken:
dnsmasq honours a static bind inside the dynamic range and excludes that address from
allocation. It is reported because it arises silently — a subnet move put 20 of 27
reservations inside the pool on real hardware — and because it explains an `AVAILABLE`
count that otherwise looks wrong.

### `gogl lan set`

Writes the LAN address or the DHCP pool.

| Flag | Meaning |
|---|---|
| `--pool-start ADDR` | New pool start |
| `--pool-end ADDR` | New pool end |
| `--ip ADDR` | New LAN address |
| `--mask MASK` | New netmask |
| `--guest` | Write the guest interface instead of the main LAN |
| `--force` | Allow a subnet move while reservations exist |
| `--dry-run` | Show the change and any refusal without writing |

**Changing only the pool** needs neither `--ip` nor `--mask`; they are read from the
device. Nothing moves, the session survives, and it is never refused:

```bash
gogl lan set --pool-start 192.168.8.50 --pool-end 192.168.8.150
```

**Moving the subnet** requires all four flags together, because a pool from the old subnet
cannot be valid in a new one:

```bash
gogl lan set --ip 192.168.4.1 --mask 255.255.255.0 \
             --pool-start 192.168.4.100 --pool-end 192.168.4.149
```

That **drops your session** — the router changes address mid-call. gogl treats a lost
connection as success and tells you where to reconnect. It is refused while reservations
exist unless `--force`; see [the guards](#the-three-guards).

A netmask change counts as moving the subnet even when the address stays put.

### `gogl lan leases`

Lists the dynamic DHCP leases. Leases are not reservations — a lease expires. This is how
you find what is worth reserving.

```console
$ gogl lan leases
IP             MAC                HOSTNAME  EXPIRES
192.168.8.135  6e:1d:47:db:54:54  iPhone    11h32m0s
```

`EXPIRES` is time remaining, which is what the firmware reports.

---

## `gogl lan reservations`

Static MAC-to-IP bindings. Aliased as `res`, since the canonical form is four words.

**Each host declaration is two writes.** A static bind for the address, and a host-file
entry for the name — the firmware stores them separately and joins them for nobody. These
commands keep both in step, so you write one declaration and get both. See
[Two things the firmware does](#two-things-the-firmware-does-that-will-surprise-you).

### `gogl lan reservations list`

Lists the reservations. Text output is ISC DHCP format; `--output json` gives the raw
records.

### `gogl lan reservations export`

Writes every reservation to stdout in ISC DHCP host-declaration format.

```console
$ gogl lan reservations export > lab.hosts
$ head -8 lab.hosts
# gogl reservations
# exported from GL.iNet router at 192.168.8.1
# lan: 192.168.8.0/24  pool: 192.168.8.100-192.168.8.149  lease: 12h
# date: 2026-07-30

host europa {
    hardware ethernet 10:51:07:1f:8d:1c;
    fixed-address 192.168.8.10;
}
```

The format is kept on its own merits, not for compatibility: it diffs, it reviews, it lives
in git, and it is how a UniFi dump from `gofips` gets in.

### `gogl lan reservations import [file]`

Imports host declarations from a file, or stdin when no file is given.

| Flag | Meaning |
|---|---|
| `--prune` | Also delete reservations and names on the router but absent from the file |
| `--force` | Proceed past conflicts |
| `--dry-run` | Show what would change without changing it |

Four phases, in a deliberate order. All file validation before any device contact, so a
malformed file never half-writes a router. All reads before any write, so the diff is
against one snapshot. Then bindings one at a time, and every name change in a single
host-file write.

**Idempotent.** A second run reports everything skipped and leaves the host file
byte-identical. That is what makes a host file usable as a checked-in description of a
network.

**Repairs drift.** Bindings and names are diffed separately, so a binding whose name went
missing gets the name back rather than being skipped because the binding matched.

Requires a configured domain — see [the guards](#the-three-guards).

### `gogl lan reservations add [declaration]`

Adds one host, by flags or as a declaration fragment.

| Flag | Meaning |
|---|---|
| `--name NAME` | Hostname |
| `--mac MAC` | MAC address |
| `--ip ADDR` | IPv4 address |
| `--force` | Proceed past conflicts |
| `--dry-run` | Show what would change without changing it |

```bash
gogl lan reservations add --name nas --mac aa:bb:cc:dd:ee:01 --ip 192.168.8.13
gogl lan reservations add 'host nas { hardware ethernet aa:bb:cc:dd:ee:01; fixed-address 192.168.8.13; }'
```

The three flags are required together. Giving both a declaration and flags is an error
rather than a silent preference. With neither, the fragment is read from stdin.

### `gogl lan reservations rm`

Removes one host and its DNS name, identified by exactly one of `--name`, `--mac`, `--ip`.

| Flag | Meaning |
|---|---|
| `--name`, `--mac`, `--ip` | Which host |
| `--force` | Skip the confirmation prompt |
| `--dry-run` | Show what would change without changing it |

The name goes first, then the binding. A leftover binding is an address with no name, which
the next import repairs; a leftover name keeps resolving to an address nothing holds, which
is the worse failure.

### `gogl lan reservations clear`

Removes **every** reservation and **every** managed DNS name.

| Flag | Meaning |
|---|---|
| `--force` | Skip the confirmation prompt |
| `--dry-run` | List what would go without removing it |

Both tables, because they are one intent stored in two places. The DNS domain survives — it
is configuration, not content.

This is also the precondition for moving the LAN subnet without `--force`.

---

## `gogl lan dns`

The DNS domain and the names the router resolves.

**A reservation does not create a DNS record.** On this firmware the name on a binding is a
label. Records live in the router's host file, which dnsmasq answers from, and that is what
these commands write.

### `gogl lan dns show`

Reports the domain and every managed name.

```console
$ gogl lan dns show
DOMAIN     lab.example

ADDRESS        NAMES
192.168.8.13   nas  nas.lab.example
192.168.8.14   pi   pi.lab.example
```

With no domain configured it says so, because that state blocks every reservation write.

### `gogl lan dns set`

Sets the DNS domain.

| Flag | Meaning |
|---|---|
| `--domain DOMAIN` | The DNS suffix, e.g. `lab.example` |

Required before any reservation write. Changing an existing domain requalifies every managed
name, so resolution does not split between two suffixes.

The domain is stored inside gogl's block in the host file. The firmware exposes no dnsmasq
domain setting, so gogl writes fully-qualified names instead — which works for any suffix —
and keeps the suffix where it travels with the device.

### `gogl lan dns add NAME ADDR`

Points a name at an address. The entry carries both the bare name and its qualified form,
so either resolves. A name already in use is replaced.

### `gogl lan dns rm NAME`

Removes a name, in either its bare or qualified form. Both spellings go — leaving the bare
name behind would keep the host resolving after being asked to stop.

### `gogl lan dns clear`

Removes every managed name, keeping the domain and everything outside gogl's block —
the file also carries the loopback and IPv6 entries the router resolves its own name from.

This leaves reservations in place. To clear both, use `gogl lan reservations clear`.

---

## `gogl radio`

Per-radio tuning: the fields the firmware scopes by radio rather than by SSID, so they
affect every network on that radio.

### Selecting a radio

`--band 2.4`, `--band 5`, or `--band 6`. Spellings like `2`, `2g`, `2.4GHz` all work.

Resolution reads the band each radio **reports** — never a static `radio0`/`radio1` map.
Nothing guarantees that ordering across models, and a three-radio device would break a
fixed map silently. If two radios report one band, `--device` becomes required rather than
gogl guessing.

`--device NAME` names a radio directly and overrides `--band`.

### `gogl radio list`

Every radio with the values it accepts.

```console
$ gogl radio list
BAND  INTERFACE       SSID       ENCRYPTION  STATE    KEY
5G    default_radio1  lab-5g     psk2        enabled  (11 characters, --show-key to reveal)

5G radio radio1: channel 44, 20, 11a/n/ac, power Max
  channels:    36, 40, 44, 48, 149, 153, 157, 161, 165, or 0 for auto
  bandwidths:  20, 40, 80
  hw modes:    11ac, 11n/ac, 11a/n/ac
  encryptions: none, psk2, psk-mixed, sae, sae-mixed
  power:       Max, High, Medium, Low
```

Those lists come from the radio, and the tuning flags are unusable without them: valid
values differ per radio and per regulatory domain. gogl validates against them, so a bad
channel is refused with the available ones named.

Note the vocabulary. Bandwidths are **widths in MHz** — `20`, `40`, `80`, plus `auto` where
offered — not the `HT20`/`VHT80` form GL.iNet's own API description implies. Hardware modes
are slash-joined combinations like `11a/n/ac`, not bare `11ac`. `sae` is WPA3.

### `gogl radio show`

Reports one radio. Takes `--band` or `--device`.

### `gogl radio set`

| Flag | Meaning |
|---|---|
| `--band`, `--device` | Which radio |
| `--channel N` | Channel, or `0` for automatic |
| `--width MHZ` | `20`, `40`, `80`, or `auto` |
| `--hwmode MODE` | e.g. `11a/n/ac` |
| `--power LEVEL` | `Max`, `High`, `Medium`, `Low` |
| `--yes` | Skip the confirmation prompt |
| `--dry-run` | Show the change and any refusal without writing |

```bash
gogl radio set --band 5 --channel 149
gogl radio set --band 2.4 --channel 0                 # automatic
gogl radio set --device radio1 --width 80 --power Low
```

Only the fields you name are sent. Refused over a wireless session — see
[the guards](#the-three-guards).

**DFS channels get a warning.** The radio must vacate one if it detects radar, taking every
client with it for the minutes it spends re-scanning. A poor choice for kit that has to come
up reliably in an unfamiliar building.

---

## `gogl wifi`

Per-interface identity: the fields the firmware scopes by SSID rather than by radio, so a
guest and a main SSID on one radio are set independently.

Selected the same way as a radio — `--band`, plus `--guest` for the guest interface on it.
`--iface NAME` names an interface directly.

### `gogl wifi list`

Every wireless interface. `--show-key` reveals passphrases.

### `gogl wifi show`

One interface, with its passphrase length rather than its value.

### `gogl wifi set`

| Flag | Meaning |
|---|---|
| `--band`, `--guest`, `--iface` | Which interface |
| `--ssid NAME` | New SSID, up to 32 characters |
| `--passphrase` | **Takes no value.** Prompts, echo off |
| `--passphrase-command CMD` | Read the passphrase from a command's first line |
| `--encryption MODE` | e.g. `psk2`, or `sae` for WPA3 |
| `--hidden=true\|false` | Whether the SSID is advertised |
| `--enabled=true\|false` | Whether the interface is up |
| `--yes` | Skip the confirmation prompt |
| `--dry-run` | Show the change and any refusal without writing |

```bash
gogl wifi set --band 5 --ssid lab-5g
gogl wifi set --band 5 --passphrase
gogl wifi set --band 2.4 --guest --enabled=true --hidden=false --ssid lab-guest
```

**Note `--hidden=false`, with the value attached.** A partial update must distinguish "set
this to false" from "leave it alone", and a bare boolean flag meaning true is how
`--enabled` ends up disabling something.

**Only the fields you name are sent.** Setting a passphrase leaves the SSID, encryption and
enabled state exactly as they were — verified on hardware.

Refused over a wireless session; `--yes` does not override that.

---

## `gogl clients`

Connected stations. Its own area rather than part of `lan`, because a station arrives over
cable, 2.4GHz or 5GHz and the useful view is all of them together.

### `gogl clients list`

| Flag | Meaning |
|---|---|
| `-a`, `--all` | Include clients the router remembers but does not currently see |
| `-w`, `--wifi` | Only wireless stations |
| `-e`, `--wired` | Only wired stations |
| `-r`, `--reserved` | Mark which stations hold a reservation |

```console
$ gogl clients list
MAC                ADDRESS        HOSTNAME  MANUFACTURER      SINCE
10:51:07:1f:8d:1c  192.168.8.10   europa    Intel Corporate   3h12m0s
6e:1d:47:db:54:54  192.168.8.135  iPhone    randomized        41m0s
```

**Only clients the router currently sees, by default.** The client list carries history: a
router renumbered from `192.168.2.0/24` to `192.168.8.0/24` was still listing a station at
`192.168.2.138`. Presenting that beside live clients with no distinction is worse than not
showing it. `--all` includes them and adds an `ONLINE` column.

`SINCE` is how long a client has been connected, where the firmware's value can be
understood. Its format is undocumented, so gogl renders it when it makes sense and leaves it
blank when it does not, rather than guessing.

Manufacturer lookup uses the IEEE OUI registry, downloaded and cached independently of the
router. `randomized` means a locally-administered address that identifies nobody by design.

### `gogl clients vendor MAC`

Looks up a MAC's manufacturer. **Entirely offline** — reads the cache and never opens a
session.

```console
$ gogl clients vendor b4:0e:cf:2a:85:6f
b4:0e:cf:2a:85:6f  Bouffalo Lab (Nanjing) Co., Ltd.
```

---

## `gogl profile`

A whole network captured to JSON, and applied back — to the same router or a different one.

**A profile is not a router image.** It carries what defines a *network*: the LAN and pool,
reservations, DNS names and domain, wireless identity, radio tuning. It omits everything
identifying a particular unit — the router's own MAC, serial, uptime, lease state. That
omission is the point: those fields are what make a full config dump useless on a second
device. Client MACs *are* included, since a reservation is a MAC-to-IP binding.

Absent for a reason: lease time, upstream DNS, firewall, VPN, VLANs. No endpoint verified
against hardware writes them.

For a byte-exact backup of one device, `sysupgrade -b` over SSH is the right tool and always
will be.

### `gogl profile export`

| Flag | Meaning |
|---|---|
| `--with-keys` | Include WiFi passphrases in cleartext |

```bash
gogl profile export > lab.json                 # no passphrases; safe to commit
gogl profile export --with-keys > lab.json     # passphrases included
```

**An omitted key is not an empty key.** On import a missing key is not written at all,
leaving whatever the target already has — so the private default is also the safe one.

### `gogl profile import [file]`

| Flag | Meaning |
|---|---|
| `--wireless` | Apply the wireless sections too; needs a wired session |
| `--force` | Allow a subnet move while reservations exist |
| `--dry-run` | Show what would change without changing it |

Applied in a fixed order, each step where it is because doing it later fails:

1. **Domain** — reservation writes are refused without one
2. **Network** — reservations must be inside the subnet before they are written
3. **Reservations**, then **DNS names**
4. **Wireless**, opt-in, last — most likely to be refused, and a refusal there must not
   undo the addressing

**If the profile's subnet differs from the router's, the run stops after step 2** and prints
how to resume. The router changes address mid-write, so nothing after that is reachable from
the same session. Reporting success for a third-applied profile would be a lie.

A **model mismatch warns rather than fails.** Addresses and names are portable; wireless is
not — interface names, radio names, channel lists and hardware modes are per-device and
per-regulatory-domain. Interfaces the target lacks are reported and skipped.

Idempotent: a second run reports nothing to do.

Unknown fields in a profile are an **error**. A file from a newer gogl may carry a section
this build would silently drop, and silently dropping part of a network is worse than
refusing the file.

---

## `gogl system`

### `gogl system info`

```console
$ gogl system info
MODEL      sft1200
FIRMWARE   4.3.28
ENDPOINT   http://192.168.8.1:80/rpc
```

---

## `gogl config`

gogl's own configuration. The only area acting on your machine rather than a router.

### `gogl config show`

Where the file is, whether it exists, and what it resolves to.

### `gogl config routers`

```console
$ gogl config routers
NAME            HOST         DOMAIN       PASSWORD
home (default)  192.168.8.1  lab.example  command
travel          192.168.8.1  -            environment or prompt
```

### `gogl config init`

Writes a commented starting file. `--force` overwrites an existing one.

---

## `gogl completion`

```bash
gogl completion bash > /etc/bash_completion.d/gogl
gogl completion zsh  > "${fpath[1]}/_gogl"
gogl completion fish > ~/.config/fish/completions/gogl.fish
```

---

## The three guards

Each refuses a well-formed request because the state is wrong, and each exits **3**.

### A reservation write requires a configured domain

`ErrDomainNotSet`. A binding with no name is an address nothing can find, and nothing in the
router's UI flags it as incomplete. Making the domain a precondition turns a silent omission
into an error where the mistake is.

Reads and deletes are ungated; only writes that create addressing are.

```bash
gogl lan dns set --domain lab.example      # once per router
```

### Moving the LAN subnet requires no reservations

`ErrReservationsExist`, unless `--force`.

The original reason was wrong and is worth knowing. It assumed the firmware would strand
reservations outside the new subnet. **It does not — it silently renumbers them, preserving
host parts.** Twenty-seven reservations moved from `192.168.2.x` to `192.168.8.x` with no
prompt.

The guard is kept on narrower grounds: an unannounced rewrite of every reservation is a
large side effect of an address flag, and the behaviour is only known for a same-size
subnet. A narrower netmask, or a move that lands addresses inside the new pool, is untested
— twenty of those twenty-seven landed in the pool.

A pool-only change is never guarded, because nothing moves.

### A wireless write requires a wired session

`ErrWirelessSession`. Changing an SSID, passphrase or channel drops every client on that
radio. Unlike a LAN renumber there is no new address to reconnect at: the network the session
was using stops existing under that name. Recovery means ethernet or the reset pin.

gogl finds the address it reaches the router from, looks it up in the router's own client
list, and reads the interface the firmware reports. Anything but `cable` is refused.

`--yes` does **not** override this. The prompt is about intent; the guard is about whether
you can still reach the device afterward, and no amount of intent makes a severed session
recoverable. A session from off-LAN is allowed, since no radio here carries it.

---

## Two things the firmware does that will surprise you

### A reservation does not create a DNS record

The admin panel shows **LAN → Address Reservation** with Name, MAC and IP, so the Name reads
as a DNS name. It is not. Tested twice on 4.3.28: a binding for an absent MAC never
resolved, and a binding for a *present, actively leased* device under a new name still did
not resolve — while that device's own lease hostname kept resolving throughout.

Three separate things:

| What | Mechanism | Who controls it |
|---|---|---|
| The address a device receives | static bind | gogl |
| A name you choose | the router's host file | gogl |
| A name a device announces | its DHCP lease hostname | the client |

`gogl` writes both of the first two from one host declaration, so this is invisible in normal
use. What it means in practice: names are batched into one host-file write, and the two
tables can drift if something else edits them — which the next import repairs.

### The firmware renumbers reservations when you move the LAN

Documented under [the second guard](#moving-the-lan-subnet-requires-no-reservations). Usually
what you want, entirely unannounced.

---

## Recipes

### Copy a network from UniFi to a travel router

```bash
# At home, against the controller:
gofips -H 192.168.4.1 -k --get > home.hosts

# On the travel router, in order:
gogl lan dns set --domain lab.example
gogl lan set --ip 192.168.4.1 --mask 255.255.255.0 \
             --pool-start 192.168.4.100 --pool-end 192.168.4.149
# reconnect at the new address
gogl -H 192.168.4.1 lan reservations import home.hosts --dry-run
gogl -H 192.168.4.1 lan reservations import home.hosts
```

### Clone one router onto another

```bash
gogl --router home   profile export --with-keys > home.json
gogl --router travel profile import home.json --wireless
```

### Rename the WiFi for a site

```bash
gogl radio list                                  # find the interface names
gogl wifi set --band 2.4 --ssid site-2g
gogl wifi set --band 5   --ssid site-5g
gogl wifi set --band 5   --passphrase            # prompts
```

Over ethernet.

### Move off a congested channel

```bash
gogl radio list                                  # what this radio offers
gogl clients list -w                             # who is actually associated
gogl radio set --band 5 --channel 149
```

### Find what is worth reserving

```bash
gogl clients list --reserved      # who has a reservation and who does not
gogl lan leases                   # what the router is handing out dynamically
```

### Version-control a network

```bash
gogl lan reservations export > lab.hosts
git add lab.hosts && git commit -m "capture lab addressing"
gogl lan reservations import lab.hosts            # idempotent; reports all skips
```

---

## Troubleshooting

### `Access denied` on every call, after a burst of commands

The firmware's brute-force protection. **Every `gogl` invocation performs a full login** —
two challenge calls and a login — and the router counts them. A script making dozens of
separate calls trips this.

gogl now says so explicitly, because a denial on the *challenge* call cannot mean a wrong
password: challenge carries only a username.

Wait a few minutes. Pace scripted calls, and prefer one `profile import` over many
individual writes.

### `refusing to change wireless over a wireless session`

Working as intended. Connect over ethernet. `--yes` will not override it.

### `DNS domain is not configured`

Run `gogl lan dns set --domain <domain>` once per router.

### `reservations exist: N reservation(s) present`

You are moving the LAN subnet with reservations in place. Either
`gogl lan reservations clear` first, or pass `--force` to accept the firmware rewriting them.

### A flag the documentation mentions does not exist

A stale binary shadowing a fresh one. `gogl --version` reports the revision and whether the
tree was dirty; `make uninstall && make install` clears both the current and legacy install
locations.

### `HTTP 418` from the OUI download

IEEE's bot filter. `gogl` sends an explicit User-Agent to avoid it; if you see this, something
is rewriting the header. A stale cache still works. With no cache and no network,
`clients list` fails rather than printing a table of blanks.

### Passphrases appear in output

Only with `--show-key`, or in `profile export --with-keys`. Both are explicit. `--show-key`
keeps a key out of your scrollback by default; it is not an access control, since anyone with
the admin password can read passphrases over the API.

---

## Related documents

- [`../README.md`](../README.md) — the problem gogl solves, and the Go library
- [`../VISION.md`](../VISION.md) — requirements and the critical rules
- [`DESIGN-V2.md`](DESIGN-V2.md) — the single-binary command tree and its reasoning
- [`DESIGN.md`](DESIGN.md) — library architecture
- [`../GL_INET_4X_API_DOCUMENTATION.md`](../GL_INET_4X_API_DOCUMENTATION.md) —
  hardware-verified API surface, and where it disagrees with its own documentation
- [`api/`](api/README.md) — all 43 groups and 313 methods, each marked verified, absent or
  untested
- [`../TODO.md`](../TODO.md) — what still needs capturing against hardware
