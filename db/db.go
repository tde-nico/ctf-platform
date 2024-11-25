package db

import (
	"database/sql"
	"io"
	"os"
	"platform/log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(path string) {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Noticef("Successfully connected to the database")
}

func CloseDB() {
	if db == nil {
		log.Error("Database not initialized")
		return
	}

	err := db.Close()
	if err != nil {
		log.Errorf("Error closing database: %v", err)
		return
	}

	log.Noticef("Database connection closed")
}

func ExecSQLFile(path string) {
	if db == nil {
		log.Error("Database not initialized")
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Errorf("SQL file does not exist: %v", err)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		log.Errorf("Error opening SQL file: %v", err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Errorf("Error reading SQL file: %v", err)
		return
	}

	_, err = db.Exec(string(data))
	if err != nil {
		log.Errorf("Error executing SQL: %v", err)
		return
	}

	log.Noticef("SQL executed successfully")
}

func CleanDB() {
	if db == nil {
		log.Error("Database not initialized")
		return
	}

	_, err := db.Exec("DELETE FROM submissions")
	if err != nil {
		log.Errorf("Error cleaning submissions table: %v", err)
		return
	}

	_, err = db.Exec("DELETE FROM solves")
	if err != nil {
		log.Errorf("Error cleaning solves table: %v", err)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE is_admin=0")
	if err != nil {
		log.Errorf("Error cleaning users table: %v", err)
		return
	}

	log.Noticef("Database cleaned successfully")
}

func PruneDB() {
	if db == nil {
		log.Error("Database not initialized")
		return
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Errorf("Error querying table names: %v", err)
		return
	}
	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			log.Errorf("Error scanning table: %v", err)
			return
		}
		tables = append(tables, table)
	}

	for _, table := range tables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			log.Errorf("Error pruning %s table: %v", table, err)
			return
		}
	}

	log.Noticef("Database pruned successfully")
}

func DropTables() {
	if db == nil {
		log.Error("Database not initialized")
		return
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Errorf("Error querying table names: %v", err)
		return
	}
	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			log.Errorf("Error scanning table: %v", err)
			return
		}
		tables = append(tables, table)
	}

	for _, table := range tables {
		_, err := db.Exec("DROP TABLE " + table)
		if err != nil {
			log.Errorf("Error dropping %s table: %v", table, err)
			return
		}
	}
}
