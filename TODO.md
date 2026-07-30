# TODO: capture pass on the SFT1200

Phase 4 of [`docs/DESIGN-V2.md`](docs/DESIGN-V2.md) waits on this. The rest of phase 1 and all of
phase 2 do not — those are wiring over endpoints already verified.

**The OpenWrt One steps are gone.** Generic OpenWrt support was dropped: `gogl` is a GL.iNet 4.x
tool and will stay one. See [Scope](VISION.md#scope). Nothing here needs a second device.

Everything below is read-only except where marked. Record into `discovery/captures/` — gitignored,
because these payloads carry WiFi passphrases in cleartext.

```bash
mkdir -p discovery/captures
cd ~/h/src/er/gogl
export GL_ROUTER_IP=192.168.8.1      # or wherever it is now
export GL_PASSWORD='...'             # leading space if HISTCONTROL=ignorespace
```

---

## 1. `wan` — 2 verified methods of 40

The largest gap in the command tree. These method names are read out of
`docs/api/reference/`, not guessed. `discovery/shape` redacts secrets.

```bash
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
  probe modem     get_status         # expect absent: no modem on this model
} | tee discovery/captures/sft1200-wan.txt
```

`-32601 Method not found` and `-32603 Internal error` are both useful answers: they tell us what
is absent on this model, which the reference cannot.

`repeater get_status` is the interesting one. `iw dev` showed a `wlan-sta0` interface associated to
`Autograph_GUEST`, so this router's uplink *is* a repeater connection — the live case, not a
hypothetical.

**Do not run `repeater connect`, `disconnect`, `forget`, or anything under `set_`.** Several would
drop your uplink, and `wlan-sta0` shares `phy1` with your 5GHz AP.

---

## 2. `access` — 0 verified methods of 9

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

## 3. `system reboot` — uncaptured

`gogl system reboot` is in the tree with nothing behind it. The `reboot` group has 2 methods.

```bash
grep -A 12 '^| \[`' docs/api/reference/reboot.md
```

Only exercise it when you are physically near the device.

---

## 4. Optional: does a radio write work?

The last unverified assumption in the wireless surface. `wifi.set_config`'s interface-scoped
fields are confirmed — writing only `ssid` left the passphrase and encryption intact — but the
**radio-scoped** fields have never been written, and one inference rests on them: the firmware
reports only the maximum channel width per hardware mode, and `gogl` offers every narrower width
on the assumption they are valid.

Over ethernet, and expect the radio's clients to drop:

```bash
goglnet                                          # note the current channel and width
goglnet --device radio1 --set-channel 149 --dry-run
goglnet --device radio1 --set-channel 149
goglnet                                          # did it take?
goglnet --device radio1 --set-htmode 40          # tests the narrower-width inference
```

If `--set-htmode 40` is rejected, that inference is wrong and `HTModes.Options` should offer only
the reported maximum.

**These commands need the old binaries.** `make build` currently produces none — the command tree
is mid-rebuild — but your installed copies in `~/.local/bin` still work.

---

## What to bring back

1. `sft1200-wan.txt` — becomes the `wan` types and command surface
2. `sft1200-acl.txt` — becomes `access`
3. Whether a radio write works, and whether `--set-htmode 40` is accepted

With 1 and 2 I can write those areas from evidence rather than from a vendor description that has
been wrong three times: `dhcp.*` does not exist, `dns.set_host` rejects `(`, `)` and `=`, and
`htmodes` is an object rather than an array.

---

## Meanwhile, needing no hardware

**The command tree.** `utilities/gogl` — the cobra tree over `reservations`, `netcfg`, `clients`
and `profile`, which are already extracted, exported and tested. This is the next thing to build,
and it restores a working CLI.

Smaller and independent:

- `netcfg`'s `optionalBool`/`optionalInt`/`optionalString` become dead once pflag's `Changed()`
  handles set-versus-unset. Delete them and their tests along with the tree.
- Wire `--version`, `--output` and `--router` as persistent flags on the root command, replacing
  the per-tool copies.
