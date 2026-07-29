# `ovpn_server`

ovpn server api

16 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","ovpn_server","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_route`](#add_route) | - | add ovpn server custom route |
| [`add_user`](#add_user) | - | add ovpn user |
| [`export_config`](#export_config) | - | Export ovpn configuration |
| [`generate_certificate`](#generate_certificate) | - | Generate ovpn cert |
| [`get_config`](#get_config) | - | Get openvpn configuration |
| [`get_route_list`](#get_route_list) | - | get ovpn server custom route list |
| [`get_setting`](#get_setting) | - | get ovpn server setting |
| [`get_status`](#get_status) | - | Get ovpn server status |
| [`get_user_list`](#get_user_list) | - | get ovpn user list |
| [`remove_route`](#remove_route) | - | remove ovpn server custom route |
| [`remove_user`](#remove_user) | - | remove ovpn user |
| [`set_config`](#set_config) | - | Set ovpn server config |
| [`set_route`](#set_route) | - | set ovpn server custom route |
| [`set_setting`](#set_setting) | - | set ovpn server setting |
| [`start`](#start) | - | start ovpn server |
| [`stop`](#stop) | - | Stop ovpn server |

---

## add_route

add ovpn server custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dest` | string | yes | destination address |
| `gateway` | string | no | gateway |
| `scope` | string | no | scope (host, link, global, or numbers greater than 0) |
| `route_flag` | number | yes | Route flag, 4 for IPv4 route, 6 for IPv6 route |
| `metric` | number | no | metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"add_route\",{\"route_flag\":4,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## add_user

add ovpn user

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `username` | string | yes | Username |
| `password` | string | yes | Password |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"add_user\",{\"username\":\"glinet\",\"password\":\"123456\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## export_config

Export ovpn configuration

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `ddns` | bool | yes | Whether to enable DDNS |

