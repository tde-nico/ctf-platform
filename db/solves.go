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

func GetUsersScores() ([]User, error) {
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

func GetUsersScoreboard() ([]UserScore, error) {
	users, err := GetUsersScores()
	if err != nil {
		return nil, err
	}

	query, err := GetStatement("GetUsersBadges")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	badges := make(map[int][]Badge)
	for rows.Next() {
		var badge Badge
		var userID int
		err = rows.Scan(&userID, &badge.Name, &badge.Desc, &badge.Extra)
		if err != nil {
			return nil, err
		}
		if _, ok := badges[userID]; !ok {
			badges[userID] = make([]Badge, 0)
		}
		badge.Char = string(badge.Name[0])
		badges[userID] = append(badges[userID], badge)
	}

	scoreUsers := make([]UserScore, len(users))
	for i, user := range users {
		scoreUsers[i].Username = user.Username
		scoreUsers[i].Score = user.Score
		bgs, ok := badges[user.ID]
		if ok {
			scoreUsers[i].Badges = bgs
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
