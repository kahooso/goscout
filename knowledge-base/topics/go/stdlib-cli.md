# stdlib для CLI: flag · os · bufio · io.Reader

## TL;DR
Четыре пакета строят CLI: `flag` парсит `--флаги` и позиционные аргументы, `os.Open`
открывает файл (`*os.File`), `bufio.Scanner` режет поток на строки, `io.Reader` —
контракт «любой источник байт», который развязывает функции от диска (→ тестируемость).

## Где разобрано подробно
- `pkg.go.dev/flag`, `pkg.go.dev/os`, `pkg.go.dev/bufio`, `pkg.go.dev/io`
- `gobyexample.com/command-line-flags`, `.../reading-files`, `.../line-filters`
- видео Николая Тузова: «Пакет flag», «Работа с файлами в Go» (искать по названию)

## Минимальный пример
```go
// main: подкоманда → свой FlagSet → источник целей → проба
switch os.Args[1] {                                   // уровень подкоманд
case "dns":
    fs := flag.NewFlagSet("dns", flag.ExitOnError)    // изолированный набор флагов
    timeout := fs.Duration("timeout", 5*time.Second, "таймаут на резолв") // *time.Duration
    wordlist := fs.String("wordlist", "", "файл со списком доменов")      // *string
    fs.Parse(os.Args[2:])                             // [0]=goscout [1]=dns отрезаны

    var targets []string
    if *wordlist != "" {
        f, err := os.Open(*wordlist)                  // (*os.File, error)
        if err != nil { ... }                         // проверка ДО defer
        defer f.Close()
        targets, err = readTargets(f)                 // *os.File → io.Reader
    } else {
        targets = fs.Args()                           // позиционные домены
    }
    runProbe(DNSProbe{}, targets, *timeout)           // *timeout — разыменовать
}

// читалка построчно — принимает io.Reader, НЕ *os.File (ключ к тестам)
func readTargets(r io.Reader) ([]string, error) {
    var out []string
    sc := bufio.NewScanner(r)
    for sc.Scan() {                                   // true пока есть строки
        line := sc.Text()                            // строка без \n (и без \r на Windows-\r\n)
        if line == "" { continue }                   // пропуск пустых
        out = append(out, line)
    }
    return out, sc.Err()                             // ОБЯЗАТЕЛЬНО: EOF vs ошибка
}
```

## Личный опыт

### Что я понял не сразу
- **`bufio.Scanner` сам режет окончания строк, ОБА варианта.** Сплиттер по умолчанию
  `bufio.ScanLines` снимает и `\n` (Unix), и `\r\n` (Windows) — `Text()` всегда отдаёт
  «чистую» строку. Не надо вручную `strings.TrimRight(s, "\r\n")` после `Text()` — это
  лишнее. (Иначе при чтении wordlist, созданного в Блокноте, домен `example.com\r`
  не нашёлся бы в DNS — резолвер бы ругался на невалидное имя.)
- **Зачем `flag.X` возвращает указатель (`*int`, `*time.Duration`).** Регистрация флага и
  его парсинг — два РАЗНЫХ момента. При регистрации `fs.Duration(...)` создаёт ячейку,
  кладёт дефолт, возвращает её **адрес**. При `fs.Parse(args)` записывает распарсенное
  значение **в ту же ячейку**. Вернись `flag` значением (копией дефолта) — `Parse` писал бы
  в свою ячейку, а твоя копия осталась бы дефолтом навсегда. Указатель = ты и `flag`
  смотрите на одну ячейку. Аналогия C: это `scanf("%d", &n)` — функция пишет по адресу,
  только адрес `flag` создаёт сам через `new(T)` (а `new` никогда не даёт nil → указатель
  валиден всегда, ещё до Parse, в нём дефолт — не nil).
- **`*timeout` при передаче дальше.** `timeout` — это `*time.Duration` (адрес). Параметр
  `runProbe(... timeout time.Duration)` — значение. Типы должны совпасть → `*timeout`.
  Указатель в параметре не нужен: функция только читает (отдаёт в `context.WithTimeout`),
  не меняет. Правило: указатель в параметре — только если функция МЕНЯЕТ значение у
  вызывающего. `int64` (= `time.Duration`) копировать дёшево, копируй по значению.
