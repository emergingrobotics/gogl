# `cloud`

This is the API related to file share.

3 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","cloud","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | verified | get cloud conf |
| [`set_config`](#set_config) | - | set cloud conf |
| [`unbind`](#unbind) | - | unbind device |

---

## get_config

get cloud conf

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `cloud_enable` | bool | - | enable connect to cloud |
| `rtty_ssh` | bool | - | enable rtty ssh |
| `rtty_web` | bool | - | enable rtty web |
| `serverzone` | string | - | server zone addr |
| `serverzones` | array | - | server zone array |
| `name` | string | - | DDNS |
| `username` | string | - | username |
| `email` | string | - | user email |
| `bindtime` | string | - | user bindtime |
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud\",\"get_config\"]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"cloud_enable\": \"true\",\"rtty_ssh\":\"true\",\"rtty_web\":\"true\",\"serverzone\":\"China\",\"serverzones\":[\"Europe\",\"America\",\"China"],\"name\": \"gclone\",\"email\": \"88666@126.com\",\"bindtime\": \"192168\"}}
```

---

## set_config

set cloud conf

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `serverzone` | string | yes | server zone addr |
| `cloud_enable` | bool | yes | enable connect to cloud |
| `rtty_ssh` | bool | yes | enable rtty ssh |
| `rtty_web` | bool | yes | enable rtty web |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud\",\"set_config\",{\"cloud_enable\": \"true\",\"rtty_ssh\":\"true\",\"rtty_web\":\"true\",\"serverzone\":\"China\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## unbind

unbind device

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `api` | string | - | unbind api |
| `code` | number | - | result code |
| `err_code` | number | - | Error code, 0:ok, -1:Network not reachable; |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud\",\"unbind\"]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"api\":\"cloud/unbind\",\"code\": -1}}
```

---
