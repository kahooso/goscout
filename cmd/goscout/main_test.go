package main

import (
	"context"
	"testing"
)

const (
	localhost = "localhost"
	invalid   = "nonexistent.invalid"
)

func TestResolveDomainTimeout(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr bool
	}{
		{
			name:    "expired ctx",
			domain:  localhost,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 0)
			defer cancel()
			ch := make(chan dnsResult, 1)
			resolveDomain(ctx, tc.domain, ch)
			res := <-ch
			if (res.err != nil) != tc.wantErr {
				t.Errorf("domain %q: err = %v, wantErr = %v", tc.domain, res.err, tc.wantErr)
			}
			if !tc.wantErr && len(res.addrs) == 0 {
				t.Errorf("domain %q: expected addrs, got nil list", tc.domain)
			}
		})
	}
}

func TestResolveDomain(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr bool
	}{
		{
			name:    "happy localhost",
			domain:  localhost,
			wantErr: false,
		},
		{
			name:    "invalid path",
			domain:  invalid,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ch := make(chan dnsResult, 1)
			resolveDomain(context.Background(), tc.domain, ch)
			res := <-ch
			if (res.err != nil) != tc.wantErr {
				t.Errorf("domain %q: err = %v, wantErr = %v", tc.domain, res.err, tc.wantErr)
			}
			if !tc.wantErr && len(res.addrs) == 0 {
				t.Errorf("domain %q: expected addrs, got nil list", tc.domain)
			}
		})
	}
}
