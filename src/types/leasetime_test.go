package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestLeaseTimeUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  LeaseTime
	}{
		{"hours", `"12h"`, LeaseTime(12 * time.Hour)},
		{"days", `"1d"`, LeaseTime(24 * time.Hour)},
		{"weeks", `"2w"`, LeaseTime(14 * 24 * time.Hour)},
		{"minutes", `"30m"`, LeaseTime(30 * time.Minute)},
		{"seconds suffix", `"90s"`, LeaseTime(90 * time.Second)},
		{"bare seconds as number", `86400`, LeaseTime(24 * time.Hour)},
		{"bare seconds as string", `"86400"`, LeaseTime(24 * time.Hour)},
		{"infinite", `"infinite"`, LeaseInfinite},
		{"compound", `"1h30m"`, LeaseTime(90 * time.Minute)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got LeaseTime
			if err := json.Unmarshal([]byte(tt.input), &got); err != nil {
				t.Fatalf("Unmarshal(%s) error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("Unmarshal(%s) = %v, want %v", tt.input, time.Duration(got), time.Duration(tt.want))
			}
		})
	}
}

func TestLeaseTimeUnmarshalError(t *testing.T) {
	for _, input := range []string{`"12x"`, `"forever"`, `""`, `[]`, `{}`, `"xd"`} {
		var got LeaseTime
		if err := json.Unmarshal([]byte(input), &got); err == nil {
			t.Errorf("Unmarshal(%s) succeeded, want error", input)
		}
	}
}

func TestLeaseTimeString(t *testing.T) {
	tests := []struct {
		lease LeaseTime
		want  string
	}{
		{LeaseTime(12 * time.Hour), "12h"},
		{LeaseTime(24 * time.Hour), "24h"},
		{LeaseTime(30 * time.Minute), "30m"},
		{LeaseTime(90 * time.Second), "1m30s"},
		{LeaseTime(45 * time.Second), "45s"},
		{LeaseInfinite, "infinite"},
		{LeaseTime(0), "0s"},
	}
	for _, tt := range tests {
		if got := tt.lease.String(); got != tt.want {
			t.Errorf("LeaseTime(%d).String() = %q, want %q", int64(tt.lease), got, tt.want)
		}
	}
}

func TestLeaseTimeMarshalRoundTrip(t *testing.T) {
	for _, want := range []LeaseTime{
		LeaseTime(12 * time.Hour),
		LeaseInfinite,
		LeaseTime(90 * time.Second),
		LeaseTime(30 * time.Minute),
		LeaseTime(0),
	} {
		b, err := json.Marshal(want)
		if err != nil {
			t.Fatalf("Marshal(%v) error: %v", want, err)
		}
		var got LeaseTime
		if err := json.Unmarshal(b, &got); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", b, err)
		}
		if got != want {
			t.Errorf("round trip via %s = %v, want %v", b, time.Duration(got), time.Duration(want))
		}
	}
}

func TestLeaseTimeSeconds(t *testing.T) {
	if got := LeaseTime(12 * time.Hour).Seconds(); got != 43200 {
		t.Errorf("Seconds() = %d, want 43200", got)
	}
	if got := LeaseInfinite.Seconds(); got != -1 {
		t.Errorf("LeaseInfinite.Seconds() = %d, want -1", got)
	}
}
