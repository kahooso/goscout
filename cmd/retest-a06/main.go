package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func ProcessAll(items []string, work func(string) string, timeout time.Duration) []string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ch := make(chan string, len(items))
	for _, item := range items {
		go func(s string) {
			ch <- work(s)
		}(item)
	}

	output := make([]string, 0, len(items))
	for range len(items) {
		select {
		case res := <-ch:
			output = append(output, res)
		case <-ctx.Done():
			return output
		}
	}
	return output
}

func main() {
	res := ProcessAll(
		[]string{"T", "E", "S", "T"},
		func(s string) string { return strings.ToUpper(s) },
		2*time.Second,
	)
	fmt.Printf("%+v\tlen=%v\tcap=%v\tres==nil=%t\n", res, len(res), cap(res), res == nil)
}
