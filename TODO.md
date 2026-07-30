# TODO: capture pass, 2026-07-31

Everything in [`docs/DESIGN-V2.md`](docs/DESIGN-V2.md) phases 3–5 waits on this. Phases 1 and 2
do not — those are wiring over endpoints already verified, and can start any time.

Work top to bottom. **Step 1 is decisive**: it determines whether the two-HTTP-backend design
survives, so do it first and stop if it fails.

Record everything into `discovery/captures/` (gitignored — payloads contain passphrases). Paste
results back into the conversation and I will turn them into types, fixtures and docs.

```bash
mkdir -p discovery/captures
```

---

## 1. Is `/ubus` routed on the OpenWrt One? (decisive, ~2 minutes)

On the device:

```bash
uci get uhttpd.main.ubus_prefix      # expect: /ubus
ls /usr/share/rpcd/acl.d/            # expect at least unauthenticated.json
opkg list-installed | grep -E 'rpcd|uhttpd|luci-base'
```

From the laptop — this is the one that matters:

```bash
export OW=192.168.1.1
export OWPASS='...'                  # leading space if HISTCONTROL=ignorespace

curl -s -X POST http://$OW/ubus -d '{"jsonrpc":"2.0","id":1,"method":"call","params":
  ["00000000000000000000000000000000","session","login",
  {"username":"root","password":"'"$OWPASS"'"}]}' | tee discovery/captures/ubus-login.json
```

**Expected:** JSON containing `ubus_rpc_session`, a 32-hex-character string.

| Result | Meaning |
|---|---|
| `ubus_rpc_session` present | Two-HTTP-backend design confirmed. Continue to step 2. |
| HTTP 404 | `/ubus` not routed. **Stop and tell me** — we revisit the SSH question. |
| `{"error":...,"code":-32002}` | Access denied: wrong password, or ACL missing. Check `acl.d/`. |

If it worked, keep the session id for the rest:

```bash
export SID=$(jq -r '.result[1].ubus_rpc_session' discovery/captures/ubus-login.json)
echo "$SID"                          # sanity check: 32 hex chars
```

Sessions expire (300s default). Re-run the login if calls start returning `-32002`.

---

## 2. Capture the four UCI configs gogl needs

One call each. These are the exact shapes the OpenWrt backend must read and write.

```bash
ucall() {  # $1 = config name
  curl -s -X POST http://$OW/ubus -d '{"jsonrpc":"2.0","id":1,"method":"call","params":
    ["'"$SID"'","uci","get",{"config":"'"$1"'"}]}' \
    | tee discovery/captures/uci-$1.json | jq .
}

ucall network      # lan interface: proto, ipaddr, netmask
ucall dhcp         # pool start/limit/leasetime, @dnsmasq[0] domain, host sections
ucall wireless     # wifi-device (radio) + wifi-iface (SSID) sections
ucall system       # hostname, timezone
```

Also the board identity, which is how `gogl system info` works here:

```bash
curl -s -X POST http://$OW/ubus -d '{"jsonrpc":"2.0","id":1,"method":"call","params":
  ["'"$SID"'","system","board",{}]}' | tee discovery/captures/system-board.json | jq .
```

**What I am looking for in these:**

- `dhcp`: does the LAN section use `start` + `limit` (a count) as expected? What is `leasetime`
  spelled as? Is there a `@dnsmasq[0]` with `domain`, `local`, `expandhosts`?
