DROP TRIGGER IF EXISTS chall_update_points_newsub;
DROP TRIGGER IF EXISTS chall_update_points_delsub;
DROP TRIGGER IF EXISTS chall_update_points_newconf;
DROP TRIGGER IF EXISTS chall_change_points;
DROP TRIGGER IF EXISTS users_update_points;


CREATE TRIGGER IF NOT EXISTS chall_update_points_newsub
AFTER INSERT ON submissions
WHEN NEW.status = 'c'
BEGIN
  UPDATE challenges
  SET solves = (
      SELECT COUNT(*)
        FROM submissions
        WHERE submissions.status = 'c'
          AND submissions.chalid = NEW.chalid
          AND (SELECT is_admin FROM users WHERE users.id = submissions.userid) == 0
    )
  WHERE challenges.id = NEW.chalid;

  UPDATE challenges
  SET points = (
      SELECT
        CASE
          WHEN max_points <= min_points
            THEN max_points
          ELSE
            MAX(
              min_points,
              CAST(((min_points - max_points) / (decay * decay) * (solves * solves) + max_points) AS INT)
            )
        END
      FROM (
        SELECT
          CAST((SELECT value FROM config WHERE key = 'chall-min-points') AS INT) AS min_points,
          CAST((SELECT value FROM config WHERE key = 'chall-points-decay') AS REAL) AS decay
      )
    )
  WHERE challenges.id = NEW.chalid;
END;

CREATE TRIGGER IF NOT EXISTS chall_update_points_delsub
AFTER DELETE ON submissions
BEGIN
  UPDATE challenges
  SET
    solves = (
      SELECT COUNT(*)
        FROM submissions
        WHERE submissions.status = 'c'
          AND submissions.chalid = OLD.chalid
          AND (SELECT is_admin FROM users WHERE users.id = submissions.userid) == 0
    )
  WHERE challenges.id = OLD.chalid;
  
  UPDATE challenges
  SET points = (
      SELECT
        CASE
          WHEN max_points <= min_points
            THEN max_points
          ELSE
            MAX(
              min_points,
              CAST(((min_points - max_points) / (decay * decay) * (solves * solves) + max_points) AS INT)
            )
        END
      FROM (
        SELECT
          CAST((SELECT value FROM config WHERE key = 'chall-min-points') AS INT) AS min_points,
          CAST((SELECT value FROM config WHERE key = 'chall-points-decay') AS REAL) AS decay
      )
    )
  WHERE challenges.id = OLD.chalid;
END;

CREATE TRIGGER IF NOT EXISTS chall_change_points
AFTER UPDATE ON challenges
WHEN NEW.max_points != OLD.max_points
BEGIN
  UPDATE challenges
  SET
    solves = (
      SELECT COUNT(*)
        FROM submissions
        WHERE submissions.status = 'c'
          AND submissions.chalid = NEW.id
          AND (SELECT is_admin FROM users WHERE users.id = submissions.userid) == 0
    )
  WHERE challenges.id = NEW.id;
  
  UPDATE challenges
  SET points = (
      SELECT
        CASE
          WHEN max_points <= min_points
            THEN max_points
          ELSE
            MAX(
              min_points,
              CAST(((min_points - max_points) / (decay * decay) * (solves * solves) + max_points) AS INT)
            )
        END
      FROM (
        SELECT
          CAST((SELECT value FROM config WHERE key = 'chall-min-points') AS INT) AS min_points,
          CAST((SELECT value FROM config WHERE key = 'chall-points-decay') AS REAL) AS decay
      )
    )
  WHERE challenges.id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS chall_update_points_newconf
AFTER UPDATE ON config
WHEN NEW.key IN ('chall-min-points', 'chall-points-decay')
BEGIN
  UPDATE challenges
    SET points = (
      SELECT
	      CASE
        WHEN max_points <= min_points
	        THEN max_points
	      ELSE
		      max(
		        min_points,
            CAST(((min_points - max_points) / (decay * decay) * (solves * solves) + max_points) AS INT)
          )
	      END
      FROM (
        SELECT
          CAST((SELECT value FROM config WHERE key = 'chall-min-points') AS INT) AS min_points,
          CAST((SELECT value FROM config WHERE key = 'chall-points-decay') AS REAL) AS decay
      )
    );
END;

CREATE TRIGGER IF NOT EXISTS users_update_points
AFTER UPDATE ON challenges
BEGIN
  UPDATE users
  SET score = (
    SELECT COALESCE(SUM(challenges.points), 0)
    FROM challenges 
    WHERE challenges.id IN (
      SELECT chalid
      FROM submissions
      WHERE submissions.status = 'c'
        AND submissions.userid = users.id
    )
  );
END;


DROP TRIGGER IF EXISTS update_badges_on_solve_insert;
DROP TRIGGER IF EXISTS update_badges_on_solve_delete;
DROP TRIGGER IF EXISTS update_badges_on_challenge_update;
DROP TRIGGER IF EXISTS update_badges_on_challenge_insert;
DROP TRIGGER IF EXISTS update_badges_on_challenge_delete;

