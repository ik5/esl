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

func TestSocketRecvError(t *testing.T) {
	socket := &Socket{}

	n, buf, err := socket.Recv(10)
	if err == nil {
		t.Errorf("Expected err, but nil returned")
	}

	if n != 0 {
		t.Errorf("Expected n to be 0, but %d returned", n)
	}

	if buf != nil {
		t.Errorf("Expected buf to be nil, but got '%s'", buf)
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

func TestSocketConnectFailedPassword(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	socket, err := Connect(eslHost, eslPasword+eslHost, 1, 30*time.Second)

	if err == nil {
		t.Errorf("An error was expected, but non provided")
		if socket != nil {
			socket.Close()
		}
	}
}

func TestSocketConnectAddressError(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	socket, err := Connect(eslHost+".511", eslPasword, 1, 30*time.Second)

	if err == nil {
		t.Errorf("An error was expected, but non provided")
		if socket != nil {
			socket.Close()
		}
	}

}
