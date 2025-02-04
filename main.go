package main

import (
	"encoding/hex"
	"flag"
	"os"
	"platform/db"
	"platform/log"
	"platform/routes"
)

type Flags struct {
	migrate  bool
	triggers bool
	clean    bool
	prune    bool
	dev      bool
}

func InitFlags() *Flags {
	var flags Flags
	flag.BoolVar(&flags.migrate, "m", false, "Migrate the database")
	flag.BoolVar(&flags.triggers, "t", false, "Reload Triggers")
	flag.BoolVar(&flags.clean, "c", false, "Clean the userland database (non admin users, solves, submissions)")
	flag.BoolVar(&flags.prune, "p", false, "Clean ALL the database")
	flag.BoolVar(&flags.dev, "d", false, "Enables Dev Mode")
	flag.Parse()
	if _, err := os.Stat("DEV"); err == nil || os.Getenv("DEV") != "" {
		flags.dev = true
	}
	return &flags
}

func EvalFlags(flags *Flags) bool {
	if flags.dev {
		log.Notice("Dev mode enabled")
	}
	if flags.migrate {
		db.DropTables()
		db.ExecSQLFile("db/schema.sql")
		db.ExecSQLFile("db/triggers.sql")
	} else if flags.triggers {
		db.ExecSQLFile("db/triggers.sql")
	} else if flags.clean {
		db.CleanDB()
	} else if flags.prune {
		db.PruneDB()
	} else {
		return true
	}
	return false
}

func LoadSecretKey() []byte {
	var err error

	secretHex, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		log.Info("SECRET_KEY not found in environment, using database")
		secretHex, err = db.GetKey("secret-key")
		if err != nil {
			log.Fatal("SECRET_KEY not found: %v", err)
		}
	}

	if len(secretHex) != 64 {
		log.Fatal("SECRET_KEY must be 32 bytes long")
	}

	secret := make([]byte, 32)
	n, err := hex.Decode(secret, []byte(secretHex))
	if err != nil {
		log.Fatal("Error decoding SECRET_KEY: %v", err)
	}
	if n != 32 {
		log.Fatal("Error decoding SECRET_KEY")
	}

	return secret
}

func main() {
	flags := InitFlags()

	db.InitDB("database.db")
	defer db.CloseDB()
	if !EvalFlags(flags) {
		return
	}
	db.LoadStatements("db/statements.sql")

	key := LoadSecretKey()
	routes.StartRouting(key)
}
