package db

import "time"

type Config struct {
	Key   string
	Type  string
	Value string
	Desc  string
}

type User struct {
	ID       int
	Username string
	Email    string
	ApiKey   string
	Score    int
	IsAdmin  bool
}

type Challenge struct {
	ID          int
	Name        string
	Description string
	Difficulty  string
	Points      int
	MaxPoints   int
	Solves      int
	Host        string
	Port        string
	Category    string
	Files       string
	Flag        string
	Hint1       string
	Hint2       string
	Hidden      bool
	IsExtra     bool
}

type Solve struct {
	UserID       int
	ChalName     string
	ChalCategory string
	ChalExtra    bool
	Timestamp    time.Time
}

type Submission struct {
	UserUsername string
	ChalName     string
	Status       string
	Flag         string
	Timestamp    time.Time
}

type Badges struct {
	Name  string
	Char  string
	Extra bool
}

type UserScore struct {
	Username string
	Score    int
	Badges   []Badges
}

type GraphData struct {
	User      string
	Points    int
	Timestamp time.Time
}

var CATEGORIES = []string{
	"Web",
	"Crypto",
	"Rev",
	"Pwn",
}

var DIFFICULTIES = []string{
	"Easy",
	"Medium",
	"Hard",
}
