package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ik5/esl"
)

const (
	maxRetries uint64        = 24
	timeout    time.Duration = 45 * time.Second
)

func main() {
	host := os.Getenv("ESLHOST")
	password := os.Getenv("ESLPASSWORD")

	// Connect and login
	socket, err := esl.Connect(host, password, maxRetries, timeout)
	if err != nil {
		panic(err)
	}

	defer socket.Close()

	msg, err := socket.API("echo", "Hello World")
	if err != nil {
		panic(err)
	}

	// Should print:
	//   Answer: Headers: Content-Type=api/response | Content-Length=11 ; Body: Hello World
	fmt.Printf("Answer: %s\n", msg)
}
