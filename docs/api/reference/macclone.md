# `macclone`

macclone settings

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","macclone","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_mac`](#get_mac) | - | Get the configuration information of mac |
| [`set_mac`](#set_mac) | - | Get the configuration information of mac |

---

## get_mac

Get the configuration information of mac

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `mac` | string | - | The MAC address of the router's WAN port.[Corresponding to Your Router (WAN)][Random] |
| `factory_mac` | string | - | The factory default MAC address of the router.[Corresponding to Factory Default][default] |
| `remote_mac` | string | - | The MAC address of the client connected to the router.[Corresponding to Your Current Client][clone] |
| `err_code` | number | - | Error code, -1: Failed to get information. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "macclone",
    "get_mac"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "mac": "08:10:7B:B3:E8:92",
    "factory_mac": "94:83:C4:0C:6D:D6",
    "remote_mac": "94:83:C4:0C:6D:D6"
  }
}
```

---

## set_mac

Get the configuration information of mac

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | The MAC address to be cloned. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1: Illegal parameter. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "macclone",
    "set_mac",
    {
      "mac": "94:83:C4:0C:6D:D6"
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
