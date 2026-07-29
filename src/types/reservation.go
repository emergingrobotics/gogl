package types

import (
	"fmt"
	"net"
	"strings"
)

const (
	maxNameLength  = 253
	maxLabelLength = 63
	macOctets      = 6
	hexPerOctet    = 2
)

// Reservation is a static DHCP binding -- what GL.iNet calls a "static bind".
//
// A reservation pins a MAC address to an IP address. It does NOT create a DNS
// record: that was tested on a GL-SFT1200 running 4.3.28 and is false. A bind's
// Name is a label -- it identifies the entry in the admin panel and in an exported
// host file, and nothing more.
//
// DNS records come from HostFile, written through HostsService. The router also
// answers from DHCP *lease* hostnames that clients supply themselves, which is why
// a name sometimes resolves with no help from gogl; that disappears with the lease
// and is not something to rely on.
//
// The three fields are exactly what lan.get_static_bind_list returns. There is no
// enable/disable flag: a binding either exists or it does not.
type Reservation struct {
	// Name labels the entry. It is NOT a DNS name -- see the type comment.
	Name string `json:"name"`

	// MAC is the client identity, lowercase colon-separated. It is the key for
	// update and delete: it is the only thing a client cannot change about
	// itself, and it is what dnsmasq keys the lease on.
	MAC string `json:"mac"`

	// IP is the reserved IPv4 address.
	IP string `json:"ip"`
}

// Validate reports whether r is fit to write, normalizing MAC to lowercase
// colon-separated form in place. The service layer calls this before every
// write so that no consumer can bypass it.
func (r *Reservation) Validate() error {
	if err := ValidateName(r.Name); err != nil {
		return err
	}

	mac, err := NormalizeMAC(r.MAC)
	if err != nil {
		return err
	}
	r.MAC = mac

	ip := net.ParseIP(r.IP)
	if ip == nil || ip.To4() == nil {
		return fmt.Errorf("%w: %q", ErrInvalidIP, r.IP)
	}

	return nil
}

// ValidateName enforces DNS label rules on the name.
//
// Two reasons, neither of which is "the router will serve it as a DNS record" --
// it will not. First, GL.iNet writes the name into its dnsmasq configuration, and
// a known firmware defect lets a quote there corrupt the file and break DHCP for
// the whole router. Second, the ISC DHCP host-declaration format this project
// exchanges is keyed by hostname, and a name that is not a legal DNS label cannot
// round-trip through it.
//
// It rejects rather than escapes, and names the offending character.
//
// This is deliberately stricter than gofips, which permits underscores. An
// underscore is legal in a UniFi record but is not a legal DNS label character.
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: empty", ErrInvalidName)
	}
	if len(name) > maxNameLength {
		return fmt.Errorf("%w: %d characters exceeds maximum of %d", ErrInvalidName, len(name), maxNameLength)
	}
	if strings.HasPrefix(name, "-") || strings.HasPrefix(name, ".") {
		return fmt.Errorf("%w: %q must not begin with %q", ErrInvalidName, name, name[:1])
	}
	if strings.HasSuffix(name, "-") || strings.HasSuffix(name, ".") {
		return fmt.Errorf("%w: %q must not end with %q", ErrInvalidName, name, name[len(name)-1:])
	}

	for _, label := range strings.Split(name, ".") {
		if label == "" {
			return fmt.Errorf("%w: %q contains an empty label", ErrInvalidName, name)
		}
		if len(label) > maxLabelLength {
			return fmt.Errorf("%w: label %q exceeds %d characters", ErrInvalidName, label, maxLabelLength)
		}
		for _, c := range label {
			if !isNameRune(c) {
				return fmt.Errorf("%w: %q contains %q, which is not permitted (allowed: letters, digits, hyphen, and dot as a separator)",
					ErrInvalidName, name, string(c))
			}
		}
	}
	return nil
}

func isNameRune(c rune) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c >= '0' && c <= '9':
		return true
	case c == '-':
		return true
	default:
		return false
	}
}

// NormalizeMAC parses any form net.ParseMAC accepts and returns the lowercase
// colon-separated form used everywhere in this module.
func NormalizeMAC(mac string) (string, error) {
	trimmed := strings.TrimSpace(mac)
	if trimmed == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidMAC)
	}

	// net.ParseMAC does not accept the unseparated form, which appears in IEEE
	// OUI files and occasionally in hand-written host files.
	if len(trimmed) == macOctets*hexPerOctet && !strings.ContainsAny(trimmed, ":-.") {
		var b strings.Builder
		for i := 0; i < len(trimmed); i += hexPerOctet {
			if i > 0 {
				b.WriteByte(':')
			}
			b.WriteString(trimmed[i : i+hexPerOctet])
		}
		trimmed = b.String()
	}

	parsed, err := net.ParseMAC(trimmed)
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidMAC, mac)
	}
	if len(parsed) != macOctets {
		return "", fmt.Errorf("%w: %q is not a 6-octet address", ErrInvalidMAC, mac)
	}
	return strings.ToLower(parsed.String()), nil
}
