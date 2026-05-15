# Пакет `testing`

## TL;DR

Стандартный пакет для тестов. Файлы `*_test.go` рядом с кодом, функции
`func TestXxx(t *testing.T)` запускаются через `go test`. В C/C++ нужен внешний
фреймворк, в Go всё встроено.

## Где разобрано подробно

- `pkg.go.dev/testing` — справочник API
- `gobyexample.com/testing-and-benchmarking` — минимальные примеры
- Видео Николая Тузова: «Генерация и использование моков в Go / Mockery» (моки — продвинутая тема)

## Структура тест-файла

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

## `*testing.T` — основные методы

| Метод | Когда |
|-------|-------|
| `t.Errorf(format, args)` | провал, но идём дальше — следующие assertions выполнятся |
| `t.Fatalf(format, args)` | провал и сразу `runtime.Goexit()` — для случаев когда дальше тестировать бессмысленно (nil pointer впереди) |
| `t.Run(name, func(*T))` | подтест — отдельная строка в выводе, можно запускать отдельно |
| `t.Helper()` | вызвать в начале helper-функции — ошибка показывает строку вызывающего |
| `t.Parallel()` | запустить параллельно с другими `t.Parallel()` тестами |
| `t.Cleanup(func())` | действие при завершении теста (закрыть файл, восстановить env) |
| `t.Skip(reason)` | пропустить тест условно |
| `t.Logf(format, args)` | лог, показывается только при `-v` или провале |

**`Fatalf` — не паника.** `panic()` нужно recover'ить; `Fatalf` — управляемая остановка
через `runtime.Goexit()`: следующие тесты продолжают работать.

**Сообщение всегда содержит `got` и `want`:**
```go
t.Errorf("got %v, want %v", got, want)   // правильно
t.Errorf("неправильно!")                  // бесполезно
```

## Table-driven — главная идиома

```go
func TestPop(t *testing.T) {
    tests := []struct {
        name      string
        input     []int   // что было в стеке ДО Pop
        wantVal   int
        wantOk    bool
        wantAfter []int   // что стало в стеке ПОСЛЕ Pop
    }{
        {name: "обычный pop",  input: []int{1, 2, 5, 7}, wantVal: 7, wantOk: true, wantAfter: []int{1, 2, 5}},
        {name: "пустой стек",  input: []int{},          wantVal: 0, wantOk: false, wantAfter: []int{}},
        {name: "один элемент", input: []int{42},        wantVal: 42, wantOk: true, wantAfter: []int{}},
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            s := &Stack{items: tc.input}
            gotVal, gotOk := s.Pop()
            if gotVal != tc.wantVal || gotOk != tc.wantOk {
                t.Errorf("Pop() = (%v, %v), want (%v, %v)", gotVal, gotOk, tc.wantVal, tc.wantOk)
            }
            if !reflect.DeepEqual(s.items, tc.wantAfter) {
                t.Errorf("state = %v, want %v", s.items, tc.wantAfter)
            }
        })
    }
}
```

**Правило:** все параметры кейса — поля структуры. Никаких хардкодов:
```go
// плохо
got := topMessages(tc.top, 2)

// хорошо — n в структуре
got := topMessages(tc.top, tc.n)
```

## Сравнение слайсов / мэпов

Через `==` нельзя (compile error для слайсов, identity для указателей у мэпов). Только `reflect.DeepEqual`:
```go
import "reflect"
if !reflect.DeepEqual(got, want) {
    t.Errorf("got %v, want %v", got, want)
}
```

Важно: `map[string]int{}` (пустая) ≠ `map[string]int(nil)` для `DeepEqual`. Если функция
возвращает `make(map[string]int)` — сравнивать с `map[string]int{}`, не с nil.

## Стандарт тестов в этом проекте

Минимум **3 кейса** на каждую функцию:
1. **Happy path** — обычный валидный ввод
2. **Граничный** — пустой, nil, размер 1, max
3. **Злой** — мусор, паника-кейс, граница

Для concurrency-кода — обязательный `go test -race`.

## Запуск

```bash
go test ./...                                    # все тесты
go test -v ./cmd/pointers                        # подробный вывод
go test -run TestPop ./cmd/pointers              # только TestPop
go test -run TestPop/пустой_стек ./cmd/pointers  # подтест (пробел → _)
go test -race ./...                              # детектор гонок (для concurrency)
go test -cover ./...                             # покрытие в процентах
```

Форматирование:
```bash
gofmt -l cmd/pointers/   # показать неправильно отформатированные
gofmt -w cmd/pointers/   # перезаписать с правильным
```

## Личный опыт

### Что я понял не сразу

- **`t.Errorf` vs `t.Fatalf`.** `Errorf` записывает провал и идёт дальше. `Fatalf` —
  останавливает subtest. Использовать `Fatalf` когда следующие проверки зависят от
  предыдущей (проверили `cfg != nil`, иначе `cfg.Host` → паника).
- **`reflect.DeepEqual` обязателен** для слайсов/мэпов. `==` либо не компилируется,
  либо сравнивает по указателю (не по содержимому).
- **Поле теста `input` ≠ `want`.** Если test struct использует одно поле и для входа
  и для ожидания — тест становится тавтологией (положили `want`, проверили `want`).
- **Без `t.Run` теряются имена кейсов.** В выводе будет только «FAIL: TestPop»
  без указания какой кейс упал.

### На чём попадался

- task-03 (`TestPush`): использовал `tc.want` как input — тест проходил всегда,
  потому что клал в стек ровно то что потом проверял (это **до сих пор** в коде, см. `cmd/pointers/pointers_test.go`)
- task-05: написал `int64(0)` и `int64(3)` в литералах test struct — лишний каст, тип выводится
- task-06: `func(s string)` где `s` не используется — должно быть `func(_ string)`
- task-05: проверка `tc.err == errors.New("...")` всегда `false` — каждый вызов `errors.New`
  создаёт новый объект. Сравнивать через `errors.Is`

### Подводные камни

- Файл должен быть `*_test.go` (с подчёркиванием) — иначе `go test` не увидит
- Функция должна начинаться с `Test` + большая буква (`TestFoo`, не `Testfoo`)
- `reflect.DeepEqual(nil, []int{})` → `false` — nil-slice и пустой slice не равны
- Запуск без `-race` для concurrency-кода — гонки не отловятся, но в проде упадут

## Другие виды тестов (для справки — освоим позже)

- `testing.B` — бенчмарки: `BenchmarkXxx(b *testing.B)`, запуск `go test -bench=.`
- `testing.F` — fuzzing: автоматическая генерация входных данных для поиска паник

## Связанные темы

[[slices-maps]] [[errors]] [[goroutines]]
