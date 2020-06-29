package handlers

import (
	"io"
	"net/http"

	"../storage"
)

// LogoutHandler is the http handler for the '/logout' endpoint.
type LogoutHandler struct {
	Datastore *storage.Datastore
}

func (handler LogoutHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil {
		io.WriteString(w, "You were not logged in")
		return
	}

	datastore := handler.Datastore
	if _, exists := datastore.ReadSession(cookie.Value); exists {
		datastore.DeleteSession(cookie.Value)
	}

	emptyCookie := &http.Cookie{Name: "session", Value: "", MaxAge: -1}
	http.SetCookie(w, emptyCookie)

	io.WriteString(w, "You are now logged out")
}
