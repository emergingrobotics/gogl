# `firewall`

Firewall

11 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","firewall","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_port_forward`](#add_port_forward) | - | Add port forwarding |
| [`add_rule`](#add_rule) | - | Add firewall rule: Open port(src="wan", proto="tcp" dest_port=80, target="ACCEPT", enabled=true) |
| [`get_dmz`](#get_dmz) | - | Get DMZ configuration |
| [`get_port_forward_list`](#get_port_forward_list) | - | Get all port forwarding configurations |
| [`get_rule_list`](#get_rule_list) | - | Get all firewall rules |
| [`get_zone_list`](#get_zone_list) | - | Get all firewall zones |
| [`remove_port_forward`](#remove_port_forward) | - | Delete port forwarding |
| [`remove_rule`](#remove_rule) | - | Delete firewall rule |
| [`set_dmz`](#set_dmz) | - | Get DMZ configuration |
| [`set_port_forward`](#set_port_forward) | - | Modify port forwarding |
| [`set_rule`](#set_rule) | - | Modify firewall rule |

---

## add_port_forward

Add port forwarding

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | no | Name |
| `proto` | string | no | Protocol (default: "tcp udp", optional values: "tcp udp", "tcp", "udp") |
| `src` | string | yes | External zone (obtained from get_zone_list interface) |
| `src_dport` | string | yes | External port |
| `dest` | string | yes | Internal zone (obtained from get_zone_list interface) |
| `dest_ip` | string | yes | Internal IP address |
| `enabled` | bool | no | Whether enabled |
| `dest_port` | number | yes | Internal port |

| Results | Type | Required | Description |
|---|---|---|---|
| `id` | string | - | Configuration item ID |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "add_port_forward",
    {
      "name": "test",
      "proto": "tcp",
      "src": "wan",
      "src_dport": 80,
      "dest": "lan",
      "dest_ip": "192.168.8.100",
      "dest_port": 80
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
    "id": "cfg153837"
  }
}
```

---

## add_rule

Add firewall rule: Open port(src="wan", proto="tcp" dest_port=80, target="ACCEPT", enabled=true)

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | no | Name |
| `src` | string | no | Source zone (obtained from get_zone_list interface) |
| `src_ip` | string | no | Source IP address |
| `src_mac` | string | no | Source MAC address |
| `proto` | string | no | Protocol (optional values: "tcp udp", "tcp", "udp") |
| `dest` | string | no | Destination zone (obtained from get_zone_list interface) |
| `dest_ip` | string | no | Destination IP address |
| `target` | string | yes | Action (optional values: ACCEPT, REJECT, DROP) |
| `enabled` | bool | no | Whether enabled |
| `src_port` | number | no | Source port (1 ~ 65535) |
| `dest_port` | number | no | Destination port (1 ~ 65535) |

| Results | Type | Required | Description |
|---|---|---|---|
| `id` | string | - | Configuration item ID |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "add_rule",
    {
      "name": "test",
      "src": "wan",
      "dest_port": 22,
      "target": "ACCEPT"
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
    "id": "cfg1392bd"
  }
}
```

---

## get_dmz

Get DMZ configuration

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enabled` | bool | - | Whether enabled |
| `dest_ip` | string | - | Internal IP address |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "get_dmz"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "enabled": true,
    "dest_ip": "192.168.8.100"
  }
}
```

---

## get_port_forward_list

Get all port forwarding configurations

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `res` | array | - | Port forwarding configuration list |
| `res.name` | string | - | Name |
| `res.proto` | string | - | Protocol (default: "tcp udp", optional values: "tcp udp", "tcp", "udp") |
| `res.src` | string | - | External zone (obtained from get_zone_list interface) |
| `res.src_dport` | string | - | External port |
| `res.dest` | string | - | Internal zone (obtained from get_zone_list interface) |
| `res.dest_ip` | string | - | Internal IP address |
| `res.dest_port` | number | - | Internal port |
| `res.enabled` | bool | - | Whether enabled |
| `res.id` | string | - | Configuration item ID |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "get_port_forward_list"
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
        "enabled": true,
        "id": "cfg143837",
        "dest_ip": "192.168.8.100",
        "src": "wan",
        "proto": "tcp"
      }
    ]
  }
}
```

