# `system`

System Operations

16 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","system","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_user`](#add_user) | - | Add user |
| [`disk_info`](#disk_info) | - | Get disk information (unit: bytes B) |
| [`get_httpd_mem_status`](#get_httpd_mem_status) | - | Get the memory usage of the HTTP server. When testing for memory leaks, this interface can be used to monitor if the API has memory leaks |
| [`get_info`](#get_info) | verified | Get basic device information, including model, firmware version, hardware version, etc. |
| [`get_load`](#get_load) | - | Get load information |
| [`get_percent`](#get_percent) | - | Get battery level |
| [`get_security_policy`](#get_security_policy) | - | Get router security policy |
| [`get_status`](#get_status) | verified | Get system network, service, resource usage and other status |
| [`get_timezone_config`](#get_timezone_config) | - | Get the current router's timezone information |
| [`get_unixtime`](#get_unixtime) | - | Get Unix timestamp |
| [`reboot`](#reboot) | - | Reboot device |
| [`remove_user`](#remove_user) | - | Delete user |
| [`reset_firmware`](#reset_firmware) | - | Reset to factory |
| [`set_password`](#set_password) | - | Set password |
| [`set_security_policy`](#set_security_policy) | - | Set the router's security policy. |
| [`set_timezone_config`](#set_timezone_config) | - | Set router timezone information [for timezones other than UTC, both parameters are required]. |

---

## add_user

Add user

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `username` | string | yes | Username |
| `password` | string | no | Password |
| `home` | string | no | User home directory |
| `interpreter` | string | no | User command interpreter |
| `uid` | number | no | User ID |
| `gid` | number | no | User group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: uid and gid must be greater than 0; -2: Username already exists; -3: uid already exists; -4: gid already exists) |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "add_user",
    {
      "username": "test",
      "password": "123"
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

## disk_info

Get disk information (unit: bytes B)

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `root` | object | - | Root file system information |
| `root.total` | number | - | Total size |
| `root.free` | number | - | Free size |
| `root.used` | number | - | Used size |
| `tmp` | object | - | Memory file system information |
| `tmp.total` | number | - | Total size |
| `tmp.free` | number | - | Free size |
| `tmp.used` | number | - | Used size |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "disk_info"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "root": {
      "total": 1231,
      "free": 100,
      "used": 1131
    },
    "tmp": {
      "total": 1231,
      "free": 100,
      "used": 1131
    }
  }
}
```

---

## get_httpd_mem_status

Get the memory usage of the HTTP server. When testing for memory leaks, this interface can be used to monitor if the API has memory leaks

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `memory_used` | number | - | Memory usage, unit is KB |
| `err_code` | number | - | Error code (-1: Failed to get memory usage information or httpd service is not running) |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "get_httpd_mem_status"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "memory_used": 11952
  }
}
```

---

## get_info

Get basic device information, including model, firmware version, hardware version, etc.

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `hardware_version` | string | - | Device hardware version [temporarily none] |
| `firmware_version` | string | - | Device firmware version |
| `firmware_type` | string | - | Type of the device's current firmware version |
| `firmware_date` | string | - | Firmware build date |
| `vendor` | string | - | Device vendor |
| `mac` | string | - | Device MAC address |
| `ddns` | string | - | ddns |
| `model` | string | - | Device model |
| `sn` | string | - | sn |
| `sn_bak` | string | - | sn backup |
| `country_code` | string | - | Country code |
| `hardware_feature` | object | - | Hardware features |
| `hardware_feature.wan` | string | - | WAN port |
| `hardware_feature.lan` | string | - | LAN port |
| `hardware_feature.build_in_modem` | string | - | Built-in modem |
| `hardware_feature.usb` | string | - | USB port |
| `hardware_feature.radio` | string | - | radio [determine if repeater is supported] |
| `hardware_feature.reset_button` | string | - | reset button |
| `hardware_feature.switch_button` | string | - | switch button |
| `hardware_feature.microsd` | string | - | SD card |
| `hardware_feature.mcu` | bool | - | Whether MCU is supported [detect battery level and temperature] |
| `hardware_feature.fan` | bool | - | Whether fan is supported [detect fan status and temperature] |
| `hardware_feature.nand` | bool | - | Whether NAND is supported |
| `software_feature` | object | - | Software features |
| `software_feature.ipv6` | bool | - | Whether IPv6 service is supported |
| `software_feature.vpn` | bool | - | Whether VPN service is supported |
| `software_feature.adguard` | bool | - | Whether AdGuard service is supported |
| `software_feature.tor` | bool | - | Whether Tor service is supported |
| `software_feature.ids_ips` | bool | - | Whether IDS/IPS service is supported |
| `disable_guest_during_scan_wifi` | bool | - | Whether to disable guest WiFi during WiFi scan |
| `board_info` | object | - | Board-level properties |
| `board_info.architecture` | string | - | CPU architecture |
| `board_info.hostname` | string | - | hostname |
| `board_info.kernel_version` | string | - | Kernel version |
| `board_info.openwrt_version` | string | - | OpenWrt version |
| `board_info.model` | string | - | OpenWrt board model |
| `err_code` | number | - | Error code,-1: Failed to get information. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "get_info"
  ]
}
```