- **`io.Reader` вместо `*os.File` в `readTargets`.** Параметр-интерфейс → на его место
  встаёт что угодно с методом `Read(p []byte)(int,error)`. В бою — `*os.File` (из `os.Open`),
  в тесте — `*strings.NewReader("a\nb\n")` (из памяти, без диска). Функция не видит разницы.
  Это и есть «программирование против интерфейса» из [[interfaces]] (A0.7), закреплённое на
  чтении. `os.Open` живёт в `main`, не внутри `readTargets` — чтобы тест не трогал диск.
- **`bufio` буфер = меньше syscall.** Голый `Read` по байту = syscall на каждый байт
  (дорого, переход в ядро). `bufio.NewScanner` тянет в буфер сразу ~64 КБ одним `Read`,
  потом режет строки в памяти без новых syscall. Аналогия C: свой `char buf[65536]` +
  один `read()`, дальше бегаешь по буферу. Дефолтный лимит строки 64 КБ (меняется
  `sc.Buffer(...)`). К процессору размер не привязан.

### Потоки: stdout vs stderr + семейство Fprint
- Каждая программа имеет 3 потока (как в C): `os.Stdin` (fd 0), `os.Stdout` (fd 1),
  `os.Stderr` (fd 2). В Go это `*os.File` → они же `io.Writer` → годятся в `Fprintf`.
- **`F`-префикс = «куда писать»**, НЕ асинхронность. `Printf(...)` ≡ `Fprintf(os.Stdout, ...)`.
  Таблица: без формата / `ln` / `f` × stdout (`Print`/`Println`/`Printf`) или куда укажешь
  (`Fprint`/`Fprintln`/`Fprintf`).
- **Зачем 2 потока:** их разделяют при запуске. `> file` уводит stdout, `2> file` — stderr,
  `| grep` берёт ТОЛЬКО stdout. Два уровня: программа выбирает ЯЩИК (`Fprintf(os.Stderr,...)`),
  shell выбирает МАРШРУТ ящика. Правило: **результат → stdout, ошибки/диагностика → stderr.**
  Ошибка в stdout испортила бы пайп/файл с результатом.
- `fmt.Errorf` СОЗДАЁТ error (не печатает!) ≠ `Fprintf(os.Stderr,...)` ПЕЧАТАЕТ.
  `errors.New(s)` — фиксированный текст; `fmt.Errorf(fmt, ...)` — с подстановкой и `%w` (wrap).

### Коды выхода (exit codes)
- Конвенция (Unix, не жёсткий стандарт): **0 успех / 1 runtime-ошибка (файл, сеть — среда
  подвела) / 2 usage-ошибка (кривой ввод пользователя)**. Технически любое 0–255, но держись 0/1/2.
- Баг в самом коде → не `os.Exit`, а `panic` + стектрейс (рантайм сам, код 2).
- BSD `sysexits.h` (64/66/77...) — попытка расширить конвенцию, почти никто не использует.
- Код выхода = API программы для скриптов/CI (`cmd && echo OK` смотрит на 0 vs не-0).
  Человек читает stderr, скрипт читает код выхода.

### nil-слайс vs пустой слайс (ловушка тестов)
- `var out []string` (без append) → **nil-слайс**. `[]string{}` → **пустой не-nil** слайс.
- `%v` печатает оба как `[]` (выглядят одинаково!), НО `reflect.DeepEqual([]string{}, nil)`
  → **false**. Поэтому `TestReadTargets/empty` падал: ждал `[]string{}`, функция вернула `nil`.
- Фикс: ждать `nil` в кейсе пустого ввода (отражает реальность функции). Либо сравнивать по
  `len` (тогда nil и `{}` равны). Частый Go-собес вопрос. См. [[slices-maps]].

### На чём попадался
- `task-09`: **`scanner.Err() == io.EOF`** — НЕТ. `bufio.Scanner` прячет нормальный EOF;
  `Err()` при штатном конце = `nil`. Уход в `else` → возврат `nil` целей на успешном чтении.
  Достаточно `return out, sc.Err()`.
- `task-09`: **копирование `io.Reader`** — `s := *strings.NewReader(...)` + `&s`. Антипаттерн:
  Reader хранит позицию чтения. Передавать `strings.NewReader(...)` напрямую.
