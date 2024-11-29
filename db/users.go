package db

import (
	"errors"
	"fmt"
	"platform/log"
	"platform/utils"
)

const APIKEY_LENGTH = 32
const SALT_LENGTH = 8

var ErrDatabaseNotInitialized = errors.New("database not initialized")

func UserExists(username string) bool {
	if db == nil {
		log.Error("Database not initialized")
		return false
	}

	rows, err := db.Query("SELECT id FROM users WHERE username = ?", username)
	if err != nil {
		log.Errorf("Error querying users for username: %v", err)
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func EmailExists(email string) bool {
	if db == nil {
		log.Error("Database not initialized")
		return false
	}

	rows, err := db.Query("SELECT id FROM users WHERE email = ?", email)
	if err != nil {
		log.Errorf("Error querying users for mail: %v", err)
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func createUser(username, email, salt, secret, apiKey string) error {
	if db == nil {
		log.Error("Database not initialized")
		return ErrDatabaseNotInitialized
	}

	_, err := db.Exec("INSERT INTO users (username, email, salt, password, apikey) VALUES (?, ?, ?, ?, ?)", username, email, salt, secret, apiKey)
	if err != nil {
		log.Errorf("Error inserting user: %v", err)
		return err
	}

	return nil
}

func RegisterUser(username, email, password string) error {
	_, apiKeyHex, err := utils.GetRand(APIKEY_LENGTH)
	if err != nil {
		return err
	}
	salt, saltHex, err := utils.GetRand(SALT_LENGTH)
	if err != nil {
		return err
	}
	secretHex := utils.HashPassword(password, salt)
	err = createUser(username, email, string(saltHex), string(secretHex), string(apiKeyHex))
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(username, password string) (*User, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.Errorf("User not found")
		return nil, fmt.Errorf("user not found")
	}

	var user User
	var saltHex, secretHex string
	err = rows.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&saltHex,
		&secretHex,
		&user.Apikey,
		&user.Score,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, err
	}

	salt, err := utils.HexToBytes(saltHex)
	if err != nil {
		return nil, err
	}

	hash := utils.HashPassword(password, salt)
	if secretHex != hash {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}
