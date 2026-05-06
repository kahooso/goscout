# Пакет testing

## Что это

Встроенный пакет для написания тестов. Тесты запускаются через `go test`, не входят в бинарник.

## Как работает

### Структура тест-файла

```go
// main_test.go — суффикс _test.go обязателен
package main   // тот же пакет что и тестируемый код

import "testing"

func TestCountByLevel(t *testing.T) {
    // тело теста
}
```

Go компилирует `*_test.go` только для `go test`, не для `go build`.
Функция теста: `TestXxx(t *testing.T)` — обязательно с заглавной буквы после `Test`.

### Основные методы t

```go
t.Errorf("got %v, want %v", got, want)
// тест провален, выполнение продолжается

t.Fatalf("got %v, want %v", got, want)
// тест провален, сразу останавливается
```

`Errorf` — когда проверки независимы. `Fatalf` — когда дальше нет смысла (nil вместо слайса и т.п.).

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

`t.Run(name, func)` — подтест с именем. Можно запустить один: `go test -run TestCountByLevel/пустой_вход`.

### Сравнение map — reflect.DeepEqual

Map нельзя сравнить через `==` — ошибка компиляции. Два варианта:

```go
import "reflect"

if !reflect.DeepEqual(got, want) {
    t.Errorf("got %v, want %v", got, want)
}
```

Или вручную — перебрать ключи и сравнить значения + проверить `len`.

## Запуск

```bash
go test ./cmd/logparse/           # тихий (только PASS/FAIL)
go test -v ./cmd/logparse/        # подробный (все t.Run с именами)
go test -run TestCountByLevel ./cmd/logparse/  # один тест
go test ./...                     # все тесты в модуле
```

## Подводные камни

- Файл `*_test.go` с объявленными но неиспользуемыми переменными → ошибка компиляции (`go vet`)
- `map` нельзя сравнить через `==` — нужен `reflect.DeepEqual` или ручное сравнение
- Table-driven удобнее отдельных функций: все кейсы видны вместе, легче добавить новый

## Связанные темы

[[slices-maps]] [[structs]] [[interfaces]]
