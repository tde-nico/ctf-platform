package routes

import (
	"fmt"
	"net/http"
	"platform/db"
	"platform/log"
	"regexp"
	"strings"
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

const SEPARATOR = "|"
const MAX_PASSWORD_LENGTH = 6

var USERNAME_REGEX = regexp.MustCompile(`[0-9a-zA-Z_!@#€\-&+]{4,32}`)

func setFlash(w http.ResponseWriter, t string, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "flash",
		Value: strings.Join([]string{t, msg}, SEPARATOR),
	})
}

// ! TODO make it return a slice of flashes
func getFlash(w http.ResponseWriter, r *http.Request) (*Flash, error) {
	cookie, err := r.Cookie("flash")
	if err != nil && err != http.ErrNoCookie {
		log.Errorf("Error getting flash cookie %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	} else if err == nil {
		flashParts := strings.SplitN(cookie.Value, SEPARATOR, 2)
		http.SetCookie(w, &http.Cookie{Name: "flash", MaxAge: -1})
		if len(flashParts) == 2 {
			return &Flash{
				Message: flashParts[1],
				Type:    flashParts[0],
			}, nil
		}
	}
	return nil, nil
}

func getTemplate(w http.ResponseWriter, page string) (*template.Template, error) {
	tmpl, err := template.New("").ParseFiles("templates/base.html", fmt.Sprintf("templates/%s.html", page))
	if err != nil {
		log.Errorf("Error parsing template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}
	return tmpl, nil
}

func executeTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data *Data) {
	flash, err := getFlash(w, r)
	if err != nil {
		return
	}
	if flash != nil {
		data.Flashes = append(data.Flashes, *flash)
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Errorf("Error executing template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func isRegistrationAllowed(w http.ResponseWriter, r *http.Request) bool {
	allowed := db.GetConfig("registration-allowed")
	if allowed == nil || *allowed == 0 {
		setFlash(w, "danger", "Registration is disabled")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return false
	}
	return true
}

func isUserDataValid(w http.ResponseWriter, username, email, password string) bool {
	if len(password) < MAX_PASSWORD_LENGTH {
		setFlash(w, "danger", "Password is too short")
		return false
	}

	if strings.ToLower(username) == "admin" {
		setFlash(w, "danger", "Username is reserved")
		return false
	}

	if !USERNAME_REGEX.MatchString(username) {
		setFlash(w, "danger", "Invalid username")
		return false
	}

	if db.UserExists(username) {
		setFlash(w, "danger", "Username already taken")
		return false
	}

	if db.EmailExists(email) {
		setFlash(w, "danger", "Email already taken")
		return false
	}

	return true
}

func loginUser(w http.ResponseWriter, username, password string) (string, error) {
	username = strings.TrimSpace(username)
	if !USERNAME_REGEX.MatchString(username) {
		setFlash(w, "danger", "Invalid username")
		return "", fmt.Errorf("invalid username")
	}

	apiKey, err := db.LoginUser(username, password)
	if err != nil {
		setFlash(w, "danger", "Invalid username or password")
		return "", err
	}

	return apiKey, nil
}
