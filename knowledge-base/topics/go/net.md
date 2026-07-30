# net — TCP-сокеты в Go (для goscout ports)

## TL;DR
`net` — низкоуровневый пакет: сырые TCP/UDP-соединения (аналог POSIX `socket()`/`connect()`).
Сканер портов = `net.DialTimeout` по каждому порту + разбор `(conn, err)` на open/closed/filtered.

## Где разобрано подробно
- `pkg.go.dev/net` — сигнатуры и контракты (открывать его, не выдумывать функции)
- `gobyexample.com` → «Worker Pools» (паттерн пула воркеров, пригодится на стадии 2)
- Фундамент под этим пакетом — [[tcp-ip]] (что такое порт, handshake, три исхода скана)

## Минимальный пример
```go
// последовательный скан списка портов на одном хосте
func (p PortProbe) Run(ctx context.Context, target string) Result {
    if len(p.Ports) == 0 {
        return Result{Probe: "ports", Target: target, Error: errors.New("no ports to scan")}
    }
    var open []string
    for _, port := range p.Ports {
        addr := net.JoinHostPort(target, strconv.Itoa(int(port))) // "host:port", корректно для IPv6
        conn, err := net.DialTimeout("tcp", addr, p.Timeout)      // handshake со своим таймаутом
        if err == nil {                                           // err==nil ⇒ порт OPEN
            open = append(open, strconv.Itoa(int(port)))
            conn.Close()                                          // только на успехе: иначе conn==nil → паника
        }
    }
    return Result{Probe: "ports", Target: target, Output: open}
}
```

## Ключевые функции
- `net.DialTimeout(network, address string, timeout time.Duration) (Conn, error)` — TCP-соединение
  со **своим** коротким потолком. `net.Dial` — без потолка, на тишине висит минуты (для сканера нельзя).
- `net.JoinHostPort(host, port string) string` — собрать `"host:port"` (оборачивает IPv6 в скобки;
  руками клеить нельзя). `port` — **строка**.
- `conn.Close()` — освободить соединение (ресурс ОС). Только при `err == nil`.
- `net.Listen("tcp", "127.0.0.1:0") (Listener, error)` — поднять слушателя; порт `0` = «ОС дай любой
  свободный». Реальный порт: `ln.Addr().(*net.TCPAddr).Port`. Для тестов сканера — эталонный «open»-порт.

## Личный опыт

### Что я понял не сразу (task-10)
- **`target` — всегда голый host.** Передал `"127.0.0.1:0"` (адрес слушателя) как `target` →
  `JoinHostPort` увидел двоеточие, решил IPv6, обернул в скобки → мусорный адрес → скан ничего не
  нашёл. Строка `host:port` и «host сканера» — разные вещи. Порт приходит из `p.Ports`, не из target.
- **`err == nil` ⇔ OPEN, и `conn` валиден только тогда.** `conn.Close()` вне `if err == nil` → на
  закрытом порту `conn == nil` → паника. `Close` строго внутри успешной ветки.
- **`net.Listen` без `Accept()` уже даёт open-порт.** Handshake (SYN-ACK) делает ядро ОС, не
  приложение — поэтому голого `Listen` в тесте хватает, чтобы порт был «открыт» для сканера. Прямое
  подтверждение теории из [[tcp-ip]] (порт open ⇔ handshake прошёл).
- **`uint16` как тип порта.** Порт 0–65535 влезает ровно в 16 бит. `[]uint16` вместо `[]int` —
  самодокументируется; расплата — `strconv.Itoa(int(port))` (Itoa хочет `int`).

### `DialTimeout` vs `DialContext` (task-12)

`net.DialTimeout` **не умеет отмену** — он подставляет контекст, который не отменяется никогда.
Исходник `net/dial.go`:

```go
func DialTimeout(network, address string, timeout time.Duration) (Conn, error) {
    d := Dialer{Timeout: timeout}
    return d.Dial(network, address)
}
func (d *Dialer) Dial(network, address string) (Conn, error) {
    return d.DialContext(context.Background(), network, address)   // <- вот здесь
}
```

Поэтому как только у функции в сигнатуре появился `ctx`, dial обязан идти через
`(&net.Dialer{Timeout: ...}).DialContext(ctx, ...)`, иначе отмена обрывается на этом месте.

- **Два независимых бюджета времени.** `Dialer.Timeout` — потолок на ОДИН dial, отсчёт с нуля
  на каждом. `ctx` — дедлайн на весь запуск. Побеждает **более ранний** срок, `dial.go:250-258`
  берёт `minNonzeroTime(now+Timeout, ctx.Deadline(), d.Deadline)`. Волна воркеров, стартующая
  на исходе общего бюджета, получит не полный `Timeout`, а остаток.
- Правило: per-dial маленький, общий дедлайн большой. Если они равны, первая же волна воркеров
  съедает весь бюджет, и остальные порты не проверяются (в `outpost` это пока один флаг
  `--timeout` на оба — записано в `BACKLOG.md`).
