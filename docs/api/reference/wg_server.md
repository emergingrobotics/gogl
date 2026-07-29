# `wg_server`

Wireguard Server API

18 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","wg_server","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_peer`](#add_peer) | - | Add new peer |
| [`add_route`](#add_route) | - | add wireguard server custom route |
| [`generate_key`](#generate_key) | - | Generate wireguard server private_key and public_key |
| [`generate_peer`](#generate_peer) | - | Generate wireguard peer config |
| [`generate_publickey`](#generate_publickey) | - | generate wireguard server public key |
| [`get_config`](#get_config) | - | Get the server config |
| [`get_peer_list`](#get_peer_list) | - | Get the peer list of the wireguard client |
| [`get_route_list`](#get_route_list) | - | get wireguard server custom route list |
| [`get_setting`](#get_setting) | - | get wireguard server setting |
| [`get_status`](#get_status) | - | Get wireguard server status |
| [`remove_peer`](#remove_peer) | - | Remove wireguard peer by name |
| [`remove_route`](#remove_route) | - | remove wireguard server custom route |
| [`set_config`](#set_config) | - | Set the wireguard server config |
| [`set_peer`](#set_peer) | - | modify peer config |
| [`set_route`](#set_route) | - | set wireguard server custom route |
| [`set_setting`](#set_setting) | - | set wireguard server setting |
| [`start`](#start) | - | Start the wireguard server |
| [`stop`](#stop) | - | Stop wireguard server |

---

## add_peer

Add new peer

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Wireguard peer name |
| `presharedkey` | string | no | Preshared key |
| `dns` | string | no | DNS server |
| `presharedkey_enable` | bool | no | Whether to use preshared key |
| `allowed_ips` | array | no | Forwardable subnets |
| `mtu` | number | no | MTU |
| `persistent_keepalive` | number | no | Keepalive interval |

| Results | Type | Required | Description |
|---|---|---|---|
| `peer_id` | number | - | Node ID |
| `err_code` | number | - | Error code (-22: Missing parameter -23: Invalid length -24: Empty name -25: UCI file missing -26: Failed to generate IP address -27: Name already exists) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"add_peer\",{\"name\":\"test\",\"presharedkey_enable\":true,\"presharedkey\":\"\",\"allowed_ips\":[\"0.0.0.0/0,::/0\"],\"dns\":\"64.6.64.6\",\"mtu\":1420,\"persistent_keepalive\":25}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"peer_id\":1921}}
```

---

## add_route

add wireguard server custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dest` | string | yes | Destination address |
| `gateway` | string | no | Gateway |
| `scope` | string | no | scope (host, link, global, or a number greater than 0) |
| `route_flag` | number | yes | Route flag: 4 for IPv4 route, 6 for IPv6 route |
| `metric` | number | no | Metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"add_route\",{\"route_flag\":4,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\":{}}
```

---

## generate_key

Generate wireguard server private_key and public_key

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"generate_key\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\":{\"private_key\":\"oOMRHz5wTf5Gq6TPnpXLKIOhy54wx5XO7/4nquChIFI=\",\"public_key\":\"m7pNpcTJxf2gqpmgSa3dCMgY0c2nP5rXHV7ume0/uzs=\"}}
```

---

## generate_peer

Generate wireguard peer config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `peer_id` | number | yes | Node ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `address` | string | - | Node address |
| `dns` | string | - | Node DNS |
| `allowed_ips` | string | - | Node's allowed IPs |
| `domain_end_point` | string | - | Node's endpoint (domain format) |
| `ip_end_point` | string | - | Node's endpoint (IP format) |
| `listen_port` | number | - | Listen port |
| `persistent_keepalive` | number | - | Node keepalive |
| `private_key` | string | - | Node private key |
| `public_key` | string | - | Node public key |
| `mtu` | number | - | MTU |
| `presharedkey` | number | - | Preshared key |
| `err_code` | number | - | Error code (-28: Missing parameter -29: UCI file missing -30: Client not found -31: Private key does not exist -32: No network -33: Public key does not exist) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"generate_peer\",{\"peer_id\":2284}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\":{\"address\":\"10.0.0.2/32\",\"listen_port\":50436,\"private_key\":\"oOMRHz5wTf5Gq6TPnpXLKIOhy54wx5XO7/4nquChIFI=\",\"dns\":\"64.6.64.6\",\"domain_end_point\":\"gwe10df.glddns.com:51820\",\"ip_end_point\":\"116.77.73.248:51820\",\"public_key\":\"m7pNpcTJxf2gqpmgSa3dCMgY0c2nP5rXHV7ume0/uzs=\",\"allowed_ips\":\"0.0.0.0/0\",\"persistent_keepalive\":25,\"mtu\":1420}}
```

