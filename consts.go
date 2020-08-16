package esl

import "errors"

// Default settings
const (
	DefaultPort                 = "8021"
	EOL                         = "\r\n"
	MaxBufferSize         int64 = 2_000_000
	AuthRequestBufferSize int64 = 32
)

var (
	CmdEOLError            = errors.New("cmd contains EOL")
	ContentLengthZeroError = errors.New("Content Length is zero")
)
