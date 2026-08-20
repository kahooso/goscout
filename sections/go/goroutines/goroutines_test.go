package main

import (
	"sync/atomic"
	"testing"
)

func TestRunAll(t *testing.T) {
	tests := []struct {
		name  string
		tasks []string
		want  int64
	}{
		{
			name:  "пустой слайс",
			tasks: []string{},
			want:  0,
		},
		{
			name:  "несколько элементов",
			tasks: []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"},
			want:  3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var counter int64
			runAll(tc.tasks, func(_ string) {
				atomic.AddInt64(&counter, 1)
			})

			if counter != tc.want {
				t.Errorf("runAll: got %v, want %v", counter, tc.want)
			}
		})
	}
}
