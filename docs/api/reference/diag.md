# `diag`

Diag

2 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","diag","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`ping`](#ping) | - | Network detection |
| [`traceroute`](#traceroute) | - | Trace network path |

---

## ping

Network detection

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `addr` | string | yes | Target address |

| Results | Type | Required | Description |
|---|---|---|---|
| `ping_result` | string | - | Ping result |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "diag",
    "ping",
    {
      "addr": "192.168.8.1"
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
    "ping_result": "string"
  }
}
```

---

## traceroute

Trace network path

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `addr` | string | yes | Target address |

| Results | Type | Required | Description |
|---|---|---|---|
| `trace_result` | string | - | Traceroute result |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "diag",
    "traceroute",
    {
      "addr": "192.168.8.1"
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
    "trace_result": "string"
  }
}
```

---
