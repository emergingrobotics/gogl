# `lan`

This is the API related to wired Internet.

6 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","lan","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_static_bind`](#add_static_bind) | verified | Set static IP binding |
| [`get_config_list`](#get_config_list) | verified | Get lan or guest ip information |
| [`get_static_bind_list`](#get_static_bind_list) | verified | Get the information list of all current static IP binding entries |
| [`remove_static_bind`](#remove_static_bind) | verified | Delete a static IP binding or delete all static IP bindings |
| [`set_config`](#set_config) | verified | Set the IP address of the LAN. |
| [`set_static_bind`](#set_static_bind) | verified | Modify static IP binding |

---

## add_static_bind

Set static IP binding

**Verified** on a GL-SFT1200 running 4.3.28.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | The MAC address of the static IP binding entry to be set. |
| `ip` | string | yes | The IP address of the static IP binding entry to be set. |
| `name` | string | no | The hostname of the device to bind static IP. [Optional parameter, default is "" if no value] |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1: Parameter not found, -3: MAC address format error, -4: UCI initialization error |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "lan",
    "add_static_bind",
    {
      "name": "DESKTOP-HO0T5C1",
      "mac": "58:41:20:0b:9b:b1",
      "ip": "192.168.8.162"
    }
  ],
  "id": 1
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

## get_config_list

Get lan or guest ip information

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `interfaces` | array | - | Interface information array |
| `interfaces.interface` | string | - | Interface name [ lan \| guest ], indicating whether the current information is lan ip information or guest ip information |
| `interfaces.ip` | string | - | The ip of the interface. |
| `interfaces.netmask` | string | - | The netmask of the interface. |
| `interfaces.start` | string | - | The start IP address of the interface. |
| `interfaces.end` | string | - | The end IP address of the interface. |
| `err_code` | number | - | Error code, -1: Parameter not found, -4: UCI initialization error, -5: No LAN interface. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "lan",
    "get_config_list",
    {}
  ],
  "id": 1
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "interfaces": [
      {
        "interface": "lan",
        "ip": "192.168.8.1",
        "protocol": "static",
        "netmask": "255.255.255.0",
        "start": "192.168.8.100",
        "end": "192.168.8.249"
      },
      {
        "interface": "guest",
        "ip": "192.168.9.1",
        "protocol": "static",
        "netmask": "255.255.255.0",
        "start": "192.168.9.100",
        "end": "192.168.9.249"
      }
    ]
  }
}
```

---

## get_static_bind_list

Get the information list of all current static IP binding entries

**Verified** on a GL-SFT1200 running 4.3.28.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `static_bind_list` | array | - | Get the list of all static IP binding entries. |
| `static_bind_list.mac` | string | - | The MAC address of the static IP binding entry. |
| `static_bind_list.ip` | string | - | The IP address of the static IP binding entry. |
| `static_bind_list.name` | string | - | The hostname of the device to bind static IP. [Default is "" if no value] |
| `err_code` | number | - | Error code, -4: UCI initialization error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "lan",
    "get_static_bind_list",
    {}
  ],
  "id": 1
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "static_bind_list": [
      {
        "name": "DESKTOP-HO0T5C1",
        "mac": "58:41:20:0b:9b:b1",
        "ip": "192.168.8.162"
      },
      {
        "name": "DESKTOP-HO0T5C3",
        "mac": "58:41:20:0b:9b:b3",
        "ip": "192.168.8.163"
      }
    ]
  }
}
```

---

## remove_static_bind

Delete a static IP binding or delete all static IP bindings

**Verified** on a GL-SFT1200 running 4.3.28.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | The MAC address of the static IP binding entry to be deleted. Find the corresponding entry based on the MAC address and delete it. [When mode is 0, the MAC address must be passed; when mode is 1, the MAC address is optional (MAC value is empty or MAC parameter is not passed)] |
| `mode` | number | yes | Mode sets the mode for deleting static IP binding entries. Mode 0 is to delete a single static IP binding (requires passing the MAC parameter), mode 1 is to delete all static IP bindings (do not pass the MAC parameter or MAC parameter is empty). |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -4: UCI initialization error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "lan",
    "remove_static_bind",
    {
      "mode": 0,
      "mac": "58:41:20:0b:9b:b1"
    }
  ],
  "id": 1
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

## set_config

Set the IP address of the LAN.

**Verified** on a GL-SFT1200 running 4.3.28.

| Params | Type | Required | Description |
|---|---|---|---|
| `interface` | string | yes | Interface name [ lan \| guest ], indicating whether to set lan ip or guest ip. |
| `ip` | string | yes | Set new IP. |
| `start` | string | yes | Set new start IP address. |
| `end` | string | yes | Set new end IP address. |
| `netmask` | string | yes | Set new netmask. |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1: Parameter not found, -2: IP format error, -4: UCI initialization error, -7: IP conflict. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "lan",
    "set_config",
    {
      "interface": "guest",
      "ip": "192.168.9.1",
      "start": "192.168.9.101",
      "end": "192.168.9.245",
      "netmask": "255.255.255.0"
    }
  ],
  "id": 1
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

## set_static_bind

Modify static IP binding

**Verified** on a GL-SFT1200 running 4.3.28.

| Params | Type | Required | Description |
|---|---|---|---|
| `mac` | string | yes | The MAC address of the static IP binding entry to be modified. [Find the corresponding static IP binding entry based on MAC] |
| `ip` | string | yes | The IP address of the static IP binding entry to be modified. [Set the IP of the static IP binding entry for the corresponding MAC] |
| `name` | string | yes | The name of the static IP binding entry to be modified. [Set the name of the static IP binding entry for the corresponding MAC] |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -4: UCI initialization error. |
| `err_msg` | string | - | Error message |

**Request**

```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": [
    "",
    "lan",
    "set_static_bind",
    {
      "mac": "58:41:20:0b:9b:b1",
      "ip": "192.168.8.80",
      "name": "pc1"
    }
  ],
  "id": 1
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
