package routes

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func home(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	tmpl, err := getTemplate(w, "home")
	if err != nil {
		return
	}

	if r.RequestURI != "/" {
		addFlash(s, "404 Page not found")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	user := getSessionUser(s)

	data := Data{User: user}

	executeTemplate(w, r, s, tmpl, &data)
}
