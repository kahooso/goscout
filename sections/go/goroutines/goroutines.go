// task-06 (A0.5): параллельный обработчик задач через sync.WaitGroup.
// Разбор: knowledge-base/topics/go/goroutines.md
package main

import (
	"fmt"
	"sync"
)

func runAll(tasks []string, fn func(string)) {
	var wg sync.WaitGroup
	for _, v := range tasks {
		wg.Add(1)           // ДО `go`, иначе гонка с Wait()
		go func(v string) { // v параметром — не закрытие на переменную цикла
			defer wg.Done() // первой строкой — гарантия даже при панике fn
			fn(v)
		}(v)
	}
	wg.Wait()
}

func main() {
	servers := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
	runAll(servers, func(s string) {
		fmt.Println("checking: " + s)
	})
}