| Results | Type | Required | Description |
|---|---|---|---|
| `file_name` | string | - | Configuration file name |
| `file_path` | string | - | Configuration file path |
| `err_code` | number | - | Error code (-10: Missing parameters -11: No network) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"export_config\",{\"ddns\":false}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"file_name\": \"client.ovpn\", \"file_path\": \"/etc/openvpn/ovpn/client.ovpn\"}}
```

---

## generate_certificate

Generate ovpn cert

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `ca` | string | no | Server root certificate |
| `dh` | string | no | Diffie Hellman parameters |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-19: No network -20: openssl missing -21: Certificate generation failed) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"generate_certificate\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## get_config

Get openvpn configuration

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `local_access` | bool | - | Server local access switch status |
| `mode` | string | - | VPN mode: tun/tap |
| `subnetv4` | string | - | IPv4 subnet |
| `mask` | string | - | Mask |
| `ipv6_enable` | bool | - | Whether to enable IPv6 |
| `subnetv6` | string | - | IPv6 subnet |
| `proto` | string | - | Transport protocol |
| `port` | number | - | Port |
| `auth` | string | - | Authentication method |
| `cipher` | string | - | Encryption method |
| `start` | string | - | IP pool start address in tap mode |
| `end` | string | - | IP pool end address in tap mode |
| `lzo` | bool | - | Whether to use LZO compression |
| `ca` | string | - | Server root certificate |
| `dh` | string | - | Diffie Hellman parameters |
| `hmac` | bool | - | HMAC authentication |
| `ta` | string | - | TLS-Auth additional key |
| `client_to_client` | bool | - | Client-to-client communication |
| `verb` | string | - | Log level |
| `access_scope` | number | - | Client using VPN 1: Access internal network (LAN) and external network (WAN) 2: Only access internal network (LAN) 3: Only access external network |
| `client_auth` | number | - | Client authentication method 1: Certificate only 2: Username/password only 3: Use both username/password and certificate |
| `cert` | string | - | Server certificate |
| `key` | string | - | Server key |
| `initialization` | bool | - | Whether initialized |
| `tap_address` | string | - | Tap mode IP address |
| `tap_mask` | string | - | Tap mode mask |
| `err_code` | number | - | Error code (-10: uci file missing -11: Failed to get IPv6 subnet -12: Failed to get TLS-Auth additional key) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"get_config\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"initialization\":true,\"local_access\":false,\"auth\":\"SHA256\",\"proto\":\"tcp\",\"port\":1194,\"subnetv4\":\"10.8.0.0\",\"ipv6_enable\":true,\"subnetv6\":\"fd00:db8:0:123::/64\",\"mask\":\"255.255.255.0\",\"cipher\":\"AES-256-GCM\",\"mode\":\"tun\",\"start\":\"10.8.0.2\",\"end\":\"10.8.0.100\",\"lzo\":false,\"ca\":\"-----BEGIN CERTIFICATE----- MIIDCzCCAfOgAwIBAgIUequqVhT2jU3wWJF988KVcIUyUtwwDQYJKoZIhvcNAQEF BQAwFTETMBEGA1UEAwwKT3BlblZQTiBDQTAeFw0yMTA2MTkwMjA4MzNaFw0zMTA2 MTcwMjA4MzNaMBUxEzARBgNVBAMMCk9wZW5WUE4gQ0EwggEiMA0GCSqGSIb3DQEB AQUAA4IBDwAwggEKAoIBAQCjLr7B1mc04E/yHqR8hv7Qng5WuaZVGYDfUoxxly0c 8qOj6LqkL5UwjtcQV2rFgucwCsTYW8/dRY/Fd6cRllNdvRuA1nuMi/JhECkmf9AB ckdXlyXiObeYNL1MofcOsjRkB1oxokOO7BN2uJ2sQqekBSc7DN+8b9w5R+gaVXY3 082bVHvIKPiwYWzpNbO3MfdBsmh9oL8CryAfDAv8aTaAeCRoHNy5hnzqWoHfv6xm MsuVfwb9Vw5MegJ011x+VMF/++o95IcL/UsMSCrhjstyaf94ucPwIi7NuD2lor6e RP0eqKryGui1r2coZweZjGuYqOnTwHZhJcjpgSK1RdbhAgMBAAGjUzBRMB0GA1Ud DgQWBBRUf/PRJrXI4eHeTBBldnTeJnrefTAfBgNVHSMEGDAWgBRUf/PRJrXI4eHe TBBldnTeJnrefTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBBQUAA4IBAQCS bAc7oMhLnAnHMPaFhXVDmvsgX7rK2BEKk/qqlHR2qAd7eSK3cteYX+Kayq+j8RmO PPsk4hXQNo+6tcaZSi05ACVYcCOkNw5dOZ8DDIo5swydamAZewp4ExsZQP4pdr40 8ikNle6Amt4psgxT0AVH56g4YO7qqdg5RMY14grIsBTL21j5OhdnwORPxstC2icW S5vzouDaIZ85IkiLaBku4fw8hBPkSmLRBexzeYwx/CoOs6olOEaKYoLy5TYFT+Ed 4Tvd7d0sw8vh1l1DiNtb5jGw7BsX4S8rzuB6o0GpxtxJprlyOPqqsyxwFNwhyAVk TFYWeM/vOcydA3zBBJbH -----END CERTIFICATE-----\",\"dh\":\"-----BEGIN DH PARAMETERS----- MIGHAoGBALUmuxsYkB+OLETWb+fz2WhAwa6PUMNzweXEA44rZ8Yv3PqfUIb/rPiW ZUwIq+ouAF1kMPnRYYhsqg5d4OeMv8R6SPlQUckpZrBELsJro9qGRf2nda9v5KAI 8DaorlIhKHpEoJ9xBz8h/DtQ85PX3UT96RZXA/SAyVAhdDpHkJ5rAgEC -----END DH PARAMETERS-----\",\"hmac\":true,\"ta\":\"-----BEGIN OpenVPN Static key V1----- 66301913fa0fa22457019b1a52e26aef 058c0405905f7e415615cc576bf1ca56 a0d12e27f181044c7ce4035afa5a7e30 1e02284208762d1ef8cd94b1989b8940 b534c469920a3012b440a4abad67af32 a7d80e0fc9447bd82b0d2ef2098a99c2 e7decaa1c5442b5ab6c0bbd413652814 c260fd76318d3abc6a65552fda44cb00 066420d3232f423f2810c861e0c3429f f566d0737d8e63c86dbf5bea5693dd51 1002a5e1e3aa95eddd388550e0022b04 8299cac19864821e3369a17e02cf02ef 0c41460146edba8748e9bccdd57cf8ac 1558f8d9a5b3430bf2229cad92625910 e09dce6311e834c2ca8c004f3a18c7d2 97188c096342674ff89c4bebda1e04db -----END OpenVPN Static key V1-----\",\"client_to_client\":false,\"verb\":\"3\",\"access_scope\":1,\"client_auth\":1,\"cert\":\"nil\",\"key\":\"nil\",\"tap_address\":\"10.8.0.1\",\"tap_mask\":\"255.255.255.255\"}}
```

---

## get_route_list

