package esl

import (
	"fmt"
	"strings"
)

// isHostWithPort try to look for name:port in host.
// if found, it returns true, if not, it will return false
func isHostWithPort(host string) (bool, error) {
	fragments := strings.Split(host, ":")
	l := len(fragments)
	if l < 2 {
		return false, nil
	}

	if l > 2 {
		return false, fmt.Errorf("Invalid host structure, expected host:port, found: %s", host)
	}

	return fragments[1] != "", nil
}

// setPort add port to host, if not found
func setPort(host string, port string) string {
	newHost := host
	b, err := isHostWithPort(host)
	if err != nil {
		return newHost
	}

	if !b {
		newHost += ":" + port
	}

	return newHost
}
