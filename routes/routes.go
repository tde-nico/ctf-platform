package routes

import (
	"encoding/gob"
	"net/http"
	"platform/db"
	"platform/log"

	"github.com/gorilla/sessions"
)

func handleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request, s *sessions.Session)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		s, ok := getSession(w, r)
		if !ok {
			return
		}
		handler(w, r, s)
	})
}

func authHandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request, s *sessions.Session)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		s, ok := getSession(w, r)
		if !ok {
			return
		}

		if getSessionUser(s) == nil {
			addFlash(s, "You must be logged in to access that page")
			if saveSession(w, r, s) {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}

		handler(w, r, s)
	})
}

func StartRouting() {
	gob.Register(&db.User{})
	gob.Register(&Flash{})

	store.Options.Path = "/"
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteDefaultMode

	fs := http.FileServer(http.Dir("static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// FILES fileserver

	handleFunc("GET /", home)
	handleFunc("GET /register", register_get)
	handleFunc("POST /register", register_post)
	handleFunc("GET /login", login_get)
	handleFunc("POST /login", login_post)

	authHandleFunc("GET /logout", logout)
	// authHandleFunc("GET /user/{username}", user)

	log.Notice("Serving on :8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Error Serving: %v", err)
	}
}
