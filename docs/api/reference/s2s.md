# `s2s`

This is the API related to s2s.

7 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","s2s","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`enable_echo_server`](#enable_echo_server) | - | enable echo listening server for testing |
| [`generate_wg_genkey`](#generate_wg_genkey) | - | Generate a new private key and return its public key. |
| [`get_status`](#get_status) | - | Get WireGuard service status. |
| [`remove_config`](#remove_config) | - | Delete WireGuard configuration |
| [`set_config`](#set_config) | - | Push the complete WireGuard configuration file. |
| [`start_wg`](#start_wg) | - | Start WireGuard service |
| [`stop_wg`](#stop_wg) | - | Stop WireGuard service |

---

## enable_echo_server

enable echo listening server for testing

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `port` | string | yes | Set the server listening port. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1:参数错误,-2:文件不存在,-3:uci 初始化错误. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"s2s\",\"enable_echo_server\",{\"port\":\"51820\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## generate_wg_genkey

Generate a new private key and return its public key.

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `public_key` | string | - | Public key |
| `err_code` | number | - | Error code,-3:uci 初始化错误. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"s2s\",\"generate_wg_genkey\",{}],\"id\":1}
```

**Response**

```json
{\"id\": 1,\"jsonrpc\": \"2.0\",\"result\": {\"public_key\": \"RYqIcKGF2Zs1obKDOqmxXH/WsnYn4/OwnExQXBgN7xE=\"}}
```

---

## get_status

Get WireGuard service status.

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `status` | string | - | Display WireGuard status |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"s2s\",\"get_status\",{}],\"id\":1}
```

**Response**

```json
{\"id\": 1,\"jsonrpc\": \"2.0\",\"result\": {\"status\": \"\"}}
```

---

## remove_config

Delete WireGuard configuration

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"s2s\",\"remove_config\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_config

Push the complete WireGuard configuration file.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `interface` | object | yes | WireGuard interface information. |
| `lan_ip` | string | no | Set. |
| `interface.private_key` | string | no | Private key. |
| `interface.listen_port` | string | no | Listening port. |
| `interface.mtu` | string | no | mtu. |
| `interface.metric` | string | no | metric. |
| `interface.fwmark` | string | no | fwmark. |
| `peer.public_key` | string | yes | Public key. |
| `peer.endpoint` | string | no | Host and port. |
| `peer.keepalive` | string | no | Keepalive interval (0 to disable). |
| `peer.preshared_key` | string | no | Preshared key. |
| `restart` | bool | no | Whether to restart the WireGuard service. |
| `interface.address` | array | yes | Address. |
| `peer` | array | yes | WireGuard Peer information. |
| `peer.allowed_ips` | array | yes | Encryption key routing. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1:参数错误,-2:文件不存在,-3:uci 初始化错误. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\": \"2.0\",\"method\": \"call\",\"params\": [\"\",\"s2s\",\"set_config\",{\"lan_ip\": \"192.168.9.1\",\"interface\": {\"mtu\": \"1560\",\"listen_port\": \"51820\",\"address\": [\"10.0.0.1\",\"0:0:0:0:0:ffff:a00:1\"]},\"peer\": [{\"public_key\": \"4Mss1+aEQMh4z8Ykgs/pwazmgjhWEzil9QreL2HZoEc=\",\"allowed_ips\": [\"192.168.8.0/24\",\"192.168.9.0/24\"]},{\"public_key\": \"6I13DX9LcF5ON15vYOfceSk3PzS2JjLyG0BbkspzDXc=\",\"allowed_ips\": [\"192.168.10.0/24\",\"0:0:0:0:0:ffff:c0a8:a00\"]}]}],\"id\": 1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## start_wg

Start WireGuard service

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"s2s\",\"start_wg\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## stop_wg

Stop WireGuard service

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"s2s\",\"stop_wg\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---
