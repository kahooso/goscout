# Указатели в Go

## Что это

Указатель хранит адрес переменной в памяти. Синтаксис как в C, но без арифметики и ручного free.

## Как работает

### Базовый синтаксис

```go
x := 42
p := &x   // p имеет тип *int — берём адрес x
*p = 100  // разыменование — меняем x через указатель
fmt.Println(x) // 100
```

`&` — взять адрес (address-of), как в C.
`*` — разыменовать (dereference), как в C.

### Отличия от C

| | C | Go |
|---|---|---|
| Арифметика указателей | `p++` — можно | нельзя |
| Освобождение памяти | `free(p)` — вручную | GC — автоматически |
| Возврат адреса локальной переменной | UB (стек умирает) | безопасно (escape analysis) |
| Нулевой указатель | `NULL` | `nil` |

### Escape analysis

В C нельзя вернуть адрес локальной переменной — она умрёт вместе со стек-фреймом:
```c
int* bad() { int x = 42; return &x; }  // UB
```

В Go компилятор сам решает где хранить переменную. Если адрес "убегает" наружу —
переменная кладётся в heap автоматически:
```go
func newInt(v int) *int {
    x := v
    return &x  // безопасно — Go положит x в heap
}
```

### Value receiver vs pointer receiver

**Value receiver** — метод получает копию структуры. Изменения не видны снаружи:
```go
func (s Stack) Push(v int) {
    s.items = append(s.items, v)  // меняем копию — оригинал не меняется
}
```

**Pointer receiver** — метод получает указатель. Изменения видны снаружи:
```go
func (s *Stack) Push(v int) {
    s.items = append(s.items, v)  // меняем оригинал через указатель
}
```

Go автоматически берёт адрес для addressable переменных:
```go
s := Stack{}
s.Push(42)  // Go неявно делает (&s).Push(42)
```

### Когда pointer receiver

| Ситуация | Receiver |
|----------|----------|
| Метод изменяет поля структуры | `*T` |
| Структура большая (копировать дорого) | `*T` |
| Только читает, структура маленькая | `T` |
| Хоть один метод типа — pointer | все `*T` |

### Три способа создать структуру

```go
var s Stack      // тип Stack (value), s.items == nil
s := Stack{}     // то же, struct literal, s.items == nil
s := NewStack()  // тип *Stack (указатель), через конструктор
```

Конструктор возвращает `*T` — явный сигнал что объект изменяемый:
```go
func NewStack() *Stack {
    return &Stack{items: make([]int, 0)}
}
```

### nil-указатель

```go
var p *int  // p == nil
*p = 42     // паника: nil pointer dereference
```

Разыменование nil → паника (в C — segfault/UB).
В методах: вызов на nil-указателе не паникует до первого обращения к полю.

### Рекурсивные структуры — только через указатель

```go
type Node struct {
    val  int
    next Node   // ошибка компиляции — бесконечный размер
}

type Node struct {
    val  int
    next *Node  // ок — указатель всегда 8 байт
}
```

## Пример — изменяемый стек

```go
type Stack struct {
    items []int
}

func NewStack() *Stack {
    return &Stack{items: make([]int, 0)}
}

func (s *Stack) Push(v int) {
    if s == nil {
        return
    }
    s.items = append(s.items, v)
}

func (s *Stack) Pop() (int, bool) {
    if s == nil || len(s.items) == 0 {
        return 0, false
    }
    last := s.items[len(s.items)-1]      // сохранить ДО обрезки
    s.items = s.items[:len(s.items)-1]   // уменьшить len
    return last, true
}

func (s *Stack) Peek() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    return s.items[len(s.items)-1], true
}
```

## Подводные камни

- Nil-check с `s = NewStack()` внутри метода — не работает. Переприсвоение локальной копии указателя не меняет оригинал у вызывающего кода.
- `Pop`: сохранять элемент ДО обрезки слайса. Иначе вернёшь не тот элемент.
- Смешивать value и pointer receivers на одном типе — плохо. Если хоть один `*T` — делай все `*T`.
- `t.Errorf("got %v", v == 0 && ok)` — условие не проверяет конкретные значения. Всегда сравнивай с `tc.wantVal` и `tc.wantOk` напрямую.

## Связанные темы

[[structs]] [[slices-maps]] [[testing]]
