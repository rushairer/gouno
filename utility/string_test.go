package utility_test

import (
	"testing"

	"github.com/rushairer/gouno/utility"
)

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"foo_bar", "FooBar"},
		{"snake_case", "SnakeCase"},
		{"simple", "Simple"},
		{"a", "A"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := utility.ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello_world"},
		{"FooBar", "foo_bar"},
		{"SnakeCase", "snake_case"},
		{"123ABC", "123_abc"},
		{"Simple", "simple"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := utility.ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}