- `task-09`: **молчаливый `os.Exit(1)`** без сообщения = проглоченная ошибка. Перед выходом —
  `Fprintf(os.Stderr, "...: %v\n", err)`.
- `task-09`: **«сколько тире — зависит от длины имени флага»** — НЕТ. Это Unix-традиция
  (`-v` короткий, `--verbose` длинный). Пакет `flag` в Go **не различает** `-` и `--`:
  `--timeout` ≡ `-timeout`. Количество букв ни при чём.
- `task-09`: **«копировать int64 не дёшево»** — наоборот, ДЁШЕВО (8 байт). Это аргумент ЗА
  передачу по значению. Дорого копировать большие структуры, не примитивы.
- `task-09`: забыть `scanner.Err()` после цикла — это не «утечка данных», а **сокрытие
  ошибки**: вернёшь неполный результат молча, как будто всё прочиталось. (`Scan()` даёт
  `false` и на EOF, и на ошибке — различает только `Err()`.)
- `task-09`: при `os.Open`-ошибке `file == nil` → `defer file.Close()` ставить ТОЛЬКО
  после проверки `err` (иначе `nil.Close()` → паника). nil-файл чистить не нужно и нельзя:
  память не выделена, дескриптора нет.

### Подводные камни
- **`os.Exit` НЕ выполняет `defer`** — рубит процесс сразу, теряет `Close()`. Звать только
  в самом верху `main` (или `log.Fatal`). В обычных функциях/библиотеках — возвращать `error`,
  не `os.Exit`. `flag.ExitOnError` делает `os.Exit` внутри — приемлемо, т.к. в начале до
  открытия файлов (терять нечего).
- **Флаги ДО позиционных.** `flag` стопорится на первом не-флаге: `dns example.com --timeout 3s`
  → `--timeout` уже не распознан. Сначала флаги, потом позиционные.
- **Сколько токенов съедает флаг:** не-булев — ровно ОДИН следующий (`--timeout` → `3s`);
  булев (`fs.Bool`) — НОЛЬ (`-v` сам по себе = true, следующий токен не съедает).
- **Сигнатура `Read` целиком:** `Read([]byte) int` (без `error`) НЕ реализует `io.Reader`.
  Нужно `Read([]byte) (int, error)` байт-в-байт.

### Темы на потом (помечено в этой сессии)
- **const в Go vs C/C++.** В Go нет `const` для указателей/значений — только для
  compile-time констант (числа/строки/bool/iota). Аналога `const int*` / `int* const` нет.
  Защита от изменения в Go: передача по значению (копия), неэкспортируемые поля,
  интерфейсы только на чтение. Вернуться если станет интересно.
- **Свой тип флага через `flag.Value`.** Готовые конструкторы есть только для типов, что
  `flag` умеет парсить (String/Int/Bool/Float/Duration). Нет «флага размера» (`10MB`), т.к.
  в stdlib нет `ParseBytes`. Любой свой тип → `fs.Var(...)` + реализация интерфейса
  `flag.Value` (`String()` + `Set(string) error`). Не лезли вглубь.

## Ключевые термины (English)
- `flag` / `FlagSet` — package for parsing CLI args; `FlagSet` = isolated set of flags (per subcommand)
- `positional argument` — unnamed arg, recognized by position (vs named `--flag`)
- `os.Open` → `*os.File` — open file for reading; `*os.File` is a file descriptor (like C `FILE*`)
- `file descriptor` — kernel handle to an open file/socket/pipe
- `io.Reader` — interface `Read(p []byte) (n int, err error)`; "any source of bytes"
- `io.EOF` — sentinel error meaning the stream ended (not a real failure)
- `bufio.Scanner` — buffered line/token reader over an `io.Reader`; `Scan()`/`Text()`/`Err()`
- `buffer` — memory block read in large chunks to avoid per-byte syscalls
- `syscall` — call into the OS kernel (e.g. each raw `read`); expensive, batched via buffering
- `time.Duration` — `int64` of nanoseconds; `flag.Duration` parses `"3s"` via `time.ParseDuration`
- `dereference` — `*p`, read the value a pointer points to

## Связанные темы
[[interfaces]] [[pointers]] [[errors]] [[context]] [[go-tooling]]
