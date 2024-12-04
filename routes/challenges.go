package routes

import (
	"net/http"
	"platform/db"
	"platform/log"

	"github.com/gorilla/sessions"
)

type DataChallenges struct {
	Data
	Challenges map[string][]*db.Challenge
	Solves     map[string]bool
}

func challenges(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	tmpl, err := getTemplate(w, "challenges")
	if err != nil {
		return
	}

	data := &DataChallenges{}
	data.User = getSessionUser(s)

	solves, err := db.GetSolvesByUser(data.User)
	if err != nil {
		log.Errorf("Error getting solves by user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Solves = make(map[string]bool)
	for _, solve := range solves {
		data.Solves[solve.ChalName] = true
	}

	data.Challenges, err = db.GetChallenges()
	if err != nil {
		log.Errorf("Error getting challenges: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// TODO: make chals decscription HTML safe

	executeTemplate(w, r, s, tmpl, data)
}
