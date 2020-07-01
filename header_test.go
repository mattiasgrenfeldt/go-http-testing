package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"./requesttesting"
)

const (
	statusOK = "HTTP/1.1 200 OK\r\n"
)

func TestBasic(t *testing.T) {
	// All tests use verbatim newline characters instead of using multiline strings to ensure that \r and \n end up in exactly the right places.
	request := []byte("GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: B\r\n" +
		"\r\n")

	req, resp, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("PerformRequest() got: %v want: nil", err)
	}

	if !bytes.HasPrefix(resp, []byte(statusOK)) {
		got := string(resp[:bytes.IndexByte(resp, '\n')+1])
		t.Errorf("status code got: %q want: %q", got, statusOK)
	}

	headers := req.Header
	want := []string{"B"}
	if diff := cmp.Diff(headers["A"], want); diff != "" {
		t.Errorf("req.Header[\"A\"] got: %v want: %v\ndiff: %v", headers["A"], want, diff)
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
		t.Errorf("PerformRequest() got: %v want: nil", err)
	}

	if !bytes.HasPrefix(resp, []byte(statusOK)) {
		got := string(resp[:bytes.IndexByte(resp, '\n')+1])
		t.Errorf("status code got: %q want: %q", got, statusOK)
	}

	headers := req.Header
	want := []string{"X", "Y"}
	if diff := cmp.Diff(headers["A"], want); diff != "" {
		t.Errorf("req.Header[\"A\"] got: %v want: %v\ndiff: %v", headers["A"], want, diff)
	}
}
