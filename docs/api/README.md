# GL.iNet Firmware 4.x JSON-RPC API Reference

43 groups, 313 methods.

## Provenance

GL.iNet's official 4.x API reference at `dev.gl-inet.com/router-4.x-api` is no
longer publicly reachable. This reference was generated from the
machine-readable API description that ships inside
[`python-glinet`](https://github.com/tomtana/python-glinet), which exports
GL.iNet's own documentation database.

That source file is **not** vendored here. `python-glinet` is GPL-3.0 and gogl is
MIT, so redistributing it would be a licensing question not worth guessing at.
Method names and signatures are functional interface facts and are documented
freely below; to obtain the source yourself:

```bash
scripts/fetch-api-description.sh
python3 scripts/generate-api-docs.py /tmp/gl-api-description.json
```

## What is actually verified

The description is GL.iNet's documentation for the firmware line as a whole, not
for any one device. 22 endpoints were confirmed by calling them
against a **GL-SFT1200 (Opal) running firmware 4.3.28**, and several documented
endpoints are absent on it. Each method below is marked accordingly.

Two payloads differ from the documentation on real hardware:

- `network.get_dhcp_leases` wraps its array in `leases`, not the documented `entries`.
- `lan.get_config_list` returns more fields than documented (`leasetime`, `dns`,
  `gateway`, `lpr`) and reports `enable` as a **number**, not a boolean.

Treat this reference as a map, and the device as the authority. See
[`../../GL_INET_4X_API_DOCUMENTATION.md`](../../GL_INET_4X_API_DOCUMENTATION.md)
for authentication, error codes, and the hardware-verified essentials.

## Groups

| Group | Methods | Verified here | Description |
|---|---|---|---|
| [`acl`](reference/acl.md) | 8 | - | Permission Management |
| [`adguardhome`](reference/adguardhome.md) | 2 | - | Adguardhome |
| [`cable`](reference/cable.md) | 4 | - | This is the API related to wired Internet. |
| [`clients`](reference/clients.md) | 5 | 2 | Client management related interfaces |
| [`cloud`](reference/cloud.md) | 3 | 1 | This is the API related to file share. |
| [`cloud_batch_manage`](reference/cloud_batch_manage.md) | 7 | - | This is the API related to cloud batch manage. |
| [`custom_dns`](reference/custom_dns.md) | 2 | - | Custom dns server api |
| [`ddns`](reference/ddns.md) | 3 | 1 | This is the API related to ddns. |
| [`diag`](reference/diag.md) | 2 | - | Diag |
| [`dlna`](reference/dlna.md) | 2 | - | This is the API related to file share. |
| [`dns`](reference/dns.md) | 6 | 3 | DNS interface |
| [`edgerouter`](reference/edgerouter.md) | 5 | - | Bypass Router Mode |
| [`fan`](reference/fan.md) | 2 | - | This is the API related to fan Internet access. |
| [`firewall`](reference/firewall.md) | 11 | - | Firewall |
| [`igmp`](reference/igmp.md) | 2 | - | This is the igmp api. |
| [`ipv6`](reference/ipv6.md) | 2 | - | IPv6 settings |
| [`lan`](reference/lan.md) | 6 | 6 | This is the API related to wired Internet. |
| [`led`](reference/led.md) | 2 | 1 | This is the API related to led Internet access. |
| [`logread`](reference/logread.md) | 8 | - | logread |
| [`macclone`](reference/macclone.md) | 2 | - | macclone settings |
| [`modem`](reference/modem.md) | 17 | - | This is the API related to modems. |
| [`nas_web`](reference/nas_web.md) | 17 | - | This is the API related to nas web config. |
| [`netmode`](reference/netmode.md) | 2 | - | Network Mode Configuration |
| [`network`](reference/network.md) | 5 | 1 | Network related operations |
| [`ovpn_client`](reference/ovpn_client.md) | 23 | - | OpenVPN Client API |
| [`ovpn_server`](reference/ovpn_server.md) | 16 | - | ovpn server api |
| [`plugins`](reference/plugins.md) | 6 | - | This is the api related to installing the package. |
| [`qos`](reference/qos.md) | 15 | 1 | QOS management related interfaces |
| [`reboot`](reference/reboot.md) | 2 | - | This is the API related to reboot Internet access. |
| [`repeater`](reference/repeater.md) | 9 | 1 | WiFi Repeater |
| [`rs485`](reference/rs485.md) | 13 | - | This is the api related to rs485. |
| [`rtty`](reference/rtty.md) | 4 | - | Rtty Configuration |
| [`s2s`](reference/s2s.md) | 7 | - | This is the API related to s2s. |
| [`samba`](reference/samba.md) | 2 | - | This is the API related to samba. |
| [`switch_button`](reference/switch_button.md) | 3 | - | Switch Button Configuration |
| [`system`](reference/system.md) | 16 | 2 | System Operations |
| [`tethering`](reference/tethering.md) | 3 | - | This is the API related to tethering Internet access. |
| [`ui`](reference/ui.md) | 6 | - | UI Operations |
| [`upgrade`](reference/upgrade.md) | 7 | 1 | This is the api related to firmware upgrade. |
| [`vpn_policy`](reference/vpn_policy.md) | 10 | - | VPN Policy |
| [`wg_client`](reference/wg_client.md) | 24 | - | Wireguard Client API |
| [`wg_server`](reference/wg_server.md) | 18 | - | Wireguard Server API |
| [`wifi`](reference/wifi.md) | 4 | 2 | WiFi related operations |
