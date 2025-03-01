-- GetConfig
SELECT value
	FROM config
	WHERE key = ?;

-- GetConfigs
SELECT *
	FROM config;

-- SetConfig
UPDATE config
	SET value = ?
	WHERE key = ?;

-- GetKey
SELECT key
	FROM keys
	WHERE name = ?;

-- ChallengeExistsID
SELECT name
	FROM challenges
	WHERE id = ?;

-- ChallengeExistsName
SELECT id
	FROM challenges
	WHERE name = ?;

-- GetChallengeName
SELECT name
	FROM challenges
	WHERE id = ?;

-- GetChallenges
SELECT *
	FROM challenges
	ORDER BY points;

-- CreateChallenge
INSERT INTO challenges (name, description, difficulty, points, max_points, solves, host, port, category, files, flag, hint1, hint2, hidden, is_extra)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- DeleteChallenge
DELETE FROM challenges
	WHERE name = ?;

-- UpdateChallenge
UPDATE challenges
	SET name = ?,
		description = ?,
		difficulty = ?,
		max_points = ?,
		host = ?,
		port = ?,
		category = ?,
		files = ?,
		flag = ?,
		hint1 = ?,
		hint2 = ?,
		hidden = ?,
		is_extra = ?
	WHERE id = ?;

-- UserExists
SELECT id
	FROM users
	WHERE username = ?;

-- EmailExists
SELECT id
	FROM users
	WHERE email = ?;

-- UpdatePassword
UPDATE users
	SET salt = ?,
		password = ?,
		apikey = ?
	WHERE username=?;

-- CreateUser
INSERT INTO users (username, email, salt, password, apikey, score, is_admin)
	VALUES (?, ?, ?, ?, ?, 0, 0);

-- LoginUser
SELECT apikey, salt, password
	FROM users
	WHERE username = ?;

-- FlagExists
SELECT id
	FROM challenges
	WHERE flag = ?;

-- GetUsersScores
SELECT u.id, u.username, u.score
	FROM users AS u
	LEFT JOIN (
		SELECT userid, MAX(timestamp) AS last_submission
		FROM submissions
		GROUP BY userid
	) AS s ON u.id = s.userid
	WHERE u.is_admin = 0
	ORDER BY u.score DESC, s.last_submission ASC;

-- GetUsersBadges
SELECT u.id, b.name, b.desc, b.extra
	FROM users AS u, badges AS b
	WHERE u.id = b.userid
		AND b.name <> 'Intro';

-- GetUserSolves
SELECT c.name, c.category, c.is_extra, s.timestamp
	FROM solves AS s, challenges AS c
	WHERE s.chalid = c.id
		AND s.userid = ?;

-- IsChallengeSolved
SELECT *
	FROM solves
	WHERE userid = ?
		AND chalid = ?;

-- GetSubmissions
SELECT u.username, c.name, s.status, s.flag, s.timestamp
	FROM users AS u, challenges AS c, submissions AS s
	WHERE s.userid = u.id
		AND s.chalid = c.id
	ORDER BY s.timestamp DESC;

-- GetChallIfCorrectFlag
SELECT name, solves
	FROM challenges
	WHERE id = ?
		AND flag = ?;

-- SubmitFlag
INSERT INTO submissions (userid, chalid, status, flag, timestamp)
	VALUES (?, ?, ?, ?, ?);

-- InsertSolve
INSERT INTO solves (userid, chalid, timestamp)
	VALUES (?, ?, ?);

-- GetUserByAPIKey
SELECT id, username, email, score, is_admin
	FROM users
	WHERE apikey = ?;

-- GetUserByUsername
SELECT id, username, email, score, is_admin
	FROM users
	WHERE username = ?;

-- GetUsers
SELECT username, email, score, is_admin
	FROM users
	ORDER BY username;

-- GetGraphData
SELECT u.username, c.points, s.timestamp
	FROM users AS u, solves AS s, challenges AS c
	WHERE u.is_admin = 0
		AND s.userid = u.id
		AND s.chalid = c.id
	ORDER BY s.timestamp;
