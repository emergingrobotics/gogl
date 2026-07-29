# `logread`

logread

8 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","logread","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`export_logs`](#export_logs) | - | Export logs |
| [`get_config`](#get_config) | - | Read log configuration [The record_size parameter is for setting the crash log function, enable and path parameters are for setting the local log saving function. Here, two situations share one interface] |
| [`get_crash_log`](#get_crash_log) | - | Get crash log [crash log], default return all logs if no parameters |
| [`get_kernel_log`](#get_kernel_log) | - | Get kernel log |
| [`get_system_log`](#get_system_log) | - | Get system log |
| [`get_uboot_log`](#get_uboot_log) | - | Get uboot log |
| [`remove_crash_log`](#remove_crash_log) | - | Delete crash log |
| [`set_config`](#set_config) | - | Configure logs [The record_size parameter is for setting the crash log function, enable and path parameters are for setting the local log saving function. Here, two situations share one interface, pass parameters as needed] |

---

## export_logs

Export logs

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `file_name` | string | - | File name |
| `file_path` | string | - | File path |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"logread\",\"export_logs\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"file_name\":\"logread.tar\",\"file_path\":\"/js/logread.tar\"}}
```

---

## get_config

Read log configuration [The record_size parameter is for setting the crash log function, enable and path parameters are for setting the local log saving function. Here, two situations share one interface]

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `record_size` | number | no | Single crash log size [default 8192, minimum 4096, this parameter is an integer multiple of 4096] |

| Results | Type | Required | Description |
|---|---|---|---|
| `enable` | bool | - | Enable or disable openwrt local log saving |
| `path` | string | - | openwrt local log save path |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"logread\",\"get_config\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"record_size\":8192,\"path\":\"/usr/share/mylog\",\"enable\":true}}
```

---

## get_crash_log

Get crash log [crash log], default return all logs if no parameters

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mode` | number | no | Set the log mode to retrieve [0 or 1] (0 is select the latest log \| 1 is select log sequence, return the nth log, in this mode log_number parameter must have value) |
| `log_number` | number | no | Set the log sequence to retrieve [when mode is 1, log_number parameter must have value] |

| Results | Type | Required | Description |
|---|---|---|---|
| `log` | string | - | Output the latest crash log [default return all logs if no parameters, with parameters only output the specified log] |
| `sum` | number | - | Output the total number of currently stored crash logs |
| `err_code` | number | - | Error code, -1: Entered wrong log sequence. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"logread\",\"get_crash_log\",{\"mode\":1,\"log_number\":2}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"log\":\"Log Entry 2 (at position 1)\n<6>[  650.605535] br-lan: port 2(wlan0) entered forwarding state\n\n\",\"sum\":3}}
```

---

## get_kernel_log

Get kernel log

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `log` | string | - | Output kernel log [return all] |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"logread\",\"get_kernel_log\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"log\":\"[    0.000000] Linux version 4.14.221 (glinet@ubuntu) (gcc version 7.5.0 (OpenWrt GCC 7.5.0 r11306-c4a6851c72)) #0 Mon Feb 15 15:22:37 2021\n\"}}
```

---

## get_system_log

Get system log

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `module` | string | no | Set the log module name to retrieve [specific module names need to be checked in the backend, if not set, no related filtering is performed] |
| `lines` | number | no | Set the number of log lines to retrieve [if the number of log lines is not set, return all matching logs] |

| Results | Type | Required | Description |
|---|---|---|---|
| `log` | string | - | Output system log [default return all logs if no parameters, return only if log exists] |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"logread\",\"get_system_log\",{\"lines\":1,\"module\":\"oui-httpd\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"log\":\"Mon Feb 15 18:48:34 2021 daemon.debug oui-httpd: (uhttpd.c:461) Listen on: [::]:443 with ssl\n\"}}
```

---

## get_uboot_log

Get uboot log

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `log` | string | - | Output uboot log [return all] |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"logread\",\"get_uboot_log\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"log\":\"uboot.info: uboot entering firmware upgrade model\n\"}}
```

---

## remove_crash_log

Delete crash log

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"logread\",\"remove_crash_log\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## set_config

Configure logs [The record_size parameter is for setting the crash log function, enable and path parameters are for setting the local log saving function. Here, two situations share one interface, pass parameters as needed]

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `path` | string | no | openwrt local log save path [if not passed, default path /usr/share/mylog] |
| `enable` | bool | yes | Enable or disable openwrt local log saving |
| `record_size` | number | no | Single crash log size [default 8192, minimum 4096, this parameter is an integer multiple of 4096] |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"logread\",\"set_config\",{\"record_size\":4096,\"enable\":true,\"path\":\"/usr/share/mylog\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":null}
```

---
