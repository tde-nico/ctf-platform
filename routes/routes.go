package routes

import (
	"net/http"
	"platform/log"
)

func StartRouting() {
	http.HandleFunc("GET /", home)
	http.HandleFunc("GET /register", register_get)
	http.HandleFunc("POST /register", register_post)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// FILES fileserver

	log.Notice("Serving on :8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Error Serving: %v", err)
	}
}
