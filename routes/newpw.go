package routes

import (
	"net/http"
	"platform/db"
	"platform/log"
	"platform/middleware"
)

func newpw_get(ctx *middleware.Ctx) {
	tmpl := getTemplate(ctx, "newpw")
	if tmpl == nil {
		return
	}

	data := Data{User: ctx.User}
	executeTemplate(ctx, tmpl, &data)
}

func newpw_post(ctx *middleware.Ctx) {
	if ctx.User == nil {
		ctx.AddFlash("You must be logged in to access that page")
		ctx.Redirect("/", http.StatusSeeOther)
		return
	}

	newpw := ctx.FormValue("password")

	if len(newpw) < MAX_PASSWORD_LENGTH {
		ctx.AddFlash("Password is too short")
		ctx.Redirect("/newpw", http.StatusSeeOther)
		return
	}

	err := db.ChangePassword(ctx.User.Username, newpw, false)
	if err != nil {
		log.Errorf("Error changing password: %v", err)
		ctx.AddFlash("Error changing password")
		ctx.Redirect("/newpw", http.StatusSeeOther)
		return
	}

	ctx.SetSessionValue("apikey", "")

	ctx.AddFlash("Password changed successfully", "success")
	ctx.Redirect("/", http.StatusSeeOther)
}