get ovpn server custom route list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `ipv4_route_rules` | array | - | IPv4 route rules |
| `ipv4_route_rules.rule_id` | number | - | Route rule ID |
| `ipv4_route_rules.dest` | string | - | destination address |
| `ipv4_route_rules.gateway` | string | - | gateway |
| `ipv4_route_rules.metric` | number | - | metric |
| `ipv4_route_rules.mtu` | number | - | MTU |
| `ipv4_route_rules.scope` | string | - | scope (host, link, global, or numbers greater than 0) |
| `ipv6_route_rules` | array | - | IPv6 route rules |
| `ipv6_route_rules.rule_id` | number | - | Route rule ID |
| `ipv6_route_rules.dest` | string | - | destination address |
| `ipv6_route_rules.gateway` | string | - | gateway |
| `ipv6_route_rules.metric` | number | - | metric |
| `ipv6_route_rules.mtu` | number | - | MTU |
| `ipv6_route_rules.scope` | string | - | scope (host, link, global, or numbers greater than 0) |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"get_route_list\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"ipv4_route_rules\":[{\"rule_id\":6512,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"ipv6_route_rules\":[{\"rule_id\":6515,\"dest\":\"ff00::/8\",\"gateway\":\"ff00::1\",\"metric\":20,\"mtu\":1500}]}}
```

---

## get_setting

get ovpn server setting

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
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"get_setting\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"local_access\":true,\"masq\":true,\"mtu\":1460}}
```

---

## get_status

Get ovpn server status

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `status` | number | - | Server status 0: Not connected 1: Connected successfully 2: Connecting |
| `initialization` | bool | - | Whether initialized |
| `log` | string | - | Log |
| `tunnel_ip` | string | - | tunnel ip |
| `rx_bytes` | number | - | Received bytes |
| `tx_bytes` | number | - | Sent bytes |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"get_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"status\":1,\"initialization\":true,\"log\":\"\",\"tunnel_ip\":\"10.8.0.2/32\"\"rx_bytes\":57,\"tx_bytes\":60}}
```

---

## get_user_list

get ovpn user list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `user_list` | array | - | User list |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"get_user_list\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"user_list\":[{\"username\":\"test1\",\"password\":\"11111111\"},{\"username\":\"test2\",\"password\":\"22222222\"}]}}
```

---

## remove_route

remove ovpn server custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `route_flag` | number | yes | Route flag, 4 for IPv4 route, 6 for IPv6 route |
| `rule_id` | number | yes | Route rule ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"remove_route\",{\"rule_id\":6512,\"route_flag\":4}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## remove_user

remove ovpn user

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `username` | string | yes | Username |
| `password` | string | yes | Password |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"remove_user\",{\"username\":\"glinet\",\"password\":\"123456\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_config

