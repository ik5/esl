package esl

import (
	"errors"
	"time"
)

// Event hold information regarding
type Event struct {
	Message chan *Message
}

// ESL is the structure for all type of events and commands
type ESL struct {
	socket *Socket
	funcs  map[string][]func(Event)

	loggedIn bool
}

// NewESL create a new ESL, and does a login
func NewESL(host string, password string, maxRetries uint64, timeout time.Duration) (*ESL, error) {
	esl := ESL{}
	socket, err := Dial(host, password, maxRetries, timeout)
	if err != nil {
		return nil, err
	}

	esl.socket = socket

	if !esl.socket.LoggedIn() {
		loggedIn, err := esl.socket.Login()
		if err != nil {
			return nil, err
		}

		if !loggedIn {
			return nil, errors.New("Unable to loggin, but no error")
		}
	}

	return &esl, nil
}
