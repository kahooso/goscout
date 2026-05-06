# Пакет testing

## Что это

Встроенный пакет для написания тестов. Тесты запускаются через `go test`, не входят в бинарник.
В C/C++ тесты требуют внешних фреймворков (Google Test, Unity). В Go — ничего внешнего.

## Как работает

### Структура тест-файла

```go
// logparse_test.go — суффикс _test.go обязателен
package main   // тот же пакет что и тестируемый код

import "testing"

func TestCountByLevel(t *testing.T) {
    // тело теста
}
```

Go компилирует `*_test.go` только для `go test` — отдельный бинарник, не входит в `go build`.
Функция теста: `TestXxx(t *testing.T)` — обязательно с заглавной буквы после `Test`.

### testing.T — это struct, не интерфейс

`*testing.T` — указатель на структуру. Методы через него меняют внутреннее состояние теста
(отмечают провал, пишут лог). Поэтому передаётся именно `*T`, а не `T`.

### Основные методы t

```go
t.Errorf("got %v, want %v", got, want)
// тест помечен провален, выполнение продолжается

t.Fatalf("got %v, want %v", got, want)
// тест помечен провален, функция останавливается через runtime.Goexit()
```

`Fatalf` — **не паника**. `panic()` нужно recover'ить. `Fatalf` — управляемая остановка
через `runtime.Goexit()`: тест завершается чисто, следующие тесты продолжают работать.

`Errorf` — когда проверки независимы. `Fatalf` — когда дальше нет смысла (nil вместо слайса и т.п.).

**Сообщение всегда должно содержать got и want:**
```go
t.Errorf("got %v, want %v", got, tc.want)  // правильно
t.Errorf("не правильно!")                  // бесполезно — при провале ничего не видно
```

### t.Log / t.Logf

```go
t.Logf("промежуточное значение: %v", got)
```

Печатается только если:
- тест **провален** (всегда показывается)
- запущен с флагом **`-v`** (явно попросил подробности)

Удобно для отладки: добавил лог, посмотрел с `-v`, убрал.

### Table-driven pattern — главная идиома

```go
func TestCountByLevel(t *testing.T) {
    tests := []struct {
        name string
        logs []LogEntry
        want map[string]int
    }{
        {
            name: "три уровня",
            logs: []LogEntry{{Level: "INFO"}, {Level: "ERROR"}, {Level: "INFO"}},
            want: map[string]int{"INFO": 2, "ERROR": 1},
        },
        {
            name: "пустой вход",
            logs: []LogEntry{},
            want: map[string]int{},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := countByLevel(tc.logs)
            if !reflect.DeepEqual(got, tc.want) {
                t.Errorf("got %v, want %v", got, tc.want)
            }
        })
    }
}
```

### t.Run — подтесты

`t.Run(name, func(t *testing.T))` запускает подтест с именем и **отдельным** `t`.

Зачем отдельный `t` внутри: если вызвать `t.Fatalf` внутри подтеста — останавливается
только он. Внешний цикл продолжается, остальные кейсы запустятся.
Без `t.Run` — `Fatalf` остановит весь `TestXxx` после первого провала.

Правило: **все параметры кейса должны быть полями структуры**, включая `n`, `timeout` и т.п.

```go
// плохо — n захардкожен, нельзя варьировать по кейсам
got := topMessages(tc.top, 2)

// хорошо — n в структуре
tests := []struct {
    name string
    top  map[string]int
    n    int        // ← поле
    want []string
}{ ... }
got := topMessages(tc.top, tc.n)
```

### Сравнение map — reflect.DeepEqual

Map нельзя сравнить через `==` — ошибка компиляции. Нужен `reflect.DeepEqual`:

```go
import "reflect"

if !reflect.DeepEqual(got, want) {
    t.Errorf("got %v, want %v", got, want)
}
```

Важно: `map[string]int{}` (пустая) и `map[string]int(nil)` — **разные** значения для `DeepEqual`.
Если функция возвращает `make(map[string]int)` — сравнивать с `map[string]int{}`, не с nil.

## Запуск

```bash
go test ./cmd/logparse/                          # тихий (только PASS/FAIL)
go test -v ./cmd/logparse/                       # подробный (все t.Run с именами)
go test -run TestCountByLevel ./cmd/logparse/    # только эта функция
go test -run TestCountByLevel/пустой_вход ./...  # только этот подтест (пробел → _)
go test ./...                                    # все тесты в модуле
```

Форматирование тест-файлов — тем же `gofmt`:
```bash
gofmt -w cmd/logparse/   # перезаписать с правильным форматированием
gofmt -l cmd/logparse/   # только показать какие файлы неправильно отформатированы
```

## Другие типы тестов (для справки — будет позже)

- `testing.B` — бенчмарки: `BenchmarkXxx(b *testing.B)`, запуск через `go test -bench=.`
- `testing.F` — fuzzing: автоматическая генерация входных данных для поиска паник

## Подводные камни

- `t.Errorf("не правильно!")` — при провале ничего не видно. Всегда включать `got` и `want`.
- `*_test.go` с объявленными но неиспользуемыми переменными → ошибка компиляции (`go vet`)
- `map` нельзя сравнить через `==` — нужен `reflect.DeepEqual`
- Параметры конкретного кейса должны быть полями struct — не хардкодить в вызове функции
- `t.Fatalf` — не паника, `runtime.Goexit()`. Другие тесты продолжат работать.

## Связанные темы

[[slices-maps]] [[structs]] [[interfaces]]