---

## generate_publickey

generate wireguard server public key

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `private_key` | string | yes | Private key |

| Results | Type | Required | Description |
|---|---|---|---|
| `public_key` | string | - | Public key |
| `err_code` | number | - | Error code (-10: Missing parameter -11: Failed to generate public key) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"generate_publickey\",{\"private_key\":\"XVpIdr+oYjTcgDwzSZmNa1nSsk8JO+tx1NBo17LDBAI=\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"public_key\":\"MANEa9SUz3kh6JujfX11rkj4Tw9WKnU7A9wmr4nmLiA=\"}}
```

---

## get_config

Get the server config

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `local_access` | bool | - | Local access switch |
| `address_v4` | string | - | IPv4 subnet |
| `address_v6` | string | - | IPv6 subnet |
| `port` | number | - | Port |
| `public_key` | string | - | Public key |
| `private_key` | string | - | Private key |
| `initialization` | bool | - | Whether initialized |
| `ipv6_enable` | bool | - | Whether IPv6 is enabled |
| `err_code` | number | - | Error code (-15: WG not installed -16: UCI file missing) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"get_config\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"port\": 51820, \"address_v4\": \"10.8.0.0/24\", \"address_v6\": \"fd00:db8:0:123::/64\", \"local_access\": false, \"public_key\":\"izv0p34WZN7p2vIgehwe33QF27ExjChrPUisk481JHU0=\",,\"initialization\":true,\"ipv6_enable\":false}}
```

---

## get_peer_list

Get the peer list of the wireguard client

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `peers` | array | - | Node information |
| `peers.name` | string | - | Node name |
| `peers.client_ip` | string | - | Client address |
| `peers.dns` | string | - | Node DNS |
| `peers.allowed_ips` | string | - | Node's allowed IPs |
| `peers.end_point` | string | - | Node's endpoint |
| `peers.listen_port` | number | - | Listen port |
| `peers.persistent_keepalive` | number | - | Node keepalive |
| `peers.private_key` | string | - | Node private key |
| `peers.public_key` | string | - | Node public key |
| `peers.mtu` | number | - | MTU |
| `peers.peer_id` | number | - | Node ID |
| `peers.presharedkey_enable` | bool | - | Whether to use preshared key |
| `peers.presharedkey` | number | - | Preshared key |
| `peers.deprecated` | number | - | Whether peer is available (0: available 1: unavailable) |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"get_peer_list\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"peers\":[{\"name\":\"sample\",\"client_ip\":\"10.0.0.2/32\",\"listen_port\":50436,\"private_key\":\"oOMRHz5wTf5Gq6TPnpXLKIOhy54wx5XO7/4nquChIFI=\",\"dns\":\"64.6.64.6\",\"end_point\":\"116.77.73.248:51820\",\"public_key\":\"m7pNpcTJxf2gqpmgSa3dCMgY0c2nP5rXHV7ume0/uzs=\",\"allowed_ips\":\"0.0.0.0/0\",\"persistent_keepalive\":25,\"mtu\":1420,\"peer_id\":2844,\"presharedkey_enable\":true,\"presharedkey\":\"wKaBacfDah7EXNzIq9AbpHso4NUncS6M7NKPViUWxc0=\",\"deprecated\":0}]}}
```

---

## get_route_list

get wireguard server custom route list

Not exercised against hardware here.

_No params._

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
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"get_route_list\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"ipv4_route_rules\":[{\"rule_id\":6512,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"ipv6_route_rules\":[{\"rule_id\":6515,\"dest\":\"ff00::/8\",\"gateway\":\"ff00::1\",\"metric\":20,\"mtu\":1500}]}}
```

---

## get_setting

