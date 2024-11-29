package routes

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func logout(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	s.Options.MaxAge = -1
	if saveSession(w, r, s) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
