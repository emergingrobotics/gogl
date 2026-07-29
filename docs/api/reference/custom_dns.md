# `custom_dns`

Custom dns server api

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","custom_dns","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_info`](#get_info) | absent | Get dns setting info |
| [`set_info`](#set_info) | absent | Set custom dns setting |

---

## get_info

Get dns setting info

**Absent** on a GL-SFT1200 running 4.3.28: returns `-32601 Method not found`.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `manual_dns` | bool | - | Identifies whether setting manual dns |
| `custom_dns` | string | - | Specifies custom dns |
| `force_dns` | bool | - | Identifies whether enable force dns |
| `auto_dns` | bool | - | Identifies whether enable auto dns |
| `cloudflare_dns` | bool | - | Identifies whether enable cloudflare dns |
| `dns_name` | string | - | Cloudfare or Nextdns |
| `nextdns_id` | string | - | Nextdns id |
| `dnscrypt_proxy` | bool | - | Identifies whether enable dnscrypt_proxy dns |
| `proxy_server` | string | - | Proxy server |
| `proxy_serverlist` | array | - | Proxy server list |
| `rebind_protection` | bool | - | Identifies whether enable rebind_protection |
| `err_code` | number | - | ERR CODE |
| `err_msg` | string | - | ERR MSG |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"custom-dns\",\"get_info\",{}],\"id\":1 }
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": {\"manual_dns\": false,\"custom_dns\":\"\",\"force_dns\": false,\"auto_dns\": true,\"cloudflare_dns\": false,\"dns_name\":\"Cloudflare\",\"nextdns_id\":\"\",\"dnscrypt_proxy\": false,\"proxy_server\":\"\",\"proxy_serverlist\": [],\"quad9_dns\": false,\"rebind_protection\": true}}
```

---

## set_info

Set custom dns setting

**Absent** on a GL-SFT1200 running 4.3.28: returns `-32601 Method not found`.

| Params | Type | Required | Description |
|---|---|---|---|
| `dns1` | string | no | Manual DNS setting First dns (dns1 is present only when manual_dns is true) |
| `dns2` | string | no | Manual DNS setting Second dns (dns2 is present only when manual_dns is true) |
| `dns_name` | string | no | Cloudfare or NextDNS (dns_name is present only when cloudflare_dns is true) |
| `nextdns_id` | string | no | Next DNS ID (nextdns_id is present only when dns_name is NextDNS) |
| `proxy_server` | string | no | Proxy Server (proxy_server is present and required only when dnscrypt_proxy is true) |
| `force_dns` | bool | no | Identifies whether enable force dns |
| `auto_dns` | bool | no | Identifies whether enable auto dns (when auto_dns is true, manual_dns, cloudflare_dns, dnscrypt_proxy, quad9_dns cannot be true) |
| `manual_dns` | bool | no | Identifies whether setting manual dns (when manual_dns is true, auto_dns, cloudflare_dns, dnscrypt_proxy, quad9_dns cannot be true) |
| `cloudflare_dns` | bool | no | Identifies whether enable cloudflare dns (when cloudflare_dns is true, auto_dns, manual_dns, dnscrypt_proxy, quad9_dns cannot be true) |
| `dnscrypt_proxy` | bool | no | Identifies whether enable dnscrypt_proxy dns (when dnscrypt_proxy is true, auto_dns, manual_dns, cloudflare_dns, quad9_dns cannot be true) |
| `quad9_dns` | bool | no | Identifies whether enable quad9 dns (when quad9_dns is true, auto_dns, manual_dns, cloudflare_dns, dnscrypt_proxy cannot be true) |
| `rebind_protection` | bool | no | Identifies whether enable rebind_protection |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | ERR CODE; -1: Parameter missing, -2: Parameter error, -3: DNS is occupied |
| `err_msg` | string | - | ERR MSG |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"custom-dns\",\"set_info\",{\"cloudflare_dns\":false}],\"id\":1 }
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": {}}
```

---