- `wireless`: confirm the `wifi-device` / `wifi-iface` split, and what `htmode` and `hwmode`
  actually contain on current OpenWrt (expect `HT20`/`VHT80`-style, unlike GL.iNet's bare widths)
- `network`: is the LAN a bridge device? `br-lan` versus a plain interface changes what "the LAN"
  means
- Section names: OpenWrt uses anonymous sections with generated ids like `cfg031234`. Writes need
  either those ids or `@type[index]` addressing, and which one is stable matters

---

## 3. Confirm a write works, and that it is reversible

Nothing destructive. Read the domain, set it, read it back, restore it.

```bash
# Read
curl -s -X POST http://$OW/ubus -d '{"jsonrpc":"2.0","id":1,"method":"call","params":
  ["'"$SID"'","uci","get",{"config":"dhcp","type":"dnsmasq"}]}' | jq .

# Set (note: no shell interpolation of user data; the value is a literal here)
curl -s -X POST http://$OW/ubus -d '{"jsonrpc":"2.0","id":1,"method":"call","params":
  ["'"$SID"'","uci","set",{"config":"dhcp","type":"dnsmasq","values":{"domain":"probe.test"}}]}' | jq .

# What is staged?
curl -s -X POST http://$OW/ubus -d '{"jsonrpc":"2.0","id":1,"method":"call","params":
  ["'"$SID"'","uci","changes",{}]}' | jq .

# Discard it — do NOT commit
curl -s -X POST http://$OW/ubus -d '{"jsonrpc":"2.0","id":1,"method":"call","params":
  ["'"$SID"'","uci","revert",{"config":"dhcp"}]}' | jq .

# Prove it is back
curl -s -X POST http://$OW/ubus -d '{"jsonrpc":"2.0","id":1,"method":"call","params":
  ["'"$SID"'","uci","changes",{}]}' | jq .    # expect empty
```

**Record whether `revert` cleanly discarded the change.** That answers whether the OpenWrt backend
can offer a real dry-run — stage, show `changes`, revert — which the GL.iNet backend cannot,
since its writes are immediate.

If you want to see a commit work, do it with the domain you actually want:

```bash
# set domain, then:
#   uci commit {"config":"dhcp"}
#   uci apply  {"rollback":true,"timeout":90}   <- LuCI's Save & Apply
#   uci confirm {}                              <- within the timeout, or it rolls back
```

`apply` + `confirm` is the rollback-on-lost-connection dance. **Note whether it exists here** — if
it does, the OpenWrt backend gets a safety property the GL.iNet one lacks, and the LAN-renumber
guard could relax on OpenWrt.

---

## 4. Back on the SFT1200: capture `wan` and `access`

These are the two areas in `DESIGN-V2.md` with nothing verified — 2 methods of 40 for WAN, 0 of 9
for access. Use `discovery/shape`, which already exists and redacts secrets.

These are the real method names, read out of `docs/api/reference/`, not guesses. All read-only.

```bash
cd ~/h/src/er/gogl
export GL_ROUTER_IP=192.168.8.1      # or wherever it is now
export GL_PASSWORD='...'

probe() { echo "=== $1.$2 ==="; go run ./discovery/shape -group "$1" -method "$2" 2>&1; }

{
  probe netmode   get_mode           # which WAN mode is active: the entry point
  probe cable     get_status         # wired WAN
  probe cable     get_config
  probe network   check_wan_cable    # is a cable plugged in
  probe network   routes
  probe network   get_arp_list
  probe repeater  get_status         # WiFi uplink — this router uses one
  probe repeater  get_config
  probe repeater  get_saved_ap_list
  probe tethering get_status
  probe modem     get_status         # expect -32603 or absent: no modem on this model
} | tee discovery/captures/sft1200-wan.txt
```

`-32601 Method not found` and `-32603 Internal error` are both useful answers — they tell us what
is absent on this model, which the reference cannot.

`repeater get_status` is the interesting one: `iw dev` showed a `wlan-sta0` interface associated to
`Autograph_GUEST`, so this router's uplink *is* a WiFi repeater connection. That is the live case.

**Do not run `repeater connect`, `disconnect`, `forget`, or anything under `set_`.** Several would
drop your uplink.

Then access:

```bash
{
  probe acl    get_group_list
  probe acl    get_acl_list
  probe system get_status
} | tee discovery/captures/sft1200-acl.txt
```

**Do not probe `system.set_password`.** It takes `old_password` and `new_password`, and a
half-understood call to it locks you out of your own router.

---

## 5. Optional: does `reboot` work as documented?

`gogl system reboot` is in the tree and uncaptured. The `reboot` group has 2 methods. Only try
this when you are physically near the device.

```bash
grep -A 12 '^| \[`' docs/api/reference/reboot.md
```

---

## What to bring back

Paste, or say the file is in `discovery/captures/`:

1. **Step 1's verdict** — `/ubus` routed, yes or no. Everything else is contingent on it.
2. `uci-dhcp.json`, `uci-wireless.json`, `uci-network.json` — these become the OpenWrt types
3. Whether `uci revert` and `uci apply`/`confirm` work
4. The WAN capture from the SFT1200
5. `system-board.json`

With 1–3 I can write `backend/openwrt` and move the three wire artifacts out of `src/types`. With
4 I can write `gogl wan` from evidence instead of from a vendor description that has been wrong
three times.

---

## Meanwhile, needing no hardware

Phases 1 and 2 of `DESIGN-V2.md` are unblocked: the cobra tree over verified endpoints, TOML
config, XDG paths, `password_command`, and the install move to `~/.local/bin`. Say the word and I
will start on those today rather than waiting for tomorrow's results.

Also outstanding and independent of all this:

- `gogl` (as `goglnet`) puts a WiFi passphrase in argv with `--set-key`, contradicting the reason
  `GL_PASSWORD` has no flag. Fix: prompt without echo, or `--passphrase-command`.
- `--version` with a build stamp. Two stale-binary confusions this week would have been one
  command each.
- `make check-docs`: diff every documented flag against `--help`, and compile the README's Go
  snippets. A stale `Network().Set` call sat in the README through several changes.
