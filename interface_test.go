package esl

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestBasicConnectionSendRecv(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	socket, err := Dial(eslHost, eslPasword, 1, 30*time.Second)
	if err != nil {
		t.Errorf("Dial error: %s", err)
		return
	}

	if socket == nil {
		t.Errorf("Socket is nil")
		return
	}
	defer socket.Close()

	n, content, err := socket.Recv(AuthRequestBufferSize)
	if err != nil {
		t.Errorf("Unable to read content: %s", err)
	}

	if content == nil {
		t.Error("Content is nil")
	}

	if len(content) == 0 {
		t.Error("Content cannot be empty")
	}

	if int64(n) >= AuthRequestBufferSize {
		t.Errorf("n cannot be 28")
	}

	found := bytes.HasPrefix(content, []byte("Content-Type: auth/request\n\n"))
	if !found {
		t.Errorf("Expected 'Content-Type: auth/request', got: '%s' | %+v | side: %t", content, content, found)
	}

	n, content, err = socket.SendRecv("auth " + eslPasword)
	if err != nil {
		t.Errorf("Login error: %s", err)
		return
	}

	if n != 54 {
		t.Errorf("Expected length of 54, got: %d", n)
	}

	expected := []byte("Content-Type: command/reply\nReply-Text: +OK accepted\n\n")

	cmp := bytes.Compare(expected, content[0:n])
	if cmp != 0 {
		t.Errorf("Expected '%s' %+v \ngot: '%s' | %+v | side: %d\n",
			expected, expected, content, content, cmp)
	}
}

func TestSocketSendRecevError(t *testing.T) {
	socket := &Socket{}

	n, buf, err := socket.SendRecv("foo")
	if err == nil {
		t.Errorf("Expected error but got nil")
	}

	if buf != nil {
		t.Errorf("Expected nil buf, but got: '%s'", buf)
	}

	if n != 0 {
		t.Errorf("Expected n 0 but got '%d'", n)
	}
}

func TestSocketAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	socket, err := Connect(eslHost, eslPasword, 1, 30*time.Second)
	if err != nil {
		t.Errorf("Unable to connect: %s", err)
		return
	}

	if socket == nil {
		t.Errorf("Socket is nil without an error")
		return
	}
	defer socket.Close()

	msg, err := socket.API("show", "api")
	if err != nil {
		t.Errorf("Unable to call API: %s", err)
		return
	}

	if msg == nil {
		t.Errorf("No message returned without an error")
		return
	}

	ct := msg.ContentType()
	if ct != ECTAPIResponse {
		t.Errorf("Invalid content type, expected '%s' got '%s'", ECTAPIResponse, ct)
		return
	}

	l := msg.Headers.GetInt("Content-Length")
	if l <= 0 {
		t.Errorf("l (%d) cannot be empty", l)
		return
	}

	bodyLength := int64(len(msg.Body))
	if bodyLength != l-1 && bodyLength != l {
		t.Errorf("Expected body length to be %d but got %d", l, bodyLength)
		return
	}

	if !bytes.HasSuffix(msg.Body, []byte(" total.")) {
		t.Errorf("Unexpected end of API commands.")
		return
	}

}

func TestSocketAPIConnectError(t *testing.T) {
	socket := &Socket{}

	msg, err := socket.API("show", "")
	if err == nil {
		t.Errorf("Expected err, but nil returned")
	}

	if msg != nil {
		t.Errorf("Expected msg to be nil, but got '%s'", msg)
	}
}

func TestSocketBgAPI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	socket, err := Connect(eslHost, eslPasword, 1, 30*time.Second)
	if err != nil {
		t.Errorf("Unable to connect: %s", err)
		return
	}

	if socket == nil {
		t.Errorf("Socket is nil without an error")
		return
	}
	defer socket.Close()

	msg, err := socket.BgAPI("show", "api")
	if err != nil {
		t.Errorf("Unable to call API: %s", err)
		return
	}

	if msg == nil {
		t.Errorf("No message returned without an error")
		return
	}

	ct := msg.ContentType()
	if ct != ECTCommandReply {
		t.Errorf("Invalid content type, expected '%s' got '%s'", ECTCommandReply, ct)
		return
	}

	if !msg.Headers.Exists("Reply-Text") {
		t.Errorf("Expected Header of 'Reply-Text', but go headers of: %v", msg.Headers.Keys())
		return
	}

	if msg.HasError() {
		t.Errorf("Expected UUI, but got error of: %s", msg.Error())
		return
	}

	text := msg.Headers.GetString("Reply-Text")
	if !strings.HasPrefix(text, "+OK Job-UUID:") {
		t.Errorf("Expected '+OK Job-UUID:' prefix but got: %s", text)
	}
}

func TestSocketBgAPIError(t *testing.T) {
	socket := &Socket{}

	msg, err := socket.BgAPI("show", "")
	if err == nil {
		t.Errorf("Expected err, but nil returned")
	}

	if msg != nil {
		t.Errorf("Expected msg to be nil, but got '%s'", msg)
	}

}
