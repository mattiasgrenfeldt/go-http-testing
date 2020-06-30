package main

import (
	"context"
	"testing"

	"./requesttesting"
)

func TestBasic(t *testing.T) {
	// All tests use verbatim newline characters instead of using multiline strings to ensure that \r and \n end up in exactly the right places.
	request := "GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: B\r\n" +
		"\r\n"

	req, resp, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("got: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("got: %v wanted: 200 OK", resp.Status)
	}

	headers := req.Header
	if len(headers["A"]) != 1 || headers["A"][0] != "B" {
		t.Errorf("got:, %v, wanted: map[A:[B]]", headers)
	}
}

func TestOrdering(t *testing.T) {
	request := "GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: X\r\n" +
		"A: Y\r\n" +
		"\r\n"

	req, resp, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("got: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("got: %v wanted: 200 OK", resp.Status)
	}

	headers := req.Header
	if len(headers["A"]) != 2 || headers["A"][0] != "X" || headers["A"][1] != "Y" {
		t.Errorf("got: %v, wanted: map[A:[X Y]]", headers)
	}
}
