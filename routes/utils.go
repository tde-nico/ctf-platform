package routes

import (
	"net/http"
	"platform/db"
	"platform/log"
	"regexp"
	"strings"
)

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
