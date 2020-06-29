package server

import (
	"context"
	"html"
	"io"
	"log"
	"net/http"

	"./handlers"
	"./storage"
)

// Server is a basic http server.
type Server struct {
	server http.Server
}

// New creates a new Server using the provided datastore.
func New(datastore *storage.Datastore) Server {
	handler := SetupHandlers(datastore)
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: handler}

	return Server{server: server}
}

// SetupHandlers sets up all of the HTTP handlers for different endpoints.
func SetupHandlers(datastore *storage.Datastore) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/greetings", greetingsHandler)

	mux.Handle("/login", handlers.LoginHandler{Datastore: datastore})
	mux.Handle("/logout", handlers.LogoutHandler{Datastore: datastore})
	mux.Handle("/register", handlers.RegisterHandler{Datastore: datastore})
	mux.Handle("/profile", handlers.ProfileHandler{Datastore: datastore})

	return mux
}

// Start starts the Server.
func (server *Server) Start() {
	go func() {
		log.Fatal(server.server.ListenAndServe())
	}()
}

// Stop stops the Server.
func (server *Server) Stop() {
	server.server.Shutdown(context.Background())
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, `Hello, world!
You can register an account at /register
Or say hello to yourself at /greetings`)
}

func greetingsHandler(w http.ResponseWriter, req *http.Request) {
	if req.FormValue("name") == "" {
		io.WriteString(w, "Please specify your name using the 'name' query parameter.\n")
		return
	}
	io.WriteString(w, "Hello "+html.EscapeString(req.FormValue("name"))+"!")
}
