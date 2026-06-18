# goscout

Сетевой security-сканер — якорный продукт. CLI с подкомандами для разведки.
Конечная цель (v2, сем 7-8): сканер + AI-модуль для классификации находок.
Общий вектор — в корневом [README](../../README.md) и `knowledge-base/university-plan.md`.

## Назначение

`goscout` — это не клон nmap или httpx. Это **учебный** инструмент с названной целью:
каждая фича добавляется когда нужная тема пройдена. Ближайшая веха — к концу A0 иметь
рабочий `dns` + `ports` (v0.1) и опубликовать на GitHub.

## Состояние подкоманд

| Команда | Тема | Статус |
|---------|------|--------|
| `goscout dns <domain>...` | A0.6 каналы + context, A0.8 CLI | ✅ работает (горутины + context, флаги `--timeout`/`--wordlist`) |
| `goscout ports <host>` | A0.7 интерфейсы + A0.9 net | каркас (`PortProbe` пустая, заглушка `not implemented yet`) |
| `goscout http <url>` | пост-A0: net/http, TLS | заглушка |
| `goscout --version` / `-v` | — | работает |
| `goscout --help` / `-h` | — | работает |

## Что осталось до v0.1

### A0.9 — `ports` (в работе)
TCP port scan через `net.DialTimeout`.
- Worker pool: буферизованный канал заданий, N воркеров
- Параллельность ограничена `--workers` флагом
- Три исхода: OPEN (SYN-ACK) / CLOSED (RST) / FILTERED (timeout) — см. `topics/networks/tcp-ip.md`

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

Покрытие на текущий момент: `Probe.Name()`, `DNSProbe.Run` (happy + invalid + истёкший
context), `PortProbe.Run` (заглушка возвращает ошибку), `readTargets` (happy / пустые
строки / пустой ввод через `strings.NewReader`).

## Сборка

```bash
go build -o goscout ./cmd/goscout      # Linux/macOS
go build -o goscout.exe ./cmd/goscout  # Windows
```
