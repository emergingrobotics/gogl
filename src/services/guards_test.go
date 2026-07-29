package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/emergingrobotics/gogl/src/mock"
	"github.com/emergingrobotics/gogl/src/services"
	"github.com/emergingrobotics/gogl/src/types"
)

// The factory host file, so tests start from the state a real router is in.
const factoryHosts = "127.0.0.1 localhost\n\n::1     localhost ip6-localhost ip6-loopback\n" +
	"ff02::1 ip6-allnodes\nff02::2 ip6-allrouters\n"

const oneInterface = `{"interfaces":[{"interface":"lan","ip":"192.168.8.1","netmask":"255.255.255.0",` +
	`"enable":1,"start":"192.168.8.100","end":"192.168.8.249","leasetime":"12h"}]}`

// guardServer starts a mock with a network and an optionally-configured domain.
func guardServer(t *testing.T, domain string, seed []types.Reservation) *mock.Server {
	t.Helper()
	s := mock.NewServer(t, mock.Options{Password: "secret"})
	s.LoadFixture(mock.NetworkGroup, mock.MethodGetConfigList, json.RawMessage(oneInterface))
	s.SetReservations(seed)

	s.SetHostFile(mock.HostFileWith(domain))
	return s
}

// ---------------------------------------------------------------------------
// Guard 1: a reservation is not written until a DNS domain exists.
// ---------------------------------------------------------------------------

func TestReservationCreateRefusedWithoutDomain(t *testing.T) {
	s := guardServer(t, "", nil)
	svc := services.NewReservationService(newTransport(t, s))

	_, err := svc.Create(context.Background(), &types.Reservation{
		Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.13",
	})
	if !errors.Is(err, types.ErrDomainNotSet) {
		t.Fatalf("error = %v, want ErrDomainNotSet", err)
	}
	// Refusing must not write.
	if got := s.Reservations(); len(got) != 0 {
		t.Errorf("device holds %d reservations after a refused create, want 0", len(got))
	}
}

