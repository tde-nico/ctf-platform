package routes

import (
	"net/http"
	"platform/db"
	"platform/log"

	"github.com/gorilla/sessions"
)

func newpw_get(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	tmpl, err := getTemplate(w, "newpw")
	if err != nil {
		return
	}
	user := getSessionUser(s)

	data := Data{User: user}

	executeTemplate(w, r, s, tmpl, &data)
}

func newpw_post(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	user := getSessionUser(s)
	if user == nil {
		addFlash(s, "You must be logged in to access that page")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Errorf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	log.Infof("form: %+v", r.Form)

	newpw := r.FormValue("password")

	if len(newpw) < MAX_PASSWORD_LENGTH {
		addFlash(s, "Password is too short")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/newpw", http.StatusSeeOther)
		}
		return
	}

	err = db.ChangePassword(user.Username, newpw, false)
	if err != nil {
		log.Errorf("Error changing password: %v", err)
		addFlash(s, "Error changing password")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/newpw", http.StatusSeeOther)
		}
		return
	}

	s.Values["apikey"] = ""

	addFlash(s, "Password changed successfully", "success")
	if saveSession(w, r, s) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
