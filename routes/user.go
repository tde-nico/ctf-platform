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

	// ! TODO: get user from session
	user := &db.User{
		Username: "test",
		Email:    "test@gmail.com",
		Score:    1337,
		IsAdmin:  true,
	}
	user = nil
	data := Data{
		User: user,
	}

	executeTemplate(w, r, tmpl, &data)
}

func register_get(w http.ResponseWriter, r *http.Request) {
	if !isRegistrationAllowed(w, r) {
		return
	}

	// ! TODO: redirect if session

	tmpl, err := getTemplate(w, "register")
	if err != nil {
		return
	}

	executeTemplate(w, r, tmpl, &Data{})
}

func register_post(w http.ResponseWriter, r *http.Request) {
	if !isRegistrationAllowed(w, r) {
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if !isUserDataValid(w, username, email, password) {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	err := db.RegisterUser(username, email, password)
	if err != nil {
		log.Errorf("Error registering user %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ! TODO create session

	setFlash(w, "success", "Registration Completed")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func login_get(w http.ResponseWriter, r *http.Request) {
	tmpl, err := getTemplate(w, "login")
	if err != nil {
		return
	}

	// ! TODO: redirect if session

	executeTemplate(w, r, tmpl, &Data{})
}

func login_post(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	apiKey, err := loginUser(w, username, password)
	if err != nil {
		log.Errorf("Error logging in: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// ! TODO create session
	log.Infof("Login post %s", apiKey)

	http.Redirect(w, r, "/challenges", http.StatusSeeOther)
}
