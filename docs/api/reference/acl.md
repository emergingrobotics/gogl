# `acl`

Permission Management

8 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","acl","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_acl`](#add_acl) | - | Add permission |
| [`add_group`](#add_group) | - | Add permission group |
| [`add_user`](#add_user) | - | Add a user to a permission group |
| [`get_acl_list`](#get_acl_list) | - | Get all permissions of a certain group |
| [`get_group_list`](#get_group_list) | - | Get permission group list |
| [`remove_acl`](#remove_acl) | - | Get all permissions of a certain group |
| [`remove_group`](#remove_group) | - | Delete permission group |
| [`remove_user`](#remove_user) | - | Delete user |

---

## add_acl

Add permission

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group` | string | yes | Permission group |
| `scope` | string | yes | Permission scope |
| `entry` | string | yes | Permission entry |
| `perm` | string | yes | Permission |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: Already exists) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "acl",
    "add_acl",
    {
      "group": "test",
      "scope": "rpc",
      "entry": "system.info",
      "perm": "x"
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

## add_group

Add permission group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group` | string | yes | The permission group to add |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: Already exists) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "acl",
    "add_group",
    {
      "group": "test"
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

## add_user

Add a user to a permission group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group` | string | yes | Permission group |
| `username` | string | yes | Username |

_No results._

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "acl",
    "add_user",
    {
      "group": "test",
      "username": "test"
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

## get_acl_list

Get all permissions of a certain group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group` | string | yes | The permission group to get |

| Results | Type | Required | Description |
|---|---|---|---|
| `acls` | array | - | Permission list |
| `acls.scope` | string | - | Permission scope |
| `acls.entry` | string | - | Permission entry |
| `acls.perm` | string | - | Permission |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "acl",
    "get_acl_list",
    {
      "group": "test"
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
    "acls": [
      {
        "scope": "rpc",
        "entry": "system.info",
        "perm": "x"
      }
    ]
  }
}
```

---

## get_group_list

Get permission group list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `groups` | array | - | Permission group list |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "acl",
    "get_group_list"
  ]
}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "groups": [
      "root",
      "test"
    ]
  }
}
```

---

## remove_acl

Get all permissions of a certain group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group` | string | yes | Permission group |
| `scope` | string | yes | Permission scope |
| `entry` | string | yes | Permission entry |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: does not exist) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "acl",
    "remove_acl",
    {
      "group": "test",
      "scope": "rpc",
      "entry": "system.info"
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

## remove_group

Delete permission group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group` | string | yes | The permission group to delete |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: Does not exist) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "acl",
    "remove_group",
    {
      "group": "test"
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

## remove_user

Delete user

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `username` | string | yes | Username |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: does not exist) |

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "call",
  "params": [
    "",
    "acl",
    "remove_user",
    {
      "group": "test",
      "username": "test"
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
