package routes

import (
	"net/http"
	"platform/db"
	"platform/log"

	"github.com/gorilla/sessions"
)

type AdminPanelData struct {
	Data
	Users        []*db.User
	Categories   []string
	Difficulties []string
	Challenges   map[string][]*db.Challenge
	Submissions  []*db.Submission
	Config       []*db.Config
}

func admin(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	tmpl, err := getTemplate(w, "admin")
	if err != nil {
		return
	}

	data := &AdminPanelData{}

	data.User = getSessionUser(s)

	data.Users, err = db.GetUsers()
	if err != nil {
		log.Errorf("Error getting users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Categories = db.CATEGORIES
	data.Difficulties = db.DIFFICULTIES

	data.Challenges, err = db.GetChallenges()
	if err != nil {
		log.Errorf("Error getting challenges: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Submissions, err = db.GetSubmissions()
	if err != nil {
		log.Errorf("Error getting submissions: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Config, err = db.GetConfigs()
	if err != nil {
		log.Errorf("Error getting configs: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Flashes = getFlashes(w, r, s)

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Errorf("Error executing template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func adminNewChall(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Errorf("Error parsing multipart form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	chal, err := getChallFromForm(w, r)
	if err != nil {
		return
	}

	if len(r.MultipartForm.File) > 0 {
		err = extractChallengeFiles(w, r, s, chal)
		if err != nil {
			return
		}
	}

	createChallenge(w, r, s, chal)
}

func adminUpdateChall(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Errorf("Error parsing multipart form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	chal, err := getChallFromForm(w, r)
	if err != nil {
		return
	}

	err = isChallengeValid(chal)
	if err != nil {
		addFlash(s, err.Error())
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		return
	}

	err = renameChallenge(chal)
	if err != nil {
		log.Errorf("Error renaming challenge: %v", err)
		addFlash(s, "Error renaming challenge")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		return
	}

	if len(r.MultipartForm.File) > 0 {
		err = deleteChallengeFiles(chal.Name)
		if err != nil {
			log.Errorf("Error deleting challenge files: %s: %v", chal.Name, err)
			addFlash(s, "Error deleting challenge")
			if saveSession(w, r, s) {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
			}
			return
		}

		err = extractChallengeFiles(w, r, s, chal)
		if err != nil {
			return
		}
	}

	err = db.UpdateChallenge(chal)
	if err != nil {
		log.Errorf("Error updating challenge: %v", err)
		addFlash(s, "Error updating challenge")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		return
	}

	addFlash(s, "Challenge Updated Successfully", "success")
	if saveSession(w, r, s) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func adminDeleteChall(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf("Error parsing form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	log.Infof("Delete Challenge: %s", name)
	err = deleteChallenge(w, r, s, name)
	if err != nil {
		return
	}

	addFlash(s, "Challenge Deleted Successfully", "success")
	if saveSession(w, r, s) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

// func adminDeleteChall(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
// 	log.Infof("adminDeleteChall")
// }

func adminResetPw(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	log.Infof("adminResetPw")
}

func adminConfig(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf("Error parsing form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	configs, err := db.GetConfigs()
	if err != nil {
		log.Errorf("Error getting configs: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, config := range configs {
		value := r.FormValue(config.Key)
		if value == "" {
			err = db.SetConfig(config.Key, "0")
		} else {
			err = db.SetConfig(config.Key, value)
		}
		if err != nil {
			log.Errorf("Error setting config %s: %v", config.Key, err)
			addFlash(s, "Error updating configuration", "danger")
			if saveSession(w, r, s) {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
			}
			return
		}
	}

	addFlash(s, "Configuration updated successfully", "success")
	if saveSession(w, r, s) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
