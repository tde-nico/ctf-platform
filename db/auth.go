package db

import (
	"fmt"
	"platform/utils"
)

const APIKEY_LENGTH = 32
const SALT_LENGTH = 8
const INVALID_PREFIX = "INVALID_"

func UserExists(username string) (bool, error) {
	query, err := GetStatement("UserExists")
	if err != nil {
		return false, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	return true, nil
}

func EmailExists(email string) (bool, error) {
	query, err := GetStatement("EmailExists")
	if err != nil {
		return false, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(email)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	return true, nil
}

func updatePassword(username, salt, secret, apiKey string) error {
	update, err := GetStatement("UpdatePassword")
	if err != nil {
		return fmt.Errorf("error getting statement: %v", err)
	}

	_, err = update.Exec(salt, secret, apiKey, username)
	return err
}

func ChangePassword(username, password string, invalid bool) error {
	_, apiKeyHex, err := utils.GetRand(APIKEY_LENGTH)
	if err != nil {
		return err
	}

	salt, saltHex, err := utils.GetRand(SALT_LENGTH)
	if err != nil {
		return err
	}
	secretHex := utils.HashPassword(password, salt)

	apiKey := string(apiKeyHex)
	if invalid {
		apiKey = INVALID_PREFIX + apiKey
	}

	err = updatePassword(username, string(saltHex), string(secretHex), apiKey)
	return err
}

func ResetPassword(username string) (string, error) {
	exists, err := UserExists(username)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("user not found")
	}

	_, password, err := utils.GetRand(SALT_LENGTH)
	if err != nil {
		return "", err
	}

	err = ChangePassword(username, password, true)
	if err != nil {
		return "", err
	}

	return password, nil
}

func createUser(username, email, salt, secret, apiKey string) error {
	insert, err := GetStatement("CreateUser")
	if err != nil {
		return fmt.Errorf("error getting statement: %v", err)
	}

	_, err = insert.Exec(username, email, salt, secret, apiKey)
	return err
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
	return err
}

func LoginUser(username, password string) (string, error) {
	query, err := GetStatement("LoginUser")
	if err != nil {
		return "", fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(username)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("user \"%s\" not found", username)
	}

	var apikey, saltHex, secretHex string
	err = rows.Scan(&apikey, &saltHex, &secretHex)
	if err != nil {
		return "", err
	}

	salt, err := utils.HexToBytes(saltHex)
	if err != nil {
		return "", err
	}

	hash := utils.HashPassword(password, salt)
	if secretHex != hash {
		return "", fmt.Errorf("invalid password")
	}

	return apikey, nil
}