**Response**

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "mac": "94:83:C4:0C:74:9A",
    "hardware_version": "",
    "software_feature": {
      "ids_ips": false,
      "ipv6": true,
      "adguard": false,
      "tor": false,
      "vpn": true
    },
    "vendor": "GL.iNet",
    "hardware_feature": {
      "switch_button": "",
      "radio": "radio0",
      "lan": "eth0",
      "nand": true,
      "usb": "1-1.3",
      "build_in_modem": "1-1.2",
      "microsd": "1-1.1",
      "wan": "eth1",
      "mcu": true,
      "fan": true,
      "reset_button": "gpio-3"
    },
    "country_code": "",
    "sn_bak": "705e8c25242eddd6",
    "board_info": {
      "architecture": "ARMv7 Processor rev 4 (v7l)",
      "hostname": "GL-AXT1800",
      "kernel_version": "4.4.60",
      "openwrt_version": "OpenWrt 21.02-SNAPSHOT r16273+114-378769b555",
      "model": "GL Technologies, Inc. AXT1800"
    },
    "firmware_date": "2022-02-16 5:42:18",
    "model": "xe300",
    "ddns": "evc749a",
    "sn": "4078493df9a88098",
    "firmware_type": "alpha4",
    "firmware_version": "4.0.0"
  }
}
```

---

## get_load

Get load information

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `load_average` | array | - | System load situation, fixed return an array containing three decimals, using decimals to represent CPU usage rate, sequentially indicating the average load values for 1/5/15 minutes |
| `memory_total` | number | - | The system's total memory, unit is byte |
| `memory_free` | number | - | The system's free memory, unit is byte |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "get_load"
  ]
}
```

**Response**

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "load_average": [
      2.01,
      0.89,
      0.33
    ],
    "flash_free": 105918464,
    "flash_total": 106278912
  }
}
```

---

## get_percent

Get battery level

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `percent` | number | - | Battery percentage. |
| `err_code` | number | - | Error code, -1: Failed to get information. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "get_percent"
  ]
}
```

**Response**

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "percent": 100
  }
}
```

---

## get_security_policy

Get router security policy

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `redirect_https` | number | - | Redirect switch |
| `security_rule` | number | - | Security regulation options [0: No security regulation requirements; 1: General security regulation requirements] |
| `err_code` | number | - | Error code, -1: Failed to get information. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "get_security_policy"
  ]
}
```

**Response**

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "redirect_https": 0,
    "security_rule": 0
  }
}
```

---

## get_status

Get system network, service, resource usage and other status

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `system` | object | - | Contains basic system status information |
| `system.mode` | number | - | System operating mode (0: router mode, 1: WDS bridge mode, 2: relay bridge mode, 3: mesh mode, 4: AP mode) [default is 0, returns 5 if not one of the above modes] |
| `system.timestamp` | number | - | System timestamp |
| `system.tzoffset` | string | - | Router's current timezone offset |
| `system.uptime` | number | - | System uptime (seconds) |
| `system.load_average` | array | - | System load, fixed return of an array with three decimals, representing CPU usage, in order of 1/5/15 minute average load values |
| `system.memory_total` | number | - | System total memory, unit is byte |
| `system.memory_free` | number | - | System free memory, unit is byte |
| `system.flash_total` | number | - | System total disk space, unit is byte |
| `system.flash_free` | number | - | System unused disk space, unit is byte |
| `system.flash_app` | number | - | Disk space used by system APP, unit is byte |
| `system.lan_ip` | string | - | LAN IP |
| `system.guest_ip` | string | - | GUEST IP |
| `system.ipv6_enabled` | bool | - | Whether the system enables IPV6 |
| `system.ddns_enabled` | bool | - | Whether the system enables ddns |
| `system.mcu` | object | - | mcu. |
| `system.mcu.temperature` | number | - | Battery gauge temperature. |
| `system.mcu.charge_percent` | number | - | Battery percentage. |
| `system.mcu.charging_status` | number | - | Charging status. [0 is not charging, 1 is charging] |
| `system.mcu.charge_cnt` | number | - | Battery charge cycle count. |
| `system.cpu.temperature` | number | - | CPU core temperature. |
| `network` | array | - | Contains system network connectivity status |
| `network.interface` | string | - | Network interface name |
| `network.up` | bool | - | Whether the network started successfully |
| `network.online` | bool | - | Whether external network is accessible |
| `network.balance` | number | - | Pre-configured network usage rate, represented as a decimal, 1 means 100% usage, 0 means unused |
| `service` | array | - | Contains system service usage, mainly including VPN, AdguardHome and other security services enabled status [returns an empty name if the service is not present] |
| `service.name` | string | - | Service name |
| `service.status` | number | - | Service enabled status [0 not enabled \| 1 connected successfully \| 2 enabled but connection not successful] |
| `service.group_id` | number | - | ovpnclient or wgclient group id |
| `service.client_id` | number | - | ovpnclient's client_id |
| `service.peer_id` | number | - | wgclient's peer_id |
| `wifi` | array | - | Contains basic system WiFi status |
| `wifi.ssid` | string | - | WiFi SSID |
| `wifi.name` | string | - | Interface name, unique identifier for WiFi |
| `wifi.band` | string | - | WiFi band ("2G" or "5G") |
| `wifi.channel` | number | - | WiFi channel (0 means auto) |
| `wifi.guest` | bool | - | Whether it belongs to guest network |
| `wifi.up` | bool | - | Whether normally enabled |
| `wifi.passwd` | string | - | WiFi password |
| `client` | array | - | Total online client status |
| `client.cable_total` | number | - | Total online wired clients. |
| `client.wireless_total` | number | - | Total online wireless clients. |
| `err_code` | number | - | Error code,-1: Failed to get information. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "get_status"
  ]
}
```

