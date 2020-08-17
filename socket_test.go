package esl

import (
	"bytes"
	"errors"
	"os"
	"strings"
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

func TestBasicConnectionSendCmdEOL(t *testing.T) {
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

	err = socket.Send(EOL)
	if err == nil {
		t.Errorf("No Error was returned")
		return
	}

	if !errors.Is(err, ErrCmdEOL) {
		t.Errorf("Unexpected error: %s", err)
		return
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

func TestDoubleAuthentication(t *testing.T) {
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

	loggedIn, err = socket.Login()
	if err != nil {
		t.Errorf("Login error: %s", err)
		return
	}

	if !loggedIn {
		t.Errorf("Expected logged in to be true")
	}
}

func TestAuthenticationBadCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	// Place bad login to test it up
	socket, err := Dial(eslHost, eslPasword+eslHost, 1, 30*time.Second)
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
	if err == nil {
		t.Errorf("Expected error, but was not received one")
		return
	}

	if loggedIn {
		t.Errorf("Expected not to be logged in")
		return
	}

	strErr := err.Error()

	if strings.HasPrefix("Unable to send/recv auth:", strErr) {
		t.Errorf("Error sending connection: %s", strErr)
		return
	}

	if strErr != "Login error: invalid" {
		t.Errorf("Unexpected error: %s", strErr)
		return
	}

}

func TestSocketConnect(t *testing.T) {
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
	socket.Close()
}
