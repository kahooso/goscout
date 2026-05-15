// Команда goscout — сетевой scout, якорный проект курса.
//
// Сейчас это каркас: парсит подкоманду, валидирует аргументы.
// Реальные проверки появляются по мере прохождения тем A0.6–A0.9:
//   - dns:   каналы + context (A0.6)
//   - ports: интерфейсы + net.DialTimeout (A0.7 + A0.9)
//   - http:  stdlib + net/http (после A0)
package main

import (
	"fmt"
	"os"
)

const version = "0.0.1"

const usage = `goscout — сетевой scout.

Использование:
  goscout dns <domain> [<domain>...]   параллельный DNS reconnaissance
  goscout ports <host>                 TCP port scanner (заглушка)
  goscout http <url>                   HTTP probe (заглушка)
  goscout --version                    вывести версию

Подкоманды реализуются по мере прохождения тем — см. cmd/goscout/README.md.
`

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
		runDNS(os.Args[2:])
	case "ports":
		runPorts(os.Args[2:])
	case "http":
		runHTTP(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "неизвестная подкоманда: %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func runDNS(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "dns: нужен минимум один домен")
		os.Exit(2)
	}
	// TODO(A0.6): параллельный resolve через горутины + каналы + context
	fmt.Println("dns: not implemented yet (ждёт A0.6)")
	fmt.Printf("домены для проверки: %v\n", args)
}

func runPorts(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "ports: нужен host")
		os.Exit(2)
	}
	// TODO(A0.9): TCP DialTimeout + worker pool через буферизованный канал
	fmt.Println("ports: not implemented yet (ждёт A0.9)")
	fmt.Printf("host: %s\n", args[0])
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
