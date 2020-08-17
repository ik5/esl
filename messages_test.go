package esl

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestSimpleMessageAutoParse(t *testing.T) {
	type messageType struct {
		Headers Headers
		Body    []byte
		Parsed  bool
	}
	type fixture struct {
		input    []byte
		expected messageType
	}

	msgContent := []fixture{
		{
			input: []byte("Event-Name: SOCKET_DATA\r\nContent-Type: auth/request\r\n"),
			expected: messageType{
				Body:   nil,
				Parsed: true,
				Headers: Headers{
					header: map[string]interface{}{
						"Event-Name":   "SOCKET_DATA",
						"Content-Type": "auth/request",
					},
				},
			},
		},
	}

	for idx, content := range msgContent {
		msg, err := NewMessage(content.input, true)
		if err != nil {
			t.Errorf("Unable to parse msg (#%d): %s", idx, err)
		}

		if msg == nil {
			t.Errorf("msg #%d is empty", idx)
		}
		if msg.Parsed != content.expected.Parsed {
			t.Errorf("Expected (#%d) Parsed to be %t got: %t", idx, content.expected.Parsed, msg.Parsed)
		}

		if !reflect.DeepEqual(msg.Body, content.expected.Body) {
			t.Errorf("Expected (%d) Body %+v is not %+v", idx, msg.Body, content.expected.Body)
		}

		for key, head := range content.expected.Headers.header {
			head2 := msg.Headers.GetString(key)
			if head != head2 {
				t.Errorf("Expected (%d) header %s %+v is not %+v", idx, key, head2, head)
			}
		}
	}

}

func TestMessageAsString(t *testing.T) {
	type fixture struct {
		input    []byte
		expected string
	}
	msgContent := []fixture{
		{
			input:    []byte("Content-Length: 27\n\n-ERR af Command not found!\n\n"),
			expected: "Headers: Content-Length=27 ; Body: -ERR af Command not found!",
		},
		{
			input:    []byte("Content-Type: auth/request\n\n"),
			expected: "Headers: Content-Type=auth/request",
		},
	}

	for idx, msg := range msgContent {
		message, err := NewMessage(msg.input, true)
		if err != nil {
			t.Errorf("Unable to parse message (%d): %s", idx, err)
		}

		msgStr := message.String()
		found := strings.Compare(msg.expected, msgStr)
		if found != 0 {
			t.Errorf("Expected (%d): \n'%s', got \n'%s' (%d)\n", idx, msgStr, msg.expected, found)
		}
	}
}

func TestMessageParseContentLengthZero(t *testing.T) {
	buf := []byte("Content-Length: 0\n")

	_, err := NewMessage(buf, true)
	if err == nil {
		t.Errorf("No error found.")
	}

	if !errors.Is(err, ErrContentLengthZero) {
		t.Errorf("Unexpected error found: %s", err)
	}
}

func TestMessageParseContentLength(t *testing.T) {
	buf := []byte("Content-Length: 1\n\nhello world")

	msg, err := NewMessage(buf, true)
	if err != nil {
		t.Errorf("Error parsing new message: %s", err)
		return
	}

	l := len(msg.Body)
	if l != 1 {
		t.Errorf("Expected len of 1, got: %d", l)
	}
}

func TestMessageHasError(t *testing.T) {
	type fixture struct {
		input    []byte
		expected bool
	}

	fixtures := []fixture{
		{
			input:    []byte("Content-Length: 18\n\n-ERR Testing error message\n"),
			expected: true,
		},
		{
			input:    []byte("Content-Type: command/reply\nReply-Text: -ERR Testing error message\n"),
			expected: true,
		},
		{
			input:    []byte("Content-Type: command/reply\nReply-Text: +OK log level 9999\n"),
			expected: false,
		},
		{
			input:    []byte("Content-Type: api/response\nContent-Length: 12\n\nlog 9999"),
			expected: false,
		},
	}

	for idx, msg := range fixtures {
		message, err := NewMessage(msg.input, true)
		if err != nil {
			t.Errorf("Unexpected error parsing (%d): %s", idx, msg.input)
			continue
		}

		hasError := message.HasError()
		if hasError != msg.expected {
			t.Errorf("Expected (%d): %t but got %t", idx, msg.expected, hasError)
			continue
		}
	}
}

func TestMessageError(t *testing.T) {
	type fixture struct {
		input    []byte
		expected error
	}

	fixtures := []fixture{
		{
			input:    []byte("Content-Type: api/response\nContent-Length: 18\n\n-ERR Testing error message\n"),
			expected: errors.New("Testing error message"),
		},
		{
			input:    []byte("Content-Type: command/reply\nReply-Text: -ERR Testing error message\n"),
			expected: errors.New("Testing error message"),
		},
		{
			input:    []byte("Content-Type: command/reply\nReply-Text: +OK log level 9999\n"),
			expected: nil,
		},
		{
			input:    []byte("Content-Type: api/response\nContent-Length: 12\n\nlog 9999"),
			expected: nil,
		},
		{
			input:    []byte("Content-Type: text/plain\nContent-Length: 25\n\n-ERR Invalid JSON format\n"),
			expected: nil,
		},
		{
			input:    []byte("Conten-Type: text/plain\nReply-Text: -ERR Invalid JSON format"),
			expected: nil,
		},
	}

	for idx, msg := range fixtures {
		message, err := NewMessage(msg.input, true)
		if err != nil {
			t.Errorf("Unexpected error parsing (%d): %s", idx, msg.input)
			continue
		}

		err = message.Error()
		if err == nil && msg.expected == nil {
			continue
		}
		if err == nil && msg.expected != nil {
			t.Errorf("Expected (%d) err to be %s but got nil", idx, msg.expected)
			continue
		}
		if errors.Is(err, msg.expected) {
			t.Errorf("Expected (%d): %s but got %s", idx, msg.expected, err)
			continue
		}
	}

}
