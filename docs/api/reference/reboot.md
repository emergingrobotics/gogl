# `reboot`

This is the API related to reboot Internet access.

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","reboot","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | - | Get config of reboot. |
| [`set_config`](#set_config) | - | Set reboot config. |

---

## get_config

Get config of reboot.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `hour` | string | yes | Restart time (hour) |
| `min` | string | yes | Restart time (minute) |
| `enable` | bool | yes | Scheduled restart status, true: enable false: disable. |
| `week` | array | yes | Which days in a week to perform restart operations |

_No results._

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"reboot\",\"get_config\",{}],\"id\":1}
```

**Response**

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "enable": false,
    "week": [
      0,
      1,
      2,
      3,
      4,
      5,
      6
    ],
    "hour": "07",
    "min": "00"
  }
}
```

---

## set_config

Set reboot config.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `hour` | string | yes | Restart time (hour) |
| `min` | string | yes | Restart time (minute) |
| `enable` | bool | yes | Scheduled restart status, true: enable false: disable. |
| `week` | array | yes | Which days in a week to perform restart operations |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1;parameter error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"reboot\",\"set_config\",{\"enable\":true,\"hour\":\"03\",\"min\":\"00\",\"week\":[0,1,2,3,4,5,6]}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---
