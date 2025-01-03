package db

import (
	"errors"
	"fmt"
	"platform/log"
	"platform/utils"
)

const APIKEY_LENGTH = 32
const SALT_LENGTH = 8

const (
	StatusWrongFlag int = iota
	StatusAlreadySolved
	StatusCorrectFlag
)

var ErrDatabaseNotInitialized = errors.New("database not initialized")

func UserExists(username string) bool {
	if db == nil {
		log.Error("Database not initialized")
		return false
	}

	rows, err := db.Query("SELECT id FROM users WHERE username = ?", username)
	if err != nil {
		log.Errorf("Error querying users for username: %v", err)
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func EmailExists(email string) bool {
	if db == nil {
		log.Error("Database not initialized")
		return false
	}

	rows, err := db.Query("SELECT id FROM users WHERE email = ?", email)
	if err != nil {
		log.Errorf("Error querying users for mail: %v", err)
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func createUser(username, email, salt, secret, apiKey string) error {
	if db == nil {
		log.Error("Database not initialized")
		return ErrDatabaseNotInitialized
	}

	_, err := db.Exec("INSERT INTO users (username, email, salt, password, apikey, score, is_admin) VALUES (?, ?, ?, ?, ?, 0, 0)", username, email, salt, secret, apiKey)
	if err != nil {
		log.Errorf("Error inserting user: %v", err)
		return err
	}

	return nil
}

func RegisterUser(username, email, password string) error {
	_, apiKeyHex, err := utils.GetRand(APIKEY_LENGTH)
	if err != nil {
		return err
	}
	salt, saltHex, err := utils.GetRand(SALT_LENGTH)
	if err != nil {
		return err
	}
	secretHex := utils.HashPassword(password, salt)
	err = createUser(username, email, string(saltHex), string(secretHex), string(apiKeyHex))
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(username, password string) (string, error) {
	if db == nil {
		return "", ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT apikey, salt, password FROM users WHERE username = ?", username)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		log.Errorf("User not found")
		return "", fmt.Errorf("user not found")
	}

	var apikey, saltHex, secretHex string
	err = rows.Scan(
		&apikey,
		&saltHex,
		&secretHex,
	)
	if err != nil {
		return "", err
	}

	salt, err := utils.HexToBytes(saltHex)
	if err != nil {
		return "", err
	}

	hash := utils.HashPassword(password, salt)
	if secretHex != hash {
		return "", fmt.Errorf("invalid password")
	}

	return apikey, nil
}

func GetUserByAPIKey(apiKey string) (*User, error) {
	if db == nil {
		return nil, ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT id, username, email, score, is_admin FROM users WHERE apikey = ?", apiKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	user := User{ApiKey: apiKey}
	err = rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Score,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	if db == nil {
		return nil, ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT id, email, score, is_admin FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	user := User{Username: username}
	err = rows.Scan(
		&user.ID,
		&user.Email,
		&user.Score,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetSolvesByUser(user *User) ([]*Solve, error) {
	if db == nil {
		return nil, ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT c.name, c.category, c.is_extra, s.timestamp FROM solves AS s, challenges AS c WHERE s.chalid = c.id AND s.userid = ?", user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	solves := make([]*Solve, 0)
	for rows.Next() {
		var solve Solve
		var timestamp *string
		err = rows.Scan(&solve.ChalName, &solve.ChalCategory, &solve.ChalExtra, &timestamp)
		if err != nil {
			return nil, err
		}
		solve.Timestamp, err = utils.ParseTime(timestamp)
		if err != nil {
			return nil, err
		}
		solves = append(solves, &solve)
	}

	return solves, nil
}

func challengeEsits(ID int) (bool, error) {
	if db == nil {
		return false, ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT id FROM challenges WHERE id = ?", ID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func isChallengeSolved(user *User, chalID int) (bool, error) {
	if db == nil {
		return false, ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT * FROM solves WHERE userid = ? AND chalid = ?", user.ID, chalID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func getChallIfCorrectFlag(chalID int, flag string) (*Challenge, error) {
	if db == nil {
		return nil, ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT name, solves FROM challenges WHERE id = ? AND flag = ?", chalID, flag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var chal Challenge
	err = rows.Scan(&chal.Name, &chal.Solves)
	if err != nil {
		return nil, err
	}

	return &chal, nil
}

func SubmitFlag(user *User, chalID int, flag string) (int, error) {
	if db == nil {
		return StatusWrongFlag, ErrDatabaseNotInitialized
	}

	now := utils.CurrentTime()

	status, err := challengeEsits(chalID)
	if err != nil {
		return StatusWrongFlag, err
	}
	if !status {
		return StatusWrongFlag, fmt.Errorf("challenge not found")
	}

	status, err = isChallengeSolved(user, chalID)
	if err != nil {
		return StatusWrongFlag, err
	}
	if status {
		_, err := db.Exec("INSERT INTO submissions (userid, chalid, status, flag, timestamp) VALUES (?, ?, ?, ?, ?)", user.ID, chalID, "r", flag, now)
		if err != nil {
			return StatusWrongFlag, err
		}

		return StatusAlreadySolved, fmt.Errorf("challenge already solved")
	}

	chal, err := getChallIfCorrectFlag(chalID, flag)
	if err != nil {
		return StatusWrongFlag, err
	}
	if chal == nil {
		_, err := db.Exec("INSERT INTO submissions (userid, chalid, status, flag, timestamp) VALUES (?, ?, ?, ?, ?)", user.ID, chalID, "w", flag, now)
		if err != nil {
			return StatusWrongFlag, err
		}

		return StatusWrongFlag, nil
	}

	if chal.Solves == 0 && !user.IsAdmin {
		log.Noticef("First Blood on %s from %s", chal.Name, user.Username)
		// TODO: bot first blood
	}

	_, err = db.Exec("INSERT INTO submissions (userid, chalid, status, flag, timestamp) VALUES (?, ?, ?, ?, ?)", user.ID, chalID, "c", flag, now)
	if err != nil {
		log.Warningf("%d", 262)
		return StatusWrongFlag, err
	}

	_, err = db.Exec("INSERT INTO solves (userid, chalid, timestamp) VALUES (?, ?, ?)", user.ID, chalID, now)
	if err != nil {
		log.Warningf("%d", 268)
		return StatusWrongFlag, err
	}

	return StatusCorrectFlag, nil
}

func getVisibleChallsByCategory() (map[string][]int, error) {
	if db == nil {
		return nil, ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT category, COUNT(*), SUM(is_extra) FROM challenges WHERE hidden=0 GROUP BY category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challs := make(map[string][]int)
	for rows.Next() {
		var name string
		var count, extra int
		err = rows.Scan(&name, &count, &extra)
		if err != nil {
			return nil, err
		}
		challs[name] = []int{count - extra, extra}
	}

	return challs, nil

}

func GetScoreUsers() ([]UserScore, error) {
	if db == nil {
		return nil, ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT id, username, score FROM users WHERE is_admin=0 ORDER BY score DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Score)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	challs, err := getVisibleChallsByCategory()
	if err != nil {
		return nil, err
	}

	scoreUsers := make([]UserScore, len(users))
	for i, user := range users {
		scoreUsers[i].Username = user.Username
		scoreUsers[i].Score = user.Score

		solves, err := GetSolvesByUser(user)
		if err != nil {
			return nil, err
		}

		counter := make(map[string][]int)
		for _, s := range solves {
			c, ok := counter[s.ChalCategory]
			if !ok {
				c = []int{0, 0}
			}
			if s.ChalExtra {
				c[1]++
			} else {
				c[0]++
			}
			counter[s.ChalCategory] = c
		}

		for category, counts := range counter {
			if category == "Intro" {
				continue
			}

			challCounts, ok := challs[category]
			if !ok {
				continue
			}
			if challCounts[0] == counts[0] {
				scoreUsers[i].Badges = append(scoreUsers[i].Badges, Badges{
					Name:  category,
					Char:  string(category[0]),
					Extra: false,
				})
				if challCounts[1] > 0 && challCounts[1] == counts[1] {
					scoreUsers[i].Badges[len(scoreUsers[i].Badges)-1].Extra = true
				}
			}
		}

	}

	return scoreUsers, nil
}

func GetGraphData() ([]GraphData, error) {
	if db == nil {
		return nil, ErrDatabaseNotInitialized
	}

	rows, err := db.Query("SELECT u.username, c.points, s.timestamp FROM users AS u, solves AS s, challenges AS c WHERE u.is_admin=0 AND s.userid=u.id AND s.chalid=c.id ORDER BY s.timestamp")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]GraphData, 0)
	for rows.Next() {
		var tmp GraphData
		var timestamp *string
		err = rows.Scan(&tmp.User, &tmp.Points, &timestamp)
		if err != nil {
			return nil, err
		}
		tmp.Timestamp, err = utils.ParseTime(timestamp)
		if err != nil {
			return nil, err
		}
		data = append(data, tmp)
	}

	return data, nil
}
