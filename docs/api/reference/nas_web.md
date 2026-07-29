# `nas_web`

This is the API related to nas web config.

17 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","nas_web","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_share`](#add_share) | - | add one share |
| [`add_user`](#add_user) | - | add samba share user |
| [`eject_disk`](#eject_disk) | - | eject one disk |
| [`get_disk_list`](#get_disk_list) | - | get nas disk info |
| [`get_file_list`](#get_file_list) | - | get file list |
| [`get_nas_ser`](#get_nas_ser) | - | get nas service info |
| [`get_proto_config`](#get_proto_config) | - | get nas proto conf |
| [`get_share_list`](#get_share_list) | - | get share list |
| [`get_status`](#get_status) | - | get nas service info |
| [`get_user_list`](#get_user_list) | - | get share user list |
| [`remove_share`](#remove_share) | - | remove one share |
| [`remove_user`](#remove_user) | - | remove samba share user |
| [`set_nas_ser`](#set_nas_ser) | - | set nas ser |
| [`set_proto_config`](#set_proto_config) | - | set nas proto conf |
| [`set_share`](#set_share) | - | set share |
| [`set_user_pwd`](#set_user_pwd) | - | set samba share user pwd |
| [`start`](#start) | - | start nas ser |

---

## add_share

add one share

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `file` | string | yes | share path |
| `proto` | string | yes | share proto |
| `users.name` | string | yes | user name |
| `users` | array | yes | file share user |
| `readonly` | number | yes | read only |
| `public` | number | yes | public |
| `users.readonly` | number | yes | user is readonly |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"add_share\",{\"file\":\"/disk1_part1\",\"proto\":\"samba\",\"readonly\":1,\"public\":1,\"users\":{\"name\":\"xxx\",\"readonly\":1}}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## add_user

add samba share user

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | user name |
| `password` | string | yes | user password |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"add_user\",{\"name\":\"test\",\"password\":123456}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## eject_disk

eject one disk

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dev_name` | string | yes | dev name |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"eject_disk\",{\"dev_name\":\"sda1\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## get_disk_list

get nas disk info

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `disk_number` | number | - | disk number |
| `disk` | array | - | disk array |
| `disk.part_num` | number | - | disk part number |
| `disk.sd_card` | number | - | disk is sd card |
| `disk.part` | array | - | disk part array |
| `part.dev_name` | string | - | part dev name |
| `part.disk_name` | string | - | part disk name |
| `part.label` | string | - | part disk label |
| `part.uid` | string | - | part disk luid |
| `part.total_len` | number | - | disk total size |
| `part.free_size` | number | - | disk free size |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"get_disk_list\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"disk_number\":1,\"disk\":[{\"part_num\":4,\"sd_card\":1,\"part\":[{\"dev_name\":\"mmcblk0p11\",\"disk_name\":\"disk1_part1\",\"fs_type\":\"EXFAT\",\"uid\":\"7005-A5B5\",\"label\":\"exFAT_131G\",\"total_len\":134663,\"free_size\":5351},{\"dev_name\":\"mmcblk0p6\",\"disk_name\":\"disk1_part2\",\"fs_type\":\"FAT\",\"uid\":\"D4E5-0DF1\",\"label\":\"FAT32_1G\",\"total_len\":1026,\"free_size\":1014},{\"dev_name\":\"mmcblk0p7\",\"disk_name\":\"disk1_part3\",\"fs_type\":\"EXFAT\",\"uid\":\"9784-6325\",\"label\":\"exFAT_1G\",\"total_len\":1024,\"free_size\":354},{\"dev_name\":\"mmcblk0p8\",\"disk_name\":\"disk1_part4\",\"fs_type\":\"EXT4\",\"uid\":\"a56d2330-dee8-4b4a-8912-cdc90be2061a\",\"label\":\"EXT4_1G\",\"total_len\":1023,\"free_size\":953}]}]}}
```

---

## get_file_list

get file list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `result_code` | number | - | result code |
| `files` | array | - | file list array |
| `files.n` | string | - | file name |
| `files.d` | number | - | file date |
| `files.l` | number | - | file length |
| `files.t` | string | - | file type |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"get_file_list\",{\"path\":\"/disk1_part1\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{"files":[{\"n\":\"/disk1_part1\",\"d\":123,\"l\":6009,\"t\":\"d\",\"have_dir\":1}]}}
```

---

## get_nas_ser

get nas service info

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enable` | number | - | nas ser if enable |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"get_nas_ser\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"res\":[{\"enable\":1,\"port\":6009}]}}
```

---

## get_proto_config

