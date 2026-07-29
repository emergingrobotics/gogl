# `samba`

This is the API related to samba.

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","samba","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | - | get samba conf |
| [`set_config`](#set_config) | - | set samba conf |

---

## get_config

get samba conf

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `share_path` | string | - | samba share dir path |
| `enable` | bool | - | samba share enable, share in lan |
| `wan_share` | bool | - | samba share in wan |
| `writeable` | bool | - | samba share is writeable |
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"samba\",\"get_config\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"share_dir\": "/mnt",\"share_on_lan\":"true",\"share_on_wan\":\"true\",\"writable\":\"true\",\"list\":[\"/mnt\",\"/mnt/sda1\",\"/mnt/sdb1\"]}}
```

---

## set_config

set samba conf

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `share_path` | string | yes | samba share path |
| `enable` | bool | yes | samba share enable, share in lan |
| `wan_share` | bool | yes | samba share in wan |
| `writeable` | bool | yes | samba share is writeable |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"samba\",\"set_config\",{\"path\":\"/mnt\",\"lan_share\":\"true\",\"wan_share\":\"true\",\"writable\":\"true\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---
