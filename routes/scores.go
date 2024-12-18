package routes

import (
	"net/http"
	"platform/db"
	"platform/log"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
)

func submit(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	challID := r.FormValue("challID")
	flag := strings.TrimSpace(r.FormValue("flag"))

	user := getSessionUser(s)
	if user == nil {
		addFlash(s, "You must be logged in to submit flags")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	chalID, err := strconv.Atoi(challID)
	if err != nil {
		log.Errorf("Error converting chalID to int: %v", err)
		http.Error(w, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	status, err := db.SubmitFlag(user, chalID, flag)
	if err != nil {
		log.Errorf("Error submitting flag: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	switch status {
	case db.StatusCorrectFlag:
		addFlash(s, "Correct flag :)", "success")
	case db.StatusAlreadySolved:
		addFlash(s, "Challenge already solved", "warning")
	case db.StatusWrongFlag:
		addFlash(s, "Wrong flag :(")
	}

	if saveSession(w, r, s) {
		http.Redirect(w, r, "/challenges", http.StatusSeeOther)
	}
}

type ScoresData struct {
	Data
	Users []db.UserScore
}

func scores(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	tmpl, err := getTemplate(w, "scores")
	if err != nil {
		return
	}

	users, err := db.GetScoreUsers()
	if err != nil {
		log.Errorf("Error getting users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := ScoresData{Data: Data{}}
	data.Users = users
	data.User = getSessionUser(s)

	executeTemplate(w, r, s, tmpl, &data)
}

func graphData(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	log.Infof("graphData")
}
