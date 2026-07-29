# `wg_client`

Wireguard Client API

24 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","wg_client","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_config`](#add_config) | - | add wireguard client config |
| [`add_group`](#add_group) | - | add wireguard client group |
| [`add_route`](#add_route) | - | add wireguard client custom route |
| [`check_config`](#check_config) | - | Check uploaded wireguard file's parameter |
| [`clear_config_list`](#clear_config_list) | - | remove all wireguard client config of a group |
| [`confirm_config`](#confirm_config) | - | confirm wireguard client config |
| [`get_all_config_list`](#get_all_config_list) | - | get wireguard client config list of all groups |
| [`get_config_list`](#get_config_list) | - | get wireguard client config list of a group |
| [`get_group_list`](#get_group_list) | - | get wireguard client group list |
| [`get_recommend_config`](#get_recommend_config) | - | get recommend config |
| [`get_route_list`](#get_route_list) | - | get wireguard client custom route list |
| [`get_setting`](#get_setting) | - | get wireguard client setting |
| [`get_status`](#get_status) | - | get the wireguard client status |
| [`get_third_config`](#get_third_config) | - | get third provider's client config |
| [`remove_config`](#remove_config) | - | remove wireguard client config |
| [`remove_group`](#remove_group) | - | remove wireguard client group |
| [`remove_route`](#remove_route) | - | remove wireguard client custom route |
| [`set_config`](#set_config) | - | modify specific client's config |
| [`set_group`](#set_group) | - | modify wireguard client group |
| [`set_proxy`](#set_proxy) | - | set wireguard client proxy |
| [`set_route`](#set_route) | - | set wireguard client custom route |
| [`set_setting`](#set_setting) | - | set wireguard client setting |
| [`start`](#start) | - | start wireguard client |
| [`stop`](#stop) | - | stop wireguard client |

---

## add_config

add wireguard client config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | configuration name |
| `address_v4` | string | yes | IPv4 subnet |
| `address_v6` | string | no | IPv6 subnet |
| `private_key` | string | yes | node private key |
| `allowed_ips` | string | yes | allowed IPs |
| `end_point` | string | yes | endpoint |
| `public_key` | string | yes | public key |
| `dns` | string | no | dns |
| `preshared_key` | string | no | preshared key |
| `ipv6_enable` | bool | no | whether to enable IPv6 |
| `presharedkey_enable` | bool | yes | whether to use preshared key |
| `group_id` | number | yes | group ID |
| `listen_port` | number | no | listen port |
| `persistent_keepalive` | number | no | persistent keepalive interval |
| `mtu` | number | no | mtu |

| Results | Type | Required | Description |
|---|---|---|---|
| `peer_id` | number | - | configuration ID |
| `err_code` | number | - | error code (-19: missing parameters -20: length too long -21: empty name -22: UCI file missing -23: invalid port) |
| `err_msg` | string | - | error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"add_config\",{\"group_id\":3212,\"name\":\"test\",\"address_v4\":\"10.8.0.0/24\",\"address_v6\":\"fd00:db8:0:123::/64\",\"private_key\":\"XVpIdr+oYjTcgDwzSZmNa1nSsk8JO+tx1NBo17LDBAI=\",\"allowed_ips\":\"0.0.0.0/0,::/0\",\"end_point\":\"103.231.88.18:3102\",\"public_key\":\"zv0p34WZN7p2vIgehwe33QF27ExjChrPUisk481JHU0=\",\"dns\":\"193.138.219.228\",\"presharedkey_enable\":false,\"listen_port\":22536,\"persistent_keepalive\":25,\"mtu\":1420,\"ipv6_enable\":true}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"peer_id\":1921}}
```

---

## add_group

