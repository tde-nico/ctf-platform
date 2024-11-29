package db

type User struct {
	Id       int
	Username string
	Email    string
	Apikey   string
	Score    int
	IsAdmin  bool
}

type Challenge struct {
	Id          int
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
	UserId    int
	ChalId    int
	Timestamp string
}

type Submission struct {
	Id        int
	UserId    int
	ChalId    int
	Status    rune
	Flag      string
	Timestamp string
}
