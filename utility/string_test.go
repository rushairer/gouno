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
	}

	for _, test := range tests {
		result := utility.ToCamelCase(test.input)
		if result != test.expected {
			t.Errorf("ToCamelCase(%s) = %s; want %s", test.input, result, test.expected)
		}
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
	}

	for _, test := range tests {
		result := utility.ToSnakeCase(test.input)
		if result != test.expected {
			t.Errorf("ToSnakeCase(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}
