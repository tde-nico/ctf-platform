package routes

import (
	"encoding/json"
	"net/http"
	"platform/db"
	"platform/log"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
)

type GraphPoint struct {
	X string `json:"x"`
	Y int    `json:"y"`
}

func submit(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	challID := r.FormValue("challID")
	flag := strings.TrimSpace(r.FormValue("flag"))

	user := getSessionUser(s)
	if user == nil {
		addFlash(s, "You must be logged in to submit flags")
		if saveSession(w, r, s) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		return
	}

	chalID, err := strconv.Atoi(challID)
	if err != nil {
		log.Errorf("Error converting challID to int: %v", err)
		http.Error(w, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	status, err := db.SubmitFlag(user, chalID, flag)
	if err != nil {
		log.Errorf("Error submitting flag: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	header := http.StatusOK
	switch status {
	case db.StatusCorrectFlag:
		header = http.StatusAccepted
	case db.StatusAlreadySolved:
		header = http.StatusConflict
		addFlash(s, "Challenge already solved", "warning")
	case db.StatusWrongFlag:
		header = http.StatusNotAcceptable
	}

	if saveSession(w, r, s) {
		w.WriteHeader(header)
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

	users, err := db.GetUsersScores()
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
	data, err := db.GetGraphData()
	if err != nil {
		log.Errorf("Error getting graph data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonData := make(map[string][]GraphPoint)
	for _, d := range data {
		if len(jsonData[d.User]) == 0 {
			point := GraphPoint{
				X: data[0].Timestamp.Format("2006-01-02 15:04:05"),
				Y: 0,
			}
			jsonData[d.User] = append(jsonData[d.User], point)
		}
		point := GraphPoint{
			X: d.Timestamp.Format("2006-01-02 15:04:05"),
			Y: jsonData[d.User][len(jsonData[d.User])-1].Y + d.Points,
		}
		jsonData[d.User] = append(jsonData[d.User], point)
	}

	j, err := json.Marshal(jsonData)
	if err != nil {
		log.Errorf("Error marshaling json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(j)
}
