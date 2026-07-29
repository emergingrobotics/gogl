# `clients`

Client management related interfaces

5 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","clients","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`block_client`](#block_client) | - | Set client's blacklist attribute |
| [`get_list`](#get_list) | verified | Get device list |
| [`get_status`](#get_status) | verified | Get client online total status |
| [`remove_offline`](#remove_offline) | - | Delete offline clients |
| [`set_info`](#set_info) | - | Set name |

---

## block_client

Set client's blacklist attribute

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | MAC address of the client to set blacklist for. |
| `block` | bool | yes | If true, the specified MAC client cannot access the Internet. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"clients\",\"block_client\",{\"mac\":\"84:7a:88:79:e5:13\",\"block\":true}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## get_list

Get device list

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `clients` | array | - | Client list |
| `clients.mac` | string | - | Client's MAC address. |
| `clients.ip` | string | - | Client's IP. |
| `clients.tx` | number | - | Client's upload speed (in Bps). |
| `clients.rx` | number | - | Client's download speed (in Bps). |
| `clients.total_tx` | number | - | Client's total upload traffic (in Bytes) |
| `clients.total_rx` | number | - | Client's total download traffic (in Bytes). |
| `clients.limit_tx` | number | - | Set client's upload speed limit (in KB/S). |
| `clients.limit_rx` | number | - | Set client's download speed limit (in KB/S). |
| `clients.blocked` | bool | - | GL firewall blacklist, true to add to blacklist, false not to add. |
| `clients.online_time` | string | - | Timestamp when client comes online. |
| `clients.online` | bool | - | true for online, false for offline. |
| `clients.iface` | string | - | Name of the interface the client is connected to. |
| `clients.type` | number | - | Client's connected interface type,0-2.4G,1-5G,2-lan,3-2.4G guest,4-5G guest,5-unknown,6-Dongle. |
| `clients.name` | string | - | Client device's hostname [not returned if no value]. |
| `clients.remote` | bool | - | If true, indicates that the current client is connecting to the router. |
| `clients.alias` | string | - | User-set name [not returned if no value]. |
| `clients.class` | string | - | User-set device type [not returned if no value]. |
| `err_code` | number | - | Error code. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"clients\",\"get_list\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"clients\":[{\"mac\":\"18:c0:4d:dc:4f:cc\",\"type\":2,\"limit_tx\":2048,\"limit_rx\":2048,\"remote\":false,\"name\":\"DESKTOP-HO0T5C1\",\"vendor\":\"unknown\",\"ip\":\"192.168.8.148\",\"tx\":0,\"total_rx\":397117728,\"rx\":0,\"blocked\":false,\"iface\":\"eth0\",\"online\":true,\"online_time\":\"1638517436\",\"total_tx\":470371319}]}}
```

---

## get_status

Get client online total status

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `cable_total` | number | - | Online total of wired clients. |
| `wireless_total` | number | - | Online total of wireless clients. |
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"clients\",\"get_status\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": {\"cable_total\":\"1\",\"wireless_total\":\"0\"}}
```

---

## remove_offline

Delete offline clients

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | MAC address of the client to delete (if value is FF:FF:FF:FF:FF:FF, delete all offline clients). |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"clients\",\"remove_offline\",{\"mac\":\"FF:FF:FF:FF:FF:FF\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## set_info

Set name

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | MAC address |
| `clients.alias` | string | no | User-set name |
| `clients.class` | string | no | User-set device type |

_No results._

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"clients\",\"set_info\",{\"mac\":\"16:c3:7b:72:19:74\",\"alias"\:\"test\",\"class\":\"tv\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---
