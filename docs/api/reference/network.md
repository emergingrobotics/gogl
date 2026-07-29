# `network`

Network related operations

5 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","network","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`check_wan_cable`](#check_wan_cable) | - | Get whether cable WAN is connected, whether MAC clone is enabled |
| [`get_arp_list`](#get_arp_list) | - | Get ARP table |
| [`get_dhcp_leases`](#get_dhcp_leases) | verified | Get DHCP lease list |
| [`routes`](#routes) | - | Get active IPv4 routing table |
| [`routes6`](#routes6) | - | Get active IPv6 routing table |

---

## check_wan_cable

Get whether cable WAN is connected, whether MAC clone is enabled

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `cable_enabled` | bool | - | Check if WAN port has cable connection [true for cable connected, false for no cable connected]. |
| `macclone_enabled` | bool | - | Whether MAC clone is enabled [true for enabled, false for not enabled]. |
| `err_code` | number | - | Error code,-1:获取信息失败 cable not in wan. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "network",
    "check_wan_cable",
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
    "cable_enabled": true,
    "macclone_enabled": false
  }
}
```

---

## get_arp_list

Get ARP table

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `entries` | array | - | ARP list |
| `entries.ip` | string | - | IP address |
| `entries.mac` | string | - | MAC address |
| `entries.device` | string | - | Network interface |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "network",
    "get_arp_list"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "entries": [
      {
        "ip": "192.168.8.2",
        "mac": "11:22:33:44:55:66",
        "device": "eth0.2"
      },
      {
        "ip": "192.168.8.3",
        "mac": "11:22:33:44:55:26",
        "device": "eth0.2"
      }
    ]
  }
}
```

---

## get_dhcp_leases

Get DHCP lease list

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `leases` | array | - | DHCP lease list |
| `leases.ip` | string | - | IP address |
| `leases.mac` | string | - | MAC address |
| `leases.hostname` | string | - | Hostname |
| `leases.expires` | number | - | Remaining lease (seconds) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "network",
    "get_dhcp_leases"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "entries": [
      {
        "ip": "192.168.8.2",
        "mac": "11:22:33:44:55:66",
        "hostname": "test",
        "expires": 12
      }
    ]
  }
}
```

---

## routes

Get active IPv4 routing table

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `routes` | array | - | IPv4 routing table |
| `entries.target` | string | - | Target |
| `entries.nexthop` | string | - | Next hop |
| `entries.metric` | number | - | Metric |
| `entries.device` | string | - | Network interface |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "network",
    "routes"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "entries": [
      {
        "target": "192.168.8.2",
        "nexthop": "192.168.113.1",
        "metric": 0,
        "device": "eth0.2"
      }
    ]
  }
}
```

---

## routes6

Get active IPv6 routing table

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `routes` | array | - | IPv4 routing table |
| `entries.target` | string | - | Target |
| `entries.source` | string | - | Source |
| `entries.nexthop` | string | - | Next hop |
| `entries.metric` | number | - | Metric |
| `entries.device` | string | - | Network interface |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "network",
    "routes6"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "entries": [
      {
        "target": "::/0",
        "source": "fd94:ba1d:265d:8::/62",
        "nexthop": "fe80::5054:ff:fe00:8145",
        "metric": 512,
        "device": "eth0.2"
      }
    ]
  }
}
```

---
