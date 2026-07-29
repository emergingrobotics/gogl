# `qos`

QOS management related interfaces

15 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","qos","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_device_group`](#add_device_group) | - | Add client group |
| [`add_speed_limit_rule`](#add_speed_limit_rule) | - | Add client device speed limit rule |
| [`delete_device_group`](#delete_device_group) | - | Delete client group |
| [`delete_speed_limit_rule`](#delete_speed_limit_rule) | - | Delete client device speed limit rule |
| [`enable_qos`](#enable_qos) | - | Enable or disable QoS function |
| [`get_client_list (在点击添加客户端限速规则时调用)`](#get_client_list-在点击添加客户端限速规则时调用) | - | Get client list |
| [`get_config`](#get_config) | verified | Get QoS configuration information |
| [`get_device_group`](#get_device_group) | - | Get client group information |
| [`modify_device_group`](#modify_device_group) | - | Modify client group information |
| [`modify_speed_limit_rule`](#modify_speed_limit_rule) | - | Modify client device speed limit rule |
| [`remove_speed_limit_rule`](#remove_speed_limit_rule) | - | Delete client device speed limit rule |
| [`set_channel_bandwidth_ratio`](#set_channel_bandwidth_ratio) | - | Configure priority channel bandwidth ratio |
| [`set_model`](#set_model) | - | Set QoS mode |
| [`set_other_client_priority`](#set_other_client_priority) | - | Set priority for other clients |
| [`set_speed_limit_rule`](#set_speed_limit_rule) | - | Set client device speed limit rule |

---

## add_device_group

Add client group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_name` | string | yes | Group name. |
| `clients.mac` | string | no | Client's MAC address. |
| `clients.name` | string | no | Client's name. |
| `clients` | array | no | Client list. |
| `priority` | number | yes | Group priority. |

| Results | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | - | Group id. |
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"block_client\",{\"mac\":\"84:7a:88:79:e5:13\",\"enable\":true}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## add_speed_limit_rule

Add client device speed limit rule

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | Client's MAC address. |
| `upload` | number | yes | Upload [Kbit]. |
| `download` | number | yes | Download [Kbit]. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"add_speed_limit_rule\",{\"mac\":\"18:c0:4d:dc:4f:cc\",\"download\":10,\"upload\":2}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## delete_device_group

Delete client group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group id. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"block_client\",{\"mac\":\"84:7a:88:79:e5:13\",\"enable\":true}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## delete_speed_limit_rule

Delete client device speed limit rule

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | Client's MAC address. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"delete_speed_limit_rule\",{\"mac\":\"18:c0:4d:dc:4f:cc\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## enable_qos

Enable or disable QoS function

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `enable` | bool | yes | true-enable QoS function, false-disable QoS function. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"set_config\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## get_client_list (在点击添加客户端限速规则时调用)

Get client list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `clients` | array | - | Client list. |
| `clients.mac` | string | - | Client's MAC address. |
| `clients.upload` | number | - | Client's upload [KB/S]. |
| `clients.download` | number | - | Client's download [KB/S]. |
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"get_client_list\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"clients\":[{\"mac\":\"18:c0:4d:dc:4f:cc\",\"download\":10,\"upload\":2},{\"mac\":\"18:c0:4d:dc:4f:dd\",\"download\":11,\"upload\":3}]}}
```

---

## get_config

Get QoS configuration information

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enable` | bool | - | QoS switch status. |
| `model` | string | - | Current working mode of QoS. (Return configuration information of the corresponding mode based on the current mode) |
| `wan_list` | array | - | QoS WAN port configuration information. |
| `wan_list.name` | string | - | WAN port name wan-wired network, wwan-wireless relay, tethering-hotspot sharing, modem1-mobile network 1, modem2-mobile network 2. |
| `wan_list.up` | bool | - | Whether the corresponding WAN port is enabled, false-not enabled, true-enabled. (Frontend does not display uplink and downlink configuration information for unenabled WAN ports) |
| `wan_list.upload` | number | - | Upload of the corresponding WAN port, in kbit. |
| `wan_list.download` | number | - | Download of the corresponding WAN port, in kbit. |
| `cake_model` | string | - | Current working mode of cake. |
| `group_device` | array | - | Device group information. |
| `group_device.name` | string | - | Name of the device group. |
| `group_device.id` | number | - | ID of the device group. |
| `group_device.number` | number | - | Total number of devices in the device group. |
| `other_clients_priority` | number | - | Priority of other clients [1-high priority, 2-medium priority, 3-low priority]. (Consider adding a button in the interface) |
| `speed_limit_clients` | array | - | List of speed-limited clients. |
| `speed_limit_clients.name` | string | - | Name of the client. |
| `speed_limit_clients.mac` | string | - | MAC address of the client. |
| `speed_limit_clients.ip` | string | - | IP address of the client. |
| `speed_limit_clients.upload` | string | - | Upload of the client. |
| `speed_limit_clients.download` | string | - | Download of the client. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"get_config\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## get_device_group

