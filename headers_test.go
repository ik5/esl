package esl

import (
	"strings"
	"testing"
)

func TestHeadersGetString(t *testing.T) {
	type fixture struct {
		key      string
		value    interface{}
		expected string
	}

	fixtures := []fixture{
		{
			key:      "int",
			value:    int(42),
			expected: "42",
		},
		{
			key:      "int32",
			value:    int32(1),
			expected: "1",
		},
		{
			key:      "uint16",
			value:    uint16(100),
			expected: "100",
		},
		{
			key:      "float64",
			value:    float64(3.14),
			expected: "3.140",
		},
		{
			key:      "string",
			value:    "hello world",
			expected: "hello world",
		},
		{
			key:      "nil",
			value:    nil,
			expected: "",
		},
	}

	headers := NewHeaders()

	for _, test := range fixtures {
		headers.Add(test.key, test.value)
	}

	for idx, test := range fixtures {
		s := headers.GetString(test.key)

		if s != test.expected {
			t.Errorf("Expected (%d) result of '%s', got '%s'", idx, s, test.expected)
		}
	}
}

func TestHeadersGetInt(t *testing.T) {
	type fixture struct {
		key      string
		value    interface{}
		expected int64
	}

	fixtures := []fixture{
		{
			key:      "int",
			value:    int(42),
			expected: 42,
		},
		{
			key:      "string",
			value:    "hello world",
			expected: 0,
		},
		{
			key:      "int32",
			value:    int32(1),
			expected: 1,
		},
		{
			key:      "uint16",
			value:    uint16(100),
			expected: 100,
		},
		{
			key:      "float64",
			value:    float64(3.14),
			expected: 3,
		},
		{
			key:      "nil",
			value:    nil,
			expected: 0,
		},
	}

	headers := NewHeaders()

	for _, test := range fixtures {
		headers.Add(test.key, test.value)
	}

	for idx, test := range fixtures {
		i := headers.GetInt(test.key)

		if i != test.expected {
			t.Errorf("Expected (%d) result of '%d', got '%d'", idx, i, test.expected)
		}
	}

}

func TestHeadersExists(t *testing.T) {
	headers := NewHeaders()

	if headers.Exists("foo") {
		t.Errorf("No key was set, but 'foo' was found")
		return
	}

	headers.Add("foo", "bar")

	if !headers.Exists("foo") {
		t.Errorf("'foo' was set, but not found")
		return
	}
}

func TestHeadersRemove(t *testing.T) {
	headers := NewHeaders()

	headers.Add("foo", "bar")
	if !headers.Exists("foo") {
		t.Errorf("expected to find key 'foo'")
		return
	}

	headers.Remove("foo")

	if headers.Exists("foo") {
		t.Errorf("'foo' needed to be removed, but still exists")
		return
	}

	headers.Remove("bar")
	if headers.Exists("bar") {
		t.Errorf("'bar' did not exists, removed, but now it is")
		return
	}
}

func TestHeadersString(t *testing.T) {
	headers := NewHeaders()

	s := headers.String()

	if s != "" {
		t.Errorf("Expected empty string for empty headers, but got '%s'", s)
		return
	}

	headers.Add("foo", "bar")

	s = headers.String()

	if s != "foo=bar" {
		t.Errorf("Expected 'foo=bar', but got: '%s'", s)
		return
	}

	headers.Add("int", 10)

	s = headers.String()
	a := strings.Contains(s, "foo=bar")
	b := strings.Contains(s, "int=10")
	sep := strings.Contains(s, " | ")
	if !a || !b || !sep {
		t.Errorf("'foo=bar': %t | 'int=10': %t | sep: %t | '%s'", a, b, sep, s)
		return
	}
}

func TestHeadersLen(t *testing.T) {
	headers := NewHeaders()

	l := headers.Len()
	if l > 0 {
		t.Errorf("Expected len of 0, got: %d", l)
		return
	}

	headers.Add("foo", "bar")
	l = headers.Len()

	if l != 1 {
		t.Errorf("Expected len of 1, got %d", l)
		return
	}

	headers.Add("int", 10)
	l = headers.Len()

	if l != 2 {
		t.Errorf("Expected len of 2, got %d", l)
		return
	}

	headers.Remove("foo")
	l = headers.Len()
	if l != 1 {
		t.Errorf("Expected len of 1, got %d", l)
		return
	}
}
