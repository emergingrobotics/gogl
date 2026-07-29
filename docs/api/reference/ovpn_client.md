# `ovpn_client`

OpenVPN Client API

23 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","ovpn_client","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`add_config`](#add_config) | - | add ovpn client config |
| [`add_group`](#add_group) | - | add ovpn client group |
| [`add_route`](#add_route) | - | add ovpn client custom route |
| [`check_config`](#check_config) | - | check uploaded ovpn config's parameter |
| [`clear_config_list`](#clear_config_list) | - | remove all ovpn client config of a group |
| [`confirm_config`](#confirm_config) | - | confirm ovpn client config |
| [`get_all_config_list`](#get_all_config_list) | - | get ovpn client config of all group |
| [`get_config_list`](#get_config_list) | - | get ovpn client config of a group |
| [`get_group_list`](#get_group_list) | - | get ovpn client group list |
| [`get_recommend_config`](#get_recommend_config) | - | get recommend config |
| [`get_route_list`](#get_route_list) | - | get ovpn client custom route list |
| [`get_setting`](#get_setting) | - | get ovpn client setting |
| [`get_status`](#get_status) | - | get ovpn client status |
| [`get_third_config`](#get_third_config) | - | get third provider's ovpn client config |
| [`remove_config`](#remove_config) | - | remove ovpn client config |
| [`remove_group`](#remove_group) | - | remove ovpn client group |
| [`remove_route`](#remove_route) | - | remove ovpn client custom route |
| [`set_config`](#set_config) | - | modify ovpn client config |
| [`set_group`](#set_group) | - | modify ovpn client group |
| [`set_route`](#set_route) | - | set ovpn client custom route |
| [`set_setting`](#set_setting) | - | set ovpn client setting |
| [`start`](#start) | - | start ovpn client |
| [`stop`](#stop) | - | stop ovpn client |

---

## add_config

add ovpn client config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Name |
| `mode` | string | yes | VPN mode: tun/tap |
| `proto` | string | yes | Transport protocol: tcp/udp |
| `remote` | string | yes | Server address/port, format: server address:port |
| `verb` | string | no | Log level |
| `auth` | string | no | Authentication algorithm |
| `cipher` | string | no | Encryption algorithm |
| `ca` | string | no | Server root certificate |
| `ta` | string | no | TLS-Auth additional key |
| `cert` | string | no | Client certificate |
| `key` | string | no | Client key |
| `username` | string | no | Username |
| `password` | string | no | Password |
| `random_remote` | bool | yes | Whether to randomly select remote server |
| `lzo` | bool | no | Whether to use lzo compression |
| `hmac` | bool | no | Whether to use additional HMAC verification |
| `group_id` | number | yes | Group ID |
| `client_auth` | number | yes | Client authentication method 1: Certificate only 2: Username/password only 3: Use both username/password and certificate |

| Results | Type | Required | Description |
|---|---|---|---|
| `client_id` | number | - | Client ID |
| `err_code` | number | - | Error code (-10: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"add_config\",{\"group_id\":3212,\"name\":\"AirVPN\",\"mode\":\"tun\",\"proto\":\"tcp\",\"remote\":\"bn76377.glddns.com:1194\",\"random_remote\":false,\"verb\":\"3\",\"auth\":\"SHA256\",\"cipher\":\"AES-256-GCM\",\"lzo\":false,\"ca\":\"-----BEGIN CERTIFICATE----- MIIDCzCCAfOgAwIBAgIUequqVhT2jU3wWJF988KVcIUyUtwwDQYJKoZIhvcNAQEF BQAwFTETMBEGA1UEAwwKT3BlblZQTiBDQTAeFw0yMTA2MTkwMjA4MzNaFw0zMTA2 MTcwMjA4MzNaMBUxEzARBgNVBAMMCk9wZW5WUE4gQ0EwggEiMA0GCSqGSIb3DQEB AQUAA4IBDwAwggEKAoIBAQCjLr7B1mc04E/yHqR8hv7Qng5WuaZVGYDfUoxxly0c 8qOj6LqkL5UwjtcQV2rFgucwCsTYW8/dRY/Fd6cRllNdvRuA1nuMi/JhECkmf9AB ckdXlyXiObeYNL1MofcOsjRkB1oxokOO7BN2uJ2sQqekBSc7DN+8b9w5R+gaVXY3 082bVHvIKPiwYWzpNbO3MfdBsmh9oL8CryAfDAv8aTaAeCRoHNy5hnzqWoHfv6xm MsuVfwb9Vw5MegJ011x+VMF/++o95IcL/UsMSCrhjstyaf94ucPwIi7NuD2lor6e RP0eqKryGui1r2coZweZjGuYqOnTwHZhJcjpgSK1RdbhAgMBAAGjUzBRMB0GA1Ud DgQWBBRUf/PRJrXI4eHeTBBldnTeJnrefTAfBgNVHSMEGDAWgBRUf/PRJrXI4eHe TBBldnTeJnrefTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBBQUAA4IBAQCS bAc7oMhLnAnHMPaFhXVDmvsgX7rK2BEKk/qqlHR2qAd7eSK3cteYX+Kayq+j8RmO PPsk4hXQNo+6tcaZSi05ACVYcCOkNw5dOZ8DDIo5swydamAZewp4ExsZQP4pdr40 8ikNle6Amt4psgxT0AVH56g4YO7qqdg5RMY14grIsBTL21j5OhdnwORPxstC2icW S5vzouDaIZ85IkiLaBku4fw8hBPkSmLRBexzeYwx/CoOs6olOEaKYoLy5TYFT+Ed 4Tvd7d0sw8vh1l1DiNtb5jGw7BsX4S8rzuB6o0GpxtxJprlyOPqqsyxwFNwhyAVk TFYWeM/vOcydA3zBBJbH -----END CERTIFICATE-----\",\"hmac\":false,\"ta\":\"-----BEGIN OpenVPN Static key V1----- 66301913fa0fa22457019b1a52e26aef 058c0405905f7e415615cc576bf1ca56 a0d12e27f181044c7ce4035afa5a7e30 1e02284208762d1ef8cd94b1989b8940 b534c469920a3012b440a4abad67af32 a7d80e0fc9447bd82b0d2ef2098a99c2 e7decaa1c5442b5ab6c0bbd413652814 c260fd76318d3abc6a65552fda44cb00 066420d3232f423f2810c861e0c3429f f566d0737d8e63c86dbf5bea5693dd51 1002a5e1e3aa95eddd388550e0022b04 8299cac19864821e3369a17e02cf02ef 0c41460146edba8748e9bccdd57cf8ac 1558f8d9a5b3430bf2229cad92625910 e09dce6311e834c2ca8c004f3a18c7d2 97188c096342674ff89c4bebda1e04db -----END OpenVPN Static key V1-----\",\"client_auth\":2,\"username\":\"glinet\",\"password\":\"123456\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"client_id\":1912}}
```

---

## add_group

add ovpn client group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_name` | string | yes | Group name |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters -11: Group name empty) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"add_group\",{\"group_name\":\"mullvad\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## add_route

add ovpn client custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dest` | string | yes | Destination address |
| `gateway` | string | no | Gateway |
| `scope` | string | no | scope (host, link, global, or a number greater than 0) |
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |
| `route_flag` | number | yes | Route flag 4 for IPv4 route, 6 for IPv6 route |
| `metric` | number | no | Metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"add_route\",{\"group_id\":3213,\"client_id\":5715,\"route_flag\":4,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## check_config

check uploaded ovpn config's parameter

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `filename` | string | yes | File name to be verified |
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `needauth` | array | - | List configurations that require authentication |
| `needinfo` | array | - | List configurations that require keys |
| `passed` | array | - | List configurations that do not require additional information |
| `unpassed` | array | - | List configurations that failed the check |
| `rootcert` | array | - | Server root certificate |
| `err_code` | number | - | Error code (-10: Missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"check_config\",{\"group_id\":2312,\"filename\":\"France-Paris_UDP53.ovpn\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"passed\": [\"AirVPN_US-Dallas-Texas_Mensa_TCP-443.ovpn\"], \"unpassed\": [\"AirVPN_US-Dallas-Texas_Vulpecula_TCP-443.ovpn\"], \"needauth\": [\"France-Paris_UDP53.ovpn\"], \"needinfo\": [], \"rootcert\": [\"mullvad_ca.crt\"]}}
```

---

## clear_config_list

remove all ovpn client config of a group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters -11: Client is running) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"clear_config_list\",{\"group_id\":3212}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## confirm_config

confirm ovpn client config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"confirm_config\",{\"group_id\":2312}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## get_all_config_list

get ovpn client config of all group

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `config_list` | array | - | Node information |
| `config_list.group_id` | number | - | Group ID |
| `config_list.group_name` | string | - | Group name |
| `config_list.group_type` | number | - | Group type 1: Default group 2: Manually added group |
| `config_list.auth_type` | number | - | Authentication type 0: No username/password or account required 1: Username/password 2: Account |
| `config_list.username` | string | - | Username |
| `config_list.password` | string | - | Password |
| `config_list.clients` | array | - | List all ovpn client configuration files |
| `config_list.clients.client_id` | number | - | Client ID |
| `config_list.clients.name` | string | - | Name |
| `config_list.clients.mode` | string | - | VPN mode: tun/tap |
| `config_list.clients.proto` | string | - | Transport protocol: tcp/udp |
| `config_list.clients.remote` | string | - | Server address/port, format: server address:port |
| `config_list.clients.random_remote` | bool | - | Whether to randomly select remote server |
| `config_list.clients.verb` | string | - | Log level |
| `config_list.clients.auth` | string | - | Authentication algorithm |
| `config_list.clients.cipher` | string | - | Encryption algorithm |
| `config_list.clients.lzo` | bool | - | Whether to use LZO compression |
| `config_list.clients.ca` | string | - | Server root certificate |
| `config_list.clients.hmac` | bool | - | Whether to use additional HMAC verification |
| `config_list.clients.ta` | string | - | TLS-Auth additional key |
| `config_list.clients.client_auth` | number | - | Client authentication method 1: Certificate only 2: Username/password only 3: Use both username/password and certificate |
| `config_list.clients.cert` | string | - | Client certificate |
| `config_list.clients.key` | string | - | Client key |
| `config_list.clients.username` | string | - | Username |
| `config_list.clients.password` | string | - | Password |
| `config_list.clients.location` | string | - | Server location |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"get_all_config_list\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\":{\"config_list\":[{\"group_id\":2312,\"group_name\":\"mullvad-ae-all\",\"group_type\":2,\"auth_type\":1,\"username\":\"test\",\"password\":\"123456\",\"clients\":[{\"client_id\":3212,\"name\":\"AirVPN\",\"mode\":\"tun\",\"proto\":\"tcp\",\"remote\":\"bn76377.glddns.com:1194\",\"random_remote\":false,\"verb\":\"3\",\"auth\":\"SHA256\",\"cipher\":\"AES-256-GCM\",\"lzo\":false,\"ca\":\"-----BEGIN CERTIFICATE----- MIIDCzCCAfOgAwIBAgIUequqVhT2jU3wWJF988KVcIUyUtwwDQYJKoZIhvcNAQEF BQAwFTETMBEGA1UEAwwKT3BlblZQTiBDQTAeFw0yMTA2MTkwMjA4MzNaFw0zMTA2 MTcwMjA4MzNaMBUxEzARBgNVBAMMCk9wZW5WUE4gQ0EwggEiMA0GCSqGSIb3DQEB AQUAA4IBDwAwggEKAoIBAQCjLr7B1mc04E/yHqR8hv7Qng5WuaZVGYDfUoxxly0c 8qOj6LqkL5UwjtcQV2rFgucwCsTYW8/dRY/Fd6cRllNdvRuA1nuMi/JhECkmf9AB ckdXlyXiObeYNL1MofcOsjRkB1oxokOO7BN2uJ2sQqekBSc7DN+8b9w5R+gaVXY3 082bVHvIKPiwYWzpNbO3MfdBsmh9oL8CryAfDAv8aTaAeCRoHNy5hnzqWoHfv6xm MsuVfwb9Vw5MegJ011x+VMF/++o95IcL/UsMSCrhjstyaf94ucPwIi7NuD2lor6e RP0eqKryGui1r2coZweZjGuYqOnTwHZhJcjpgSK1RdbhAgMBAAGjUzBRMB0GA1Ud DgQWBBRUf/PRJrXI4eHeTBBldnTeJnrefTAfBgNVHSMEGDAWgBRUf/PRJrXI4eHe TBBldnTeJnrefTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBBQUAA4IBAQCS bAc7oMhLnAnHMPaFhXVDmvsgX7rK2BEKk/qqlHR2qAd7eSK3cteYX+Kayq+j8RmO PPsk4hXQNo+6tcaZSi05ACVYcCOkNw5dOZ8DDIo5swydamAZewp4ExsZQP4pdr40 8ikNle6Amt4psgxT0AVH56g4YO7qqdg5RMY14grIsBTL21j5OhdnwORPxstC2icW S5vzouDaIZ85IkiLaBku4fw8hBPkSmLRBexzeYwx/CoOs6olOEaKYoLy5TYFT+Ed 4Tvd7d0sw8vh1l1DiNtb5jGw7BsX4S8rzuB6o0GpxtxJprlyOPqqsyxwFNwhyAVk TFYWeM/vOcydA3zBBJbH -----END CERTIFICATE-----\",\"hmac\":false,\"ta\":\"-----BEGIN OpenVPN Static key V1----- 66301913fa0fa22457019b1a52e26aef 058c0405905f7e415615cc576bf1ca56 a0d12e27f181044c7ce4035afa5a7e30 1e02284208762d1ef8cd94b1989b8940 b534c469920a3012b440a4abad67af32 a7d80e0fc9447bd82b0d2ef2098a99c2 e7decaa1c5442b5ab6c0bbd413652814 c260fd76318d3abc6a65552fda44cb00 066420d3232f423f2810c861e0c3429f f566d0737d8e63c86dbf5bea5693dd51 1002a5e1e3aa95eddd388550e0022b04 8299cac19864821e3369a17e02cf02ef 0c41460146edba8748e9bccdd57cf8ac 1558f8d9a5b3430bf2229cad92625910 e09dce6311e834c2ca8c004f3a18c7d2 97188c096342674ff89c4bebda1e04db -----END OpenVPN Static key V1-----\",\"client_auth\":2,\"username\":\"glinet\",\"password\":\"123456\",\"location\":\"United States, Chicago; United States, New York\"}]}]}}
```

---

## get_config_list

get ovpn client config of a group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `clients` | array | - | List all ovpn client configuration files |
| `clients.client_id` | number | - | Client ID |
| `clients.name` | string | - | Name |
| `clients.mode` | string | - | VPN mode: tun/tap |
| `clients.proto` | string | - | Transport protocol: tcp/udp |
| `clients.remote` | string | - | Server address/port, format: server address:port |
| `clients.random_remote` | bool | - | Whether to randomly select remote server |
| `clients.verb` | string | - | Log level |
| `clients.auth` | string | - | Authentication algorithm |
| `clients.cipher` | string | - | Encryption algorithm |
| `clients.lzo` | bool | - | Whether to use LZO compression |
| `clients.ca` | string | - | Server root certificate |
| `clients.hmac` | bool | - | Whether to use additional HMAC verification |
| `clients.ta` | string | - | TLS-Auth additional key |
| `clients.client_auth` | number | - | Client authentication method 1: Certificate only 2: Username/password only 3: Use both username/password and certificate |
| `clients.cert` | string | - | Client certificate |
| `clients.key` | string | - | Client key |
| `clients.username` | string | - | Username |
| `clients.password` | string | - | Password |
| `clients.location` | string | - | Server location |
| `err_code` | number | - | Error code (-10: Missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"get_config_list\",{\"group_id\":2312}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\": 1,\"result\":{\"clients\":[{\"client_id\":3212,\"name\":\"AirVPN\",\"mode\":\"tun\",\"proto\":\"tcp\",\"remote\":\"bn76377.glddns.com:1194\",\"random_remote\":false,\"verb\":\"3\",\"auth\":\"SHA256\",\"cipher\":\"AES-256-GCM\",\"lzo\":false,\"ca\":\"-----BEGIN CERTIFICATE----- MIIDCzCCAfOgAwIBAgIUequqVhT2jU3wWJF988KVcIUyUtwwDQYJKoZIhvcNAQEF BQAwFTETMBEGA1UEAwwKT3BlblZQTiBDQTAeFw0yMTA2MTkwMjA4MzNaFw0zMTA2 MTcwMjA4MzNaMBUxEzARBgNVBAMMCk9wZW5WUE4gQ0EwggEiMA0GCSqGSIb3DQEB AQUAA4IBDwAwggEKAoIBAQCjLr7B1mc04E/yHqR8hv7Qng5WuaZVGYDfUoxxly0c 8qOj6LqkL5UwjtcQV2rFgucwCsTYW8/dRY/Fd6cRllNdvRuA1nuMi/JhECkmf9AB ckdXlyXiObeYNL1MofcOsjRkB1oxokOO7BN2uJ2sQqekBSc7DN+8b9w5R+gaVXY3 082bVHvIKPiwYWzpNbO3MfdBsmh9oL8CryAfDAv8aTaAeCRoHNy5hnzqWoHfv6xm MsuVfwb9Vw5MegJ011x+VMF/++o95IcL/UsMSCrhjstyaf94ucPwIi7NuD2lor6e RP0eqKryGui1r2coZweZjGuYqOnTwHZhJcjpgSK1RdbhAgMBAAGjUzBRMB0GA1Ud DgQWBBRUf/PRJrXI4eHeTBBldnTeJnrefTAfBgNVHSMEGDAWgBRUf/PRJrXI4eHe TBBldnTeJnrefTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBBQUAA4IBAQCS bAc7oMhLnAnHMPaFhXVDmvsgX7rK2BEKk/qqlHR2qAd7eSK3cteYX+Kayq+j8RmO PPsk4hXQNo+6tcaZSi05ACVYcCOkNw5dOZ8DDIo5swydamAZewp4ExsZQP4pdr40 8ikNle6Amt4psgxT0AVH56g4YO7qqdg5RMY14grIsBTL21j5OhdnwORPxstC2icW S5vzouDaIZ85IkiLaBku4fw8hBPkSmLRBexzeYwx/CoOs6olOEaKYoLy5TYFT+Ed 4Tvd7d0sw8vh1l1DiNtb5jGw7BsX4S8rzuB6o0GpxtxJprlyOPqqsyxwFNwhyAVk TFYWeM/vOcydA3zBBJbH -----END CERTIFICATE-----\",\"hmac\":false,\"ta\":\"-----BEGIN OpenVPN Static key V1----- 66301913fa0fa22457019b1a52e26aef 058c0405905f7e415615cc576bf1ca56 a0d12e27f181044c7ce4035afa5a7e30 1e02284208762d1ef8cd94b1989b8940 b534c469920a3012b440a4abad67af32 a7d80e0fc9447bd82b0d2ef2098a99c2 e7decaa1c5442b5ab6c0bbd413652814 c260fd76318d3abc6a65552fda44cb00 066420d3232f423f2810c861e0c3429f f566d0737d8e63c86dbf5bea5693dd51 1002a5e1e3aa95eddd388550e0022b04 8299cac19864821e3369a17e02cf02ef 0c41460146edba8748e9bccdd57cf8ac 1558f8d9a5b3430bf2229cad92625910 e09dce6311e834c2ca8c004f3a18c7d2 97188c096342674ff89c4bebda1e04db -----END OpenVPN Static key V1-----\",\"client_auth\":2,\"username\":\"glinet\",\"password\":\"123456\",\"location\":\"United States, Chicago; United States, New York\"}]}
```

---

## get_group_list

get ovpn client group list

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `groups` | array | - | Group nodes |
| `groups.group_id` | number | - | Group ID |
| `groups.group_type` | number | - | Group type 1: Default group 2: Manually added group 3: App group |
| `groups.group_name` | string | - | Group name |
| `groups.askpass_exist` | bool | - | Whether the certificate file key exists |
| `groups.askpass` | string | - | Certificate file key |
| `groups.auth_type` | number | - | Authentication type 0: No username/password or account required 1: Username/password 2: Account |
| `groups.client_count` | number | - | Number of configurations in the group |
| `groups.username` | string | - | Username |
| `groups.password` | string | - | Password |
| `groups.procedure` | number | - | Configuration download process 0: Direct download 1: Step-by-step download |
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"get_group_list\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"groups\": [{\"group_id\":3212,\"group_type\":1,\"group_name\":\"mullvad\",\"askpass_exist\":false,\"auth_type\":1,\"client_count\":20,\"procedure\":1},{\"group_id\":3213,\"group_type\"1,\"group_name\":\"azirevpn\",\"client_count\":25,\"procedure\":0}]}}
```

---

## get_recommend_config

get recommend config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `country_id` | array | no | Country ID |
| `servers` | array | no | Server information |
| `group_id` | number | yes | Group ID |
| `proto` | number | yes | Protocol 0:tcp and udp 1:tcp 2:udp |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters -11: UCI configuration missing -12: Group information missing -13: Missing password -14: Country information empty -15: Unknown service provider) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"get_recommend_config\",{\"group_id\": 2312, \"country_id\":[2,28], \"proto\": 1, \"servers\":[{\"country_name\":\"Germany\",\"city_name\":\"Frankfurt\",\"hostname\":[\"dk205.nordvpn.com\",\"cz150.nordvpn.com\"]}]}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## get_route_list

get ovpn client custom route list

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `ipv4_route_rules` | array | - | IPv4 route rules |
| `ipv4_route_rules.rule_id` | string | - | Route rule ID |
| `ipv4_route_rules.dest` | string | - | Destination address |
| `ipv4_route_rules.gateway` | string | - | Gateway |
| `ipv4_route_rules.metric` | number | - | Metric |
| `ipv4_route_rules.mtu` | number | - | MTU |
| `ipv4_route_rules.scope` | string | - | scope (host, link, global, or a number greater than 0) |
| `ipv6_route_rules` | array | - | IPv4 route rules |
| `ipv6_route_rules.rule_id` | string | - | Route rule ID |
| `ipv6_route_rules.dest` | string | - | Destination address |
| `ipv6_route_rules.gateway` | string | - | Gateway |
| `ipv6_route_rules.metric` | number | - | Metric |
| `ipv6_route_rules.mtu` | number | - | MTU |
| `ipv6_route_rules.scope` | string | - | scope (host, link, global, or a number greater than 0) |
| `err_code` | number | - | Error code (-10: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"get_route_list\",{\"group_id\":3213,\"client_id\":5715}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"ipv4_route_rules\":[{\"rule_id\":6512,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"ipv6_route_rules\":[{\"rule_id\":6515,\"dest\":\"ff00::/8\",\"gateway\":\"ff00::1\",\"metric\":20,\"mtu\":1500}]}}
```

---

## get_setting

get ovpn client setting

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `local_access` | bool | - | Whether to enable local access |
| `masq` | bool | - | Whether to enable masq |
| `default_metric` | number | - | Default metric |
| `mtu` | number | - | MTU |
| `err_code` | number | - | Error code (-10: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"get_setting\",{\"group_id\":3213,\"client_id\":5715}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"local_access\":true,\"masq\":true,\"default_metric\":30,\"mtu\":1460}}
```

---

## get_status

get ovpn client status

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | - | Client's group ID |
| `client_id` | number | - | Client ID |
| `name` | string | - | Client name |
| `log` | string | - | Log |
| `status` | number | - | Connection status 0: Not connected 1: Connected successfully 2: Connecting |
| `domain` | string | - | Server domain |
| `port` | number | - | Server listening port |
| `ipv4` | string | - | IPv4 address |
| `ipv6` | string | - | IPv6 address |
| `tx_bytes` | number | - | Transmitted bytes |
| `rx_bytes` | number | - | Received bytes |
| `err_code` | string | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"get_status\",{}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"group_id\":3212, \"client_id\":1254, \"name\":\"mullvad_ae_all\", \"log\":\"\", \"status\":1, \"tx_bytes\":56, \"rx_bytes\":50, \"ipv4\":\"10.70.33.3\", \"ipv6\":\"fc00:bbbb:bb01::11\", \"port\":1194, \"domain\":\"ae-dxb-001.mullvad.net\"}}
```

---

## get_third_config

get third provider's ovpn client config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `username` | string | yes | Username |
| `password` | string | no | Password. If authentication type is username/password, this field is required |
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `provider` | string | - | Provider |
| `server_info` | object | - | Server information |
| `err_code` | number | - | Error code (-10: Missing parameters -11: UCI configuration missing -12: Group information missing -13: Missing password -14: Country information empty -15: Unknown service provider) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"get_third_config\",{\"group_id\":2312,\"username\":\"glinet\",\"password\":\"123456\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {\"provider\":\"nordvpn\",\"server_info\":\"\"}}
```

---

## remove_config

remove ovpn client config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters -11: Client is running) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"remove_config\",{\"group_id\":3212,\"client_id\":1254}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## remove_group

remove ovpn client group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters -11: Client is running) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"remove_group\",{\"group_id\":3212}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## remove_route

remove ovpn client custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |
| `route_flag` | number | yes | Route flag, 4 for IPv4 route, 6 for IPv6 route |
| `rule_id` | number | yes | Route rule ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"remove_route\",{\"group_id\":3213,\"client_id\":5715,\"route_flag\":4,\"rule_id\":6512}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_config

modify ovpn client config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Name |
| `mode` | string | no | VPN mode: tun/tap |
| `proto` | string | no | Transport protocol: tcp/udp |
| `remote` | string | no | Server address/port, format: server address:port |
| `verb` | string | no | Log level |
| `auth` | string | no | Authentication algorithm |
| `cipher` | string | no | Encryption algorithm |
| `ca` | string | no | Server root certificate |
| `ta` | string | no | TLS-Auth additional key |
| `cert` | string | no | Client certificate |
| `key` | string | no | Client key |
| `username` | string | no | Username |
| `password` | string | no | Password |
| `random_remote` | bool | no | Whether to randomly select remote server |
| `lzo` | bool | no | Whether to use lzo compression |
| `hmac` | bool | no | Whether to use additional HMAC verification |
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |
| `client_auth` | number | no | Client authentication method 1: Certificate only 2: Username/password only 3: Use both username/password and certificate |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"set_config\",{\"group_id\":3212,\"client_id\":1254,\"name\":\"AirVPN\",\"mode\":\"tun\",\"proto\":\"tcp\",\"remote\":\"bn76377.glddns.com:1194\",\"random_remote\":false,\"verb\":\"3\",\"auth\":\"SHA256\",\"cipher\":\"AES-256-GCM\",\"lzo\":false,\"ca\":\"-----BEGIN CERTIFICATE----- MIIDCzCCAfOgAwIBAgIUequqVhT2jU3wWJF988KVcIUyUtwwDQYJKoZIhvcNAQEF BQAwFTETMBEGA1UEAwwKT3BlblZQTiBDQTAeFw0yMTA2MTkwMjA4MzNaFw0zMTA2 MTcwMjA4MzNaMBUxEzARBgNVBAMMCk9wZW5WUE4gQ0EwggEiMA0GCSqGSIb3DQEB AQUAA4IBDwAwggEKAoIBAQCjLr7B1mc04E/yHqR8hv7Qng5WuaZVGYDfUoxxly0c 8qOj6LqkL5UwjtcQV2rFgucwCsTYW8/dRY/Fd6cRllNdvRuA1nuMi/JhECkmf9AB ckdXlyXiObeYNL1MofcOsjRkB1oxokOO7BN2uJ2sQqekBSc7DN+8b9w5R+gaVXY3 082bVHvIKPiwYWzpNbO3MfdBsmh9oL8CryAfDAv8aTaAeCRoHNy5hnzqWoHfv6xm MsuVfwb9Vw5MegJ011x+VMF/++o95IcL/UsMSCrhjstyaf94ucPwIi7NuD2lor6e RP0eqKryGui1r2coZweZjGuYqOnTwHZhJcjpgSK1RdbhAgMBAAGjUzBRMB0GA1Ud DgQWBBRUf/PRJrXI4eHeTBBldnTeJnrefTAfBgNVHSMEGDAWgBRUf/PRJrXI4eHe TBBldnTeJnrefTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBBQUAA4IBAQCS bAc7oMhLnAnHMPaFhXVDmvsgX7rK2BEKk/qqlHR2qAd7eSK3cteYX+Kayq+j8RmO PPsk4hXQNo+6tcaZSi05ACVYcCOkNw5dOZ8DDIo5swydamAZewp4ExsZQP4pdr40 8ikNle6Amt4psgxT0AVH56g4YO7qqdg5RMY14grIsBTL21j5OhdnwORPxstC2icW S5vzouDaIZ85IkiLaBku4fw8hBPkSmLRBexzeYwx/CoOs6olOEaKYoLy5TYFT+Ed 4Tvd7d0sw8vh1l1DiNtb5jGw7BsX4S8rzuB6o0GpxtxJprlyOPqqsyxwFNwhyAVk TFYWeM/vOcydA3zBBJbH -----END CERTIFICATE-----\",\"hmac\":false,\"ta\":\"-----BEGIN OpenVPN Static key V1----- 66301913fa0fa22457019b1a52e26aef 058c0405905f7e415615cc576bf1ca56 a0d12e27f181044c7ce4035afa5a7e30 1e02284208762d1ef8cd94b1989b8940 b534c469920a3012b440a4abad67af32 a7d80e0fc9447bd82b0d2ef2098a99c2 e7decaa1c5442b5ab6c0bbd413652814 c260fd76318d3abc6a65552fda44cb00 066420d3232f423f2810c861e0c3429f f566d0737d8e63c86dbf5bea5693dd51 1002a5e1e3aa95eddd388550e0022b04 8299cac19864821e3369a17e02cf02ef 0c41460146edba8748e9bccdd57cf8ac 1558f8d9a5b3430bf2229cad92625910 e09dce6311e834c2ca8c004f3a18c7d2 97188c096342674ff89c4bebda1e04db -----END OpenVPN Static key V1-----\",\"client_auth\":2,\"username\":\"glinet\",\"password\":\"123456\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_group

modify ovpn client group

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_name` | string | yes | Group name |
| `username` | string | no | Username |
| `password` | string | no | Password |
| `askpass` | string | no | Certificate file key |
| `group_id` | number | yes | Group ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"set_group\",{\"group_id\":3212,\"group_name\":\"mullvad\",\"username\":\"glinet\",\"password\":\"123456\"}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_route

set ovpn client custom route

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `dest` | string | no | Destination address |
| `gateway` | string | no | Gateway |
| `scope` | string | no | scope (host, link, global, or a number greater than 0) |
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |
| `route_flag` | number | yes | Route flag 4 for IPv4 route, 6 for IPv6 route |
| `rule_id` | number | yes | Route rule ID |
| `metric` | number | no | Metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameters) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"set_route\",{\"group_id\":3213,\"client_id\":5715,\"route_flag\":4,\"rule_id\":6512,\"dest\":\"192.168.9.13/32\",\"gateway\":\"10.8.3.1\",\"metric\":20,\"mtu\":1500}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## set_setting

set ovpn client setting

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `local_access` | bool | no | Whether to enable local access |
| `masq` | bool | no | Whether to enable masq |
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |
| `default_metric` | number | no | Default metric |
| `mtu` | number | no | MTU |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: missing parameter) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"set_setting\",{\"group_id\":3213,\"client_id\":5715,\"local_access\":true,\"masq\":true,\"default_metric\":30,\"mtu\":1460}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## start

start ovpn client

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code (-10: Missing parameter -11: Global mode client conflict) |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"start\",{\"group_id\":3212,\"client_id\":1254}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---

## stop

stop ovpn client

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `group_id` | number | yes | Group ID |
| `client_id` | number | yes | Client ID |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"method\":\"call\",\"params\":[\"\",\"ovpn-client\",\"stop\",{\"group_id\":3212,\"client_id\":1254}],\"id\":1}
```

**Response**

```json
{\"jsonrpc\": \"2.0\", \"id\": 1, \"result\": {}}
```

---
