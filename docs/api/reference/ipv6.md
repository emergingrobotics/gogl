# `ipv6`

IPv6 settings

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","ipv6","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_ipv6`](#get_ipv6) | - | Get IPv6 configuration information |
| [`set_ipv6`](#set_ipv6) | - | Set IPv6 |

---

## get_ipv6

Get IPv6 configuration information

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enable` | bool | - | Enable or disable IPv6 |
| `lan_mode` | string | - | LAN interface mode is relay/nat6/static [ relay is displayed as native in the frontend ] . |
| `lan_ip` | string | - | When lan_mode is static, it is the IPv6 address of the LAN interface. |
| `lan_dns_mode` | bool | - | Mode of the LAN interface DNS server, auto is true, manual is false. |
| `lan_dns1` | string | - | When lan_dnsmode is manual, the first DNS server address of the LAN interface. |
| `lan_dns2` | string | - | When lan_dnsmode is manual, the second DNS server address of the LAN interface. |
| `err_code` | number | - | Error code,-1: Failed to obtain information. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "ipv6",
    "get_ipv6"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "enable": false,
    "lan_dns_mode": true,
    "lan_mode": "nat6"
  }
}
```

---

## set_ipv6

Set IPv6

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `lan_mode` | string | yes | LAN interface mode is relay/nat6/static [ relay is displayed as native in the frontend ]. |
| `lan_ip` | string | yes | When lan_mode is static, it is the IPv6 address of the LAN interface. |
| `lan_dns1` | string | yes | When lan_dnsmode is manual, the first DNS server address of the LAN interface. |
| `lan_dns2` | string | yes | When lan_dnsmode is manual, the second DNS server address of the LAN interface. |
| `enable` | bool | yes | Enable or disable IPv6.[After enabling IPv6, it will also enable IPv6 for Repeater, Tethering, and Modem accordingly] |
| `lan_dns_mode` | bool | yes | Mode of the LAN interface DNS server, auto is true, manual is false. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "ipv6",
    "set_ipv6",
    {
      "enable": true,
      "lan_mode": "static",
      "lan_ip": "fdab:6a57:e8f5:10:d14f:d19b:a63c:d7a5",
      "lan_dns_mode": true,
      "lan_dns1": "fdab:6a57:e8f5:10::1",
      "lan_dns2": "fdab:6a57:e8f5:10::2"
    }
  ],
  "id": 1
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": null
}
```

---
