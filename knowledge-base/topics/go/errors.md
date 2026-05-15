# Обработка ошибок в Go

## Что это

Ошибки в Go — обычные возвращаемые значения, не исключения. Нет try/catch.
`error` — встроенный интерфейс с одним методом: `Error() string`.

## Как работает

### error — это интерфейс

```go
type error interface {
    Error() string
}
```

Любой тип с методом `Error() string` автоматически реализует `error`.
Нет явного `implements` как в Java — только наличие метода.

### bool vs error

**`bool`** — когда оба исхода ожидаемы в нормальной работе:
```go
val, ok := m["ключ"]   // ключа может не быть — это норма
val, ok := s.Pop()     // стек может быть пустым — это норма
```

**`error`** — когда что-то пошло не так и нужно объяснить причину:
```go
f, err := os.Open("file.txt")      // файл недоступен — это ошибка
n, err := strconv.Atoi("abc")      // не число — это ошибка
```

Правило: если вызывающий код хочет знать **почему** не получилось — нужен `error`.

### Создание ошибок

```go
// Простая статическая ошибка:
err := errors.New("not found")

// С форматированием:
err := fmt.Errorf("user %d not found", id)

// С оборачиванием (%w сохраняет оригинал):
err := fmt.Errorf("parseConfig: %w", originalErr)
```

### Sentinel errors

Именованные переменные уровня пакета — аналог errno-кодов в C (`ENOENT`, `EACCES`):

```go
var ErrNotFound = errors.New("not found")
var ErrMissingField = errors.New("missing field")
```

Каждый вызов `errors.New` создаёт уникальный объект — даже с одинаковым текстом.
Поэтому объявляют один раз, не внутри функций.

### Оборачивание ошибок — %w

```go
func readConfig(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("readConfig %s: %w", path, err)
    }
    return nil
}
```

Цепочка в сообщении: `"readConfig config.json: open config.json: no such file"`.
`%w` сохраняет оригинальный объект — `errors.Is`/`errors.As` могут его найти.

### errors.Is

Ищет конкретное **значение** в цепочке обёрток:

```go
var ErrNotFound = errors.New("not found")

err := fmt.Errorf("layer2: %w", fmt.Errorf("layer1: %w", ErrNotFound))
errors.Is(err, ErrNotFound)  // true — нашёл на любом уровне
err == ErrNotFound           // false — прямое сравнение не работает через обёртки
```

### errors.As

Ищет конкретный **тип** в цепочке и извлекает объект:

```go
type ValidationError struct{ Field string }
func (e *ValidationError) Error() string { return "invalid: " + e.Field }

err := fmt.Errorf("handler: %w", &ValidationError{Field: "email"})

var ve *ValidationError
if errors.As(err, &ve) {    // передаём **T, не *T
    fmt.Println(ve.Field)   // "email" — достали объект со всеми полями
}
```

`errors.As` принимает `**T` — указатель на указатель — чтобы записать найденный объект.

### Кастомный тип ошибки

```go
type ParseError struct {
    Field string
    Value string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("invalid value for %s: %q", e.Field, e.Value)
}
```

Возвращается как `error`, извлекается через `errors.As`.

### Главный паттерн

```go
result, err := someFunc()
if err != nil {
    return fmt.Errorf("myFunc: %w", err)  // добавить контекст и передать выше
}
// дальше используем result
```

### fmt — где что выводить

| Функция | Куда |
|---------|------|
| `fmt.Printf` | stdout (терминал) |
| `fmt.Sprintf` | возвращает строку |
| `fmt.Fprintf` | в любой `io.Writer` (файл, сеть, буфер) |
| `fmt.Errorf` | возвращает `error` |

## Пример из задания

```go
var ErrMissingField = errors.New("missing field")

type ParseError struct {
    Field string
    Value string
}
func (e *ParseError) Error() string {
    return fmt.Sprintf("invalid value for %s: %q", e.Field, e.Value)
}

func parseConfig(input string) (*Config, error) {
    // ... парсинг ...
    if c.Host == "" || c.Port == 0 || c.Timeout == 0 {
        return nil, fmt.Errorf("parseConfig: %w", ErrMissingField)
    }
    // при ошибке Atoi:
    return nil, fmt.Errorf("parseConfig: %w", &ParseError{Field: "port", Value: value})
}

// В вызывающем коде:
_, err := parseConfig("host=localhost timeout=30")
if errors.Is(err, ErrMissingField) { ... }   // sentinel

_, err = parseConfig("host=localhost port=abc timeout=30")
var pe *ParseError
if errors.As(err, &pe) { fmt.Println(pe.Field) }  // кастомный тип
```

## Подводные камни

- `errors.Is` не работает без `%w` — если завернуть через `%v`, оригинал теряется
- `errors.As` требует `**T` — передавай `&pe`, не `pe`
- Sentinel errors объявлять на уровне пакета — не внутри функций, иначе каждый вызов создаёт новый объект
- `err.Error() == "text"` — хрупко, сломается при любом изменении текста. Используй `errors.Is`

## Связанные темы

[[testing]] [[structs]] [[interfaces]]
