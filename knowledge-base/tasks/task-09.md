---
id: task-09
block: A0.8
topic: stdlib для CLI (os, bufio, flag, io.Reader)
status: completed
date_started: 2026-06-10
date_completed: 2026-06-15
---

# Task-09: флаги и чтение wordlist (CLI для goscout)

## Задание

Дать goscout настоящий CLI: флаги `--timeout` и `--wordlist` для подкоманды `dns`,
чтение списка доменов из файла построчно. Заменить захардкоженный таймаут и ручной
сбор целей.

**Файл:** `cmd/goscout/main.go`. **Тесты:** `cmd/goscout/main_test.go`, `-race`.

### Что должно получиться

```
goscout dns --timeout 3s example.com github.com      # домены позиционно + флаг таймаута
goscout dns --wordlist domains.txt --timeout 2s       # домены из файла построчно
```

### Требования

1. **Функция чтения wordlist** — `readTargets(r io.Reader) ([]string, error)`:
   - принимает **`io.Reader`** (НЕ `*os.File`, НЕ имя файла — это ключ к тестам);
   - читает построчно через `bufio.Scanner`;
   - возвращает слайс непустых строк (пропускать пустые строки);
   - после цикла проверить `scanner.Err()`.

2. **Флаги для `dns`** через `flag.FlagSet` (отдельный набор на подкоманду):
   ```go
   fs := flag.NewFlagSet("dns", flag.ExitOnError)
   timeout := fs.Duration("timeout", 5*time.Second, "таймаут на резолв")
   wordlist := fs.String("wordlist", "", "файл со списком доменов")
   fs.Parse(args)            // args = os.Args[2:] (после "dns")
   targets := fs.Args()      // позиционные домены (не-флаги)
   ```
   Почему `FlagSet`, а не глобальный `flag`: у `dns`/`ports` свои флаги, глобальный
   `flag.Parse()` не знает про подкоманды. `FlagSet` = изолированный набор.

3. **Источник целей:** если задан `--wordlist` → открыть файл (`os.Open`, `defer Close`),
   прочитать через `readTargets`. Иначе → взять `fs.Args()` (позиционные). Если пусто —
   ошибка «нужен домен или --wordlist».

4. **Прокинуть таймаут:** `runProbe` сейчас хардкодит `5*time.Second`. Добавить параметр
   `timeout time.Duration`, передавать из флага. `runProbe(p, targets, timeout)`.

5. **Вывод** — как было (домен → адреса / ошибка).

### Тесты (table-driven, -race)

Главное — `readTargets` тестируется БЕЗ файла, через `strings.NewReader`:
- happy: `strings.NewReader("a.com\nb.com\n")` → `["a.com", "b.com"]`;
- пустые строки: `"a.com\n\n\nb.com\n"` → `["a.com", "b.com"]` (пустые отброшены);
- пустой ввод: `""` → пустой слайс, без ошибки.

Существующие тесты проб (TestProbeName, TestDNSProbeRun...) не ломать.

---

## Моё решение

- `readTargets(r io.Reader) ([]string, error)` — `bufio.NewScanner(r)`, цикл
  `for sc.Scan()`, пропуск пустых (`if line == ""`), в конце `if err := sc.Err(); err != nil`.
  Принимает `io.Reader`, НЕ `*os.File` → тестируется без диска.
- Ветка `dns`: `flag.NewFlagSet("dns", ExitOnError)` → `fs.Duration("timeout",...)` +
  `fs.String("wordlist",...)` → `fs.Parse(os.Args[2:])`. Источник целей: если `*wordlist != ""`
  → `os.Open` + `defer Close` + `readTargets(f)`; иначе → `fs.Args()` (позиционные).
- Ошибки файла → `fmt.Fprintf(os.Stderr, "cannot open/read ...: %v", err)` + `os.Exit(1)`.
- `runProbe` получил параметр `timeout time.Duration`, в `main` передаём `*timeout`
  (разыменование указателя из flag). Хардкод `5*time.Second` в `context.WithTimeout` убран.
