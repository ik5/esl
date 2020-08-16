package esl

import (
	"bufio"
	"bytes"
	"errors"
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

	if m.Headers.Exists("Content-Length") {
		contentLength := m.Headers.GetInt("Content-Length")
		if contentLength == 0 {
			return errors.New("Content Length is zero")
		}

		l := int(m.Headers.GetInt("Content-Length"))
		lines := make([]byte, 0, l)

		for err == nil {
			line, err := m.tr.ReadLineBytes()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			lines = append(lines, line...)
			if len(lines) >= l {
				lines = lines[:l]
				break
			}
		}
		m.Body = lines
		return nil
	}

	return nil
}

// ContentType returns the content type arrived, or empty string if not found
func (m *Message) ContentType() EventContentType {
	contentType := m.Headers.GetString("Content-Type")

	return EventContentType(contentType)
}

// HasError return true if a message got error from freeswitch message
func (m *Message) HasError() bool {
	if m.Headers.Exists("Content-Length") {
		return bytes.HasPrefix(m.Body, []byte("-ERR"))
	}
	err := m.Headers.GetString("Reply-Text")
	return strings.HasPrefix(err, "-ERR")
}

// Error return the message error msg or empty string if non found
func (m *Message) Error() string {
	if !m.HasError() {
		return ""
	}

	if m.Headers.Exists("Content-Length") {
		return string(m.Body[5:])
	}

	switch m.ContentType() {
	case ECTCommandReply:
		err := m.Headers.GetString("Reply-Text")
		return err[5:]
	case ECTAPIResponse:
		return string(m.Body)
	}

	return ""
}
