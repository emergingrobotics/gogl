# `igmp`

This is the igmp api.

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","igmp","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get`](#get) | - | Get IGMP configuration |
| [`set`](#set) | - | Set IGMP |

---

## get

Get IGMP configuration

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enable` | bool | - | Whether to enable IGMP |
| `version` | number | - | IGMP version [1/2/3]. |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "igmp",
    "get",
    {}
  ],
  "id": 1
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "enable": false,
    "version": 3
  }
}
```

---

## set

Set IGMP

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `enable` | bool | yes | Whether to enable IGMP. |
| `version` | number | yes | IGMP version [1/2/3]. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1: Parameter error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "igmp",
    "set",
    {
      "version": 3,
      "enable": false
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
  "result": {}
}
```

---