add wireguard client group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_name` | string | yes | Group name |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-40: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"add_group\",{\"group_name\":\"mullvad\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## add_route

add wireguard client custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dest` | string | yes | Destination address |
| `gateway` | string | no | Gateway |
| `scope` | string | no | scope (host, link, global, or number greater than 0) |
| `group_id` | number | yes | Group ID |
| `peer_id` | number | yes | Configuration ID |
| `route_flag` | number | yes | Route flag 4 for IPv4 route, 6 for IPv6 route |
| `metric` | number | no | Metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code(-40: missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"add_route\",{\"group_id\":3213,\"peer_id\":5715,\"route_flag\":4,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## check_config

Check uploaded wireguard file's parameter

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `filename` | string | yes | File name to be verified |
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `passed` | array | - | List configurations that do not require additional information |
| `unpassed` | array | - | List configurations that failed the check |
| `err_code` | number | - | Error code (-40: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"check_config\",{\"group_id\":2312,\"filename\":\"mullvad-at4.conf\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"passed\": [], \"unpassed\": []}}
```

---

## clear_config_list

remove all wireguard client config of a group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-40: Missing parameter -41: Client is running) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"clear_config_list\",{\"group_id\":3212}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## confirm_config

confirm wireguard client config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-40: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"confirm_config\",{\"group_id\":2312}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## get_all_config_list

get wireguard client config list of all groups

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `config_list` | array | - | Node information |
| `config_list.group_id` | number | - | Group ID |
| `config_list.group_name` | string | - | Group name |
| `config_list.auth_type` | number | - | Authentication type 0: no username/password or account required 1: username/password 2: account |
| `config_list.username` | string | - | Username |
| `config_list.password` | string | - | Password |
| `config_list.peers` | array | - | Node information |
| `config_list.peers.peer_id` | number | - | Configuration ID |
| `config_list.peers.address_v4` | string | - | IPv4 subnet |
| `config_list.peers.ipv6_enable` | bool | - | Whether to use IPv6 |
| `config_list.peers.address_v6` | string | - | IPv6 subnet |
| `config_list.peers.name` | string | - | Configuration name |
| `config_list.peers.allowed_ips` | string | - | allowips |
| `config_list.peers.dns` | string | - | dns |
| `config_list.peers.end_point` | string | - | endpoint |
| `config_list.peers.listen_port` | number | - | Listen port |
| `config_list.peers.mtu` | number | - | mtu |
| `config_list.peers.persistent_keepalive` | number | - | Keepalive interval |
| `config_list.peers.presharedkey_enable` | bool | - | Whether to use pre-shared key |
| `config_list.peers.presharedkey` | string | - | Pre-shared key |
| `config_list.peers.private_key` | string | - | Private key |
| `config_list.peers.public_key` | string | - | Public key |
| `config_list.peers.location` | string | - | Server location |
| `err_code` | number | - | Error code(-49: empty output -50: client configuration not found) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"get_all_config_list\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"config_list\":[{\"group_id\":2312,\"group_name\":mullvad,\"auth_type\":1,\"username\":\"glinet\",\"password\":123456,\"peers\":[{\"peer_id\":1254,\"address_v4\":\"10.8.0.0/24\",\"ipv6_enable\":true,\"address_v6\":\"fd00:db8:0:123::/64\",\"name\":\"mullvad_au4\",\"allowed_ips\":\"0.0.0.0/0,::/0\",\"dns\":\"193.138.219.228\",\"end_point\":\"103.231.88.2:3100\",\"listen_port\":22536,\"mtu\":1380,\"persistent_keepalive\":25,\"private_key\":\"EHwE7Aiiup/lxVEAKlofJmmMfQZkHTmmEy3xcAcIn2I=\",\"public_key\":\"kXXykjh6KqiE/pvtmTV8kCB+jhhkl9kT0Dg+yyDz8hg=\",\"presharedkey_enable\":false,\"location\":\"Australia, Melbourne\"}]}]}}
```

---

## get_config_list

get wireguard client config list of a group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `auth_type` | number | - | Authentication type 0: No username/password or account required 1: Username/password 2: Account |
| `public_key` | string | - | Mullvad public key |
| `peers` | array | - | Node information |
| `peers.peer_id` | number | - | Configuration ID |
| `peers.address_v4` | string | - | IPv4 subnet |
| `peers.ipv6_enable` | bool | - | Whether to use IPv6 |
| `peers.address_v6` | string | - | IPv6 subnet |
| `peers.name` | string | - | Configuration name |
| `peers.allowed_ips` | string | - | allowips |
| `peers.dns` | string | - | dns |
| `peers.end_point` | string | - | endpoint |
| `peers.listen_port` | number | - | Listen port |
| `peers.mtu` | number | - | mtu |
| `peers.persistent_keepalive` | number | - | Keepalive interval |
| `peers.presharedkey_enable` | bool | - | Whether to use preshared key |
| `peers.presharedkey` | string | - | Preshared key |
| `peers.private_key` | string | - | Private key |
| `peers.public_key` | string | - | Public key |
| `peers.location` | string | - | Server location |
| `err_code` | number | - | Error code (-40: Missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"get_config_list\",{\"group_id\":3212}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"auth_type\":2,\"public_key\":\"NjUxYTE3NzZiMjczNGMzNjg2MmNiMmNhMTBkZGYyYTQ=\",\"peers\":[{\"peer_id\":1254,\"address_v4\":\"10.8.0.0/24\",\"ipv6_enable\":true,\"address_v6\":\"fd00:db8:0:123::/64\",\"name\":\"mullvad_au4\",\"allowed_ips\":\"0.0.0.0/0,::/0\",\"dns\":\"193.138.219.228\",\"end_point\":\"103.231.88.2:3100\",\"listen_port\":22536,\"mtu\":1380,\"persistent_keepalive\":25,\"private_key\":\"EHwE7Aiiup/lxVEAKlofJmmMfQZkHTmmEy3xcAcIn2I=\",\"public_key\":\"kXXykjh6KqiE/pvtmTV8kCB+jhhkl9kT0Dg+yyDz8hg=\",\"presharedkey_enable\":false,\"location\":\"Australia, Melbourne\"}]}}
```

---

## get_group_list

get wireguard client group list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `groups` | array | - | Group nodes |
| `groups.group_id` | number | - | Group ID |
| `groups.group_type` | number | - | Group type 1: default group 2: manually added group 3: app group |
| `groups.group_name` | string | - | Group name |
| `groups.auth_type` | number | - | Authentication type 0: no username/password or account required 1: username/password 2: account |
| `groups.username` | string | - | Username |
| `groups.password` | string | - | Password |
| `groups.peer_count` | number | - | Number of configurations in the group |
| `groups.procedure` | number | - | Configuration download process 0: direct download 1: step-by-step download |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"get_group_list\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"groups\": [{\"group_id\":3212,\"group_name\":\"mullvad\",\"peer_count\":20,\"group_type\":1,\"procedure\":1},{\"group_id\":3213,\"group_name\":\"azirevpn\",\"peer_count\":25,\"group_type\":1,\"procedure\":0}]}}
```

---

## get_recommend_config

get recommend config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `country_name` | array | yes | country name |
| `group_id` | number | yes | group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | error code (-10: missing parameters -11: UCI configuration missing -12: group information missing -13: missing password -14: country information empty -15: unknown service provider) |
| `err_msg` | string | - | error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"get_recommend_config\",{\"group_id\": 2312, \"country_name\":[\"Australia\",\"Brazil\"]}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## get_route_list

get wireguard client custom route list

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `peer_id` | number | yes | Configuration ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `ipv4_route_rules` | array | - | IPv4 route rules |
| `ipv4_route_rules.rule_id` | number | - | Route rule ID |
| `ipv4_route_rules.dest` | string | - | Destination address |
| `ipv4_route_rules.gateway` | string | - | Gateway |
| `ipv4_route_rules.metric` | number | - | Metric |
| `ipv4_route_rules.mtu` | number | - | MTU |
| `ipv4_route_rules.scope` | string | - | scope (host, link, global, or a number greater than 0) |
| `ipv6_route_rules` | array | - | IPv6 route rules |
| `ipv6_route_rules.rule_id` | number | - | Route rule ID |
| `ipv6_route_rules.dest` | string | - | Destination address |
| `ipv6_route_rules.gateway` | string | - | Gateway |
| `ipv6_route_rules.metric` | number | - | Metric |
| `ipv6_route_rules.mtu` | number | - | MTU |
| `ipv6_route_rules.scope` | string | - | scope (host, link, global, or a number greater than 0) |
| `err_code` | number | - | Error code (-40: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"get_route_list\",{\"group_id\":3213,\"peer_id\":5715}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"ipv4_route_rules\":[{\"route_id\":6512,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"ipv6_route_rules\":[{\"route_id\":6515,\"dest\":\"ff00::/8\",\"gateway\":\"ff00::1\",\"metric\":20,\"mtu\":1500}]}}
```

---

## get_setting

get wireguard client setting

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `peer_id` | number | yes | Client ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `local_access` | bool | - | Whether local access is enabled |
| `masq` | bool | - | Whether masq is enabled |
| `default_metric` | number | - | Default metric |
| `mtu` | number | - | MTU |
| `err_code` | number | - | Error code (-40: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"get_setting\",{\"group_id\":3213,\"peer_id\":5715}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"local_access\":true,\"masq\":true,\"default_metric\":30,\"mtu\":1460}}
```

---

## get_status

get the wireguard client status

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | - | Client's group ID |
| `peer_id` | number | - | Client ID |
| `name` | string | - | Client name |
| `proxy` | bool | - | Global proxy |
| `log` | string | - | Log |
| `ipv4` | string | - | IPv4 address |
| `ipv6` | string | - | IPv6 address |
| `domain` | string | - | Server domain |
| `port` | number | - | Server listen port |
| `status` | number | - | Connection status 0: Not connected 1: Connected successfully 2: Connecting |
| `rx_bytes` | number | - | Received bytes |
| `tx_bytes` | number | - | Sent bytes |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"get_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"group_id\":2312, \"peer_id\":5712, \"name\":\"mullvad_ae_all\", \"proxy\":false, \"log\":\"\", \"status\":1, \"ipv4\":\"10.70.33.3\", \"ipv6\":\"fc00:bbbb:bb01::11\", \"rx_bytes\":50, \"tx_bytes\":45, \"port\":41348, \"domain\":\"ae-dxb-001.mullvad.net\"}}
```

---

## get_third_config

get third provider's client config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `username` | string | yes | username |
| `password` | string | no | password (required if authentication type is username/password) |
| `group_id` | number | yes | group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `countries` | array | - | node information |
| `countries.country_name` | string | - | country name |
| `countries.server_count` | number | - | server count |
| `err_code` | number | - | error code (-10: missing parameters -11: invalid name length (exceeds 64 bytes) -12: name empty -13: UCI configuration missing -14: missing password -15: failed to get AzireVPN configuration -16: invalid username or password -17: unknown service provider) |
| `err_msg` | string | - | error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"get_third_config\",{\"group_id\":2312,\"username\":\"glinet\",\"password\":\"123456\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"countries\":[{\"country_name\":\"Australia\",\"server_count\":11},{\"country_name\":\"Brazil\",\"server_count\":5}]}}
```

---

## remove_config

remove wireguard client config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | group ID |
| `peer_id` | number | yes | configuration ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | error code (-43: missing parameters -44: client is running) |
| `err_msg` | string | - | error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"remove_config\",{\"group_id\":3212,\"peer_id\":1254}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## remove_group

remove wireguard client group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-40: missing parameter -41: client is running) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"remove_group\",{\"group_id\":3212}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## remove_route

remove wireguard client custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `peer_id` | number | yes | Configuration ID |
| `route_flag` | number | yes | Route flag, 4 for IPv4 route, 6 for IPv6 route |
| `rule_id` | number | yes | Route rule ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-40: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"remove_route\",{\"group_id\":3213,\"peer_id\":5715,\"route_flag\":4,\"rule_id\":6512}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_config

modify specific client's config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | node name |
| `address_v4` | string | yes | node IPv4 subnet |
| `address_v6` | string | no | node IPv6 subnet |
| `private_key` | string | yes | node private key |
| `allowed_ips` | string | yes | node's allowed IPs |
| `end_point` | string | yes | node's endpoint |
| `public_key` | string | yes | node public key |
| `dns` | string | no | node's dns |
| `preshared_key` | string | no | preshared key |
| `ipv6_enable` | bool | no | whether to enable IPv6 |
| `presharedkey_enable` | bool | yes | whether to use preshared key |
| `group_id` | number | yes | group ID |
| `peer_id` | number | yes | configuration ID |
| `listen_port` | number | no | listen port |
| `persistent_keepalive` | number | no | node keepalive |
| `mtu` | number | no | node's mtu |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | error code (-31: missing parameters -32: name too long -33: name empty -34: UCI file missing -35: invalid port) |
| `err_msg` | string | - | error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"set_config\",{\"group_id\":3212,\"peer_id\":1254,\"name\":\"test\",\"address_v4\":\"10.8.0.0/24\",\"address_v6\":\"fd00:db8:0:123::/64\",\"private_key\":\"XVpIdr+oYjTcgDwzSZmNa1nSsk8JO+tx1NBo17LDBAI=\",\"allowed_ips\":\"0.0.0.0/0,::/0\",\"end_point\":\"103.231.88.18:3102\",\"public_key\":\"zv0p34WZN7p2vIgehwe33QF27ExjChrPUisk481JHU0=\",\"dns\":\"193.138.219.228\",\"presharedkey_enable\":false,\"listen_port\":22536,\"persistent_keepalive\":25,\"mtu\":1420,\"ipv6_enable\":true}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_group

modify wireguard client group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_name` | string | yes | Group name |
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-40: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"set_group\",{\"group_id\":3212,\"group_name\":\"mullvad\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_proxy

set wireguard client proxy

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `enable` | bool | yes | Enable or disable global proxy |

_No results._

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"set_proxy\",{\"enable\":true}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_route

set wireguard client custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dest` | string | no | Destination address |
| `gateway` | string | no | Gateway |
| `scope` | string | no | scope (host, link, global, or number greater than 0) |
| `group_id` | number | yes | Group ID |
| `peer_id` | number | yes | Configuration ID |
| `route_flag` | number | yes | Route flag 4 for IPv4 route, 6 for IPv6 route |
| `rule_id` | number | yes | Route rule ID |
| `metric` | number | no | Metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code(-40: missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"set_route\",{\"group_id\":3213,\"peer_id\":5715,\"route_flag\":4,\"rule_id\":6512,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_setting

set wireguard client setting

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `local_access` | bool | no | Whether to enable local access |
| `masq` | bool | no | Whether to enable masq |
| `group_id` | number | yes | Group ID |
| `peer_id` | number | yes | Client ID |
| `default_metric` | number | no | Default metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-40: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"set_setting\",{\"group_id\":3213,\"peer_id\":5715,\"local_access\":true,\"masq\":true,\"default_metric\":30,\"mtu\":1460}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## start

start wireguard client

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `peer_id` | number | yes | Configuration ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-51: Missing parameter -52: VPN conflict -53: Please stop WG server first -54: WG not installed -55: Missing startup script -56: UCI file missing -57: Client not found) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"start\",{\"group_id\":3212,\"peer_id\":1254}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## stop

stop wireguard client

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-client\",\"stop\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---
