package routes

import (
	"fmt"
	"net/http"
	"platform/db"
	"platform/log"
	"regexp"
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

const MAX_PASSWORD_LENGTH = 6

var USERNAME_REGEX = regexp.MustCompile(`[0-9a-zA-Z_!@#€\-&+]{4,32}`)

// ! TODO: change
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
		flashes[i] = *flash.(*Flash)
	}
	err := s.Save(r, w)
	if err != nil {
		log.Errorf("Error saving session: %v", err)
	}
	return flashes
}

func getTemplate(w http.ResponseWriter, page string) (*template.Template, error) {
	// TODO: parse all templates at startup
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
