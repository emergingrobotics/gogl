# `wifi`

WiFi related operations

4 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","wifi","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | verified | Get all WiFi configuration information of the system |
| [`get_status`](#get_status) | verified | Get WiFi device status |
| [`set_config`](#set_config) | - | Set WiFi parameters |
| [`set_txpower`](#set_txpower) | - | Set RF transmit power |

---

## get_config

Get all WiFi configuration information of the system

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `dfs_support` | bool | - | Whether DFS is supported |
| `res` | array | - | All WiFi configuration information |
| `res.hwmodes` | array | - | List of supported hardware modes |
| `res.htmodes` | array | - | List of supported bandwidths |
| `res.channels` | array | - | List of supported channels |
| `res.encryptions` | array | - | List of supported encryptions |
| `res.band` | string | - | WiFi band ("2G" or "5G") |
| `res.hwmode` | string | - | Currently configured hardware mode |
| `res.htmode` | string | - | Currently configured bandwidth |
| `res.channel` | number | - | Currently configured channel (0 means auto) |
| `res.txpower` | string | - | Currently configured power |
| `res.device` | string | - | WiFi device name |
| `res.ifaces` | array | - | List of WiFi interface configurations |
| `res.ifaces.enabled` | bool | - | Whether enabled |
| `res.ifaces.ssid` | string | - | WiFi name |
| `res.ifaces.hidden` | bool | - | Whether to hide SSID |
| `res.ifaces.encryption` | string | - | Encryption method |
| `res.ifaces.key` | string | - | Password |
| `res.ifaces.guest` | bool | - | Whether it is guest WiFi |
| `res.ifaces.name` | string | - | Interface name |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "wifi",
    "get_config"
  ]
}
```

**Response**

```json
{
  "id": null,
  "jsonrpc": "2.0",
  "result": {
    "res": [
      {
        "band": "5G",
        "ifaces": [
          {
            "enabled": true,
            "ssid": "test",
            "encryption": "psk2",
            "name": "default_radio0",
            "guest": false,
            "hidden": false,
            "key": "goodlife"
          },
          {
            "enabled": false,
            "ssid": "test-5G-Guest",
            "encryption": "psk2",
            "name": "guest5g",
            "guest": true,
            "hidden": false,
            "key": "goodlife"
          }
        ],
        "channels": [
          {
            "dfs": false,
            "channel": 36
          },
          {
            "dfs": true,
            "channel": 64
          },
          {
            "dfs": true,
            "channel": 56
          },
          {
            "dfs": false,
            "channel": 44
          },
          {
            "dfs": false,
            "channel": 48
          },
          {
            "dfs": true,
            "channel": 116
          }
        ],
        "device": "radio0",
        "hwmode": "11ac/ax",
        "txpower": "Max",
        "channel": 40,
        "hwmodes": [
          "11ac/ax",
          "11n/ac/ax",
          "11a/n/ac/ax"
        ],
        "htmode": "auto",
        "htmodes": [
          "20",
          "40",
          "80",
          "auto"
        ],
        "encryptions": [
          "none",
          "psk2",
          "psk-mixed",
          "sae",
          "sae-mixed"
        ],
        "ready": false
      },
      {
        "band": "2G",
        "ifaces": [
          {
            "enabled": true,
            "ssid": "test-c8f",
            "encryption": "psk2",
            "name": "default_radio1",
            "guest": false,
            "hidden": false,
            "key": "goodlife"
          },
          {
            "enabled": false,
            "ssid": "test-Guest",
            "encryption": "psk2",
            "name": "guest2g",
            "guest": true,
            "hidden": false,
            "key": "goodlife"
          }
        ],
        "channels": [
          {
            "dfs": false,
            "channel": 1
          },
          {
            "dfs": false,
            "channel": 3
          },
          {
            "dfs": false,
            "channel": 2
          }
        ],
        "device": "radio1",
        "hwmode": "11n/ax",
        "txpower": "Max",
        "channel": 0,
        "hwmodes": [
          "11n/ax",
          "11g/n/ax",
          "11b/g/n/ax"
        ],
        "htmode": "20",
        "htmodes": [
          "20",
          "40",
          "auto"
        ],
        "encryptions": [
          "none",
          "psk2",
          "psk-mixed",
          "sae",
          "sae-mixed"
        ],
        "ready": false
      }
    ]
  }
}
```

---

## get_status

Get WiFi device status

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `res` | array | - | WiFi device status information |
| `res.name` | string | - | device name |
| `res.state` | string | - | WiFi device status: starting: starting, ready: normal, conflict: conflict with relay channel |
| `res.channel` | number | - | current channel |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "wifi",
    "get_status"
  ]
}
```

**Response**

```json
{
  "id": null,
  "jsonrpc": "2.0",
  "result": {
    "res": [
      {
        "name": "radio0",
        "state": "ready",
        "channel": 49
      },
      {
        "name": "radio1",
        "state": "starting"
      }
    ]
  }
}
```

---

## set_config

Set WiFi parameters

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `device` | string | no | Device name (selectable values obtained from get interface) |
| `hwmode` | string | no | Hardware mode (requires device parameter, selectable values obtained from get interface) |
| `htmode` | string | no | Bandwidth (requires device parameter, selectable values obtained from get interface) |
| `txpower` | string | no | Transmit power (requires device parameter, optional values: Max, High, Medium, Low) |
| `iface_name` | string | no | Interface name (obtained from get interface) |
| `ssid` | string | no | WiFi name (no more than 32 characters, requires iface_name parameter) |
| `encryption` | string | no | Encryption method (requires iface_name parameter, selectable values obtained from get interface) |
| `key` | string | no | Password (requires iface_name parameter, 8 ~ 63 characters) |
| `enabled` | bool | no | Whether to enable WiFi (requires iface_name parameter) |
| `hidden` | bool | no | Whether to hide WiFi (requires iface_name parameter) |
| `channel` | number | no | Channel (requires device parameter, 0 means auto, selectable values obtained from get interface) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | number | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "wifi",
    "set_config",
    {
      "device": "radio0",
      "channel": 1
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

## set_txpower

Set RF transmit power

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `device` | string | yes | Device name (obtained from get_config interface) |
| `txpower` | string | yes | Transmit power (optional values: Max, High, Medium, Low) |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "wifi",
    "set_txpower",
    {
      "device": "radio0",
      "txpower": "Max"
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
