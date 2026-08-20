---
id: task-05
block: A0.4
topic: Обработка ошибок — errors, fmt.Errorf, errors.Is/As
status: completed
date_started: 2026-05-10
date_completed: 2026-05-10
---

# Task-05: Парсер конфига с типизированными ошибками

## Задание

**Файл:** `sections/go/errors-demo/errors.go`
**Тесты:** `sections/go/errors-demo/errors_test.go`

Написать `parseConfig(input string) (*Config, error)` — парсит строку вида
`host=localhost port=8080 timeout=30` в структуру `Config`.

Реализовать:
1. Sentinel error `ErrMissingField` — если обязательное поле отсутствует
2. Кастомный тип `ParseError{Field, Value}` — если значение не конвертируется в int
3. Оборачивать ошибки через `fmt.Errorf("parseConfig: %w", err)`
4. В main показать три кейса с `errors.Is` и `errors.As`

---

## Моё решение

Разбил строку через `strings.Fields` (по пробелам), каждую часть через
`strings.SplitN(field, "=", 2)` на ключ и значение. Switch по ключу заполняет Config.

После цикла проверка `c.Host == "" || c.Port == 0 || c.Timeout == 0` → `ErrMissingField`.

В main: три вызова — happy path, отсутствующее поле (`errors.Is`), неверное значение (`errors.As`).

---

## Ошибки и трудности

- Сложнее всего дались `errors.Is` и `errors.As` — особенно `errors.As` с двойным указателем `&pe`.
- В тестах: перевёрнутое условие `reflect.DeepEqual` вместо `!reflect.DeepEqual` — тест падал на правильном результате.
- Поле `Field` в `ParseError` было с заглавной буквы (`"Port"`) — не совпадало с ключами из input (`"port"`).
- Тест структура с `err error` — не очевидно как использовать `errors.Is`/`errors.As` через поле структуры; решение: type assertion `tc.err.(*ParseError)` для определения ветки проверки.
- Тема в целом далась тяжело — механика оборачивания и разворачивания ошибок пока не интуитивна.

---

## Что бы сделал иначе

Сразу называть поля `ParseError` в том же регистре что ключи в input.
В тестах — сначала написать проверку условия на бумаге, потом в код.

---

## Ключевые пакеты и паттерны

- `errors.New("text")` — sentinel error; каждый вызов создаёт уникальный объект
- `fmt.Errorf("ctx: %w", err)` — оборачивание с контекстом; `%w` сохраняет оригинал
- `errors.Is(err, target)` — ищет target в цепочке обёрток по значению
- `errors.As(err, &target)` — ищет тип в цепочке, кладёт в target; нужен `**T`
- `strings.Fields(s)` — разбивка по пробелам (умнее чем `Split(s, " ")`)
- `strings.SplitN(s, sep, 2)` — разбивка не более чем на 2 части
- `strconv.Atoi(s)` — строка → int, возвращает `(int, error)`
- Type assertion `v, ok := x.(T)` — проверить конкретный тип интерфейса

---

## Связанные темы

[[errors]] [[testing]] [[structs]]
