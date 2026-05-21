package twosum

import (
	"testing"
)

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		lst    []int
		target int
		v1     int
		v2     int
		ok     bool
	}{
		{
			name:   "андер",
			lst:    []int{2, 9, 6, 3, 5},
			target: 9,
			v1:     2,
			v2:     3,
			ok:     true,
		},
		{
			name:   "булдыр",
			lst:    []int{3, 3},
			target: 6,
			v1:     0,
			v2:     1,
			ok:     true,
		},
		{
			name:   "бойдер",
			lst:    []int{2, 9, 6, 3, 5},
			target: 10,
			v1:     0,
			v2:     0,
			ok:     false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v1, v2, ok := TwoSum(tc.lst, tc.target)
			if ok != tc.ok {
				t.Errorf("want ok=%t, got ok=%t", ok, tc.ok)
			}
			if tc.ok && (v1 != tc.v1 || v2 != tc.v2) {
				t.Errorf("want v1=%d, v2=%d, ok=%t\tgot v1=%d, v2=%d, ok=%t",
					tc.v1, tc.v2, tc.ok, v1, v2, ok)
			}
		})
	}
}
