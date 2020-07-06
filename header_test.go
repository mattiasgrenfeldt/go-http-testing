package main

import (
	"bytes"
	"context"
	"testing"

	"./requesttesting"
	"github.com/google/go-cmp/cmp"
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

	res, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("PerformRequest() got: %v want: nil", err)
	}

	if !bytes.HasPrefix(res.Resp, []byte(statusOK)) {
		got := string(res.Resp[:bytes.IndexByte(res.Resp, '\n')+1])
		t.Errorf("status code got: %q want: %q", got, statusOK)
	}

	want := map[string][]string{"A": []string{"B"}}
	if diff := cmp.Diff(want, map[string][]string(res.Req.Header)); diff != "" {
		t.Errorf("res.Req.Header mismatch (-want +got):\n%s", diff)
	}
}

func TestOrdering(t *testing.T) {
	request := []byte("GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: X\r\n" +
		"A: Y\r\n" +
		"\r\n")

	res, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("PerformRequest() got: %v want: nil", err)
	}

	if !bytes.HasPrefix(res.Resp, []byte(statusOK)) {
		got := string(res.Resp[:bytes.IndexByte(res.Resp, '\n')+1])
		t.Errorf("status code got: %q want: %q", got, statusOK)
	}

	want := map[string][]string{"A": []string{"X", "Y"}}
	if diff := cmp.Diff(want, map[string][]string(res.Req.Header)); diff != "" {
		t.Errorf("res.Req.Header mismatch (-want +got):\n%s", diff)
	}
}

func TestContentLength(t *testing.T) {
	request := []byte("GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"Content-Length: 5\r\n" +
		"\r\n" +
		"ABCDE\r\n" +
		"\r\n")

	res, err := requesttesting.PerformRequest(context.Background(), request)
	if err != nil {
		t.Errorf("PerformRequest() got: %v want: nil", err)
	}

	if !bytes.HasPrefix(res.Resp, []byte(statusOK)) {
		got := string(res.Resp[:bytes.IndexByte(res.Resp, '\n')+1])
		t.Errorf("status code got: %q want: %q", got, statusOK)
	}

	want := map[string][]string{"Content-Length": []string{"5"}}
	if diff := cmp.Diff(want, map[string][]string(res.Req.Header)); diff != "" {
		t.Errorf("res.Req.Header mismatch (-want +got):\n%s", diff)
	}

	if res.Req.ContentLength != 5 {
		t.Errorf("res.Req.ContentLength got: %v want: 5", res.Req.ContentLength)
	}

	if string(res.ReqBody) != "ABCDE" {
		t.Errorf(`res.ReqBody got: %q want: "ABCDE"`, res.ReqBody)
	}
}
