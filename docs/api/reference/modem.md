# `modem`

This is the API related to modems.

17 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","modem","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`disconnect`](#disconnect) | - | Disconnect network connection |
| [`get_cells_info`](#get_cells_info) | - | Get cellular network information |
| [`get_config`](#get_config) | absent | Read the current configuration of the modem |
| [`get_debug_msg`](#get_debug_msg) | - | Get debug information of the modem |
| [`get_info`](#get_info) | - | Get the modems hardware information |
| [`get_sms_list`](#get_sms_list) | - | Get SMS list |
| [`get_status`](#get_status) | - | Get the modems status |
| [`get_traffic_config`](#get_traffic_config) | - | Set automatic saving of traffic statistics data |
| [`reboot_modem`](#reboot_modem) | - | Reboot modem |
| [`remove_sms`](#remove_sms) | - | Delete SMS |
| [`reset_traffic_count`](#reset_traffic_count) | - | Clear traffic statistics data |
| [`send_at_command`](#send_at_command) | - | Send AT command to modem |
| [`send_sms`](#send_sms) | - | Send SMS |
| [`set_auto_connect`](#set_auto_connect) | - | Auto connect to the network |
| [`set_connect`](#set_connect) | - | Set the modem's dial-up configuration and connect |
| [`set_sms`](#set_sms) | - | Set SMS status |
| [`set_traffic_auto_save`](#set_traffic_auto_save) | - | Set automatic saving of traffic statistics data |

---

## disconnect

Disconnect network connection

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | The mounting point of the modem on the USB bus (can be obtained through get_info). |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1: No related module information found. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"disconnect\",{\"bus\":\"1-1\"}],\"id\":1}
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

## get_cells_info

Get cellular network information

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | The mounting point of the modem on the USB bus (can be obtained through get_info). |

| Results | Type | Required | Description |
|---|---|---|---|
| `cells` | array | - | Cellular network list |
| `cells.type` | string | - | Cellular network type (possible values: servingcell, neighbourcell) |
| `cells.id` | number | - | Cellular network ID (this parameter is returned only when type is servingcell) |
| `cells.mode` | string | - | Network mode (possible values: LTE, WCDMA, TDSCDMA, CDMA, HDR, GSM) |
| `cells.band` | number | - | Operating band (this parameter is returned only when type is servingcell and mode is LTE) |
| `cells.dl_bandwidth` | string | - | Download bandwidth (this parameter is returned only when type is servingcell and mode is LTE) |
| `cells.ul_bandwidth` | string | - | Upload bandwidth (this parameter is returned only when type is servingcell and mode is LTE) |
| `cells.rssi` | number | - | Received Signal Strength Indicator (this parameter is returned only when type is servingcell and mode is LTE) |
| `cells.rsrp` | number | - | Reference Signal Received Power (this parameter is returned only when type is servingcell and mode is LTE) |
| `cells.rsrq` | number | - | Reference Signal Received Quality (this parameter is returned only when type is servingcell and mode is LTE) |
| `cells.sinr` | number | - | Signal to Interference-plus-Noise Ratio (this parameter is returned only when type is servingcell and mode is LTE) |
| `cells.ecio` | number | - | Energy to Interference Ratio (this parameter is returned only when type is servingcell and mode is TDSCDMA, CDMA, HDR, etc.) |
| `err_code` | number | - | Error code, (-1: No related module information found, -2: Failed to get cells information). |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"get_cells_info\",{\"bus\":\"1-1\"}],\"id\":1}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "cells": [
      {
        "type": "servingcell",
        "id": 1234,
        "band": 41,
        "dl_bandwidth": "20M",
        "ul_bandwidth": "20M",
        "mode": "LTE",
        "rssi": -71,
        "rerp": -61,
        "rsrq": -12,
        "sinr": 13
      }
    ]
  }
}
```

---

## get_config

Read the current configuration of the modem

**Absent** on a GL-SFT1200 running 4.3.28: returns `-32601 Method not found`.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | Mount point of the modem on the USB bus (can be obtained via get_info). |

| Results | Type | Required | Description |
|---|---|---|---|
| `protocol` | string | - | Currently used protocol |
| `device` | string | - | Currently used connection device |
| `apn` | string | - | APN configuration |
| `service` | string | - | Currently used service type (this parameter is returned only when protocol is 3g) |
| `auth` | string | - | Currently used authentication method (possible values: NONE, PAP, CHAP, PAP/CHAP) |
| `username` | string | - | Currently used authentication username (this parameter is returned only when auth is not NONE) |
| `password` | string | - | Currently used authentication password (this parameter is returned only when auth is not NONE) |
| `dial_number` | string | - | Currently used dial number (this parameter is returned only when protocol is 3g and auth is not NONE) |
| `ttl` | number | - | Current TTL setting value |
| `pin_code` | string | - | Current PIN code setting |
| `err_code` | number | - | Error code,-1: No related module information found. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"get_config\",{\"bus\":\"1-1\"}],\"id\":1}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocol": "3g",
    "device": "ttyUSB3",
    "apn": "3gnet",
    "service": "CDC-WDMA",
    "auth": "PAP",
    "username": "gl-inet",
    "password": "goodlife1",
    "dial_number": "#9999",
    "ttl": "65",
    "pin_code": "1234"
  }
}
```

---

## get_debug_msg

Get debug information of the modem

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | The mounting point of the modem on the USB bus (can be obtained through get_info). |

| Results | Type | Required | Description |
|---|---|---|---|
| `msgs` | array | - | An array of objects, each object represents a status information, the specific content is variable, determined by the module characteristics |
| `msgs.cmd` | string | - | Executed AT command |
| `msgs.result` | string | - | Return result of the command |
| `err_code` | number | - | Error code, -1: No related module information found, -2: Failed to obtain information. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"get_debug_msg\",{\"bus\":\"1-1\"}],\"id\":1}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "msgs": [
      {
        "cmd": "AT",
        "result": "OK"
      }
    ]
  }
}
```

---

## get_info

Get the modems hardware information

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `modems` | array | - | Hardware-related information for all modems on the device. |
| `modems.devices` | array | - | Port device list for a single modem |
| `modems.type` | number | - | Modem type (0: Built-in compatible modem, 1: External compatible modem, 2: Incompatible modem). |
| `modems.bus` | string | - | Mount point of the modem on the USB bus |
| `modems.at_port` | string | - | Modem's AT command port |
| `modems.data_port` | string | - | Modem's data port, used with 3G protocol |
| `modems.sms_support` | bool | - | Whether the modem supports SMS function |
| `modems.imei` | string | - | Modem's IMEI number |
| `modems.name` | string | - | Modem's name (some modules include software version information) |
| `modems.vendor` | string | - | Modem's manufacturer |
| `modems.protocols` | array | - | Connection protocols supported by the modem (backend directly returns the list, possible protocols: 3g, ncm, qmi, qcm, directip, callip) |
| `modems.simcard` | object | - | SIM card related information (if no SIM card is detected, this object is not returned) |
| `modems.simcard.iccid` | string | - | SIM card's ICCID information (if the SIM card requires PIN unlock, this value may not be obtainable) |
| `modems.simcard.phone_number` | string | - | SIM card's phone number (if the SIM card requires PIN unlock or the SIM card does not store the local number, this value may not be obtainable) |
| `modems.simcard.mcc` | string | - | Mobile Country Code (used in combination with MNC to determine SIM card identity information, if the SIM card requires PIN unlock, this value may not be obtainable) |
| `modems.simcard.mnc` | string | - | Mobile Network Code (used in combination with MCC to determine SIM card identity information, if the SIM card requires PIN unlock, this value may not be obtainable) |
| `err_code` | number | - | Error code,-1: No related module information found. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"get_info\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"modems\":[{\"devices\":[\"/dev/ttyUSB1\",\"/dev/cdc-wdm0\"],\"protocols\":[\"3g",\"qmi\"],\"simcard\":{\"iccid\":\"1234567890\",\"phone_number\":\"8617603079735\",\"mcc\":\"123\",\"mnc\":\"456\"},\"type\":0,\"bus\":\"1-1\",\"sms_support\":true,\"name\":\" EC25EFAR06A04M4G\",\"vendor\":\"quectel\",\"imei\":\"860548047084717\"}]}}
```

---

## get_sms_list

Get SMS list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `list` | array | - | SMS list |
| `list.name` | string | - | Unique identifier of the SMS in the background, used for deletion, setting read status, etc. |
| `list.type` | number | - | SMS type (0:international, 1:national, 2:unknown) |
| `list.phone_number` | string | - | SMS source number |
| `list.bus` | string | - | The module bus number to which the SMS belongs, used to distinguish which module the SMS comes from |
| `list.date` | string | - | Date of receiving the SMS (format: 21-07-27 18:23:15) |
| `list.body` | string | - | SMS content (encoding format is UTF-8) |
| `list.status` | number | - | SMS status (0:Unread SMS, 1:Read SMS, 2:Sent SMS, 3:Sending, 4:Waiting to send, 5:Send failed) |
| `err_code` | number | - | Error code, (-1: No related module information found). |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"get_sms_list\",{}],\"id\":1}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "list": [
      {
        "type": "national",
        "name": "xy7rz",
        "phone_number": "15219475678",
        "modem": "1-1.2",
        "date": "21-07-27 18:23:15",
        "body": "hello",
        "status": 1
      }
    ]
  }
}
```

---

## get_status

Get the modems status

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `modems` | array | - | Status-related information for all modems on the device. |
| `modems.bus` | string | - | Mount point of the modem on the USB bus |
| `modems.simcard` | object | - | SIM card related status (if no SIM card is detected, this object is not returned) |
| `modems.simcard.status` | number | - | SIM card status (0: SIM card registered, 1: SIM card not registered, 2: SIM locked, requires PIN to unlock) |
| `modems.simcard.carrier` | string | - | Carrier name (if the SIM card requires PIN unlock or the SIM card does not store the local number, this value may not be obtainable) |
| `modems.simcard.signal` | object | - | SIM card signal related status (if the device is not properly registered to the network, this object is not returned) |
| `modems.simcard.signal.mode` | number | - | Network mode (2: 2G mode, 3: 3G mode, 4: 4G mode, 5: Reserved for 5G, 41: 4G+ mode) |
| `modems.simcard.signal.strength` | number | - | Signal strength (signal divided into 4 levels, represented by numbers 1, 2, 3, 4, larger number indicates stronger signal) |
| `modems.simcard.signal.rssi` | number | - | Received Signal Strength Indicator |
| `modems.simcard.signal.rsrp` | number | - | Reference Signal Received Power (this parameter is returned only in 4G and 4G+ modes) |
| `modems.simcard.signal.rsrq` | number | - | Reference Signal Received Quality (this parameter is returned only in 4G and 4G+ modes) |
| `modems.simcard.signal.sinr` | number | - | Signal to Interference-plus-Noise Ratio (this parameter is returned only in 4G and 4G+ modes) |
| `modems.simcard.signal.ecio` | number | - | Energy to Interference Ratio (this parameter is returned only in 3G mode) |
| `modems.network` | object | - | Network data connection related status (if the modem's data connection is not enabled, this object is not returned) |
| `modems.network.status` | number | - | Network data connection status (0: Connected, 1: Connecting) |
| `modems.network.protocol` | string | - | Protocol used for network data connection (this parameter is returned only when network.status is 0) |
| `modems.network.tx` | number | - | Amount of data sent, in bytes (this parameter is returned only when network.status is 0) |
| `modems.network.rx` | number | - | Amount of data received, in bytes (this parameter is returned only when network.status is 0) |
| `modems.network.ipv4` | object | - | Network data connection related status (if modems.network.status is not 0, this object is not returned) |
| `modems.network.ipv4.ip` | string | - | IPv4 address obtained by the modem |
| `modems.network.ipv4.netmask` | string | - | IPv4 subnet mask obtained by the modem |
| `modems.network.ipv4.gateway` | string | - | IPv4 gateway address obtained by the modem |
| `modems.network.ipv4.dns` | array | - | IPv4 DNS server list obtained by the modem |
| `modems.network.ipv6` | object | - | Network data connection related status (if IPv6 network is not enabled or modems.network.status is not 0, this object is not returned) |
| `modems.network.ipv6.ip` | string | - | IPv6 address and mask information obtained by the modem |
| `modems.network.ipv6.gateway` | string | - | IPv6 gateway address obtained by the modem |
| `modems.network.ipv6.dns` | array | - | IPv6 DNS server list obtained by the modem |
| `new_sms_count` | number | - | Number of unread SMS |
| `err_code` | number | - | Error code,-1: Failed to get status information. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"get_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"modems\":[{\"bus\":\"1-1\",\"simcard\":{\"status\":0,\"carrier\":\"AT&T\",\"signal\":{\"mode\":4,\"strength\":4,\"rssi\":-67,\"rsrp\":-60,\"rsrq\":-10,\"sinr\":13}},\"network\":{\"status\":0,\"protocol\":"3g",\"tx\":1024,\"rx\":4096,\"new_sms_count\":5,\"ipv4\":{\"ip\":\"103.25.102.101\",\"netmask\":\"255.255.255.0\",\"gateway\":\"103.25.102.100\",\"dns\":[\"114.114.114.114\",\"8.8.8.8\"]},\"ipv6\":{\"ip\":\"2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/64\",\"gateway\":\"2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b\",\"dns\":[\"2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b\",\"2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b\"]}}}]}}
```

---

## get_traffic_config

Set automatic saving of traffic statistics data

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | The mounting point of the modem on the USB bus (can be obtained through get_info). |

| Results | Type | Required | Description |
|---|---|---|---|
| `auto_save` | bool | - | Whether to enable automatic saving of traffic statistics data. |
| `err_code` | number | - | Error code, -1: No related module information found |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"get_traffic_config\",{\"bus\":\"1-1\"}],\"id\":1}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "auto_save": true
  }
}
```

---

## reboot_modem

Reboot modem

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | The mounting point of the modem on the USB bus (can be obtained through get_info). |
| `hw_reboot` | bool | no | Whether to use hardware power off/on to complete the reboot. |

| Results | Type | Required | Description |
|---|---|---|---|
| `wait_time` | number | - | The time to wait for the modem to complete the reboot, in seconds |
| `err_code` | number | - | Error code, -1: No related module information found, -2: Reboot failed |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"reboot_modem\",{\"bus\":\"1-1\"}],\"id\":1}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "wait_time": 25
  }
}
```

---

## remove_sms

Delete SMS

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | no | Unique identifier of the SMS in the background (can be obtained through get_sms_list, this parameter is only needed when scope value is 10) |
| `scope` | number | yes | The range of SMS to delete (possible values: 0: unread SMS, 1: read SMS, 2: sent SMS, 3: sending, 4: waiting to send, 5: send failed, 10: single SMS, 11: all SMS, 12: all received SMS, 13: all SMS in sending area) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1: No related module information found. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"remove_sms\",{\"bus\":\"1-1\",\"scope\":10,\"name\":\"xy7rz\"}],\"id\":1}
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

## reset_traffic_count

Clear traffic statistics data

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | Mount point of the modem on the USB bus (can be obtained through get_info). |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1: No related module information found |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"reset_traffic_count\",{\"bus\":\"1-1\"}],\"id\":1}
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

## send_at_command

Send AT command to modem

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | The mounting point of the modem on the USB bus (can be obtained through get_info). |
| `port` | string | yes | The port where the modem receives commands (the device's port list can be obtained through the ports parameter of get_info, only ttyUSBX ports support AT commands). |
| `command` | string | yes | AT command (command length does not exceed 256 bytes). |

| Results | Type | Required | Description |
|---|---|---|---|
| `response` | string | - | AT command execution result |
| `err_code` | number | - | Error code: -1: No related module information found, -2: command exceeds length limit, -3: Operation not supported, -4: Unknown error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"send_at_command\",{\"bus\":\"1-1\",\"port\":\"ttyUSB2\",\"command\":\"AT\"}],\"id\":1}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "response": "ok"
  }
}
```

---

## send_sms

Send SMS

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | Mount point of the modem on the USB bus (can be obtained through get_info). |
| `phone_number` | string | yes | Target number to receive SMS (note, needs to include country code, e.g., 8617603079726) |
| `body` | string | yes | SMS content (length not exceeding 140 bytes) |
| `timeout` | number | yes | Timeout time (in seconds), if set to 0, means do not wait for send result, return immediately |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: No related module information found, -2: Unknown error, -3: Send failed, -4: Send timeout, -5: Current device does not support SMS function, -6: Current network environment does not support SMS function). |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"get_sms_list\",{\"bus\":\"1-1\",\"phone_number\":\"8617603079726\",\"body\":\"hello\"}],\"id\":1}
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

## set_auto_connect

Auto connect to the network

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | Mount point of the modem on the USB bus (can be obtained via get_info). |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1: No related module information found. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"set_auto_connect\",{\"bus\":\"1-1\"}],\"id\":1}
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

## set_connect

Set the modem's dial-up configuration and connect

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `config` | object | yes | Configuration information |
| `bus` | string | yes | Mount point of the modem on the USB bus (can be obtained via get_info). |
| `config.protocol` | string | yes | Protocol to be used (list of supported protocols can be obtained via get_info) |
| `config.device` | string | yes | Connection device to be used (list of supported protocols can be obtained via get_info) |
| `config.apn` | string | no | APN configuration |
| `config.service` | string | no | Service type (this parameter is required only when protocol is 3g) |
| `config.auth` | string | no | Authentication method (possible values: PAP, CHAP, PAP/CHAP) |
| `config.username` | string | no | Authentication username (this parameter is valid only when auth is not NONE) |
| `config.password` | string | no | Authentication password (this parameter is valid only when auth is not NONE) |
| `config.dial_number` | string | no | Dial number (this parameter is valid only when protocol is 3g) |
| `config.pin_code` | string | no | PIN code required for SIM card unlock |
| `config.ttl` | number | no | Data packet TTL setting (if this parameter is set, the device will lock the TTL value of the data packet to the specified value) |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1: No related module information found. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"set_connect\",{"bus":"1-1","config":{"protocol":"3g","device":"ttyUSB3","apn":"3gnet","service":"CDCWDMA","auth":"PAP","username":"gl-inet","password":"goodlife1","dial_number":"#9999","ttl":65,"pin_code":"1234"}}],\"id\":1}
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

## set_sms

Set SMS status

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | no | Unique identifier of the SMS in the background, only needs to be passed when status value is not equal to 6 (can be obtained through get_sms_list) |
| `status` | number | yes | The SMS status to set (0: unread SMS, 1: "read SMS", 2: "sent SMS", 3: sending, 4: "waiting to send", 5: "send failed", 6: "set all SMS as read") |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, (-1: No related module information found, -2: SMS does not exist, -3: Failed to set SMS) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"set_sms\",{\"bus\":\"1-1\",\"name\":\"xy7rz\",\"status\":0}],\"id\":1}
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

## set_traffic_auto_save

Set automatic saving of traffic statistics data

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `bus` | string | yes | Mount point of the modem on the USB bus (can be obtained through get_info). |
| `enable` | bool | yes | Whether to enable automatic saving of traffic statistics data, if enabled, the automatic save interval is fixed at 30S. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1: No related module information found |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"modem\",\"set_traffic_auto_save\",{\"bus\":\"1-1\",\"enable\":true}],\"id\":1}
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
