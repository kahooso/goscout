package main

import (
	"testing"
)

type TestEntry struct {
	name string
	logs []LogEntry
	want map[string]int
}

func TestCountByLevel(t *testing.T) {
	// tests := []TestEntry{
	// 	{
	// 		name: "три уровня",
	// 		logs: []LogEntry{
	// 			{Level: "INFO"},
	// 			{Level: "ERROR"},
	// 			{Level: "INFO"},
	// 		},
	// 		want: map[string]int{"INFO": 2, "ERROR": 1},
	// 	},
	// 	{
	// 		name: "пустой вход",
	// 		logs: []LogEntry{},
	// 		want: map[string]int{},
	// 	},
	// }
	// for _, tc := range tests {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		got := countByLevel(tc.logs)
	// 	})
	// }
}