func TestReservationUpdateRefusedWithoutDomain(t *testing.T) {
	s := guardServer(t, "", []types.Reservation{
		{Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.13"},
	})
	svc := services.NewReservationService(newTransport(t, s))

	_, err := svc.Update(context.Background(), &types.Reservation{
		Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.20",
	})
	if !errors.Is(err, types.ErrDomainNotSet) {
		t.Errorf("error = %v, want ErrDomainNotSet", err)
	}
}

func TestReservationCreateAllowedWithDomain(t *testing.T) {
	s := guardServer(t, "lab.example", nil)
	svc := services.NewReservationService(newTransport(t, s))

	if _, err := svc.Create(context.Background(), &types.Reservation{
		Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.13",
	}); err != nil {
		t.Fatalf("Create with a domain set: %v", err)
	}
	if got := s.Reservations(); len(got) != 1 {
		t.Errorf("device holds %d reservations, want 1", len(got))
	}
}

// Reads and deletes are not gated: only writes that create addressing are.
func TestReservationReadsAndDeletesWorkWithoutDomain(t *testing.T) {
	s := guardServer(t, "", []types.Reservation{
		{Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.13"},
	})
	svc := services.NewReservationService(newTransport(t, s))
	ctx := context.Background()

	if _, err := svc.List(ctx); err != nil {
		t.Errorf("List without a domain: %v", err)
	}
	if err := svc.Delete(ctx, "aa:bb:cc:dd:ee:01"); err != nil {
		t.Errorf("Delete without a domain: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Guard 2: the network is not renumbered while reservations exist.
// ---------------------------------------------------------------------------

func newNetwork() *types.Network {
	return &types.Network{
		Interface: types.InterfaceLAN,
		LANIP:     "192.168.2.1",
		Netmask:   "255.255.255.0",
		DHCPStart: "192.168.2.100",
		DHCPStop:  "192.168.2.149",
	}
}

func TestNetworkSetRefusedWhileReservationsExist(t *testing.T) {
	s := guardServer(t, "lab.example", []types.Reservation{
		{Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.13"},
	})
	svc := services.NewNetworkService(newTransport(t, s))

	err := svc.Set(context.Background(), newNetwork())
	if !errors.Is(err, types.ErrReservationsExist) {
		t.Fatalf("error = %v, want ErrReservationsExist", err)
	}
	// The count belongs in the message: "clear them" is useless without "how many".
	if !strings.Contains(err.Error(), "1 reservation") {
		t.Errorf("error does not say how many are in the way: %v", err)
	}
	// And nothing may have changed.
	if got := s.Network(); len(got) != 1 || got[0].LANIP != "192.168.8.1" {
		t.Errorf("a refused Set changed the network: %+v", got)
	}
}

func TestNetworkSetAllowedOnceCleared(t *testing.T) {
	s := guardServer(t, "lab.example", []types.Reservation{
		{Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.13"},
	})
	tr := newTransport(t, s)
	reservations := services.NewReservationService(tr)
	network := services.NewNetworkService(tr)
	ctx := context.Background()

	if err := reservations.DeleteAll(ctx); err != nil {
		t.Fatalf("DeleteAll: %v", err)
	}
	if err := network.Set(ctx, newNetwork()); err != nil {
		t.Fatalf("Set after clearing: %v", err)
	}

	got := s.Network()
	if len(got) != 1 || got[0].LANIP != "192.168.2.1" {
		t.Errorf("network = %+v, want the new address", got)
	}
	if got[0].DHCPStart != "192.168.2.100" || got[0].DHCPStop != "192.168.2.149" {
		t.Errorf("pool not written: %+v", got[0])
	}
}

// A pool outside the new subnet yields a DHCP server that hands out nothing, with
// no error from the firmware to explain it. Catch it before the write.
func TestNetworkSetRejectsPoolOutsideSubnet(t *testing.T) {
	s := guardServer(t, "lab.example", nil)
	svc := services.NewNetworkService(newTransport(t, s))

	n := newNetwork()
	n.DHCPStart, n.DHCPStop = "10.0.0.100", "10.0.0.149"

	err := svc.Set(context.Background(), n)
	if err == nil {
		t.Fatal("Set accepted a pool outside the subnet")
	}
	if !strings.Contains(err.Error(), "outside") {
		t.Errorf("error does not explain the problem: %v", err)
	}
}

func TestNetworkSetRejectsInvertedPool(t *testing.T) {
	s := guardServer(t, "lab.example", nil)
	svc := services.NewNetworkService(newTransport(t, s))

	n := newNetwork()
	n.DHCPStart, n.DHCPStop = n.DHCPStop, n.DHCPStart

	if err := svc.Set(context.Background(), n); err == nil {
		t.Error("Set accepted a pool running backwards")
	}
}

func TestNetworkSetRejectsUnusableSubnetAndMissingInterface(t *testing.T) {
	s := guardServer(t, "lab.example", nil)
	svc := services.NewNetworkService(newTransport(t, s))
	ctx := context.Background()

	bad := newNetwork()
	bad.LANIP = "nonsense"
	if err := svc.Set(ctx, bad); err == nil {
		t.Error("Set accepted an unusable address")
	}

	noIface := newNetwork()
	noIface.Interface = ""
	if err := svc.Set(ctx, noIface); err == nil {
		t.Error("Set accepted an empty interface name")
	}
}

// ---------------------------------------------------------------------------
// DeleteAll
// ---------------------------------------------------------------------------

func TestReservationDeleteAll(t *testing.T) {
	s := guardServer(t, "lab.example", []types.Reservation{
		{Name: "a", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.10"},
		{Name: "b", MAC: "aa:bb:cc:dd:ee:02", IP: "192.168.8.11"},
		{Name: "c", MAC: "aa:bb:cc:dd:ee:03", IP: "192.168.8.12"},
	})
	svc := services.NewReservationService(newTransport(t, s))

	if err := svc.DeleteAll(context.Background()); err != nil {
		t.Fatalf("DeleteAll: %v", err)
	}
	if got := s.Reservations(); len(got) != 0 {
		t.Errorf("device holds %d reservations after DeleteAll", len(got))
	}
}

// "Make sure there are none" is a reasonable request, so an empty table is a no-op
// rather than an error -- unlike Delete, where a missing MAC is a typo.
func TestReservationDeleteAllOnEmptyDevice(t *testing.T) {
	s := guardServer(t, "lab.example", nil)
	svc := services.NewReservationService(newTransport(t, s))

	if err := svc.DeleteAll(context.Background()); err != nil {
		t.Errorf("DeleteAll on an empty device: %v", err)
	}
}
