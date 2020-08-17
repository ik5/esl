package esl

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
