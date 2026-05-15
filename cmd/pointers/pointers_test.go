package main

import (
	"reflect"
	"testing"
)

func TestPop(t *testing.T) {
	tests := []struct {
		name    string
		items   []int
		want    []int
		wantVal int
		wantOk  bool
	}{
		{
			name:    "pop с элементами",
			items:   []int{1, 2, 5, 7},
			want:    []int{1, 2, 5},
			wantVal: 7,
			wantOk:  true,
		},
		{
			name:    "pop на пустом стеке",
			items:   []int{},
			want:    []int{},
			wantVal: 0,
			wantOk:  false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Stack{items: tc.items}
			v, ok := s.Pop()
			if v != tc.wantVal || ok != tc.wantOk || !reflect.DeepEqual(s.items, tc.want) {
				t.Errorf("got (%v, %v, %v), want (%v, %v, %v)", v, ok, s.items, tc.wantVal, tc.wantOk, tc.want)
			}
		})
	}
}

func TestPush(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  []int
	}{
		{
			name:  "несколько push",
			items: []int{},
			want:  []int{1, 2, 3},
		},
		{
			name:  "push на пустой stack",
			items: []int{},
			want:  []int{1},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Stack{items: tc.items}
			for i := range tc.want {
				s.Push(tc.want[i])
			}
			if !reflect.DeepEqual(s.items, tc.want) {
				t.Errorf("got %v, want %v", s.items, tc.want)
			}
		})
	}
}

func TestPeek(t *testing.T) {
	tests := []struct {
		name    string
		items   []int
		wantVal int
		wantOk  bool
	}{
		{
			name:    "peek с элементами",
			items:   []int{1, 2, 4, 9},
			wantVal: 9,
			wantOk:  true,
		},
		{
			name:    "peek на пустом стеке",
			items:   []int{},
			wantVal: 0,
			wantOk:  false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Stack{items: tc.items}
			gotVal, gotOk := s.Peek()
			if gotVal != tc.wantVal || gotOk != tc.wantOk {
				t.Errorf("got (%v, %v), want (%v, %v)", gotVal, gotOk, tc.wantVal, tc.wantOk)
			}
		})
	}
}
