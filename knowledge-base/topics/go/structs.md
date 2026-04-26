# Структуры и методы в Go

## Что это

Тип-агрегат полей. Аналог `struct` в C — буквально то же самое по памяти и назначению.

## Как работает

```go
// Go
type Task struct {
    ID     int
    Title  string
    IsDone bool
}
```

```c
// C — то же самое
typedef struct {
    int   id;
    char *title;
    int   is_done;
} Task;
```

Методы определяются снаружи структуры через receiver:

```go
// Pointer receiver — метод получает указатель на оригинал, может менять поля
func (t *Task) Complete() {
    t.IsDone = true
}

// Value receiver — метод получает копию, поля менять бессмысленно
func (t Task) String() string {
    return fmt.Sprintf("[...] %s", t.Title)
}
```

Аналог в C:
```c
// pointer receiver:
void task_complete(Task *t) { t->is_done = true; }

// value receiver:
char *task_string(Task t) { /* работает с копией */ }
```

Go добавляет синтаксический сахар: `t.Complete()` вместо `task_complete(&t)`.

## Конструктор

```go
func NewTask(id int, title string) *Task {
    return &Task{
        ID:    id,
        Title: title,
        // IsDone не указываем — нулевое значение bool уже false
    }
}
```

Нулевые значения в Go: `int → 0`, `bool → false`, `string → ""`, указатель → `nil`.

## Интерфейс fmt.Stringer

Если тип реализует метод `String() string` — `fmt.Println` использует его автоматически.
Это неявная реализация интерфейса (тема A0.8):

```go
func (t *Task) String() string {
    sgn := " "
    if t.IsDone { sgn = "x" }
    return fmt.Sprintf("[%s] %d. %s", sgn, t.ID, t.Title)
}

fmt.Println(task) // вызывает task.String() сам
```

## Подводные камни

- **Value receiver у мутирующего метода** — изменения потеряются, Go не предупредит
- **Nil receiver** — вызов метода на nil-указателе вызовет панику если метод обращается к полям
- **Лишние публичные методы** — публичный метод это обязательство для всех, кто использует пакет

## Что ещё нужно разобрать

Модель памяти: `[]Task` vs `[]*Task` vs `[][]Task` vs `*[]Task` — отдельная тема, связана с указателями (A0.3).

## Связанные темы

[[pointers]] [[interfaces]] [[fmt]]
