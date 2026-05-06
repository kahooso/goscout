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
	// A: make — встроенная полиморфная функция, работает для трёх типов: slice, map, channel.
	// Параметры у каждого типа свои, компилятор проверяет их на этапе компиляции.
	// make([]T, len, cap) — slice: len элементов заполненных нулём, cap зарезервировано.
	// make(map[K]V, hint) — map: hint задаёт начальную ёмкость (опционально).
	// make(chan T, buf)    — channel: buf = 0 → unbuffered, buf > 0 → буферизованный.
	// Если вставить неправильные параметры — ошибка компиляции, не паника.

	for _, l := range logs {
		m[l.Level]++
	}

	return m
}

func topMessages(mTop map[string]int, n int) []string {
	top := make([]string, 0, len(mTop))
	for k := range mTop {
		top = append(top, k)
	}
	// A: В Go только camelCase — ни snake_case, ни ALL_CAPS.
	// Экспортируемое (публичное) → MTop; приватное → mTop.
	// Короткие имена (m, l, i, v) нормальны в маленьком скоупе — это идиома Go, не лень.
	// В широком скоупе — описательное: mTop понятнее чем m, потому что map не одна.

	sort.Slice(top, func(i, j int) bool {
		return mTop[top[i]] > mTop[top[j]]
	})
	// A: Можно сортировать []LogEntry по Message полю — но тогда нужно после сортировки
	// извлекать Message из каждого LogEntry. Сложнее и дублируем сообщения.
	// Текущий подход — итерируем ключи из map — проще: ключи уже уникальные строки,
	// счётчики доступны через mTop, слайс готов к возврату без дополнительного маппинга.

	if n > len(top) {
		return top
	}
	return top[:n]
	// В данном случае, я как бы возвращаю срез слайса
	// Получается, что мой слайс в main будет ссылаться
	// не на весь слайс top, а только от 0 до n элементов.
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
