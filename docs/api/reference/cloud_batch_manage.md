# `cloud_batch_manage`

This is the API related to cloud batch manage.

7 method(s). Call shape:

```json
{"jsonrpc":"2.0","id":1,"method":"call","params":["<sid>","cloud_batch_manage","<method>",{}]}
```

| Method | Status | Description |
|---|---|---|
| [`bind_info`](#bind_info) | - | set cloud bind info |
| [`designated_customer`](#designated_customer) | - | set customer info |
| [`get_2b_config`](#get_2b_config) | - | get 2b config |
| [`get_batch_config`](#get_batch_config) | - | get batch config |
| [`send_router_info`](#send_router_info) | - | device send router info |
| [`set_2b_config`](#set_2b_config) | - | set 2b config |
| [`set_batch_config`](#set_batch_config) | - | set batch config |

---

## bind_info

set cloud bind info

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `username` | string | yes | cloud user name |
| `email` | string | yes | cloud user email |
| `bindtime` | string | yes | router device in cloud bind time |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud-batch-manage\",\"bind_info\",{\"username\":\"glcloud\",\"email\":\"126@qq.com\",\"bindtime\":\"192168\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## designated_customer

set customer info

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `server_name` | string | yes | server name |
| `server_url` | string | yes | server url |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud-batch-manage\",\"designated_customer\",{\"server_name\":\"glcloud\",\"server_url\":\"http://xxx\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## get_2b_config

get 2b config

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `lan_ip` | string | - | device lan ip |
| `autoupdate_firmeare_path` | string | - | device autoupdate firmare path |
| `autoupdate_time` | string | - | device autoupdate time |
| `autoupdate_enable` | bool | - | enable autoupdate |
| `password` | string | - | system password |
| `system_timezone` | string | - | system timezone |
| `hostname` | string | - | hostname |
| `auto_timezone` | string | - | auto sync timezone |
| `server_name` | string | - | server name |
| `server_url` | string | - | server url |
| `language` | string | - | languate |
| `image_url` | string | - | image url |
| `customer_name` | string | - | customer name |
| `help_url` | string | - | help url |
| `ssid_2g` | string | - | 2.4G ssid |
| `disable_2g` | bool | - | 2.4g disable |
| `hidden_2g` | bool | - | 2.4G hidden |
| `channel_2g` | string | - | 2.4G channle |
| `txpower_2g` | string | - | 2.4G txpower |
| `htmode_2g` | string | - | 2.4g htmode |
| `encrytion_2g` | string | - | 2.4G encryption |
| `key_2g` | string | - | 2.4G key |
| `guest_2g_ssid` | string | - | 2.4G guest ssid |
| `guest_2g_encrytion` | string | - | 2.4G guest encryption |
| `guest_2g_key` | string | - | 2.4G guest key |
| `guest_2g_disable` | bool | - | 2.4G guest disable |
| `ssid_5g` | string | - | 5G ssid |
| `channel_5g` | string | - | 5G channle |
| `txpower_5g` | string | - | 5G txpower |
| `htmode_5g` | string | - | 5G htmode |
| `encrytion_5g` | string | - | 5G encryption |
| `disable_5g` | bool | - | 5G disable |
| `hidden_5g` | bool | - | 5G hidden |
| `key_5g` | string | - | 5G key |
| `guest_5g_ssid` | string | - | 5G guest ssid |
| `guest_5g_encrytion` | string | - | 5G guest encryption |
| `guest_5g_key` | string | - | 5G guest key |
| `guest_5g_disable` | bool | - | 5G guest disable |
| `reboot_flag` | string | - | sysem reboot flag |
| `private_disabled` | string | - | fb private disabled |
| `private_controller` | string | - | fb private controller |
| `guest_disabled` | bool | - | fb guest disable |
| `guest_controller` | string | - | fb guest controller |
| `mode` | string | - | fb mode |
| `ssl_id` | string | - | fb ssl id |
| `private_wired` | string | - | fb private wired |
| `private_datapath_id` | string | - | fb pravate data path id |
| `private_datapath_desc` | string | - | fb private data path describte |
| `private_ipaddr` | string | - | fb private ipaddr |
| `private_netmask` | string | - | fb private dhcp lease |
| `private_dhcp_lease` | string | - | fb private wlans |
| `private_wlans` | string | - | fb private wlans |
| `guest_wired` | string | - | fb guest wired |
| `guest_datapath_id` | string | - | fb guest datapath id |
| `guest_datapath_desc` | string | - | fb guest datapath desc |
| `guest_ipaddr` | string | - | fb guest ipaddr |
| `guest_netmask` | string | - | fb guest netmask |
| `guest_dhcp_lease` | string | - | fb dhcp lease |
| `guest_wlans` | string | - | gb guest wlans |
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud-batch-manage\",\"get_2b_config\"]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"lan_ip:\"192.168.8.1\",\"autoupdate_firmware_path:\"http://xxx\",\"autoupdate_time:\"128169\",\"autoupdate_enable:\"true\",\"Password:\"abc\",\"system_timezone:\"china\",\"hostname:\"ABC\",\"auto_timezone:\"china\",\"server_name:\"GL\",\"server_url:\"http://xxx\",\"language:\"ZH\",\"image_url:\"http://xxx\",\"customer_name:\"ABC\",\"help_url:\"http://xxx\",\"ssid_2g:\"2G-wifi\",\"disable_2g:\"true\",\"hidden_2g:\"true\",\"channel_2g:\"9\",\"txpower_2g:\"9\",\"htmode_2g:\"30\",\"encrytion_2g:\"none\",\"key_2g:\"123456\",\"guest_2g_ssid:\"2G-guest\",\"guest_2g_encrytion:\"none\",\"guest_2g_key:\"123456\",\"guest_2g_disable:\"true\",\"ssid_5g:\"5G-wifi\",\"channel_5g:\"30\",\"txpower_5g:\"30\",\"htmode_5g:\"30\",\"encrytion_5g:\"none\",\"disable_5g:\"true\",\"hidden_5g:\"false\",\"key_5g:\"123456\",\"guest_5g_ssid:\"5G-wifi\",\"guest_5g_encrytion:\"none\",\"guest_5g_key:\"123456\",\"guest_5g_disable:\"true\",\"reboot_flag:\"false\",\"private_disabled:\"true\",\"private_controller:\"192.168.8.1\",\"guest_disabled:\"flase\",\"guest_controller:\"192.168.8.1\",\"mode:\"abc\",\"ssl_id:\"XXX\",\"private_wired:\"abc\",\"private_datapath_id:\"ABC\",\"private_datapath_desc:\"xxx\",\"private_ipaddr:\"192.168.8.1\",\"private_netmask:\"255.255.255.0\",\"private_dhcp_lease:\"192.168.9.100\",\"private_wlans:\"AAA\",\"guest_wired:\"XXX\",\"guest_datapath_id:\"YYY\",\"guest_datapath_desc:\"ZZZ\",\"guest_ipaddr:\"192.168.8.1\"\"guest_netmask:\"255.255.255.0\",\"guest_dhcp_lease:\"192.168.9.100\",\"guest_wlans:\"192.168.8.26\"}}
```

---

## get_batch_config

get batch config

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `wan_proto` | string | - | wan proto |
| `wan_username` | string | - | wan proto username |
| `wan_password` | string | - | wan proto password |
| `wan_ipaddr` | string | - | wan ipaddr |
| `wan_gateway` | string | - | wan gateway |
| `wan_netmask` | string | - | wan netmask |
| `wan_dns` | string | - | wan dns |
| `reboot_flag` | bool | - | reboot_flag |
| `ssid_2g` | string | - | 2.4G ssid |
| `disable_2g` | bool | - | 2.4g disable |
| `hidden_2g` | bool | - | 2.4G hidden |
| `channel_2g` | string | - | 2.4G channle |
| `txpower_2g` | string | - | 2.4G txpower |
| `htmode_2g` | string | - | 2.4g htmode |
| `encrytion_2g` | string | - | 2.4G encryption |
| `key_2g` | string | - | 2.4G key |
| `guest_2g_ssid` | string | - | 2.4G guest ssid |
| `guest_2g_encrytion` | string | - | 2.4G guest encryption |
| `guest_2g_key` | string | - | 2.4G guest key |
| `guest_2g_disable` | bool | - | 2.4G guest disable |
| `ssid_5g` | string | - | 5G ssid |
| `channel_5g` | string | - | 5G channle |
| `txpower_5g` | string | - | 5G txpower |
| `htmode_5g` | string | - | 5G htmode |
| `encrytion_5g` | string | - | 5G encryption |
| `disable_5g` | bool | - | 5G disable |
| `hidden_5g` | bool | - | 5G hidden |
| `key_5g` | string | - | 5G key |
| `guest_5g_ssid` | string | - | 5G guest ssid |
| `guest_5g_encrytion` | string | - | 5G guest encryption |
| `guest_5g_key` | string | - | 5G guest key |
| `guest_5g_disable` | bool | - | 5G guest disable |
| `probe_data_path1` | string | - | probe data path |
| `probe_data_path2` | string | - | probe data path |
| `probe_max_lines` | string | - | probe max lines |
| `probe_interval` | string | - | probe interval |
| `probe_json` | string | - | probe json |
| `probe_virmac` | string | - | probe virmac |
| `probe_type` | string | - | probe type |
| `probe_tcpurl` | string | - | probe tcp url |
| `probe_filter` | string | - | probe filter |
| `probe_sflag` | string | - | probe sflag |
| `probe_protocol` | string | - | probe protocol |
| `probe_channel` | string | - | probe channel |
| `probe_channel5g` | string | - | probe 5G channel |
| `probe_induce` | string | - | probe induce |
| `lan_ip` | string | - | lan ip |
| `guest_lan_ip` | string | - | guest lan ip |
| `autoupdate_firmware_path` | string | - | autoupdate firmware path |
| `autoupdate_time` | string | - | autoupdate time |
| `autoupdate_enable` | bool | - | atuoupdate enable |
| `password` | string | - | system password |
| `system_timezone` | string | - | system timezone |
| `auto_timezone` | string | - | auto timezone |
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud-batch-manage\",\"get_batch_config\"]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"wan_proto\":\"dhcp\",\"wan_username\":\"glcloud\",\"wan_password\":\"abc\",\"wan_ipaddr\":\"192.168.8.1\",\"wan_gateway\":\"255.255.255.1\",\"wan_netmask\":\"255.255.255.0\"\"wan_dns\":\"8.8.8.8\",\"reboot_flag\":\"true\",\"ssid_2g\":\"GL\",\"disable_2g\":\"true\",\"hidden_2g\":\"true\",\"channel_2g\":\"8\",\"txpower_2g\":\"10\",\"htmode_2g\":\"10\",\"encrytion_2g\":\"none\",\"key_2g\":\"abc\",\"guest_2g_ssid\":\"GL-guest\",\"guest_2g_encrytion\":\"none\",\"guest_2g_key\":\"123456\",\"guest_2g_disable\":\"true\",\"ssid_5g\":\"GL-5G\",\"channel_5g\":\"8\",\"txpower_5g\":\"30\",\"htmode_5g\":\"30\",\"encrytion_5g\":\"none\",\"disable_5g\":\"true\",\"hidden_5g\":\"false\",\"key_5g\":\"abc\",\"guest_5g_ssid\":\"GL-5Guest\",\"guest_5g_encrytion\":\"none\",\"guest_5g_key\":\"abc123\",\"guest_5g_disable\":\"true\",\"probe_data_path1\":\"XXX\",\"probe_data_path2\":\"YYY\",\"probe_max_lines\":\"30\",\"probe_interval\":\"100\",\"probe_json\":\"json\",\"probe_virmac\":\"AABBCCEEDDFF\",\"probe_type\":\"3\",\"probe_tcpurl\":\"http://xxx\",\"probe_filter\":\"abc\",\"probe_sflag\":\"XYZ\",\"probe_protocol\":\"DHCP\",\"probe_channel\":\"10\",\"probe_channel5g\":\"30\",\"probe_induce\":\"ABC\",\"lan_ip\":\"192.168.8.1\",\"guest_lan_ip\":\"192.168.9.1\",\"autoupdate_firmware_path\":\"http://xxx\",\"autoupdate_time\":\"192168\",\"autoupdate_enable\":\"false\",\"Password\":\"abc\",\"system_timezone\":\"china\",\"auto_timezone\":\"china\"}}
```

---

## send_router_info

device send router info

Not exercised against hardware here.

_No params._

| Results | Type | Required | Description |
|---|---|---|---|
| `model` | string | - | device model |
| `rtty_ssh` | bool | - | enable rtty ssh |
| `rtty_web` | bool | - | enable rtty web |
| `ddns` | string | - | device ddns |
| `sn` | string | - | device sn |
| `mac` | string | - | device mac |
| `fw_type` | string | - | framewrok type |
| `version` | string | - | framewrok version |
| `dataupload` | bool | - | enable data upload |
| `mode` | string | - | bridge mode |
| `mesh` | string | - | mesh type |
| `services` | array | - | server address |
| `firmware_path` | string | - | firmware update path |
| `data_path` | string | - | wifi probe data address |
| `ip` | string | - | device ip |
| `guest_ip` | string | - | device guest ip |
| `wan_ip` | string | - | device wan ip |
| `ssid` | string | - | 2.4G wifi ssid |
| `ssid5g` | string | - | 5G wifi ssid |
| `bssid2g` | string | - | 2.4G wifi bssid |
| `bssid5g` | string | - | 5G wifi bssid |
| `reboot_flag` | bool | - | reboot flag |
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud-batch-manage\",\"send_router_info\"]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"model\":\"x300b\",\"rtty_ssh\":\"true\",\"rtty_web\":\"true\",\"ddns\":\"12345678\",\"sn\":\"123456\",\"mac\":\"AABBCCDDEEFF\",\"fw_type\":\"1\",\"version\":\"VB_3.023\",\"dataupload\":\"true\",\"model\":\"bm\",\"mesh\":\"mesh\",\"services\":[\"abc\","www"],\"firmware_path\":\"http://xxx\",\"data_path\":\"http://xxx\",\"type\":\"1\",\"ip\":\"192.168.8.216\",\"guest_ip\":\"192.168.8.215\",\"wan_ip\":\"192.168.8.3\",\"ssid\":\"GL-wifi\",\"ssid5g\":\"GL-wifi-5G\",\"bssid2g\":\"XXX\",\"bssid5g\":\"YYY\",\"reboot_flag\":\"true\"}}
```

---

## set_2b_config

set 2b config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `lan_ip` | string | no | device lan ip |
| `autoupdate_firmeare_path` | string | no | device autoupdate firmare path |
| `autoupdate_time` | string | no | device autoupdate time |
| `password` | string | no | system password |
| `system_timezone` | string | no | system timezone |
| `hostname` | string | no | hostname |
| `auto_timezone` | string | no | auto sync timezone |
| `server_name` | string | no | server name |
| `server_url` | string | no | server url |
| `language` | string | no | languate |
| `image_url` | string | no | image url |
| `customer_name` | string | no | customer name |
| `help_url` | string | no | help url |
| `ssid_2g` | string | no | 2.4G ssid |
| `channel_2g` | string | no | 2.4G channle |
| `txpower_2g` | string | no | 2.4G txpower |
| `htmode_2g` | string | no | 2.4g htmode |
| `encrytion_2g` | string | no | 2.4G encryption |
| `key_2g` | string | no | 2.4G key |
| `guest_2g_ssid` | string | no | 2.4G guest ssid |
| `guest_2g_encrytion` | string | no | 2.4G guest encryption |
| `guest_2g_key` | string | no | 2.4G guest key |
| `ssid_5g` | string | no | 5G ssid |
| `channel_5g` | string | no | 5G channle |
| `txpower_5g` | string | no | 5G txpower |
| `htmode_5g` | string | no | 5G htmode |
| `encrytion_5g` | string | no | 5G encryption |
| `key_5g` | string | no | 5G key |
| `guest_5g_ssid` | string | no | 5G guest ssid |
| `guest_5g_encrytion` | string | no | 5G guest encryption |
| `guest_5g_key` | string | no | 5G guest key |
| `reboot_flag` | string | no | sysem reboot flag |
| `private_disabled` | string | no | fb private disabled |
| `private_controller` | string | no | fb private controller |
| `guest_controller` | string | no | fb guest controller |
| `mode` | string | no | fb mode |
| `ssl_id` | string | no | fb ssl id |
| `private_wired` | string | no | fb private wired |
| `private_datapath_id` | string | no | fb pravate data path id |
| `private_datapath_desc` | string | no | fb private data path describte |
| `private_ipaddr` | string | no | fb private ipaddr |
| `private_netmask` | string | no | fb private dhcp lease |
| `private_dhcp_lease` | string | no | fb private wlans |
| `private_wlans` | string | no | fb private wlans |
| `guest_wired` | string | no | fb guest wired |
| `guest_datapath_id` | string | no | fb guest datapath id |
| `guest_datapath_desc` | string | no | fb guest datapath desc |
| `guest_ipaddr` | string | no | fb guest ipaddr |
| `guest_netmask` | string | no | fb guest netmask |
| `guest_dhcp_lease` | string | no | fb dhcp lease |
| `guest_wlans` | string | no | gb guest wlans |
| `autoupdate_enable` | bool | no | enable autoupdate |
| `disable_2g` | bool | no | 2.4g disable |
| `hidden_2g` | bool | no | 2.4G hidden |
| `guest_2g_disable` | bool | no | 2.4G guest disable |
| `disable_5g` | bool | no | 5G disable |
| `hidden_5g` | bool | no | 5G hidden |
| `guest_5g_disable` | bool | no | 5G guest disable |
| `guest_disabled` | bool | no | fb guest disable |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud-batch-manage\",\"set_2b_config\",{\"lan_ip:\"192.168.8.1\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---

## set_batch_config

set batch config

Not exercised against hardware here.

| Params | Type | Required | Description |
|---|---|---|---|
| `wan_proto` | string | no | wan proto |
| `wan_username` | string | no | wan proto username |
| `wan_password` | string | no | wan proto password |
| `wan_ipaddr` | string | no | wan ipaddr |
| `wan_gateway` | string | no | wan gateway |
| `wan_netmask` | string | no | wan netmask |
| `wan_dns` | string | no | wan dns |
| `ssid_2g` | string | no | 2.4G ssid |
| `channel_2g` | string | no | 2.4G channle |
| `txpower_2g` | string | no | 2.4G txpower |
| `htmode_2g` | string | no | 2.4g htmode |
| `encrytion_2g` | string | no | 2.4G encryption |
| `key_2g` | string | no | 2.4G key |
| `guest_2g_ssid` | string | no | 2.4G guest ssid |
| `guest_2g_encrytion` | string | no | 2.4G guest encryption |
| `guest_2g_key` | string | no | 2.4G guest key |
| `ssid_5g` | string | no | 5G ssid |
| `channel_5g` | string | no | 5G channle |
| `txpower_5g` | string | no | 5G txpower |
| `htmode_5g` | string | no | 5G htmode |
| `encrytion_5g` | string | no | 5G encryption |
| `key_5g` | string | no | 5G key |
| `guest_5g_ssid` | string | no | 5G guest ssid |
| `guest_5g_encrytion` | string | no | 5G guest encryption |
| `guest_5g_key` | string | no | 5G guest key |
| `probe_data_path1` | string | no | probe data path |
| `probe_data_path2` | string | no | probe data path |
| `probe_max_lines` | string | no | probe max lines |
| `probe_interval` | string | no | probe interval |
| `probe_json` | string | no | probe json |
| `probe_virmac` | string | no | probe virmac |
| `probe_type` | string | no | probe type |
| `probe_tcpurl` | string | no | probe tcp url |
| `probe_filter` | string | no | probe filter |
| `probe_sflag` | string | no | probe sflag |
| `probe_protocol` | string | no | probe protocol |
| `probe_channel` | string | no | probe channel |
| `probe_channel5g` | string | no | probe 5G channel |
| `probe_induce` | string | no | probe induce |
| `lan_ip` | string | no | lan ip |
| `guest_lan_ip` | string | no | guest lan ip |
| `autoupdate_firmware_path` | string | no | autoupdate firmware path |
| `autoupdate_time` | string | no | autoupdate time |
| `password` | string | no | system password |
| `system_timezone` | string | no | system timezone |
| `auto_timezone` | string | no | auto timezone |
| `reboot_flag` | bool | no | reboot_flag |
| `disable_2g` | bool | no | 2.4g disable |
| `hidden_2g` | bool | no | 2.4G hidden |
| `guest_2g_disable` | bool | no | 2.4G guest disable |
| `disable_5g` | bool | no | 5G disable |
| `hidden_5g` | bool | no | 5G hidden |
| `guest_5g_disable` | bool | no | 5G guest disable |
| `autoupdate_enable` | bool | no | atuoupdate enable |

| Results | Type | Required | Description |
|---|---|---|---|
| `err_code` | number | - | Error code, -1:argument err |
| `err_msg` | string | - | Error message |

**Request**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"call\",\"params\":[\"\",\"cloud-batch-manage\",\"set_batch_config\",{\"wan_proto\":\"dhcp\"}]}
```

**Response**

```json
{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}
```

---
