# gogl - Go GL.iNet Travel Router

A Go module for programmatic control of GL.iNet travel routers running firmware 4.x,
targeting the **GL-SFT1200 (Opal)**.

`gogl` is the travel-router counterpart to
[`gofi`](https://github.com/emergingrobotics/gofi). It exposes the same shape of API and
the same command-line ergonomics, so that a network described on a UniFi UDM Pro can be
reproduced on a pocket router. Dump the fixed-IP and DNS assignments from a UniFi site
with `gofips`, hand the file to `goglps`, and the travel router serves the same addresses
and the same names.

## Target Device

| Property | Value |
|----------|-------|
| Model | GL-SFT1200 ("Opal") |
| Firmware track | 4.x (4.8.3 at time of writing) |
| SoC | SiFlower SF19A28 |
| OpenWrt base | 18.06 |
| Ethernet | 1 x GbE WAN, 2 x GbE LAN |
| Wireless | AC1200 dual-band (2.4 GHz 300 Mbps, 5 GHz 867 Mbps) |
| Default LAN | `192.168.8.1/24` |
| Admin user | `root` (password set during first-run setup; there is no default) |

Firmware 4.x is a hard requirement. GL.iNet replaced the older REST API with JSON-RPC at
firmware 4.0, and `gogl` speaks only JSON-RPC. Firmware 3.x devices are out of scope.

Other GL.iNet 4.x models should work, since the API is shared across the firmware line,
but the SFT1200 is the only device this module is specified and tested against.

## Project Resources

- **API Reference**: `GL_INET_4X_API_DOCUMENTATION.md` - Endpoint documentation, captured
  and verified against a live SFT1200 (see [API Discovery](#api-discovery)).
- **Architecture**: `docs/DESIGN.md` - Design with mermaid diagrams, service interfaces,
  type definitions.
- **Implementation Plan**: `docs/plan.md` - Phased plan with progress tracking.

## Critical Rules

1. **Every function MUST have a test** - No exceptions. Run `make test` to verify.
2. **Every endpoint MUST be supported in the mock server** - Tests use the mock, not real hardware.
3. **No phase advancement without 100% test coverage** - Complete and test each phase before moving on.
4. **Phases are sequential** - Follow `docs/plan.md` in order.
5. **No SSH, no UCI, no shell** - Every operation goes through the JSON-RPC API. Anything
   the API cannot express is a documented gap, never a shell workaround. This keeps the
   entire surface reachable by the mock server.

## Scope

`gogl` v1 manages, on a single flat LAN with no VLANs:

- **Network information** - LAN address, netmask, DHCP pool boundaries, lease time, DNS servers.
- **Static IP reservations** - MAC-to-IP bindings.
- **Local DNS** - names for the hosts whose addresses it reserves.
- **Connected clients** - with independent IEEE OUI manufacturer lookup.

Explicitly **out of scope for v1**:

- **Wireless configuration.** SSIDs, passphrases, and radio settings are managed through
  the GL.iNet admin panel. `gogl` never writes them. Writing wireless config over a
  wireless management session is a good way to lock yourself out of the device.
- **VLANs.** The Opal exposes no VLAN configuration through its API. VLANs are reachable
  only via LuCI/UCI and `swconfig` on this SoC, which rule 5 forbids and which users
  report as unreliable on the SiFlower switch.
- **Guest network.** Reported by `goglnet` for address planning if the API exposes it,
  never written.
- **VPN, firewall, port forwarding, traffic rules.**
- **Changing the router's LAN address.** See [Subnet Mismatch](#subnet-mismatch).

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
│   │   ├── reservation.go # Static leases (DHCP binding + DNS name)
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
└── utilities/         # CLI tools: goglps, goglmac, goglnet, goglsync
```

### Differences from gofi's service layer

- **No site concept.** GL.iNet routers have no equivalent of a UniFi site. Every `site`
  parameter present in `gofi`'s method signatures is absent here.
- **No separate DNS service.** On the Opal a reservation *is* the DNS record; see
  [The Unified Reservation](#the-unified-reservation). `gofi` needs `Users()` and `DNS()`
  kept in sync. `gogl` needs neither.
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

## The Unified Reservation

This is the most important structural difference from `gofi`, and it makes `gogl` simpler
than its counterpart.

The Opal runs dnsmasq. A static lease in dnsmasq carries a name alongside the MAC and IP,
and dnsmasq answers DNS queries for that name. GL.iNet's admin panel surfaces this as
**LAN → Address Reservation**, with Name, MAC Address, and IP Address fields, and writes
all three into the dnsmasq configuration.

So one reservation entry provides both the DHCP binding and the DNS record, atomically:

```
Reservation{Name: "myserver", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.10"}
```

...yields a device that always receives `192.168.8.10`, and a name `myserver` that
resolves to `192.168.8.10` for every client of the router.

Consequences, all of which the implementation depends on:

- **Drift is impossible.** `gofips` must keep a UniFi user record's fixed IP and a
  separate static DNS record consistent, and can find them disagreeing. `goglps` writes
  one entry. There is nothing to reconcile.
- **`--keep-dns` is meaningless and is not implemented.** You cannot retain the name
  while removing the reservation; they are the same object.
- **Names are not optional in practice.** A reservation with an empty name is a valid
  DHCP binding with no DNS. `goglps` always writes a name, because the ISC DHCP format it
  consumes is keyed by hostname.
- **DNS is limited to reserved hosts.** A host with a hand-configured static IP outside
  the reservation table gets no name. This is a documented limitation, and it matches the
  intended model: the router serves DNS for exactly the hosts whose DHCP it manages.

### Domain Suffix

The router's dnsmasq domain defaults to `lan`, so a reservation named `myserver` resolves
as both `myserver` and `myserver.lan`. `goglnet` reports the configured suffix. `goglps`
stores and emits bare names, never fully-qualified ones.

### Name Validation

GL.iNet writes the reservation name directly into the dnsmasq configuration file, and a
known firmware defect lets characters such as `"` corrupt that file, breaking DHCP and DNS
for the entire router until it is repaired.

`gogl` therefore validates every name before any write, and rejects rather than escapes:

- Permitted characters: `[a-zA-Z0-9-]`, plus `.` as a label separator.
- Must not begin or end with `-` or `.`.
- Each label at most 63 characters; total at most 253.
- Must not be empty.

This validation lives in the library, not in the CLI, so no consumer can bypass it.

**This is deliberately stricter than `gofips`.** `gofips` accepts `_` in hostnames, which
is legal in a UniFi record but is not a legal DNS label character and has no business in a
dnsmasq name. So a host file exported from a UniFi controller that contains a name like
`my_server` will be **rejected** by `goglps`, with the offending name and character
reported, rather than silently rewritten to `my-server`.

This is the one place where a `gofips` file is not guaranteed to import unchanged. Renaming
is a decision about what your hosts are called, and `gogl` will not make it silently. Fix
it on the UniFi side, or edit the file.

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
containing a name that is legal on UniFi but not in DNS is rejected rather than converted.
See [Name Validation](#name-validation).

### End-to-End Workflow

This example assumes a home network of `192.168.4.0/24` behind a UDM Pro at
`192.168.4.1`, and a travel router whose LAN address has already been set to
`192.168.4.1/24`. That last part is a prerequisite, not something `goglsync` does; see
[Subnet Mismatch](#subnet-mismatch).

```bash
# On the UniFi side: dump the network shape and the host bindings.
gofinet -H 192.168.4.1 -k -j    > home.net.json
gofips  -H 192.168.4.1 -k --get > home.hosts

# Inspect what the travel router currently looks like. Note the different address:
# the travel router is reached over its own LAN, not the UDM's.
goglnet -H 192.168.4.1

# Preview every change, touching nothing.
goglsync -H 192.168.4.1 --net home.net.json --hosts home.hosts --dry-run

# Apply.
goglsync -H 192.168.4.1 --net home.net.json --hosts home.hosts
```

The two `-H 192.168.4.1` values refer to different devices on different occasions: the UDM
Pro while you are at home taking the dump, and the travel router later. They coincide
because the travel router is impersonating the UDM's LAN address, which is the entire
point. Do not run both halves of this workflow with both devices on the same segment.

The result is a travel router that hands out the same addresses to the same MAC addresses
and answers the same names as the network it was copied from. Plug it into hotel ethernet,
and devices that expect `myserver` at `192.168.4.10` find it there.

### What Is Reproduced

| From the UniFi dump | On the travel router | Notes |
|---------------------|----------------------|-------|
| `dhcpd_start` / `dhcpd_stop` | DHCP pool start / limit | dnsmasq expresses the pool as start plus a count; `goglsync` converts. |
| `dhcpd_leasetime` | Lease time | Seconds converted to dnsmasq duration form. |
| `dhcpd_dns_1` .. `dhcpd_dns_4` | DNS servers | The Opal accepts fewer; excess entries are reported, not silently dropped. |
| `domain_name` | dnsmasq domain suffix | |
| Each `host {}` block | One address reservation | Supplies both the DHCP binding and the DNS name. |

### What Is Not Reproduced

`goglsync` reports each of these explicitly rather than failing silently. A dump
containing them still applies; the report names what was left behind.

| From the UniFi dump | Why not |
|---------------------|---------|
| VLANs / additional networks | Out of scope. A dump with more than one LAN-purpose network is an error, not a partial apply. |
| Wireless config | Out of scope for v1. |
| Firewall rules, port forwards, traffic rules | Out of scope for v1. |
| WAN configuration | The travel router's WAN is whatever it is plugged into. |
| IPv6, mDNS, IGMP snooping, DHCP guard, ARP inspection | No API equivalent within v1 scope. |
| DNS records not backed by a reservation | See [The Unified Reservation](#the-unified-reservation). |

### Subnet Mismatch

A UniFi LAN is typically something like `192.168.4.0/24`. The Opal ships on
`192.168.8.0/24`. Reservations at `192.168.4.x` are meaningless on a router whose LAN is
`192.168.8.0/24`, so the two must agree before any reservation can be written.

**`gogl` v1 does not resolve this automatically. It detects it and refuses.**

The reason is that both possible fixes are disruptive, and neither should happen as a side
effect of a sync:

- **Re-IP the router** to match the dump. This changes the address you are managing the
  router at, so it drops your management session mid-operation, and your own machine holds
  a DHCP lease on the old subnet until it renews. A tool that does this has to survive
  losing contact with the device it is configuring, partway through a multi-step change.
- **Renumber the host file** to the router's subnet, rewriting the network part of every
  address while preserving host parts. Non-destructive to the router, but it silently
  changes the addresses you asked for, which is not something to do without being asked.

So when `goglsync` finds that the dump's subnet and the router's LAN subnet differ, it
prints both remedies and exits 1 before writing anything:

```
error: subnet mismatch
  dump subnet:   192.168.4.0/24 (from home.net.json)
  router LAN:    192.168.8.1/24
  38 of 38 reservations fall outside the router's LAN.

Resolve by either:
  - Setting the router's LAN address to 192.168.4.1/24 in the GL.iNet admin panel
    (LAN -> Router IP Address), then re-running. Your management session will drop
    and you will need to reconnect at the new address.
  - Renumbering home.hosts into 192.168.8.0/24 by hand before re-running.
```

`goglps` applies the same rule per entry: a reservation whose IP falls outside the current
LAN subnet is rejected with the specific conflict named. Neither tool ever writes the
router's LAN address.

### Future Work

Two tools would close the gap, and the design above leaves room for both:

- **`goglreip`** - change the router's LAN address as a deliberate, single-purpose,
  confirmed operation, with the reconnection handling that requires. Kept separate from
  `goglsync` so that the destructive step is never incidental to a sync.
- **`goglrenumber`** - rewrite the network part of every address in an ISC DHCP file to a
  target subnet, preserving host parts, emitting a new file rather than editing in place.
  A pure text transformation with no device access, and therefore the cheaper of the two
  to build and trust.

Neither is in v1. `goglsync`'s error message points at the manual equivalent of each.

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

A command-line tool for managing static IP reservations and their DNS names on a GL.iNet
travel router, using ISC DHCP host declaration format. Lives in `utilities/goglps/`.
Built on the gogl module. The counterpart to `gofips`.

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
| `--prune` | `-P` | On `--set`, delete reservations present on the router but absent from the file |
| `--dry-run` | n/a | Show what would be done without making changes |

`gofips`'s `--keep-dns` flag has no analogue. On this device the DNS name and the
reservation are one object, so retaining one while deleting the other is not expressible.

### Behavior: `--get` Mode

1. Connect to the router.
2. List all reservations via `client.Reservations().List()`.
3. For each reservation, the hostname is the reservation's own name. There is no fallback
   chain and no cross-referencing, because there is no second record that could disagree.
   A reservation with an empty name is emitted with the MAC as its hostname, colons
   replaced by hyphens (e.g. `aa-bb-cc-dd-ee-ff`), and a comment noting that the router
   serves no DNS for it. Be aware that this does not round-trip inertly: feeding that file
   back through `--set` gives the reservation that MAC-derived name for real, and the
   router starts answering DNS for it. The emitted comment says so.
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
   - IP must fall within the LAN subnet. If any entry fails, report the mismatch as in
     [Subnet Mismatch](#subnet-mismatch) and exit 1 without writing.
   - IP must fall outside the DHCP dynamic pool, or the router may hand the same address
     to another client. Report as an error per entry; `--force` downgrades it to a warning.
   - IP must not be the router's own LAN address.
6. For each host declaration:
   - **Skip if unchanged**: MAC already reserved with the same IP and name. Print skip to stderr.
   - **Update if changed**: MAC exists but IP or name differs. Update the reservation.
     Print update to stderr.
   - **Create if new**: MAC has no existing reservation. Create it. Print create to stderr.
7. If `--prune`, delete reservations on the router whose MAC does not appear in the file.
   Without `--prune`, extra reservations on the router are left alone and counted in the
   summary.
8. Print summary to stderr: `N processed, N skipped, N created, N updated, N pruned, N errors`.
9. Exit 1 if any errors occurred.

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
   - Is the IP inside the DHCP dynamic pool, outside the LAN subnet, or the router's own
     address?
4. Create or update the reservation.
5. Print the result to stdout:
   ```
   Created: mydevice aa:bb:cc:dd:ee:ff 192.168.8.50
   ```

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
6. Delete the reservation. This removes the DHCP binding and the DNS name together; there
   is no way to remove one and keep the other.
7. Print the result to stdout:
   ```
   Deleted: mydevice aa:bb:cc:dd:ee:ff 192.168.8.50
     Removed DHCP reservation and DNS name
   ```

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
| IP inside the DHCP dynamic pool | Error per entry; warning with `--force` |
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

A command-line tool for reporting the travel router's LAN address, DHCP pool, and DNS
settings. Lives in `utilities/goglnet/`. Built on the gogl module. The counterpart to
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
   pool start and end, lease time, domain suffix, DNS servers, reservation count, and the
   count of addresses in the subnet that are neither pooled nor reserved nor the router
   itself.
6. Output text or JSON (`-j`).

Where the API exposes a guest network, its subnet is reported read-only, because it
constrains address planning. `gogl` never writes it.

### Text Output Format

```
MODEL       GL-SFT1200
FIRMWARE    4.8.3
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

## goglsync Tool

A command-line tool that applies a complete `gofi` dump to a travel router in one
operation. Lives in `utilities/goglsync/`. Built on the gogl module. It has no `gofi`
counterpart; it is the tool that makes the two modules useful together.

### Purpose

`goglps --set` handles the host bindings, but a dump also carries the network shape: pool
boundaries, lease time, DNS servers, domain suffix. `goglsync` applies both, in the right
order, with one plan and one report. It is the single command that answers "make this
travel router behave like my home network."

### CLI Interface

```
goglsync [connection flags] --net <file> --hosts <file> [--dry-run]
goglsync [connection flags] --hosts <file>
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--net` | `-N` | Network shape, as produced by `gofinet -j` |
| `--hosts` | `-f` | Host bindings, as produced by `gofips --get` |
| `--dry-run` | n/a | Print the full plan and exit without writing |
| `--prune` | `-P` | Delete reservations on the router that are absent from the host file |
| `--yes` | `-y` | Skip the confirmation prompt |

At least one of `--net` or `--hosts` is required. Given only `--hosts`, `goglsync`
behaves as `goglps --set` plus the fidelity report.

### Behavior

1. Parse both input files completely before connecting. A `gofinet -j` array containing
   more than one LAN-purpose network is an error: v1 handles a single flat network, and
   silently picking one of several would be a guess.
2. Validate the host file exactly as `goglps --set` does.
3. Connect to the router and fetch current LAN config, reservations, and system info.
4. **Reconcile subnets.** If the dump's subnet differs from the router's LAN subnet, print
   the subnet mismatch report (see [Subnet Mismatch](#subnet-mismatch)) and exit 1 without
   writing anything.
5. Build a plan: every setting to change and every reservation to create, update, prune,
   or skip, plus every item in the dump that cannot be reproduced.
6. Print the plan. If `--dry-run`, exit 0 here.
7. Unless `--yes`, prompt for confirmation. If stdout is not a terminal, `--yes` is
   required rather than assumed; a sync is a bulk write and should not happen by accident
   in a pipeline.
8. Apply network settings first, then reservations. Network settings first because a pool
   change can free an address that a reservation needs.
9. Print the fidelity report and summary.

The router's LAN address is never written, in any mode. That is the one change `goglsync`
will not make.

### Plan Output

```
plan for GL-SFT1200 at 192.168.8.1 (firmware 4.8.3)

network:
  dhcp pool     192.168.8.100-192.168.8.249  ->  192.168.8.100-192.168.8.189
  lease time    12h                          ->  24h
  dns servers   192.168.8.1                  ->  1.1.1.1, 8.8.8.8
  domain        lan                          ->  unchanged

reservations:
  create   myserver     aa:bb:cc:dd:ee:01  192.168.8.10
  update   printer      aa:bb:cc:dd:ee:02  192.168.8.11  (was 192.168.8.12)
  skip     nas          aa:bb:cc:dd:ee:03  192.168.8.13  (unchanged)
  prune    oldlaptop    aa:bb:cc:dd:ee:99  192.168.8.44  (absent from host file)

not reproduced:
  dhcpd_dns_3, dhcpd_dns_4    router accepts 2 DNS servers; 3rd and 4th dropped
  igmp_snooping               no API equivalent
  mdns_enabled                no API equivalent

37 host declarations: 12 create, 3 update, 22 skip
1 router reservation pruned
3 network settings change
3 dump fields not reproduced
```

Host-file counts and router-side counts are reported on separate lines because they count
different things: `create`/`update`/`skip` partition the declarations in the file, while
`prune` counts reservations that exist only on the router. Summing them would be
meaningless.

### Fidelity Report

The `not reproduced` section is not optional output and is not suppressible. Every field
present in the dump that `gogl` did not write appears there with the reason. A sync that
reproduces everything prints `not reproduced: nothing`. The report is what makes the lossy
projection honest: you always know what your travel router is not doing.

### Error Handling

| Condition | Behavior |
|-----------|----------|
| Neither `--net` nor `--hosts` given | Print usage, exit 1 |
| Input file unreadable or malformed | Print error with file and line, exit 1 |
| `--net` file contains more than one LAN network | Print the networks found, exit 1 |
| Subnet mismatch between dump and router | Print both remedies, exit 1 before writing |
| Host file validation failure | Print all errors, exit 1 before connecting |
| Not a terminal and `--yes` absent | Print error, exit 1 |
| Network settings write failure | Print error, exit 1 without touching reservations |
| Individual reservation failure | Print error to stderr, continue, exit 1 at end |

### Project Layout

```
utilities/
  goglsync/
    main.go           # Entry point, flag parsing
    dump.go           # gofinet -j and gofips ISC DHCP input parsing
    plan.go           # Diff current state against desired state
    apply.go          # Ordered application of the plan
    report.go         # Plan output and fidelity report
  docs/
    goglsync/
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

`UTILITIES := goglmac goglnet goglps goglsync`

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
