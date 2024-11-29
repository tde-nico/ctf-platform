package routes

import (
	"fmt"
	"net/http"
	"platform/db"
	"platform/log"
	"strings"

	"github.com/gorilla/sessions"
)

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

func login_get(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	tmpl, err := getTemplate(w, "login")
	if err != nil {
		return
	}

	if getSessionUser(s) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	executeTemplate(w, r, s, tmpl, &Data{})
}

func login_post(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if getSessionUser(s) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := loginUser(s, username, password)
	if err != nil {
		log.Errorf("Error logging in: %v", err)
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		return
	}

	s.Values["user"] = user

	if saveSession(w, r, s) {
		http.Redirect(w, r, "/challenges", http.StatusSeeOther)
	}
}
