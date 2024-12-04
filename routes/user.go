package routes

import (
	"net/http"
	"platform/db"
	"platform/log"

	"github.com/gorilla/sessions"
)

type DataUserInfo struct {
	Data
	UserProfile *db.User
	Solves      []*db.Solve
}

func userInfo(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	tmpl, err := getTemplate(w, "user")
	if err != nil {
		return
	}

	data := &DataUserInfo{}
	data.User = getSessionUser(s)

	username := r.PathValue("username")
	if username == "" {
		addFlash(s, "Invalid username")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	data.UserProfile, err = db.GetUserByUsername(username)
	if err != nil {
		log.Errorf("Error getting user by username: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Solves, err = db.GetSolvesByUser(data.UserProfile)
	if err != nil {
		log.Errorf("Error getting solves by user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	executeTemplate(w, r, s, tmpl, data)
}
