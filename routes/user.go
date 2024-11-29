package routes

import (
	"net/http"
	"platform/db"
	"platform/log"
	"text/template"
)

type Flash struct {
	Message string
	Type    string
}

type Data struct {
	User    *db.User
	Flashes []Flash
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/base.html", "templates/home.html")
	if err != nil {
		log.Errorf("Error parsing home template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	flash, err := getFlash(w, r)
	if err != nil {
		return
	}
	if flash != nil {
		data.Flashes = append(data.Flashes, *flash)
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Errorf("Error executing home template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func register_get(w http.ResponseWriter, r *http.Request) {
	if !isRegistrationAllowed(w, r) {
		return
	}

	// ! TODO: redirect if session

	tmpl, err := template.New("").ParseFiles("templates/base.html", "templates/register.html")
	if err != nil {
		log.Errorf("Error parsing register template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := Data{}
	flash, err := getFlash(w, r)
	if err != nil {
		return
	}
	if flash != nil {
		data.Flashes = append(data.Flashes, *flash)
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Errorf("Error executing home template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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