- Тесты: `TestReadTargets` — table-driven через `strings.NewReader` (happy / пустые строки /
  пустой ввод). Старые тесты проб не тронуты. Все `t.Errorf` приведены к единому канону
  `Call(input) = got, want want` (порядок got→want, английский). 11 подтестов, `-race` зелёные.

---

## Ошибки и трудности

**Тема прошла легче предыдущих** (со слов, 2026-06-15): «не так уж сложно, потому что
основательно затронули теорию заранее». Контраст с A0.7 интерфейсами, где «зачем» долго
не сходилось. Подтверждает принцип «теория ПЕРЕД заданием».

Конкретные ловушки по ходу:
1. **`scanner.Err() == io.EOF` — НЕВЕРНО.** Думал «Scan кончился → Err вернёт io.EOF».
   `bufio.Scanner` прячет нормальный EOF внутри: `Err()` при штатном конце даёт `nil`, не EOF.
   Из-за этого первая версия уходила в `else` и возвращала `nil` целей на успешном чтении →
   симптом «нужен минимум один домен». Исправил на `if err := sc.Err(); err != nil`.
2. **nil-слайс vs пустой `[]string{}`** — кейс `empty` падал на `reflect.DeepEqual`:
   `want []` и `got []` печатаются одинаково, но `DeepEqual([]string{}, nil) == false`.
   `readTargets("")` возвращает nil (`var out []string` без append). Поправил ожидание
   на `nil`. Это частый Go-собес вопрос — занесено в [[stdlib-cli]] и [[slices-maps]].
3. **Копирование `io.Reader`** в тесте — `s := *strings.NewReader(...)` потом `&s`:
   разыменовать→скопировать→снова взять адрес. Антипаттерн (Reader хранит позицию чтения).
   Исправил на `readTargets(strings.NewReader(tc.input))` напрямую.
4. **Молчаливый `os.Exit(1)`** — выходил без сообщения. Добавил `Fprintf(os.Stderr,...)`
   с причиной. Правило: не глотать ошибку — печатать в stderr перед выходом.
5. **`io.EOF` руками не нужен** при `bufio` — обёртка делает это за тебя; EOF как ошибка
   приходит только при прямом `Read`.

**Концептуальные вопросы, разобранные в теории (вынесены в [[stdlib-cli]]):** зачем `flag.X`
возвращает указатель (регистрация ≠ парсинг, ячейка + `new` → не nil), `*timeout` при
передаче, потоки stdout/stderr и семейство `Fprint`, коды выхода 0/1/2, `errors.New` vs
`fmt.Errorf`/`%w`. `os`/потоки — базовое понимание есть, не автоматизм (пометка для ретеста).

---

## Что бы сделал иначе

- Сразу помнить: `bufio.Scanner.Err()` НЕ возвращает `io.EOF` — `return out, sc.Err()` достаточно.
- Не копировать `io.Reader` в тесте — передавать `strings.NewReader(...)` напрямую.
- Не выходить молча — `Fprintf(os.Stderr, ...)` перед `os.Exit` с самого начала.
- Держать в голове nil vs `[]string{}` при сравнении слайсов в тестах.

---

## Ключевые пакеты и паттерны

- **`flag.FlagSet`** — изолированный набор флагов на подкоманду; `Duration`/`String` → указатели.
- **`os.Open` + `defer Close`** после проверки `err`; **`bufio.Scanner`** построчно.
- **`io.Reader` в сигнатуре** ради тестируемости (`strings.NewReader` вместо файла) — закрепление [[interfaces]].
- **stdout (результат) / stderr (ошибки)** — `Fprintf(os.Stderr, ...)`.
- **Коды выхода 0/1/2** (успех / runtime / usage).
- **table-driven тесты** + `reflect.DeepEqual` (с учётом nil≠empty); единый канон `t.Errorf`.

---

## Связанные темы

[[interfaces]] [[channels]]
