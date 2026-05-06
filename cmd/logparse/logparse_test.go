package main

import (
	"reflect"
	"testing"
)

func TestCountByMessage(t *testing.T) {
	tests := []struct {
		name string
		logs []LogEntry
		want map[string]int
	}{
		{
			name: "одно сообщение - несколько раз",
			logs: []LogEntry{
				{Message: "Hey there"},
				{Message: "Hey there"},
				{Message: "I am using WhatsApp"},
			},
			want: map[string]int{"Hey there": 2, "I am using WhatsApp": 1},
		},
		{
			name: "пустой вход",
			logs: []LogEntry{},
			want: map[string]int{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := countByMessage(tc.logs)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCountByLevel(t *testing.T) {
	tests := []struct {
		name string
		logs []LogEntry
		want map[string]int
	}{
		{
			name: "несколько уровней",
			logs: []LogEntry{
				{Level: "INFO"},
				{Level: "ERROR"},
				{Level: "INFO"},
			},
			want: map[string]int{"INFO": 2, "ERROR": 1},
		},
		{
			name: "пустой вход",
			logs: []LogEntry{},
			want: map[string]int{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := countByLevel(tc.logs)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTopMessages(t *testing.T) {
	tests := []struct {
		name string
		top  map[string]int
		n    int
		want []string
	}{
		{
			name: "n < len(top)",
			top:  map[string]int{"Hey there": 2, "I am using...": 1},
			n:    1,
			want: []string{"Hey there"},
		},
		{
			name: "n > len(top)",
			top:  map[string]int{"Hey there": 2},
			n:    2,
			want: []string{"Hey there"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := topMessages(tc.top, tc.n)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
