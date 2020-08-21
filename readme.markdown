# ESL

The following package create support for [Freeswitch](https://www.freeswitch.org) [ESL](https://freeswitch.org/confluence/display/FREESWITCH/Event+Socket+Library) in Go Language.

## Why ?

There already good implementation for ESL, but on some cases it didn't work well for me.

The way I need a support for ESL:

  * I need to separate events and actions.
  * Ability to connect using backoff to find the best way to connect if it fails.
  * Easy to debug, log agnostic way.
  * Raw vs Parsed support.
  * Minimal dependencies as possible.
  * [mod_commands](https://freeswitch.org/confluence/display/FREESWITCH/mod_commands) (and ESL interface) support.
  * Support for [dptools](https://freeswitch.org/confluence/display/FREESWITCH/mod_dptools) using API's ["execute"](https://freeswitch.org/confluence/display/FREESWITCH/mod_event_socket#mod_event_socket-3.9.1.1execute) command.
  * Test driven development (most of the time) for all my commands.
  * Testing over my [Freeswitch's docker](https://github.com/ik5/freeswitch-docker) repo.

## Work in progress

At the moment the package is work in progress that will add more support, when I
need it for my projects.

# How it works?

```go
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
```

# TODO:

[ ] Add debug support using callbacks.
[ ] Finish interface support.
[ ] Work on supporting events (Dual connection commands and for events).
[ ] Parse events
[ ] Work on supporting callbacks for registered events.
[ ] Examples
[ ] Better documentation

# License

The following project release under Apache license Version 2.0
