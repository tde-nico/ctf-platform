package routes

import (
	"fmt"
	"net/http"
	"platform/db"
	"platform/log"
	"platform/middleware"
	"strings"
)

func isUserDataValid(ctx *middleware.Ctx, username, email, password string) bool {
	if len(password) < MAX_PASSWORD_LENGTH {
		ctx.AddFlash("Password is too short")
		return false
	}

	if strings.ToLower(username) == "admin" {
		ctx.AddFlash("Username is reserved")
		return false
	}

	if !USERNAME_REGEX.MatchString(username) {
		ctx.AddFlash("Invalid username")
		return false
	}

	exists, err := db.UserExists(username)
	if err != nil {
		log.Errorf("Error checking username: %v", err)
		return false
	}
	if exists {
		ctx.AddFlash("Username already taken")
		return false
	}

	exists, err = db.EmailExists(email)
	if err != nil {
		log.Errorf("Error checking email: %v", err)
		return false
	}
	if exists {
		ctx.AddFlash("Email already taken")
		return false
	}

	return true
}

func isRegistrationAllowed(ctx *middleware.Ctx) bool {
	allowed, err := db.GetConfig("registration-allowed")
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting registration-allowed config: %v", err))
		return false
	}
	if allowed == 0 {
		ctx.AddFlash("Registration is disabled")
		ctx.Redirect("/", http.StatusSeeOther)
		return false
	}
	return true
}

func register_get(ctx *middleware.Ctx) {
	if ctx.User != nil {
		ctx.Redirect("/", http.StatusSeeOther)
		return
	}

	if !isRegistrationAllowed(ctx) {
		return
	}

	tmpl := getTemplate(ctx, "register")
	if tmpl == nil {
		return
	}

	executeTemplate(ctx, tmpl, &Data{})
}

func register_post(ctx *middleware.Ctx) {
	if ctx.User != nil {
		ctx.Redirect("/", http.StatusSeeOther)
		return
	}

	if !isRegistrationAllowed(ctx) {
		return
	}

	username := ctx.FormValue("username")
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	if !isUserDataValid(ctx, username, email, password) {
		ctx.Redirect("/register", http.StatusSeeOther)
		return
	}

	err := db.RegisterUser(username, email, password)
	if err != nil {
		ctx.InternalError(fmt.Errorf("error registering user: %v", err))
		return
	}

	ctx.AddFlash("Registration Completed", "success")
	ctx.Redirect("/login", http.StatusSeeOther)
}
