package esl

import "errors"

// Default settings
const (
	DefaultPort                 = "8021"
	EOL                         = "\r\n"
	BufferSize            int64 = 2048
	AuthRequestBufferSize int64 = 32
)

var (
	CmdEOLError = errors.New("cmd contains EOL")
)
