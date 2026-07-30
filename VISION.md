# gogl - Go GL.iNet Travel Router

A Go module for programmatic control of GL.iNet travel routers running firmware 4.x,
targeting the **GL-SFT1200 (Opal)**.

`gogl` is the travel-router counterpart to
[`gofi`](https://github.com/emergingrobotics/gofi). It exposes the same shape of API and
the same command-line ergonomics, so that a network described on a UniFi UDM Pro can be
reproduced on a pocket router. Dump the fixed-IP assignments from a UniFi site with
`gofips`, hand the file to `goglps`, and the travel router hands out the same addresses to
the same MAC addresses.

Addresses, not names: a reservation on this firmware does not create a DNS record. See
[The Reservation Model](#the-reservation-model).

## Target Device

| Property | Value |
|----------|-------|
| Model | GL-SFT1200 ("Opal") |
| Firmware verified | **4.3.28** (model string `sft1200`) |
| SoC | SiFlower SF19A28 |
| OpenWrt base | 18.06 |
| Ethernet | 1 x GbE WAN, 2 x GbE LAN |
| Wireless | AC1200 dual-band (2.4 GHz 300 Mbps, 5 GHz 867 Mbps) |
| Default LAN | `192.168.8.1/24` |
| Admin user | `root` (password set during first-run setup; there is no default) |

Firmware 4.x is a hard requirement. GL.iNet replaced the older REST API with JSON-RPC at
firmware 4.0, and `gogl` speaks only JSON-RPC. Firmware 3.x devices are out of scope.

Other GL.iNet 4.x models should work, since the API is shared across the firmware line, but
the SFT1200 is the only device this module is tested against — and it is a **reduced build**:
six endpoints GL.iNet documents are absent on it. See [`docs/api/`](docs/api/README.md),
where every method is marked verified, absent, or untested.

## Project Resources

- **API Reference**: [`GL_INET_4X_API_DOCUMENTATION.md`](GL_INET_4X_API_DOCUMENTATION.md) -
  authentication, error codes, and every endpoint verified against a live SFT1200.
- **Full API**: [`docs/api/`](docs/api/README.md) - all 43 groups and 313 methods, generated
  from GL.iNet's own API description.
- **Architecture**: `docs/DESIGN.md` - Design with mermaid diagrams, service interfaces,
  type definitions.
- **Implementation Plan**: `docs/plan.md` - Phased plan with progress tracking.

## Critical Rules

1. **Every function MUST have a test** - No exceptions. Run `make test` to verify.
2. **Every endpoint MUST be supported in the mock server** - Tests use the mock, not real hardware.
3. **No phase advancement without 100% test coverage** - Complete and test each phase before moving on.
4. **Phases are sequential** - Follow `docs/plan.md` in order.
5. **Structured transports only** - Every operation goes through GL.iNet's JSON-RPC API, so
   the whole surface stays reachable by an in-process mock and no hardware appears in the test
   suite. **Never build a command string from user-supplied data.** Anything the API cannot
   express is a documented gap, never a shell workaround.

   This rule was inherited from `gofi`, where it read "no SSH, no UCI, no shell", and the
   justification originally written for it here -- that SSH and shell are unmockable -- was
   false. `uci show` is text in and text out, easier to mock than the HTTP server this project
   already mocks. The rule is kept for the two reasons that actually hold.

   A structured API is a capability boundary: it can set a DHCP pool and cannot delete `/etc`,
   which bounds the blast radius of a bug in a tool built for unattended provisioning. And it
   keeps user data out of shell command strings, which matters here because a quote in a
   hostname can corrupt dnsmasq's configuration -- the reason `gogl` rejects such names rather
   than escaping them.

## Scope

**GL.iNet firmware 4.x only. This is not a portable OpenWrt tool and will not become one.**

The devices run OpenWrt underneath, and roughly 59% of GL.iNet's 313 API methods wrap standard
OpenWrt or Linux subsystems, so a UCI-over-`/ubus` backend was considered and specified. It was
dropped deliberately: supporting arbitrary OpenWrt hardware means two implementations of every
service, two credential paths, two sets of field-name translations, and a name that no longer
describes the tool -- in exchange for devices this project does not target. Focus is worth more.

The practical consequences, all of them simplifications:

- `src/types` may hold GL.iNet's wire shapes without apology. `IntBool` exists because the
  firmware sends `enable: 1`; `HTModes` is its capability object. These were once described here
  as vendor artifacts leaking into a neutral layer. There is no neutral layer, so they are simply
  the domain.
- The DNS domain living in a marker line inside the host file is the answer, not a workaround
  pending a better backend. GL.iNet's API exposes no dnsmasq domain setting and `POST /ubus`
  returns 404 on this hardware, so there is nothing better to wait for.
- `wan` can model GL.iNet's own concepts -- `repeater`, `tethering`, `modem`, `netmode` -- with
  no question of how vendor-specific commands would fit a portable tool.
- The name and the logo stay.

`gogl` v1 manages, on a single flat LAN with no VLANs:

- **Network configuration** - LAN address, netmask, DHCP pool. Read and **written**
  (`lan.set_config`). Writing is refused while reservations exist; see
  [Ordering Rules](#ordering-rules).
- **Static IP reservations** - MAC-to-IP bindings, created, updated, deleted individually or
  cleared wholesale.
- **DNS names** - through the router's host file (`dns.set_host`), which dnsmasq answers from.
  A reservation does **not** create a DNS record; see
  [The Reservation Model](#the-reservation-model).
- **The DNS domain** - carried by gogl inside its block in the host file, because the firmware
  exposes no dnsmasq domain setting. Required before any reservation write.
- **DHCP leases** - read-only, and the way to discover what is worth reserving.
- **Whole-network profiles** - every writable section above, captured to JSON and applied back,
  onto the same router or another one. See [Profiles](#profiles).
- **Connected clients** - with independent IEEE OUI manufacturer lookup.
- **Wireless identity** - SSID, passphrase, encryption, hidden and enabled state, per
  interface. Read and **written**.
- **Radio tuning** - channel, bandwidth, hardware mode, transmit power, per radio. Read and
  **written**. Both under the safety rules in [Wireless Writes](#wireless-writes).

Explicitly **out of scope for v1**:

- **Nothing in the wireless stack.** Radio tuning was excluded on the grounds that it is
  tuning rather than provisioning; that was wrong for a travel router, where the channel a
  site's existing WiFi occupies is exactly the thing you need to move off.
- **VLANs.** The Opal exposes no VLAN configuration through its API.
- **Guest network writes.** Reported by `goglnet` for address planning; `lan.set_config` can
  address it, but v1 only ever writes the main LAN.
- **VPN, firewall, port forwarding, traffic rules.**
- **dnsmasq's own `domain` / `local` / `expandhosts` settings.** No endpoint exposes them and
  `/ubus` is unavailable. gogl writes fully-qualified names into the host file instead, which
  works for any suffix.

An earlier version of this document made "reservations are the only thing gogl writes" a
design rule, on the grounds that a bulk import should never be able to change the address the
router is managed at. That concern is real, and it is now handled by the ordering rules below
rather than by refusing to write at all.

## Ordering Rules

Two preconditions, both enforced in the library so no consumer can bypass them, and both
covering states that are tedious to recover from.

### A reservation write requires a configured domain

`Create` and `Update` return `ErrDomainNotSet` until a domain exists.

A reservation creates no DNS record on its own. Writing one before the domain is set yields an
address with no name and nothing to indicate that was unintended -- you find out later, when
something cannot resolve. Making the domain a precondition forces the pairing to be
deliberate. Reads and deletes are ungated: only writes that create addressing are.

### Moving the subnet requires no reservations, unless forced

`Network().Set` returns `ErrReservationsExist` while any reservation is present, unless called
with `WriteForced`.

The original reason was wrong and is worth recording. It was written expecting the firmware to
leave static binds alone, stranding all of them outside the new subnet. Observed on hardware
2026-07-29: the firmware silently renumbers every bind into the new subnet, preserving host
parts. Twenty-seven reservations moved from `192.168.2.x` to `192.168.8.x` with no prompt.

The guard is kept on narrower grounds. An unannounced rewrite of every reservation is a large
side effect of an address flag, and the behavior is only characterized for a same-size subnet:
a narrower netmask, where a host part no longer fits, is untested, as is a move that lands
addresses inside the new DHCP pool -- which happened to twenty of those twenty-seven. `--force`
exists because the rewrite is usually what the operator wants.

**A pool-only change is never guarded.** The subnet does not move, so the session survives and
no reservation can be renumbered by it. Guarding it would be guarding against something that
cannot happen. Only a change to the address or the netmask counts as moving the subnet.

The correct order is therefore:

1. `goglps --domain <domain>`
2. `goglps --clear` if reservations already exist, or pass `--force` in step 3
3. `goglnet --set-ip ... --set-mask ... --set-start ... --set-end ...`
4. `goglps --set <file>`

Step 3 moves the router, so the session drops. That is inherent, not a defect. Both writing
steps take `--dry-run`, and a dry run runs every check the real write runs, including the
refusals: a preview that approves what the write would reject is worse than no preview.

`--clear` removes the managed DNS names along with the reservations, because they are one
intent stored in two tables. Clearing only the bindings would leave names resolving to
addresses the router no longer reserves, and since clearing is what unblocks step 3, those
names would end up answering for a subnet that no longer exists — exactly the stranded state
the guard exists to prevent. The domain from step 1 survives a clear.

## Wireless Writes

Reproducing a network means reproducing what devices connect *to*, not only what addresses they
receive. A signage player or robot configured for a given SSID does not care that the DHCP
reservations are right if it cannot associate. So `gogl` writes wireless identity: SSID,
passphrase, hidden flag, and enabled flag, per interface.

This was out of scope in v1 for a reason that has not gone away, so the reason becomes a
guard rather than a prohibition.

### The lock-yourself-out problem

Changing an SSID or passphrase drops every client on that radio, including the session issuing
the request. Unlike a LAN renumber, there is no new address to reconnect at: the network the
management session was using has ceased to exist under that name. Recovery means ethernet or
the reset pin.

These rules apply to **every** wireless write, tuning included. Retuning a radio drops its
clients for at least a re-association, and a DFS channel change for the minutes the radio
spends re-scanning, so the distinction between "identity" and "tuning" does not matter to
whoever loses their session.

Two rules:

1. **The session must not be arriving over WiFi.** `gogl` determines the local address it
   reaches the router from, finds that address in `client.get_list`, and reads the `iface`
   field the firmware reports for it -- `cable`, `2.4G` or `5G`. Anything but `cable` is
   refused with `ErrWirelessSession`.

   If the address is not in the client list at all, the session is arriving from off-LAN
   through a router, so the radio cannot be the path. That is allowed, and noted.

   A stricter reading would refuse only when the session arrives over the *same* radio being
   modified, since changing the 2.4G SSID from a 5G association is safe. That refinement is
   deliberately not in this version: the simple rule is the one that cannot be got wrong, and
   the cost of it is having to plug in a cable.

2. **A human must confirm.** Interactive invocations prompt `y/N` and show the before and
   after. `--yes` skips the prompt for scripted use, and a non-terminal stdout does not
   silently proceed -- unlike `--del`, where proceeding is the safe assumption. Here the
   failure mode is losing the device, so the default is to stop.

### Partial updates, and why the flags look the way they do

`wifi.set_config` leaves any field it is not given alone, and `gogl` sends only the fields
asked for. Sending the unchanged values back would work, but would turn every write into a
chance to clobber a concurrent edit from the admin panel.

That requires telling "set this to false" apart from "do not mention it", which Go's `flag`
package cannot do: both leave a `*bool` at `false`. So the wireless flags take an explicit
value -- `--set-hidden=false`, not a bare `--set-hidden`. A bare boolean flag meaning true is
exactly how `--set-enabled` ends up disabling something.

The same applies to `--set-channel=0`, where 0 is a real value meaning "choose
automatically" and so cannot double as a not-set sentinel.

### The firmware scopes writes two ways

Interface-scoped fields -- SSID, passphrase, encryption, hidden, enabled -- require
`iface_name`. Radio-scoped fields -- channel, bandwidth, hardware mode, transmit power --
require `device`. A write carrying both goes out as two calls, because that is how the
firmware separates them, and because a failure then names which half did not apply.

`goglnet` reports what each radio advertises as supported: its channels with the DFS ones
marked, its bandwidths, its hardware modes, its encryptions. Without that the tuning flags are
unusable, since the valid values differ per radio and per regulatory domain. `gogl` validates
against those lists before writing, so a bad channel is refused with the available ones named
rather than answered with a bare error code.

### The capability payload is not what the docs say

`htmodes` is the field to be careful with. GL.iNet's description calls it a "List of supported
bandwidths", an array of strings; the device sends an object keyed by hardware mode whose values
are the maximum channel width in MHz, plus an `auto` key holding a bool:

```json
"htmodes": {"11b/g/n": 40, "11g/n": 40, "11n": 40, "auto": true}
```

Typing it from the description made every read of `wifi.get_config` fail outright. `hwmodes` was
wrong too -- slash-joined combinations such as `11a/n/ac`, not the bare `11b`/`11g`/`11n` the
description implies -- and the encryption list contains `sae` and `sae-mixed` for WPA3 while
containing no bare `psk`.

Two consequences for the design. The settable `htmode` is a width string (`auto`, `20`, `40`,
`80`), and the narrower widths are *inferred* to be valid from the reported maximum, since the
firmware states only the maximum. And `src/mock`'s wireless fixture is a verbatim capture rather
than a composed example, with a test asserting it decodes through the same types the library
uses -- the check whose absence let a fixture and a type agree with each other while both
disagreed with the hardware.

### The passphrase is readable, and that is accepted

`wifi.get_config` and `system.get_status` both return passphrases in cleartext over plain HTTP
on port 80. Reading them requires being on the router's network and holding the admin password,
which is the same bar as opening the admin panel, so this is not treated as a defect to work
around. `gogl` does not log the raw payloads, and `goglnet` masks the key in its own output
unless asked for it, because a passphrase on a terminal is a passphrase in a scrollback buffer.

## Profiles

`goglcfg` captures the writable sections as a JSON profile and applies one back. It is the
fourth utility and has no `gofi` counterpart: the work spans `goglps`, `goglnet` and `goglmac`
and belongs in none of them. The three-tool mirror exists so that knowing one set means knowing
the other, not as a cap on what the project may contain.

### A profile is a network, not a router

It carries the LAN address and pool, reservations, DNS names and domain, wireless identity and
radio tuning. It omits the router's own MAC, serial number, uptime and lease state -- the fields
that make a full configuration dump worthless on a second unit, which is the case a profile
exists to serve. Client MAC addresses are included: a reservation is a MAC-to-IP binding and a
profile without them reproduces nothing.

Every section comes from an endpoint verified against hardware. The API exposes 110 getters and
23 are verified; a profile built on the remainder would be guesswork, and guessing from GL.iNet's
descriptions has been wrong three times in this project. Lease time, upstream DNS servers,
firewall, VPN and VLANs are absent for that reason, not because they were forgotten.

### Passphrases are omitted unless asked for

A profile is a file people commit. `--get` writes no WiFi keys; `--get --with-keys` includes
them in cleartext.

An omitted key is not an empty key. On apply, a missing key is not written at all, leaving
whatever the target already has -- so the private default is also the safe one. This depends on
`wifi.set_config` leaving unmentioned fields alone, verified on hardware 2026-07-29.

### Apply order

Fixed, and each step is where it is because doing it later fails:

1. **Domain.** Reservation writes are refused without one.
2. **Network.** Reservations must be inside the subnet before they are written.
3. **Reservations**, then **DNS names**.
4. **Wireless**, opt-in via `--wireless`. It needs a wired session and is the section least
   likely to be wanted, so a refusal there must not undo the addressing.

### A subnet move ends the run

When the profile's subnet differs from the router's, the router changes address during step 2 and
nothing after it is reachable from that session. `goglcfg` stops, says the router has moved, and
prints the command to resume at the new address.

Reporting success for a run that applied a third of a profile would be a lie, and continuing
would mean writing reservations over a connection that no longer exists. A pool-only difference
is not a move and runs straight through.

### Cross-model applies warn rather than fail

Addresses and names are portable. Wireless is not: interface names, radio names, channel lists
and hardware modes are per-device and per-regulatory-domain. A profile from another model warns,
and any interface or radio the target lacks is reported and skipped rather than failing the load.

## Concurrency Requirements

- **Session ID**: Use `atomic.Value` for thread-safe storage and updates.
- **Re-authentication**: Use `sync.Mutex` so that a burst of calls hitting an expired
  session triggers exactly one re-login, not one per call.
- **Keepalive**: A single background goroutine issues `alive`; it must be cancellable via
  the client's context and must not outlive `Close()`.
- **Rate limiting**: Implement semaphore-based request limiting. The SFT1200 is a small
  SoC and will drop requests under load.
- **Connection pooling**: Configure `http.Transport` appropriately.

## Architecture Overview

The module path is `github.com/emergingrobotics/gogl`; the library source lives under
`src/`, so the root package is imported as `github.com/emergingrobotics/gogl/src`
(package name is `gogl`). `examples/` and `utilities/` remain at the repo top level.

This mirrors `gofi`'s layout deliberately, so the two modules read side by side.

```
gogl/
├── src/               # Library source (package gogl + sub-packages)
│   ├── client.go      # Main client, session lifecycle, request handling
│   ├── types/         # All domain types (Network, Reservation, Client, ...)
│   ├── errors.go      # Sentinel errors, RPCError type
│   ├── services/      # Service implementations
│   │   ├── network.go     # LAN address, DHCP server settings
│   │   ├── reservation.go # Static MAC-to-IP bindings
│   │   ├── client.go      # Connected clients
│   │   └── system.go      # Model, firmware version, uptime
│   ├── mock/          # Mock JSON-RPC server for testing
│   │   ├── server.go
│   │   ├── handlers.go
│   │   ├── fixtures/
│   │   └── scenarios/
│   ├── auth/          # Challenge/response login, unix crypt hashing
│   ├── transport/     # JSON-RPC envelope, sid injection, retry on expiry
│   └── internal/      # Internal helpers
├── examples/          # Runnable example programs (consumers)
└── utilities/         # CLI tools: goglps, goglmac, goglnet
```

### Differences from gofi's service layer

- **No site concept.** GL.iNet routers have no equivalent of a UniFi site. Every `site`
  parameter present in `gofi`'s method signatures is absent here.
- **No DNS service at all.** Not because a reservation doubles as one — it does not, see
  [The Reservation Model](#the-reservation-model) — but because the firmware exposes no way to
  set a per-host DNS record. `gofi` needs `Users()` and `DNS()` kept in sync; `gogl` has
  neither to offer.
- **No device service.** The router is the only device. `System()` reports it.
- **No WebSocket.** Firmware 4.x has no event stream equivalent to UniFi's.

## Key Technical Details

- **Endpoint**: A single JSON-RPC 2.0 endpoint at `POST /rpc`. There are no per-resource
  paths and no API versions in the URL.
- **Auth**: Challenge/response. The password is never transmitted, in cleartext or
  otherwise; only a digest derived from it.
- **Session**: A `sid` returned by `login`, passed as the first element of every
  subsequent call's `params` array.
- **Scheme**: GL.iNet routers serve the admin interface over **HTTP on port 80** by
  default, unlike the UDM Pro's HTTPS on 443. So `gogl` defaults to HTTP port 80 and
  requires `--https` to use TLS.
- **TLS**: Where HTTPS is available it uses a self-signed certificate, so under `--https`
  certificate verification is **off by default** and enabled with `--secure`/`-k`. This
  matches `gofi`'s inverted `-k` semantics. Without `--https`, `-k` has no effect and
  passing it is an error rather than a silent no-op.

### Authentication Flow

```
POST /rpc
{"jsonrpc":"2.0","id":1,"method":"challenge","params":{"username":"root"}}

<- {"jsonrpc":"2.0","id":1,"result":{"alg":6,"salt":"<salt>","nonce":"<nonce>"}}
```

`alg` selects the unix crypt algorithm:

| `alg` | Algorithm |
|-------|-----------|
| `1` | MD5 |
| `5` | SHA-256, 5000 rounds |
| `6` | SHA-512, 5000 rounds |

Firmware has been observed returning `alg` as both a JSON number and a JSON string, so it
must be decoded tolerantly. An unrecognized `alg` is a hard error, never a fallback to a
weaker algorithm.

The login digest is then computed in two steps:

1. `cipher = crypt(password, "$" + alg + "$" + salt + "$")` - standard unix crypt,
   equivalent to `openssl passwd -6 -salt <salt> <password>`. The full crypt output is
   used, including the `$6$<salt>$` prefix.
2. `hash = md5_hex(username + ":" + cipher + ":" + nonce)`

```
POST /rpc
{"jsonrpc":"2.0","id":2,"method":"login","params":{"username":"root","hash":"<hash>"}}

<- {"jsonrpc":"2.0","id":2,"result":{"sid":"<sid>","username":"root"}}
```

**The nonce is valid for roughly one second.** The transport MUST issue `challenge`
immediately before `login` and MUST NOT cache a challenge across calls. Because unix
crypt with 5000 rounds of SHA-512 is deliberately slow, the implementation should compute
the expensive `cipher` step first, then re-issue `challenge` to obtain a fresh nonce, then
compute the cheap MD5 and log in. Otherwise the crypt cost alone can expire the nonce.
The `cipher` value depends only on the password and salt, so it can be cached for the
lifetime of the client; the nonce and the resulting `hash` cannot.

### Authenticated Calls

Every operation is the `call` method, with the `sid`, a functional group, a method name,
and an argument object:

```
POST /rpc
{"jsonrpc":"2.0","id":3,"method":"call","params":["<sid>","<group>","<method>",{...}]}
```

### Session Lifetime

The `sid` expires after approximately 35 seconds of inactivity, and every successful call
resets the timer. Two mechanisms keep it alive:

1. A background goroutine issues `alive` at an interval safely under the expiry window
   (default 20 seconds).
2. Any call that fails with an access-denied error triggers exactly one transparent
   re-login and one retry. A second failure is returned to the caller.

Because sessions expire on a timescale shorter than a human typing a command, the CLI
utilities log in at the start of each invocation and do not attempt to persist sessions
between runs.

### API Discovery

GL.iNet's official 4.x API reference at `dev.gl-inet.com/router-4.x-api` is no longer
publicly reachable. The authentication flow above is verified, but **exact functional
group and method names must be captured from a live SFT1200 before the typed services are
frozen.**

Phase 0 of `docs/plan.md` is therefore a discovery phase:

1. The transport exposes a generic, untyped escape hatch:
   ```go
   func (c *Client) Call(ctx context.Context, group, method string, args, out any) error
   ```
2. Every request and response observed against the live device is recorded verbatim into
   `src/mock/fixtures/`.
3. `GL_INET_4X_API_DOCUMENTATION.md` is written from those recordings, not from
   third-party documentation.
4. Only then are the typed service interfaces written against the recorded fixtures.

The groups expected to matter for v1 scope, to be confirmed during discovery, cover LAN
and DHCP settings, static lease management, connected clients, and system information.
Do not hard-code a group or method name into a service until it appears in a recorded
fixture. The generic `Call` remains public after Phase 0, both as an escape hatch and
because it is how future discovery happens.

## The Reservation Model

This section originally asserted that one reservation provides both the DHCP binding **and**
the DNS record, atomically. **That was tested against a GL-SFT1200 on firmware 4.3.28 and is
false.** The corrected model follows; the original claim is recorded because a good deal of
the design was built on it.

### What a reservation actually is

The Opal runs dnsmasq, and GL.iNet's admin panel presents **LAN → Address Reservation** with
Name, MAC and IP fields. The natural reading is that the Name becomes a dnsmasq hostname.
It does not.

```
Reservation{Name: "myserver", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.10"}
```

...yields a device that always receives `192.168.8.10`. It does **not** yield a name that
resolves.

Two experiments established this:

1. A reservation for a MAC no device was using. The name resolved neither bare nor with the
   `.lan` suffix.
2. A reservation for a **present, actively leased** device under a *different* name. The new
   name still did not resolve, while that device's original lease hostname continued to
   resolve throughout.

### Where DNS names actually come from

The router answers from **DHCP leases**. The hostname in a lease is the one the client
announced about itself. So:

| | Source | Controlled by |
|---|--------|---------------|
| Address a device receives | static bind (reservation) | **gogl** |
| Name that resolves | DHCP lease hostname | **the client** |

### Consequences for this project

- **gogl reproduces addresses, not names.** That is the honest scope.
- **A reservation's Name is a label.** It identifies the entry in the admin panel and keys it
  in an exported host file, which is genuinely useful — but nothing resolves it.
- **`--keep-dns` still has no analogue**, for a duller reason than originally argued: there is
  no DNS record to keep.
- **Drift is still impossible**, but trivially so. `gofips` must keep a UniFi user record and
  a separate DNS record consistent; gogl writes one object that has no DNS half.
- **Names often work anyway, and that is not gogl's doing.** Most devices announce a sensible
  hostname over DHCP, so `nas.lan` frequently resolves. gogl neither causes that nor can
  arrange it for a device that stays silent.

### Name Validation

The name still gets strict DNS-label validation, for two reasons that survive the correction:

GL.iNet writes the name into its dnsmasq configuration file, and a known firmware defect lets
a character such as `"` corrupt that file, breaking DHCP for the entire router until it is
repaired by hand. And the ISC DHCP host-declaration format this project exchanges is keyed by
hostname, so a name that is not a legal DNS label cannot round-trip through it.

Validation rejects rather than escapes:

- Permitted characters: `[a-zA-Z0-9-]`, plus `.` as a label separator.
- Must not begin or end with `-` or `.`.
- Each label at most 63 characters; total at most 253.
- Must not be empty.

This lives in the library, not the CLI, so no consumer can bypass it.

**This is deliberately stricter than `gofips`**, which accepts `_`. A host file containing
`my_server` is **rejected** with the offending character named, rather than silently rewritten
to `my-server`. It is the one case where a `gofips` file may not import unchanged, and
renaming is a decision about what your hosts are called — not one to make silently.

## Type Patterns

GL.iNet's JSON-RPC responses are well-typed, unlike UniFi's. `gofi`'s `FlexInt` and
`FlexBool` workarounds are deliberately **not** ported. Add them only if a recorded
fixture proves they are needed.

One conversion is genuinely required. UniFi expresses `dhcpd_leasetime` as an integer
number of seconds (`86400`); dnsmasq expresses lease time as a duration string (`12h`,
`1d`) or bare seconds. A `LeaseTime` type handles both directions:

```go
// LeaseTime is a DHCP lease duration. It unmarshals from a dnsmasq duration
// string ("12h", "1d", "infinite") or a bare number of seconds, and marshals
// to the dnsmasq duration form.
type LeaseTime time.Duration
```

`infinite` maps to a sentinel value rather than an overflowing duration.

## Mock Server

The mock server must:

- Implement the full JSON-RPC envelope, including `challenge`, `login`, `alive`, and `call`.
- Reproduce the real challenge/response flow with a configurable `alg`, so that the crypt
  implementation is exercised for MD5, SHA-256, and SHA-512.
- Enforce nonce expiry, so that a client which caches a challenge fails against the mock
  exactly as it would against hardware.
- Enforce session expiry, so that keepalive and transparent re-login are exercised.
- Support all endpoints with realistic responses, loaded from fixtures recorded from a
  live device.
- Support error scenarios: bad password, expired nonce, expired session, malformed
  request, and JSON-RPC error objects.

## Replicating a gofi Network

This is the reason `gogl` exists. The workflow takes a network as configured on a UniFi
UDM Pro and reproduces its addressing and naming on a travel router.

### The Interchange Format Is the ISC DHCP File

No new format is introduced. `gofips --get` already emits exactly the three fields an Opal
reservation needs, keyed by hostname:

```
host myserver {
    hardware ethernet aa:bb:cc:dd:ee:01;
    fixed-address 192.168.8.10;
}
```

`goglps` parses and emits the identical format with identical rules, so files move between
the UniFi controller and the travel router with no conversion step. The format is
diffable, version-controllable, and hand-editable, which is the whole argument for it.

One caveat: `goglps` validates hostnames more strictly than `gofips` does, so a file
containing a name that is legal on UniFi but not a legal DNS label is rejected rather than
converted.
See [Name Validation](#name-validation).

### End-to-End Workflow

Two commands. This example assumes a home network of `192.168.4.0/24` behind a UDM Pro at
`192.168.4.1`, and a travel router whose LAN address has already been set to
`192.168.4.1/24` by hand. That last part is a prerequisite; see
[The Router's Own Configuration](#the-routers-own-configuration).

```bash
# On the UniFi side, at home:
gofips -H 192.168.4.1 -k --get > home.hosts

# Later, against the travel router:
goglps -H 192.168.4.1 --set home.hosts
```

The two `-H 192.168.4.1` values refer to different devices on different occasions: the UDM
Pro while you are at home taking the dump, and the travel router later. They coincide
because the travel router is standing in for the UDM's LAN address, which is the entire
point. Do not run both halves with both devices on the same segment.

Use `goglnet` beforehand to confirm the router's subnet and pool, and `goglps --set
--dry-run` to see exactly what would change before it changes.

The result is a travel router that hands out the same addresses to the same MAC addresses
and answers the same names as the network it was copied from. Plug it into hotel ethernet,
and devices that expect `myserver` at `192.168.4.10` find it there.

### What Crosses Over

| From the UniFi dump | On the travel router |
|---------------------|----------------------|
| Each `host {}` block, address | One static bind: a MAC-to-IP binding |
| Each `host {}` block, hostname | One host-file entry, answering both bare and qualified |

Each host declaration therefore becomes **two** writes, to two independent tables that the
firmware does not join for you. `goglps` performs both and reports them as one row, which is
the honest presentation: they are one intent, and either can be present without the other.

Nothing else crosses over, because `gofips --get` emits host declarations and nothing else.
The network's *shape* — subnet and pool boundaries — is not in the file, so you either set it
with `goglnet --set-*` or in the admin panel. Lease time and upstream DNS servers stay
whatever the router has; see the warning below for why you want that for DNS in particular.

### The Router's Own Configuration

`goglnet` reports the subnet, pool, and lease time, and writes the first two with
`--set-ip`, `--set-mask`, `--set-start`, and `--set-end`. Lease time and upstream DNS
servers are read-only: no endpoint on this firmware sets them.

Two things to get right, because both cause failures that are hard to read from the symptom:

**Match the subnet to the network you are standing in for.** A reservation at
`192.168.4.10` is meaningless on a router whose LAN is `192.168.8.0/24`. Set the LAN address
first, either with `goglnet --set-ip 192.168.4.1 --set-mask 255.255.255.0 --set-start ...
--set-end ...` or in the admin panel under LAN → Router IP Address. The management session
drops when the router moves, and you reconnect at the new address; `goglnet` treats a
connection loss after the request as success, because that is what a successful renumber
looks like from this end.

`goglnet` refuses the change while any reservation exists, and `goglps` refuses any
reservation outside the current subnet rather than writing something inert. Between them there
is no order of operations that leaves a reservation stranded, and neither failure is silent.

**Do not copy your home network's DNS servers.** This is the trap. If your home network
hands out a Pi-hole at `192.168.4.5`, and you configure the travel router to advertise that
same resolver, then every client at the customer site is told to use a DNS server that does
not exist there. It presents as "the internet is broken," not as "DNS is misconfigured."
Leave the router advertising **itself**: dnsmasq answers your reservation names locally and
forwards everything else to whatever resolver the WAN handed it. That works in a hotel, at
a customer site, and on a bench with no uplink at all. Public resolvers like `1.1.1.1` are
safe to set if you want them, since they work anywhere.

### Future Work

- **`goglrenumber`** - rewrite the network part of every address in an ISC DHCP file to a
  target subnet, preserving host parts, emitting a new file rather than editing in place. A
  pure text transformation with no device access, which makes it cheap to build and easy to
  trust. This is the alternative to re-IPing the router: instead of moving the router to
  the file's subnet, move the file to the router's.
Renumbering the router is no longer future work: `goglnet --set-*` does it. It stayed out of
`goglps` deliberately, so that losing contact with the device is never a side effect of
importing a host file.

`goglrenumber` is not in v1.

## Common CLI Conventions

All utilities share these.

### Connection Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--host` | `-H` | `$GL_ROUTER_IP` | Router host address |
| `--port` | `-p` | `80` | Router port |
| `--https` | n/a | `false` | Use HTTPS instead of HTTP |
| `--secure` | `-k` | `false` | Under `--https`, enforce TLS certificate verification (default: accept self-signed) |

The port default is 80, not 443, because that is what GL.iNet firmware serves. Passing
`--https` without changing `--port` is almost certainly wrong; the tools warn when
`--https` is combined with port 80. Passing `-k` without `--https` is an error.

There is no `--site` flag. GL.iNet routers have no sites.

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GL_PASSWORD` | Yes | Router admin password |
| `GL_USERNAME` | No | Router admin username (default `root`) |
| `GL_ROUTER_IP` | No | Router host address (fallback if `-H` not given) |

### Output Conventions

- Data goes to stdout. Status, progress, and warnings go to stderr.
- Addresses sort numerically by IP (uint32 conversion), never lexically.
- MAC addresses are lowercase and colon-separated on output; case-insensitive on input.

## goglps Tool

A command-line tool for managing static IP reservations and DNS names on a GL.iNet travel
router, using ISC DHCP host declaration format. Lives in `utilities/goglps/`. Built on the gogl
module. The counterpart to `gofips`.

Each host declaration becomes **two** writes, because the firmware stores the two facts
separately: a static bind for the address, and a host-file entry for the name. The hostname on a
bind is a label and resolves nothing; the DNS record comes from the host file. `goglps` keeps
both in step, so a caller writes one declaration and gets both.

### Purpose

A travel router needs the same MAC-to-IP-to-name bindings as the network it stands in for.
The ISC DHCP `dhcpd.conf` host declaration format is the industry-standard way to
represent these bindings, and is what `gofips` emits. `goglps` reads and writes that same
format so that administrators can:

- Export the router's reservation table to a file that can be version-controlled, diffed,
  and edited.
- Import a file of host declarations to provision the router in bulk, including files
  produced by `gofips` from a UniFi controller.
- Add or delete individual hosts from the command line using the same format fragment.
- Set the DNS domain the names are qualified under, and clear both tables wholesale.

### ISC DHCP Host Declaration Format

The canonical format for each host entry is:

```
host myhostname {
    hardware ethernet aa:bb:cc:dd:ee:ff;
    fixed-address 192.168.8.10;
}
```

Rules for parsing and emitting this format, identical to `gofips` so that files are
portable between the two tools:

- `host <hostname> {` opens a declaration. The hostname is a single token (no spaces,
  DNS-safe characters: `[a-zA-Z0-9.-]`).
- `hardware ethernet <mac>;` specifies the MAC address. Colon-separated hex,
  case-insensitive on input, lowercase on output.
- `fixed-address <ip>;` specifies the IPv4 address. Only IPv4 is supported.
- `}` closes the declaration.
- Blank lines between declarations are ignored.
- Lines starting with `#` are comments and are ignored on input. On output, a header
  comment is emitted.
- Whitespace is flexible on input: leading/trailing spaces, tabs, and varying indentation
  are all accepted. On output, use 4-space indentation as shown above.
- Semicolons after `hardware ethernet` and `fixed-address` values are required.
- Declarations may appear in any order in the file. On output, sort by IP address
  numerically.
- Other dhcpd.conf directives (subnet, option, etc.) are silently ignored during parsing;
  only `host {}` blocks are extracted.

### CLI Interface

```
goglps [connection flags] --get
goglps [connection flags] --set [filename]
goglps [connection flags] --add <host-declaration>
goglps [connection flags] --del --name <hostname>
goglps [connection flags] --del --mac <mac>
goglps [connection flags] --del --ip <ip>
```

### Mode Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--get` | `-g` | Export all reservations to stdout in ISC DHCP format |
| `--set` | `-s` | Import host declarations from a file or stdin |
| `--add` | `-a` | Add a single host from a declaration fragment (string argument or stdin) |
| `--del` | `-d` | Delete a host identified by `--name`, `--mac`, or `--ip` |
| `--domain` | n/a | Set the DNS domain. Required before any reservation write |
| `--clear` | n/a | Delete every reservation and every managed DNS name |

Exactly one mode flag must be specified. Specifying none, or more than one, prints usage
and exits with error.

### Identifier Flags (used with `--del`)

| Flag | Short | Description |
|------|-------|-------------|
| `--name` | `-n` | Hostname to delete |
| `--mac` | `-m` | MAC address to delete |
| `--ip` | `-i` | IP address to delete |

Exactly one identifier must be given with `--del`. If a MAC or IP matches multiple
entries (should not happen, but defensive), report all matches and require `--force` to
proceed.

### Other Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-f` | Skip conflict checks; proceed even if multiple matches found on delete |
| `--prune` | `-P` | On `--set`, delete reservations present on the router but absent from the file, and their DNS names |
| `--dry-run` | n/a | Show what would be done without making changes. Runs every check the real write runs, including the refusals |

`gofips`'s `--keep-dns` flag has no analogue, for a duller reason than first assumed: a
reservation creates no DNS record, so there is nothing to keep. See
[The Reservation Model](#the-reservation-model).

### Behavior: `--get` Mode

1. Connect to the router.
2. List all reservations via `client.Reservations().List()`.
3. For each reservation, the hostname is the reservation's own name. There is no fallback
   chain and no cross-referencing, because there is no second record that could disagree.
   A reservation with an empty name is emitted with the MAC as its hostname, colons replaced
   by hyphens (e.g. `aa-bb-cc-dd-ee-ff`), and a comment saying so. Be aware that this does not
   round-trip inertly: feeding that file back through `--set` gives the reservation that
   MAC-derived label for real. The emitted comment says so.
4. Sort entries by IP address numerically.
5. Output to stdout in ISC DHCP format with a header comment:

```
# goglps reservations
# exported from GL-SFT1200 at 192.168.8.1
# lan: 192.168.8.1/24  pool: 192.168.8.100-192.168.8.249  domain: lan
# date: 2026-07-27

host myserver {
    hardware ethernet aa:bb:cc:dd:ee:01;
    fixed-address 192.168.8.10;
}

host printer {
    hardware ethernet aa:bb:cc:dd:ee:02;
    fixed-address 192.168.8.11;
}
```

The header records the LAN subnet and pool, so that a file's intended network is evident
from the file itself.

6. If no reservations exist, output a commented example showing the expected format.
7. Status/progress messages go to stderr.

### Behavior: `--set` Mode

1. Parse the entire input file (or stdin) before connecting. Extract all `host {}` blocks.
   Ignore any non-host directives.
2. Validate every entry:
   - Hostname must pass the library's name validation (see [Name Validation](#name-validation)).
   - MAC must be valid colon-separated hex.
   - IP must be valid IPv4.
   - No duplicate hostnames within the file.
   - No duplicate MACs within the file.
   - No duplicate IPs within the file.
   - On any validation error, print all errors with context (line numbers) and exit before
     connecting.
3. Connect to the router.
4. Fetch the current LAN configuration and reservation table.
5. Validate every entry against the device:
   - IP must fall within the LAN subnet. If any entry fails, print the subnet mismatch
     report (below) and exit 1 without writing.
   - IP must not be the router's own LAN address.
   - An IP inside the DHCP dynamic pool is a **warning, not an error**. dnsmasq honors a
     static host entry whose address lies inside the dynamic range and excludes that
     address from dynamic allocation, so this is safe on this device even though it would
     be a conflict under ISC dhcpd. UniFi permits fixed IPs inside the pool as well, so
     rejecting them would fail valid dumps over a hazard that does not exist here. The
     warning is worth printing because it is still untidy, and because a later pool change
     is easier to reason about when reservations sit outside it.

On subnet mismatch, print both remedies and write nothing:

```
error: subnet mismatch
  host file:   192.168.4.0/24 (38 of 38 entries)
  router LAN:  192.168.8.1/24

Resolve by either:
  - Setting the router's LAN address to 192.168.4.1/24, with
    goglnet --set-ip 192.168.4.1 --set-mask 255.255.255.0 --set-start ... --set-end ...
    or in the GL.iNet admin panel (LAN -> Router IP Address), then re-running. Your
    management session will drop and you will need to reconnect at the new address.
  - Renumbering the host file into 192.168.8.0/24 before re-running.
```
6. Refuse with `ErrDomainNotSet` if the host file carries no domain, naming
   `goglps --domain <domain>` as the remedy. This happens before any write.
7. Diff the **bindings**, keyed by MAC:
   - **Skip if unchanged**: MAC already bound to the same IP. Print skip to stderr.
   - **Update if changed**: MAC exists but the IP differs. Print update to stderr.
   - **Create if new**: MAC has no existing binding. Print create to stderr.
8. Diff the **names** independently, keyed by name, against the router's host file. A name that
   is absent, or present with a different address, is written. The diff is separate because a
   binding creates no DNS record: an entry whose binding already matches can still be missing
   its name, and folding the two diffs together would suppress that.
9. If `--prune`, delete bindings on the router whose MAC does not appear in the file, and remove
   their names. Without `--prune`, extras are left alone and counted.
10. Write: bindings one at a time, then every name change in a single `dns.set_host` call. The
    batch is not an optimization -- that endpoint replaces the whole file, so per-name writes
    would be one read-modify-write cycle each, racing one another.
11. Print summary to stderr:
    `N host declarations: N created, N updated, N skipped; N pruned; N DNS name(s) written, N removed; N errors`.
12. Exit 1 if any errors occurred.

### Behavior: `--add` Mode

Add a single host using an ISC DHCP declaration fragment. The fragment can be provided as:

- A positional argument (quoted string):
  ```bash
  goglps -H 192.168.8.1 --add 'host mydevice {
      hardware ethernet aa:bb:cc:dd:ee:ff;
      fixed-address 192.168.8.50;
  }'
  ```
- Stdin (if no positional argument):
  ```bash
  echo 'host mydevice {
      hardware ethernet aa:bb:cc:dd:ee:ff;
      fixed-address 192.168.8.50;
  }' | goglps -H 192.168.8.1 --add
  ```

Behavior:

1. Parse the single host declaration. Validate hostname, MAC, IP.
2. Connect to the router.
3. Check for conflicts (unless `--force`):
   - Is the IP already reserved for a different MAC?
   - Is the MAC already reserved a different IP?
   - Is another reservation already using this name?
   - Is the IP outside the LAN subnet, or the router's own address?

   An IP inside the DHCP dynamic pool warns but is not a conflict, and `--force` is not
   needed to proceed past it. See `--set` step 5 for why.
4. Refuse with `ErrDomainNotSet` if no domain is configured, naming the remedy.
5. Create or update the binding, then write the DNS name.
6. Print the result to stdout:
   ```
   Created: mydevice aa:bb:cc:dd:ee:ff 192.168.8.50
     DNS name mydevice.lab.example -> 192.168.8.50
   ```
   If the binding is written and the name write then fails, say which half succeeded: the
   difference determines whether the operator re-runs `--add` or repairs with `--set`.

### Behavior: `--del` Mode

Delete a host identified by one of `--name`, `--mac`, or `--ip`.

1. Connect to the router.
2. Find the matching reservation:
   - `--name <hostname>`: match the reservation name.
   - `--mac <mac>`: match the reservation MAC.
   - `--ip <ip>`: match the reservation IP.
3. If no match found, print error and exit 1.
4. If multiple matches found, list all matches and exit 1 unless `--force` is set.
5. Display what will be deleted and ask for confirmation (unless stdout is not a terminal,
   in which case proceed without prompting).
6. Remove the DNS name first, then the binding. The order matters: a leftover binding is an
   address with no name, which `--set` reports and repairs, while a leftover name keeps
   resolving to an address nothing holds.
7. Print the result to stdout:
   ```
   Deleted: mydevice aa:bb:cc:dd:ee:ff 192.168.8.50
     Removed the DHCP reservation and its DNS name
   ```

### Behavior: `--domain` Mode

Set the DNS domain. Required before any reservation write.

```bash
goglps -H 192.168.8.1 --domain lab.example
```

Behavior:

1. Validate the domain by the same rules as a hostname. Reject rather than escape: a quote in
   the marker line would corrupt the host file, and the host file is what the router resolves
   its own name from.
2. Read the host file and note the existing domain, if any.
3. Write the new domain into gogl's begin-marker line, and requalify every managed entry so no
   name keeps the old suffix. Resolution split across two domains is worse than either.
4. Report which of the three cases happened, because they have different consequences:
   ```
   DNS domain set to "lab.example"
   DNS domain changed from "old.example" to "lab.example"; existing host entries requalified
   DNS domain already "lab.example"
   ```

The domain lives on the router rather than in a config file, so it travels with the device. The
firmware has no dnsmasq domain setting to use, and `/ubus` -- the standard OpenWrt route to one
-- returns 404 on this model.

### Behavior: `--clear` Mode

Delete every reservation and every managed DNS name.

```bash
goglps -H 192.168.8.1 --clear
goglps -H 192.168.8.1 --clear --dry-run
```

Behavior:

1. List both the bindings and the managed host entries.
2. If both are empty, print `nothing to clear` and exit 0. "Make sure there is nothing there"
   is a reasonable request, so an empty device is a no-op rather than an error.
3. Print every binding and every name that will go.
4. With `--dry-run`, stop here.
5. Confirm interactively unless `--force` or stdout is not a terminal.
6. Remove the names first, then the bindings, for the reason given in `--del` step 6.
7. Print `Deleted N reservations and N DNS names`.

Both tables, because they are one intent stored in two places. The domain survives: it is
configuration rather than content, and re-setting it after every clear would be a papercut with
no purpose.

This is also the precondition for `goglnet --set-*`, which refuses to renumber the LAN while
any reservation exists.

### Error Handling

| Condition | Behavior |
|-----------|----------|
| No mode flag or multiple mode flags | Print usage, exit 1 |
| Missing password or host | Print error, exit 1 |
| Connection failure | Print error, exit 1 |
| Authentication failure | Print error distinguishing bad password from unreachable, exit 1 |
| Invalid host declaration syntax | Print error with line number, exit 1 |
| Hostname fails name validation | Print error naming the offending character, exit 1 |
| Duplicate hostname/MAC/IP in input file | Print all duplicates, exit 1 before connecting |
| Any IP outside the LAN subnet | Print subnet mismatch report, exit 1 before writing |
| IP inside the DHCP dynamic pool | Warning to stderr, proceed; never an error |
| Conflict on `--add` without `--force` | Print conflict details, exit 1 |
| No match on `--del` | Print error, exit 1 |
| Multiple matches on `--del` without `--force` | List matches, exit 1 |
| Individual create/update failure in `--set` | Print error to stderr, continue, exit 1 at end |

### Project Layout

```
utilities/
  goglps/
    main.go           # Entry point, flag parsing, mode dispatch
    parse.go          # ISC DHCP format parser
    format.go         # ISC DHCP format output
    operations.go     # get/set/add/del business logic
  docs/
    goglps/
      DESIGN.md       # Detailed design document (generated from this spec)
```

## goglnet Tool

A command-line tool for reporting the travel router's LAN address, DHCP pool, and the DNS
resolvers it advertises to clients. Lives in `utilities/goglnet/`. Built on the gogl module. The counterpart to
`gofinet`.

### Purpose

Before reserving addresses you need to know which addresses are safe to reserve. The DHCP
pool boundaries define what the router hands out dynamically; everything else in the
subnet is available for static reservation. `goglnet` is the companion to `goglps` in
exactly the way `gofinet` is the companion to `gofips`.

It is also the first thing to run when deciding whether a `gofi` dump can be applied at
all, since it reports the subnet that a dump's addresses must fall within.

### CLI Interface

```
goglnet [connection flags]
goglnet [connection flags] -j
```

### Output Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | Output in JSON format instead of text |

### Behavior

1. Connect to the router.
2. Fetch the LAN configuration via `client.Network().Get()`.
3. Fetch the reservation count via `client.Reservations().List()`.
4. Fetch model and firmware version via `client.System().Info()`.
5. Report: model, firmware, LAN address and netmask, subnet in CIDR form, DHCP enabled,
   pool start and end, lease time, DNS servers, reservation count, and the count of
   addresses in the subnet that are neither pooled nor reserved nor the router itself.
   Also list any reservation whose address falls inside the DHCP pool: dnsmasq honors those and
   excludes them from allocation, so nothing is broken, but they are missing from the available
   count and the discrepancy is otherwise unexplainable. They arise by accident -- a subnet move
   put twenty of twenty-seven inside the pool on real hardware, because the firmware rewrites
   host parts without regard to where the pool is.
6. Output text or JSON (`-j`).
7. With `--set-start` and `--set-end` alone, write the DHCP pool. The address and netmask are
   read from the device. Never refused, and the session survives: nothing moves.
8. With `--set-ip` and `--set-mask` (which require both pool bounds too), move the subnet.
   Refused while reservations exist unless `--force`; the session drops. Warn when the new pool
   would enclose existing reservations.

   Unlike a wireless write, this has no `y/N` gate -- only `--dry-run`. The asymmetry is
   deliberate but arguable: a wireless change can leave the device unreachable with no address
   to try, while a moved router is reachable at an address `goglnet` prints. If that reasoning
   stops holding, add the prompt here rather than removing it there.

9. Report each wireless interface: band, interface name, SSID, encryption, guest flag, hidden
   and enabled state. The passphrase is masked unless `--show-key` is given.
10. Report what each radio advertises: selectable channels with DFS marked, bandwidths,
   hardware modes, encryptions, and the transmit-power levels.
11. Write wireless identity with `--set-ssid`, `--set-key`, `--set-encryption`,
    `--set-hidden=true|false`, `--set-enabled=true|false`, each requiring `--iface`.
12. Write radio tuning with `--set-channel`, `--set-htmode`, `--set-hwmode`, `--set-txpower`,
    each requiring `--device`.
13. Refuse any wireless write when the session arrives over WiFi, and confirm with `y/N`
    unless `--yes`. Warn when the chosen channel is a DFS channel. See
    [Wireless Writes](#wireless-writes).

The domain is not part of this. It lives in the host file and is `goglps --domain`, because
the firmware has no domain setting to report or write.

Wireless lives in `goglnet` rather than a fourth utility because it is network configuration,
and because the three-tool mirror of `gofi` is worth more than a tidier separation. `goglnet`
is the tool that reports and writes what the network *is*; an SSID is part of that.

Where the API exposes a guest network, its subnet is reported, because it constrains address
planning. `gogl` only ever writes the main LAN.

### Text Output Format

```
MODEL       GL-SFT1200
FIRMWARE    4.3.28
LAN         192.168.8.1/24  (255.255.255.0)
DHCP        enabled
POOL        192.168.8.100 - 192.168.8.249  (150 addresses)
LEASE       12h
DOMAIN      lan
DNS         192.168.8.1
RESERVED    38
AVAILABLE   65
```

Rules:
- `AVAILABLE` counts host addresses in the subnet excluding the pool, existing
  reservations, the router's own address, and the network and broadcast addresses.
- When DHCP is disabled, `POOL` and `LEASE` show `(disabled)`.
- Columns are space-aligned via `text/tabwriter`.
- Status/progress messages go to stderr.

### JSON Output Format

A single JSON object with `model`, `firmware`, `lan_ip`, `netmask`, `subnet`,
`dhcp_enabled`, `dhcp_start`, `dhcp_stop`, `dhcp_lease`, `domain`, `dns`,
`reserved_count`, `available_count`, and `guest_subnet` where applicable. Empty/zero
fields are omitted (`omitempty`).

`gofinet` emits an array because a UDM Pro has many networks. `goglnet` emits an object
because the travel router has one LAN. Consumers of both must handle the difference.

### Project Layout

```
utilities/
  goglnet/
    main.go           # Entry point, flag parsing, connection
    operations.go     # LAN and DHCP fetch and flattening
    format.go         # Text and JSON output
```

## goglmac Tool

A command-line tool for listing connected clients on a GL.iNet travel router, filtered by
connection type (wired or WiFi), with independent OUI manufacturer lookup. Lives in
`utilities/goglmac/`. Built on the gogl module. The counterpart to `gofimac`.

### Purpose

To reserve an address for a device you need its MAC address, and to tell devices apart you
need to know who made them. `goglmac` performs its own OUI lookup using the IEEE OUI
database rather than trusting whatever the router reports, so manufacturer identification
is current regardless of the firmware's vintage. On an OpenWrt 18.06 base this matters
more than on the UniFi side, not less.

Its practical role in the replication workflow is discovering the MAC addresses to put
into a host file in the first place.

### OUI Database

Identical behavior to `gofimac`, so the two tools share a cache format but not a cache
directory.

The IEEE publishes the canonical OUI database at
`https://standards-oui.ieee.org/oui/oui.txt`. The first 3 octets of a MAC address identify
the manufacturer.

**Storage location**: `$XDG_DATA_HOME/goglmac/oui.txt`, falling back to
`~/.local/share/goglmac/oui.txt` if `$XDG_DATA_HOME` is not set. This follows the XDG Base
Directory Specification and requires no root access.

**Freshness check**: On every invocation, before performing any lookups:

1. Check if the OUI file exists at the storage location.
2. If it exists, check the file modification time. If the file is older than 30 days,
   re-download it.
3. If it does not exist, download it.
4. Download from `https://standards-oui.ieee.org/oui/oui.txt`.
5. If the download fails and a cached file exists (even if stale), use the cached file and
   print a warning to stderr.
6. If the download fails and no cached file exists, exit with an error.

**Parsing**: The IEEE OUI file format has entries like:

```
AA-BB-CC   (hex)		Acme Corporation
AABBCC     (base 16)		Acme Corporation
				123 Main Street
				Springfield IL 12345
				US
```

Parse only the `(hex)` lines. Extract the 3-octet prefix (normalized to lowercase
colon-separated, e.g. `aa:bb:cc`) and the manufacturer name.

**Lookup**: Given a MAC address, extract the first 3 octets and look up the manufacturer.
If not found, return `unknown`.

Randomized MAC addresses (locally-administered bit set in the first octet) resolve to
`randomized` rather than `unknown`, since they will never appear in the IEEE database and
are a poor choice to reserve an address for. `--get` output from `goglps` is unaffected;
this is a reporting nicety.

### CLI Interface

```
goglmac [connection flags] --wifi
goglmac [connection flags] --wired
goglmac [connection flags] --all
goglmac [connection flags] --wifi -j
```

### Mode Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--wifi` | `-w` | List only WiFi-connected clients |
| `--wired` | `-e` | List only wired (ethernet) clients |
| `--all` | `-a` | List all connected clients (default if no mode specified) |

If no mode flag is given, `--all` is assumed.

### Output Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | Output in JSON format instead of text |
| `--reserved` | `-r` | Mark which clients already have a reservation |

### Text Output Format

Default output is one line per client, tab-separated:

```
MAC              IP              HOSTNAME        OUI-MANUFACTURER
```

Concrete example:

```
aa:bb:cc:dd:ee:01	192.168.8.10	myserver	Dell Inc.
aa:bb:cc:dd:ee:02	192.168.8.11	printer 	Hewlett Packard
aa:bb:cc:dd:ee:03	192.168.8.112	unknown 	randomized
```

Rules:
- MAC addresses are lowercase, colon-separated.
- If the client has a name set, use it as the hostname. Otherwise use `unknown`.
- The OUI manufacturer comes from our own OUI database lookup, never from the router's
  reported value.
- Sort output by IP address numerically.
- Clients without an IP address are listed at the end with IP shown as `-`.
- With `--reserved`, a trailing column shows `reserved` or `dynamic`.
- Status/progress messages (OUI download progress, etc.) go to stderr.

### JSON Output Format

When `--json` or `-j` is specified, output a JSON array to stdout:

```json
[
  {
    "mac": "aa:bb:cc:dd:ee:01",
    "ip": "192.168.8.10",
    "hostname": "myserver",
    "manufacturer": "Dell Inc.",
    "is_wired": false,
    "reserved": true,
    "online": true,
    "rx_bytes": 123456789,
    "tx_bytes": 987654321
  }
]
```

JSON fields:
- Always present: `mac`, `hostname`, `manufacturer`, `is_wired`, `online`.
- Present when known: `ip`, `reserved`, `rx_bytes`, `tx_bytes`, `signal`, `band`.
- Omit fields with zero/empty values (`omitempty`).

The field set is narrower than `gofimac`'s because the router reports less than a UniFi
controller does. Fields are added only when a recorded fixture shows the API returning
them; nothing is invented to match `gofimac`'s output shape.

### Behavior

1. Check and update the OUI database (see Freshness check above).
2. Parse the OUI database into a lookup map (3-octet prefix -> manufacturer name).
3. Connect to the router.
4. Fetch clients via `client.Clients().List()`.
5. Filter by connection type.
6. If `--reserved`, fetch the reservation table and mark each client.
7. For each client, look up the manufacturer from the OUI map using the first 3 octets of
   the MAC.
8. Sort by IP address numerically (clients without IPs sort last).
9. Output in the requested format (text or JSON).

### Error Handling

| Condition | Behavior |
|-----------|----------|
| Missing password or host | Print error, exit 1 |
| Connection failure | Print error, exit 1 |
| OUI download fails, no cache | Print error, exit 1 |
| OUI download fails, stale cache exists | Print warning to stderr, continue with cached data |
| OUI parse error | Print error, exit 1 |
| No clients found matching filter | Print empty output (empty line in text, `[]` in JSON), exit 0 |

### Project Layout

```
utilities/
  goglmac/
    main.go           # Entry point, flag parsing
    oui.go            # OUI database download, parse, lookup
    format.go         # Text and JSON output formatting
    operations.go     # Client listing and filtering
  docs/
    goglmac/
      DESIGN.md       # Detailed design document
```

## Commands

```bash
make test      # Run all tests
make lint      # Run linter
make build     # Build the module and utilities into bin/
make coverage  # Generate coverage report
make install   # Install utilities to ~/bin (override INSTALL_DIR)
```

`UTILITIES := goglmac goglnet goglps`

Which mirrors `gofi`'s `UTILITIES := gofimac gofinet gofips` one-for-one.

## Reference Implementations

Study these for patterns. None is authoritative; the recorded fixtures from a live
SFT1200 are.

- `github.com/emergingrobotics/gofi` - the counterpart module (imported as
  `github.com/unifi-go/gofi`; the repository path and the module path differ). Layout,
  service interface style, CLI conventions, and the ISC DHCP format implementation all come
  from here.
- `github.com/tomtana/python-glinet` - Python JSON-RPC client. The most complete public
  implementation of the challenge/login flow, including the `alg` map and the keepalive
  thread.
- `github.com/ryanrishi/glinet-client-go` - Go client for firmware 4.x, MIT licensed.
  Useful for the crypt digest implementation in Go specifically.
- `github.com/metril/ha-glinet` - Home Assistant integration for firmware 4.x. Useful for
  observed group and method names, to be confirmed against fixtures rather than trusted.
