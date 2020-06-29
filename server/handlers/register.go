package handlers

import (
	"io"
	"net/http"

	"../storage"
)

// RegisterHandler is the http handler for the '/register' endpoint
type RegisterHandler struct {
	Datastore *storage.Datastore
}

func (handler RegisterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		handler.postHandler(w, req)
	} else if req.Method == http.MethodGet {
		handler.getHandler(w, req)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler RegisterHandler) postHandler(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")
	if username == "" || password == "" {
		io.WriteString(w, "You have to fill in a username and a password.")
		return
	}

	datastore := handler.Datastore
	if _, exists := datastore.ReadDatabase(username); exists {
		io.WriteString(w, "User already exists")
		return
	}

	newUserData := storage.UserData{Password: password, Note: "Buy food"}
	datastore.WriteDatabase(username, newUserData)

	io.WriteString(w, "User registered! Now goto /login")
}

func (RegisterHandler) getHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, `<html><body><form method="POST" action="/register">
<input type="text" name="username" placeholder="Username">
<input type="password" name="password" placeholder="Password">
<input type="submit" value="Register">
</form></body></html>`)
}
