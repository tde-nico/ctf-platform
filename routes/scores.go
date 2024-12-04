package routes

import (
	"net/http"
	"platform/log"

	"github.com/gorilla/sessions"
)

func submit(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	log.Infof("submit")
}

func scores(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	log.Infof("scores")
}

func graphData(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	log.Infof("graphData")
}
