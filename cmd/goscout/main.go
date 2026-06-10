// Команда goscout — сетевой scout, якорный проект курса.
//
// Сейчас это каркас: парсит подкоманду, валидирует аргументы.
// Реальные проверки появляются по мере прохождения тем A0.6–A0.9:
//   - dns:   каналы + context (A0.6)
//   - ports: интерфейсы + net.DialTimeout (A0.7 + A0.9)
//   - http:  stdlib + net/http (после A0)
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	version = "0.0.1"
	usage   = `goscout — сетевой scout.

	Использование:
	goscout dns <domain> [<domain>...]   параллельный DNS reconnaissance
	goscout ports <host>                 TCP port scanner (заглушка)
	goscout http <url>                   HTTP probe (заглушка)
	goscout --version                    вывести версию

	Подкоманды реализуются по мере прохождения тем — см. cmd/goscout/README.md.
	`
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
type PortProbe struct{}

func runProbe(p Probe, targets []string) {
	// A0.6 - параллельный resolve через горутины + каналы + context
	if len(targets) == 0 {
		fmt.Fprintf(os.Stderr, "%s: нужен минимум один домен\n", p.Name())
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
				fmt.Printf("%v\n", res.Error)
				continue
			}
			fmt.Printf("%s -> %s\n", res.Target, res.Output)
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
	return Result{Probe: "ports", Target: target, Error: errors.New("ports: not implemented yet")}
}

func runHTTP(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "http: нужен URL")
		os.Exit(2)
	}
	// TODO(post-A0): net/http клиент с таймаутом, разбор security headers, TLS info
	fmt.Println("http: not implemented yet (ждёт post-A0)")
	fmt.Printf("url: %s\n", args[0])
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
		runProbe(DNSProbe{}, os.Args[2:])
	case "ports":
		runProbe(PortProbe{}, os.Args[2:])
	case "http":
		runHTTP(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcomm: %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}
