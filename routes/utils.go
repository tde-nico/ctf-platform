package routes

import (
	"fmt"
	"net/http"
	"platform/db"
	"platform/log"
	"regexp"
	"strings"
	"text/template"

	"github.com/gorilla/sessions"
)

type Flash struct {
	Message string
	Type    string
}

type Data struct {
	User        *db.User
	UserProfile *db.User
	Solves      []db.Solve
	Flashes     []Flash
}

const SEPARATOR = "|"
const MAX_PASSWORD_LENGTH = 6

var USERNAME_REGEX = regexp.MustCompile(`[0-9a-zA-Z_!@#€\-&+]{4,32}`)

var store = sessions.NewCookieStore([]byte("GrazieDarioGrazieDarioGrazieDP_1"))

func getSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, bool) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Errorf("Error getting session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, false
	}
	return session, true
}

func saveSession(w http.ResponseWriter, r *http.Request, s *sessions.Session) bool {
	err := s.Save(r, w)
	if err != nil {
		log.Errorf("Error saving session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}
	return true
}

func getSessionUser(s *sessions.Session) *db.User {
	val := s.Values["user"]
	var user = &db.User{}
	user, ok := val.(*db.User)
	if !ok {
		return nil
	}
	return user
}

func addFlash(s *sessions.Session, args ...string) {
	if len(args) < 1 || len(args) > 2 {
		return
	}
	var flashType string
	if len(args) == 1 {
		flashType = "danger"
	} else {
		flashType = args[1]
	}
	s.AddFlash(&Flash{args[0], flashType})
}

func getFlashes(w http.ResponseWriter, r *http.Request, s *sessions.Session) []Flash {
	tmp := s.Flashes()
	flashes := make([]Flash, len(tmp))
	for i, flash := range tmp {
		log.Infof("%+v", flash)
		flashes[i] = *flash.(*Flash)
	}
	err := s.Save(r, w)
	if err != nil {
		log.Errorf("Error saving session: %v", err)
	}
	return flashes
}

func getTemplate(w http.ResponseWriter, page string) (*template.Template, error) {
	tmpl, err := template.New("").ParseFiles("templates/base.html", fmt.Sprintf("templates/%s.html", page))
	if err != nil {
		log.Errorf("Error parsing template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}
	return tmpl, nil
}

func executeTemplate(w http.ResponseWriter, r *http.Request, s *sessions.Session, tmpl *template.Template, data *Data) {
	data.Flashes = getFlashes(w, r, s)

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Errorf("Error executing template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func isRegistrationAllowed(w http.ResponseWriter, r *http.Request, s *sessions.Session) bool {
	allowed := db.GetConfig("registration-allowed")
	if allowed == nil || *allowed == 0 {
		addFlash(s, "Registration is disabled")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return false
	}
	return true
}

func isUserDataValid(s *sessions.Session, username, email, password string) bool {
	if len(password) < MAX_PASSWORD_LENGTH {
		addFlash(s, "Password is too short")
		return false
	}

	if strings.ToLower(username) == "admin" {
		addFlash(s, "Username is reserved")
		return false
	}

	if !USERNAME_REGEX.MatchString(username) {
		addFlash(s, "Invalid username")
		return false
	}

	if db.UserExists(username) {
		addFlash(s, "Username already taken")
		return false
	}

	if db.EmailExists(email) {
		addFlash(s, "Email already taken")
		return false
	}

	return true
}

func loginUser(s *sessions.Session, username, password string) (*db.User, error) {
	username = strings.TrimSpace(username)
	if !USERNAME_REGEX.MatchString(username) {
		addFlash(s, "Invalid username")
		return nil, fmt.Errorf("invalid username")
	}

	user, err := db.LoginUser(username, password)
	if err != nil {
		addFlash(s, "Invalid username or password")
		return nil, err
	}

	return user, nil
}