Get client group information

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group id. |

| Results | Type | Required | Description |
|---|---|---|---|
| `group_name` | string | - | Group name. |
| `priority` | number | - | Group priority. |
| `clients` | array | - | Client list. |
| `clients.mac` | string | - | Client's MAC address. |
| `clients.name` | string | - | Client's name. |
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"block_client\",{\"mac\":\"84:7a:88:79:e5:13\",\"enable\":true}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## modify_device_group

Modify client group information

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_name` | string | yes | Group name. |
| `clients.mac` | string | yes | Client's MAC address. |
| `clients.name` | string | yes | Client's name. |
| `clients` | array | yes | Client list. |
| `group_id` | number | yes | Group id. |
| `priority` | number | yes | Group priority. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"block_client\",{\"mac\":\"84:7a:88:79:e5:13\",\"enable\":true}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## modify_speed_limit_rule

Modify client device speed limit rule

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | Client's MAC address. |
| `upload` | number | yes | Upload [Kbit]. |
| `download` | number | yes | Download [Kbit]. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-4:Mac no exist. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"modify_speed_limit_rule\",{\"mac\":\"18:c0:4d:dc:4f:cc\",\"download\":10,\"upload\":2}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## remove_speed_limit_rule

Delete client device speed limit rule

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | Client's MAC address. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"remove_speed_limit_rule\",{\"mac\":\"18:c0:4d:dc:4f:cc\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## set_channel_bandwidth_ratio

Configure priority channel bandwidth ratio

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `upload` | array | yes | Upload bandwidth configuration. |
| `download` | array | yes | Download bandwidth configuration. |
| `upload.priority` | number | yes | Upload bandwidth priority. |
| `upload.min` | number | yes | Minimum upload bandwidth for the corresponding priority. |
| `upload.max` | number | yes | Maximum upload bandwidth for the corresponding priority. |
| `download.priority` | number | yes | Download bandwidth priority. |
| `download.min` | number | yes | Minimum download bandwidth for the corresponding priority. |
| `download.max` | number | yes | Maximum download bandwidth for the corresponding priority. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"block_client\",{\"mac\":\"84:7a:88:79:e5:13\",\"enable\":true}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## set_model

Set QoS mode

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `model` | string | yes | QoS startup mode [ speed_limit \| priority \| cake ]. (For limiting client rate mode, should add an apply button) |
| `wan_list.name` | string | no | WAN port name wan-wired network, wwan-wireless relay, tethering-hotspot sharing, modem1-mobile network 1, modem2-mobile network 2. |
| `cake_model` | string | no | Current working mode of cake. |
| `wan_list.up` | bool | no | Whether the corresponding WAN port is enabled, false-not enabled, true-enabled. (Frontend does not display uplink and downlink configuration information for unenabled WAN ports) |
| `wan_list` | array | no | QoS WAN port configuration information. |
| `wan_list.upload` | number | no | Upload of the corresponding WAN port, in kbit. |
| `wan_list.download` | number | no | Download of the corresponding WAN port, in kbit. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"switch_model\",{\"model\":\"cake\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## set_other_client_priority

Set priority for other clients

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `priority` | number | yes | Priority of other clients [1-High priority, 2-Medium priority, 3-Low priority]. (Consider adding a button in the interface) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"set_config\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---

## set_speed_limit_rule

Set client device speed limit rule

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | Client's MAC address. |
| `upload` | number | yes | Upload [KB/S]. When set to 0, no speed limit |
| `download` | number | yes | Download [KB/S]. When set to 0, no speed limit |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:Invalid user,-5:No parameter found. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"id\":1,\"params\":[\"\",\"qos\",\"set_speed_limit_rule\",{\"mac\":\"18:c0:4d:dc:4f:cc\",\"download\":10,\"upload\":2}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\": null}
```

---