get nas proto conf

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `res` | array | - | all proto config array |
| `res.name` | string | - | proto name |
| `res.enable` | bool | - | proto if enable |
| `res.port` | number | - | listen port |
| `res.wan_access` | bool | - | if wan access |
| `res.ssl_enable` | bool | - | webdav ssl enable |
| `res.media_dir` | string | - | dlna media dir |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"get_proto_config\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"res\":[{\"name\":\"webdav\",\"enable\":true,\"port\":6009,\"wan_access\":true}]}}
```

---

## get_share_list

get share list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `result_code` | number | - | result code |
| `share` | array | - | share array |
| `share.n` | string | - | share path |
| `share.share_id` | string | - | share id |
| `share.file_ok` | number | - | share file if ok |
| `share.protos` | string | - | protos |
| `protos.name` | string | - | proto name |
| `protos.enable` | number | - | proto share enable |
| `protos.public` | number | - | proto is public |
| `protos.public_readonly` | number | - | proto is public readonly |
| `protos.share_name` | string | - | proto share file name |
| `protos.users` | array | - | proto user list |
| `users.name` | string | - | user name |
| `users.readonly` | string | - | user is readonly |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"get_share_list\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"result_code\":0,\"share\":[{\"n\":\"/disk1_part4/lost+found\",\"owner_readonly\":0,\"share_id\":\"bf90e29f2e3ab0091b1d922e81294dba\",\"t\":\"d\",\"share_time\":1653547578,\"d\":1650649376,\"l\":4096,\"file_ok\":1,\"disk_uid\":\"a56d2330-dee8-4b4a-8912-cdc90be2061a\",\"protos\":[{\"name\":\"webdav\",\"enable\":1,\"public\":1,\"public_readonly\":1,\"share_name\":\"lost+found\",\"users\":[]}]},{\"n\":\"/disk1_part2/www\",\"owner_readonly\":0,\"share_id\":\"70334cc85ab58f13363f8a7d906cef22\",\"t\":\"d\",\"share_time\":1653547602,\"d\":1653548241,\"l\":4096,\"file_ok\":1,\"disk_uid\":\"D4E5-0DF1\",\"protos\":[{\"name\":\"webdav\",\"enable\":1,\"public\":1,\"public_readonly\":0,\"share_name\":\"www\",\"users\":[]}]},{\"n\":\"/disk1_part2/test\",\"owner_readonly\":0,\"share_id\":\"59297d2089f011a05f84dfc3a261c979\",\"t\":\"d\",\"share_time\":1653548194,\"d\":1653477986,\"l\":4096,\"file_ok\":1,\"disk_uid\":\"D4E5-0DF1\",\"protos\":[{\"name\":\"webdav\",\"enable\":1,\"public\":1,\"public_readonly\":1,\"share_name\":\"test\",\"users\":[]}]}]}}
```

---

## get_status

get nas service info

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `enable` | number | - | nas ser is enable |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"get_status\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"res\":[{\"enable\":1}]}}
```

---

## get_user_list

get share user list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `list` | array | - | user list |
| `list.name` | string | - | user name |
| `list.password` | string | - | user password |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"get_user_list\",{}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{"list":[{\"name\":\"test\",\"password\":\"123\"}]}}
```

---

## remove_share

remove one share

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `share_id` | string | yes | share id |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"remove_share\",{\"share_id\":\"xxx\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## remove_user

remove samba share user

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | proto name |
| `all` | string | no | remove all user |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"remove_user\",{\"name\":\"test1\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## set_nas_ser

set nas ser

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `enable` | number | yes | if enable nas ser |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"set_nas_ser\",{\"enable\":1}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## set_proto_config

set nas proto conf

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | proto name |
| `media_dir` | string | no | dlna media dir |
| `enable` | bool | yes | proto enable |
| `wan_access` | bool | yes | if wan access |
| `ssl_enable` | bool | no | webdav ssl enable |
| `port` | number | yes | listen port |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"set_proto_config\",{\"name\":\"webdav\",\"enable\":true,\"port\":6009,\"wan_access\":true}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## set_share

set share

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `share_id` | string | yes | share id |
| `proto` | string | yes | share proto name |
| `users.name` | string | yes | user name |
| `users` | array | yes | file share user |
| `readonly` | number | yes | read only |
| `public` | number | yes | public |
| `users.readonly` | number | yes | user is readonly |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"add_share\",{\"share_id\":\"xxx\",\"proto\":\"samba\",\"readonly\":1,\"public\":1,\"users\":{\"name\":\"xxx\",\"readonly\":1}}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## set_user_pwd

set samba share user pwd

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | user name |
| `password` | string | yes | user password |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"set_user_pwd\",{\"name\":\"test\",\"password\":123456}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## start

start nas ser

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"nas-web\",\"start\"]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---
