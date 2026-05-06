# goscout

Учебный репозиторий для подготовки к роли **AppSec инженера / Security Go developer**.
Связка Claude Code + Obsidian как обучающая среда: Claude ведёт по плану, Obsidian хранит конспекты.

## Цель

К концу 3 курса (лето–осень 2028) — junior-оффер в Positive Technologies, Bi.Zone, Яндекс Безопасность или Kaspersky.
Промежуточная цель этапа 1 (конец 2026) — рабочий **HTTP Security Scanner на Go**.

## Структура

```
goscout/
├── go.mod                 — единый Go модуль github.com/kahooso/goscout
├── cmd/                   — исполняемые CLI-инструменты, по одному на бинарник
│   ├── structs/           — task-01: тип Task с методами (A0.1)
│   └── collections/       — task-02: счётчик слов (A0.2)
│
├── knowledge-base/        — Obsidian vault: теория и журнал задач
│   ├── 00-roadmap.md      — журнал выполненных задач
│   ├── tasks/             — разбор каждого задания (опыт, ошибки, рефлексия)
│   ├── topics/            — концепции с [[wikilinks]]
│   │   ├── _index.md      — карта всех тем (MOC)
│   │   ├── go/            — Go: пакеты, паттерны, идиомы
│   │   ├── networks/      — сетевой стек, протоколы (наполнится в блоке C)
│   │   └── security/      — AppSec: уязвимости, атаки, защита (наполнится в блоке B)
│   └── interviews/        — мок-собесы каждые ~8 задач
│
└── .claude/               — инструкции для Claude Code
    ├── CLAUDE.md          — мастер-контекст (загружается автоматически)
    ├── strategy/          — учебный план и карьерная стратегия
    ├── skills/            — поведение Claude в конкретных ситуациях
    └── templates/         — шаблоны новых файлов
```

Папки `cmd/scanner/`, `topics/networks/`, `topics/security/` появятся когда дойдём до соответствующих блоков плана.

## Текущий статус

- **Этап:** 1 (фундамент)
- **Блок:** A0 — Go фундамент
- **Готово:** A0.1 (структуры), A0.2 (слайсы и мэпы) — закрепление в работе
- **Подробности:** [`learning-strategy.md`](.claude/strategy/learning-strategy.md) → раздел «Текущая позиция»

## Запуск

```bash
go run ./cmd/structs
go run ./cmd/collections
go test ./...
go vet ./...
```

## Как читать репозиторий

- Чтобы понять **где сейчас идёт обучение** — `.claude/strategy/learning-strategy.md`.
- Чтобы понять **как устроен учебный процесс** — `.claude/CLAUDE.md` и скиллы рядом.
- Чтобы посмотреть **разобранные темы** — `knowledge-base/topics/_index.md`.
- Чтобы увидеть **какие задачи выполнены и с какими граблями** — `knowledge-base/tasks/`.
