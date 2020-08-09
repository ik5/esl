package esl

import (
	"bytes"
	"os"
	"testing"
	"time"
)

var (
	eslHost    = os.Getenv("ESLHOST")
	eslPasword = os.Getenv("ESLPASSWORD")
)

func TestBasicConnection(t *testing.T) {
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

	err = socket.Close()
	if err != nil {
		t.Errorf("Unable to close connection: %s", err)
	}
}

func TestBasicConnectionRecv(t *testing.T) {
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

	err = socket.Close()
	if err != nil {
		t.Errorf("Unable to close connection: %s", err)
	}
}

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

func TestBasicAuthentication(t *testing.T) {
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

	loggedIn, err := socket.Login()
	if err != nil {
		t.Errorf("Login error: %s", err)
		return
	}

	if !loggedIn {
		t.Errorf("Not logged in based on return")
		return
	}

	if !socket.LoggedIn() {
		t.Errorf("Not marked as logged in")
		return
	}
}
