package db

import (
	"database/sql"
	"fmt"
	"platform/log"
)

func GetConfig(key string) *int {
	if db == nil {
		log.Error("Database not initialized")
		return nil
	}
	rows, err := db.Query("SELECT value FROM config WHERE key = ?", key)
	if err != nil {
		log.Errorf("Error querying config: %v", err)
		return nil
	}
	defer rows.Close()

	if !rows.Next() {
		log.Errorf("Error scanning config: %v", err)
		return nil
	}

	var value int
	err = rows.Scan(&value)
	if err != nil {
		log.Errorf("Error scanning config: %v", err)
		return nil
	}
	return &value
}

func GetConfigs() ([]*Config, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT * FROM config")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configs := make([]*Config, 0)
	for rows.Next() {
		var config Config
		err = rows.Scan(&config.Key, &config.Type, &config.Value, &config.Desc)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}
	return configs, nil
}

func SetConfig(key, value string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := db.Exec("UPDATE config SET value = ? WHERE key = ?", value, key)
	return err
}

func GetUsers() ([]*User, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT username, email, score, is_admin FROM users ORDER BY username")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Username, &user.Email, &user.Score, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetChallenges() (map[string][]*Challenge, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT * FROM challenges")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := make(map[string][]*Challenge, 0)
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
		challenges[chall.Category] = append(challenges[chall.Category], &chall)
	}
	return challenges, nil
}

func GetSubmissions() ([]*Submission, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT u.username, c.name, s.status, s.flag, s.timestamp FROM users AS u, challenges AS c, submissions AS s WHERE s.userid=u.id AND s.chalid=c.id ORDER BY s.timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]*Submission, 0)
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
		submissions = append(submissions, &sub)
	}
	return submissions, nil
}
