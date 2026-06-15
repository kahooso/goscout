---
id: task-08
block: A0.7
topic: интерфейсы (Probe — общий контракт для проб сканера)
status: completed
date_started: 2026-06-09
date_completed: 2026-06-10
---

# Task-08: интерфейс Probe в goscout

## Задание

Ввести интерфейс `Probe` (`Run(ctx, target) Result` + `Name()`), обернуть DNS-резолвер
в `DNSProbe` (первая реализация), сделать заглушку `PortProbe`, прогнать обе через
общий раннер `runProbe(p Probe, targets)`. Рефакторинг A0.6-кода под интерфейс.

**Файл:** `cmd/goscout/main.go`. **Тесты:** `cmd/goscout/main_test.go`, `-race`.

---

## Моё решение

- `Result{Probe, Target, Output []string, Error}` — общий результат любой пробы
  (заменил `dnsResult`).
- `Probe interface { Run(ctx, target) Result; Name() string }`.
- `DNSProbe struct{}` (value receiver): `Run` = `net.DefaultResolver.LookupHost`, возвращает
  `Result`. `Name()` → "dns". (Старый `resolveDomain` стал телом `Run`, но `return` вместо `ch<-`.)
- `PortProbe struct{}` — заглушка: `Run` возвращает `Result{Error: errors.New("ports: not
  implemented yet")}`. Реальный scan — A0.9.
- `runProbe(p Probe, targets)` — ОДИН раннер: A0.6-обвязка (ctx + buffered-канал + горутины +
  select) переехала сюда. Горутина: `ch <- p.Run(ctx, t)` — раннер зовёт Run и сам шлёт в канал.
  Валидация/имя через `p.Name()`.
- `main`: `runProbe(DNSProbe{}, args)` / `runProbe(PortProbe{}, args)` — выбор пробы.
- Тесты: `TestProbeName` (таблица `[]Probe` — полиморфно Name всех проб), `TestDNSProbeRun`
  (happy/invalid), `TestDNSProbeRunTimeout` (истёкший ctx), `TestPortProbeRun` (заглушка → Error).
  `gofmt`/`vet`/`build`/`-race` зелёные. Ручной прогон dns+ports ОК.

---

## Ошибки и трудности

**Главная трудность — концептуальная, не синтаксис: «зачем интерфейс вообще».**
Долго не сходилось: «структуры пустые и одинаковые, зачем 2 типа, почему не звать функции
напрямую». Разобрали (см. [[interfaces]]): типы разные не данными, а методом `Run` (разное
поведение); интерфейс отделяет общий механизм (раннер) от конкретной операции (проба) —
новая проба = новый тип + строка в switch, обвязку не трогаешь. Аналогия: `file_operations`
в ядре Linux. Пустота структур временная — появятся настройки (ports[], timeout).

Конкретные баги по ходу ревью:
1. **`ch` в сигнатуре `Run`** (дважды!) — сначала `DNSProbe.Run(..., ch chan Result)`, потом
   `PortProbe.Run(..., ch chan Result)`. Лишний параметр → не совпадает с интерфейсом → тип
   не реализует Probe. Канал живёт в РАННЕРЕ, не в Run. Run возвращает Result, горутина шлёт.
2. **`res.Output[]`** — невалидный синтаксис (пустые скобки). → `res.Output`.
3. **`HttpProbe` лишний** — добавил третий тип с пустым телом Run (не компилится) + неверный
   вызов `runHTTP(HttpProbe{}, ...)`. Убрал целиком, http не трогаем до post-A0.
4. **`Fprintln` с `%s`** — Fprintln не форматирует, `%s` печатается буквально. → `Fprintf`.
5. **мёртвая `runPorts`** — заменена `runProbe(PortProbe{})`, удалил.
6. **тесты ждали старый `dnsResult`/`resolveDomain`** — переписал под `DNSProbe{}.Run`/`Result`.
7. косметика: `DnsProbe`→`DNSProbe` (аббревиатуры капсом), `"port"`→`"ports"`, убрал `Output: nil`.

**Самое сложное (со слов, 2026-06-10):** не синтаксис, а **рефакторинг всего под
интерфейсы** — переписать рабочий код, понять как структура НЕЯВНО имплементирует
интерфейс, увязать это с тем, что классического ООП в Go нет. Осталось не до конца:
**как интерфейсы устроены в ПАМЯТИ** (itable/data-ptr) и «логически работают» — понято
частично, закрывать постепенно на след. темах. **Тесты пишутся тяжело** — пока не автоматизм.

**Понимание, которое донеслось не сразу (но донеслось):**
- value vs pointer receiver В КОНТЕКСТЕ интерфейсов: value receiver → и `T`, и `*T`
  реализуют; pointer receiver → только `*T` (значение в `[]Probe` не пройдёт).
- `[]Probe` — слайс ИНТЕРФЕЙСОВ (не `*Probe`!). Указатель на данные уже встроен в
  интерфейс-значение → `[]*Probe` почти никогда не нужен. value/pointer определяется тем,
  что КЛАДЁШЬ в `[]Probe` (`DNSProbe{}` vs `&DNSProbe{}`), а не типом слайса.
- typed-nil ловушка (приготовлена, в коде не выстрелила — `err` уже типа `error`).

---

## Что бы сделал иначе

- Сразу держать в голове «Run чистая (target→Result), канал в раннере» — не пихать ch в Run.
- Не добавлять `HttpProbe` вперёд плана.
- Таблицу из ОДНОГО кейса не делать (TestDNSProbeRunTimeout / TestPortProbeRun) — для одного
  кейса проще прямое тело без `[]struct{}`+цикла. Имя кейса заглушки "happy" → неправда.

---

## Ключевые пакеты и паттерны

- **Интерфейс как контракт** + неявная реализация (структура с методами нужных сигнатур).
- **Полиморфный раннер** `runProbe(p Probe, ...)` — один механизм на все реализации.
- value receiver для проб без состояния; `var _ Probe = DNSProbe{}` — compile-time страховка.
- `[]Probe` в тесте — полиморфная проверка `Name()` таблицей.
- A0.6 обвязка (каналы+context+select) переиспользована внутри раннера.

---

## Связанные темы

[[interfaces]] [[channels]] [[context]] [[pointers]] [[go-tooling]]
