package db

type User struct {
	Id       int
	Username string
	Email    string
	Salt     string
	Password string
	Apikey   string
	Score    int
	IsAdmin  bool
}
