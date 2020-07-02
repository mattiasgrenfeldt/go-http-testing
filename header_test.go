package main

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"./requesttesting"
)

const (
	statusOK = "HTTP/1.1 200 OK\r\n"
)

func checkStatusMessage(t *testing.T, response []byte, statusMessage string) {
	if !bytes.HasPrefix(response, []byte(statusMessage)) {
		got := string(response[:bytes.IndexByte(response, '\n')+1])
		t.Errorf("status code got: %q want: %q", got, statusMessage)
	}
}

func checkEqualHeaders(t *testing.T, got http.Header, want http.Header) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("req.Header mismatch (-want +got):\n%s", diff)
	}
}

func TestBasic(t *testing.T) {
	// All tests use verbatim newline characters instead of using multiline strings to ensure that \r and \n end up in exactly the right places.
	request := []byte("GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: B\r\n" +
		"\r\n")

	req, _, resp, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("PerformRequest() got: %v want: nil", err)
	}

	checkStatusMessage(t, resp, statusOK)

	want := http.Header{"A": []string{"B"}}
	checkEqualHeaders(t, req.Header, want)
}

func TestOrdering(t *testing.T) {
	request := []byte("GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: X\r\n" +
		"A: Y\r\n" +
		"\r\n")

	req, _, resp, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("PerformRequest() got: %v want: nil", err)
	}

	checkStatusMessage(t, resp, statusOK)

	want := http.Header{"A": []string{"X", "Y"}}
	checkEqualHeaders(t, req.Header, want)
}

func TestContentLength(t *testing.T) {
	request := []byte("GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"Content-Length: 5\r\n" +
		"\r\n" +
		"ABCDE\r\n" +
		"\r\n")

	req, reqBody, resp, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("PerformRequest() got: %v want: nil", err)
	}

	checkStatusMessage(t, resp, statusOK)

	want := http.Header{"Content-Length": []string{"5"}}
	checkEqualHeaders(t, req.Header, want)

	if req.ContentLength != 5 {
		t.Errorf("req.ContentLength got: %v want: 5", req.ContentLength)
	}

	if string(reqBody) != "ABCDE" {
		t.Errorf("req.Body got: %q want: \"ABCDE\"", reqBody)
	}
}
