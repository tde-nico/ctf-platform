package db

import (
	"fmt"
	"platform/log"
	"platform/utils"
)

const (
	StatusWrongFlag int = iota
	StatusAlreadySolved
	StatusCorrectFlag
)

func GetSubmissions() ([]Submission, error) {
	query, err := GetStatement("GetSubmissions")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]Submission, 0)
	for rows.Next() {
		var sub Submission
		err = rows.Scan(
			&sub.UserUsername,
			&sub.ChalName,
			&sub.Status,
			&sub.Flag,
			&sub.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, sub)
	}
	return submissions, nil
}

func getChallIfCorrectFlag(chalID int, flag string) (*Challenge, error) {
	query, err := GetStatement("GetChallIfCorrectFlag")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(chalID, flag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var chal Challenge
	err = rows.Scan(&chal.Name, &chal.Solves)
	if err != nil {
		return nil, err
	}

	return &chal, nil
}

func SubmitFlag(user *User, chalID int, flag string) (int, error) {
	now := utils.CurrentTime()

	exists, err := ChallengeExistsID(chalID)
	if err != nil {
		return StatusWrongFlag, err
	}
	if !exists {
		return StatusWrongFlag, fmt.Errorf("challenge not found")
	}

	insertSubmit, err := GetStatement("SubmitFlag")
	if err != nil {
		return StatusWrongFlag, fmt.Errorf("error getting statement: %v", err)
	}

	status, err := isChallengeSolved(user, chalID)
	if err != nil {
		return StatusWrongFlag, err
	}
	if status {
		_, err := insertSubmit.Exec(user.ID, chalID, "r", flag, now)
		if err != nil {
			return StatusWrongFlag, err
		}

		return StatusAlreadySolved, fmt.Errorf("challenge already solved")
	}

	chal, err := getChallIfCorrectFlag(chalID, flag)
	if err != nil {
		return StatusWrongFlag, err
	}
	if chal == nil {
		_, err := insertSubmit.Exec(user.ID, chalID, "w", flag, now)
		if err != nil {
			return StatusWrongFlag, err
		}

		return StatusWrongFlag, nil
	}

	if chal.Solves == 0 && !user.IsAdmin {
		log.Noticef("First Blood on %s from %s", chal.Name, user.Username)
		// TODO: bot first blood
	}

	_, err = insertSubmit.Exec(user.ID, chalID, "c", flag, now)
	if err != nil {
		return StatusWrongFlag, err
	}

	insertSolve, err := GetStatement("InsertSolve")
	if err != nil {
		return StatusWrongFlag, fmt.Errorf("error getting statement: %v", err)
	}

	_, err = insertSolve.Exec(user.ID, chalID, now)
	if err != nil {
		return StatusWrongFlag, err
	}

	return StatusCorrectFlag, nil
}
