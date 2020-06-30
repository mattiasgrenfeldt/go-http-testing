package main

import (
	"context"
	"io"
	"net/http"
	"testing"
)

type ToyHandler struct {
	Req  *http.Request
	Done chan bool
}

func (h *ToyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Req = r
	io.WriteString(w, "Cool")
	h.Done <- true
}

func performRequest(request string) (*http.Request, error) {
	done := make(chan bool, 1)
	handler := &ToyHandler{
		Req:  nil,
		Done: done,
	}

	srv := http.Server{Handler: handler}
	listener := NewToyListener(request)

	var err error
	go func() {
		err = srv.Serve(&listener)
		if err != http.ErrServerClosed {
			done <- true
		}
	}()
	<-done
	srv.Shutdown(context.Background())

	if err != http.ErrServerClosed {
		return handler.Req, err
	}

	return handler.Req, nil
}

func TestBasic(t *testing.T) {
	request := "GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: B\r\n" +
		"\r\n"

	req, err := performRequest(request)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	headers := req.Header
	if len(headers["A"]) != 1 || headers["A"][0] != "B" {
		t.Errorf("got:, %v", headers)
	}
}

func TestOrdering(t *testing.T) {
	request := "GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: X\r\n" +
		"A: Y\r\n" +
		"\r\n"

	req, err := performRequest(request)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	headers := req.Header
	if len(headers["A"]) != 2 || headers["A"][0] != "X" || headers["A"][1] != "Y" {
		t.Errorf("got:, %v", headers)
	}
}
