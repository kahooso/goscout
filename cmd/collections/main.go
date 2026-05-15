// task-02 (A0.2): счётчик слов на map + сортировка по убыванию частоты.
// Подробный разбор: knowledge-base/topics/go/slices-maps.md, topics/go/strings-sort.md
package main

import (
	"fmt"
	"sort"
	"strings"
)

const msg = "Hey there! I am using WhatsApp! WhatsApp"
const cutset = ",.!?:;\"'()`"

func main() {
	if len(msg) == 0 {
		fmt.Println("msg string is empty")
		return
	}

	counts := countWords(msg)
	keys := sortByFreq(counts)
	for _, v := range keys {
		fmt.Printf("%s : %d\n", v, counts[v])
	}
}

func countWords(text string) map[string]int {
	counts := make(map[string]int)
	if len(text) == 0 {
		return counts
	}
	for _, v := range strings.Fields(text) {
		v = strings.Trim(strings.ToLower(v), cutset)
		counts[v]++
	}
	return counts
}

func sortByFreq(counts map[string]int) []string {
	if counts == nil {
		return nil
	}
	keys := make([]string, 0, len(counts)) // pre-allocate — знаем размер
	for k := range counts {
		keys = append(keys, k)
	}
	// less строгое: > и < нельзя смешивать, иначе сортировка ломается
	sort.Slice(keys, func(i, j int) bool {
		return counts[keys[i]] > counts[keys[j]]
	})
	return keys
}
