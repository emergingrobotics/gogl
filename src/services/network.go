package services

import (
	"context"
	"fmt"

	"github.com/emergingrobotics/gogl/src/transport"
	"github.com/emergingrobotics/gogl/src/types"
)

// CONFIRMED against a GL-SFT1200 on firmware 4.3.28, 2026-07-28.
//
// The group is "lan" but the method is get_config_list, not get_config -- the
// latter does not exist and returns -32601, which is what blocked this service
// until the endpoint was captured.
const (
	networkGroup        = "lan"
	networkGetConfig    = "get_config_list"
	networkSetConfig    = "set_config"
	networkDHCPGroup    = "network"
	networkGetDHCPLease = "get_dhcp_leases"
)

// interfaceList is the wire shape of lan.get_config_list: both the LAN and the
// guest network come back in one call.
type interfaceList struct {
	Interfaces []types.Network `json:"interfaces"`
}

// leaseList is the wire shape of network.get_dhcp_leases.
//
// The published documentation calls this field "entries"; the device actually
// sends "leases". The device wins.
type leaseList struct {
	Leases []types.DHCPLease `json:"leases"`
}

type networkService struct {
	transport transport.Transport

	// reservations is consulted before a write, to enforce the ordering rule that
	// the LAN cannot be renumbered underneath live reservations.
	reservations ReservationService
}

// NewNetworkService returns the LAN and DHCP configuration service.
func NewNetworkService(t transport.Transport) NetworkService {
	return &networkService{transport: t, reservations: NewReservationService(t)}
}

// Get reads the main LAN's address and DHCP configuration.
func (s *networkService) Get(ctx context.Context) (*types.Network, error) {
	all, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range all {
		if all[i].Interface == types.InterfaceLAN {
			return &all[i], nil
		}
	}
	return nil, fmt.Errorf("%w: router reported no %q interface", types.ErrNotFound, types.InterfaceLAN)
}

// List returns every interface the router reports, main LAN and guest alike.
// Guest is useful for address planning even though gogl never writes it.
func (s *networkService) List(ctx context.Context) ([]types.Network, error) {
	var list interfaceList
	if err := s.transport.Call(ctx, networkGroup, networkGetConfig, nil, &list); err != nil {
		return nil, fmt.Errorf("gogl: read network config: %w", err)
	}
	return list.Interfaces, nil
}

// Leases returns the dynamic DHCP leases the router is currently holding.
//
// These are not reservations: a lease expires, a reservation is a permanent
// MAC-to-IP binding. Leases are how you discover what is worth reserving, and
// their hostnames are what the router's DNS actually answers.
func (s *networkService) Leases(ctx context.Context) ([]types.DHCPLease, error) {
	var list leaseList
	if err := s.transport.Call(ctx, networkDHCPGroup, networkGetDHCPLease, nil, &list); err != nil {
		return nil, fmt.Errorf("gogl: read DHCP leases: %w", err)
	}
	return list.Leases, nil
}

// Set writes an interface's address and DHCP pool via lan.set_config.
//
// Refused while any reservation exists. Renumbering the LAN underneath live
// reservations leaves every one of them outside the new subnet -- present in the
// table, silently inert, and easy to miss for a long time. Clear them first, apply
// the new network, then import the addresses you want.
//
// This changes the address the router is managed at, so the session making the call
// will not survive it. That is inherent, not a defect: expect to reconnect.
func (s *networkService) Set(ctx context.Context, n *types.Network) error {
	if err := n.ValidateForWrite(); err != nil {
		return fmt.Errorf("gogl: %w", err)
	}

	existing, err := s.reservations.List(ctx)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return fmt.Errorf("%w: %d reservation(s) present; clear them before changing the network",
			types.ErrReservationsExist, len(existing))
	}

	args := map[string]any{
		"interface": n.Interface,
		"ip":        n.LANIP,
		"netmask":   n.Netmask,
		"start":     n.DHCPStart,
		"end":       n.DHCPStop,
	}
	if err := s.transport.Call(ctx, networkGroup, networkSetConfig, args, nil); err != nil {
		return fmt.Errorf("gogl: write network config: %w", err)
	}
	return nil
}
