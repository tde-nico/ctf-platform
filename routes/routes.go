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

func adminHandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request, s *sessions.Session)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		s, ok := getSession(w, r)
		if !ok {
			return
		}

		user := getSessionUser(s)
		if user == nil || !user.IsAdmin {
			addFlash(s, "You must be admin to access that page")
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

	static := http.FileServer(http.Dir("static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", static))

	files := http.FileServer(http.Dir("files"))
	http.Handle("GET /files/", http.StripPrefix("/files/", files))

	handleFunc("GET /", home)
	handleFunc("GET /register", register_get)
	handleFunc("POST /register", register_post)
	handleFunc("GET /login", login_get)
	handleFunc("POST /login", login_post)

	authHandleFunc("GET /logout", logout)

	authHandleFunc("GET /user/{username}", userInfo)
	authHandleFunc("GET /challenges", challenges)
	authHandleFunc("POST /submit", submit)        // TODO
	authHandleFunc("GET /scores", scores)         // TODO
	authHandleFunc("POST /graph_data", graphData) // TODO
	// authHandleFunc("GET /newpw", home)
	// authHandleFunc("POST /newpw", home)

	adminHandleFunc("GET /admin", admin)
	adminHandleFunc("POST /admin/newchal", adminNewChall)
	adminHandleFunc("POST /admin/updatechal", adminUpdateChall)
	adminHandleFunc("POST /admin/deletechal", adminDeleteChall)
	adminHandleFunc("POST /admin/resetpw", adminResetPw) // TODO
	adminHandleFunc("POST /admin/config", adminConfig)

	log.Notice("Serving on :8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Error Serving: %v", err)
	}
}
