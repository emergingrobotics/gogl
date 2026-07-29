package ipmath

import (
	"net"
	"testing"
)

func TestToUint32(t *testing.T) {
	tests := []struct {
		ip   string
		want uint32
	}{
		{"0.0.0.0", 0},
		{"0.0.0.1", 1},
		{"0.0.1.0", 256},
		{"192.168.8.1", 3232237569},
		{"255.255.255.255", 4294967295},
	}
	for _, tt := range tests {
		if got := ToUint32(net.ParseIP(tt.ip)); got != tt.want {
			t.Errorf("ToUint32(%s) = %d, want %d", tt.ip, got, tt.want)
		}
	}
}

func TestToUint32NonIPv4(t *testing.T) {
	if got := ToUint32(net.ParseIP("fe80::1")); got != 0 {
		t.Errorf("ToUint32(fe80::1) = %d, want 0", got)
	}
	if got := ToUint32(nil); got != 0 {
		t.Errorf("ToUint32(nil) = %d, want 0", got)
	}
}

// Sorting by uint32 must order numerically, not lexically: the whole point is
// that 192.168.8.9 precedes 192.168.8.10.
func TestToUint32OrdersNumerically(t *testing.T) {
	nine := ToUint32(net.ParseIP("192.168.8.9"))
	ten := ToUint32(net.ParseIP("192.168.8.10"))
	if nine >= ten {
		t.Errorf("192.168.8.9 (%d) should sort before 192.168.8.10 (%d)", nine, ten)
	}
}

func TestInRange(t *testing.T) {
	start, stop := net.ParseIP("192.168.8.100"), net.ParseIP("192.168.8.249")
	tests := []struct {
		ip   string
		want bool
	}{
		{"192.168.8.99", false},
		{"192.168.8.100", true},
		{"192.168.8.175", true},
		{"192.168.8.249", true},
		{"192.168.8.250", false},
	}
	for _, tt := range tests {
		if got := InRange(net.ParseIP(tt.ip), start, stop); got != tt.want {
			t.Errorf("InRange(%s) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestSubnetFrom(t *testing.T) {
	tests := []struct {
		ip   string
		mask string
		want string
	}{
		{"192.168.8.1", "255.255.255.0", "192.168.8.0/24"},
		{"192.168.8.130", "255.255.255.0", "192.168.8.0/24"},
		{"10.1.2.3", "255.255.0.0", "10.1.0.0/16"},
	}
	for _, tt := range tests {
		n, err := SubnetFrom(tt.ip, tt.mask)
		if err != nil {
			t.Errorf("SubnetFrom(%s, %s) error: %v", tt.ip, tt.mask, err)
			continue
		}
		if got := n.String(); got != tt.want {
			t.Errorf("SubnetFrom(%s, %s) = %s, want %s", tt.ip, tt.mask, got, tt.want)
		}
	}
}

func TestSubnetFromRejects(t *testing.T) {
	for _, tt := range []struct{ ip, mask string }{
		{"", "255.255.255.0"},
		{"192.168.8.1", ""},
		{"192.168.8.1", "not-a-mask"},
		{"nope", "255.255.255.0"},
		{"fe80::1", "255.255.255.0"},
		{"192.168.8.1", "255.0.255.0"},
	} {
		if _, err := SubnetFrom(tt.ip, tt.mask); err == nil {
			t.Errorf("SubnetFrom(%q, %q) succeeded, want error", tt.ip, tt.mask)
		}
	}
}

func TestUsableHosts(t *testing.T) {
	tests := []struct {
		cidr string
		want int
	}{
		{"192.168.8.0/24", 254},
		{"192.168.0.0/16", 65534},
		{"192.168.8.0/30", 2},
		{"192.168.8.0/31", 0},
		{"192.168.8.1/32", 0},
	}
	for _, tt := range tests {
		_, n, err := net.ParseCIDR(tt.cidr)
		if err != nil {
			t.Fatalf("ParseCIDR(%s): %v", tt.cidr, err)
		}
		if got := UsableHosts(n); got != tt.want {
			t.Errorf("UsableHosts(%s) = %d, want %d", tt.cidr, got, tt.want)
		}
	}
}
