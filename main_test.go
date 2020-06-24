package main

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// TestHelloWorld ensures that runHelloWorldServer returns the string "Hello, world!\n"
func TestHelloWorld(t *testing.T) {
	go runHelloWorldServer()
	// Time to let the server get up and running before sending requests
	time.Sleep(10 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Error(err)
		return
	}

	got, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		t.Error(err2)
		return
	}

	const expected = "Hello, world!\n"
	if len(got) != len(expected) {
		t.Error("Bad length")
		return
	}

	for i := 0; i < len(expected); i++ {
		if got[i] != expected[i] {
			t.Error("Strings don't match")
			return
		}
	}
}
