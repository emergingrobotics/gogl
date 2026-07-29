# `plugins`

This is the api related to installing the package.

6 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","plugins","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_list`](#get_list) | - | Get package list |
| [`get_package_info`](#get_package_info) | - | Get detailed information of a single package |
| [`get_repository_status`](#get_repository_status) | - | Get software repository status |
| [`install_package`](#install_package) | - | Install package |
| [`remove_package`](#remove_package) | - | Uninstall software package |
| [`update_repository`](#update_repository) | - | Update software repository |

---

## get_list

Get package list

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `search_condition` | string | no | Search condition, used when searching packages, regular expression like *test*, no need to add * in input |
| `search_initials` | string | no | Search condition, used when searching packages, regular expression like a*, no need to add * in input, for filtering packages starting with initial letter, these two search conditions are mutually exclusive |
| `status` | number | yes | Get package list of specified status, 0: not installed, 1: installed, 2: updatable, 3: user reserved, 4: all |
| `limit` | number | yes | Maximum number per page |
| `page` | number | yes | Page number, set to 1 when searching, set to specific page when paging |

| Results | Type | Required | Description |
|---|---|---|---|
| `packages` | array | - | Package information |
| `packages.name` | string | - | Package name |
| `packages.status` | number | - | Package status, 0: not installed, 1: installed, 2: updatable, 3: user reserved |
| `packages.uninstallable` | bool | - | Whether uninstall is allowed, true: allowed, false: not allowed (packages that affect system operation) |
| `packages.version` | string | - | Package version, not returned when repository not updated and status is user reserved |
| `packages.new_version` | string | - | New package version, returned when package status is updatable |
| `packages.size` | number | - | Package size, returned after repository has been updated |
| `count_search` | number | - | Number of packages matching search condition, returned when searching |
| `err_code` | number | - | Error code,-7: database error |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"plugins\",\"get_list\",{\"status\":1,\"limit\":3,"page\":1}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"packages\":[{\"name\":\"464xlat\",\"version\":\"12\",\"size\":4919,\"status\":1},{\"name\":\"zram-swap\",\"version\":\"1.1-3\",\"size\":2890,\"status\":1},{\"name\":\"zstd\",\"version\":\"1.4.5-2\",\"size\":47101,\"status\":1}]}}
```

---

## get_package_info

Get detailed information of a single package

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | The name of software. |

| Results | Type | Required | Description |
|---|---|---|---|
| `info` | string | - | opkg command execution result |
| `err_code` | number | - | Error code,-2:network unreachable,-3:flash not enough storage space.-6:opkg stderr. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"plugins\",\"get_package_info\",{\"name\":\"kmod-usb-serial-pl2303\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"info\":\"Package: kmod-usb-serial-pl2303\nVersion: 4.14.221-1\nDepends: kernel (= 4.14.221-1-0a66574dfa39734270ca25227d9a8ac6), kmod-usb-serial\nStatus: unknown ok not-installed\nSection: kernel\nArchitecture: mips_24kc\nSize: 5798\nFilename: kmod-usb-serial-pl2303_4.14.221-1_mips_24kc.ipk\nDescription: Kernel support for Prolific PL2303 USB-to-Serial converters\n\n\"}}
```

---

## get_repository_status

Get software repository status

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `status` | number | - | Software repository status, 0: not updated, 1: updating, 2: updated. |
| `time_last_update` | number | - | Timestamp of the last update, 0 means never updated |
| `time_current` | number | - | Current timestamp |
| `count_all` | number | - | Total number of packages (counts are not returned during update) |
| `count_installed` | number | - | Number of installed packages |
| `count_not_installed` | number | - | Number of not installed packages |
| `count_updatable` | number | - | Number of updatable packages |
| `count_reserve` | number | - | Number of user-reserved packages, returned after upgrading user-reserved packages |
| `err_code` | number | - | Error code,-7: database error |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"plugins\",\"get_repository_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"status\":2,\"count_all\":7220,\"count_installed\":200,\"count_not_installed\":7000,\"count_updatable\":20,\"count_reserve\":0}}
```

---

## install_package

Install package

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | array | yes | Package name to install, batch install can be used when installing user-reserved packages |

| Results | Type | Required | Description |
|---|---|---|---|
| `info` | string | - | opkg command execution result |
| `err_code` | number | - | Error code,-1:opkg resource temporarily unavailable,-2:network unreachable,-3:flash not enough storage space.-6:opkg stderr. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"plugins\",\"install_package\",{\"name\":\["tcpliveplay\"]}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"info\":\"Configuring tcpliveplay.\"}}
```

---

## remove_package

Uninstall software package

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Package name. |
| `force` | bool | yes | Whether to uninstall dependent packages |

| Results | Type | Required | Description |
|---|---|---|---|
| `info` | string | - | opkg command execution result |
| `err_code` | number | - | Error code,-1:opkg resource temporarily unavailable,-4:parameter error,-5:can't uninstall.-6:opkg stderr. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"plugins\",\"remove_package\",{\"name\":\"tcpliveplay\",\"force\":false}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"info\":\"Removing package tcpliveplay from root...\"}}
```

---

## update_repository

Update software repository

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1:opkg resource temporarily unavailable,-2:network unreachable. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"plugins\",\"update_repository\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---
