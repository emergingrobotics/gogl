# `ddns`

This is the API related to ddns.

3 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","ddns","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | verified | Get the ddns cloud control status of the router. |
| [`get_status`](#get_status) | - | get ddns status. |
| [`set_config`](#set_config) | - | Set the ddns cloud control status of the router. |

---

## get_config

Get the ddns cloud control status of the router.

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enable_ddns` | bool | - | Whether to enable the DDNS function. |
| `enable_ssh_access` | bool | - | Whether to allow access through SSH. |
| `enable_http_access` | bool | - | Whether to allow access through HTTP. |
| `enable_https_access` | bool | - | Whether to allow access through HTTPS. |
| `device_id` | string | - | ddns device id. |
| `err_code` | number | - | Error code,-2:uci init error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ddns\",\"get_config\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"device_id\":\"aa628d6\",\"enable_ddns\":false,\"enable_ssh_access\":false,\"enable_http_access\":false,\"enable_https_access\":false}}
```

---

## get_status

get ddns status.

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `ip_public` | string | - | public ip. |
| `ip_nslookup` | string | - | nslookup ip. |
| `ip_wan` | string | - | wan ip. |
| `status` | number | - | ddns status,0:DDNS resolution service is normal,1:DDNS resolution failed,2:Behind the firewall or the IP provided by the operator is not a public IP. |
| `err_code` | number | - | Error code,-2:uci init error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ddns\",\"get_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"ip_public\":\"61.141.113.215\",\"ip_nslookup\":\"No nslookup ip\",\"ip_wan\":\"192.168.113.137\",\"status\":1}}
```

---

## set_config

Set the ddns cloud control status of the router.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `enable_ddns` | bool | yes | Whether to enable the DDNS function. |
| `enable_ssh_access` | bool | yes | Whether to allow access through SSH. |
| `enable_http_access` | bool | yes | Whether to allow access through HTTP. |
| `enable_https_access` | bool | yes | Whether to allow access through HTTPS. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1:missing parameters,-2:uci init error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ddns\",\"set_config\",{\"enable_ddns\":true,\"enable_ssh_access\":true,\"enable_http_access\":true,\"enable_https_access\":true}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---
