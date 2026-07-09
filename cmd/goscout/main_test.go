package main

import (
	"context"
	"net"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
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
	res := PortProbe{}.Run(context.Background(), placeholder)
	if res.Error == nil {
		t.Errorf("Run(%q).Error = nil, want error", placeholder)
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
		})
	}
}

func TestPortProbeRunOpen(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listener cannot be created: %v", err)
	}
	defer ln.Close()

	port := uint16(ln.Addr().(*net.TCPAddr).Port)

	p := PortProbe{Ports: []uint16{port}, Timeout: time.Second}
	res := p.Run(context.Background(), placeholder)

	if res.Error != nil {
		t.Fatalf("PortProbe.Run(): %v", res.Error)
	}

	want := strconv.Itoa(int(port))
	if len(res.Output) != 1 || res.Output[0] != want {
		t.Errorf("Output = %v, want = [%s]", res.Output, want)
	}
}

func TestPortProbeRunClosed(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listener cannot be created: %v", err)
	}
	port := uint16(ln.Addr().(*net.TCPAddr).Port)
	ln.Close()

	p := PortProbe{Ports: []uint16{port}, Timeout: time.Second}
	res := p.Run(context.Background(), placeholder)

	if res.Error != nil {
		t.Fatalf("closed port: Error = %v, want nil", res.Error)
	}
	if len(res.Output) != 0 {
		t.Errorf("closed port: Output = %v, want empty", res.Output)
	}
}

func TestPortProbeRunSorted(t *testing.T) {
	var listeners []net.Listener
	for id := range 4 {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listener #%d cannot be created: %v", id, err)
		}
		listeners = append(listeners, ln)
	}

	var ports []uint16
	for _, v := range listeners {
		ports = append(ports, uint16(v.Addr().(*net.TCPAddr).Port))
	}

	for _, ln := range listeners {
		defer ln.Close()
	}

	p := PortProbe{Ports: ports, Timeout: time.Second}

	res := p.Run(context.Background(), placeholder)
	if res.Error != nil {
		t.Fatalf("PortProbe.Run(): %v", res.Error)
	}

	var temp []uint16
	for _, p := range ports {
		temp = append(temp, p)
	}
	slices.Sort(temp)

	var output []string
	for _, tp := range temp {
		output = append(output, strconv.Itoa(int(tp)))
	}
	if !reflect.DeepEqual(output, res.Output) {
		t.Errorf("Output = %v, want = [%v]", res.Output, output)
	}
}
