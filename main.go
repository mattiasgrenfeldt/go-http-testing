package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// maps username to userdata
var database map[string]userData

// maps session token to username
var session map[string]string

type userData struct {
	Password string
	Note     string
}

func main() {
	fmt.Println("Listening on localhost:8080...")
	runHelloWorldServer()
}

func runHelloWorldServer() {
	database = make(map[string]userData)
	session = make(map[string]string)

	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/greetings", helloUserHandler)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/profile", profileHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func helloWorldHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
	io.WriteString(w, "You can register an account at /register\n")
	io.WriteString(w, "Or say hello to yourself at /greetings\n")
}

func helloUserHandler(w http.ResponseWriter, req *http.Request) {
	if req.FormValue("name") == "" {
		io.WriteString(w, "Please specify your name using the 'name' query parameter.\n")
	} else {
		io.WriteString(w, "Hello "+html.EscapeString(req.FormValue("name"))+"!")
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
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

		userdata, exists := database[username]
		if !exists || password != userdata.Password {
			io.WriteString(w, "User doesn't exist or the password is incorrect")
			return
		}

		token := uuid.New().String()
		session[token] = username

		cookie := &http.Cookie{Name: "session", Value: token}
		http.SetCookie(w, cookie)

		io.WriteString(w, "Logged in, proceed to /profile")
	} else {
		io.WriteString(w, "<html><body><form method=\"POST\" action=\"/login\">"+
			"<input type=\"text\" name=\"username\"placeholder=\"Username\">"+
			"<input type=\"password\" name=\"password\"placeholder=\"Password\">"+
			"<input type=\"submit\" value=\"Log in\">"+
			"</form></body></html")
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	if cookie, err := req.Cookie("session"); err == nil {
		if _, exists := session[cookie.Value]; exists {
			delete(session, cookie.Value)
		}
	}

	emptyCookie := &http.Cookie{Name: "session", Value: "", MaxAge: -1}
	http.SetCookie(w, emptyCookie)

	io.WriteString(w, "You are now logged out")
}

func registerHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		username := req.FormValue("username")
		password := req.FormValue("password")
		if username == "" || password == "" {
			io.WriteString(w, "You have to fill in a username and a password.")
			return
		}

		if _, exists := database[username]; exists {
			io.WriteString(w, "User already exists")
			return
		}

		newUserData := userData{Password: password, Note: "Buy food"}
		database[username] = newUserData

		io.WriteString(w, "User registered! Now goto /login")
	} else {
		io.WriteString(w, "<html><body><form method=\"POST\" action=\"/register\">"+
			"<input type=\"text\" name=\"username\"placeholder=\"Username\">"+
			"<input type=\"password\" name=\"password\"placeholder=\"Password\">"+
			"<input type=\"submit\" value=\"Register\">"+
			"</form></body></html")
	}
}

func profileHandler(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil {
		io.WriteString(w, "You have to log in before visiting /profile. Goto /login")
		return
	}
	username, exists := session[cookie.Value]
	if !exists {
		io.WriteString(w, "Unknown session token. HACKER DETECTED!")
		return
	}

	note := database[username].Note

	io.WriteString(w, "Welcome "+html.EscapeString(username)+". Don't forget: "+html.EscapeString(note))
}
