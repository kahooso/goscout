package main

type alphametic struct {
	addends []string
	result  string
	letters []rune
	lead    map[rune]bool
}
