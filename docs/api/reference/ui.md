# `ui`

UI Operations

6 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","ui","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`check_initialized`](#check_initialized) | - | Query if initialized |
| [`get_lang`](#get_lang) | - | Get language |
| [`get_menu_list`](#get_menu_list) | - | Get menu list |
| [`init`](#init) | - | Initialization |
| [`load_locales`](#load_locales) | - | Load translation objects |
| [`set_lang`](#set_lang) | - | Get language |

---

## check_initialized

Query if initialized

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `initialized` | bool | - | Whether initialized |
| `model` | string | - | Device model |
| `hostname` | string | - | Device hostname |
| `mac` | string | - | MAC address on the device label |
| `firmware_version` | string | - | Firmware version |
| `security_rule` | number | - | Security regulation options, used to meet different security requirements of different regulations (0: No security regulation requirements; 1: General security regulation requirements) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "ui",
    "check_initialized"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "initialized": true,
    "model": "x750",
    "hostname": "GL-X750",
    "mac": "e4:95:6e:30:24:38",
    "firmware_version": "1.2",
    "security_rule": 1
  }
}
```

---

## get_lang

Get language

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `lang` | string | - | Language (auto means automatic) |
| `list` | array | - | Supported language list |
| `list.value` | string | - | Language code |
| `list.lable` | string | - | Language description |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "ui",
    "get_lang"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "lang": "auto",
    "list": [
      {
        "value": "en",
        "lable": "English"
      },
      {
        "value": "zh-cn",
        "lable": "\u7b80\u4f53\u4e2d\u6587"
      },
      {
        "value": "zh-tw",
        "lable": "\u7e41\u4f53\u4e2d\u6587"
      },
      {
        "value": "ja",
        "lable": "\u65e5\u672c\u8a9e"
      }
    ]
  }
}
```

---

## get_menu_list

Get menu list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `menus` | array | - | Current menu list |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "ui",
    "get_menu_list"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "menus": [
      {
        "title": "Network",
        "index": 30,
        "path": "/network",
        "children": [
          {
            "title": "DHCP/DNS",
            "path": "/network/dhcp",
            "index": 40,
            "view": "oui-app-dhcp"
          }
        ]
      },
      {
        "title": "Test",
        "index": 31,
        "path": "/test",
        "view": "oui-app-test"
      }
    ]
  }
}
```

---

## init

Initialization

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `lang` | string | yes | Language |
| `username` | string | yes | Username |
| `password` | string | yes | Password |
| `ssid` | string | no | 2.4G SSID |
| `encryption_wifi` | string | no | 2.4G WiFi Encryption Method |
| `password_wifi` | string | no | 2.4G WiFi Password |
| `ssid_5g` | string | no | 5G SSID |
| `encryption_wifi_5g` | string | no | 5G WiFi Encryption Method |
| `password_wifi_5g` | string | no | 5G WiFi Password |
| `location` | string | no | Router's Placement Position |
| `redirect_https` | bool | no | Whether to Redirect HTTP to HTTPS |
| `security_rule` | number | no | Security Regulation Options, Used to Meet Different Security Requirements of Different Regulations (0: No Security Regulation Requirements; 1: General Security Regulation Requirements) |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "ui",
    "init",
    {
      "lang": "auto",
      "username": "admin",
      "password": "123",
      "ssid": "gl-x750",
      "password_wifi": "goodlife",
      "location": "bedroom",
      "security_rule": 1
    }
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

## load_locales

Load translation objects

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `locale` | string | yes | Language code |

| Results | Type | Required | Description |
|---|---|---|---|
| `locales` | array | - | Translation object list (frontend needs to merge this array into a JSON object) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "ui",
    "load_locales",
    {
      "locale": "en"
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
    "locales": [
      {
        "IPv4-Address": "IPv4-Address",
        "Services": "Services"
      },
      {
        "Output": "Output",
        "Month Days": "Month Days"
      }
    ]
  }
}
```

---

## set_lang

Get language

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `lang` | string | yes | Language to set (auto means automatic) |

| Results | Type | Required | Description |
|---|---|---|---|
| `lang` | string | - | Current language |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "ui",
    "set_lang",
    {
      "lang": "auto"
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
    "lang": "auto"
  }
}
```

---
