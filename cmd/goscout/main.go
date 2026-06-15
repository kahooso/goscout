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

/*

 */

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

func runProbe(p Probe, targets []string, timeout time.Duration) {
	// A0.6 - параллельный resolve через горутины + каналы + context
	if len(targets) == 0 {
		fmt.Fprintf(os.Stderr, "%s: нужен минимум один домен\n", p.Name())
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

func readTargets(r io.Reader) ([]string, error) {
	var out []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text() /* have no -> '\n', '\r' */
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
		// 3. Разобрать флаги этой команды → создаю FlagSet, регистрирую --timeout, --wordlist, Parse
		// 4. Понять ОТКУДА брать домены:
		//      - если задан --wordlist → открыть файл, прочитать через readTargets
		//      - иначе                 → взять позиционные fs.Args()
		// 5. Если доменов нет вообще     → ошибка "нужен домен или --wordlist"
		fs := flag.NewFlagSet("dns", flag.ExitOnError)

		timeout := fs.Duration("timeout", 5*time.Second, "resolve timeout")
		wordlist := fs.String("wordlist", "", "wordlist")
		fs.Parse(os.Args[2:]) // A: да — [0]=goscout, [1]=dns отрезаны, дальше флаги+позиционные

		var targets []string
		if *wordlist != "" {
			f, err := os.Open(*wordlist)
			if err != nil {
				// A: да, ошибку проверяем всегда: файла нет/нет прав → f невалиден,
				// читать нечего. Печатаем причину в stderr и выходим (runtime error → 1).
				fmt.Fprintf(os.Stderr, "cannot open wordlist: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()

			targets, err = readTargets(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot read targets: %v\n", err)
				os.Exit(1)
			}
		} else {
			targets = fs.Args()
		}
		// A: коды выхода — конвенция 0/1/2: 0 успех, 1 runtime (файл/сеть),
		// 2 usage (кривой ввод). Подробно: topics/go/stdlib-cli.md
		runProbe(DNSProbe{}, targets, *timeout)
	case "ports":
		runProbe(PortProbe{}, os.Args[2:], 5*time.Second)
	case "http":
		runHTTP(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcomm: %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}
