package main

import (
	"context"
	"testing"

	"./requesttesting"
)

const (
	statusOK = "HTTP/1.1 200 OK"
)

// hasPrefix checks whether 'bytes' has the prefix 'prefix'.
func hasPrefix(bytes []byte, prefix string) bool {
	prefixSlice := []byte(prefix)
	for i := range prefixSlice {
		if bytes[i] != prefixSlice[i] {
			return false
		}
	}
	return true
}

func TestBasic(t *testing.T) {
	// All tests use verbatim newline characters instead of using multiline strings to ensure that \r and \n end up in exactly the right places.
	request := []byte("GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: B\r\n" +
		"\r\n")

	req, resp, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("err, got: %v, wanted: nil", err)
	}

	if !hasPrefix(resp, statusOK) {
		t.Errorf("Status code, got: %v, wanted: %v", resp[:len(statusOK)], statusOK)
	}

	headers := req.Header
	if len(headers["A"]) != 1 || headers["A"][0] != "B" {
		t.Errorf("req.Header, got:, %v, wanted: map[A:[B]]", headers)
	}
}

func TestOrdering(t *testing.T) {
	request := []byte("GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: X\r\n" +
		"A: Y\r\n" +
		"\r\n")

	req, resp, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("err, got: %v, wanted: nil", err)
	}

	if !hasPrefix(resp, statusOK) {
		t.Errorf("Status code, got: %v, wanted: %v", resp[:len(statusOK)], statusOK)
	}

	headers := req.Header
	if len(headers["A"]) != 2 || headers["A"][0] != "X" || headers["A"][1] != "Y" {
		t.Errorf("req.Header, got: %v, wanted: map[A:[X Y]]", headers)
	}
}
