# goscout

Учебный репозиторий + якорный проект для подготовки к **AppSec инженеру / Security Go developer**.
Связка Claude Code + Obsidian как обучающая среда: Claude ведёт по плану, Obsidian хранит конспекты,
`cmd/goscout/` — реальный инструмент который растёт с каждой задачей.

## Цель

- **2026:** `goscout` v0.1 с командами `dns` + `ports`, тесты + CI
- **2027:** стажировка AppSec / Security Go (Positive Technologies, Bi.Zone, Яндекс, Kaspersky)
- **2028:** junior-оффер

## `goscout` — якорный проект

Сетевой scout, CLI инструмент:

```bash
go run ./cmd/goscout dns example.com google.com   # параллельный DNS (A0.6)
go run ./cmd/goscout ports example.com            # TCP scan (A0.9)
go run ./cmd/goscout http example.com             # HTTP probe (post-A0)
go run ./cmd/goscout --version
```

Каждая фича добавляется когда нужная тема пройдена — см. [cmd/goscout/README.md](cmd/goscout/README.md).

## Структура

```
goscout/
├── go.mod                — модуль github.com/kahooso/goscout
├── .github/workflows/    — CI: go vet + go test -race на push
├── .gitignore            — .exe, .obsidian, vendor и т.д.
│
├── cmd/                  — исполняемые бинарники
│   ├── goscout/          — якорный проект (см. cmd/goscout/README.md)
│   ├── structs/          — task-01 (A0.1) — артефакт
│   ├── collections/      — task-02 (A0.2) — артефакт
│   ├── logparse/         — task-03 (A0.2-bis) — артефакт
│   ├── pointers/         — task-04 (A0.3) — артефакт
│   ├── errors-demo/      — task-05 (A0.4) — артефакт
│   └── goroutines/       — task-06 (A0.5) — артефакт
│
├── knowledge-base/       — Obsidian vault
│   ├── 00-roadmap.md     — журнал задач
│   ├── tasks/            — разбор каждого задания (опыт, ошибки)
│   ├── topics/           — концепции с [[wikilinks]]
│   │   ├── _index.md     — карта всех тем
│   │   └── go/           — Go: пакеты, идиомы, личный опыт
│   └── interviews/       — мок-собесы каждые ~8 задач
│
└── .claude/              — инструкции для Claude Code
    ├── CLAUDE.md         — мастер-контекст (загружается автоматически)
    ├── skills/           — реакции на конкретные триггеры
    ├── strategy/         — учебный план
    └── templates/        — шаблоны
```

## Запуск

```bash
go test -race ./...   # все тесты с детектором гонок
go vet ./...          # статический анализ
go run ./cmd/goscout  # якорный проект
```

## Как читать репозиторий

- **Как работать со связкой Claude + проект** — [.claude/HOWTO.md](.claude/HOWTO.md) — шпаргалка
- **Где сейчас идёт обучение** — [.claude/strategy/learning-strategy.md](.claude/strategy/learning-strategy.md) → «Текущая позиция»
- **Как устроен учебный процесс** — [.claude/CLAUDE.md](.claude/CLAUDE.md) + скиллы рядом
- **Разобранные темы** — [knowledge-base/topics/_index.md](knowledge-base/topics/_index.md)
- **Выполненные задачи и грабли** — [knowledge-base/tasks/](knowledge-base/tasks/)
