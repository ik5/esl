package esl

import (
	"fmt"
	"strings"
)

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
	_, msg, err := s.SendCommands("api", cmd, args)

	if err != nil {
		return nil, err
	}

	return msg, err
}

// BgAPI sends the bgapi commands
func (s Socket) BgAPI(cmd string, args string) (*Message, error) {
	_, msg, err := s.SendCommands("bgapi", cmd, args)

	if err != nil {
		return nil, err
	}

	return msg, err
}

// Filter supports the simple filter
func (s Socket) Filter(eventName, valueToFilter string) (*Message, error) {
	_, msg, err := s.SendCommands("filter", eventName, valueToFilter)
	return msg, err
}

// FilterWithOutput execute filter command with output type (plain -
// default, XML and JSON)
func (s Socket) FilterWithOutput(outputType EventOutputType, eventName, valueToFilter string) (*Message, error) {
	_, msg, err := s.SendCommands("filter", string(outputType), fmt.Sprintf("%s %s", eventName, valueToFilter))
	return msg, err
}

// FilterDelete Specify the events which you want to revoke the filter.
// filter delete can be used when some filters are applied wrongly or when
// there is no use of the filter.
func (s Socket) FilterDelete(eventName, valueToFilter string) (*Message, error) {
	_, msg, err := s.SendCommands("filter", "delete", fmt.Sprintf("%s %s", eventName, valueToFilter))
	return msg, err
}

// SendEvent Send an event into the event system (multi line input for headers).
func (s Socket) SendEvent(eventName string, headers Headers, body string) (*Message, error) {

	var hdrs strings.Builder
	for _, header := range headers.Keys() {
		hdrs.WriteString(header + ": ")
		hdrs.WriteString(headers.GetString(header))
		hdrs.WriteString(EOL)
	}

	// Add content length if ther is a body, but no content length, so the
	// action will work well.
	if body != "" && !headers.Exists("Content-Length") {
		hdrs.WriteString(fmt.Sprintf("Content-Length: %d%s", len(body), EOL))
	}

	toSend := fmt.Sprintf("%s%s", EOL, hdrs)

	if body != "" {
		toSend = fmt.Sprintf("%s%s%s", toSend, EOL, body)
	}

	_, msg, err := s.SendCommands("sendevent", eventName+EOL, toSend)
	return msg, err
}
