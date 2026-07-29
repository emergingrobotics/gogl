package types

import (
	"encoding/json"
	"errors"
	"net"
	"testing"
	"time"
)

func testNetwork() *Network {
	return &Network{
		LANIP:       "192.168.8.1",
		Netmask:     "255.255.255.0",
		DHCPEnabled: true,
		DHCPStart:   "192.168.8.100",
		DHCPStop:    "192.168.8.249",
		DHCPLease:   LeaseTime(12 * time.Hour),
		Interface:   InterfaceLAN,
		DNS:         []string{"192.168.8.1"},
	}
}

func TestNetworkSubnet(t *testing.T) {
	n, err := testNetwork().Subnet()
	if err != nil {
		t.Fatalf("Subnet() error: %v", err)
	}
	if got := n.String(); got != "192.168.8.0/24" {
		t.Errorf("Subnet() = %s, want 192.168.8.0/24", got)
	}
}

func TestNetworkSubnetError(t *testing.T) {
	n := &Network{LANIP: "nope", Netmask: "255.255.255.0"}
	if _, err := n.Subnet(); err == nil {
		t.Error("Subnet() succeeded on a bad LAN address")
	}
	if got := n.UsableHosts(); got != 0 {
		t.Errorf("UsableHosts() = %d on a bad subnet, want 0", got)
	}
}

func TestNetworkContains(t *testing.T) {
	n := testNetwork()
	tests := []struct {
		ip   string
		want bool
	}{
		{"192.168.8.10", true},
		{"192.168.8.1", true},
		{"192.168.8.254", true},
		{"192.168.4.10", false},
		{"10.0.0.1", false},
	}
	for _, tt := range tests {
		got, err := n.Contains(net.ParseIP(tt.ip))
		if err != nil {
			t.Fatalf("Contains(%s) error: %v", tt.ip, err)
		}
		if got != tt.want {
			t.Errorf("Contains(%s) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestNetworkInDHCPPool(t *testing.T) {
	n := testNetwork()
	tests := []struct {
		ip   string
		want bool
	}{
		{"192.168.8.10", false},
		{"192.168.8.99", false},
		{"192.168.8.100", true},
		{"192.168.8.200", true},
		{"192.168.8.249", true},
		{"192.168.8.250", false},
	}
	for _, tt := range tests {
		got, err := n.InDHCPPool(net.ParseIP(tt.ip))
		if err != nil {
			t.Fatalf("InDHCPPool(%s) error: %v", tt.ip, err)
		}
		if got != tt.want {
			t.Errorf("InDHCPPool(%s) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

// A disabled DHCP server has no pool, so nothing is inside it.
func TestNetworkInDHCPPoolWhenDisabled(t *testing.T) {
	n := testNetwork()
	n.DHCPEnabled = false
	got, err := n.InDHCPPool(net.ParseIP("192.168.8.150"))
	if err != nil {
		t.Fatalf("InDHCPPool error: %v", err)
	}
	if got {
		t.Error("InDHCPPool = true with DHCP disabled, want false")
	}
}

func TestNetworkInDHCPPoolWithBadBoundaries(t *testing.T) {
	n := testNetwork()
	n.DHCPStart = "garbage"
	got, err := n.InDHCPPool(net.ParseIP("192.168.8.150"))
	if err != nil {
		t.Fatalf("InDHCPPool error: %v", err)
	}
	if got {
		t.Error("InDHCPPool = true with an unparseable pool start, want false")
	}
}

func TestNetworkPoolSize(t *testing.T) {
	if got := testNetwork().PoolSize(); got != 150 {
		t.Errorf("PoolSize() = %d, want 150", got)
	}

	disabled := testNetwork()
	disabled.DHCPEnabled = false
	if got := disabled.PoolSize(); got != 0 {
		t.Errorf("PoolSize() with DHCP disabled = %d, want 0", got)
	}

	inverted := testNetwork()
	inverted.DHCPStart, inverted.DHCPStop = inverted.DHCPStop, inverted.DHCPStart
	if got := inverted.PoolSize(); got != 0 {
		t.Errorf("PoolSize() with stop before start = %d, want 0", got)
	}

	bad := testNetwork()
	bad.DHCPStop = "nope"
	if got := bad.PoolSize(); got != 0 {
		t.Errorf("PoolSize() with a bad boundary = %d, want 0", got)
	}
}

func TestNetworkUsableHosts(t *testing.T) {
	if got := testNetwork().UsableHosts(); got != 254 {
		t.Errorf("UsableHosts() = %d, want 254", got)
	}
}

func TestNetworkUnmarshalsLeaseTime(t *testing.T) {
	const payload = `{"interface":"lan","ip":"192.168.8.1","netmask":"255.255.255.0","leasetime":"12h","enable":1}`
	var n Network
	if err := json.Unmarshal([]byte(payload), &n); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if n.DHCPLease != LeaseTime(12*time.Hour) {
		t.Errorf("DHCPLease = %v, want 12h", time.Duration(n.DHCPLease))
	}
}

// ValidateForWrite is the only thing standing between an operator and a DHCP
// server that accepts its configuration and then serves nothing.
func TestValidateForWrite(t *testing.T) {
	valid := func() *Network {
		return &Network{
			Interface: InterfaceLAN,
			LANIP:     "192.168.2.1",
			Netmask:   "255.255.255.0",
			DHCPStart: "192.168.2.100",
			DHCPStop:  "192.168.2.149",
		}
	}

	if err := valid().ValidateForWrite(); err != nil {
		t.Fatalf("a valid network was rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*Network)
		want   error
	}{
		{"no interface", func(n *Network) { n.Interface = "" }, nil},
		{"unparseable address", func(n *Network) { n.LANIP = "nonsense" }, nil},
		{"unparseable mask", func(n *Network) { n.Netmask = "255.255.0" }, nil},
		{"pool start outside subnet", func(n *Network) { n.DHCPStart = "10.0.0.5" }, ErrOutsideSubnet},
		{"pool end outside subnet", func(n *Network) { n.DHCPStop = "10.0.0.5" }, ErrOutsideSubnet},
		{"unparseable pool bound", func(n *Network) { n.DHCPStart = "x" }, ErrInvalidIP},
		{"pool runs backwards", func(n *Network) {
			n.DHCPStart, n.DHCPStop = n.DHCPStop, n.DHCPStart
		}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := valid()
			tt.mutate(n)

			err := n.ValidateForWrite()
			if err == nil {
				t.Fatal("accepted an unusable network")
			}
			if tt.want != nil && !errors.Is(err, tt.want) {
				t.Errorf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

// A /32 has no usable host addresses, so a pool inside it is not describable.
// It is a plausible typo for a mask, which is why it is worth pinning.
func TestValidateForWriteRejectsSingleHostSubnet(t *testing.T) {
	n := &Network{
		Interface: InterfaceLAN,
		LANIP:     "192.168.2.1",
		Netmask:   "255.255.255.255",
		DHCPStart: "192.168.2.100",
		DHCPStop:  "192.168.2.149",
	}
	if err := n.ValidateForWrite(); err == nil {
		t.Error("accepted a pool inside a /32")
	}
}

func TestIsGuest(t *testing.T) {
	if (&Network{Interface: InterfaceGuest}).IsGuest() != true {
		t.Error("the guest interface does not report as guest")
	}
	if (&Network{Interface: InterfaceLAN}).IsGuest() != false {
		t.Error("the LAN interface reports as guest")
	}
}
