# GL.iNet Firmware 4.x JSON-RPC API

Captured from a **GL-SFT1200 (Opal)** on 2026-07-28. The device reports model
`sft1200`, firmware **4.3.28**. (An earlier draft of this file said 4.8.3, which is
the newest release GL.iNet lists for the model, not what this unit runs.)

Everything here was observed against that device. GL.iNet's official 4.x reference is no
longer publicly reachable, and the three public client libraries
([`python-glinet`](https://github.com/tomtana/python-glinet),
[`ryanrishi/glinet-client-go`](https://github.com/ryanrishi/glinet-client-go),
[`metril/ha-glinet`](https://github.com/metril/ha-glinet)) all predate at least one
behavior documented below. Where this file and any of those disagree, this file was tested
against hardware.

## Endpoint

A single JSON-RPC 2.0 endpoint:

```
POST http://192.168.8.1/rpc
Content-Type: application/json
```

Plain **HTTP on port 80** by default. Not 443 — assuming the usual HTTPS admin port fails
here.

---

## Authentication

Challenge/response. The password is never transmitted, only a digest derived from it.

### Step 1: challenge

```json
{"jsonrpc":"2.0","id":1,"method":"challenge","params":{"username":"root"}}
```

Observed response:

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "nonce": "cVm4yDf1F2X7KUH2y0gbJESZpI4NmN9f",
    "alg": 5,
    "salt": "j2K6qjQJi8fCtAzO",
    "hash-method": "sha256"
  }
}
```

| Field | Observed | Meaning |
|-------|----------|---------|
| `alg` | `5` (number) | unix crypt(3) variant for the password: 1 = MD5, 5 = SHA-256, 6 = SHA-512. Some firmware sends this as a **string**, so decode tolerantly. |
| `salt` | 16 bytes, base64-ish alphabet | crypt salt, taken from the stored shadow entry |
| `nonce` | 32 bytes | single-use, short-lived |
| **`hash-method`** | `"sha256"` | **the digest applied to the login string** |

### `hash-method` is the field everyone misses

**This is the single most important finding here.** `hash-method` selects the hash applied
to the login string, and it is a *separate choice* from `alg`, which selects the crypt
applied to the password. On this firmware they differ in kind: `alg: 5` (SHA-256 crypt) and
`hash-method: "sha256"` (SHA-256 digest).

Every public client library hardcodes **MD5** for the digest, because they were written
against firmware that predates this field. Against 4.3.28 they all fail with
`-32000 Access denied`, which is indistinguishable from a wrong password and sends you
hunting in the wrong place for hours.

Rules:

- **`hash-method` absent** → MD5. This is what pre-4.8 firmware expects and what the public
  libraries implement.
- **`hash-method: "md5"`** → MD5.
- **`hash-method: "sha256"`** → SHA-256.
- **Anything else** → fail loudly. Never fall back to another hash: a wrong digest is
  reported identically to a wrong password.

### Step 2: derive the digest

```
cipher = crypt(password, "$" + alg + "$" + salt)        # full crypt output, with prefix
digest = H(username + ":" + cipher + ":" + nonce)       # H per hash-method
```

`cipher` is the **complete** crypt string including its `$alg$salt$` prefix, not just the
hash portion after the last `$`. Verified: the tail-only variant is rejected.

Worked example, from the capture above with a test password:

```
$ openssl passwd -5 -salt j2K6qjQJi8fCtAzO testpassword
$5$j2K6qjQJi8fCtAzO$<hash>

$ printf 'root:%s:%s' "$cipher" "$nonce" | sha256sum
```

### Step 3: login

```json
{"jsonrpc":"2.0","id":2,"method":"login","params":{"username":"root","hash":"<digest>"}}
```

Returns `{"sid": "..."}`.

### Timing

The nonce is valid for roughly **one second**, while SHA-256/SHA-512 crypt is deliberately
slow (5000 rounds). A single-challenge implementation races against its own hashing cost and
fails intermittently.

**Issue `challenge` twice**: once to learn the salt, then again after the crypt completes to
obtain a nonce that is fresh when the cheap digest runs over it. The salt is stable across
challenges, so the expensive cipher can be cached; the nonce cannot.

### Session

`sid` is passed as the **first element** of every subsequent call's `params` array. It idles
out after roughly **35 seconds**; any successful call resets the timer. An `alive` method
keeps it open.

---

## Brute-force lockout

**Observed the hard way.** After roughly a dozen failed logins, the firmware locks the
account and refuses even `challenge`:

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "error": {
    "message": "Login fail number over limit",
    "data": {"wait": 593},
    "code": -32003
  }
}
```

- `data.wait` is the remaining lockout in **seconds** — about ten minutes when first
  tripped.
- A **correct** password is refused while locked out.
- Retrying does not shorten it.

Consequences for any client: do not sweep authentication variants against real hardware, and
surface this code as its own distinct error. Reporting it as a generic authentication failure
is actively misleading, because the fix is to wait rather than to check the password.

---

## Calls

Every operation is the `call` method, with the group and method as data rather than in the
URL:

```json
{"jsonrpc":"2.0","id":3,"method":"call","params":["<sid>","<group>","<method>",{}]}
```

### Confirmed endpoints

All of these were called against the device. The full reference for all 43 groups and 313
methods is in [`docs/api/`](docs/api/README.md).

| Group | Method | Returns |
|-------|--------|---------|
| `system` | `get_info` | `model`, `firmware_version`, `mac`, `sn`, `board_info`, feature flags |
| `system` | `get_status` | `network[]` of `{interface, up, online}`, plus `wifi[]`, `service[]`, `client` |
| `clients` | `get_list` | `clients[]` — see the field notes below |
| `clients` | `get_status` | `cable_total`, `wireless_total` |
| `lan` | `get_config_list` | `interfaces[]` for **both** `lan` and `guest` |
| `lan` | `get_static_bind_list` | `static_bind_list[]` of `{name, mac, ip}` |
| `lan` | `add_static_bind` | `{}` — args `{mac, ip, name}` |
| `lan` | `set_static_bind` | `{}` — args `{mac, ip, name}`, matched on `mac` |
| `lan` | `remove_static_bind` | `{}` — args `{mode, mac}`; mode 0 = one, mode 1 = **all** |
| `network` | `get_dhcp_leases` | `leases[]` of `{mac, ip, hostname, expires}` |
| `wifi` | `get_config`, `get_status` | radio configuration and state |
| `dns` | `get_config` | `mode`, `server[]`, `force_dns`, `rebind_protection` |
| `dns` | `get_host` | `content` — the whole hosts(5) file as one string |
| `dns` | `set_host` | `[]` — args `{content}`; **replaces the entire file** |
| `repeater`, `led`, `qos`, `igmp`, `ddns`, `cloud`, `tor`, `upgrade` | `get_config` | per-feature settings |
| `iwinfo` | `devices`, `info` | ubus passthrough; `info` needs `{"device":"wlan0"}` |

### Where the device disagrees with the documentation

- `network.get_dhcp_leases` wraps its array in **`leases`**. The documentation says
  `entries`.
- `lan.get_config_list` returns more than documented — `leasetime`, `dns[]`, `gateway`,
  `lpr[]` — and reports `enable` as a **number** (`1`), not a boolean.
- `clients.get_list` reports the connection as **`iface`**: `"cable"`, `"2.4G"` or `"5G"`.
  There is no `is_wired` boolean. Cumulative counters are `total_rx`/`total_tx`; `rx`/`tx`
  are instantaneous rates. There is no `signal` field.

### A static bind does NOT create a DNS record

**Tested directly, and it matters.** GL.iNet's admin panel labels the static-bind `name`
field as a hostname, and it is tempting to assume dnsmasq serves it. It does not.

Two experiments on 4.3.28:

1. A bind for an absent MAC (`gogl-test-entry` -> `192.168.8.251`): the name did not resolve,
   bare or with the `.lan` suffix.
2. A bind for a **present, actively leased** device under a *new* name (`gogl-dnsprobe` on
   the Bouffalo's real MAC and current address): the new name still did not resolve, while
   the device's original **lease** hostname `Bouffalolab_bl606p-2a856f.lan` continued to
   resolve throughout.

So a static bind pins MAC to IP and nothing more. Its `name` is a label for the admin panel.

### DNS names DO work, through the host file

`dns.get_host` and `dns.set_host` read and write the router's `hosts(5)` file as one string,
and dnsmasq answers from it. This is how a client can create real DNS records.

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","dns","get_host",{}]}
-> {"result":{"content":"127.0.0.1 localhost\n\n::1  localhost ip6-localhost ip6-loopback\n..."}}

{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","dns","set_host",
  {"content":"127.0.0.1 localhost\n192.168.2.251 nas nas.lab.example\n"}]}
-> {"result":[]}
```

Verified on hardware:

- A written entry resolves by its **bare name** and by an **arbitrary FQDN**. `gogl-x` and
  `gogl-x.lab.example` both answered `192.168.2.251`, so the suffix need not match dnsmasq's
  own `domain` setting — host-file entries are literal name-to-address mappings.
- Propagation takes a few seconds. Immediately after a write, expect a brief window where the
  name does not resolve yet.

Two cautions:

- **`set_host` validates the content, and rejects `(`, `)` and `=`.** Anywhere in the file, on
  any line, comment or not. It answers `-32602 Invalid params` and does not say which character
  offended it. Isolated one character at a time on 2026-07-29:

  | Content | Result |
  |---|---|
  | byte-identical round trip | accepted |
  | `192.168.2.253 goglprobe` | accepted |
  | `# a plain comment` | accepted |
  | `# a comment with a colon: here` | accepted |
  | `# a comment with a dot example.test` | accepted |
  | a 78-character comment, no punctuation | accepted |
  | a trailing blank line | accepted |
  | `# gogl probe parens (here)` | **rejected** |
  | `# gogl probe equals=here` | **rejected** |

  So comments are fine, and so are colons, dots, hyphens and long lines. Length is not a factor.
  Only those three characters are known to be refused; the full accepted set is not mapped, and
  nothing else was needed. `gogl` checks content against this rule before writing
  (`types.ValidateContent`) so the caller gets the offending line rather than `Invalid params`,
  and `src/mock` enforces the same rule so the test suite can catch a violation.

  This cost real time. gogl's managed block originally carried its domain as
  `# BEGIN gogl managed hosts (domain: lab.example)`, which made **every** host-file write fail
  on hardware -- `--domain`, `--set`, `--add`, `--del`, `--clear`. The mock accepted it, so all
  404 tests passed. The marker is now `# BEGIN gogl managed hosts domain lab.example`.

- **`set_host` replaces the entire file.** It ships loopback and IPv6 boilerplate that the
  router itself relies on; read, modify, and write back. Do not send only your own entries.
- **Watch trailing newlines.** Appending to content whose final newline was stripped — easy to
  do in a shell, where `$(...)` eats it — concatenates your entry onto the previous line. That
  turned `ff02::2 ip6-allrouters` into the answer for a hostname during this investigation.

### There is no way to set the dnsmasq domain

Searched all 313 documented methods: no endpoint exposes a `domain`, `local`, or
`expandhosts` setting. `dns.set_config` covers mode, servers, `force_dns` and
`rebind_protection` only. `custom_dns` does not exist on this model.

So a client wanting a domain must either accept dnsmasq's default (`lan`) or write
fully-qualified names into the host file, which works for any suffix.

### Writes restart services

Each static-bind write appears to reload the router's DHCP/DNS services. Observed effects:

- A brief DNS outage immediately after a write — `connection refused` on port 53 for a few
  seconds.
- On one occasion the LAN link dropped long enough that the managing host failed over to
  another network and could not reach the router again.

Batch your writes, expect a short interruption, and do not assume the session survives one.

### Wireless writes

`wifi.get_config` is verified: it returns one entry per radio under `res`, each carrying
`band`, `device`, `channel`, `htmode`, `hwmode`, `txpower` and an `ifaces` array. Each interface
has `name` (the `iface_name` a write requires), `ssid`, `key`, `encryption`, `enabled`, `hidden`
and `guest`.

**Passphrases come back in cleartext**, here and in `system.get_status`, over plain HTTP on port
80. Reading them requires LAN access and the admin password, which is the same bar as opening
the admin panel, so this is recorded rather than treated as a defect. `gogl` masks keys in
`gogl lan show` and `gogl wifi show` output unless `--show-key` is given, which keeps them out
of terminal scrollback and pasted bug reports; it is not an access control.

Each radio also advertises what it supports. **Three of these fields disagree with GL.iNet's own
description, captured 2026-07-29:**

| Field | Description says | Device sends |
|---|---|---|
| `htmodes` | array of bandwidth strings | **object** keyed by hardware mode, values are the max channel width in MHz, plus an `auto` key holding a **bool** |
| `hwmodes` | array, implying `11b`/`11g`/`11n` | array of slash-joined combinations: `["11n","11g/n","11b/g/n"]`, `["11ac","11n/ac","11a/n/ac"]` |
| `encryptions` | unspecified | `["none","psk2","psk-mixed","sae","sae-mixed"]` -- includes WPA3, and no bare `psk` |

Verbatim, from the 2.4GHz radio:

```json
"htmode": "auto",
"htmodes": {"11b/g/n": 40, "11g/n": 40, "11n": 40, "auto": true},
"hwmode": "11g/n",
"hwmodes": ["11n", "11g/n", "11b/g/n"]
```

and the 5GHz radio:

```json
"htmode": "20",
"htmodes": {"11a/n/ac": 80, "11ac": 80, "11n/ac": 80, "auto": false},
"hwmode": "11a/n/ac",
"hwmodes": ["11ac", "11n/ac", "11a/n/ac"]
```

Note the current `htmode` is a **width string** -- `auto` or `20` were observed -- not the
`HT20`/`VHT80` form the description implies. Since the firmware reports only the maximum width
per hardware mode, `gogl` infers that narrower widths are also settable, and lists `auto`, `20`,
`40`, `80` accordingly. That inference is unconfirmed by a successful write.

#### `channels` reports firmware policy, not hardware capability

`channels` is an array of objects carrying `channel` and a `dfs` flag, with a top-level
`dfs_support`. **Both describe what the firmware will do, not what the radio can do**, and the
difference is large.

The captured unit reports `dfs_support: false` and lists nine 5GHz channels:

```
36, 40, 44, 48, 149, 153, 157, 161, 165
```

The same device's driver, asked directly over SSH with `iw phy`, reports **twenty-five**, and the
missing sixteen are precisely the DFS ones:

```
5180 [36]  5200 [40]  5220 [44]  5240 [48]              <- offered by the API
5260 [52] (radar detection)   5280 [56] (radar detection)
5300 [60] (radar detection)   5320 [64] (radar detection)
5500 [100] ... 5720 [144]     (all radar detection)     <- absent from the API
5745 [149] 5765 [153] 5785 [157] 5805 [161] 5825 [165]  <- offered by the API
```

`iw phy` also lists radar detect widths in its interface combinations, so the hardware and the
`mac80211` driver both support DFS. GL.iNet's API simply does not offer those channels.

Consequences for a client:

- `wifi.set_config` presumably will not accept a DFS channel either, so validating against the
  API's list is the right behavior. But the refusal means "this firmware will not", not "this
  radio cannot", and an error message that says the latter is wrong.
- `dfs_support` cannot be used to decide whether DFS handling matters on a given model. A `false`
  here is a policy statement.
- gogl's DFS warning path is therefore unreachable against this firmware and is exercised only
  against a synthetic fixture (`mock.DFSWireless`). It is kept because other models and other
  regulatory domains do expose DFS channels, and vacating one takes every client with it.

An earlier version of this document said the unit "lists no DFS channel", which was true of the
API and false of the device. The lesson is the same one `htmodes` taught: what the API says about
the hardware is a claim by the API.

`gogl` validates writes against these lists rather than hardcoding values, since they differ per
radio and per regulatory domain. Typing `htmodes` from the description instead of a capture made
every read of `wifi.get_config` fail with a decode error -- the second time in this project that
trusting the vendored description cost more than capturing would have.

`wifi.set_config` is **VERIFIED**, 2026-07-29: `{"iface_name": "default_radio0", "ssid": "..."}`
and the same for `default_radio1` both applied, on a session arriving over ethernet. The write
takes effect immediately -- **no apply or commit step** -- matching static binds and host-file
writes. It is scoped two ways:

| Scope | Key | Fields |
|---|---|---|
| Interface | `iface_name` | `ssid`, `key`, `encryption`, `enabled`, `hidden` |
| Radio | `device` | `channel`, `htmode`, `hwmode`, `txpower` |

`gogl` sends the two scopes as two calls, and sends only the fields the caller named.

**Partial updates leave unmentioned fields alone. VERIFIED** the same day: after writing only
`ssid` to both main interfaces, `wifi.get_config` still reported the original 8-character
passphrase and `psk2` encryption on each, with the guest interfaces untouched. This is the
assumption the whole partial-update design rests on, and it holds.

Still unverified:

1. **The radio-scoped fields.** `channel`, `htmode`, `hwmode` and `txpower` have not been written.
   The inference that narrower channel widths are settable when only the maximum is reported is
   also untested.
3. **Whether both scopes travel in one call.** `gogl` sends two calls and does not rely on it.

### The hardware exceeds what the API exposes

Confirmed by SSH to the same unit, which runs **OpenWrt 18.06 / LEDE** with a real
`mac80211`/`nl80211` stack -- `iw` works fully and reports complete HT/VHT capability tables.
Two places where the API understates the device:

- **DFS channels**, described above.
- **Monitor mode, sort of.** Both radios list `monitor` in their supported interface modes.
  Adding one fails with `-22 EINVAL` while the APs are up -- consistent with
  `valid interface combinations` listing only `{managed, AP}` and no monitor. After `wifi down`
  the interface *is* created and reports `type monitor`, and then **bringing it up wedged the
  device**, requiring a power cycle. Tested once, 2026-07-29.

  So the practical answer is that monitor mode is advertised, attachable, and not usable. The
  driver declares `NL80211_IFTYPE_MONITOR` without a working implementation behind it. Packet
  injection was never reached. Use a USB adapter with `ath9k_htc` or `mt7612u` instead.

  None of this is reachable over the JSON-RPC API and all of it needs SSH, which `gogl` excludes
  by rule -- which is the useful part of the finding: an operation that hangs the device is
  exactly what belongs outside a tool meant for unattended bulk provisioning.

The device is also capable of an upstream WiFi uplink: `iw dev` showed a `wlan-sta0` interface in
`managed` mode associated to another SSID, alongside the AP on the same radio. `gogl` does not
model repeater mode.

### The ubus / LuCI path, and why it is unavailable here

**Recorded as a closed door, not a route forward.** `gogl` targets GL.iNet 4.x only and will not
grow a generic OpenWrt backend; see [Scope](VISION.md#scope). This section stays because it
answers a question that will keep coming up — "the device runs OpenWrt, why not just use UCI?" —
and because it is the evidence behind the domain-in-a-marker-line design.

OpenWrt exposes UCI over HTTP through `uhttpd` at `POST /ubus`, JSON-RPC 2.0, brokered by
`rpcd` with ACLs from `/usr/share/rpcd/acl.d/*.json`. That path *can* set the dnsmasq domain,
which the GL.iNet API cannot:

```sh
SID=$(curl -s -d '{"jsonrpc":"2.0","id":1,"method":"call","params":[
  "00000000000000000000000000000000","session","login",
  {"username":"root","password":"PASS"}]}' http://ROUTER/ubus \
  | jq -r '.result[1].ubus_rpc_session')

# dnsmasq is an anonymous section, so resolve its name first
curl -s -d '{"jsonrpc":"2.0","id":1,"method":"call","params":["'$SID'","uci","get",
  {"config":"dhcp","type":"dnsmasq"}]}' http://ROUTER/ubus

curl -s -d '{"jsonrpc":"2.0","id":1,"method":"call","params":["'$SID'","uci","set",
  {"config":"dhcp","section":"cfg01411c",
   "values":{"domain":"home.arpa","local":"/home.arpa/","expandhosts":"1"}}]}' http://ROUTER/ubus

curl -s -d '{"jsonrpc":"2.0","id":1,"method":"call","params":["'$SID'","uci","commit",
  {"config":"dhcp"}]}' http://ROUTER/ubus
```

`uci` methods: `configs, get, set, add, delete, rename, order, changes, revert, commit, apply,
confirm, rollback, reload_config`. LuCI's "Save & Apply" is `uci.apply {"rollback":true,
"timeout":90}` followed by `uci.confirm` — the rollback-if-connection-lost dance. Service
restart is `luci.setInitAction {"name":"dnsmasq","action":"restart"}`, from `rpcd-mod-luci`.
List options take an array value. Locally the JSON-RPC frame is unnecessary:
`ubus call uci get '{"config":"dhcp","type":"dnsmasq"}'`.

**On this device `/ubus` returns 404.** nginx 1.17.7 fronts the admin interface rather than
uhttpd, so there is no `ubus_prefix` to route. Confirmed by request; both `POST /ubus` and a
`session.login` attempt return nginx's 404 page.

Enabling it would mean `uci set uhttpd.main.ubus_prefix='/ubus'` plus a uhttpd restart, over
SSH — and with nginx serving the panel, that would not obviously take effect anyway. A
non-root service account additionally needs write ACL on `uci: ["dhcp"]`, or `set` returns
status 6.

For older firmware: LuCI 18.06 and earlier have no JSON API by default; the DHCP/DNS page is a
server-rendered CBI form posting `cbid.dhcp.<sid>.domain=...` with a CSRF token. With
`luci-mod-rpc` installed there is JSON-RPC 1.0 at `/cgi-bin/luci/rpc/uci`, authenticated via
`/cgi-bin/luci/rpc/auth`. Firmware 3.x uses a REST-ish `/cgi-bin/api/<module>/<action>`.
`cat /etc/glversion` identifies the generation.

### Endpoints documented but absent on this model

`system.get_config`, `lan.get_config`, `network.get_config`, `network.get_status`,
`firewall.get_config`, `firewall.get_status` all return `-32601 Method not found`. The
description covers the whole firmware line; the SFT1200 is a reduced build.

These six came from GL.iNet's original public reference and do **not** appear in the API
description that [`docs/api/`](docs/api/README.md) is generated from, so they have no entry
there to mark. Three further absences *are* in that description and are marked in it:
`custom_dns.get_info`, `custom_dns.set_info` and `modem.get_config`. Nine known absences in
total.

### Not reachable at all

No introspection method exists — `list`, `methods`, `help`, `api`, `describe` are all
`-32601` at the top level. `uci`, `file`, `service`, and `luci-rpc` are not exposed. Dotted
ubus object names such as `network.interface.lan` return `-32602`, and the admin panel's own
JavaScript bundle 404s at the path its `index.html` references, so the group names cannot be
mined from the UI.

---

## Error codes

| Code | Meaning | Notes |
|------|---------|-------|
| `-32000` | Access denied | Stale/absent `sid`, **or** a wrong login digest. The overloading is why a wrong `hash-method` is so hard to diagnose. |
| `-32003` | Login fail number over limit | Brute-force lockout; carries `data.wait` in seconds |
| `-32601` | Method not found | The group exists but not the method, or neither. Indistinguishable between the two. |
| `-32602` | Invalid params | The endpoint **exists** but needs arguments you did not supply. Proven with `iwinfo.info`, which fails bare and succeeds given `{"device":"wlan0"}`. Useful for telling "exists, needs args" from "does not exist". |
| `-32603` | Internal error | Seen from `modem.get_config` on a model with no modem |

---

## Capturing more

`discovery/probe.go` in this repository does the login and probes endpoints:

```bash
export GL_PASSWORD='...'
go run ./discovery -H 192.168.8.1                          # diagnose the login
go run ./discovery -H 192.168.8.1 -batch candidates.txt     # probe many, one login
go run ./discovery -H 192.168.8.1 -group system -method get_status
```

It honors `hash-method` by default. `-variants` sweeps digest alternatives instead, but
**each wrong guess counts against the lockout** — use it only against a device you can afford
to be locked out of for ten minutes.
