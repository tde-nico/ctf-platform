package db

import (
	"fmt"
	"platform/utils"
)

// TODO: collapse GetUserByAPIKey and GetUserByUsername into one function

func GetUserByAPIKey(apiKey string) (*User, error) {
	query, err := GetStatement("GetUserByAPIKey")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(apiKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	user := User{ApiKey: apiKey}
	err = rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Score,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	query, err := GetStatement("GetUserByUsername")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	user := User{Username: username}
	err = rows.Scan(
		&user.ID,
		&user.Email,
		&user.Score,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUsers() ([]User, error) {
	query, err := GetStatement("GetUsers")
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
		err = rows.Scan(
			&user.Username,
			&user.Email,
			&user.Score,
			&user.IsAdmin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetGraphData() ([]GraphData, error) {
	query, err := GetStatement("GetGraphData")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]GraphData, 0)
	for rows.Next() {
		var tmp GraphData
		var timestamp *string
		err = rows.Scan(&tmp.User, &tmp.Points, &timestamp)
		if err != nil {
			return nil, err
		}
		tmp.Timestamp, err = utils.ParseTime(timestamp)
		if err != nil {
			return nil, err
		}
		data = append(data, tmp)
	}

	return data, nil
}
