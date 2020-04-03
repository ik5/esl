package esl

import (
	"reflect"
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
