package routes

import (
	"encoding/gob"
	"net/http"
	"platform/db"
	"platform/log"
)

func authHandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

		// ! TODO: get user from session

		handler(w, r)
	})
}

func StartRouting() {
	gob.Register(&db.User{})
	gob.Register(&Flash{})

	store.Options.Path = ""
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteDefaultMode

	fs := http.FileServer(http.Dir("static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// FILES fileserver

	http.HandleFunc("GET /", home)
	http.HandleFunc("GET /register", register_get)
	http.HandleFunc("POST /register", register_post)
	http.HandleFunc("GET /login", login_get)
	http.HandleFunc("POST /login", login_post)

	authHandleFunc("GET /logout", logout)
	// authHandleFunc("GET /user/{username}", user)

	log.Notice("Serving on :8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Error Serving: %v", err)
	}
}
