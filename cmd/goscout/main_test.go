package main

import (
	"context"
	"reflect"
	"strings"
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
				t.Errorf("Run(%q).Error = %v, wantErr %v", tc.target, res.Error, tc.wantErr)
			}
			if !tc.wantErr && len(res.Output) == 0 {
				t.Errorf("Run(%q).Output is empty, want non-empty", tc.target)
			}
		})
	}
}

func TestDNSProbeRunTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	res := DNSProbe{}.Run(ctx, localhost)
	if res.Error == nil {
		t.Errorf("Run(%q) with expired ctx: Error = nil, want error", localhost)
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
				t.Errorf("Run(%q).Error = nil, want error", placeholder)
			}
		})
	}
}

func TestReadTargets(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			"happy",
			"a.com\nb.com\n",
			[]string{"a.com", "b.com"},
		},
		{
			"50/50",
			"a.com\n\n\nb.com\n",
			[]string{"a.com", "b.com"},
		},
		{
			"empty",
			"",
			nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			targets, err := readTargets(strings.NewReader(tc.input))
			if err != nil {
				t.Errorf("readTargets(%q): unexpected error %v", tc.input, err)
				return
			}
			if !reflect.DeepEqual(targets, tc.want) {
				t.Errorf("readTargets(%q) = %v, want %v", tc.input, targets, tc.want)
			}
			/*
				f := func(x, y []string) bool {
					if len(x) != len(y) {
						return false
					}
					for i := range x {
						if x[i] != y[i] {
							return false
						}
					}
					return true
				}
				if !reflect.DeepEqual(tc.want, targets) {
					t.Errorf("ReadTargets: want %v, got %v", tc.want, targets)
				}
			*/
		})
	}
}
