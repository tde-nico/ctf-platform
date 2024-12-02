package db

import (
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
		log.Error("Database not initialized")
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query("SELECT * FROM config")
	if err != nil {
		log.Errorf("Error querying config: %v", err)
		return nil, err
	}
	defer rows.Close()

	configs := make([]*Config, 0)
	for rows.Next() {
		var config Config
		err = rows.Scan(&config.Key, &config.Type, &config.Value, &config.Desc)
		if err != nil {
			log.Errorf("Error scanning config: %v", err)
			return nil, err
		}
		configs = append(configs, &config)
	}
	return configs, nil
}

func SetConfig(key, value string) error {
	if db == nil {
		log.Error("Database not initialized")
		return fmt.Errorf("database not initialized")
	}

	_, err := db.Exec("UPDATE config SET value = ? WHERE key = ?", value, key)
	if err != nil {
		log.Errorf("Error updating config: %v", err)
		return err
	}
	return nil
}
