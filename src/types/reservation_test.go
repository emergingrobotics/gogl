package types

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateNameAccepts(t *testing.T) {
	valid := []string{
		"nas",
		"my-server",
		"host1",
		"a",
		"aa-bb-cc-dd-ee-ff",
		"printer.lab",
		"MixedCase",
		strings.Repeat("a", 63),
	}
	for _, name := range valid {
		if err := ValidateName(name); err != nil {
			t.Errorf("ValidateName(%q) = %v, want nil", name, err)
		}
	}
}

func TestValidateNameRejects(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"underscore is legal on UniFi but not in DNS", "my_server"},
		{"double quote can corrupt dnsmasq config", `my"server`},
		{"single quote", "my'server"},
		{"space", "my server"},
		{"semicolon", "my;server"},
		{"newline", "my\nserver"},
		{"brace", "my{server"},
		{"leading hyphen", "-nas"},
		{"trailing hyphen", "nas-"},
		{"leading dot", ".nas"},
		{"trailing dot", "nas."},
		{"empty label", "a..b"},
		{"label too long", strings.Repeat("a", 64)},
		{"total too long", strings.Repeat("a.", 127) + "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if err == nil {
				t.Fatalf("ValidateName(%q) = nil, want error", tt.input)
			}
			if !errors.Is(err, ErrInvalidName) {
				t.Errorf("ValidateName(%q) error = %v, want ErrInvalidName", tt.input, err)
			}
		})
	}
}

// The error must name the offending character so the operator can find it in a
// host file without guessing.
func TestValidateNameErrorNamesCharacter(t *testing.T) {
	err := ValidateName("my_server")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "_") {
		t.Errorf("error %q does not name the offending character", err)
	}
}

func TestNormalizeMAC(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"AA:BB:CC:DD:EE:FF", "aa:bb:cc:dd:ee:ff"},
		{"aa:bb:cc:dd:ee:ff", "aa:bb:cc:dd:ee:ff"},
		{"aa-bb-cc-dd-ee-ff", "aa:bb:cc:dd:ee:ff"},
		{"AABBCCDDEEFF", "aa:bb:cc:dd:ee:ff"},
		{"  aa:bb:cc:dd:ee:ff  ", "aa:bb:cc:dd:ee:ff"},
	}
	for _, tt := range tests {
		got, err := NormalizeMAC(tt.input)
		if err != nil {
			t.Errorf("NormalizeMAC(%q) error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("NormalizeMAC(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNormalizeMACRejects(t *testing.T) {
	for _, input := range []string{
		"",
		"   ",
		"aa:bb:cc",
		"zz:bb:cc:dd:ee:ff",
		"aa:bb:cc:dd:ee:ff:00",
		"not-a-mac",
	} {
		if _, err := NormalizeMAC(input); !errors.Is(err, ErrInvalidMAC) {
			t.Errorf("NormalizeMAC(%q) error = %v, want ErrInvalidMAC", input, err)
		}
	}
}

func TestReservationValidate(t *testing.T) {
	good := &Reservation{Name: "nas", MAC: "AA:BB:CC:DD:EE:01", IP: "192.168.8.10"}
	if err := good.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	// Validate normalizes in place so the service layer always writes lowercase.
	if good.MAC != "aa:bb:cc:dd:ee:01" {
		t.Errorf("Validate did not normalize MAC: got %q", good.MAC)
	}
}

func TestReservationValidateRejects(t *testing.T) {
	tests := []struct {
		name string
		res  Reservation
		want error
	}{
		{"bad name", Reservation{Name: "my_nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.10"}, ErrInvalidName},
		{"bad mac", Reservation{Name: "nas", MAC: "nope", IP: "192.168.8.10"}, ErrInvalidMAC},
		{"bad ip", Reservation{Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "999.1.1.1"}, ErrInvalidIP},
		{"ipv6 rejected", Reservation{Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "fe80::1"}, ErrInvalidIP},
		{"empty ip", Reservation{Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: ""}, ErrInvalidIP},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.res
			if err := res.Validate(); !errors.Is(err, tt.want) {
				t.Errorf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}
