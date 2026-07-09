// Команда goscout — сетевой scout, якорный проект курса.
//
// Сейчас это каркас: парсит подкоманду, валидирует аргументы.
// Реальные проверки появляются по мере прохождения тем A0.6–A0.9:
//   - dns:   каналы + context (A0.6)
//   - ports: интерфейсы + net.DialTimeout (A0.7 + A0.9)
//   - http:  stdlib + net/http (после A0)
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	version = "0.0.1"
	usage   = `goscout — a network scout.

Usage:
  goscout dns [--wordlist <file>] <domain> [<domain>...]   parallel DNS reconnaissance
  goscout ports --ports <list> <host>                      TCP port scanner
  goscout http <url>                                       HTTP probe (not implemented)
  goscout --version                                        print version

Flags must come before the host/domain arguments.
See cmd/goscout/README.md for details.
`
)

const (
	Workers = 100
	MinPort = 1
	MaxPort = 65535
)

type Result struct {
	Probe  string
	Target string
	Output []string
	Error  error
}

type Probe interface {
	Run(ctx context.Context, target string) Result
	Name() string
}

type DNSProbe struct{}
type PortProbe struct {
	Ports   []uint16
	Timeout time.Duration
}

type portResult struct {
	Port uint16
	Open bool
}

func runProbe(p Probe, targets []string, timeout time.Duration) {
	if len(targets) == 0 {
		fmt.Fprintf(os.Stderr, "%s: no targets provided\n", p.Name())
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ch := make(chan Result, len(targets))

	for _, target := range targets {
		go func(t string) {
			ch <- p.Run(ctx, t)
		}(target)
	}

	for range targets {
		select {
		case res := <-ch:
			if res.Error != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", res.Probe, res.Error)
				continue
			}
			fmt.Printf("%s -> %s\n", res.Target, strings.Join(res.Output, ", "))
		case <-ctx.Done():
			return
		}
	}
}

func (d DNSProbe) Name() string {
	return "dns"
}

func (d DNSProbe) Run(ctx context.Context, target string) Result {
	addrs, err := net.DefaultResolver.LookupHost(ctx, target)
	return Result{Probe: "dns", Target: target, Output: addrs, Error: err}
}

func (p PortProbe) Name() string {
	return "ports"
}

func (p PortProbe) Run(ctx context.Context, target string) Result {
	if len(p.Ports) == 0 {
		return Result{Probe: "ports", Target: target, Error: errors.New("no ports to scan")}
	}

	jobs := make(chan uint16, len(p.Ports))
	results := make(chan portResult, len(p.Ports))

	for _, port := range p.Ports {
		jobs <- port
	}
	close(jobs)

	var open []uint16
	for range Workers {
		go func() {
			for port := range jobs {
				addr := net.JoinHostPort(target, strconv.Itoa(int(port)))
				conn, err := net.DialTimeout("tcp", addr, p.Timeout)
				isOpen := err == nil
				if isOpen {
					conn.Close()
				}
				results <- portResult{Port: port, Open: isOpen}
			}
		}()
	}

	for i := 0; i < len(p.Ports); i++ {
		r := <-results
		if r.Open {
			open = append(open, r.Port)
		}
	}
	slices.Sort(open)

	var output []string
	for _, v := range open {
		output = append(output, strconv.Itoa(int(v)))
	}

	return Result{Probe: "ports", Target: target, Output: output, Error: nil}
}

func runHTTP(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "http: no url provided")
		os.Exit(2)
	}
	// TODO(post-A0): net/http клиент с таймаутом, разбор security headers, TLS info
	fmt.Println("http: not implemented yet")
	fmt.Printf("url: %s\n", args[0])
}

func readTargets(r io.Reader) ([]string, error) {
	var out []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func parsePorts(s string) ([]uint16, error) {
	ports := strings.Split(s, ",")

	var out []uint16
	for _, port := range ports {
		conv, err := strconv.Atoi(strings.TrimSpace(port))
		if err != nil {
			return nil, fmt.Errorf("invalid port %q: %w", port, err)
		}
		if conv >= MinPort && conv <= MaxPort {
			out = append(out, uint16(conv))
		} else {
			return nil, fmt.Errorf("port %d out of range (%d-%d)", conv, MinPort, MaxPort)
		}
	}
	return out, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "--version", "-v":
		fmt.Println(version)
	case "--help", "-h":
		fmt.Print(usage)
	case "dns":
		fs := flag.NewFlagSet(DNSProbe{}.Name(), flag.ExitOnError)
		timeout := fs.Duration("timeout", 5*time.Second, "resolve timeout")
		wordlist := fs.String("wordlist", "", "wordlist")
		fs.Parse(os.Args[2:])

		var targets []string
		if *wordlist != "" {
			f, err := os.Open(*wordlist)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dns: cannot open wordlist: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()

			targets, err = readTargets(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dns: cannot read targets: %v\n", err)
				os.Exit(1)
			}
		} else {
			targets = fs.Args()
		}
		runProbe(DNSProbe{}, targets, *timeout)
	case "ports":
		fs := flag.NewFlagSet(PortProbe{}.Name(), flag.ExitOnError)
		timeout := fs.Duration("timeout", 5*time.Second, "resolve timeout")
		portsRaw := fs.String("ports", "", "comma-separated ports, e.g. 80,443,22")
		fs.Parse(os.Args[2:])

		ports, err := parsePorts(*portsRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ports: %v\n", err)
			os.Exit(2)
		}
		targets := fs.Args()

		runProbe(PortProbe{Ports: ports, Timeout: *timeout}, targets, *timeout)
	case "http":
		runHTTP(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}
