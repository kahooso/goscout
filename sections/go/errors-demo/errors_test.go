package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name  string
		input string
		cfg   *Config
		err   error
	}{
		{
			name:  "happy path",
			input: "host=localhost port=8080 timeout=30",
			cfg: &Config{
				Host:    "localhost",
				Port:    8080,
				Timeout: 30,
			},
			err: nil,
		},
		{
			name:  "отсутствующее поле",
			input: "host=localhost port=8080",
			cfg:   nil,
			err:   ErrMissingField,
		},
		{
			name:  "неверное значение port",
			input: "host=localhost port=abc timeout=30",
			cfg:   nil,
			err: &ParseError{
				Field: "port",
				Value: "abc",
			},
		},
		{
			name:  "неверное значение timeout",
			input: "host=localhost port=8080 timeout=abc",
			cfg:   nil,
			err: &ParseError{
				Field: "timeout",
				Value: "abc",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := parseConfig(tc.input)

			if tc.err == nil {
				if err != nil {
					t.Fatalf("got error %v, want nil", err)
				}
				if !reflect.DeepEqual(cfg, tc.cfg) {
					t.Errorf("got %+v, want %+v", cfg, tc.cfg)
				}
				return
			}

			if err == nil {
				t.Fatalf("got nil, want error")
			}

			if errors.Is(tc.err, ErrMissingField) {
				if !errors.Is(err, ErrMissingField) {
					t.Errorf("Want ErrMissingField, got %v", err)
				}
				return
			}

			if wantPe, ok := tc.err.(*ParseError); ok {
				var gotPe *ParseError
				if !errors.As(err, &gotPe) {
					t.Errorf("want *ParseError, got %T", err)
					return
				}
				if gotPe.Field != wantPe.Field {
					t.Errorf("got field=%s, want=%s", gotPe.Field, wantPe.Field)
				}
			}
		})
	}
}
