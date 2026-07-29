package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"testing"

	gogl "github.com/emergingrobotics/gogl/src"
	"github.com/emergingrobotics/gogl/src/mock"
	"github.com/emergingrobotics/gogl/src/types"
)

// These drive runSetNetwork through a real client against the mock router,
// rather than calling its helpers directly.
//
// That distinction is not academic here. An earlier version of goglps had a
// --dry-run that performed a live write, because Go's flag package stops parsing
// at the first operand. Every unit test passed: they constructed the flag struct
// directly and so could not observe the parse. Only exercising the real path
// caught it.

const writeLANFixture = `{
  "interfaces": [
    {
      "interface": "lan",
      "ip": "192.168.8.1",
      "netmask": "255.255.255.0",
      "enable": 1,
      "start": "192.168.8.100",
      "end": "192.168.8.249",
      "leasetime": "12h"
    }
  ]
}`

func mockClient(t *testing.T, seed []types.Reservation) (*mock.Server, *gogl.Client) {
	t.Helper()

	s := mock.NewServer(t, mock.Options{Password: "secret"})
	s.LoadFixture(mock.NetworkGroup, mock.MethodGetConfigList, json.RawMessage(writeLANFixture))
	s.SetReservations(seed)

	u, err := url.Parse(s.URL())
	if err != nil {
		t.Fatalf("parse mock URL: %v", err)
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		t.Fatalf("parse port: %v", err)
	}

	c, err := gogl.New(gogl.Config{
		Host: u.Hostname(), Port: port, Password: "secret",
		KeepaliveInterval: -1,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})
	return s, c
}

func targetNetwork() *types.Network {
	return &types.Network{
		Interface: types.InterfaceLAN,
		LANIP:     "192.168.2.1",
		Netmask:   "255.255.255.0",
		DHCPStart: "192.168.2.100",
		DHCPStop:  "192.168.2.149",
	}
}

func TestRunSetNetworkWrites(t *testing.T) {
	s, c := mockClient(t, nil)

	if err := runSetNetwork(context.Background(), c, targetNetwork(), false); err != nil {
		t.Fatalf("runSetNetwork: %v", err)
	}

	got := s.Network()
	if len(got) != 1 {
		t.Fatalf("device reports %d interfaces, want 1", len(got))
	}
	if got[0].LANIP != "192.168.2.1" || got[0].Netmask != "255.255.255.0" {
		t.Errorf("address = %s/%s, want 192.168.2.1/255.255.255.0", got[0].LANIP, got[0].Netmask)
	}
	if got[0].DHCPStart != "192.168.2.100" || got[0].DHCPStop != "192.168.2.149" {
		t.Errorf("pool = %s-%s, want 192.168.2.100-192.168.2.149", got[0].DHCPStart, got[0].DHCPStop)
	}
}

// The whole point of a dry run: it must reach the device to report what would
// change, and must not change it.
func TestRunSetNetworkDryRunDoesNotWrite(t *testing.T) {
	s, c := mockClient(t, nil)

	if err := runSetNetwork(context.Background(), c, targetNetwork(), true); err != nil {
		t.Fatalf("dry run: %v", err)
	}

	got := s.Network()
	if len(got) != 1 || got[0].LANIP != "192.168.8.1" {
		t.Fatalf("a dry run changed the address: %+v", got)
	}
	if got[0].DHCPStart != "192.168.8.100" || got[0].DHCPStop != "192.168.8.249" {
		t.Errorf("a dry run changed the pool: %+v", got[0])
	}
}

func TestRunSetNetworkRefusedWithReservations(t *testing.T) {
	s, c := mockClient(t, []types.Reservation{
		{Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.13"},
		{Name: "pi", MAC: "aa:bb:cc:dd:ee:02", IP: "192.168.8.14"},
	})

	err := runSetNetwork(context.Background(), c, targetNetwork(), false)
	if !errors.Is(err, types.ErrReservationsExist) {
		t.Fatalf("error = %v, want ErrReservationsExist", err)
	}
	if got := s.Network(); got[0].LANIP != "192.168.8.1" {
		t.Errorf("a refused write moved the router to %s", got[0].LANIP)
	}
	// The operator needs to know the way out, not just that there is a wall.
	if !strings.Contains(err.Error(), "goglps --clear") {
		t.Errorf("error does not say how to proceed: %v", err)
	}
	if !strings.Contains(err.Error(), "2 reservation") {
		t.Errorf("error does not say how many are in the way: %v", err)
	}
}

// A dry run that approves what the real write would refuse is worse than none.
func TestRunSetNetworkDryRunReportsTheRefusal(t *testing.T) {
	_, c := mockClient(t, []types.Reservation{
		{Name: "nas", MAC: "aa:bb:cc:dd:ee:01", IP: "192.168.8.13"},
	})

	if err := runSetNetwork(context.Background(), c, targetNetwork(), true); !errors.Is(err, types.ErrReservationsExist) {
		t.Errorf("dry run error = %v, want ErrReservationsExist", err)
	}
}

// Likewise for validation: the firmware accepts a pool outside its subnet and
// then serves nothing, so the dry run has to be the thing that catches it.
func TestRunSetNetworkDryRunRejectsBadPool(t *testing.T) {
	_, c := mockClient(t, nil)

	n := targetNetwork()
	n.DHCPStart, n.DHCPStop = "10.0.0.100", "10.0.0.149"

	err := runSetNetwork(context.Background(), c, n, true)
	if !errors.Is(err, types.ErrOutsideSubnet) {
		t.Errorf("error = %v, want ErrOutsideSubnet", err)
	}
}
