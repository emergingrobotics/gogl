# `cable`

This is the API related to wired Internet.

4 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","cable","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`change_interface`](#change_interface) | - | Set wan port mode |
| [`get_config`](#get_config) | - | Get the wan port config |
| [`get_status`](#get_status) | - | Get the wan port status |
| [`set_config`](#set_config) | - | Set the WAN port access mode. |

---

## change_interface

Set wan port mode

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mode` | number | yes | Set WAN port mode,0: Set as WAN port,1: Set as LAN port. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1;parameter error,-6:remove internet cable. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"cable\",\"change_interface\",{\"mode\":1}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## get_config

Get the wan port config

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `protocol` | string | - | Identifies how protocol the WAN port obtains IP [dhcp/static/pppoe]. |
| `ipv6` | object | - | IPv6 information, returned when IPv6 is enabled and static protocol is set |
| `ipv6.ip` | string | - | ipv6 address |
| `ipv6.gateway` | string | - | ipv6 Gateway address |
| `ipv6.dns` | array | - | ipv6 DNS address |
| `ipv4` | object | - | IPv4 information, returned when static protocol is set |
| `ipv4.ip` | string | - | ipv4 address |
| `ipv4.netmask` | string | - | ipv4 Mask address |
| `ipv4.gateway` | string | - | ipv4 gateway address |
| `ipv4.dns` | array | - | ipv4 DNS address |
| `username` | string | - | PPPoE username, returned after setting to PPPoE dialing |
| `password` | string | - | PPPoE password, returned after setting to PPPoE dialing |
| `err_code` | number | - | Error code,-4: Device has no physical WAN port,-5: No virtual WAN port. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"cable\",\"get_config\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"protocol\":\"static\",\"ipv4\":{\"ip\":\"192.168.113.137\",\"netmask\":\"255.255.255.0\",\"gateway\":\"192.168.113.1\",\"dns\":[\"8.8.8.8\",\"8.8.4.4\"]}}}
```

---

## get_status

Get the wan port status

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `mode` | number | - | 0: Being used as WAN port,1: Being used as LAN port,2: Bridge mode. |
| `protocol` | string | - | Identifies how protocol the WAN port obtains IP [dhcp/static/pppoe]. |
| `status` | number | - | Status code,0: Connection failed,1: Connection successful,2: Connecting,3: Physical device not connected. |
| `ipv6` | object | - | IPv6 information, returned when status is 1 and IPv6 is enabled |
| `ipv6.ip` | string | - | ipv6 address |
| `ipv6.gateway` | string | - | ipv6 Gateway address |
| `ipv6.dns` | array | - | ipv6 DNS address |
| `ipv4` | object | - | IPv4 information, returned when status is 1 |
| `ipv4.ip` | string | - | ipv4 address |
| `ipv4.gateway` | string | - | ipv4 gateway address |
| `ipv4.dns` | array | - | ipv4 DNS address |
| `log` | string | - | Information during PPPoE dialing process, returned when set to PPPoE dialing |
| `err_code` | number | - | Error code,-4: Device has no physical WAN port,-5: No virtual WAN port. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"cable\",\"get_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"mode\": 0,\"status\":1,\"protocol\":\"static\",\"ipv4\":{\"ip\":\"192.168.113.137/24\",\"gateway\":\"192.168.113.1\",\"dns\":[\"8.8.8.8\",\"8.8.4.4\"]}}}
```

---

## set_config

Set the WAN port access mode.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `ipv6` | object | no | IPv6 information, input when setting static protocol and IPv6 is enabled |
| `ipv4` | object | no | IPv4 information, input when setting static protocol |
| `protocol` | string | yes | Identifies how protocol the WAN port obtains IP[dhcp/static/pppoe] |
| `username` | string | no | PPPoE username, input when setting PPPoE dialing |
| `password` | string | no | PPPoE password, input when setting PPPoE dialing |
| `ipv6.ip` | string | no | ipv6 address |
| `ipv6.gateway` | string | no | ipv6 Gateway address |
| `ipv4.ip` | string | no | ipv4 address |
| `ipv4.netmask` | string | no | ipv4 Mask address |
| `ipv4.gateway` | string | no | ipv4 gateway address |
| `ipv6.dns` | array | no | ipv6 DNS address |
| `ipv4.dns` | array | no | ipv4 DNS address |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1;parameter error,-2:ip format error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"cable\",\"set_config\",{\"protocol\":\"pppoe\",\"username\":\"test\",\"password\":\"123456\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---
