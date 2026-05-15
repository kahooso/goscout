// task-03 (A0.2-bis): CSV-логи — парсинг, агрегация, топ-N.
// Учебная реализация — в проде использовать encoding/csv.
// Разбор: knowledge-base/topics/go/slices-maps.md
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

type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

func parseLogs(lines []string) []LogEntry {
	logs := make([]LogEntry, len(lines))
	for i, v := range lines {
		parts := strings.Split(v, csvSep)
		// TODO(post-A0.6): len(parts) < 3 → возвращать ошибку, не молча индексировать
		logs[i] = LogEntry{
			Timestamp: parts[0],
			Level:     parts[1],
			Message:   parts[2],
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
	for _, l := range logs {
		m[l.Level]++
	}
	return m
}

func topMessages(messages map[string]int, n int) []string {
	keys := make([]string, 0, len(messages))
	for k := range messages {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return messages[keys[i]] > messages[keys[j]]
	})
	if n > len(keys) {
		return keys
	}
	return keys[:n]
}

func main() {
	logs := parseLogs(csvLog)
	levels := countByLevel(logs)
	messages := countByMessage(logs)
	top := topMessages(messages, 1)

	fmt.Println("Уровни:")
	for k, v := range levels {
		fmt.Printf("%s\t: %d\n", k, v)
	}
	fmt.Println()

	fmt.Println("Топ сообщений:")
	for _, v := range top {
		fmt.Printf("%s\t: %d\n", v, messages[v])
	}
}
