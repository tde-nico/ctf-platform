package db

import "platform/log"

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
