# `dlna`

This is the API related to file share.

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","dlna","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | - | get samba conf |
| [`set_config`](#set_config) | - | set dlna conf |

---

## get_config

get samba conf

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enabled` | bool | - | dlna enable status |
| `name` | string | - | dlna server name |
| `path` | string | - | dlna use path |
| `list` | array | - | dlna server can use path |
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "dlna",
    "get_config",
    {}
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "enabled": "true",
    "name": "DLNA Server",
    "path": "/mnt",
    "list": [
      "/mnt",
      "/mnt/sda1",
      "/mnt/sdb1"
    ]
  }
}
```

---

## set_config

set dlna conf

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | dlan use path |
| `enabled` | bool | yes | enable dlna |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "dlna",
    "set_config",
    {
      "enabled": "true",
      "path": "/mnt"
    }
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {}
}
```

---
