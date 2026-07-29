# `led`

This is the API related to led Internet access.

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","led","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | verified | Get config of LED. |
| [`set_config`](#set_config) | - | Set led config. |

---

## get_config

Get config of LED.

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `led_enable` | bool | - | LED status, true: on false: off. |
| `timer_enable` | bool | - | LED timer status, true: on false: off. |
| `turnon_hour` | string | - | LED timer set LED turn-on time (hour) |
| `turnon_min` | string | - | LED timer set LED turn-on time (minute) |
| `turnoff_hour` | string | - | LED timer set LED turn-off time (hour) |
| `turnoff_min` | string | - | LED timer set LED turn-off time (minute) |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"led\",\"get_config\",{}],\"id\":1}
```

**Response**

```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "turnon_hour": "07",
    "turnoff_min": "00",
    "turnon_min": "00",
    "led_enable": true,
    "timer_enable": true,
    "turnoff_hour": "22"
  }
}
```

---

## set_config

Set led config.

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `turnon_hour` | string | yes | LED timer set LED turn-on time (hour) |
| `turnon_min` | string | yes | LED timer set LED turn-on time (minute) |
| `turnoff_hour` | string | yes | LED timer set LED turn-off time (hour) |
| `turnoff_min` | string | yes | LED timer set LED turn-off time (minute) |
| `led_enable` | bool | yes | LED status, true: on false: off. |
| `timer_enable` | bool | yes | LED timer status, true: on false: off. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code,-1;parameter error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"led\",\"set_config\",{\"led_enable\":true,\"timer_enable\":true,\"turnon_hour\":\"07\",\"turnon_min\":\"00\",\"turnoff_hour\":\"22\",\"turnoff_min\":\"00\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---
