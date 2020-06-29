package main

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"./server"
	"./server/storage"
)

func TestIndex(t *testing.T) {
	datastore := storage.New()
	handler := server.SetupHandlers(&datastore)

	req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	resp := w.Result()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	got := string(body)
	expected := "Hello, world!\nYou can register an account at /register\nOr say hello to yourself at /greetings"
	if string(got) != expected {
		t.Errorf("Failed: got '%s', wanted '%s'.", got, expected)
	}
}

func TestGreetingNoName(t *testing.T) {
	datastore := storage.New()
	handler := server.SetupHandlers(&datastore)

	req := httptest.NewRequest("GET", "http://localhost:8080/greetings", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	resp := w.Result()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	got := string(body)
	if !strings.Contains(got, "query parameter") {
		t.Errorf("Failed: got '%s'.", got)
	}
}
func TestGreetingWithName(t *testing.T) {
	datastore := storage.New()
	handler := server.SetupHandlers(&datastore)

	req := httptest.NewRequest("GET", "http://localhost:8080/greetings?name=Pelle", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	resp := w.Result()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	got := string(body)
	if !strings.Contains(got, "Pelle") {
		t.Errorf("Failed: got '%s'.", got)
	}
}