---

## get_rule_list

Get all firewall rules

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `res` | array | - | Rule list |
| `res.name` | string | - | Name |
| `res.src` | string | - | Source zone (obtained from get_zone_list interface) |
| `res.src_ip` | string | - | Source IP address |
| `res.src_mac` | string | - | Source MAC address |
| `res.src_port` | number | - | Source port |
| `res.proto` | string | - | Protocol |
| `res.dest` | string | - | Destination zone (obtained from get_zone_list interface) |
| `res.dest_ip` | string | - | Destination IP address |
| `res.dest_port` | number | - | Destination port |
| `res.enabled` | bool | - | Whether enabled |
| `res.target` | string | - | Action (optional values: ACCEPT, REJECT, DROP) |
| `res.id` | string | - | Configuration item ID |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "get_rule_list"
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
        "id": "cfg1392bd",
        "dest_port": 22,
        "name": "test",
        "target": "ACCEPT",
        "enabled": true,
        "src": "wan"
      }
    ]
  }
}
```

---

## get_zone_list

Get all firewall zones

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `internals` | array | - | Internal zone list |
| `externals` | array | - | External zone list |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "get_zone_list"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "internals": [
      "lan",
      "guest"
    ],
    "externals": [
      "wan"
    ]
  }
}
```

---

## remove_port_forward

Delete port forwarding

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `id` | string | yes | Configuration item ID (obtained from get_port_forward_list or add_port_forward interface) |
| `all` | bool | no | Whether to delete all (do not pass id when deleting all) |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "remove_port_forward",
    {
      "id": "cfg153837"
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

## remove_rule

Delete firewall rule

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `id` | string | no | UCI identifier (obtained from get_rule_list or add_rule interface) |
| `all` | bool | no | Whether to delete all (do not pass id when deleting all) |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "remove_rule",
    {
      "id": "cfg1392bd"
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

## set_dmz

Get DMZ configuration

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dest_ip` | string | no | Internal IP address (required when enabled) |
| `enabled` | bool | yes | Whether enabled |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "set_dmz",
    {
      "dest_ip": "192.168.8.100"
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

## set_port_forward

Modify port forwarding

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | no | Name |
| `proto` | string | no | Protocol (default: "tcp udp", optional values: "tcp udp", "tcp", "udp") |
| `src` | string | yes | External zone (obtained from get_zone_list interface) |
| `src_dport` | string | yes | External port |
| `dest` | string | yes | Internal zone (obtained from get_zone_list interface) |
| `dest_ip` | string | yes | Internal IP address |
| `id` | string | yes | Configuration item ID (obtained from get_port_forward_list or add_port_forward interface) |
| `enabled` | bool | yes | Whether enabled |
| `dest_port` | number | yes | Internal port |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "set_port_forward",
    {
      "id": "cfg153837",
      "name": "test",
      "proto": "tcp",
      "src": "wan",
      "src_dport": 80,
      "dest": "lan",
      "dest_ip": "192.168.8.100",
      "dest_port": 80
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

## set_rule

Modify firewall rule

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | no | Name |
| `src` | string | no | Source zone (obtained from get_zone_list interface) |
| `src_ip` | string | no | Source IP address |
| `src_mac` | string | no | Source MAC address |
| `proto` | string | no | Protocol (optional values: "tcp udp", "tcp", "udp") |
| `dest` | string | no | Destination zone (obtained from get_zones interface) |
| `dest_ip` | string | no | Destination IP address |
| `target` | string | no | Action (optional values: ACCEPT, REJECT, DROP) |
| `id` | string | yes | Configuration item ID (obtained from get_rules or add_rule interface) |
| `enabled` | bool | no | Whether enabled |
| `src_port` | number | no | Source port (1 ~ 65535) |
| `dest_port` | number | no | Destination port (1 ~ 65535) |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "firewall",
    "set_rule",
    {
      "id": "cfg1392bd",
      "src": "wan",
      "dest_port": 80
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
