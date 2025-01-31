package main

import (
	"flag"
	"os"
	"platform/db"
	"platform/log"
	"platform/routes"
)

type Flags struct {
	migrate bool
	clean   bool
	prune   bool
	dev     bool
}

func main() {
	// Flags
	var flags Flags
	flag.BoolVar(&flags.migrate, "m", false, "Migrate the database")
	flag.BoolVar(&flags.clean, "c", false, "Clean the userland database (non admin users, solves, submissions)")
	flag.BoolVar(&flags.prune, "p", false, "Clean ALL the database")
	flag.BoolVar(&flags.dev, "d", false, "Enables Dev Mode")
	flag.Parse()
	if _, err := os.Stat("DEV"); err == nil || os.Getenv("DEV") != "" {
		flags.dev = true
	}

	// DB
	db.InitDB("database.db")
	defer db.CloseDB()

	if flags.dev {
		log.Notice("Dev mode enabled")
	}
	if flags.migrate {
		db.DropTables()
		db.ExecSQLFile("db/schema.sql")
		db.ExecSQLFile("db/triggers.sql")
		return
	} else if flags.clean {
		db.CleanDB()
		return
	} else if flags.prune {
		db.PruneDB()
		return
	}
	db.LoadStatements("db/statements.sql")

	// Server
	routes.StartRouting([]byte("GrazieDarioGrazieDarioGrazieDP_1"))
}
