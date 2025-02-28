package db

import (
	"fmt"
)

func GetConfig(key string) (int64, error) {
	query, err := GetStatement("GetConfig")
	if err != nil {
		return 0, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query(key)
	if err != nil {
		return 0, fmt.Errorf("error querying config: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("error scanning config: %v", err)
	}

	var value int64
	err = rows.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("error scanning config: %v", err)
	}
	return value, nil
}

func GetConfigs() ([]Config, error) {
	query, err := GetStatement("GetConfigs")
	if err != nil {
		return nil, fmt.Errorf("error getting statement: %v", err)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configs := make([]Config, 0)
	for rows.Next() {
		var config Config
		err = rows.Scan(&config.Key, &config.Type, &config.Value, &config.Desc)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}

func SetConfig(key, value string) error {
	update, err := GetStatement("SetConfig")
	if err != nil {
		return fmt.Errorf("error getting statement: %v", err)
	}

	_, err = update.Exec(value, key)
	return err
}
