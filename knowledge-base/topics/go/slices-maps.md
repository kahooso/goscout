# Слайсы и мэпы в Go

## Слайсы

### Что это

Слайс — структура из трёх полей ("заголовок"), хранится на стеке. Данные — в heap.

```
заголовок (стек):          массив (heap):
┌─────────────────────┐    ┌──────────────────┐
│ ptr *T  │ len │ cap │ →  │ [0] [1] [2] ... │
└─────────────────────┘    └──────────────────┘
```

Аналог в C:
```c
int *ptr = malloc(n * sizeof(int));
int len = n;
int cap = n;
```

### Нулевое значение и инициализация

```go
var s []int          // nil slice: {ptr: nil, len: 0, cap: 0}
s = append(s, 1)    // Go аллоцирует массив, возвращает новый заголовок
s = append(s, 2, 3) // append вариативный — несколько элементов сразу
```

Когда `len == cap` и делаешь `append` — Go аллоцирует новый массив ~2x, копирует, возвращает новый `ptr`. Как `realloc` в C.

### Варианты типов

| Тип         | Что хранит в массиве    | Аналог в C                           |
| ----------- | ----------------------- | ------------------------------------ |
| `[]string`  | строки (ptr+len каждая) | нет прямого аналога                  |
| `[]*string` | указатели на строки     | `char *arr[]` или `char **arr`       |
| `*[]string` | указатель на заголовок  | нет аналога (в C управляешь вручную) |

### string — не массив символов

`string` в Go — read-only slice байт: `{ptr *byte, len int}`. Иммутабельна — нельзя изменить символ по индексу.

`[N]T` — статический массив, размер фиксирован на этапе компиляции. Аналог `int arr[5]` в C.

### Pre-allocate когда знаешь размер

```go
// медленнее — несколько реаллокаций
var keys []string
for k := range m { keys = append(keys, k) }

// быстрее — один раз выделяем нужный размер
keys := make([]string, 0, len(m))
for k := range m { keys = append(keys, k) }
```

---

## Мэпы

### Что это

`map[K]V` — указатель на hash-таблицу в heap. Аналог `std::unordered_map<K,V>*` в C++.

```go
var m map[string]int     // nil, указатель ни на что
m = make(map[string]int) // Go аллоцирует hash-таблицу
```

### Нулевое значение

```go
var m map[string]int
v := m["key"]   // ok — вернёт 0 (нулевое значение int)
m["key"] = 1    // ПАНИКА — запись в nil map
```

### Основные операции

```go
counts := make(map[string]int)
counts["hello"]++          // ключа нет → создаётся с нулём, потом +1
counts["hello"] = 5        // явное присвоение
delete(counts, "hello")    // удаление ключа

for word, n := range counts { // перебор всех пар
    fmt.Println(word, n)
}
```

### Порядок не гарантирован

`range` по мэпу каждый раз даёт разный порядок. Чтобы отсортировать — извлечь ключи в слайс:

```go
keys := make([]string, 0, len(counts))
for k := range counts { keys = append(keys, k) }
sort.Slice(keys, func(i, j int) bool {
    return counts[keys[i]] > counts[keys[j]] // по убыванию частоты
})
```

---

## Подводные камни

- `var s []int` — nil slice, но `append` работает корректно
- `var m map[string]int` — nil map, запись вызовет панику
- `sort.Slice` требует строгое упорядочение: `>=` и `<=` в `less` сломают сортировку
- `make(map[string]int, 0)` — лишний `0`, то же что `make(map[string]int)`

## Связанные темы

[[strings]] [[sort]] [[pointers]] [[structs]]
