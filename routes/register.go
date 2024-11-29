package routes

import (
	"net/http"
	"platform/db"
	"platform/log"
	"strings"

	"github.com/gorilla/sessions"
)

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

func isRegistrationAllowed(w http.ResponseWriter, r *http.Request, s *sessions.Session) bool {
	allowed := db.GetConfig("registration-allowed")
	if allowed == nil || *allowed == 0 {
		addFlash(s, "Registration is disabled")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return false
	}
	return true
}

func register_get(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if getSessionUser(s) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !isRegistrationAllowed(w, r, s) {
		return
	}

	tmpl, err := getTemplate(w, "register")
	if err != nil {
		return
	}

	executeTemplate(w, r, s, tmpl, &Data{})
}

func register_post(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	if getSessionUser(s) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !isRegistrationAllowed(w, r, s) {
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if !isUserDataValid(s, username, email, password) {
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
		}
		return
	}

	err := db.RegisterUser(username, email, password)
	if err != nil {
		log.Errorf("Error registering user %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	addFlash(s, "Registration Completed", "success")
	if saveSession(w, r, s) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
