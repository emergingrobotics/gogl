# `fan`

This is the API related to fan Internet access.

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","fan","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_status`](#get_status) | - | Get status of fan. |
| [`set_status`](#set_status) | - | Set led status. |

---

## get_status

Get status of fan.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `get_speed` | bool | yes | Whether to get fan speed, true: yes false: no. |

| Results | Type | Required | Description |
|---|---|---|---|
| `fan_speed` | number | - | Fan speed. |
| `fan_status` | bool | - | Fan status, true: on false: off. |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"fan\",\"get_status\",{\"get_speed\":true}],\"id\":1}
```

**Response**

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "fan_speed": 2000,
    "fan_status": true
  }
}
```

---

## set_status

Set led status.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `test_fan` | bool | yes | Test fan startup. |
| `test_time` | number | yes | Test fan startup time in seconds, default is 10s. |

_No results._

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"fan\",\"set_status\",{\"test_fan\":true,\"test_time\":5}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---
