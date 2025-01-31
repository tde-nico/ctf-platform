package routes

import (
	"net/http"
	"platform/middleware"
)

func home(ctx *middleware.Ctx) {
	tmpl := getTemplate(ctx, "home")
	if tmpl == nil {
		return
	}

	if ctx.RequestURI != "/" {
		ctx.AddFlash("404 Page not found")
		ctx.Redirect("/", http.StatusSeeOther)
		return
	}

	data := Data{User: ctx.User}

	executeTemplate(ctx, tmpl, &data)
}
