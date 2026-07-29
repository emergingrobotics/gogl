# `vpn_policy`

VPN Policy

10 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","vpn_policy","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_domain_policy`](#get_domain_policy) | - | Get domain policy |
| [`get_global_policy`](#get_global_policy) | - | Get global policy |
| [`get_mac_policy`](#get_mac_policy) | - | Get MAC policy |
| [`get_proxy_mode`](#get_proxy_mode) | - | Get VPN proxy mode |
| [`get_vlan_policy`](#get_vlan_policy) | - | Get VLAN policy |
| [`set_domain_policy`](#set_domain_policy) | - | Set domain policy |
| [`set_global_policy`](#set_global_policy) | - | Set global policy |
| [`set_mac_policy`](#set_mac_policy) | - | Set MAC policy |
| [`set_proxy_mode`](#set_proxy_mode) | - | Set VPN proxy mode |
| [`set_vlan_policy`](#set_vlan_policy) | - | Set VLAN policy |

---

## get_domain_policy

Get domain policy

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `default_policy` | number | - | Default policy (0: do not use VPN, 1: use VPN) |
| `domain_list` | array | - | Domain/IP list |
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
    "vpn-policy",
    "get_domain_policy",
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
    "default_policy": 0,
    "domain_list": [
      "google.com",
      "facebook.com"
    ]
  }
}
```

---

## get_global_policy

Get global policy

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `kill_switch` | number | - | Whether to block traffic not going through VPN (0: Disable killswitch, do not block any traffic, 1: Block non-VPN traffic under any circumstances, 2: Block non-VPN traffic only when VPN client is enabled) |
| `wan_access` | number | - | When VPN is enabled, whether to allow LAN->WAN forwarding data (0: Not allowed, 1: Allowed) |
| `service_policy` | number | - | Whether to use VPN for GL.iNet services on the router (0: Follow system settings, 1: Do not use VPN) |
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
    "vpn-policy",
    "get_global_policy",
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
    "kill_switch": 0,
    "wan_access": 0,
    "service_policy": 1
  }
}
```

---

## get_mac_policy

Get MAC policy

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `default_policy` | number | - | Default policy (0: do not use VPN, 1: use VPN) |
| `mac_list` | array | - | MAC list |
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
    "vpn-policy",
    "get_mac_policy",
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
    "default_policy": 0,
    "mac_list": [
      "94:83:C4:0E:10:DF",
      "94:83:C4:0E:20:DF"
    ]
  }
}
```

---

## get_proxy_mode

Get VPN proxy mode

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `mode` | number | - | Proxy mode (0: Global proxy, 1: Auto detect route, 2: Custom route, 3: Based on destination site, 4: Based on client device, 5: Based on VLAN) |
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
    "vpn-policy",
    "get_proxy_mode",
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
    "mode": 0
  }
}
```

---

## get_vlan_policy

Get VLAN policy

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `vlans` | array | - | VLAN information |
| `vlans.id` | number | - | VLAN ID (1: private, 2: guest) |
| `vlans.vpn` | number | - | Whether VLAN uses VPN (0: disabled, 1: enabled) |
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
    "vpn-policy",
    "get_vlan_policy",
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
    "vlans": [
      {
        "id": 1,
        "vpn": 0
      },
      {
        "id": 2,
        "vpn": 0
      }
    ]
  }
}
```

---

## set_domain_policy

Set domain policy

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `domain_list` | array | yes | Domain/IP list |
| `default_policy` | number | yes | Default policy (0: Do not use VPN, 1: Use VPN) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "vpn-policy",
    "set_domain_policy",
    {
      "default_policy": 0,
      "domain_list": [
        "google.com",
        "facebook.com"
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

## set_global_policy

Set global policy

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `kill_switch` | number | yes | Whether to block traffic not going through VPN (0: Disable killswitch, do not block any traffic, 1: Block non-VPN traffic under any circumstances) |
| `wan_access` | number | yes | When VPN is enabled, whether to allow LAN->WAN forwarding data (0: Not allowed, 1: Allowed) |
| `service_policy` | number | yes | Whether to use VPN for GL.iNet services on the router (0: Follow system settings, 1: Do not use VPN) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: Missing required parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "vpn-policy",
    "set_global_policy",
    {
      "kill_switch": 0,
      "wan_access": 0,
      "service_policy": 1
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

## set_mac_policy

Set MAC policy

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac_list` | array | yes | MAC list |
| `default_policy` | number | yes | Default policy (0: do not use VPN, 1: use VPN) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "vpn-policy",
    "set_mac_policy",
    {
      "default_policy": 0,
      "mac_list": [
        "94:83:C4:0E:10:DF",
        "94:83:C4:0E:20:DF"
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

## set_proxy_mode

Set VPN proxy mode

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mode` | number | yes | Proxy mode (0: Global proxy, 1: Auto detect route, 2: Custom route, 3: Based on destination site, 4: Based on client device, 5: Based on VLAN) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "vpn-policy",
    "set_proxy_mode",
    {
      "mode": 0
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

## set_vlan_policy

Set VLAN policy

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `vlans` | array | yes | VLAN information |
| `vlans.id` | number | no | VLAN ID (1: private, 2: guest) |
| `vlans.vpn` | number | no | Whether VLAN uses VPN (0: disabled, 1: enabled) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "vpn-policy",
    "set_vlan_policy",
    {
      "vlans": [
        {
          "id": 1,
          "vpn": 0
        },
        {
          "id": 2,
          "vpn": 0
        }
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
