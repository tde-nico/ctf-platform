package routes

import (
	"fmt"
	"platform/db"
	"platform/middleware"
	"sort"
)

type CategoryChallenges struct {
	Category   string
	Challenges []db.Challenge
}

type DataChallenges struct {
	Data
	// Challenges map[string][]db.Challenge
	Challenges []CategoryChallenges
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

	challenges, err := db.GetChallenges()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting challenges: %v", err))
		return
	}

	for cat, chall := range challenges {
		data.Challenges = append(data.Challenges, CategoryChallenges{
			Category:   cat,
			Challenges: chall,
		})
	}

	sort.Slice(data.Challenges, func(i, j int) bool {
		if data.Challenges[i].Category == "Intro" {
			return true
		} else if data.Challenges[j].Category == "Intro" {
			return false
		}
		return data.Challenges[i].Category < data.Challenges[j].Category
	})

	executeTemplate(ctx, tmpl, data)
}
