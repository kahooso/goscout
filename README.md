# goscout

Учебный репозиторий + якорный продукт студента **Информационной безопасности** (СПбПУ, профиль
«Технологии ИИ в кибербезопасности»). Связка Claude Code + Obsidian как обучающая среда:
Claude ведёт по плану (синхронно с университетом), Obsidian хранит конспекты,
`cmd/goscout/` — реальный инструмент, который растёт с каждой задачей.

## Вектор

**Цель к выпуску (2029):** Security Go developer со специализацией **AI/ML-безопасность** —
безопасный сетевой/прикладной софт на Go + понимание атак и защиты ML/LLM-систем.
Совпадает с профилем университета.

**Модель обучения:** универ = дорога (что и когда учить), Claude = зал (то же руками на Go,
без нейронки-за-меня). Ритм — следовать за университетом. Точка входа — стажировка
(AppSec / Security Go) ~2027. Раскладка по семестрам: [knowledge-base/university-plan.md](knowledge-base/university-plan.md).

## `goscout` — якорный продукт

**Сетевой security-сканер на Go + AI-модуль для анализа находок.** Конечная форма (к 7-8 семестру):
сканирует сеть (DNS / порты / HTTP / TLS) и применяет ML для классификации и приоритезации
результатов. Идём по версиям:

| Версия | Когда | Что умеет |
|--------|-------|-----------|
| **v0.1** | сем 4 (Сети) | dns + ports, тесты, CI — первый портфолио-проект |
| **v1.0** | сем 5 | + HTTP probe, TLS, security headers |
| **v2.0** | сем 7-8 (AI-security) | + ML-модуль: классификация/приоритезация находок |

Сейчас — Go-фундамент (блок A0), каждая фича добавляется когда тема пройдена.
См. [cmd/goscout/README.md](cmd/goscout/README.md).

```bash
go run ./cmd/goscout dns example.com google.com   # параллельный DNS (A0.6)
go run ./cmd/goscout ports example.com            # TCP scan (A0.9)
go run ./cmd/goscout http example.com             # HTTP probe (v1)
go run ./cmd/goscout --version
```

## Структура

```
goscout/
├── go.mod                — модуль github.com/kahooso/goscout
├── .github/workflows/    — CI: go vet + go test -race на push
│
├── cmd/                  — исполняемые бинарники
│   ├── goscout/          — якорный продукт (см. cmd/goscout/README.md)
│   ├── tasks/<имя>/      — algo-практика (package = имя, через go test)
│   └── <учебные>/        — артефакты задач A0.1–A0.5 (structs, collections,
│                           logparse, pointers, errors-demo, goroutines)
│
├── knowledge-base/       — Obsidian vault
│   ├── 00-roadmap.md       — журнал задач, мок-собесов, CTF
│   ├── university-plan.md  — учебный план ↔ синхронизация с goscout
│   ├── tasks/              — разбор каждого задания (опыт, ошибки)
│   ├── topics/             — концепции с [[wikilinks]] (_index.md — карта)
│   ├── ctf/                — picoCTF: сессии + writeup'ы (Web+General)
│   └── interviews/         — мок-собесы
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
go run ./cmd/goscout  # якорный продукт
```

## Как читать репозиторий

- **Как работать со связкой Claude + проект** — [.claude/HOWTO.md](.claude/HOWTO.md)
- **Куда идёт обучение** — [.claude/strategy/learning-strategy.md](.claude/strategy/learning-strategy.md) → «Главный вектор»
- **Где сейчас в плане** — там же → «Текущая позиция»
- **Синхронизация с универом** — [knowledge-base/university-plan.md](knowledge-base/university-plan.md)
- **Как Claude себя ведёт** — [.claude/CLAUDE.md](.claude/CLAUDE.md) + скиллы рядом
- **Разобранные темы** — [knowledge-base/topics/_index.md](knowledge-base/topics/_index.md)
- **Выполненные задачи и грабли** — [knowledge-base/tasks/](knowledge-base/tasks/)
