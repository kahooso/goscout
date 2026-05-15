# goscout

Сетевой scout — якорный учебный проект. CLI с подкомандами для разведки.

## Назначение

`goscout` — это не клон nmap или httpx. Это **учебный** инструмент: каждая фича
добавляется когда нужная тема пройдена в блоке A0/A. Цель — к концу A0 иметь
рабочий `dns` + `ports` (v0.1) и публиковать на GitHub.

## Состояние подкоманд

| Команда | Тема | Статус |
|---------|------|--------|
| `goscout dns <domain>...` | A0.6 каналы + context | каркас, TODO |
| `goscout ports <host>` | A0.7 интерфейсы + A0.9 net | каркас, TODO |
| `goscout http <url>` | пост-A0: net/http, TLS | каркас, TODO |
| `goscout --version` | — | работает |
| `goscout --help` | — | работает |

## Что появится в каждой задаче

### A0.6 — `dns`
Параллельный resolve N доменов с таймаутом.
- N горутин-резолверов, по одной на домен
- Канал результатов `chan dnsResult`
- `context.WithTimeout` для остановки зависших lookups
- `select { case r := <-results: ... case <-ctx.Done(): ... }`

### A0.7 — `Probe` интерфейс
Общий тип для всех проверок:
```go
type Probe interface {
    Run(ctx context.Context, target string) (Result, error)
    Name() string
}
```
`DNSProbe` уже есть; `PortProbe`, `HTTPProbe` идут позже.

### A0.8 — флаги и wordlist
`--timeout`, `--wordlist` для DNS перебора (subdomain enum).
`bufio.Scanner` для чтения wordlist построчно.

### A0.9 — `ports`
TCP port scan через `net.DialTimeout`.
- Worker pool: буферизованный канал заданий, N воркеров
- Параллельность ограничена `--workers` флагом

## Запуск

```bash
go run ./cmd/goscout --version
go run ./cmd/goscout dns example.com google.com
go run ./cmd/goscout ports example.com
```

## Тесты

```bash
go test -race ./cmd/goscout
```

Тестов пока нет — появятся вместе с реальной логикой в задаче A0.6.
