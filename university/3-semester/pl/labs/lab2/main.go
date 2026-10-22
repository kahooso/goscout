package main

import (
	"fmt"
	"time"
)

func main() {
	puzzles := []string{
		"BE + BE = MOO",
		"SEND + MORE = MONEY",
		"ELEVEN + NINE + FIVE + FIVE = THIRTY",
	}

	for _, p := range puzzles {
		start := time.Now()
		res := solveV0(p)
		fmt.Printf("%-40s → %-35s %v\n", p, res, time.Since(start))
	}
}
