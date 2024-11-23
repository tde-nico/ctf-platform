DROP TRIGGER IF EXISTS chall_update_points_newsub;
DROP TRIGGER IF EXISTS chall_update_points_delsub;
DROP TRIGGER IF EXISTS chall_update_points_newconf;
DROP TRIGGER IF EXISTS chall_change_points;
DROP TRIGGER IF EXISTS users_update_points;


CREATE TRIGGER IF NOT EXISTS chall_update_points_newsub
AFTER INSERT ON submissions WHEN NEW.status = 'c'
BEGIN
  UPDATE challenges SET
    solves = (
      SELECT count(*) FROM submissions
      WHERE submissions.status = 'c'
      AND submissions.chalid = NEW.chalid
      AND (SELECT is_admin FROM users WHERE users.id = userid) == 0
    )
  WHERE challenges.id = NEW.chalid;
  UPDATE challenges SET
    points = (
	  CASE WHEN challenges.max_points <= CAST((SELECT value FROM config WHERE key = 'chall-min-points') AS INT)
	  THEN
		challenges.max_points
	  ELSE
		(SELECT max(
		  CAST( (SELECT value FROM config WHERE key = 'chall-min-points') AS INT),
          CAST( (((SELECT value FROM config WHERE key = 'chall-min-points')-max_points)  / power((SELECT value FROM config WHERE key = 'chall-points-decay'), 2) * power(solves, 2) + max_points ) AS INT)
        ))
	  END
    )
  WHERE challenges.id = NEW.chalid;
END;

CREATE TRIGGER IF NOT EXISTS chall_update_points_delsub
AFTER DELETE ON submissions
BEGIN
  UPDATE challenges SET
    solves = (
      SELECT count(*) FROM submissions
      WHERE submissions.status = 'c'
      AND submissions.chalid = OLD.chalid
      AND (SELECT is_admin FROM users WHERE users.id = userid) == 0
    )
  WHERE challenges.id = OLD.chalid;
  UPDATE challenges SET
    points = (
	  CASE WHEN challenges.max_points <= CAST((SELECT value FROM config WHERE key = 'chall-min-points') AS INT)
	  THEN
		challenges.max_points
	  ELSE
		(SELECT max(
		  CAST( (SELECT value FROM config WHERE key = 'chall-min-points') AS INT),
          CAST( (((SELECT value FROM config WHERE key = 'chall-min-points')-max_points)  / power((SELECT value FROM config WHERE key = 'chall-points-decay'), 2) * power(solves, 2) + max_points ) AS INT)
        ))
	  END
    )
  WHERE challenges.id = OLD.chalid;
END;

CREATE TRIGGER IF NOT EXISTS chall_change_points
AFTER UPDATE ON challenges WHEN NEW.max_points != OLD.max_points
BEGIN
  UPDATE challenges SET
    solves = (
      SELECT count(*) FROM submissions
      WHERE submissions.status = 'c'
      AND submissions.chalid = NEW.id
      AND (SELECT is_admin FROM users WHERE users.id = userid) == 0
    )
  WHERE challenges.id = NEW.id;
  UPDATE challenges SET
    points = (
	  CASE WHEN challenges.max_points <= CAST((SELECT value FROM config WHERE key = 'chall-min-points') AS INT)
	  THEN
		challenges.max_points
	  ELSE
		(SELECT max(
		  CAST( (SELECT value FROM config WHERE key = 'chall-min-points') AS INT),
          CAST( (((SELECT value FROM config WHERE key = 'chall-min-points')-max_points)  / power((SELECT value FROM config WHERE key = 'chall-points-decay'), 2) * power(solves, 2) + max_points ) AS INT)
        ))
	  END
    )
  WHERE challenges.id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS chall_update_points_newconf
AFTER UPDATE ON config WHEN NEW.key IN ('chall-min-points', 'chall-points-decay')
BEGIN
  UPDATE challenges SET
    points = (
	  CASE WHEN challenges.max_points <= CAST((SELECT value FROM config WHERE key = 'chall-min-points') AS INT)
	  THEN
		challenges.max_points
	  ELSE
		(SELECT max(
		  CAST( (SELECT value FROM config WHERE key = 'chall-min-points') AS INT),
          CAST( (((SELECT value FROM config WHERE key = 'chall-min-points')-max_points)  / power((SELECT value FROM config WHERE key = 'chall-points-decay'), 2) * power(solves, 2) + max_points ) AS INT)
        ))
	  END
    );
END;

CREATE TRIGGER IF NOT EXISTS users_update_points
AFTER UPDATE ON challenges
BEGIN
  UPDATE users SET score = (
    SELECT IFNULL(CAST(SUM(challenges.points) as INT), 0) FROM challenges WHERE challenges.id IN (
      SELECT chalid FROM submissions
      WHERE submissions.status = 'c'
      AND submissions.userid = users.id
    )
  );
END;