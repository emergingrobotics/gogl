# `switch_button`

Switch Button Configuration

3 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","switch_button","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | - | Get current configured function |
| [`get_funcs`](#get_funcs) | - | Get function list |
| [`set_config`](#set_config) | - | Get current configured function |

---

## get_config

Get current configured function

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `func` | string | - | Current configured function, "none" means not configured |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "switch-button",
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
    "func": "none"
  }
}
```

---

## get_funcs

Get function list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `funcs` | array | - | Function list |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "switch-button",
    "get_funcs"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "funcs": [
      "wireguard",
      "openvpn",
      "adguardhome"
    ]
  }
}
```

---

## set_config

Get current configured function

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `func` | string | yes | The function to configure (optional values obtained from get_funcs interface), currently optional openvpn,wireguard,adguardhome |

| Results | Type | Required | Description |
|---|---|---|---|
| `func` | string | - | Current configured function |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "switch-button",
    "set_config",
    {
      "func": "openvpn"
    }
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "func": "openvpn"
  }
}
```

---