get wireguard server setting

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `local_access` | bool | - | Whether to enable local access |
| `masq` | bool | - | Whether to enable masq |
| `mtu` | number | - | MTU |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"get_setting\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"local_access\":true,\"masq\":true,\"mtu\":1460}}
```

---

## get_status

Get wireguard server status

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `server` | object | - | Node information |
| `status` | number | - | Server status 0: Not connected 1: Connected successfully 2: Connecting |
| `initialization` | bool | - | Whether initialized |
| `log` | string | - | Log |
| `tunnel_ip` | string | - | tunnel ip |
| `rx_bytes` | number | - | Received bytes |
| `tx_bytes` | number | - | Sent bytes |
| `peers` | array | - | Node information |
| `peers.name` | string | - | Client name |
| `peers.status` | number | - | Client status 0: offline 1: online |
| `peers.private_ip` | string | - | Private address |
| `peers.public_ip` | string | - | Public address |
| `peers.latest_handshake` | string | - | Last handshake time |
| `peers.rx_bytes` | number | - | Received bytes |
| `peers.tx_bytes` | number | - | Sent bytes |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"get_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"server\":{\"status\": 1,\"initialization\":true,\"log\":\"\",\"tunnel_ip\":\"10.9.0.0/24, fc00:bbbb:bb01::/96\",\"rx_bytes\":245,\"tx_bytes\":168},\"peers\":[{\"name\":\"sample\",\"status\":1,\"private_ip\":\"10.8.0.2/32\",\"public_ip\":\"145.155.54.184\",\"latest_handshake\":\"24 seconds ago\",\"rx_bytes\":149,\"tx_bytes\":92}]}}
```

---

## remove_peer

Remove wireguard peer by name

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `peer_id` | number | yes | Node ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-19: Missing parameter -20: Client not found) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"remove_peer\",{\"peer_id\":2284}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## remove_route

remove wireguard server custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `route_flag` | number | yes | Route flag: 4 for IPv4 route, 6 for IPv6 route |
| `rule_id` | number | yes | Route rule ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"remove_route\",{\"rule_id\":6512,\"route_flag\":4}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\":{}}
```

---

## set_config

Set the wireguard server config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `address_v4` | string | yes | IPv4 subnet |
| `address_v6` | string | no | IPv6 subnet |
| `private_key` | string | no | Private key |
| `ipv6_enable` | bool | no | Whether to enable IPv6 |
| `port` | number | yes | Listening port |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-9: Missing parameters -10: UCI file missing -11: Invalid port -12: Invalid IP address -13: Server IP conflicts with local IP -14: IPv6 address error -15: Port occupied -16: Failed to generate public key) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"set_config\",{\"ipv6_enable\":true,\"address_v4\":\"10.8.0.0/24\",\"address_v6\":\"fd00:db8:0:abc::/64\",\"port\":51820,\"private_key\":\"XVpIdr+oYjTcgDwzSZmNa1nSsk8JO+tx1NBo17LDBAI=\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_peer

modify peer config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Wireguard peer name |
| `presharedkey` | string | no | Preshared key |
| `dns` | string | no | DNS server |
| `presharedkey_enable` | bool | yes | Whether to use preshared key |
| `allowed_ips` | array | no | Forwardable subnets |
| `peer_id` | number | yes | ID |
| `mtu` | number | no | MTU |
| `persistent_keepalive` | number | no | Keepalive interval |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-22: Missing parameter -23: Invalid length -24: Empty name -25: UCI file missing -26: Failed to generate IP address -27: Name already exists) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"set_peer\",{\"name\":\"test\",\"peer_id\":2844,\"presharedkey_enable\":true,\"presharedkey\":\"\",\"allowed_ips\":[\"0.0.0.0/0,::/0\"],\"dns\":\"64.6.64.6\",\"mtu\":1420,\"persistent_keepalive\":25}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_route

set wireguard server custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dest` | string | no | Destination address |
| `gateway` | string | no | Gateway |
| `scope` | string | no | scope (host, link, global, or a number greater than 0) |
| `route_flag` | number | yes | Route flag: 4 for IPv4 route, 6 for IPv6 route |
| `rule_id` | number | yes | Route rule ID |
| `metric` | number | no | Metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"set_route\",{\"rule_id\":6512,\"route_flag\":4,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\":{}}
```

---

## set_setting

set wireguard server setting

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `local_access` | bool | no | Whether to enable local access |
| `masq` | bool | no | Whether to enable masq |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"set_setting\",{\"local_access\":true,\"masq\":true,\"mtu\":1460}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## start

Start the wireguard server

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-4: VPN conflict -5: Please stop WG client first -6: Program is already running -7: WG not installed -8: Startup script missing) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"start\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## stop

Stop wireguard server

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"wg-server\",\"stop\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---
