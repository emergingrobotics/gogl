# `netmode`

Network Mode Configuration

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","netmode","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_mode`](#get_mode) | - | Query current network mode |
| [`set_mode`](#set_mode) | - | Set network mode |

---

## get_mode

Query current network mode

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `mode` | string | - | Current network mode |
| `lanip` | string | - | When switching to router mode, the returned LAN IP |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "netmode",
    "get_mode"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "mode": "router"
  }
}
```

---

## set_mode

Set network mode

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mode` | string | yes | Network mode, optional values: ap,router,relay,wds |
| `ssid` | string | no | SSID (required for relay and wds modes) |
| `identity` | string | no | EAP identity (required if WPA Enterprise encryption) |
| `key` | string | no | Password (required for relay and wds modes) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code: 1: Connection failed, 2: Failed to obtain IP, 3: IP segment conflict (LAN interface must be set to an address not in the same segment as the WiFi you want to connect to) |
| `err_msg` | number | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "netmode",
    "set_mode",
    {
      "mode": "router"
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
