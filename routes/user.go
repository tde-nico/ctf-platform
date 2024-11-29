package routes

import (
	"net/http"
	"platform/db"
	"platform/log"
)

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := getTemplate(w, "home")
	if err != nil {
		return
	}

	session, ok := getSession(w, r)
	if !ok {
		return
	}

	user := getSessionUser(session)

	data := Data{User: user}
	if r.RequestURI != "/" {
		data.Flashes = append(data.Flashes, Flash{
			Type:    "danger",
			Message: "404 Page not found",
		})
	}

	executeTemplate(w, r, session, tmpl, &data)
}

func register_get(w http.ResponseWriter, r *http.Request) {
	session, ok := getSession(w, r)
	if !ok {
		return
	}

	// ! TODO: redirect if session

	if !isRegistrationAllowed(w, r, session) {
		return
	}

	tmpl, err := getTemplate(w, "register")
	if err != nil {
		return
	}

	executeTemplate(w, r, session, tmpl, &Data{})
}

func register_post(w http.ResponseWriter, r *http.Request) {
	session, ok := getSession(w, r)
	if !ok {
		return
	}

	// ! TODO: redirect if session

	if !isRegistrationAllowed(w, r, session) {
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if !isUserDataValid(session, username, email, password) {
		if !saveSession(w, r, session) {
			return
		}
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	err := db.RegisterUser(username, email, password)
	if err != nil {
		log.Errorf("Error registering user %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	addFlash(session, "Registration Completed", "success")

	if !saveSession(w, r, session) {
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func login_get(w http.ResponseWriter, r *http.Request) {
	session, ok := getSession(w, r)
	if !ok {
		return
	}

	tmpl, err := getTemplate(w, "login")
	if err != nil {
		return
	}

	// ! TODO: redirect if session

	executeTemplate(w, r, session, tmpl, &Data{})
}

func login_post(w http.ResponseWriter, r *http.Request) {
	// username := r.FormValue("username")
	// password := r.FormValue("password")
	username := "nonnon"
	password := "nonnon"

	session, ok := getSession(w, r)
	if !ok {
		return
	}

	// ! TODO: redirect if session

	user, err := loginUser(session, username, password)
	if err != nil {
		log.Errorf("Error logging in: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session.Values["user"] = user

	if !saveSession(w, r, session) {
		return
	}

	http.Redirect(w, r, "/challenges", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, ok := getSession(w, r)
	if !ok {
		return
	}

	session.Options.MaxAge = -1

	if !saveSession(w, r, session) {
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

// func user(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := getTemplate(w, "user")
// 	if err != nil {
// 		return
// 	}

// 	// ! TODO get users
// 	log.Infof("User %s", r.PathValue("username"))

// 	// ! TODO: FINISH

// 	executeTemplate(w, r, tmpl, &Data{User: user, UserProfile: user2, Solves: solves})
// }
