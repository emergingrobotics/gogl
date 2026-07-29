// Package services implements the typed operations gogl offers against a
// GL.iNet router. Interfaces are small and single-purpose so each is
// independently mockable.
//
// Three of the four are read-only. Only ReservationService writes: network
// configuration is set in the GL.iNet admin panel, not here. That boundary is
// what keeps a bulk reservation import from ever leaving the router unreachable.
//
// No method takes a site parameter. GL.iNet routers have no equivalent of a
// UniFi site.
package services

import (
	"context"

	"github.com/emergingrobotics/gogl/src/types"
)

// NetworkService reads the router's LAN and DHCP configuration. Read-only.
type NetworkService interface {
	// Get returns the main LAN interface.
	Get(ctx context.Context) (*types.Network, error)

	// List returns every interface the router reports, main LAN and guest alike.
	List(ctx context.Context) ([]types.Network, error)

	// Leases returns the dynamic DHCP leases currently held. These are not
	// reservations: a lease expires, a reservation is a permanent MAC-to-IP
	// binding. Leases are how you discover what is worth reserving.
	Leases(ctx context.Context) ([]types.DHCPLease, error)

	// Set writes an interface's address and DHCP pool.
	//
	// Refused with ErrReservationsExist while any reservation is present:
	// renumbering the LAN underneath live reservations would leave every one of
	// them pointing outside the new subnet, inert and easy to miss. Clear them
	// first.
	//
	// This changes the address the router is being managed at, so the calling
	// session will not survive it.
	Set(ctx context.Context, n *types.Network) error
}

// HostsService manages DNS names through the router's hosts file.
//
// This is the only way gogl creates DNS records. A reservation does not: on this
// firmware its name is a label. gogl owns a delimited block in the file and
// preserves everything outside it.
type HostsService interface {
	// Get returns the parsed host file, managed and unmanaged parts alike.
	Get(ctx context.Context) (*types.HostFile, error)

	// Put writes the host file back.
	Put(ctx context.Context, f *types.HostFile) error

	// Domain returns the configured DNS domain, empty if never set.
	Domain(ctx context.Context) (string, error)

	// SetDomain configures the domain and requalifies existing entries.
	SetDomain(ctx context.Context, domain string) error

	// List returns the managed entries.
	List(ctx context.Context) ([]types.HostEntry, error)

	// Set points name at ip. Requires a configured domain.
	Set(ctx context.Context, name, ip string) error

	// Remove drops the entry for name.
	Remove(ctx context.Context, name string) error

	// Clear removes every managed entry.
	Clear(ctx context.Context) error
}

// ReservationService manages static DHCP bindings.
//
// A binding pins a MAC to an IP. It does not create a DNS record -- see
// types.Reservation, and HostsService for names.
type ReservationService interface {
	List(ctx context.Context) ([]types.Reservation, error)

	// GetByMAC returns the reservation for mac, or ErrNotFound.
	GetByMAC(ctx context.Context, mac string) (*types.Reservation, error)

	// GetByIP returns every reservation holding ip. More than one indicates
	// inconsistent device state rather than normal operation, so the caller
	// decides whether to tolerate it.
	GetByIP(ctx context.Context, ip string) ([]types.Reservation, error)

	// GetByName returns the reservation named name, or ErrNotFound.
	GetByName(ctx context.Context, name string) (*types.Reservation, error)

	// Create writes a new reservation. Returns ErrConflict if the MAC is already
	// reserved, and ErrDomainNotSet if no DNS domain has been configured.
	// Validates before touching the device.
	Create(ctx context.Context, r *types.Reservation) (*types.Reservation, error)

	// Update replaces the reservation identified by r.MAC.
	Update(ctx context.Context, r *types.Reservation) (*types.Reservation, error)

	// Delete removes the reservation for mac.
	Delete(ctx context.Context, mac string) error

	// DeleteAll removes every reservation in one call.
	//
	// Deliberately explicit rather than a mode flag on Delete: this is the one
	// operation that can discard a whole network's addressing, and it should be
	// impossible to reach by passing the wrong argument.
	DeleteAll(ctx context.Context) error
}

// ClientService reads stations known to the router. Read-only.
type ClientService interface {
	List(ctx context.Context) ([]types.Client, error)
}

// SystemService reads device identity. Read-only.
type SystemService interface {
	Info(ctx context.Context) (*types.SystemInfo, error)
}
