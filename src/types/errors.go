package types

import "errors"

// Sentinels live here rather than in the root package because types is a leaf:
// services and utilities both need them, and the root package re-exports them
// so consumers need not import types just to test for a failure.
var (
	// ErrInvalidName means a reservation name failed validation. Returned rather
	// than escaped: GL.iNet writes the name into dnsmasq's configuration file,
	// and a bad character there breaks DHCP and DNS for the whole router.
	ErrInvalidName = errors.New("gogl: invalid reservation name")

	// ErrOutsideSubnet means an address does not fall inside the router's LAN.
	ErrOutsideSubnet = errors.New("gogl: address outside LAN subnet")

	// ErrInvalidMAC means a MAC address was not parseable.
	ErrInvalidMAC = errors.New("gogl: invalid MAC address")

	// ErrInvalidIP means an address was not a valid IPv4 address.
	ErrInvalidIP = errors.New("gogl: invalid IPv4 address")

	// ErrNotFound means no matching record exists on the device.
	ErrNotFound = errors.New("gogl: not found")

	// ErrConflict means the write would collide with an existing record.
	ErrConflict = errors.New("gogl: conflict")

	// ErrDomainNotSet means the DNS domain has not been configured on the router.
	// Reservation writes refuse to proceed without it: a reservation alone creates
	// no DNS record, so writing one before the domain exists produces addresses
	// with no names and no way to tell that was unintended.
	ErrDomainNotSet = errors.New("gogl: DNS domain is not configured")

	// ErrReservationsExist means a network change was refused because reservations
	// are still present. Renumbering the LAN under live reservations would leave
	// every one of them into the new subnet without saying so; see NetworkService.Set.
	ErrReservationsExist = errors.New("gogl: reservations exist")

	// ErrInvalidInput means a caller-supplied value is unfit to write. Distinct from
	// the field-specific sentinels, for values with no sentinel of their own.
	ErrInvalidInput = errors.New("gogl: invalid input")

	// ErrWirelessSession means a wireless write was refused because the session
	// issuing it arrives over WiFi, and applying it would sever that session with no
	// address to reconnect at.
	ErrWirelessSession = errors.New("gogl: refusing to change wireless over a wireless session")

	// ErrUnwritableContent means host-file content contains a character the
	// firmware's dns.set_host refuses. It reports -32602 Invalid params without
	// saying which character, so gogl checks first and names it.
	ErrUnwritableContent = errors.New("gogl: host file content the firmware will not accept")
)
