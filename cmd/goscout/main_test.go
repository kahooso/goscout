package main

import (
	"context"
	"testing"
)

const (
	localhost   = "localhost"
	invalid     = "nonexistent.invalid"
	placeholder = "127.0.0.1"
)

func TestProbeName(t *testing.T) {
	tests := []struct {
		probe Probe
		want  string
	}{
		{probe: DNSProbe{}, want: "dns"},
		{probe: PortProbe{}, want: "ports"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.probe.Name(); got != tc.want {
				t.Errorf("Name() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDNSProbeRun(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{
			name:    "happy localhost",
			target:  localhost,
			wantErr: false,
		},
		{
			name:    "invalid path",
			target:  invalid,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := DNSProbe{}.Run(context.Background(), tc.target)
			if (res.Error != nil) != tc.wantErr {
				t.Errorf("domain %q: err = %v, wantErr = %v", tc.target, res.Error, tc.wantErr)
			}
			if !tc.wantErr && len(res.Output) == 0 {
				t.Errorf("domain %q: expected addrs, got nil list", tc.target)
			}
		})
	}
}

func TestDNSProbeRunTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	res := DNSProbe{}.Run(ctx, localhost)
	if res.Error == nil {
		t.Errorf(" ctx timeout: expected error, got nil")
	}
}

func TestPortProbeRun(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{
			name: "not implemented",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := PortProbe{}.Run(context.Background(), placeholder)
			if res.Error == nil {
				t.Errorf("PortProbe.Run: expected error, got nil")
			}
		})
	}
}
