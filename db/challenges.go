package db

import (
	"database/sql"
	"fmt"
)

func GetCategories() ([]string, error) {
	query, err := GetStatement("GetCategories")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]string, 0)
	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func ChallengeExistsID(id int) (bool, error) {
	query, err := GetStatement("ChallengeExistsID")
	if err != nil {
		return false, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(id)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func ChallengeExistsName(name string) (bool, error) {
	query, err := GetStatement("ChallengeExistsName")
	if err != nil {
		return false, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(name)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func GetChallengeName(id int) (string, error) {
	query, err := GetStatement("GetChallengeName")
	if err != nil {
		return "", fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("challenge not found")
	}

	var name string
	err = rows.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetChallenges() (map[string][]Challenge, error) {
	query, err := GetStatement("GetChallenges")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := make(map[string][]Challenge, 0)
	for rows.Next() {
		var chall Challenge
		var host, port, files, hint1, hint2 sql.NullString
		err = rows.Scan(
			&chall.ID,
			&chall.Name,
			&chall.Description,
			&chall.Difficulty,
			&chall.Points,
			&chall.MaxPoints,
			&chall.Solves,
			&host,
			&port,
			&chall.Category,
			&files,
			&chall.Flag,
			&hint1,
			&hint2,
			&chall.Hidden,
			&chall.IsExtra,
		)
		if err != nil {
			return nil, err
		}
		chall.Host = host.String
		chall.Port = port.String
		chall.Files = files.String
		chall.Hint1 = hint1.String
		chall.Hint2 = hint2.String
		challenges[chall.Category] = append(challenges[chall.Category], chall)
	}
	return challenges, nil
}

func CreateChallenge(chal *Challenge) error {
	insert, err := GetStatement("CreateChallenge")
	if err != nil {
		return fmt.Errorf("error getting statement: %v", err)
	}

	_, err = insert.Exec(
		chal.Name,
		chal.Description,
		chal.Difficulty,
		chal.Points,
		chal.MaxPoints,
		chal.Solves,
		chal.Host,
		chal.Port,
		chal.Category,
		chal.Files,
		chal.Flag,
		chal.Hint1,
		chal.Hint2,
		chal.Hidden,
		chal.IsExtra,
	)
	return err
}

func DeleteChallenge(name string) error {
	delete, err := GetStatement("DeleteChallenge")
	if err != nil {
		return fmt.Errorf("error getting statement: %v", err)
	}

	_, err = delete.Exec(name)
	return err
}

func UpdateChallenge(chal *Challenge) error {
	update, err := GetStatement("UpdateChallenge")
	if err != nil {
		return fmt.Errorf("error getting statement: %v", err)
	}

	_, err = update.Exec(
		chal.Name,
		chal.Description,
		chal.Difficulty,
		chal.MaxPoints,
		chal.Host,
		chal.Port,
		chal.Category,
		chal.Files,
		chal.Files,
		chal.Flag,
		chal.Hint1,
		chal.Hint2,
		chal.Hidden,
		chal.IsExtra,
		chal.ID,
	)
	return err
}
