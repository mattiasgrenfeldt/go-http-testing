package handlers

import (
	"html"
	"io"
	"net/http"

	"../storage"
)

// ProfileHandler is the http handler for the '/profile' endpoint
type ProfileHandler struct {
	Datastore *storage.Datastore
}

func (handler ProfileHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil {
		io.WriteString(w, "You have to log in before visiting /profile. Goto /login")
		return
	}

	datastore := handler.Datastore
	username, exists := datastore.ReadSession(cookie.Value)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "Unknown session token. HACKER DETECTED!")
		return
	}

	userdata, exists := datastore.ReadDatabase(username)
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	note := userdata.Note

	io.WriteString(w, "Welcome "+html.EscapeString(username)+". Don't forget: "+html.EscapeString(note))
}
