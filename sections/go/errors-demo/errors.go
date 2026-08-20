// task-05 (A0.4): парсер конфига с типизированными ошибками.
// Разбор: knowledge-base/topics/go/errors.md
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Sentinel error — аналог errno-кодов в C (ENOENT, EACCES).
// Сравнение через errors.Is, не через строки.
var ErrMissingField = errors.New("missing field")

// ParseError — кастомный тип ошибки для невалидного значения поля.
type ParseError struct {
	Field string
	Value string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("invalid value for %s: %q", e.Field, e.Value)
}

type Config struct {
	Host    string
	Port    int
	Timeout int
}

func parseConfig(input string) (*Config, error) {
	if len(input) == 0 {
		return &Config{}, nil
	}

	c := &Config{}
	for _, field := range strings.Fields(input) {
		kv := strings.SplitN(field, "=", 2)
		if len(kv) != 2 {
			continue // TODO: после A0.6 переделать на ошибку
		}
		key, value := kv[0], kv[1]

		switch key {
		case "host":
			c.Host = value
		case "port":
			port, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("parseConfig: %w", &ParseError{Field: "port", Value: value})
			}
			c.Port = port
		case "timeout":
			timeout, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("parseConfig: %w", &ParseError{Field: "timeout", Value: value})
			}
			c.Timeout = timeout
		}
	}

	if c.Host == "" || c.Port == 0 || c.Timeout == 0 {
		return nil, fmt.Errorf("parseConfig: %w", ErrMissingField)
	}
	return c, nil
}

func main() {
	c, err := parseConfig("host=localhost port=8080 timeout=30")
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("config: %+v\n", *c)
	}

	_, err = parseConfig("host=localhost timeout=30")
	if errors.Is(err, ErrMissingField) {
		fmt.Println("missing required field:", err)
	}

	_, err = parseConfig("host=localhost port=abc timeout=30")
	var pe *ParseError
	if errors.As(err, &pe) {
		fmt.Printf("incorrect value: field=%s value=%s\n", pe.Field, pe.Value)
	}
}
