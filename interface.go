package esl

import "fmt"

// Current file contains the implementation of ESL.i interface.
// the file is part of sockets.go but focuses only on the interface

// SendRecv sends a content and wait for returns an answer
func (s Socket) SendRecv(cmd string) (int, []byte, error) {
	err := s.Send(cmd)
	if err != nil {
		return 0, nil, err
	}

	return s.Recv(MaxBufferSize)
}

// API sends the api commands
func (s Socket) API(cmd string, args string) (*Message, error) {
	_, buffer, err := s.SendRecv(fmt.Sprintf("api %s %s", cmd, args))

	if err != nil {
		return nil, err
	}

	message, err := NewMessage(buffer, true)
	return message, err
}

// BgAPI sends the bgapi commands
func (s Socket) BgAPI(cmd string, args string) (*Message, error) {
	_, buffer, err := s.SendRecv(fmt.Sprintf("bgapi %s %s", cmd, args))

	if err != nil {
		return nil, err
	}

	message, err := NewMessage(buffer, true)
	return message, err
}
