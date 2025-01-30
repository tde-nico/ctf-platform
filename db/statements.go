package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"platform/log"
	"strings"
)

var STATEMENTS = make(map[string]*sql.Stmt)

func LoadStatements(path string) {
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

	stmts := strings.Split(string(data), "\n-- ")
	if strings.HasPrefix(stmts[0], "-- ") {
		stmts[0] = stmts[0][3:]
	} else {
		stmts = stmts[1:]
	}

	for _, stmt := range stmts {
		tmp := strings.SplitN(stmt, "\n", 2)
		name := strings.TrimSpace(tmp[0])
		statementStr := strings.TrimSpace(tmp[1])
		statement, err := db.Prepare(statementStr)
		if err != nil {
			log.Errorf("Error preparing statement '%s': %v", name, err)
			continue
		}
		STATEMENTS[name] = statement
	}

	log.Noticef("Statements loaded successfully")
}

func GetStatement(name string) (*sql.Stmt, error) {
	if stmt, ok := STATEMENTS[name]; ok {
		return stmt, nil
	}
	return nil, fmt.Errorf("statement not found: %s", name)
}
