# `dns`

DNS interface

6 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","dns","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | verified | Get DNS configuration |
| [`get_host`](#get_host) | verified | Get host file content |
| [`get_info`](#get_info) | - | Get DNS information |
| [`set_config`](#set_config) | - | Set DNS |
| [`set_host`](#set_host) | verified | Set host file |
| [`set_info`](#set_info) | - | Set DNS |

---

## get_config

Get DNS configuration

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `force_dns` | bool | - | Whether to override all clients |
| `rebind_protection` | bool | - | Whether to enable rebind_protection |
| `mode` | string | - | Mode: "auto", "secure", "manual", "proxy" |
| `proto` | string | - | Encryption protocol |
| `nextdns_id` | string | - | When selecting dot, if this field is present, it indicates NextDns, otherwise it indicates Cloudflare |
| `server` | array | - | Represents custom DNS server list or selected encrypted server list or proxy address or DNS server address in auto mode. For auto mode, the format is: ["wan 192.168.113.1","wwan 192.168.1.1"]. Interface names in auto mode: wan(6), wwan(6), tethering(6), ovpn(6), wg(6), modem(6) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "dns",
    "get_config",
    {}
  ]
}
```

**Response**

```json
{
  "id": null,
  "jsonrpc": "2.0",
  "result": {
    "force_dns": true,
    "proxy_server": "fvz-anyone",
    "proxy_serverlist": [
      "adguard-dns-family-ns1",
      "adguard-dns-family-ns2"
    ],
    "rebind_protection": "1"
  }
}
```

---

## get_host

Get host file content

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `content` | string | - | File content |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "dns",
    "get_host",
    {}
  ]
}
```

**Response**

```json
{
  "id": null,
  "jsonrpc": "2.0",
  "result": {
    "content": "185.199.108.154\u00a0github.githubassets.com"
  }
}
```

---

## get_info

Get DNS information

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `protos` | array | - | Supported encryption protocols: "DoT", "DoH", "DNSCrypt", "oDoH" |
| `dnscrypt_version` | number | - | DNS encryption version |
| `serverlist` | array | - | "DoH", "DNSCrypt", "oDoH" supported server list |
| `serverlist.name` | string | - | Name |
| `serverlist.proto` | string | - | Protocol |
| `serverlist.ipv6` | bool | - | IPv6 |
| `serverlist.dnssec` | bool | - | dnssec |
| `serverlist.nolog` | bool | - | No log |
| `serverlist.nofilter` | bool | - | Does not support filtering |
| `err_code` | number | - | Error code: -1 Failed to get server list |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "dns",
    "get_info",
    {}
  ]
}
```

**Response**

```json
{
  "id": null,
  "jsonrpc": "2.0",
  "result": {
    "protos": [
      "DoT",
      "DoH",
      "DNSCrypt",
      "oDoH"
    ],
    "serverlist": [
      {
        "nolog": true,
        "proto": "DoH",
        "name": "ahadns-doh-ny-ipv6",
        "nofilter": false,
        "ipv6": true,
        "dnssec": true
      },
      {
        "nolog": true,
        "proto": "DNSCrypt",
        "name": "ffmuc.net",
        "nofilter": true,
        "ipv6": false,
        "dnssec": true
      },
      {
        "nolog": true,
        "proto": "DNSCrypt",
        "name": "techsaviours.org-dnscrypt",
        "nofilter": true,
        "ipv6": false,
        "dnssec": true
      }
    ]
  }
}
```

---

## set_config

Set DNS

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mode` | string | yes | Mode: "auto", "secure", "manual", "proxy" |
| `proto` | string | no | Encryption protocol |
| `nextdns_id` | string | no | When selecting dot, if this field is present, it indicates NextDns (must be 6 hexadecimal characters or empty string), otherwise it indicates Cloudflare |
| `force_dns` | bool | yes | Whether to override all clients |
| `rebind_protection` | bool | yes | Whether to enable rebind_protection |
| `server` | array | no | Represents custom DNS server list or selected encrypted server list or proxy address |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code: -1: DNS occupied, -2: Parameter error |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "method": "call",
  "params": [
    "",
    "dns",
    "set_config",
    {
      "type": "auto",
      "rebind_protection": true
    }
  ]
}
```

**Response**

```json
{
  "id": null,
  "jsonrpc": "2.0",
  "result": []
}
```

---

## set_host

Set host file

**Verified** on a GL-SFT1200 running 4.3.28.

| Params | Type | Required | Description |
|---|---|---|---|
| `content` | string | yes | File content |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "dns",
    "set_host",
    {
      "content": "185.199.108.154\u00a0github.githubassets.com"
    }
  ]
}
```

**Response**

```json
{
  "id": null,
  "jsonrpc": "2.0",
  "result": null
}
```

---

## set_info

Set DNS

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mode` | string | yes | Mode: "auto", "secure", "manual", "proxy" |
| `proto` | string | yes | Encryption protocol |
| `nextdns_id` | string | no | When selecting dot, if this field is present, it indicates NextDns, otherwise it indicates Cloudflare |
| `force_dns` | bool | yes | Whether to override all clients |
| `rebind_protection` | bool | yes | Whether to enable rebind_protection |
| `server` | array | yes | Represents custom DNS server list or selected encrypted server list or proxy address |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code: -1: Missing parameter, -2: Parameter error, -3: DNS occupied |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "method": "call",
  "params": [
    "",
    "dns",
    "set_info",
    {
      "type": "auto",
      "rebind_protection": true
    }
  ]
}
```

**Response**

```json
{
  "id": null,
  "jsonrpc": "2.0",
  "result": []
}
```

---
