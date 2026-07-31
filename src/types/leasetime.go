package types

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// LeaseInfinite represents dnsmasq's "infinite" lease. It is a sentinel rather
// than a very large duration so that arithmetic on it cannot silently overflow.
const LeaseInfinite = LeaseTime(math.MaxInt64)

const (
	hoursPerDay  = 24
	daysPerWeek  = 7
	infiniteText = "infinite"
)

// LeaseTime is a DHCP lease duration. It unmarshals from a dnsmasq duration
// string ("12h", "1d", "2w", "infinite") or a bare count of seconds, and
// marshals to the dnsmasq duration form.
//
// Both forms are accepted because DHCP tooling disagrees: some expresses a lease
// as an integer count of seconds, dnsmasq as a duration string. Only unmarshalling
// is exercised against a device — the firmware exposes no verified endpoint that
// writes the lease time, so gogl reads it and cannot set it.
type LeaseTime time.Duration

// UnmarshalJSON accepts a JSON number of seconds or a duration string. Both
// forms occur in practice, so tolerating each avoids a class of bug that would
// only appear against one firmware.
func (l *LeaseTime) UnmarshalJSON(b []byte) error {
	var asNumber int64
	if err := json.Unmarshal(b, &asNumber); err == nil {
		*l = LeaseTime(time.Duration(asNumber) * time.Second)
		return nil
	}

	var asString string
	if err := json.Unmarshal(b, &asString); err != nil {
		return fmt.Errorf("lease time: not a number or string: %s", b)
	}

	parsed, err := parseLeaseString(asString)
	if err != nil {
		return err
	}
	*l = parsed
	return nil
}

// parseLeaseString handles dnsmasq's duration vocabulary. time.ParseDuration
// covers s, m and h but rejects d and w, which dnsmasq uses, so those are
// converted before delegating.
func parseLeaseString(s string) (LeaseTime, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("lease time: empty")
	}
	if s == infiniteText {
		return LeaseInfinite, nil
	}

	if seconds, err := strconv.ParseInt(s, 10, 64); err == nil {
		return LeaseTime(time.Duration(seconds) * time.Second), nil
	}

	value, unit := s[:len(s)-1], s[len(s)-1]
	switch unit {
	case 'd', 'w':
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("lease time: bad quantity in %q", s)
		}
		hours := n * hoursPerDay
		if unit == 'w' {
			hours *= daysPerWeek
		}
		return LeaseTime(time.Duration(hours) * time.Hour), nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("lease time: unrecognized duration %q", s)
	}
	return LeaseTime(d), nil
}

// MarshalJSON emits the duration string form. gogl never writes network
// configuration to the router, so this serves only the utilities' own JSON
// output, where the friendlier form is the useful one.
func (l LeaseTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

// String renders the duration in dnsmasq's compact form. time.Duration's own
// String would give "12h0m0s", which is correct but noisy in a report.
func (l LeaseTime) String() string {
	if l == LeaseInfinite {
		return infiniteText
	}
	d := time.Duration(l)
	switch {
	case d == 0:
		return "0s"
	case d%time.Hour == 0:
		return fmt.Sprintf("%dh", d/time.Hour)
	case d%time.Minute == 0:
		return fmt.Sprintf("%dm", d/time.Minute)
	case d < time.Minute && d%time.Second == 0:
		return fmt.Sprintf("%ds", d/time.Second)
	default:
		return d.String()
	}
}

// Seconds returns the lease in whole seconds, or -1 for an infinite lease.
func (l LeaseTime) Seconds() int64 {
	if l == LeaseInfinite {
		return -1
	}
	return int64(time.Duration(l).Seconds())
}