CREATE TRIGGER IF NOT EXISTS update_badges_on_solve_insert
AFTER INSERT ON solves
BEGIN
  INSERT OR IGNORE INTO badges (name, desc, extra, userid)
  SELECT
      c.category,
      'Completed all ' || c.category || ' challenges',
      0,
      NEW.userid
    FROM challenges AS c
      WHERE c.id = NEW.chalid
        AND (
          SELECT COUNT(*)
            FROM challenges AS c2
            WHERE c2.category = c.category
              AND c2.is_extra = 0
              AND c2.hidden = 0
              AND NOT EXISTS (
                SELECT 1
                  FROM solves AS s
                  WHERE s.chalid = c2.id
                    AND s.userid = NEW.userid
              )
        ) = 0;

  UPDATE badges
    SET extra = 1
    WHERE userid = NEW.userid
      AND extra = 0
      AND (
        SELECT COUNT(*)
          FROM challenges AS c
          WHERE c.category = badges.name
            AND c.hidden = 0
            AND c.is_extra = 1
      ) > 0
      AND (
        SELECT COUNT(*)
          FROM challenges AS c
          WHERE c.category = badges.name
            AND c.hidden = 0
            AND NOT EXISTS (
              SELECT 1
                FROM solves AS s
                WHERE s.chalid = c.id
                  AND s.userid = NEW.userid
              )
      ) = 0;
END;

CREATE TRIGGER IF NOT EXISTS update_badges_on_solve_delete
AFTER DELETE ON solves
BEGIN
  DELETE FROM badges
    WHERE userid = OLD.userid
      AND name = (
        SELECT category
          FROM challenges
          WHERE id = OLD.chalid
      );

  INSERT OR IGNORE INTO badges (name, desc, extra, userid)
  SELECT
      c.category,
      'Completed all ' || c.category || ' challenges',
      0,
      OLD.userid
    FROM challenges AS c
      WHERE c.id = OLD.chalid
        AND (
          SELECT COUNT(*)
            FROM challenges AS c2
            WHERE c2.category = c.category
              AND c2.is_extra = 0
              AND c2.hidden = 0
              AND NOT EXISTS (
                SELECT 1
                  FROM solves AS s
                  WHERE s.chalid = c2.id
                    AND s.userid = OLD.userid
              )
        ) = 0;
  
  UPDATE badges
    SET extra = 1
    WHERE userid = OLD.userid
      AND extra = 0
      AND (
        SELECT COUNT(*)
          FROM challenges AS c
          WHERE c.category = badges.name
            AND c.hidden = 0
            AND c.is_extra = 1
      ) > 0
      AND (
        SELECT COUNT(*)
          FROM challenges AS c
          WHERE c.category = badges.name
            AND c.hidden = 0
            AND NOT EXISTS (
              SELECT 1
                FROM solves AS s
                WHERE s.chalid = c.id
                  AND s.userid = OLD.userid
              )
      ) = 0;
END;

CREATE TRIGGER IF NOT EXISTS update_badges_on_challenge_update
AFTER UPDATE ON challenges
WHEN NEW.category <> OLD.category OR NEW.is_extra <> OLD.is_extra OR NEW.hidden <> OLD.hidden
BEGIN
  DELETE FROM badges
    WHERE name = OLD.category;

  INSERT OR IGNORE INTO badges (name, desc, extra, userid)
  SELECT
      NEW.category,
      'Completed all ' || NEW.category || ' challenges',
      0,
      u.id
    FROM users AS u
      WHERE (
          SELECT COUNT(*)
            FROM challenges AS c
            WHERE c.category = NEW.category
              AND c.is_extra = 0
              AND c.hidden = 0
              AND NOT EXISTS (
                SELECT 1
                  FROM solves AS s
                  WHERE s.chalid = c.id
                    AND s.userid = u.id
              )
        ) = 0;

  UPDATE badges
    SET extra = 1
    WHERE name = NEW.category
      AND extra = 0
      AND (
        SELECT COUNT(*)
          FROM challenges AS c
          WHERE c.category = badges.name
            AND c.hidden = 0
            AND c.is_extra = 1
      ) > 0
      AND (
        SELECT COUNT(*)
          FROM challenges AS c
          WHERE c.category = badges.name
            AND c.hidden = 0
            AND NOT EXISTS (
              SELECT 1
                FROM solves AS s
                WHERE s.chalid = c.id
                  AND s.userid = badges.userid
              )
      ) = 0;
END;

CREATE TRIGGER IF NOT EXISTS update_badges_on_challenge_insert
AFTER INSERT ON challenges
BEGIN
  DELETE FROM badges
    WHERE name = NEW.category
      AND NEW.is_extra = 0;

  UPDATE badges
    SET extra = 0
    WHERE name = NEW.category;
END;

CREATE TRIGGER IF NOT EXISTS update_badges_on_challenge_delete
AFTER DELETE ON challenges
BEGIN
  DELETE FROM badges
    WHERE name = OLD.category;

  INSERT OR IGNORE INTO badges (name, desc, extra, userid)
  SELECT
      OLD.category,
      'Completed all ' || OLD.category || ' challenges',
      0,
      u.id
    FROM users AS u
      WHERE (
          SELECT COUNT(*)
            FROM challenges AS c
            WHERE c.category = OLD.category
              AND c.is_extra = 0
              AND c.hidden = 0
              AND NOT EXISTS (
                SELECT 1
                  FROM solves AS s
                  WHERE s.chalid = c.id
                    AND s.userid = u.id
              )
        ) = 0;

  UPDATE badges
    SET extra = 1
    WHERE name = OLD.category
      AND extra = 0
      AND (
        SELECT COUNT(*)
          FROM challenges AS c
          WHERE c.category = badges.name
            AND c.hidden = 0
            AND c.is_extra = 1
      ) > 0
      AND (
        SELECT COUNT(*)
          FROM challenges AS c
          WHERE c.category = badges.name
            AND c.hidden = 0
            AND NOT EXISTS (
              SELECT 1
                FROM solves AS s
                WHERE s.chalid = c.id
                  AND s.userid = badges.userid
              )
      ) = 0;
END;
