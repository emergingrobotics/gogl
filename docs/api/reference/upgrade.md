# `upgrade`

This is the api related to firmware upgrade.

7 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","upgrade","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`check_firmware_local`](#check_firmware_local) | - | local firmware verify |
| [`check_firmware_online`](#check_firmware_online) | - | Get firmware information on the server. |
| [`get_config`](#get_config) | verified | Get upgrade configuration. |
| [`get_online_upgrade_status`](#get_online_upgrade_status) | - | Get online upgrade status. |
| [`set_config`](#set_config) | - | Set upgrade configuration file |
| [`upgrade_local`](#upgrade_local) | - | Local firmware upgrade. |
| [`upgrade_online`](#upgrade_online) | - | Online one-click upgrade |

---

## check_firmware_local

local firmware verify

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `version` | string | - | Uploaded firmware version (value is unknown, or normal version information) |
| `status` | number | - | Firmware status, 0: normal, 1: uploaded firmware version lower than current, 2: major kernel version upgrade, suggest not keeping configuration, 3: unknown version, third-party firmware or version too low, 4: uploaded firmware too large, 5: firmware does not exist, 6: sysupgrade -T check failed, 7: version not supported, 8: hardware mismatch, 9: firmware signature verification failed (UI warns this firmware is not released by GL), 10: unknown error. |
| `release_note` | string | - | Firmware change log |
| `sha256` | string | - | Firmware hash value |
| `not_keep_config` | string | - | Parse the reason for forcing not to keep configuration from meta data, prompt for forced upgrade without keeping configuration and its reason |
| `not_keep_config.version` | string | - | If current version is lower than this, force upgrade without keeping configuration |
| `not_keep_config.description` | string | - | Reason for forced upgrade without keeping configuration |
| `not_keep_config_part` | array | - | Parse the reason for not keeping partial configurations from meta data, prompt the paths of configurations not kept and their reasons |
| `not_keep_config_part.path` | string | - | Path of configuration file not kept |
| `not_keep_config_part.description` | string | - | Reason for not keeping the configuration file |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"upgrade\",\"check_firmware_local\",{}],\"id\":1 }
```

**Response**

```json
{\"id\":1,\"jsonrpc\":\"2.0\",\"result\":{\"status\":9,\"release_note\":\"###  System：\n* Based on openwrt 19.07.7  (AR750,AR750S,X750,XE300,MT1300,MT300N-V2)\n\n###  SDK\n* Refactoring the API Framework\n* Refactoring the UI Framework\n\n### Change log\n* WIFI，客户端，防火墙及升级页面细节优化\n* AXT1800目标兼容（目前暂未编译）\n* 修复部分modem及repeater后台问题\n* 完善部分语言翻译\n\",\"sha256\":\"264ce49c1f63e19960f5aa5e06b45d05db219ea560011504e6ca83fee5980500\",\"version\":\"4.0.0\"}}
```

---

## check_firmware_online

Get firmware information on the server.

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `current_version` | string | - | Current firmware version |
| `current_compile_time` | string | - | Compile time of the current version firmware |
| `firmware_type` | string | - | Type of the current version firmware |
| `version_new` | string | - | New firmware version, returned when network is reachable and new version exists on server |
| `new_compile_time` | string | - | Compile time of the new version firmware, returned when network is reachable and new version exists on server |
| `release_note` | string | - | Firmware change log, returned when network is reachable and new version exists on server |
| `upgrade_control` | object | - | Upgrade_control information from firmware meta data, returned when network is reachable and new version exists on server |
| `upgrade_control.not_keep_config` | object | - | Parse the reason for forcing not to keep configuration from meta data, prompt for forced upgrade without keeping configuration and its reason |
| `upgrade_control.not_keep_config.version` | string | - | For versions lower than this, force not to keep configuration, and prompt the corresponding description information |
| `upgrade_control.not_keep_config.description` | string | - | Reason for forced upgrade without keeping configuration |
| `upgrade_control.supported_version` | string | - | Minimum supported version, do not upgrade if lower than this, and prompt version not supported |
| `upgrade_control.not_keep_config_part` | array | - | Some configurations not kept, parse the reason for not keeping partial configurations from meta data, prompt the paths of configurations not kept and their reasons |
| `upgrade_control.not_keep_config_part.path` | string | - | Absolute path of the configuration file not kept |
| `upgrade_control.not_keep_config_part.description` | string | - | Reason for not keeping the configuration file |
| `upgrade_control.user` | string | - | For 2c it's general, for 2b it's customer company name, can be used for later permission management |
| `err_code` | number | - | Error code,-2:get url error. |
| `err_msg` | string | - | Error message. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"upgrade\",\"check_firmware_online\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"current_version\":\"4.000\",\"firmware_type\":\"release\",\"current_compile_time\":\"2021-08-11 9:40:45\",\"version_new\":\"4.001\",\"new_compile_time\":\"2021-08-13 09:41:35\",\"release_note\":\"System:Based on openwrt 19.07.7\",\"upgrade_control\":{\"not_keep_config\":{\"version\":\"4.003\",\"description\":\"Solve the xxx problem\"},\"supported_version\":\"4.001\",\"not_keep_config_part\":[{\"path\":\"/etc/config/ssids\",\"description\":\"Solve the xxx problem\"},{\"path\":\"/etc/config/wireless\",\"description\":\"Solve the xxx problem\"}],\"user\":\"general\"}}}
```

---

## get_config

Get upgrade configuration.

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `auto_upgrade` | bool | - | Auto upgrade function enable status. |
| `time` | string | - | Auto upgrade time. |
| `url` | string | - | Firmware upgrade URL. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"upgrade\",\"get_config\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"auto_upgrade\":false,\"time\":\"04:00\",\"url\":\"https://fw.gl-inet.com/firmware/ar750s/release\"}}
```

---

## get_online_upgrade_status

Get online upgrade status.

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `status` | number | - | Status code, 1: Downloading, 2: Download failed, 3: Download complete, 4: Verification passed, 5: Verification failed, 6: Upgrading, 7: Upgrade failed, 9: Firmware signature verification failed (this firmware is not released by GL). |
| `status_msg` | string | - | Firmware upgrade status information |
| `percent` | number | - | Returns download progress when downloading firmware |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"upgrade\",\"get_online_upgrade_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"status\":\"2\",\"status_msg\":\"download_failed\"}}
```

---

## set_config

Set upgrade configuration file

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `time` | string | yes | Auto upgrade time. |
| `url` | string | yes | Firmware upgrade URL. |
| `auto_upgrade` | bool | yes | Auto upgrade function. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1:parameter error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"upgrade\",\"set_config\",{\"auto_upgrade\":false,\"time\":\"04:00\",\"url\":\"https://fw.gl-inet.com/firmware/ar750s/release\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## upgrade_local

Local firmware upgrade.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mesh` | bool | yes | This parameter is required for upgrading devices that support mesh. |
| `keep_config` | bool | yes | True to upgrade keeping configuration. |
| `keep_package` | bool | yes | True to upgrade keeping installed packages. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1:参数错误,-7:内存不足,-8:mesh子节点升级错误,-9:mesh升级错误. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"upgrade\",\"upgrade_local\",{\"mesh\":false,\"keep_config\":false,\"keep_package\":false}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## upgrade_online

Online one-click upgrade

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `keep_config` | bool | yes | If true, keep configuration during upgrade. |
| `keep_package` | bool | yes | If true, keep installation package during upgrade. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -2: Error acquiring firmware connection. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"upgrade\",\"upgrade_online\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---
