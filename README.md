# goscout

Учебный репозиторий студента **Информационной безопасности** (СПбПУ, профиль
«Технологии ИИ в кибербезопасности»). Связка Claude Code + Obsidian как обучающая среда:
Claude ведёт по плану (синхронно с университетом), Obsidian хранит конспекты, каждая
задача разбирается и фиксируется в vault.

Продукт, который из этого вырос, живёт отдельно: **[outpost](https://github.com/kahooso/outpost)**
— сетевой recon CLI на Go. Вынесен из этого репозитория 2026-07-28.

## Вектор

**Дорога:** product / infra / deep-tech backend на Go. Кибербезопасность — сильный сквозной
угол (профиль университета — «Технологии ИИ в кибербезопасности») и реальная точка входа:
стажировка AppSec / Security Go (~2027). Один из путей на широкой дороге, не единственная цель.
Карьерный вектор подробно — в отдельном приватном пространстве размышлений (goagent).

**Модель обучения:** университет даёт теорию (контекст тем, что сейчас в потоке, не драйвер),
работа с Claude — практика на Go руками, без нейронки-за-меня. Темы беру по вектору, синхронизируя с университетом.
Раскладка по семестрам: [knowledge-base/university-plan.md](knowledge-base/university-plan.md).

## `outpost` — якорный продукт (отдельный репозиторий)

**Сетевой recon / диагностический CLI на Go с security-уклоном.** Сканирует сеть
(DNS / порты / HTTP / TLS), показывает что открыто, как настроено, где небезопасно.

Жил здесь как `cmd/goscout/` с мая по июль 2026, вынесен в собственный публичный
репозиторий: **github.com/kahooso/outpost**. Здесь остаётся только учебная часть —
конспекты, разбор задач, мок-собесы.

| Веха | Что умеет | Статус |
|------|-----------|--------|
| **v0.1** | dns + ports, тесты, CI | ✅ работает |
| **v1.0** | + HTTP probe, TLS, security headers | дальше |
| **дальше** | открыто — по мере тем и интереса | растёт со мной |

## Структура

```
goscout/
├── go.mod                — модуль github.com/kahooso/goscout
├── .github/workflows/    — CI: go vet + go test -race на push
│
├── cmd/                  — учебные артефакты (продукт переехал в outpost)
│   ├── tasks/<имя>/      — algo-практика (package = имя, через go test)
│   ├── retest-a06/       — blank-slate ретест A0.6 (гейт закрытия блока A0)
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
```

## Как читать репозиторий

- **Как работать со связкой Claude + проект** — [.claude/HOWTO.md](.claude/HOWTO.md)
- **Куда идёт обучение** — [.claude/strategy/learning-strategy.md](.claude/strategy/learning-strategy.md) → «Главный вектор»
- **Где сейчас в плане** — там же → «Текущая позиция»
- **Синхронизация с универом** — [knowledge-base/university-plan.md](knowledge-base/university-plan.md)
- **Как Claude себя ведёт** — [.claude/CLAUDE.md](.claude/CLAUDE.md) + скиллы рядом
- **Разобранные темы** — [knowledge-base/topics/_index.md](knowledge-base/topics/_index.md)
- **Выполненные задачи и грабли** — [knowledge-base/tasks/](knowledge-base/tasks/)
