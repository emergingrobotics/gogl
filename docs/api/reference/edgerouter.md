# `edgerouter`

Bypass Router Mode

5 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","edgerouter","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | - | Get configurations related to the bypass router |
| [`scan`](#scan) | - | Scan the list of devices in the local area network |
| [`set_config`](#set_config) | - | Set configurations related to bypass router mode |
| [`set_devices`](#set_devices) | - | Set devices whose traffic passes through the bypass router |
| [`status`](#status) | - | Get the status of device traffic in the local area network passing through the bypass router |

---

## get_config

Get configurations related to the bypass router

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enable` | bool | - | Whether enabled; false: disabled, true: enabled |
| `mode` | number | - | Working mode; 0: cover all devices in the upstream LAN, 1: only cover specified devices in the upstream LAN, specified device list is returned in the devices parameter |
| `force_dns` | bool | - | Whether to force use of the bypass router's DNS; false: not used, true: used |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "edgerouter",
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
    "enable": true,
    "mode": 1,
    "force_dns": false
  }
}
```

---

## scan

Scan the list of devices in the local area network

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `devices` | array | - | Stores the list of scanned devices |
| `devices.mac` | string | - | Device MAC |
| `devices.ip` | string | - | Device IP |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "edgerouter",
    "scan",
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
    "devices": [
      {
        "mac": "00:11:22:33:44:55",
        "ip": "192.168.1.2"
      },
      {
        "mac": "00:11:22:33:44:56",
        "ip": "192.168.1.3"
      }
    ]
  }
}
```

---

## set_config

Set configurations related to bypass router mode

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `enable` | bool | yes | Whether enabled; false: disabled, true: enabled |
| `force_dns` | bool | yes | Whether to force use of the bypass router's DNS; false: not used, true: used |
| `mode` | number | yes | Working mode; 0: cover all devices in the upstream LAN, 1: only cover specified devices in the upstream LAN, specified devices are designated via the devices parameter |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "edgerouter",
    "set_config",
    {
      "enable": true,
      "mode": 1,
      "force_dns": false
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

## set_devices

Set devices whose traffic passes through the bypass router

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `devices` | array | yes | Specified device list; elements are the MAC addresses of the devices |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "edgerouter",
    "set_devices",
    {
      "devices": [
        "01:04:2d:1a:6e:10",
        "01:04:2d:1a:6e:11"
      ]
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

## status

Get the status of device traffic in the local area network passing through the bypass router

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `systime` | number | - | Current system timestamp |
| `devices` | array | - | Device status list |
| `devices.mac` | string | - | Device MAC |
| `devices.ip` | string | - | Device IP |
| `devices.lastalive` | number | - | The last time the device traffic passed through the bypass router, 0 means no traffic has ever passed (unit: seconds) |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "edgerouter",
    "status",
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
    "devices": [
      {
        "systime": 1656415756,
        "mac": "00:11:22:33:44:55",
        "ip": "192.168.1.2",
        "lastalive": 0
      },
      {
        "mac": "00:11:22:33:44:56",
        "ip": "192.168.1.3",
        "lastalive": 1656304044
      }
    ]
  }
}
```

---
