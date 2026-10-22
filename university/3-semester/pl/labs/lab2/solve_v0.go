package main

import (
	"strconv"
	"strings"
)

func wordValue(w string, digit map[rune]int) int {
	n := 0
	for _, c := range w {
		n = n*10 + digit[c]
	}
	return n
}

func check(a alphametic, digit map[rune]int) bool {
	for c := range a.lead {
		if digit[c] == 0 {
			return false
		}
	}

	sum := 0
	for _, w := range a.addends {
		sum += wordValue(w, digit)
	}
	return sum == wordValue(a.result, digit)
}

func solveV0(s string) string {
	a := parse(s)
	digit := make(map[rune]int)
	used := make(map[int]bool)

	if assign(a, 0, digit, used) {
		return answer(a, digit)
	}
	return ""
}

func assign(a alphametic, i int, digit map[rune]int, used map[int]bool) bool {
	if i == len(a.letters) {
		return check(a, digit)
	}

	c := a.letters[i]
	for d := 0; d <= 9; d++ {
		if used[d] {
			continue
		}
		digit[c] = d
		used[d] = true
		if assign(a, i+1, digit, used) {
			return true
		}
		used[d] = false
	}
	return false
}

func answer(a alphametic, digit map[rune]int) string {
	var nums []string
	for _, w := range a.addends {
		nums = append(nums, strconv.Itoa(wordValue(w, digit)))
	}
	return strings.Join(nums, " + ") + " = " + strconv.Itoa(wordValue(a.result, digit))
}
