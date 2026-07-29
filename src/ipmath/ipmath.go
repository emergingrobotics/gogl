// Package ipmath provides the IPv4 arithmetic shared by gogl's services and
// utilities: numeric ordering, range containment, and subnet derivation.
//
// It is public rather than internal because the utilities are separate main
// packages outside src/, and they need the same numeric ordering the services
// use. Duplicating a sort comparator per utility would be worse than exporting
// it once.
package ipmath

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	ipv4Bits = 32
	// networkAndBroadcast accounts for the two addresses in any subnet wider
	// than a /31 that cannot be assigned to a host.
	networkAndBroadcast = 2
)

// ToUint32 converts an IPv4 address to its numeric value so that addresses sort
// numerically rather than lexically. A non-IPv4 address yields 0.
func ToUint32(ip net.IP) uint32 {
	v4 := ip.To4()
	if v4 == nil {
		return 0
	}
	return binary.BigEndian.Uint32(v4)
}

// InRange reports whether ip lies within [start, stop] inclusive.
func InRange(ip, start, stop net.IP) bool {
	n := ToUint32(ip)
	return n >= ToUint32(start) && n <= ToUint32(stop)
}

// SubnetFrom derives the CIDR network containing ip under mask, which is how
// GL.iNet reports LAN configuration: an address plus a dotted-quad netmask.
func SubnetFrom(ip, mask string) (*net.IPNet, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || parsedIP.To4() == nil {
		return nil, fmt.Errorf("ipmath: %q is not an IPv4 address", ip)
	}
	parsedMask := net.ParseIP(mask)
	if parsedMask == nil || parsedMask.To4() == nil {
		return nil, fmt.Errorf("ipmath: %q is not an IPv4 netmask", mask)
	}

	m := net.IPMask(parsedMask.To4())
	// Size reports 0,0 for a non-contiguous mask, which is not a usable subnet.
	if ones, bits := m.Size(); ones == 0 && bits == 0 {
		return nil, fmt.Errorf("ipmath: %q is not a contiguous netmask", mask)
	}

	return &net.IPNet{IP: parsedIP.Mask(m), Mask: m}, nil
}

// UsableHosts returns the count of assignable host addresses in n, excluding
// the network and broadcast addresses.
func UsableHosts(n *net.IPNet) int {
	ones, bits := n.Mask.Size()
	if bits != ipv4Bits {
		return 0
	}
	hostBits := bits - ones
	if hostBits < 2 {
		return 0
	}
	return (1 << hostBits) - networkAndBroadcast
}
