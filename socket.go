package esl

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Socket is low level ESL connection.
// Socket will generate keep-alive for a connection, to keep it open in order for
// a single connection will not be dropped after sending/receiving a payload.
type Socket struct {
	conn       *net.TCPConn
	host       string
	password   string
	maxRetries uint64
	timeout    time.Duration
	ctx        *context.Context
	loggedin   bool
	reader     *bufio.Reader
	writer     *bufio.Writer
	lock       *sync.RWMutex
}

// Dial open an new connection for Freeswitch, with retries until it maxRetries
// is due.
// If host does not contain port (e.g. freeswitch.example.com:8021), the default
// port will be assigned (8021).
// password is a clear text password that is sent to the ESL auth request.
// timeout is the amount of waiting until dialing to ESL will fail if no answer was provided.
//
// If maxRetries is 0, it will not retry if failed.
// The retry is using Backoff algorithm.
func Dial(host string, password string, maxRetries uint64, timeout time.Duration) (*Socket, error) {
	socket := Socket{
		host:       setPort(host, DefaultPort),
		password:   password,
		maxRetries: maxRetries,
		timeout:    timeout,
		lock:       &sync.RWMutex{},
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	socket.ctx = &ctx

	var conn *net.TCPConn
	var remoteAddr *net.TCPAddr
	var err error

	remoteAddr, err = net.ResolveTCPAddr("tcp", socket.host)
	if err != nil {
		return nil, err
	}

	bo := backoff.WithContext(
		backoff.WithMaxRetries(
			backoff.NewExponentialBackOff(), maxRetries,
		), *socket.ctx)

	err = backoff.Retry(func() error {
		// conn, err = net.DialTimeout("tcp", socket.host, socket.timeout)
		conn, err = net.DialTCP("tcp", nil, remoteAddr)
		if err == nil {
			socket.conn = conn
		}
		return err
	}, bo)

	if err != nil {
		return nil, err
	}

	socket.reader = bufio.NewReader(conn)
	socket.writer = bufio.NewWriter(conn)
	// make sure the connection stay open if possible
	socket.conn.SetKeepAlive(true)
	socket.conn.SetKeepAlivePeriod(timeout)

	return &socket, nil
}

// Connect Connect to ESL and does a login.
// If an error occurs, it will disconnect and return an error
func Connect(host string, password string, maxRetries uint64, timeout time.Duration) (*Socket, error) {
	socket, err := Dial(host, password, maxRetries, timeout)
	if err != nil {
		return nil, err
	}

	if socket == nil {
		return nil, ErrUnableToGetConnectedSocket
	}

	loggedIn, err := socket.Login()
	if err != nil {
		socket.Close()
		return nil, err
	}

	if !loggedIn {
		socket.Close()
		return nil, ErrUnableToLogInNoErrorReturned
	}

	return socket, nil
}

// Close a connection
func (s Socket) Close() error {
	err := s.writer.Flush()
	if err != nil {
		return err
	}

	err = s.conn.SetKeepAlive(false)
	if err != nil {
		return err
	}

	err = s.conn.CloseRead()
	if err != nil {
		return err
	}

	err = s.conn.CloseWrite()
	if err != nil {
		return err
	}
	return s.conn.Close()
}

// Send a request to ESL.
// If cmd contains EOL
func (s Socket) Send(cmd string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if strings.HasSuffix(cmd, EOL) {
		return ErrCmdEOL
	}

	buf := cmd + EOL + EOL
	l := len(buf)

	n, err := s.writer.WriteString(buf)
	if err != nil {
		return err
	}
	defer s.writer.Flush()

	if n < l && s.writer.Buffered() == 0 {
		return fmt.Errorf("Wrote %d bytes, expected %d", l, n)
	}

	return nil
}

// Recv a content from the server
func (s Socket) Recv(maxBuff int64) (int, []byte, error) {
	buf := make([]byte, maxBuff)
	n, err := s.reader.Read(buf)
	return n, buf, err
}

// Login into the ESL server
func (s *Socket) Login() (bool, error) {
	if s.loggedin {
		return true, nil
	}

	n, content, err := s.Recv(AuthRequestBufferSize)
	if err != nil {
		return false, err
	}
	if int64(n) >= AuthRequestBufferSize {
		return false, fmt.Errorf("Auth length %d is too big", n)
	}
	auth, err := NewMessage(content, true)
	if err != nil {
		return false, err
	}

	contentType := auth.Headers.GetString("Content-Type")
	if contentType != "auth/request" {
		return false, fmt.Errorf("Invalid Content-Type: %s", contentType)
	}

	n, content, err = s.SendRecv("auth " + s.password)
	if err != nil {
		return false, fmt.Errorf("Unable to send/recv auth: %s", err)
	}

	if int64(n) <= AuthRequestBufferSize {
		return false, fmt.Errorf("Invalid msg length: %d for %s", n, content)
	}

	msg, err := NewMessage(content, true)
	if err != nil {
		return false, err
	}

	if msg.HasError() {
		return false, fmt.Errorf("Login error: %s", msg.Error())
	}

	headers := msg.Headers

	answer := headers.GetString("Reply-Text")
	loggedIn := strings.Compare("+OK accepted", answer)

	s.loggedin = loggedIn == 0

	return s.loggedin, nil
}

// LoggedIn is true if a login was made successfully
func (s *Socket) LoggedIn() bool {
	return s.loggedin
}