- **Мёртвый ctx не тормозит dial**: замер на заведомо открытом порту — 0 с / 552 мкс вместо
  `Timeout: 5s`. Ошибка обёрнута в `*net.OpError` (`err == ctx.Err()` → false), но
  `errors.Is(err, context.DeadlineExceeded)` → true.
- **Один `*net.Dialer` на 100 воркеров безопасен.** `DialContext` копирует значение себе
  (`sd := &sysDialer{Dialer: *d, ...}`, `dial.go:565`), общее состояние не мутируется.
  На самодельную заглушку эта гарантия не распространяется.

### Worker pool (стадия 2a, task-10)
Последовательный `for` по портам → пул воркеров. Ключевая структура:
```go
const workers = 100
jobs := make(chan uint16, len(p.Ports))        // буфер = число портов, не число воркеров
results := make(chan portResult, len(p.Ports)) // portResult{Port uint16; Open bool}

for _, port := range p.Ports { jobs <- port }  // producer: все порты в очередь
close(jobs)                                     // "задач больше не будет" (не стирает буфер!)

for range workers {                             // N воркеров СНАРУЖИ
    go func() {
        for port := range jobs {                // каждый воркер разбирает очередь
            conn, err := net.DialTimeout("tcp", net.JoinHostPort(...), p.Timeout)
            isOpen := err == nil
            if isOpen { conn.Close() }
            results <- portResult{port, isOpen} // слать на КАЖДЫЙ порт (open И closed)
        }
    }()
}

for i := 0; i < len(p.Ports); i++ {             // collector: ровно N результатов
    if r := <-results; r.Open { open = append(open, ...) }
}
```
- **Буфер канала ≠ число воркеров.** Воркеры = сколько горутин работают параллельно (обработка).
  Буфер = сколько задач влезет без блокировки (хранение). `jobs`/`results` буферим на `len(p.Ports)`.
- **`go func` × N снаружи, `for range jobs` внутри.** Не «горутина на порт», а фиксированный пул
  переиспользуемых воркеров: закончил порт → берёт следующий из очереди, пока не опустеет.
- **`close(jobs)` не стирает буфер** — воркеры дочерпывают остаток и только потом `range` выходит.
- **`close(results)` не нужен**, т.к. читаем счётным циклом ровно `len(p.Ports)`, не через `range`.
- **Гонки нет:** `append` в `open` только в collector (одна горутина); воркеры пишут лишь в канал.

### Сортировка результатов (стадия 2b, task-10)
Worker pool отдаёт открытые порты в **непредсказуемом порядке** (кто первый прислал в `results`).
Собираю их в `open []uint16` (числами!), сортирую `slices.Sort(open)` **до** конвертации в строки:
```go
var open []uint16
for range p.Ports {                       // ровно len(p.Ports) результатов
    if r := <-results; r.Open { open = append(open, r.Port) }
}
slices.Sort(open)                          // ЧИСЛЕННО: 80 < 443 (строки дали бы "443" < "80")
var output []string
for _, v := range open { output = append(output, strconv.Itoa(int(v))) }
```
- **Зачем:** без сортировки тест на `reflect.DeepEqual(Output, ...)` моргал бы — порядок из пула
  недетерминирован. Сортировка = детерминизм ради тестируемости.
- **Численно, не строкой:** сортируй пока `[]uint16`, конвертируй в строку в конце. Лексикографика
  строк дала бы `"443" < "80"`. Разбор — [[strings-sort]].

### На чём попадался
- `task-10`: `conn.Close()` вне `if` (дважды при правках); `target` = `host:port` вместо host.
- `task-10` (2a): отправка `results <- ...` **внутри** `if isOpen` → на закрытом порту недобор
  результатов → collector виснет → **deadlock**. Слать вне `if`, на каждый порт. Ловится только
  тестом с ЗАКРЫТЫМ портом (открытый порт баг маскирует — ложно-зелёный).

### Подводные камни
- `net.Dial` без таймаута на filtered-порте висит минуты → goroutine leak. Всегда `DialTimeout`.
- Сканировать чужие хосты без разрешения — незаконно. Легальные мишени для тренировки:
  `scanme.nmap.org`, свой `localhost`.

## Ключевые термины (English)
- `DialTimeout` — establish TCP connection with a timeout (no cancellation, see task-12)
- `DialContext` — same, but honours context cancellation; the only cancellable form
- `per-operation timeout` vs `deadline` — потолок на одну операцию против бюджета всего запуска
- `JoinHostPort` — build `host:port`, IPv6-safe
- `Listener` / `Listen` — accept incoming TCP connections; `:0` = OS-assigned free port
- `Conn.Close` — release the socket
- `worker pool` — fixed set of goroutines draining a jobs channel (реализовано, стадия 2a, [[channels]])
- `deadlock` — all goroutines blocked forever; here: collector waits for results nobody sends

## Связанные темы
[[tcp-ip]] [[channels]] [[context]] [[interfaces]]
