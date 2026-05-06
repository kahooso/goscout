package main

import (
	"fmt"
	"sort"
	"strings"
)

var csvLog = []string{
	"2026-01-01 10:00:00,INFO,Server started",
	"2026-01-01 10:01:00,ERROR,Connection refused",
	"2026-01-01 10:02:00,INFO,Request processed",
	"2026-01-01 10:03:00,WARN,High memory usage",
	"2026-01-01 10:04:00,ERROR,Timeout",
	"2026-01-01 10:05:00,INFO,Request processed",
}

const csvSep = ","

// A: В Go одиночные кавычки '' — это rune-литерал (один символ, тип int32), не строка.
// Для строк только два варианта:
//   "строка"  — интерпретируемая: понимает \n, \t, \"
//   `строка`  — raw string: \n это буквально два символа \ и n, без escape.
// Raw строки удобны для регулярок — не нужно экранировать \d, \w и т.д.

type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

func parseLogs(lines []string) []LogEntry {
	logs := make([]LogEntry, len(lines))
	for i, v := range lines {
		temp := strings.Split(v, csvSep)
		logs[i] = LogEntry{
			Timestamp: temp[0],
			Level:     temp[1],
			Message:   temp[2],
		}
	}

	return logs
}

func countByMessage(logs []LogEntry) map[string]int {
	m := make(map[string]int)
	for _, l := range logs {
		m[l.Message]++
	}

	return m
}

func countByLevel(logs []LogEntry) map[string]int {
	m := make(map[string]int)
	// A: make — встроенная функция для трёх типов: slice, map, channel. Полиморфная.
	// Параметры для каждого типа свои, и компилятор проверяет — неправильные дадут
	// ошибку компиляции, не панику.
	//   make([]T, len, cap)       — slice: len обязателен, cap опционален
	//   make(map[K]V)             — map: параметр опционален (начальный размер хранилища)
	//   make(chan T, bufSize)      — channel: bufSize=0 или пусто → unbuffered
	// Смешивать параметры нельзя: make(map[string]int, 0, 10) — ошибка компиляции.

	for _, l := range logs {
		m[l.Level]++
	}

	return m
}

func topMessages(mTop map[string]int, n int) []string {
	top := make([]string, 0, len(mTop)) // УЛУЧШЕНИЕ: var top []string = make(...) — C-стиль; в Go просто :=
	for k := range mTop {
		top = append(top, k)
	}
	// A: В Go всё camelCase. snake_case и ALL_CAPS — не принято (C/Python стиль).
	//   m, v, k     — нормально в коротком scope (тело цикла, 3-5 строк)
	//   mTop        — хорошо для параметра функции: коротко и понятно
	//   messages    — хорошо для переменной в main: читаемо
	// Экспортируемое (доступно снаружи пакета): CamelCase с большой буквы — LogEntry, ParseLogs.
	// Приватное (только внутри пакета): camelCase с маленькой — csvLog, countByLevel.

	sort.Slice(top, func(i, j int) bool {
		return mTop[top[i]] > mTop[top[j]]
	})
	// A: sort.Slice(logs, ...) можно, но тогда возвращать []LogEntry, а не []string.
	// Текущий подход (ключи из mTop) чище: функция возвращает именно строки-сообщения,
	// а данные о частоте уже есть у вызывающего (mTop передан снаружи).

	if n > len(top) {
		return top
	}
	return top[:n]
	// A: Правильно. top[:n] возвращает новый заголовок {ptr: тот же, len: n, cap: len(top)}.
	// Данные не копируются — оба слайса указывают на один массив в heap.
	// Изменение элементов через возвращённый срез было бы видно в top и наоборот.
}

func main() {
	logs := parseLogs(csvLog)
	levels := countByLevel(logs)
	messages := countByMessage(logs)
	top := topMessages(messages, 1)

	fmt.Print("Уровни:\n")
	for k, v := range levels {
		fmt.Printf("%s\t: %d\n", k, v)
	}
	fmt.Println()

	fmt.Print("Топ сообщений:\n")
	for _, v := range top {
		fmt.Printf("%s\t: %d\n", v, messages[v])
	}
	fmt.Println()
}
