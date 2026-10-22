package main

import "strings"

func parse(s string) alphametic {
	var a alphametic

	parts := strings.Split(s, "=")

	for _, w := range strings.Split(parts[0], "+") {
		a.addends = append(a.addends, strings.TrimSpace(w))
	}

	a.result = strings.TrimSpace(parts[1])

	words := append([]string{}, a.addends...)
	words = append(words, a.result)

	a.lead = make(map[rune]bool)
	seen := make(map[rune]bool)

	for _, w := range words {
		if len(w) > 1 {
			a.lead[rune(w[0])] = true
		}
		for _, c := range w {
			if !seen[c] {
				seen[c] = true
				a.letters = append(a.letters, c)
			}
		}
	}

	return a
}
