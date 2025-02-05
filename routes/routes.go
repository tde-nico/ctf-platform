package routes

import (
	"net/http"
	"platform/log"
	"platform/middleware"
)

func StartRouting(key []byte) {
	middleware.InitStore(key)

	static := http.FileServer(http.Dir("static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", static))

	middleware.HandleFunc("GET /files/", download)

	middleware.HandleFunc("GET /", home)
	middleware.HandleFunc("GET /register", register_get)
	middleware.HandleFunc("POST /register", register_post)
	middleware.HandleFunc("GET /login", login_get)
	middleware.HandleFunc("POST /login", login_post)
	middleware.HandleFunc("GET /newpw", newpw_get)
	middleware.HandleFunc("POST /newpw", newpw_post)

	middleware.AuthHandleFunc("GET /logout", logout)
	middleware.AuthHandleFunc("GET /user/{username}", userInfo)
	middleware.AuthHandleFunc("GET /challenges", challenges)
	middleware.AuthHandleFunc("POST /submit", submit)
	middleware.AuthHandleFunc("GET /scores", scores)
	middleware.AuthHandleFunc("POST /graph_data", graphData)

	middleware.AdminHandleFunc("GET /admin", admin)
	middleware.AdminHandleFunc("POST /admin/newchal", adminNewChall)
	middleware.AdminHandleFunc("POST /admin/updatechal", adminUpdateChall)
	middleware.AdminHandleFunc("POST /admin/deletechal", adminDeleteChall)
	middleware.AdminHandleFunc("POST /admin/resetpw", adminResetPw)
	middleware.AdminHandleFunc("POST /admin/config", adminConfig)

	log.Notice("Serving on :8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Error Serving: %v", err)
	}
}
