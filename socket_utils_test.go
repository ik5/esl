package esl

import "testing"

func TestIsHostWithPortWithoutPort(t *testing.T) {
	b, err := isHostWithPort("foo")
	if err != nil {
		t.Errorf("Expected false without errors, but got: %s", err)
		return
	}

	if b {
		t.Errorf("Expected false, but got true")
	}
}

func TestIsHostWithPortWithTooManyPorts(t *testing.T) {
	b, err := isHostWithPort("foo:1:2")
	if err == nil {
		t.Errorf("Expected false with err, but no error returned")
		return
	}

	if b {
		t.Errorf("Expected false, but got true")
	}

}

func TestSetPortWithError(t *testing.T) {
	host := "foo:1:2"
	newHost := setPort(host, "3")

	if newHost != host {
		t.Errorf("Expected '%s' but got '%s'", host, newHost)
	}
}

func TestSetPortWithPort(t *testing.T) {
	host := "foo"
	newHost := setPort(host, "1")

	if newHost != host+":1" {
		t.Errorf("Expected '%s:1' got '%s'", host, newHost)
	}
}

func TestSetPortEqual(t *testing.T) {
	host := "foo:1"
	newHost := setPort(host, "1")

	if newHost != host {
		t.Errorf("Expected '%s' got '%s'", host, newHost)
	}

}
