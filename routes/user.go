package routes

import (
	"fmt"
	"net/http"
	"platform/db"
	"platform/log"
	"platform/middleware"
)

type DataUserInfo struct {
	Data
	UserProfile *db.User
	Solves      []db.Solve
}

func userInfo(ctx *middleware.Ctx) {
	tmpl := getTemplate(ctx, "user")
	if tmpl == nil {
		return
	}

	data := &DataUserInfo{}
	data.User = ctx.User

	username := ctx.PathValue("username")
	if username == "" {
		ctx.AddFlash("Invalid username")
		ctx.Redirect("/", http.StatusSeeOther)
		return
	}

	var err error
	data.UserProfile, err = db.GetUserByUsername(username)
	if err != nil {
		log.Warningf("Error getting user by username: %v", err)
		ctx.AddFlash("Invalid username")
		ctx.Redirect("/", http.StatusSeeOther)
		return
	}

	data.Solves, err = db.GetUserSolves(data.UserProfile)
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting solves by user: %v", err))
		return
	}

	executeTemplate(ctx, tmpl, data)
}
