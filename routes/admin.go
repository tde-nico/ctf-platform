package routes

import (
	"fmt"
	"net/http"
	"platform/db"
	"platform/log"
	"platform/middleware"
)

type AdminPanelData struct {
	Data
	Users        []db.User
	Categories   []string
	Difficulties []string
	Challenges   map[string][]db.Challenge
	Submissions  []db.Submission
	Config       []db.Config
}

func admin(ctx *middleware.Ctx) {
	tmpl := getTemplate(ctx, "admin")
	if tmpl == nil {
		return
	}

	data := &AdminPanelData{}
	data.User = ctx.User

	var err error
	data.Users, err = db.GetUsers()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting users: %v", err))
		return
	}

	data.Categories = db.CATEGORIES
	data.Difficulties = db.DIFFICULTIES

	data.Challenges, err = db.GetChallenges()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting challenges: %v", err))
		return
	}

	data.Submissions, err = db.GetSubmissions()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting submissions: %v", err))
		return
	}

	data.Config, err = db.GetConfigs()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting configs: %v", err))
		return
	}

	executeTemplate(ctx, tmpl, data)
}

func adminNewChall(ctx *middleware.Ctx) {
	ctx.ParseMultipartForm()

	chal := getChallFromForm(ctx)
	if chal == nil {
		return
	}

	if len(ctx.MultipartForm().File) > 0 {
		if !extractChallengeFiles(ctx, chal) {
			return
		}
	}

	createChallenge(ctx, chal)
}

func adminUpdateChall(ctx *middleware.Ctx) {
	ctx.ParseMultipartForm()

	chal := getChallFromForm(ctx)
	if chal == nil {
		return
	}

	err := isChallengeValid(chal)
	if err != nil {
		ctx.AddFlash(err.Error())
		ctx.Redirect("/admin", http.StatusSeeOther)
		return
	}

	err = renameChallenge(chal)
	if err != nil {
		log.Errorf("Error renaming challenge: %v", err)
		ctx.AddFlash("Error renaming challenge")
		ctx.Redirect("/admin", http.StatusSeeOther)
		return
	}

	if len(ctx.MultipartForm().File) > 0 {
		err = deleteChallengeFiles(chal.Name)
		if err != nil {
			log.Errorf("Error deleting challenge files: %s: %v", chal.Name, err)
			ctx.AddFlash("Error deleting challenge")
			ctx.Redirect("/admin", http.StatusSeeOther)
			return
		}

		if !extractChallengeFiles(ctx, chal) {
			return
		}
	}

	err = db.UpdateChallenge(chal)
	if err != nil {
		log.Errorf("Error updating challenge: %v", err)
		ctx.AddFlash("Error updating challenge")
		ctx.Redirect("/admin", http.StatusSeeOther)
		return
	}

	ctx.AddFlash("Challenge Updated Successfully", "success")
	ctx.Redirect("/admin", http.StatusSeeOther)
}

func adminDeleteChall(ctx *middleware.Ctx) {
	name := ctx.FormValue("name")
	log.Infof("Delete Challenge: %s", name)

	if !deleteChallenge(ctx, name) {
		return
	}

	ctx.AddFlash("Challenge Deleted Successfully", "success")
	ctx.Redirect("/admin", http.StatusSeeOther)
}

func adminResetPw(ctx *middleware.Ctx) {
	username := ctx.FormValue("username")
	log.Warningf("Reset Password Username: %s", username) //! FIX: username empty
	password, err := db.ResetPassword(username)
	if err != nil {
		log.Errorf("Error resetting password: %v", err)
		ctx.AddFlash("Error resetting password")
		ctx.Redirect("/admin", http.StatusSeeOther)
		return
	}

	msg := fmt.Sprintf("Temporary password [%s]: %s", username, password)
	ctx.AddFlash(msg, "success")
	ctx.Redirect("/admin", http.StatusSeeOther)
}

func adminConfig(ctx *middleware.Ctx) {
	configs, err := db.GetConfigs()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting configs: %v", err))
		return
	}

	for _, config := range configs {
		value := ctx.FormValue(config.Key)
		if value == "" {
			err = db.SetConfig(config.Key, "0")
		} else {
			err = db.SetConfig(config.Key, value)
		}
		if err != nil {
			log.Errorf("Error setting config %s: %v", config.Key, err)
			ctx.AddFlash("Error updating configuration", "danger")
			ctx.Redirect("/admin", http.StatusSeeOther)
			return
		}
	}

	ctx.AddFlash("Configuration updated successfully", "success")
	ctx.Redirect("/admin", http.StatusSeeOther)
}
