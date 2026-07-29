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

Three tools, mirroring `gofi`'s three one-for-one.

| gogl | gofi counterpart | What it does | Writes? |
|------|------------------|--------------|---------|
| `goglps` | `gofips` | DHCP reservations, in ISC DHCP format | **yes** |
| `goglnet` | `gofinet` | LAN address, DHCP pool, lease time, resolvers | **yes**, with `--set-*` |
| `goglmac` | `gofimac` | Connected clients with IEEE OUI manufacturer lookup | no |

**What gogl writes:** reservations, host-file entries and the DNS domain (`goglps`), and the
LAN address and DHCP pool (`goglnet --set-*`). **What it does not:** wireless, VLANs, firewall,
VPN. Set those in the GL.iNet admin panel.

Two ordering rules protect you from the states that are painful to recover from:

- **A reservation write needs a domain first.** Otherwise you get addresses with no names and
  no hint that was accidental.
- **A network change is refused while reservations exist.** Renumbering the LAN underneath
  them leaves every reservation outside the new subnet — still listed, silently inert. Clear
  them with `goglps --clear`, apply the new network, then import.

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
goglps --clear                   # delete ALL reservations AND DNS names (prompts unless --force)
goglps --clear --dry-run         # list what would go
```

Changing the domain rewrites every managed host entry to the new suffix, so resolution does
not split between two domains.

`--clear` removes the DNS names as well as the reservations, because they are one intent stored
in two places. Leaving the names behind would strand records pointing at addresses the router
no longer reserves — and since clearing is what unblocks a renumber, those names would then be
answering for a subnet that no longer exists. The domain survives a clear: it is configuration,
not content.

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
them leaves every one outside the new subnet.

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
| Wireless configuration | Writing SSID or passphrase over a wireless management session is a good way to lock yourself out. Use the admin panel. |
| VLANs | The Opal exposes no VLAN configuration through its API. Reaching them means LuCI, UCI, and `swconfig` on a SiFlower switch, which is both outside this tool's rules and unreliable in practice. |
| Changing LAN address, DHCP pool, lease time, DNS servers | Read-only by design. Keeps a bulk import from ever leaving the router unreachable. |
| Setting dnsmasq's `domain` / `local` / `expandhosts` | No endpoint exposes them, and `/ubus` is 404 on this model. gogl carries its own domain and writes FQDNs into the host file instead, which works for any suffix. |
| DNS names via reservations | Reservations do not create DNS records. Use the host file — `goglps --domain` plus host entries. |
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
err := client.Network().Set(ctx, &types.Network{
    Interface: types.InterfaceLAN,
    LANIP:     "192.168.2.1",
    Netmask:   "255.255.255.0",
    DHCPStart: "192.168.2.100",
    DHCPStop:  "192.168.2.149",
})
```

Returns `ErrReservationsExist` if any reservation is present. The pool is checked against the
new subnet first: the firmware accepts a pool outside it and then silently hands out nothing.

Expect the call to fail with a lost connection *on success* — the router moves to the new
address mid-request. `goglnet` treats that as expected and tells you where to reconnect.

### Errors

```go
switch {
case errors.Is(err, gogl.ErrDomainNotSet):      // set a domain before writing reservations
case errors.Is(err, gogl.ErrReservationsExist): // clear reservations before changing the network
case errors.Is(err, gogl.ErrNotFound):
case errors.Is(err, gogl.ErrConflict):        // MAC already reserved
case errors.Is(err, gogl.ErrInvalidName):     // name would be unsafe in dnsmasq config
case errors.Is(err, gogl.ErrOutsideSubnet):
case errors.Is(err, gogl.ErrUnauthorized):
}

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
