# `repeater`

WiFi Repeater

9 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","repeater","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`connect`](#connect) | - | Connect to AP |
| [`disconnect`](#disconnect) | - | Disconnect the currently connected WiFi |
| [`forget`](#forget) | - | Forget a certain AP |
| [`get_config`](#get_config) | verified | Get configuration |
| [`get_saved_ap_list`](#get_saved_ap_list) | - | Get saved AP list |
| [`get_status`](#get_status) | - | Get connection status |
| [`remove_saved_ap`](#remove_saved_ap) | - | Delete a saved AP |
| [`scan`](#scan) | - | Scan AP list |
| [`set_config`](#set_config) | - | Set |

---

## connect

Connect to AP

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `ssid` | string | yes | SSID |
| `identity` | string | no | EAP identity (required if it's WPA Enterprise encryption) |
| `key` | string | no | Password (required if connecting to an encrypted AP) |
| `bssid` | string | no | BSSID (lock BSSID, empty means not locked) |
| `protocol` | string | no | Networking method (dhcp, static) |
| `ip` | string | no | IP address (required when protocol is static) |
| `netmask` | string | no | Subnet mask (required when protocol is static) |
| `gateway` | string | no | Subnet mask (required when protocol is static) |
| `remember` | bool | no | Save network |
| `manual` | bool | no | Manual connection |
| `dns` | array | no | DNS (required when protocol is static) |

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
    "repeater",
    "connect",
    {
      "ssid": "test",
      "key": "goodlife"
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

## disconnect

Disconnect the currently connected WiFi

Not exercised against hardware here.

_No params._

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "repeater",
    "disconnect"
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

## forget

Forget a certain AP

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `id` | string | yes | Obtained from the get_saved_aps interface |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "repeater",
    "forget",
    {
      "id": "cfg12312432"
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

## get_config

Get configuration

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `dfs_support` | bool | - | Whether DFS is supported |
| `auto` | bool | - | Allow switching to other saved networks |
| `antijam` | bool | - | Anti-interference mode |
| `lock_band` | string | - | Control connection to same-frequency hotspot |
| `dfs` | bool | - | Whether to allow relaying DFS channels |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "repeater",
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
    "auto": true,
    "lock_band": "2g"
  }
}
```

---

## get_saved_ap_list

Get saved AP list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `res` | array | - | AP list |
| `res.ssid` | string | - | SSID |
| `res.bssid` | string | - | BSSID (lock BSSID, empty means not locked) |
| `res.identity` | string | - | Identity ID (for EAP encryption) |
| `res.key` | string | - | Password |
| `res.protocol` | string | - | Networking method (dhcp, static) |
| `res.ip` | string | - | IP address (static networking) |
| `res.netmask` | string | - | Subnet mask (static networking) |
| `res.gateway` | string | - | Gateway (static networking) |
| `res.dns` | array | - | DNS (static networking) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "repeater",
    "get_saved_ap_list"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "res": [
      {
        "ssid": "test",
        "key": "12132432423"
      },
      {
        "ssid": "x",
        "key": "sefrsedtfsde"
      }
    ]
  }
}
```

---

## get_status

Get connection status

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `state` | number | - | Relay status (0: Idle; 1: Connecting; 2: Connected, 3: Connection failed) |
| `device` | string | - | WiFi device name |
| `ssid` | string | - | Connected SSID or SSID being connected |
| `bssid` | string | - | Connected BSSID |
| `channel` | number | - | Connected channel |
| `signal` | number | - | Signal strength (dBm) |
| `wait_dev` | string | - | If present, it means waiting for the radio to become available |
| `ipv4` | object | - | IPv4 address information |
| `ipv4.ip` | string | - | IPv4 address (CIDR format) |
| `ipv4.dns` | array | - | IPv4 DNS |
| `ipv4.gateway` | string | - | IPv4 gateway |
| `ipv6` | object | - | IPv6 address information |
| `ipv6.ip` | array | - | IPv6 address (CIDR format) |
| `ipv6.dns` | array | - | IPv6 DNS |
| `ipv6.gateway` | string | - | IPv6 gateway |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "repeater",
    "get_status"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "state": 3,
    "ssid": "test",
    "bssid": "11:22:33:44:55:66",
    "channel": 12,
    "signal": 20,
    "ipv4": {
      "ip": "192.168.8.12/24",
      "gateway": "192.168.1.1"
    },
    "ipv6": {
      "ip": [
        "fd4d:7210:ad91:10::/64"
      ],
      "dns": [
        "fd5c:d00d:fa8e::/64"
      ],
      "gateway": "fd5c:d00d:fa8e::/64"
    }
  }
}
```

---

## remove_saved_ap

Delete a saved AP

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `ssid` | string | yes | SSID |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "repeater",
    "remove_saved_ap",
    {
      "ssid": "test"
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

## scan

Scan AP list

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `all_band` | bool | yes | Force scan all bands, if false, only scan locked bands |
| `dfs` | bool | yes | Get APs including DFS channels, if false, return based on relay settings |

| Results | Type | Required | Description |
|---|---|---|---|
| `res` | array | - | AP list |
| `res.ssid` | string | - | SSID |
| `res.bssid` | string | - | BSSID |
| `res.signal` | number | - | Signal strength |
| `res.band` | number | - | Band (2g or 5g) |
| `res.saved` | bool | - | Whether it is saved |
| `res.device` | string | - | WiFi device (corresponding to "wifi-device" in /etc/config/wireless) |
| `res.encryption` | object | - | Security information |
| `res.encryption.enabled` | bool | - | Whether encrypted |
| `res.encryption.wpa` | number | - | WPA version (bit 0 indicates whether WPA1 is enabled; bit 1 indicates whether WPA2 is enabled; bit 3 indicates whether WPA3 is enabled) |
| `res.encryption.auth_suites` | array | - | Key management methods (PSK: Pre-shared key, SAE: Simultaneous Authentication of Equals, 802.1X: Enterprise) |
| `res.encryption.pair_ciphers` | array | - | Encryption algorithms (CCMP, TKIP) |
| `res.encryption.description` | string | - | Security description |
| `res.encryption.uci` | string | - | UCI configuration string |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "wifi",
    "scan"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "res": [
      {
        "ssid": "GL-AR750S-4e2-5G",
        "band": "5g",
        "encryption": {
          "enabled": true,
          "description": "WPA2 PSK (CCMP)",
          "auth_suites": [
            "PSK"
          ],
          "wpa": 2,
          "pair_ciphers": [
            "CCMP"
          ],
          "uci": "psk2"
        },
        "bssid": "94:83:C4:0C:54:E3",
        "channel": 36,
        "signal": -60
      }
    ]
  }
}
```

---

## set_config

Set

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `lock_band` | string | no | Control connection to same-frequency hotspot ("2g" or "5g", not passing means automatic) |
| `auto` | bool | yes | Allow switching to other saved networks (If this option is enabled, when the currently connected WiFi is unavailable, the router will attempt to connect to other saved networks) |
| `antijam` | bool | yes | Anti-interference mode (When this option is enabled, the device will force into low bandwidth mode, improving connection stability at the cost of reducing connection speed) |
| `dfs` | bool | no | Whether to allow relaying DFS channels |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "repeater",
    "set_config",
    {
      "auto": true,
      "lock_band": "2g"
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
