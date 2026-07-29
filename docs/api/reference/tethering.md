# `tethering`

This is the API related to tethering Internet access.

3 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","tethering","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`disconnect`](#disconnect) | - | disconnect tethering. |
| [`get_status`](#get_status) | - | Get tethering info. |
| [`set_connect`](#set_connect) | - | connect or disconnect tethering. |

---

## disconnect

disconnect tethering.

Not exercised against hardware here.

_No params._

_No results._

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"tethering\",\"disconnect\"],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## get_status

Get tethering info.

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `devices` | array | - | Device list |
| `devices.device` | string | - | Device name |
| `devices.type` | number | - | Device type 0 Android 1 Apple 2 Unknown |
| `devices.use` | bool | - | Usage status, true: in use, false: not in use |
| `status` | number | - | Status, 0: not connected, 1: connected successfully, 2: connecting |
| `ipv6` | object | - | ipv6 info |
| `ipv6.ip` | string | - | ipv6 address |
| `ipv6.gateway` | string | - | ipv6 Gateway address |
| `ipv6.dns` | array | - | ipv6 DNS address |
| `ipv4` | object | - | ipV4 info |
| `ipv4.ip` | string | - | ipv4 address |
| `ipv4.gateway` | string | - | ipv4 gateway address |
| `ipv4.dns` | array | - | ipv4 DNS address |
| `err_code` | number | - | Error code,-1:This feature is not supported,-2:The device is not plugged in. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"tethering\",\"get_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"devices\":[{\"name\":\"usb0\",\"tyep\":0,\"use\":true},{\"device\":\"usb1\",\"tyep\":0,\"use\":false}],\"status\":1,\"ipv4\":{\"ip\":\"192.168.42.38/24\",\"gateway\":\"192.168.42.129\",\"dns\":[\"192.168.42.129\"]}}}
```

---

## set_connect

connect or disconnect tethering.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `device` | string | yes | Device name. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-3:missing parameters. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"tethering\",\"set_connect\",{\"device\":\"usb0\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---
