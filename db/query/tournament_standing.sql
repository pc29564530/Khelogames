-- name: CreateTournamentStanding :one
INSERT INTO tournament_standing (
    tournament_id,
    group_id,
    team_id,
    wins,
    loss,
    draw,
    goal_for,
    goal_against,
    goal_difference,
    points
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ) RETURNING *;


-- name: GetTournamentStanding :many
SELECT 
    ts.standing_id, ts.tournament_id, ts.group_id, ts.team_id,
    ts.wins, ts.loss, ts.draw, ts.goal_for, ts.goal_against, ts.goal_difference, ts.points,
    t.tournament_name, t.sports,
    c.name
FROM 
    tournament_standing ts
JOIN 
    group_league tg ON ts.group_id = tg.group_id
JOIN 
    tournaments t ON ts.tournament_id = t.id
JOIN 
    teams c ON ts.team_id = c.id
WHERE 
    ts.tournament_id = $1
    AND tg.group_id = $2
    AND t.sports = $3;



-- name: UpdateTournamentStanding :one
UPDATE tournament_standing AS ts
SET 
    goal_for = COALESCE((
        SELECT SUM(CASE 
            WHEN ms.home_team_id = ts.team_id THEN fs.goals
            WHEN ms.away_team_id = ts.team_id THEN fs.goals
            ELSE 0
        END)
        FROM football_score AS fs
        JOIN matches AS ms ON fs.match_id = ms.id
        WHERE fs.team_id = ts.team_id
    ), 0),
    goal_against = COALESCE((
        SELECT SUM(CASE 
            WHEN ms.home_team_id = ts.team_id THEN (
                SELECT SUM(fs.goals) 
                FROM football_score AS fs 
                WHERE fs.match_id = ms.id AND fs.team_id = ms.away_team_id
            )
            WHEN ms.away_team_id = ts.team_id THEN (
                SELECT SUM(fs2.goals) 
                FROM football_score AS fs2
                WHERE fs2.match_id = ms.id AND fs2.team_id = ms.home_team_id
            )
        END)
        FROM matches AS ms
        WHERE ms.home_team_id = ts.team_id OR ms.away_team_id = ts.team_id
    ), 0),
    goal_difference = COALESCE(goal_for, 0) - COALESCE(goal_against, 0),
    wins = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN football_score AS fs_home ON ms.id = fs_home.match_id AND ms.home_team_id = ts.team_id
        LEFT JOIN football_score AS fs_away ON ms.id = fs_away.match_id AND ms.away_team_id = fs_away.team_id
        WHERE (ms.home_team_id = ts.team_id AND fs_home.goals > fs_away.goals)
        OR (ms.away_team_id = ts.team_id AND fs_away.goals > fs_home.goals)
    ), 0),
    loss = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN football_score fs_home ON ms.id = fs_home.match_id AND ms.home_team_id = fs_home.team_id
        LEFT JOIN football_score fs_away ON ms.id = fs_away.match_id AND ms.away_team_id = fs_away.team_id
        WHERE (ms.home_team_id = ts.team_id AND fs_home.goals < fs_away.goals)
        OR (ms.away_team_id = ts.team_id AND fs_away.goals < fs_home.goals)
    ), 0),
    draw = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN football_score AS fs_home ON ms.id = fs_home.match_id AND ms.home_team_id = ts.team_id
        LEFT JOIN football_score AS fs_away ON ms.id = fs_away.match_id AND ms.away_team_id = fs_away.team_id
        WHERE (ms.home_team_id = ts.team_id AND fs_home.goals = fs_away.goals)
        OR (ms.away_team_id = ts.team_id AND fs_away.goals = fs_home.goals)
    ), 0),
    points = ((wins * 3) + draw)
WHERE ts.tournament_id = $1 AND ts.team_id = $2
RETURNING *;
