package db

import (
	"fmt"
	"platform/utils"
)

func FlagExists(flag string) error {
	query, err := GetStatement("FlagExists")
	if err != nil {
		return fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(flag)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return fmt.Errorf("flag already exists")
	}
	return nil
}

func getScores() ([]User, error) {
	query, err := GetStatement("GetUsersScores")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Score)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func getVisibleChallsByCategory() (map[string][]int, error) {
	query, err := GetStatement("GetVisibleChallengesByCategory")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challs := make(map[string][]int)
	for rows.Next() {
		var name string
		var count, extra int
		err = rows.Scan(&name, &count, &extra)
		if err != nil {
			return nil, err
		}
		challs[name] = []int{count - extra, extra}
	}

	return challs, nil
}

func GetUsersScores() ([]UserScore, error) {
	users, err := getScores()
	if err != nil {
		return nil, err
	}

	challs, err := getVisibleChallsByCategory()
	if err != nil {
		return nil, err
	}

	scoreUsers := make([]UserScore, len(users))
	for i, user := range users {
		scoreUsers[i].Username = user.Username
		scoreUsers[i].Score = user.Score

		solves, err := GetUserSolves(&user)
		if err != nil {
			return nil, err
		}

		counter := make(map[string][]int)
		for _, s := range solves {
			c, ok := counter[s.ChalCategory]
			if !ok {
				c = []int{0, 0}
			}
			if s.ChalExtra {
				c[1]++
			} else {
				c[0]++
			}
			counter[s.ChalCategory] = c
		}

		for category, counts := range counter {
			if category == "Intro" {
				continue
			}

			challCounts, ok := challs[category]
			if !ok {
				continue
			}
			if challCounts[0] == counts[0] {
				scoreUsers[i].Badges = append(scoreUsers[i].Badges, Badges{
					Name:  category,
					Char:  string(category[0]),
					Extra: false,
				})
				if challCounts[1] > 0 && challCounts[1] == counts[1] {
					scoreUsers[i].Badges[len(scoreUsers[i].Badges)-1].Extra = true
				}
			}
		}
	}
	return scoreUsers, nil
}

func GetUserSolves(user *User) ([]Solve, error) {
	query, err := GetStatement("GetUserSolves")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	solves := make([]Solve, 0)
	for rows.Next() {
		var solve Solve
		var timestamp *string
		err = rows.Scan(&solve.ChalName, &solve.ChalCategory, &solve.ChalExtra, &timestamp)
		if err != nil {
			return nil, err
		}
		solve.Timestamp, err = utils.ParseTime(timestamp)
		if err != nil {
			return nil, err
		}
		solves = append(solves, solve)
	}

	return solves, nil
}

func isChallengeSolved(user *User, chalID int) (bool, error) {
	query, err := GetStatement("IsChallengeSolved")
	if err != nil {
		return false, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(user.ID, chalID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}
