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

## Личный опыт

### task-13: когда sentinel НЕ нужен

Критерий один: **sentinel заводят под вызывающего, а не под себя.** Он нужен, только если
кто-то будет ветвиться по этой конкретной причине через `errors.Is`. Нет такого
вызывающего — нет sentinel, хватит `errors.New` / `fmt.Errorf` на месте.

В outpost семь разных причин отказа, но все ведут в одно место (напечатать в stderr и выйти).
Ветвить нечего → sentinel не заводили ни одного. Смысл несёт текст, решение несёт поле `code`
в типе. Семь разных смыслов не требуют семи разных значений.

Реальные примеры, где `errors.Is` оправдан: `fs.ErrNotExist` (нет файла — создадим),
`context.DeadlineExceeded` (таймаут — повторим), `io.EOF` (конец данных — не ошибка).
Во всех трёх ошибка — развилка в логике, а не конец программы.

### task-13: Is против As — по значению против по типу

- **`Is` — сравнение по значению**: «это ошибка *именно та самая*?» Нужен sentinel,
  с которым сравнивать. Две разные `errors.New` с одинаковым текстом **не равны**.
- **`As` — по типу + извлечение**: «есть ли в цепочке ошибка *такого сорта*, и если да —
  дай её мне, достану поле».

Нужно достать данные из ошибки → нужен `As`. Нужно узнать причину → `Is`.

### task-13: почему `As` принимает адрес и терпит nil

```go
var ee *exitError        // значение nil — не используется никак
errors.As(err, &ee)      // берётся только ТИП *exitError; ответ записывается в ee
```

`As` переменную **не читает — пишет** в неё. От переменной берётся только её тип как
образец поиска, поэтому `nil` там нормален. Адрес нужен потому, что Go передаёт аргументы
по значению: получив `ee`, функция записала бы найденное в копию. Тот же механизм, что
у `json.Unmarshal(data, &v)`. Дженериков в `errors` нет — API родом из Go 1.13.

Идёт `As` по цепочке **не через `%w`**, а через метод `Unwrap()`; `%w` — лишь один из
способов создать звено с этим методом. Сначала проверяется сама `err`, и только потом
вызывается `Unwrap()` — свою обёртку верхнего уровня `As` находит на нулевом шаге.

### task-13: три инструмента создания ошибки — прогон, не рассуждение

```
errors.New("текст")            → *errors.errorString   лист, без вложений
fmt.Errorf("...: %v", err)     → *errors.errorString   вложенная вклеена КАК ТЕКСТ, потеряна как значение
fmt.Errorf("...: %w", err)     → *fmt.wrapError        вложенная сохранена, есть Unwrap()
```

Прогон на реальной ошибке от `os.Open`:

```
%w текст:  dns: cannot open wordlist: open no-such-file.txt: The system cannot find...
%v текст:  dns: cannot open wordlist: open no-such-file.txt: The system cannot find...
                                    ↑ побайтово одинаковые, на глаз не отличить

%w  -> errors.Is(err, fs.ErrNotExist) = true
%v  -> errors.Is(err, fs.ErrNotExist) = false
```

Разница только в типе. `%v` — тупик цепочки, и это видно не в тексте, а в поведении
`errors.Is`/`errors.As`.

### task-13: `Error()` делегирует ВЛОЖЕННОЙ ошибке

```go
func (e *exitError) Error() string { return e.err.Error() }   // метод ПОЛЯ, другого объекта
func (e *exitError) Unwrap() error { return e.err }           // просто отдаёт поле, вызова нет
```

`Unwrap()` ничего не распаковывает и никого не вызывает. Рекурсия возможна только при
вызове метода **самой обёртки**:

```go
func (e *bad) Error() string { return e.Error() }   // → fatal error: stack overflow
```

`go vet` на этом молчит (exit=0), падает только в рантайме. Более коварная форма того же:
`fmt.Sprintf("exit %d: %v", e.code, e)` — `fmt` при `%v` сам вызовет `Error()`, снова
бесконечная рекурсия. Вот её `vet` как раз ловит (анализатор `printf`). Передавать `e.err`, не `e`.

### task-13: тип ошибки с полезной нагрузкой — образцы в stdlib

Ошибка, везущая дополнительные данные, — принятый в Go дизайн, не самодеятельность:

```go
// $GOROOT/src/io/fs/fs.go:251 — весь тип целиком, 9 строк
type PathError struct { Op, Path string; Err error }
func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }
func (e *PathError) Unwrap() error { return e.Err }
```

`os/exec.ExitError` встраивает `*os.ProcessState` с методом `ExitCode() int` — ровно та же
задача «ошибка + код выхода». Деталь: `ExitCode()` возвращает `-1`, если процесс не стартовал —
сигнальное значение вместо второго возврата.

Контракт `PathError` **требует** непустого `Err` (`Error()` зовёт `e.Err.Error()` без проверки
на nil). Если у части путей вложенной ошибки нет — либо всегда класть `errors.New(...)`,
либо учить `Error()` жить с nil.

### task-13: каждый уровень добавляет СВОЙ контекст

Двойной префикс (`ports: ports: no targets provided`) — симптом размытой ответственности:
низ добавил имя команды, верх добавил ещё раз.

| Уровень | Что знает | Что пишет |
|---|---|---|
| `parsePorts` | конкретное значение | `port 99999 out of range (1-65535)` |
| `runProbe` | что целей нет | `no targets provided` — **без** имени команды |
| `run` | какая команда шла | `ports: %w` |
| `main` | что вывод финальный | перевод строки, поток stderr |

Проверка при написании ошибки: **знает ли эта функция то, что пишет, из своих аргументов?**
`runProbe` получает `Probe` и достать имя может — но это знание принадлежит уровню,
который выбирал команду по `args[1]`.

### Конвенция ST1005 — текст ошибки

Со строчной буквы, без точки, **без `\n`**. Причина: ошибки оборачивают, получается
`outer: inner: leaf` — заглавная буква и точка окажутся в середине, `\n` разорвёт строку.
Перевод строки добавляет тот, кто печатает. Ловилось дважды за одну сессию (task-13).

## Ключевые термины (English)

| Термин | Значение |
|---|---|
| **sentinel error** | ошибка-значение уровня пакета (`ErrNotFound`), с которой сравнивают через `errors.Is` |
| **error wrapping** | оборачивание ошибки с сохранением оригинала — `fmt.Errorf` с `%w` |
| **error chain** | цепочка вложенных ошибок, по которой идут `errors.Is` / `errors.As` |
| **to unwrap** | развернуть обёртку, достать вложенную ошибку; метод `Unwrap() error` |
| **leaf error** | ошибка в самом низу цепочки, без вложений |
| **custom error type** | свой тип с методом `Error()`, несущий дополнительные поля |
| **payload** | полезная нагрузка, которую ошибка везёт помимо текста (код выхода, имя поля) |
| **pointer receiver** | получатель-указатель (`func (e *T) Error()`) — тогда интерфейс реализует `*T` |
| **exit code** | код выхода процесса: 0 успех, 1 runtime failure, 2 usage error |
| **ST1005** | правило staticcheck: текст ошибки со строчной буквы, без точки и `\n` |

## Связанные темы

[[testing]] [[structs]] [[interfaces]] [[stdlib-cli]]
