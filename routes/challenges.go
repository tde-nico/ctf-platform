package routes

import (
	"fmt"
	"platform/db"
	"platform/middleware"
)

type DataChallenges struct {
	Data
	Challenges map[string][]db.Challenge
	Solves     map[string]bool
}

func challenges(ctx *middleware.Ctx) {
	tmpl := getTemplate(ctx, "challenges")
	if tmpl == nil {
		return
	}

	data := &DataChallenges{}
	data.User = ctx.User

	solves, err := db.GetUserSolves(data.User)
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting solves by user: %v", err))
		return
	}

	data.Solves = make(map[string]bool)
	for _, solve := range solves {
		data.Solves[solve.ChalName] = true
	}

	data.Challenges, err = db.GetChallenges()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting challenges: %v", err))
		return
	}

	// TODO: make chals decscription HTML safe

	executeTemplate(ctx, tmpl, data)
}
