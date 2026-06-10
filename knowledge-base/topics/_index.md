# Карта тем (MOC)

Индекс всех концепций в vault. Обновляется при каждом добавлении файла в `topics/`.

## Go

| Тема | Статус | Связанный task |
|------|--------|----------------|
| [[structs]] — структуры и методы | ✅ есть | task-01 |
| [[slices-maps]] — слайсы и мэпы | ✅ есть | task-02 |
| [[strings-sort]] — пакеты strings и sort | ✅ есть | task-02 |
| [[hash-table]] — хеш-таблица vs дерево | ✅ есть | task-02 (бонус) |
| [[pointers]] — указатели, value vs pointer receiver | ✅ есть | task-04 |
| [[errors]] — `errors`, `fmt.Errorf`, `errors.Is/As` | ✅ есть | task-05 |
| [[goroutines]] — горутины, sync.WaitGroup | ✅ есть (новый формат) | task-06 |
| [[testing]] — пакет `testing`, table-driven тесты | ✅ есть | task-03+ |
| [[runtime]] — Go runtime, компиляция, stdlib | ✅ есть | сквозная |
| [[channels]] — каналы, `select`, buffered/unbuffered | ✅ есть (новый формат) | task-07 |
| [[context]] — `context.WithTimeout` и lifecycle | ✅ есть (новый формат) | task-07 |
| [[go-tooling]] — команды go, ритуал перед коммитом | ✅ есть (новый формат) | task-07 |
| [[interfaces]] — интерфейсы, неявная реализация | ✅ есть (новый формат) | task-08 |
| [[stdlib-cli]] — `os`, `bufio`, `flag` для CLI | 🔲 в плане | A0.8 |
| [[net]] — `net.Dial`, `net.DialTimeout` | 🔲 в плане | A0.9 |

## Алгоритмы

Параллельная ветка — мини-задачи для закрепления базы знаний.
Подробнее — [strategy/learning-strategy.md](../../.claude/strategy/learning-strategy.md) → раздел «Algo-практика».

| Паттерн | Статус | Связанный algo-task |
|---------|--------|---------------------|
| [[map-as-index]] — поиск O(N) через map[T]int | ✅ есть | algo-01 (Two Sum) |
| [[two-pointers]] — два указателя на одном слайсе | 🔲 в плане | будет |
| [[sliding-window]] — окно фиксированного/переменного размера | 🔲 в плане | будет |
| [[prefix-sum]] — префиксные суммы для range queries | 🔲 в плане | будет |
| [[linked-list]] — связные списки на указателях | 🔲 в плане | будет |
| [[bst]] — бинарные деревья поиска | 🔲 в плане | будет |

## Сети

> Заполнится после A0 (фичи `dns`/`ports`/`http` в goscout) и из сетевых CTF-челленджей.
> Сейчас: пусто.

## Безопасность

> Питается из двух источников: фича `goscout http` (security headers, TLS) после A0
> и **CTF-трек** (Web Exploitation — SQLi/XSS/SSRF/JWT). Каждый web-writeup привязан
> к OWASP-категории; когда наберётся 5–6 — свести в `topics/security/owasp-top10.md`.
> CTF-журнал и сессии — `knowledge-base/ctf/picoctf/`. Сейчас: пусто.

---

## Соглашения

- Wikilink `[[name]]` ведёт на файл `<name>.md` в той же папке `topics/<категория>/`.
- Если ссылка ведёт на ещё не созданный файл — это нормально, помечается 🔲.
- Каждая тема рождается из конкретной задачи или вопроса — не пишем «впрок».
- **Новый формат записи** (с goroutines.md и testing.md): TL;DR → ссылка на ресурсы → код → Личный опыт.
  Старые файлы будут переписаны при возвращении к теме.