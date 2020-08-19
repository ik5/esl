package esl

import "errors"

// Default settings
const (
	DefaultPort                 = "8021"
	EOL                         = "\r\n"
	MaxBufferSize         int64 = 2_000_000
	AuthRequestBufferSize int64 = 32
)

// Error instances that are used and can be validated using errors.Is
var (
	ErrConnectionIsNotInitialized   = errors.New("Connection is not initialized")
	ErrCmdEOL                       = errors.New("cmd contains EOL")
	ErrContentLengthZero            = errors.New("Content Length is zero")
	ErrUnableToGetConnectedSocket   = errors.New("Unable to get connected socket")
	ErrUnableToLogInNoErrorReturned = errors.New("Unable to log in, no error returned")
)
