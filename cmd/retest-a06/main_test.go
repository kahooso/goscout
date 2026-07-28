package main

import (
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

const timeout = 2 * time.Second
const duration = 10 * time.Second

type retest struct {
	name     string
	items    []string
	want     []string
	duration time.Duration
}

func work(s string) string {
	return strings.ToUpper(s)
}

func workTimeout(s string) string {
	time.Sleep(duration)
	return strings.ToUpper(s)
}

func TestProcessAll(t *testing.T) {
	tests := []retest{
		{
			name:     "good",
			items:    []string{"t", "e", "s", "t"},
			want:     []string{"E", "S", "T", "T"},
			duration: timeout,
		},
		{
			// nil на входе, а не []string{} — разные формы; фиксируем, что
			// на выходе в обоих случаях non-nil пустой слайс (см. make в ProcessAll)
			name:     "no items - slice",
			items:    []string(nil),
			want:     []string{},
			duration: timeout,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ProcessAll(tc.items, work, tc.duration)
			slices.Sort(got)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestProcessAllTimeout(t *testing.T) {
	tests := []retest{
		{
			name:     "long timeout",
			items:    []string{"T", "T", "E", "S"},
			duration: timeout,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			got := ProcessAll(tc.items, workTimeout, tc.duration)
			end := time.Since(start)

			// три независимых свойства: не сдалась раньше лимита, не дождалась всех
			// work (они спят дольше лимита), результат пуст. Проверено мутациями — см.
			// knowledge-base/interviews/retest-01.md
			if end < tc.duration {
				t.Errorf("func result has returned earlier than timeout (%v)", end)
			}
			if end >= tc.duration*2 {
				t.Errorf("func result has returned later than timeout*2 (%v)", end)
			}
			if len(got) != 0 {
				t.Errorf("func result slice length isn't zero")
			}
		})
	}
}
