package handlers

import (
	"io"
	"net/http"

	"../storage"
	"github.com/google/uuid"
)

// LoginHandler is the http handler for the '/login' endpoint.
type LoginHandler struct {
	Datastore *storage.Datastore
}

func (handler LoginHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		handler.postHandler(w, req)
	} else if req.Method == http.MethodGet {
		handler.getHandler(w, req)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler LoginHandler) postHandler(w http.ResponseWriter, req *http.Request) {
	if _, err := req.Cookie("session"); err == nil {
		io.WriteString(w, "Please log out before loggin in again. Goto /logout")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	if username == "" || password == "" {
		io.WriteString(w, "You have to fill in a username and a password.")
		return
	}

	datastore := handler.Datastore
	userdata, exists := datastore.ReadDatabase(username)
	if !exists || password != userdata.Password {
		io.WriteString(w, "User doesn't exist or the password is incorrect")
		return
	}

	token := uuid.New().String()
	datastore.WriteSession(token, username)

	cookie := &http.Cookie{Name: "session", Value: token}
	http.SetCookie(w, cookie)

	io.WriteString(w, "Logged in, proceed to /profile")
}

func (LoginHandler) getHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, `<html><body><form method="POST" action="/login">
<input type="text" name="username" placeholder="Username">
<input type="password" name="password" placeholder="Password">
<input type="submit" value="Log in">
</form></body></html>`)
}
