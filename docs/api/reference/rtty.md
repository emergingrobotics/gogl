# `rtty`

Rtty Configuration

4 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","rtty","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`get_config`](#get_config) | - | Query current configuration |
| [`run`](#run) | - | Run rtty |
| [`set_config`](#set_config) | - | Set rtty function |
| [`stop`](#stop) | - | Stop rtty |

---

## get_config

Query current configuration

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `ssh_enabled` | bool | - | Whether SSH is allowed |
| `web_enabled` | bool | - | Whether web is allowed |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "rtty",
    "get_config"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "ssh_enabled": true,
    "web_enabled": true
  }
}
```

---

## run

Run rtty

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `token` | string | no | token |
| `host` | string | yes | Server address |
| `ssl` | bool | no | Whether the server has enabled SSL |
| `port` | number | yes | Server port |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code: 1: Device has no valid DDNS ID, 2: rtty connection to server failed |
| `err_msg` | number | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "rtty",
    "run",
    {
      "token": "sfsedtretger124233as",
      "host": "192.168.3.1",
      "port": 5912,
      "ssl": true
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

## set_config

Set rtty function

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `ssh_enable` | bool | yes | Whether SSH is allowed |
| `web_enable` | bool | yes | Whether web is allowed |

| Results | Type | Required | Description |
|---|---|---|---|
| `ssh_enabled` | bool | - | Whether SSH is allowed |
| `web_enabled` | bool | - | Whether web is allowed |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "rtty",
    "set_config",
    {
      "ssh_enable": true,
      "web_enable": true
    }
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "ssh_enabled": true,
    "web_enabled": true
  }
}
```

---

## stop

Stop rtty

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
    "rtty",
    "stop"
  ]
}
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