**Response**

```json
{"id":1,"jsonrpc":"2.0","result":{"network":[{"online":false,"up":false,"interface":"wan"},{"online":false,"up":false,"interface":"wwan"},{"online":false,"up":false,"interface":"tethering"},{"online":false,"up":false,"interface":"modem_1_1_2"}],"wifi":[{"guest":false,"ssid":"GL-XE300-49a","up":true,"channel":0,"band":"2G","name":"default_radio0","passwd":"goodlife"},{"guest":true,"ssid":"GL-XE300-49a-Guest","up":false,"channel":0,"band":"2G","name":"guest2g","passwd":"goodlife"}],"service":[{"name":"wgclient","status":0},{"name":"wgserver","status":0},{"name":"ovpnclient","status":0},{"name":"ovpnserver","status":0}],"client":[{"wireless_total":0,"cable_total":1}],"system":{"lan_ip":"192.168.8.1","guest_ip":"192.168.9.1","flash_total":106278912,"memory_total":126943232,"memory_free":78471168,"ipv6_enabled":false,"ddns_enabled":false,"uptime":111,"load_average":[2.01,0.89,0.33],"flash_free":105918464i,"flash_app":1025,"mode":0,"mcu":{"charge_cnt":49,"temperature":27.6,"charge_percent":100,"charging_status":1},"cpu":{"temperature":82},"timestamp":1613402877}}}
```

---

## get_timezone_config

Get the current router's timezone information

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `zonename` | string | - | Router's current timezone name [zonename will return empty when timezone is UTC] |
| `timezone` | string | - | Router's current timezone |
| `tzoffset` | string | - | Router's current timezone offset |
| `localtime` | number | - | Router's current system timestamp |
| `auto_timezone_enabled` | bool | - | Indicates whether the router has auto timezone enabled. |
| `err_code` | number | - | Error code,-1: Failed to get information. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "get_timezone_config"
  ]
}
```

**Response**

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "zonename": "Asia/Shanghai",
    "tzoffset": "+0800",
    "autotimezone_enabled": true,
    "localtime": 1643200134,
    "timezone": "CST-8"
  }
}
```

---

## get_unixtime

Get Unix timestamp

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `time` | number | - | Unix timestamp |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "get_unixtime"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "time": 1629106981
  }
}
```

---

## reboot

Reboot device

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `delay` | number | yes | Delay reboot time (in seconds) |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "reboot"
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

## remove_user

Delete user

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `username` | string | yes | Username |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "remove_user",
    {
      "username": "test"
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

## reset_firmware

Reset to factory

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
    "system",
    "reset_firmware"
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

## set_password

Set password

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `username` | string | yes | Username |
| `old_password` | string | no | Old password |
| `new_password` | string | yes | New password |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: No permission, -2: Command not found, -3: Other unknown errors) |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "system",
    "set_password",
    {
      "username": "admin",
      "old_password": "123",
      "new_password": "456"
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

## set_security_policy

Set the router's security policy.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `redirect_https` | number | yes | Whether to redirect |
| `security_rule` | number | yes | Specify the router's security regulation options. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1: Invalid parameter. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "system",
    "set_security_policy",
    {
      "redirect_https": 1,
      "security_rule": 1
    }
  ],
  "id": 1
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

## set_timezone_config

Set router timezone information [for timezones other than UTC, both parameters are required].

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `zonename` | string | yes | Specify the timezone name to set for the router. |
| `timezone` | string | yes | Specify the timezone to set for the router. |
| `localtime` | number | yes | Specify the system timestamp to set for the router |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1: Invalid parameter. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "system",
    "set_timezone_config",
    {
      "zonename": "",
      "timezone": "UTC",
      "localtime": 1653622981
    }
  ],
  "id": 1
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
