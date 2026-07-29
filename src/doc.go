// Package gogl provides programmatic control of GL.iNet travel routers running
// firmware 4.x, targeting the GL-SFT1200 (Opal).
//
// The router exposes a single JSON-RPC 2.0 endpoint at POST /rpc. Authentication
// is challenge/response: the password is never transmitted, only a digest.
//
// Reservations are the only thing this package writes. Network configuration,
// clients, and system information are read-only; set network configuration in
// the GL.iNet admin panel.
//
// A reservation pins a MAC address to an IP address and nothing more. It does not
// create a DNS record: verified against a GL-SFT1200 on firmware 4.3.28, a static
// bind's name is a label and the router's DNS answers from DHCP lease hostnames,
// which clients report themselves.
//
// So this package reproduces a network's addresses, not its names. Addresses are
// the part devices depend on; names follow only if the clients announce them.
package gogl