Set ovpn server config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `mode` | string | yes | VPN mode: tun/tap |
| `subnetv4` | string | yes | IPv4 subnet |
| `mask` | string | yes | Mask |
| `subnetv6` | string | no | IPv6 subnet |
| `proto` | string | yes | Transport protocol |
| `auth` | string | no | Authentication method |
| `cipher` | string | no | Encryption method |
| `start` | string | no | IP pool start address in tap mode |
| `end` | string | no | IP pool end address in tap mode |
| `ca` | string | no | Server root certificate |
| `dh` | string | no | Diffie Hellman parameters |
| `ta` | string | no | TLS-Auth additional key |
| `verb` | string | no | Log level |
| `cert` | string | no | Server certificate |
| `key` | string | no | Server key |
| `tap_address` | string | no | Tap mode IP address |
| `tap_mask` | string | no | Tap mode mask |
| `ipv6_enable` | bool | no | Whether to enable IPv6 |
| `lzo` | bool | no | Whether to use LZO compression |
| `hmac` | bool | no | HMAC authentication |
| `client_to_client` | bool | no | Client-to-client communication |
| `port` | number | yes | Port |
| `access_scope` | number | no | Client using VPN 1: Access internal network (LAN) and external network (WAN) 2: Only access internal network (LAN) 3: Only access external network |
| `client_auth` | number | no | Client authentication method 1: Certificate only 2: Username/password only 3: Use both username/password and certificate |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-11: uci configuration missing -12: Invalid subnet -13: Invalid mask -14: Subnet and mask mismatch -15: Invalid port -16: Invalid IPv6 address -17: No network -18: ovpn configuration missing -19: Missing IPv6 -20: Port occupied) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"set_config\",{\"auth\":\"SHA256\",\"proto\":\"tcp\",\"port\":1194,\"subnetv4\":\"10.8.0.0\",\"ipv6_enable\":true,\"subnetv6\":\"fd00:db8:0:123::/64\",\"mask\":\"255.255.255.0\",\"cipher\":\"AES-256-GCM\",\"mode\":\"tun\",\"start\":\"10.8.0.2\",\"end\":\"10.8.0.100\",\"lzo\":false,\"ca\":\"-----BEGIN CERTIFICATE----- MIIDCzCCAfOgAwIBAgIUequqVhT2jU3wWJF988KVcIUyUtwwDQYJKoZIhvcNAQEF BQAwFTETMBEGA1UEAwwKT3BlblZQTiBDQTAeFw0yMTA2MTkwMjA4MzNaFw0zMTA2 MTcwMjA4MzNaMBUxEzARBgNVBAMMCk9wZW5WUE4gQ0EwggEiMA0GCSqGSIb3DQEB AQUAA4IBDwAwggEKAoIBAQCjLr7B1mc04E/yHqR8hv7Qng5WuaZVGYDfUoxxly0c 8qOj6LqkL5UwjtcQV2rFgucwCsTYW8/dRY/Fd6cRllNdvRuA1nuMi/JhECkmf9AB ckdXlyXiObeYNL1MofcOsjRkB1oxokOO7BN2uJ2sQqekBSc7DN+8b9w5R+gaVXY3 082bVHvIKPiwYWzpNbO3MfdBsmh9oL8CryAfDAv8aTaAeCRoHNy5hnzqWoHfv6xm MsuVfwb9Vw5MegJ011x+VMF/++o95IcL/UsMSCrhjstyaf94ucPwIi7NuD2lor6e RP0eqKryGui1r2coZweZjGuYqOnTwHZhJcjpgSK1RdbhAgMBAAGjUzBRMB0GA1Ud DgQWBBRUf/PRJrXI4eHeTBBldnTeJnrefTAfBgNVHSMEGDAWgBRUf/PRJrXI4eHe TBBldnTeJnrefTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBBQUAA4IBAQCS bAc7oMhLnAnHMPaFhXVDmvsgX7rK2BEKk/qqlHR2qAd7eSK3cteYX+Kayq+j8RmO PPsk4hXQNo+6tcaZSi05ACVYcCOkNw5dOZ8DDIo5swydamAZewp4ExsZQP4pdr40 8ikNle6Amt4psgxT0AVH56g4YO7qqdg5RMY14grIsBTL21j5OhdnwORPxstC2icW S5vzouDaIZ85IkiLaBku4fw8hBPkSmLRBexzeYwx/CoOs6olOEaKYoLy5TYFT+Ed 4Tvd7d0sw8vh1l1DiNtb5jGw7BsX4S8rzuB6o0GpxtxJprlyOPqqsyxwFNwhyAVk TFYWeM/vOcydA3zBBJbH -----END CERTIFICATE-----\",\"dh\":\"-----BEGIN DH PARAMETERS----- MIGHAoGBALUmuxsYkB+OLETWb+fz2WhAwa6PUMNzweXEA44rZ8Yv3PqfUIb/rPiW ZUwIq+ouAF1kMPnRYYhsqg5d4OeMv8R6SPlQUckpZrBELsJro9qGRf2nda9v5KAI 8DaorlIhKHpEoJ9xBz8h/DtQ85PX3UT96RZXA/SAyVAhdDpHkJ5rAgEC -----END DH PARAMETERS-----\",\"hmac\":true,\"ta\":\"-----BEGIN OpenVPN Static key V1----- 66301913fa0fa22457019b1a52e26aef 058c0405905f7e415615cc576bf1ca56 a0d12e27f181044c7ce4035afa5a7e30 1e02284208762d1ef8cd94b1989b8940 b534c469920a3012b440a4abad67af32 a7d80e0fc9447bd82b0d2ef2098a99c2 e7decaa1c5442b5ab6c0bbd413652814 c260fd76318d3abc6a65552fda44cb00 066420d3232f423f2810c861e0c3429f f566d0737d8e63c86dbf5bea5693dd51 1002a5e1e3aa95eddd388550e0022b04 8299cac19864821e3369a17e02cf02ef 0c41460146edba8748e9bccdd57cf8ac 1558f8d9a5b3430bf2229cad92625910 e09dce6311e834c2ca8c004f3a18c7d2 97188c096342674ff89c4bebda1e04db -----END OpenVPN Static key V1-----\",\"client_to_client\":false,\"verb\":\"3\",\"access_scope\":1,\"client_auth\":1,\"cert\":\"nil\",\"key\":\"nil\",\"tap_address\":\"10.8.0.1\",\"tap_mask\":\"255.255.255.255\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_route

set ovpn server custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dest` | string | no | destination address |
| `gateway` | string | no | gateway |
| `scope` | string | no | scope (host, link, global, or numbers greater than 0) |
| `route_flag` | number | yes | Route flag, 4 for IPv4 route, 6 for IPv6 route |
| `rule_id` | number | yes | Route rule ID |
| `metric` | number | no | metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"set_route\",{\"rule_id\":6512,\"route_flag\":4,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_setting

set ovpn server setting

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
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"set_setting\",{\"local_access\":true,\"masq\":true,\"mtu\":1460}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## start

start ovpn server

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-1: VPN conflict -2: openvpn package missing -3: certificate missing -4: configuration missing -5: startup script missing -6: uci configuration missing -7: service already started -8: port conflict) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"start\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## stop

Stop ovpn server

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-server\",\"stop\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---
