<img src="images/gl-gopher-antennas-transparent.png" width="33%" alt="gogl logo">

# gogl - Go GL.iNet Travel Router Client

Programmatic control of GL.iNet travel routers running firmware 4.x, targeting the
**GL-SFT1200 (Opal)**.

Disclaimer: This works for me — that's the entire guarantee. Built with AI in the loop, so check your own biases before you love it or hate it on principle. Use at your own risk, fork freely, and don't @ me when it explodes. (But do drop me a note if it helps — pay it forward.)

> **Status: working against real hardware.** Verified on a GL-SFT1200 (Opal)
> reporting model `sft1200`, firmware **4.3.28**. All three utilities read and write
> the live device; every endpoint gogl uses is confirmed, and the full API is
> documented in [`docs/api/`](docs/api/README.md).
>
> One design assumption was wrong and has been corrected: a DHCP reservation does
> **not** create a DNS record. Names work, but through the router's host file
> instead — see [Addresses and names are separate](#addresses-and-names-are-separate).

---

## The problem

A network's addressing is the part your devices actually depend on. A robot, a signage
player, or a build machine configured to reach `nas` at `192.168.4.13` stops working the
moment it lands on a network that disagrees.

Take that kit somewhere else — a customer site, a trade show, a hotel room, a test bench —
and you have to rebuild the addressing by hand, one DHCP reservation at a time, through a
web UI. Then do it again the next time, and again after a factory reset.

`gogl` removes that. Its sibling [`gofi`](https://github.com/emergingrobotics/gofi) exports
a UniFi site's fixed-IP bindings as an ISC DHCP host file. `gogl` imports that same file into
a GL.iNet travel router, so the pocket router hands out the same addresses to the same MAC
addresses as the network you copied — and can serve the same names.

The file is plain text, so it diffs, reviews, and lives in git. Your kit's network becomes
reproducible from a commit rather than from memory.

Addresses and names come from two different mechanisms on this firmware. `goglps` writes both
in one command, but the split explains a couple of things that would otherwise look odd — see
[Addresses and names are separate](#addresses-and-names-are-separate).

```bash
# At home, against the UniFi controller:
gofips -H 192.168.4.1 -k --get > home.hosts

# Later, against the travel router:
goglps --domain lab.example      # once per router
goglps --set home.hosts
```

### Addresses and names are separate

The obvious assumption here is wrong, and it took hardware testing to find out.

GL.iNet's admin panel shows **LAN → Address Reservation** with Name, MAC and IP fields, so it
reads as though the Name becomes a DNS name. **It does not.** Tested twice on 4.3.28: a
reservation for an absent MAC never resolved, and a reservation for a *present, actively
leased* device under a new name still didn't resolve — while that device's own lease hostname
kept resolving the whole time.

So there are three separate things:

| What | Mechanism | Who controls it |
|------|-----------|-----------------|
| The address a device receives | static bind (reservation) | **gogl** — `goglps` |
| A name you choose | the router's **host file** | **gogl** — `goglps --domain` and host entries |
| A name a device announces about itself | its DHCP lease hostname | the client |

Names come from `dns.set_host`, which writes the router's `hosts(5)` file — dnsmasq answers
from it, verified by writing an entry and resolving both the bare name and an arbitrary FQDN
against the router.

gogl owns a delimited block in that file and never touches anything outside it. The file also
holds the loopback and IPv6 lines the router itself needs.

**`goglps` writes both, and you do not have to think about it.** One host declaration becomes a
static bind *and* a host-file entry; `--set`, `--add`, `--del`, `--prune` and `--clear` each
keep the two in step. What the split means in practice is only this:

- Names are batched. `dns.set_host` takes the whole file as one string, so a `--set` run makes
  one write for every name rather than one per name.
- The two can drift if something else edits them — a hand edit in the admin panel, or a run
  that died between the two writes. A binding with no name is repaired by the next `--set`,
  which diffs the names independently rather than skipping an entry whose binding already
  matches. `goglps --get` reports what it finds in each.

### Set the domain first

The firmware exposes no way to set dnsmasq's domain — no endpoint among all 313 has a
`domain` field, and `/ubus`, which could do it, returns 404 on this model. So gogl carries its
own domain, stored inside its block in the host file, and writes fully-qualified names. Any
suffix works, because host-file entries are literal name-to-address mappings.

```bash
goglps --domain lab.example
```

**Reservation writes are refused until you do this.** A reservation alone produces an address
with no name and nothing to signal that was unintended; making the domain a precondition
forces the pairing to be deliberate.

`gofips`'s `--keep-dns` has no analogue, and here it would mean something: since the two are
separate objects, keeping a name while dropping its reservation is expressible. `gogl` does not
offer it, because a name resolving to an address the router no longer reserves is a trap rather
than a feature. `--del` and `--clear` remove both.

---

# Part 1: Command-line programs

Four tools. The first three mirror `gofi`'s one-for-one, so that knowing one set means
knowing the other. `goglcfg` has no counterpart: reproducing a whole network spans all
three and belongs in none of them.

| gogl | gofi counterpart | What it does | Writes? |
|------|------------------|--------------|---------|
| `goglps` | `gofips` | DHCP reservations, in ISC DHCP format | **yes** |
| `goglnet` | `gofinet` | LAN address, DHCP pool, resolvers, wireless | **yes**, with `--set-*` |
| `goglmac` | `gofimac` | Connected clients with IEEE OUI manufacturer lookup | no |
| `goglcfg` | — | A whole network as a JSON profile: capture and apply | **yes**, on `--set` |

**What gogl writes:** reservations, host-file entries and the DNS domain (`goglps`); the LAN
address and DHCP pool, wireless identity and radio tuning (`goglnet --set-*`). **What it does
not:** VLANs, firewall, VPN. Set those in the GL.iNet admin panel.

Three rules protect you from the states that are painful to recover from:

- **A reservation write needs a domain first.** Otherwise you get addresses with no names and
  no hint that was accidental.
- **Moving the LAN subnet is refused while reservations exist**, unless `--force`. The firmware
  silently rewrites every reservation into the new subnet, keeping host parts — usually what you
  want, but unannounced. Changing only the DHCP pool is never guarded, since nothing moves.
- **An SSID change is refused over a WiFi session.** It would drop the session with no address
  to reconnect at. Connect over ethernet. See [`goglnet`](#goglnet--network-and-wireless).

## Prerequisites

- A GL.iNet router on **firmware 4.x**. GL.iNet replaced its REST API with JSON-RPC at
  firmware 4.0, and `gogl` speaks only JSON-RPC. Firmware 3.x will not work.
- The admin password you set during the router's first-run setup. There is no default.
- Go 1.22 or later, to build.
- Network reachability to the router's LAN address.

Verified against a GL-SFT1200 (Opal) reporting firmware **4.3.28**. Other GL.iNet 4.x models
share the API, but note that this one is a reduced build: six endpoints GL.iNet documents are
absent on it. See [`docs/api/`](docs/api/README.md), where each method is marked.

## Quick start

```bash
git clone https://github.com/emergingrobotics/gogl
cd gogl
make build          # builds to bin/
make install        # or install to ~/bin
```

Set your credentials. The password comes from the environment rather than a flag, so it
never lands in a process listing:

```bash
export GL_ROUTER_IP=192.168.8.1
read -rsp 'router password: ' GL_PASSWORD; export GL_PASSWORD; echo
```

`read -rsp` keeps it out of your shell history too. If you would rather export it directly,
use single quotes so a `$` or backtick in the password is not mangled:

```bash
export GL_PASSWORD='your-router-admin-password'
```

See what you are working with:

```bash
goglnet
```

```
MODEL       gl-sft1200
FIRMWARE    4.3.28
LAN         192.168.8.0/24  (255.255.255.0)
DHCP        enabled
POOL        192.168.8.100 - 192.168.8.249  (150 addresses)
LEASE       12h
DOMAIN      lan
DNS         192.168.8.1
RESERVED    0
AVAILABLE   103
```

`AVAILABLE` counts addresses that are neither in the DHCP pool, nor reserved, nor the router
itself — the ones safe to hand out statically.

Find the MAC address of something you want to pin:

```bash
goglmac
```

```
aa:bb:cc:dd:ee:01	192.168.8.101	nas	    Synology Incorporated
aa:bb:cc:dd:ee:02	192.168.8.102	laptop	Dell Inc.
02:1a:2b:3c:4d:5e	192.168.8.103	unknown	randomized
```

The manufacturer comes from `goglmac`'s own copy of the IEEE OUI database, never from
whatever the router reports — on an OpenWrt 18.06 base, the router's table is likely years
stale. A locally-administered address shows `randomized` rather than `unknown`, because it
will never be in the IEEE database and is a poor thing to pin an address to.

Pin it:

```bash
goglps --add 'host nas {
    hardware ethernet aa:bb:cc:dd:ee:01;
    fixed-address 192.168.8.13;
}'
```

```
Created: nas aa:bb:cc:dd:ee:01 192.168.8.13
```

The router now always hands `192.168.8.13` to that MAC. Confirm with:

```bash
goglps --get
```

It does **not** create a DNS record for `nas`. If that device announces `nas` as its own
DHCP hostname, the router will resolve `nas.lan` from the resulting lease — but that is the
client's doing, not the reservation's. Check what the router actually knows:

```bash
goglnet          # the resolvers the router advertises
nslookup <name> <router-ip>
```

## Configuration

### Environment variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GL_PASSWORD` | yes | Router admin password. Note the full spelling: `GL_PASSWD` is a common slip, and the tools will name it if you set that instead. |
| `GL_ROUTER_IP` | no | Router address, used when `-H` is absent |
| `GL_USERNAME` | no | Router admin username (default `root`) |
| `XDG_DATA_HOME` | no | Base for `goglmac`'s OUI cache (default `~/.local/share`) |

### Connection flags

Shared by all three tools.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--host` | `-H` | `$GL_ROUTER_IP` | Router address |
| `--port` | `-p` | `80` | Router port |
| `--https` | | `false` | Use HTTPS instead of HTTP |
| `--secure` | `-k` | `false` | Under `--https`, enforce TLS certificate verification |

Note the port default is **80**, not 443. GL.iNet firmware serves its admin interface over
plain HTTP; that is a difference from the UDM Pro, so a habit carried over from `gofi` will
fail here. Passing `-k` without `--https` is an error rather than a silent no-op — there is
no certificate to verify over HTTP.

There is no `--site` flag. GL.iNet routers have no equivalent of a UniFi site.

## Usage

Three utilities. Every one takes the same connection flags and reads the same environment
variables; see [Configuration](#configuration).

| Tool | Reads | Writes |
|------|-------|--------|
| [`goglps`](#goglps--reservations-and-dns-names) | reservations, DNS names, the domain | all three |
| [`goglnet`](#goglnet--network-and-wireless) | LAN, DHCP pool, leases, wireless | LAN address and pool, wireless identity and tuning |
| [`goglmac`](#goglmac--connected-clients) | connected clients, OUI vendors | nothing |
| [`goglcfg`](#goglcfg--whole-network-profiles) | everything the other three do | everything the other three do |

Every write takes `--dry-run`, and a dry run performs every check the real write performs —
including the refusals. A preview that approves what the write would reject is worse than no
preview.

---

### `goglps` — reservations and DNS names

Manages static DHCP bindings and the DNS names that go with them, in ISC DHCP host-declaration
format. One declaration becomes two writes, because the firmware stores the address and the name
in separate tables; `goglps` keeps them in step. See
[Addresses and names are separate](#addresses-and-names-are-separate).

Exactly one mode flag per invocation.

| Mode | Short | What it does |
|------|-------|--------------|
| `--get` | `-g` | Export every reservation to stdout in ISC DHCP format |
| `--set [file]` | `-s` | Import declarations from a file, or stdin if omitted |
| `--add '<fragment>'` | `-a` | Add one host from a declaration fragment, or stdin |
| `--del` | `-d` | Delete one host, identified by `--name`, `--mac` or `--ip` |
| `--domain <domain>` | | Set the DNS domain. Required before any reservation write |
| `--clear` | | Delete **every** reservation and DNS name |

| Modifier | Short | Applies to | What it does |
|----------|-------|-----------|--------------|
| `--dry-run` | | all writes | Show what would change, change nothing |
| `--prune` | `-P` | `--set` | Also delete reservations and names on the router but absent from the file |
| `--force` | `-f` | `--add`, `--del`, `--clear` | Proceed past conflicts; skip the confirmation prompt |
| `--name` | `-n` | `--del` | Identify the target by hostname |
| `--mac` | `-m` | `--del` | Identify the target by MAC |
| `--ip` | `-i` | `--del` | Identify the target by address |

```bash
goglps --domain herlein.me                 # once per router, before any write
goglps --get > player-test.hosts           # export
goglps --set player-test.hosts --dry-run   # preview an import
goglps --set player-test.hosts             # apply
goglps --set player-test.hosts --prune     # apply, and delete what the file omits

goglps --add 'host nas { hardware ethernet aa:bb:cc:dd:ee:01; fixed-address 192.168.8.13; }'
goglps --del --name nas
goglps --clear                             # everything, prompts unless --force
```

`--set` is idempotent: run it twice and the second run reports everything skipped and leaves the
host file byte-identical. That is what makes a host file usable as a checked-in description of a
network.

Two refusals you will meet:

- **`--set`, `--add`: domain not configured.** Run `--domain` first. A reservation with no name
  is an address nothing can find, and nothing in the router's UI marks it as incomplete.
- **Any address outside the LAN subnet.** Either renumber the router with `goglnet --set-ip`, or
  renumber the file. See [Match the subnet](#match-the-subnet).

---

### `goglnet` — network and wireless

With no flags, reports the LAN, the DHCP pool, the reservation count, and every wireless radio
with what it supports. With any `--set-*` flag, writes.

#### Reporting

| Flag | Short | What it does |
|------|-------|--------------|
| *(none)* | | Report the LAN, pool, reservation counts, and every radio |
| `--json` | `-j` | JSON instead of text |
| `--show-key` | | Print WiFi passphrases instead of masking them |

```bash
goglnet
goglnet -j
goglnet --show-key
```

If any reservation sits inside the DHCP pool, the report says so and lists them:

```
RESERVED     27
  IN POOL    20  (honored by dnsmasq, excluded from the pool)
AVAILABLE    96

20 reservation(s) fall inside the DHCP pool 192.168.8.100-192.168.8.249:
  192.168.8.208   3c:6a:d2:85:1c:d5   ganymede
  ...
```

Nothing is broken — dnsmasq honors a static bind inside the dynamic range and excludes that
address from allocation. But it is reported because it arises silently: a LAN renumber moved 20
of 27 reservations into the pool on a real device, since the firmware rewrites host parts without
considering where the pool is. It is also the explanation for an `AVAILABLE` count that otherwise
looks wrong, those addresses being neither dynamically assignable nor counted as free.

The wireless part reports each radio's selectable values, because the tuning flags are unusable
without them — they differ per radio and per regulatory domain:

```
5G radio radio1: channel 44, 20, 11a/n/ac, power Max
  channels:    36, 40, 44, 48, 149, 153, 157, 161, 165, or 0 for auto
  bandwidths:  20, 40, 80
  hw modes:    11ac, 11n/ac, 11a/n/ac
  encryptions: none, psk2, psk-mixed, sae, sae-mixed
  power:       Max, High, Medium, Low
```

Bandwidths are widths in MHz — `20`, `40`, `80`, and `auto` where the radio offers it — not the
`HT20`/`VHT80` form GL.iNet's own API description implies. Hardware modes are slash-joined
combinations (`11a/n/ac`), not bare `11ac`. `sae` is WPA3. These are what the device reports;
an earlier version of this README printed the description's values, which are wrong.

#### Changing just the DHCP pool

```bash
goglnet --set-start 192.168.8.50 --set-end 192.168.8.150 --dry-run
goglnet --set-start 192.168.8.50 --set-end 192.168.8.150
```

The address and netmask are read from the device, so you do not restate them. Nothing moves:
the router keeps its address, **the session survives**, and no reservation is touched. This is
never refused, even with reservations present.

#### Moving the LAN address

All four flags together here, because a pool from the old subnet cannot be valid in a new one
and guessing a replacement would be worse than asking.

| Flag | What it sets |
|------|--------------|
| `--set-ip` | LAN address, e.g. `192.168.8.1` |
| `--set-mask` | Netmask |
| `--set-start` | DHCP pool start |
| `--set-end` | DHCP pool end |
| `--set-interface` | `lan` (default) or `guest` |
| `--force` | Proceed despite existing reservations |
| `--dry-run` | Show the change and any refusal, without writing |

```bash
goglnet --set-ip 192.168.8.1 --set-mask 255.255.255.0 \
        --set-start 192.168.8.100 --set-end 192.168.8.249 --dry-run
```

**The router moves mid-call, so the session drops.** That is success, not failure; gogl says so
and tells you the address to reconnect at.

**Refused while reservations exist, unless `--force`.** The firmware silently rewrites every
reservation into the new subnet, preserving host parts — `192.168.2.10` becomes `192.168.8.10`.
That is usually what you want, which is why `--force` exists. It is refused by default because
the rewrite is unannounced, and because it is only known to work for a same-size subnet: a
netmask change, or a move that lands addresses inside the new pool, is untested.

A netmask change counts as moving the subnet even when the address stays put, so it is guarded
too.

#### Writing wireless identity, with `--iface`

| Flag | What it sets |
|------|--------------|
| `--set-ssid` | The network name, up to 32 characters |
| `--set-key` | WPA passphrase, 8 to 63 characters |
| `--set-encryption` | e.g. `psk2`, or `sae` for WPA3; validated against the radio's list |
| `--set-hidden=true\|false` | Whether the SSID is advertised |
| `--set-enabled=true\|false` | Whether the interface is up |

```bash
goglnet --iface default_radio0 --set-ssid player-test
goglnet --iface default_radio0 --set-key 'a better passphrase'
goglnet --iface guest2g --set-enabled=true --set-hidden=false --set-ssid player-guest
```

#### Writing radio tuning, with `--device`

| Flag | What it sets |
|------|--------------|
| `--set-channel` | Channel, or `0` for auto; validated against the radio's list |
| `--set-htmode` | Channel width in MHz: `20`, `40`, `80`, or `auto` where offered |
| `--set-hwmode` | Hardware mode, e.g. `11a/n/ac`; run `goglnet` for the radio's list |
| `--set-txpower` | `Max`, `High`, `Medium` or `Low` |

```bash
goglnet --device radio1 --set-channel 149
goglnet --device radio1 --set-htmode 80 --set-txpower Low
goglnet --device radio0 --set-channel 0      # auto
```

Both scopes combine in one invocation; gogl issues the two calls the firmware needs and reports
which applied.

#### Four things about the wireless flags

**They take an explicit value: `--set-hidden=false`, not a bare `--set-hidden`.** A partial
update has to distinguish "set this to false" from "leave it alone", and Go's flag package
treats both as `false`. A bare boolean flag meaning true is exactly how `--set-enabled` ends up
disabling something. Same reason `--set-channel=0` works: 0 means auto, so it cannot double as
"unset".

**Only the fields you name are sent.** Setting a passphrase leaves the SSID, encryption, and
enabled state exactly as they were.

**Ethernet only.** Every wireless write is refused when the session arrives over WiFi:

```
gogl: refusing to change wireless over a wireless session: this session is on 5G
  changing wireless would drop it, with no address to reconnect at
  connect over ethernet and try again
```

gogl finds the address it reaches the router from in the router's own client list and reads the
interface the firmware reports for it. `--yes` does **not** override this — the prompt is about
intent, the guard is about whether you can still reach the device afterward. A session from
off-LAN is allowed, since no radio here carries it. Tuning is gated too: retuning drops the
radio's clients just as thoroughly as renaming it.

**They prompt.** Every wireless write shows the before and after and asks `y/N`. `--yes` skips
it for scripts. Unlike `--del`, a non-terminal stdout does not imply consent: there the worst
case is a lost reservation, here it is a device you have to walk to. DFS channels get an extra
warning — the radio must vacate one if it detects radar, taking every client with it for the
minutes it spends re-scanning. A poor choice for kit that has to come up reliably in an
unfamiliar building.

---

### `goglmac` — connected clients

Read-only. Lists what is connected, with IEEE OUI manufacturer lookup done independently of the
router.

| Flag | Short | What it does |
|------|-------|--------------|
| `--all` | `-a` | Every client (default) |
| `--wifi` | `-w` | Wireless clients only |
| `--wired` | `-e` | Wired clients only |
| `--reserved` | `-r` | Mark which clients hold a reservation |
| `--json` | `-j` | JSON instead of text |

```bash
goglmac                 # everything
goglmac -w              # who is on the radios
goglmac -r              # which clients are worth reserving
```

The OUI database is downloaded from IEEE and cached. A download failure is not fatal if a cache
exists; if there is no cache, `goglmac` exits 1 rather than printing a table of blanks that
looks like every device is from an unknown vendor.

---

### `goglcfg` — whole-network profiles

Captures a router's reproducible configuration to JSON, and applies one back — onto the same
router or a different one.

```bash
goglcfg --get > lab.json                 # capture
goglcfg --set lab.json --dry-run         # preview
goglcfg --set lab.json                   # apply
goglcfg --set lab.json --wireless        # apply the wireless sections too
```

| Flag | Short | What it does |
|------|-------|--------------|
| `--get` | `-g` | Capture a profile to stdout |
| `--set [file]` | `-s` | Apply a profile from a file, or stdin |
| `--with-keys` | | With `--get`, include WiFi passphrases in cleartext |
| `--wireless` | | With `--set`, apply the wireless sections; needs a wired session |
| `--dry-run` | | Show what would change, change nothing |
| `--force` | | Allow a subnet move while reservations exist |

#### What a profile is, and is not

**It is not a router image.** It carries what defines a *network* — LAN address and DHCP pool,
reservations, DNS names and domain, wireless identity and radio tuning — and omits everything
identifying a particular unit: the router's own MAC, serial number, uptime, lease state.

That omission is the point. Those fields are exactly what makes a full config dump useless on a
second router. Client MAC addresses *are* included, since a reservation is a MAC-to-IP binding
and a profile without them would reproduce nothing.

Every section comes from an endpoint verified against hardware. That is why the file is small:
the API exposes 110 getters and 23 are verified, and a profile built on the rest would be
guesswork. Lease time, upstream DNS servers, firewall, VPN and VLANs are absent because no
verified endpoint writes them.

#### Passphrases are omitted by default

```bash
goglcfg --get > lab.json              # no keys; safe to commit
goglcfg --get --with-keys > lab.json  # keys in cleartext
```

An omitted key is not an empty key. On apply, a missing key is simply not written, which leaves
whatever the target router already has — so a key-less profile is safe rather than destructive.
That relies on `wifi.set_config` leaving unmentioned fields alone, which is verified on hardware.

#### The apply order, and why a subnet move stops the run

`--set` applies in a fixed order, each step placed where it is because doing it later fails:

1. **Domain**, because reservation writes are refused without one.
2. **Network**, because reservations must be inside the subnet before they are written.
3. **Reservations**, then **DNS names**.
4. **Wireless** last, opt-in — it is the step most likely to be refused, and failing there should
   not undo the addressing having been applied.

**If the profile's subnet differs from the router's, the run stops after step 2.** The router
changes address mid-write, so nothing after it is reachable from that session:

```
network: 192.168.8.1/255.255.255.0 -> 192.168.4.1/255.255.255.0

the router has moved to 192.168.4.1. The rest of the profile was not applied.
resume with:  goglcfg -H 192.168.4.1 --set lab.json
```

Re-run at the new address and it completes. Reporting success for a run that wrote a third of the
profile would be a lie. A pool-only difference is not a move, so that case runs straight through.

The subnet move is refused while reservations exist unless `--force`, same as `goglnet`.

#### Applying to a different model

A model mismatch warns rather than fails:

```
warning: profile is from "mt3000", this router is "sft1200"
  addresses and names are portable; wireless interface and radio names
  may not exist on this model, and will be reported and skipped
```

Addresses and names carry over cleanly. Wireless is where models diverge — interface names, radio
names, channel lists and hardware modes are all per-device and per-regulatory-domain. Interfaces
the target does not have are reported and skipped, not treated as errors.

`--set` is idempotent: a second run reports nothing to do and leaves the host file byte-identical.

---

## Common tasks

### Copy a whole network from UniFi

```bash
# On the UniFi side, at home.
gofips -H 192.168.4.1 -k --get > home.hosts
git add home.hosts && git commit -m "capture lab network"

# On the travel router, in order:
goglps --domain lab.example                    # 1. domain first, or writes are refused
goglnet --set-ip 192.168.4.1 --set-mask 255.255.255.0 \
        --set-start 192.168.4.100 --set-end 192.168.4.149   # 2. match the subnet
goglps --set home.hosts --dry-run              # 3. preview
goglps --set home.hosts                        # 4. apply
```

Step 2 moves the router to a new address, so **the session drops** and you reconnect at the
new one. It is refused if reservations already exist — clear them with `goglps --clear` first.
Step 1 only needs doing once per router.

`--set` is idempotent: run it twice and the second run reports everything as skipped. That is
what makes a host file usable as a checked-in description of a network.

By default, reservations on the router but absent from your file are left alone and counted.
`--prune` deletes them; `--clear` deletes everything.

### Managing the domain and clearing out

```bash
goglps --domain lab.example      # set it, or change it and requalify existing names
goglps --clear --dry-run         # list what would go
goglps --clear                   # delete ALL reservations AND DNS names
```

Changing the domain rewrites every managed host entry to the new suffix, so resolution does
not split between two domains.

`--clear` removes the DNS names as well as the reservations, because they are one intent stored
in two places. Leaving the names behind would strand records pointing at addresses the router
no longer reserves — and since clearing is what unblocks a renumber, those names would then be
answering for a subnet that no longer exists. The domain survives a clear: it is configuration,
not content.

### Rename the WiFi to match the site

```bash
goglnet                                                   # find the interface name
goglnet --iface default_radio0 --set-ssid player-test --dry-run
goglnet --iface default_radio0 --set-ssid player-test
```

Must be over ethernet, and it prompts. Full flag reference in
[`goglnet`](#goglnet--network-and-wireless).

### Move off a congested channel

```bash
goglnet                                    # see the channels this radio supports
goglnet --device radio1 --set-channel 149
```

The most common reason to touch tuning on a travel router: the site's existing WiFi is sitting
on the channel yours picked. Combine with `goglmac -w` to see who is actually associated.

### Export, edit, push back

```bash
goglps --get > current.hosts
$EDITOR current.hosts
goglps --set current.hosts
```

### Delete a reservation

```bash
goglps --del --name nas
goglps --del --mac aa:bb:cc:dd:ee:01
goglps --del --ip 192.168.8.13
```

This removes the DHCP reservation. It does not affect DNS, because the reservation never
provided any.

### List clients by connection type

```bash
goglmac --wired
goglmac --wifi
goglmac -j -r          # JSON, marking which clients already have a reservation
```

## Before you import: two things to get right

`gogl` will not change the router's own network configuration, which means two settings are
yours to get right in the GL.iNet admin panel. Both cause failures that are hard to read
back from the symptom.

### Match the subnet

A reservation at `192.168.4.10` is meaningless on a router whose LAN is `192.168.8.0/24`.
Either set the router's LAN address under **LAN → Router IP Address** in the admin panel, or
let gogl do it:

```bash
goglps --clear                                  # required first: see below
goglnet --set-ip 192.168.4.1 --set-mask 255.255.255.0 \
        --set-start 192.168.4.100 --set-end 192.168.4.149 --dry-run
goglnet --set-ip 192.168.4.1 --set-mask 255.255.255.0 \
        --set-start 192.168.4.100 --set-end 192.168.4.149
```

Either way your management session drops when the router moves, and you reconnect at the new
address. `goglnet` refuses the change while reservations exist, because renumbering underneath
the firmware silently rewrites every one of them into the new subnet instead.

`goglps` refuses any reservation outside the current subnet rather than writing something
inert, and tells you both ways to fix it:

```
error: subnet mismatch
  host file:   192.168.4.0/24 (38 of 38 entries)
  router LAN:  192.168.8.1/24

Resolve by either:
  - Setting the router's LAN address to 192.168.4.1/24 in the GL.iNet admin panel
    (LAN -> Router IP Address), then re-running. Your management session will drop
    and you will need to reconnect at the new address.
  - Renumbering the host file into 192.168.8.0/24 before re-running.
```

### Do not copy your home network's DNS servers

This is the trap. If your home network hands out a Pi-hole at `192.168.4.5` and you
configure the travel router to advertise the same resolver, every client at the customer
site is told to use a DNS server that does not exist there. It presents as "the internet is
broken," not as "DNS is misconfigured."

Leave the router advertising **itself**. dnsmasq answers your reservation names locally and
forwards everything else to whatever resolver the WAN handed it — which works in a hotel, at
a customer site, and on a bench with no uplink at all. Public resolvers like `1.1.1.1` are
safe if you want them, since they work anywhere.

## Limitations

Deliberate, and documented rather than worked around.

| Not supported | Why |
|---------------|-----|
| Wireless writes over a wireless session | Refused, because applying one severs the session with no address to reconnect at. Connect over ethernet. |
| Radio tuning beyond channel, bandwidth, hardware mode and power | Nothing else is exposed by `wifi.set_config`. |
| VLANs | The Opal exposes no VLAN configuration through its API. Reaching them means LuCI, UCI, and `swconfig` on a SiFlower switch, which is both outside this tool's rules and unreliable in practice. |
| Changing lease time or upstream DNS servers | No verified endpoint sets them. The LAN address and DHCP pool *are* writable, with `goglnet --set-*`. |
| Setting dnsmasq's `domain` / `local` / `expandhosts` | No endpoint exposes them, and `/ubus` is 404 on this model. gogl carries its own domain and writes FQDNs into the host file instead, which works for any suffix. |
| DNS names via reservations | Reservations do not create DNS records. Use the host file — `goglps --domain` plus host entries. |
| DFS channels are not offered | `wifi.get_config` lists nine 5GHz channels; the driver supports twenty-five. The sixteen missing are the DFS ones, and `dfs_support: false` reports GL.iNet's policy rather than the radio's capability. `goglnet --set-channel 52` is refused because the firmware will not accept it, not because the hardware cannot. |
| Monitor mode and packet injection | Not reachable over the API, and not usable on this hardware anyway: the driver advertises monitor mode, refuses to add an interface while the APs are up, and hung the device when one was brought up with them down. Recovered by power cycle. Use a USB adapter. |
| VPN, firewall, port forwarding | Out of scope for v1. |
| Firmware 3.x | Different API entirely (REST, not JSON-RPC). |
| Hostnames containing `_` | Not a legal DNS label character. `gofips` permits it; `goglps` rejects it rather than silently rewriting it to `-`. This is the one case where a `gofips` file will not import unchanged. |
| `goglmac` with no cached OUI data and no internet | Hard failure, matching `gofimac`. Run it once with a network connection and the 30-day cache covers you afterward. |
| An OUI download that returns HTTP 418 | IEEE's bot filter rejecting the client. `goglmac` sends an explicit `User-Agent` to avoid it; if you see this, something is rewriting the header. A stale cache still works. |
| Repeated failed logins locking you out | The firmware's brute-force protection. After roughly a dozen failures it refuses even a **correct** password for about ten minutes, and reports the remaining wait. gogl surfaces this as its own error with the countdown, because retrying makes it worse. |

One tool would close a real gap and is not in v1:

- **`goglrenumber`** — rewrite the network part of every address in a host file to a target
  subnet, preserving host parts. The alternative to re-IPing the router: move the file to
  the router's subnet instead of the router to the file's. A pure text transformation, so
  cheap to trust.

Re-IPing the router itself is covered: `goglnet --set-ip`. It stayed out of `goglps` so that
losing contact with the device is never a side effect of importing a host file.

---

# Part 2: Using gogl in your Go program

```bash
go get github.com/emergingrobotics/gogl
```

The library lives under `src/`, so the import path carries that suffix while the package
name is `gogl`:

```go
import (
    gogl "github.com/emergingrobotics/gogl/src"
    "github.com/emergingrobotics/gogl/src/types"
)
```

## Connecting

```go
client, err := gogl.New(gogl.Config{
    Host:     "192.168.8.1",
    Password: os.Getenv("GL_PASSWORD"),
})
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

`New` does not contact the router — the first call authenticates lazily, so constructing a
client is cheap and cannot fail on a network error. The client is safe for concurrent use.

`Config`'s zero value is safe: port 80, user `root`, TLS verification **on**. Note that the
library defaults to verifying certificates while the CLIs default to accepting self-signed
ones. That inversion is intentional: a library that is insecure at its zero value is a
hazard, and a CLI that cannot reach a self-signed device out of the box is useless. Set
`InsecureSkipVerify` explicitly if you need it.

## Services

```go
client.Network()      // LAN and DHCP configuration       (read/write)
client.Reservations() // static MAC-to-IP bindings        (read/write)
client.Hosts()        // DNS names via the host file      (read/write)
client.Clients()      // connected stations               (read-only)
client.System()       // model, firmware, uptime          (read-only)
```

No method takes a `site` parameter.

### Reading the network

```go
network, err := client.Network().Get(ctx)
if err != nil {
    log.Fatal(err)
}

subnet, err := network.Subnet()          // 192.168.8.0/24
fmt.Println(subnet, network.PoolSize())  // pool size in addresses
fmt.Println(network.DHCPLease)           // "12h"

inside, err := network.Contains(net.ParseIP("192.168.8.13"))
pooled, err := network.InDHCPPool(net.ParseIP("192.168.8.150"))
```

`InDHCPPool` is informational. dnsmasq honors a static lease whose address sits inside the
dynamic range and excludes it from dynamic allocation, so it is untidy rather than broken.
It *would* be a genuine conflict under ISC dhcpd, which is where the opposite intuition
comes from.

### Managing reservations

```go
reservations := client.Reservations()

created, err := reservations.Create(ctx, &types.Reservation{
    Name: "nas",
    MAC:  "AA:BB:CC:DD:EE:01",   // normalized to lowercase on write
    IP:   "192.168.8.13",
})

all, err := reservations.List(ctx)
one, err := reservations.GetByName(ctx, "nas")
one, err = reservations.GetByMAC(ctx, "aa:bb:cc:dd:ee:01")
many, err := reservations.GetByIP(ctx, "192.168.8.13")  // slice: duplicates are possible

_, err = reservations.Update(ctx, &types.Reservation{
    Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.20",
})

err = reservations.Delete(ctx, "aa:bb:cc:dd:ee:01")
```

MAC is the identity for `Update` and `Delete`: it is the only thing a client cannot change
about itself, and it is what dnsmasq keys the lease on.

### DNS names, and the domain

```go
hosts := client.Hosts()

// The domain is a precondition for writing reservations.
if err := hosts.SetDomain(ctx, "lab.example"); err != nil {
    log.Fatal(err)
}

// Creates "nas" and "nas.lab.example", both resolving to the address.
if err := hosts.Set(ctx, "nas", "192.168.8.13"); err != nil {
    log.Fatal(err)
}

file, err := hosts.Get(ctx)          // parsed, managed and unmanaged parts
ip, ok := file.Lookup("nas")         // either spelling works
err = hosts.Remove(ctx, "nas")
err = hosts.Clear(ctx)               // drops managed entries, keeps the domain
```

Everything outside gogl's block in the host file is preserved, including the loopback and IPv6
lines the router relies on.

### Writing the network

```go
// Moving the subnet. Refused while reservations exist; pass services.WriteForced to
// accept the firmware rewriting them into the new subnet.
err := client.Network().Set(ctx, &types.Network{
    Interface: types.InterfaceLAN,
    LANIP:     "192.168.2.1",
    Netmask:   "255.255.255.0",
    DHCPStart: "192.168.2.100",
    DHCPStop:  "192.168.2.149",
}, services.WriteGuarded)
```

The third argument is a `services.WriteMode`, named rather than a bool so the call site says
what it does: `WriteGuarded` refuses a subnet move while reservations exist, `WriteForced`
proceeds. Forcing waives the guard, never the validation.

A pool-only change is never guarded, because nothing moves — but the firmware takes all four
fields in one call, so read the current address and mask and pass them through:

```go
current, err := client.Network().Get(ctx)
if err != nil {
    return err
}
current.DHCPStart, current.DHCPStop = "192.168.2.50", "192.168.2.150"
err = client.Network().Set(ctx, current, services.WriteGuarded)
```

The pool is checked against the subnet first: the firmware accepts a pool outside it and then
silently hands out nothing.

When the subnet *does* move, expect the call to fail with a lost connection *on success* — the
router changes address mid-request. `goglnet` treats that as expected and tells you where to
reconnect. A pool-only change has no such problem.

### Errors

```go
switch {
case errors.Is(err, gogl.ErrDomainNotSet):      // set a domain before writing reservations
case errors.Is(err, gogl.ErrReservationsExist): // clear them, or pass services.WriteForced
case errors.Is(err, gogl.ErrNotFound):
case errors.Is(err, gogl.ErrConflict):        // MAC already reserved
case errors.Is(err, gogl.ErrInvalidName):     // name would be unsafe in dnsmasq config
case errors.Is(err, gogl.ErrInvalidMAC):
case errors.Is(err, gogl.ErrInvalidIP):
case errors.Is(err, gogl.ErrLoginRateLimited):   // locked out; the message carries the wait
case errors.Is(err, gogl.ErrOutsideSubnet):
case errors.Is(err, gogl.ErrUnauthorized):
case errors.Is(err, gogl.ErrWirelessSession):    // wireless write over a wireless session
case errors.Is(err, gogl.ErrUnwritableContent):  // host file holds a character set_host rejects
case errors.Is(err, gogl.ErrInvalidInput):       // SSID, passphrase, channel, bandwidth
}

// ErrSessionExpired, ErrNonceExpired, ErrUnsupportedAlgorithm and
// ErrUnsupportedHashMethod are exported for inspection but handled inside the
// transport: a caller normally never sees them.

var rpcErr *gogl.RPCError
if errors.As(err, &rpcErr) {
    log.Printf("%s.%s failed with %d: %s", rpcErr.Group, rpcErr.Method, rpcErr.Code, rpcErr.Message)
}
```

### Name validation happens in the library

```go
_, err := reservations.Create(ctx, &types.Reservation{
    Name: `my"host`, MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.13",
})
// err is ErrInvalidName; nothing was sent to the router
```

Two reasons, and neither is "the router will serve it as DNS" — it will not. GL.iNet writes
reservation names into its dnsmasq config file, and a known firmware defect lets a quote
character corrupt that file, taking down DHCP for the whole router until it is repaired by
hand. And the ISC DHCP host format this project exchanges is keyed by hostname, so a name
that is not a legal DNS label cannot round-trip through it.

Names are validated on the type, before any write, and rejected rather than escaped. No
consumer of the library can bypass it.

### Reaching an endpoint that isn't modelled

```go
var out map[string]any
err := client.Call(ctx, "some_group", "some_method", nil, &out)
```

`Call` is the generic escape hatch. It is also how the router's API was mapped in the first
place, since GL.iNet's official 4.x reference is no longer public.

## Project structure

```
gogl/
├── src/                    # Library (package gogl)
│   ├── client.go           # Client, Config, service accessors
│   ├── errors.go           # Sentinels, RPCError
│   ├── types/              # Reservation, Network, Client, SystemInfo, LeaseTime
│   ├── services/           # Network, Reservations, Clients, System
│   ├── auth/               # Challenge/response login, unix crypt
│   ├── transport/          # JSON-RPC envelope, session, keepalive, retry
│   ├── mock/               # In-process router for tests
│   └── ipmath/             # IPv4 arithmetic
├── utilities/              # goglps, goglnet, goglmac
├── examples/               # basic, list, reservations
└── docs/                   # DESIGN.md, plan.md
```

## Build targets

```bash
make build      # module + utilities into bin/
make test       # go test -v -race -cover ./...
make lint       # golangci-lint
make coverage   # HTML coverage report
make examples   # build the examples
make install    # utilities into ~/bin (override INSTALL_DIR)
make all        # lint, test, build
```

## How authentication works

Worth knowing if you are debugging a login, because it is the least obvious part of the
module.

Firmware 4.x uses challenge/response at `POST /rpc`. The password is never transmitted:

1. `challenge` returns a crypt algorithm, a salt, a nonce, and a `hash-method`.
2. The password is hashed with unix `crypt(3)` under that salt.
3. `H(username:cipher:nonce)` produces the login digest, where `H` is whatever
   `hash-method` named.
4. `login` returns a session id, passed as the first parameter of every later call.

**Step 3 is where every other client library gets it wrong.** `hash-method` selects the
login digest, separately from `alg` which selects the crypt on the password. Firmware 4.3.28
on the SFT1200 reports `sha256`; `python-glinet`, `glinet-client-go` and `ha-glinet` all
hardcode MD5 because they predate the field, and against 4.3.28 they fail with
`-32000 Access denied` — indistinguishable from a wrong password. gogl reads the field,
treats its absence as MD5 for older firmware, and refuses to guess if it sees a value it
does not implement.

The other subtlety: **the nonce lives about one second, and step 2 is deliberately slow.**
So `gogl` calls `challenge` *twice* — once to learn the salt, then again after the crypt
completes, to get a nonce that is still alive when the cheap digest runs over it. A
single-challenge implementation races against its own hashing cost and fails
intermittently, which is far worse than failing consistently.

The session then idles out after roughly 35 seconds. A background keepalive holds it open,
and any call rejected for a stale session triggers exactly one transparent re-login and
retry. Bounded at one attempt on purpose: an unbounded loop against a wrong password is a
login flood aimed at a very small SoC.

## Documentation

- [`VISION.md`](VISION.md) — requirements and per-tool specifications
- [`docs/DESIGN.md`](docs/DESIGN.md) — architecture, diagrams, interfaces, decision log
- [`docs/plan.md`](docs/plan.md) — phased implementation plan (continues in `plan-part2.md` through `plan-part4.md`)
- [`GL_INET_4X_API_DOCUMENTATION.md`](GL_INET_4X_API_DOCUMENTATION.md) — authentication,
  error codes, and every endpoint verified against real hardware
- [`docs/api/`](docs/api/README.md) — full reference for all 43 groups and 313 methods,
  generated from GL.iNet's own API description

## Related

- [`gofi`](https://github.com/emergingrobotics/gofi) — the UniFi UDM Pro counterpart. Same
  layout, same CLI conventions, same ISC DHCP file format. `gofips --get` output is
  `goglps --set` input.

## License

MIT

---

The Go gopher was designed by [Renée French](https://reneefrench.blogspot.com/) and is licensed under [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/).
