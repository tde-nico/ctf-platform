package db

import "fmt"

func GetKey(key string) (string, error) {
	query, err := GetStatement("GetKey")
	if err != nil {
		return "", fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(key)
	if err != nil {
		return "", fmt.Errorf("error querying key: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("error scanning key: %v", err)
	}

	var value string
	err = rows.Scan(&value)
	if err != nil {
		return "", fmt.Errorf("error scanning key: %v", err)
	}
	return value, nil
}
