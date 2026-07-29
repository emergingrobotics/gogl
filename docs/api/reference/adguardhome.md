# `adguardhome`

Adguardhome

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","adguardhome","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | - | Query current configuration |
| [`set_config`](#set_config) | - | Set |

---

## get_config

Query current configuration

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enabled` | bool | - | Whether to enable |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "adguardhome",
    "get_config"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "enabled": true
  }
}
```

---

## set_config

Set

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `enabled` | bool | yes | Whether to enable |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code: 1: Other DNS not closed |
| `err_msg` | number | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "adguardhome",
    "set_config",
    {
      "enabled": true
    }
  ]
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
