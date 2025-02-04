package routes

import (
	"net/http"
	"platform/middleware"
)

func logout(ctx *middleware.Ctx) {
	ctx.ExpireCookie()
	ctx.Redirect("/", http.StatusSeeOther)
}
