package esl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/textproto"
	"strings"
)

// Message hold information about a command
type Message struct {
	MessageType EventType
	Headers     Headers
	Body        []byte
	Parsed      bool

	buf []byte
	tr  *textproto.Reader
}

// NewMessage - Will build and execute parsing against received freeswitch message.
// As return will give brand new Message{} for you to use it.
func NewMessage(buf []byte, autoParse bool) (*Message, error) {

	reader := bufio.NewReader(bytes.NewReader(buf))

	msg := Message{
		buf:     buf,
		tr:      textproto.NewReader(reader),
		Headers: NewHeaders(),
		Parsed:  false,
	}

	if autoParse {
		if err := msg.Parse(); err != nil {
			return &msg, err
		}
	}

	return &msg, nil
}

func (m *Message) String() string {
	return fmt.Sprintf("%s | %s", m.Headers, m.Body)
}

// Parse out message received from ESL and make it Go friendly.
func (m *Message) Parse() error {
	var err error
	var mime textproto.MIMEHeader

	mime, err = m.tr.ReadMIMEHeader()
	if err != nil {
		if err != io.EOF {
			return err
		}
	}
	err = nil
	m.Parsed = true

	for key, values := range mime {
		value := strings.Join(values, ";")
		m.Headers.Add(key, value)
	}

	contentLength := m.Headers.GetInt("Content-Length")
	if contentLength == 0 {
		return
	}

	return nil
}
