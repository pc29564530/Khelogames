package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const updateCricketInningsStatus = `
WITH resolve_ids AS (
  SELECT 
    m.id AS match_id, 
    t.id AS team_id
  FROM matches m, teams t
  WHERE m.public_id = $1 AND t.public_id = $2
)
UPDATE cricket_score
SET inning_status = 'in_progress'
WHERE match_id = (SELECT match_id FROM resolve_ids)
  AND team_id = (SELECT team_id FROM resolve_ids)
  AND inning_number = $3
RETURNING *;
`

func (q *Queries) UpdateCricketInningsStatus(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, inningNumber int) (*models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketInningsStatus, matchPublicID, teamPublicID, inningNumber)
	var i models.CricketScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.InningNumber,
		&i.Score,
		&i.Wickets,
		&i.Overs,
		&i.RunRate,
		&i.TargetRunRate,
		&i.FollowOn,
		&i.IsInningCompleted,
		&i.Declared,
		&i.InningStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan: ", err)
	}

	return &i, nil

}

const addCricketBall = `
WITH resolve_ids AS (
	SELECT 
		m.id AS match_id, 
		t.id AS team_id, 
		p.id AS player_id
	FROM matches m
	JOIN teams t ON t.public_id = $2
	JOIN players p ON p.public_id = $3
	WHERE m.public_id = $1
)
INSERT INTO bowler_score (
    match_id,
    team_id,
    bowler_id,
	inning_number,
    ball_number,
    runs,
    wickets,
    wide,
    no_ball,
	bowling_status,
	is_current_bowler
)
SELECT
	ri.match_id,
	ri.team_id,
	ri.player_id,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	TRUE,
	TRUE
FROM resolve_ids ri
RETURNING *;
`

type AddCricketBallParams struct {
	MatchPublicID  uuid.UUID `json:"match_public_id"`
	TeamPublicID   uuid.UUID `json:"team_public_id"`
	BowlerPublicID uuid.UUID `json:"bowler_public_id"`
	InningNumber   int       `json:"inning_number"`
	BallNumber     int32     `json:"ball_number"`
	Runs           int32     `json:"runs"`
	Wickets        int32     `json:"wickets"`
	Wide           int32     `json:"wide"`
	NoBall         int32     `json:"no_ball"`
}

func (q *Queries) AddCricketBall(ctx context.Context, arg AddCricketBallParams) (models.BowlerScore, error) {
	row := q.db.QueryRowContext(ctx, addCricketBall,
		arg.MatchPublicID,
		arg.TeamPublicID,
		arg.BowlerPublicID,
		arg.InningNumber,
		arg.BallNumber,
		arg.Runs,
		arg.Wickets,
		arg.Wide,
		arg.NoBall,
	)
	var i models.BowlerScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.BowlerID,
		&i.InningNumber,
		&i.BallNumber,
		&i.Runs,
		&i.Wickets,
		&i.Wide,
		&i.NoBall,
		&i.BowlingStatus,
		&i.IsCurrentBowler,
	)
	return i, err
}

const getCricketStricker = `
	SELECT b.* FROM batsman_score b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN teams t ON t.id = b.team_id
	WHERE m.public_id=$1 AND t.public_id=$2 AND is_currently_batting=true AND is_striker=true AND b.inning_number= $3;
`

func (q *Queries) GetCricketStricker(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, inningNumber int) (*models.BatsmanScore, error) {
	row := q.db.QueryRowContext(ctx, getCricketStricker, matchPublicID, teamPublicID, inningNumber)

	var i models.BatsmanScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.BatsmanID,
		&i.TeamID,
		&i.MatchID,
		&i.InningNumber,
		&i.Position,
		&i.RunsScored,
		&i.BallsFaced,
		&i.Fours,
		&i.Sixes,
		&i.BattingStatus,
		&i.IsStriker,
		&i.IsCurrentlyBatting,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, err
}

const addCricketBatsScore = `
WITH resolve_ids AS (
	SELECT m.id AS match_id, t.id AS team_id, p.id AS player_id FROM matches m, teams t, players p
	WHERE m.public_id = $1 AND t.public_id = $2 AND p.public_id = $3
)
INSERT INTO batsman_score (
    match_id,
    team_id,
	batsman_id,
	inning_number,
    position,
    runs_scored,
    balls_faced,
    fours,
    sixes,
	batting_status,
	is_striker,
	is_currently_batting
)
SELECT
	ri.match_id, ri.team_id, ri.player_id, $4, $5, $6, $7, $8, $9, $10, $11, $12
FROM resolve_ids ri
RETURNING *;
`

type AddCricketBatsScoreParams struct {
	MatchPublicID      uuid.UUID `json:"match_public_id"`
	TeamPublicID       uuid.UUID `json:"team_public_id"`
	BatsmanPublicID    uuid.UUID `json:"batsman_public_id"`
	InningNumber       int       `json:"inning_number"`
	Position           string    `json:"position"`
	RunsScored         int32     `json:"runs_scored"`
	BallsFaced         int32     `json:"balls_faced"`
	Fours              int32     `json:"fours"`
	Sixes              int32     `json:"sixes"`
	BattingStatus      bool      `json:"batting_status"`
	IsStriker          bool      `json:"is_striker"`
	IsCurrentlyBatting bool      `json:"is_currently_batting"`
}

func (q *Queries) AddCricketBatsScore(ctx context.Context, arg AddCricketBatsScoreParams) (models.BatsmanScore, error) {
	row := q.db.QueryRowContext(ctx, addCricketBatsScore,
		arg.MatchPublicID,
		arg.TeamPublicID,
		arg.BatsmanPublicID,
		arg.InningNumber,
		arg.Position,
		arg.RunsScored,
		arg.BallsFaced,
		arg.Fours,
		arg.Sixes,
		arg.BattingStatus,
		arg.IsStriker,
		arg.IsCurrentlyBatting,
	)
	var i models.BatsmanScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.BatsmanID,
		&i.InningNumber,
		&i.Position,
		&i.RunsScored,
		&i.BallsFaced,
		&i.Fours,
		&i.Sixes,
		&i.BattingStatus,
		&i.IsStriker,
		&i.IsCurrentlyBatting,
	)
	return i, err
}

const getCricketBalls = `
SELECT b.* FROM bowler_score b
JOIN matches m ON m.id = b.match_id
JOIN teams t ON t.id = b.team_id
WHERE m.public_id=$1 AND t.public_id = $2
ORDER BY b.id, b.inning_number
`

func (q *Queries) GetCricketBalls(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) ([]models.BowlerScore, error) {
	rows, err := q.db.QueryContext(ctx, getCricketBalls, matchPublicID, teamPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.BowlerScore
	for rows.Next() {
		var i models.BowlerScore
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.TeamID,
			&i.MatchID,
			&i.BowlerID,
			&i.InningNumber,
			&i.BallNumber,
			&i.Runs,
			&i.Wickets,
			&i.Wide,
			&i.NoBall,
			&i.BowlingStatus,
			&i.IsCurrentBowler,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCricketPlayerScore = `
SELECT * FROM batsman_score b
JOIN matches m ON m.id = b.match_id
JOIN players p ON p.id = b.batsman_id
WHERE m.public_id=$1 AND p.public_id=$2 LIMIT 1
`

// Not used function
func (q *Queries) GetCricketPlayerScore(ctx context.Context, matchPublicID, batsmanPublicID uuid.UUID) (models.BatsmanScore, error) {
	row := q.db.QueryRowContext(ctx, getCricketPlayerScore, matchPublicID, batsmanPublicID)
	var i models.BatsmanScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.BatsmanID,
		&i.InningNumber,
		&i.Position,
		&i.RunsScored,
		&i.BallsFaced,
		&i.Fours,
		&i.Sixes,
		&i.BattingStatus,
		&i.IsStriker,
		&i.IsCurrentlyBatting,
	)
	return i, err
}

const getCricketPlayersScore = `
SELECT b.* FROM batsman_score b
JOIN matches m ON m.id = b.match_id
JOIN teams t ON t.id = b.team_id
WHERE m.public_id = $1 AND t.public_id = $2
ORDER BY b.id, b.inning_number
`

func (q *Queries) GetCricketPlayersScore(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) ([]models.BatsmanScore, error) {
	rows, err := q.db.QueryContext(ctx, getCricketPlayersScore, matchPublicID, teamPublicID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var items []models.BatsmanScore
	for rows.Next() {
		var i models.BatsmanScore
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.MatchID,
			&i.TeamID,
			&i.BatsmanID,
			&i.InningNumber,
			&i.Position,
			&i.RunsScored,
			&i.BallsFaced,
			&i.Fours,
			&i.Sixes,
			&i.BattingStatus,
			&i.IsStriker,
			&i.IsCurrentlyBatting,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// not used function
const getCricketWicket = `
SELECT * FROM wickets
WHERE match_id=$1 AND batsman_id=$2 LIMIT 1
`

type GetCricketWicketParams struct {
	MatchID   int64 `json:"match_id"`
	BatsmanID int64 `json:"batsman_id"`
}

func (q *Queries) GetCricketWicket(ctx context.Context, arg GetCricketWicketParams) (models.Wicket, error) {
	row := q.db.QueryRowContext(ctx, getCricketWicket, arg.MatchID, arg.BatsmanID)
	var i models.Wicket
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.BatsmanID,
		&i.BowlerID,
		&i.WicketsNumber,
		&i.WicketType,
		&i.BallNumber,
		&i.FielderID,
		&i.Score,
		&i.InningNumber,
	)
	return i, err
}

const getCricketWickets = `
SELECT * FROM wickets w
JOIN matches m ON m.id = w.match_id
JOIN teams t ON t.id = w.team_id
WHERE m.public_id=$1 AND t.public_id=$2
ORDER BY w.id, w.inning_number
`

const getWickets = `
SELECT json_build_object(
    'id', w.id,
    'public_id', w.public_id,
    'match_id', w.match_id,
    'team_id', w.team_id,
    'batsman_id', w.batsman_id,
    'bowler_id', w.bowler_id,
    'wicket_number', w.wicket_number,
    'wicket_type', w.wicket_type,
    'ball_number', w.ball_number,
    'fielder_id', w.fielder_id,
    'score', w.score,
    'batsman_player', json_build_object(
        'id', bp.id,
        'name', bp.name,
        'slug', bp.slug,
        'shortName', bp.short_name,
        'position', bp.positions,
        'username', bp.username
    ),
    'bowler_player', json_build_object(
        'id', bowp.id,
        'name', bowp.name,
        'slug', bowp.slug,
        'shortName', bowp.short_name,
        'position', bowp.positions,
        'username', bowp.username
    ),
    'fielder_player', CASE 
        WHEN w.fielder_id IS NOT NULL THEN json_build_object(
            'id', fp.id,
            'name', fp.name,
            'slug', fp.slug,
            'shortName', fp.short_name,
            'position', fp.positions,
            'username', fp.username
        )
        ELSE NULL
    END
) AS wicket_data
FROM wickets w
JOIN matches m ON m.id = w.match_id
JOIN teams t ON t.id = w.team_id
JOIN players bp ON bp.id = w.batsman_id
JOIN players bowp ON bowp.id = w.bowler_id
LEFT JOIN players fp ON fp.id = w.fielder_id
WHERE m.public_id = $1 AND t.public_id = $2
ORDER BY w.id, w.innings_number;
`

func (q *Queries) GetWickets(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getWickets, matchPublicID, teamPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var jsonByte []byte
		if err := rows.Scan(&jsonByte); err != nil {
			return nil, err
		}
		var item map[string]interface{}
		err := json.Unmarshal(jsonByte, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) GetCricketWickets(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) ([]models.Wicket, error) {
	rows, err := q.db.QueryContext(ctx, getCricketWickets, matchPublicID, teamPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Wicket
	for rows.Next() {
		var i models.Wicket
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.MatchID,
			&i.TeamID,
			&i.BatsmanID,
			&i.BowlerID,
			&i.WicketsNumber,
			&i.WicketType,
			&i.BallNumber,
			&i.FielderID,
			&i.Score,
			&i.InningNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCricketBowler = `
UPDATE bowler_score
SET 
    ball_number = $1,
    runs = $2,
    wickets = $3,
    wide = $4,
    no_ball = $5
WHERE match_id = $6 AND bowler_id = $7 AND team_id=$8 AND inning_number= $9
RETURNING *;
`

type UpdateCricketBowlerParams struct {
	BallNumber   int32 `json:"ball_number"`
	Runs         int32 `json:"runs"`
	Wickets      int32 `json:"wickets"`
	Wide         int32 `json:"wide"`
	NoBall       int32 `json:"no_ball"`
	MatchID      int64 `json:"match_id"`
	BowlerID     int64 `json:"bowler_id"`
	TeamID       int64 `json:"team_id"`
	InningNumber int   `json:"inning_number"`
}

// Not used function
func (q *Queries) UpdateCricketBowler(ctx context.Context, arg UpdateCricketBowlerParams) (models.BowlerScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketBowler,
		arg.BallNumber,
		arg.Runs,
		arg.Wickets,
		arg.Wide,
		arg.NoBall,
		arg.MatchID,
		arg.BowlerID,
		arg.TeamID,
		arg.InningNumber,
	)
	var i models.BowlerScore
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.MatchID,
		&i.BowlerID,
		&i.BallNumber,
		&i.Runs,
		&i.Wickets,
		&i.Wide,
		&i.NoBall,
		&i.InningNumber,
	)
	return i, err
}

// not used function
const updateCricketRunsScored = `
UPDATE batsman_score
SET runs_scored = $1,
    balls_faced = $2,
    fours = $3,
    sixes = $4
WHERE match_id = $5 AND batsman_id = $6 AND team_id=$7 AND inning_number= $8
RETURNING *
`

type UpdateCricketRunsScoredParams struct {
	RunsScored   int32 `json:"runs_scored"`
	BallsFaced   int32 `json:"balls_faced"`
	Fours        int32 `json:"fours"`
	Sixes        int32 `json:"sixes"`
	MatchID      int64 `json:"match_id"`
	BatsmanID    int64 `json:"batsman_id"`
	TeamID       int64 `json:"team_id"`
	InningNumber int   `json:"inning_number"`
}

func (q *Queries) UpdateCricketRunsScored(ctx context.Context, arg UpdateCricketRunsScoredParams) (models.BatsmanScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketRunsScored,
		arg.RunsScored,
		arg.BallsFaced,
		arg.Fours,
		arg.Sixes,
		arg.MatchID,
		arg.BatsmanID,
		arg.TeamID,
		arg.InningNumber,
	)
	var i models.BatsmanScore
	err := row.Scan(
		&i.ID,
		&i.BatsmanID,
		&i.TeamID,
		&i.MatchID,
		&i.Position,
		&i.RunsScored,
		&i.BallsFaced,
		&i.Fours,
		&i.Sixes,
		&i.BattingStatus,
		&i.IsStriker,
		&i.IsCurrentlyBatting,
		&i.InningNumber,
	)
	return i, err
}

// not used function
const updateBowlingStats = `
UPDATE bowler_score b
SET runs = runs + $1,
    ball_number = ball_number + 1
FROM matches m, players bw
WHERE m.public_id = $2 AND bw.public_id = $3 AND is_current_bowler=true AND inning_number= $4 AND m.id = b.match_id AND bw.id = b.bowler_id
RETURNING *;
`

func (q *Queries) UpdateBowlerStats(ctx context.Context, runs int32, matchPublicID, bowlerPublicID uuid.UUID, inningNumber int) (models.BowlerScore, error) {
	row := q.db.QueryRowContext(ctx, updateBowlingStats,
		runs,
		matchPublicID,
		bowlerPublicID,
		inningNumber,
	)
	var i models.BowlerScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.BowlerID,
		&i.InningNumber,
		&i.BallNumber,
		&i.Runs,
		&i.Wickets,
		&i.Wide,
		&i.NoBall,
		&i.BowlingStatus,
		&i.IsCurrentBowler,
	)
	return i, err
}

const getCurrentPlayingBatsman = `
	SELECT * FROM batsman_score bs
	JOIN matches m ON m.id = bs.match_id
	WHERE bs.match_id = $1 AND bs.batting_status = true AND inning_number= $2;
`

func (q *Queries) GetCurrentPlayingBatsmen(ctx context.Context, matchPublicID uuid.UUID, inningNumber int) ([]models.BatsmanScore, error) {
	rows, err := q.db.QueryContext(ctx, getCurrentPlayingBatsman, matchPublicID, inningNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batsmen []models.BatsmanScore
	for rows.Next() {
		var i models.BatsmanScore
		err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.MatchID,
			&i.TeamID,
			&i.BatsmanID,
			&i.InningNumber,
			&i.Position,
			&i.RunsScored,
			&i.BallsFaced,
			&i.Fours,
			&i.Sixes,
			&i.BattingStatus,
			&i.IsStriker,
			&i.IsCurrentlyBatting,
		)
		if err != nil {
			return nil, err
		}
		batsmen = append(batsmen, i)
	}

	// Check for any error that occurred during the iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return batsmen, nil
}

func (q *Queries) ToggleCricketStricker(ctx context.Context, matchPublicID uuid.UUID, inningNumber int) ([]models.BatsmanScore, error) {
	const query = `
		UPDATE batsman_score b
		SET is_striker = NOT b.is_striker
		FROM matches m
		WHERE b.match_id = m.id
		AND m.public_id = $1
		AND b.is_currently_batting = true
		AND b.inning_number = $2
		RETURNING b.*;
	`

	rows, err := q.db.QueryContext(ctx, query, matchPublicID, inningNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batsmen []models.BatsmanScore
	for rows.Next() {
		var bat models.BatsmanScore
		err := rows.Scan(
			&bat.ID,
			&bat.PublicID,
			&bat.MatchID,
			&bat.TeamID,
			&bat.BatsmanID,
			&bat.InningNumber,
			&bat.Position,
			&bat.RunsScored,
			&bat.BallsFaced,
			&bat.Fours,
			&bat.Sixes,
			&bat.BattingStatus,
			&bat.IsStriker,
			&bat.IsCurrentlyBatting,
		)
		if err != nil {
			return nil, err
		}

		batsmen = append(batsmen, bat)
	}

	// if len(batsmen) != 2 {
	// 	return nil, fmt.Errorf("unexpected number of batsmen updated: expected 2, got %d", len(batsmen))
	// }

	return batsmen, nil
}

const updateWideRun = `
	WITH get_match AS (
		SELECT * FROM matches WHERE public_id = $1
	),
	get_bowler AS (
		SELECT * FROM players WHERE public_id = $3
	),
	get_team AS (
		SELECT * FROM teams WHERE public_id = $2
	),
	update_bowler AS (
		UPDATE bowler_score
		SET 
			wide = wide + 1, 
			runs = runs + $4 + 1
		WHERE 
			match_id = (SELECT id FROM get_match) 
			AND bowler_id = (SELECT id FROM get_bowler) 
			AND is_current_bowler = true 
			AND inning_number = $5
		RETURNING *
	),
	update_inning_score AS (
		UPDATE cricket_score
		SET
			score = score + $4 + 1
		WHERE 
			match_id = (SELECT id FROM get_match)
			AND team_id = (SELECT id FROM get_team)
			AND inning_number = $5
		RETURNING *
	)
	SELECT 
		bs.*,
		ubl.*, 
		uis.*
	FROM batsman_score bs
	JOIN update_inning_score uis
			ON bs.match_id = uis.match_id 
			AND bs.inning_number = uis.inning_number
	JOIN update_bowler ubl 
		ON bs.match_id = ubl.match_id 
	AND bs.inning_number = ubl.inning_number
	WHERE bs.match_id = (SELECT id FROM get_match) AND bs.team_id = (SELECT id FROM get_team) AND bs.inning_number = $5 AND bs.is_striker = true`

func (q *Queries) UpdateWideRuns(ctx context.Context, matchPublicID, battingTeamPublicID, bowlerPublicID uuid.UUID, runsScored int32, inningNumber int) (*models.BatsmanScore, *models.BowlerScore, *models.CricketScore, error) {
	var bowler models.BowlerScore
	var batsman models.BatsmanScore
	var inningScore models.CricketScore
	row := q.db.QueryRowContext(ctx, updateWideRun, matchPublicID, battingTeamPublicID, bowlerPublicID, runsScored, inningNumber)
	err := row.Scan(
		&batsman.ID,
		&batsman.PublicID,
		&batsman.MatchID,
		&batsman.TeamID,
		&batsman.BatsmanID,
		&batsman.InningNumber,
		&batsman.Position,
		&batsman.RunsScored,
		&batsman.BallsFaced,
		&batsman.Fours,
		&batsman.Sixes,
		&batsman.BattingStatus,
		&batsman.IsStriker,
		&batsman.IsCurrentlyBatting,
		&bowler.ID,
		&bowler.PublicID,
		&bowler.MatchID,
		&bowler.TeamID,
		&bowler.BowlerID,
		&bowler.InningNumber,
		&bowler.BallNumber,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
		&inningScore.ID,
		&inningScore.PublicID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.InningNumber,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
		&inningScore.InningStatus,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to exec query: %w", err)
	}
	return &batsman, &bowler, &inningScore, nil
}

const updateNoBallRun = `
	WITH get_match AS (
		SELECT * FROM matches WHERE public_id = $1
	),
	get_bowler AS (
		SELECT * FROM players WHERE public_id = $3
	),
	get_team AS (
		SELECT * FROM teams WHERE public_id = $2
	),
	update_bowler AS (
		UPDATE bowler_score
		SET 
			no_ball = no_ball + 1, 
			runs = runs + $4 + 1
		WHERE 
			match_id = (SELECT id FROM get_match) 
			AND bowler_id = (SELECT id FROM get_bowler) 
			AND is_current_bowler = true 
			AND inning_number = $5
		RETURNING *
	),
	update_inning_score AS (
		UPDATE cricket_score
		SET
			score = score + $4 + 1
		WHERE 
			match_id = (SELECT id FROM get_match)
			AND team_id = (SELECT id FROM get_team)
			AND inning_number = $5
		RETURNING *
	)
	SELECT 
		bs.*,
		ubl.*, 
		uis.*
	FROM batsman_score bs
	JOIN update_inning_score uis
			ON bs.match_id = uis.match_id 
			AND bs.inning_number = uis.inning_number
	JOIN update_bowler ubl 
		ON bs.match_id = ubl.match_id 
	AND bs.inning_number = ubl.inning_number
	WHERE bs.match_id = (SELECT id FROM get_match) AND bs.team_id = (SELECT id FROM get_team) AND bs.inning_number = $5 AND is_striker = true
`

func (q *Queries) UpdateNoBallsRuns(ctx context.Context, matchPublicID, battingTeamPublicID, bowlerPublicID uuid.UUID, runsScored int32, inningNumber int) (*models.BatsmanScore, *models.BowlerScore, *models.CricketScore, error) {
	var bowler models.BowlerScore
	var batsman models.BatsmanScore
	var inningScore models.CricketScore

	rows, err := q.db.QueryContext(ctx, updateNoBallRun, matchPublicID, battingTeamPublicID, bowlerPublicID, runsScored, inningNumber)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to exec query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil, nil, fmt.Errorf("no rows returned - check if striker exists and updates succeeded")
	}

	err = rows.Scan(
		&batsman.ID,
		&batsman.PublicID,
		&batsman.MatchID,
		&batsman.TeamID,
		&batsman.BatsmanID,
		&batsman.InningNumber,
		&batsman.Position,
		&batsman.RunsScored,
		&batsman.BallsFaced,
		&batsman.Fours,
		&batsman.Sixes,
		&batsman.BattingStatus,
		&batsman.IsStriker,
		&batsman.IsCurrentlyBatting,
		&bowler.ID,
		&bowler.PublicID,
		&bowler.MatchID,
		&bowler.TeamID,
		&bowler.BowlerID,
		&bowler.InningNumber,
		&bowler.BallNumber,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
		&inningScore.ID,
		&inningScore.PublicID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.InningNumber,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
		&inningScore.InningStatus,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &batsman, &bowler, &inningScore, nil
}

// Enhance about the no ball_number
const addCricketWicket = `
WITH resolved_ids AS (
  SELECT 
    m.id AS match_id,
    t.id AS team_id,
    batsman.id AS batsman_id,
    bowler.id AS bowler_id,
    fielder.id AS fielder_id
  FROM matches m
  JOIN teams t ON t.public_id = $2
  JOIN players batsman ON batsman.public_id = $3
  JOIN players bowler ON bowler.public_id = $4
  LEFT JOIN players fielder ON fielder.public_id = $8  -- fielder may be null
  WHERE m.public_id = $1
),
add_wicket AS (
  INSERT INTO wickets (
    match_id,
    team_id,
    batsman_id,
    bowler_id,
    wickets_number,
    wicket_type,
    ball_number,
    fielder_id,
    score,
    inning_number
  )
  SELECT 
    r.match_id, r.team_id, r.batsman_id, r.bowler_id,
    $5, $6, $7, r.fielder_id, $9, $10
  FROM resolved_ids r
  RETURNING *
),
update_out_batsman AS (
  UPDATE batsman_score b
  SET 
    balls_faced = balls_faced + 1,
    runs_scored = runs_scored + (CASE WHEN is_striker THEN GREATEST($9, 0) ELSE 0 END),
    is_currently_batting = false,
    is_striker = false
  FROM resolved_ids r
  WHERE b.match_id = r.match_id
    AND b.team_id = r.team_id
    AND b.batsman_id = r.batsman_id
    AND b.inning_number = $10
  RETURNING *
),
update_not_out_batsman AS (
  UPDATE batsman_score b
  SET 
    balls_faced = balls_faced + 1,
    runs_scored = runs_scored + (CASE WHEN is_striker THEN GREATEST($9, 0) ELSE 0 END)
  FROM resolved_ids r
  WHERE b.match_id = r.match_id
    AND b.team_id = r.team_id
    AND b.batsman_id <> r.batsman_id
    AND b.is_currently_batting = true
    AND b.inning_number = $10
  RETURNING *
),
update_bowler AS (
  UPDATE bowler_score bl
  SET 
    wickets = CASE WHEN $6 != 'Run Out' THEN bl.wickets + 1 ELSE bl.wickets END,
    runs = bl.runs + GREATEST($9, 0),
    ball_number = bl.ball_number + 1
  FROM resolved_ids r
  WHERE bl.match_id = r.match_id
    AND bl.bowler_id = r.bowler_id
    AND bl.is_current_bowler = true
    AND bl.inning_number = $10
  RETURNING *
),
update_inning_score AS (
  UPDATE cricket_score cs
  SET 
    overs = cs.overs + 1,
    wickets = cs.wickets + 1,
    score = cs.score + GREATEST($9, 0)
  FROM resolved_ids r
  WHERE cs.match_id = r.match_id
    AND cs.team_id = r.team_id
    AND cs.inning_number = $10
  RETURNING *
)
SELECT 
  o.*, n.*, b.*, sc.*, w.*
FROM add_wicket w
JOIN update_out_batsman o ON w.match_id = o.match_id AND w.team_id = o.team_id AND o.inning_number = w.inning_number
JOIN update_not_out_batsman n ON w.match_id = n.match_id AND w.team_id = n.team_id AND n.inning_number = w.inning_number
JOIN update_bowler b ON w.match_id = b.match_id AND w.bowler_id = b.bowler_id AND b.inning_number = w.inning_number
JOIN update_inning_score sc ON w.match_id = sc.match_id AND w.team_id = sc.team_id AND sc.inning_number = w.inning_number;
`

func (q *Queries) AddCricketWicket(ctx context.Context, matchPublicID, teamPublicID, batsmanPublicID, bowlerPublicID uuid.UUID, wicketNumber int, wicketType string, ballNumber int, fielderID uuid.UUID, score int, runsScored int32, inningNumber int) (*models.BatsmanScore, *models.BatsmanScore, *models.BowlerScore, *models.CricketScore, *models.Wicket, error) {
	var outBatsman models.BatsmanScore
	var notOutBatsman models.BatsmanScore
	var bowler models.BowlerScore
	var inningScore models.CricketScore
	var wickets models.Wicket

	row := q.db.QueryRowContext(ctx, addCricketWicket, matchPublicID, teamPublicID, batsmanPublicID, bowlerPublicID, wicketNumber, wicketType, ballNumber, fielderID, score, runsScored, inningNumber)
	err := row.Scan(
		&outBatsman.ID,
		&outBatsman.BatsmanID,
		&outBatsman.TeamID,
		&outBatsman.MatchID,
		&outBatsman.Position,
		&outBatsman.RunsScored,
		&outBatsman.BallsFaced,
		&outBatsman.Fours,
		&outBatsman.Sixes,
		&outBatsman.BattingStatus,
		&outBatsman.IsStriker,
		&outBatsman.IsCurrentlyBatting,
		&outBatsman.InningNumber,
		&notOutBatsman.ID,
		&notOutBatsman.BatsmanID,
		&notOutBatsman.TeamID,
		&notOutBatsman.MatchID,
		&notOutBatsman.Position,
		&notOutBatsman.RunsScored,
		&notOutBatsman.BallsFaced,
		&notOutBatsman.Fours,
		&notOutBatsman.Sixes,
		&notOutBatsman.BattingStatus,
		&notOutBatsman.IsStriker,
		&notOutBatsman.IsCurrentlyBatting,
		&notOutBatsman.InningNumber,
		&bowler.ID,
		&bowler.TeamID,
		&bowler.MatchID,
		&bowler.BowlerID,
		&bowler.BallNumber,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
		&bowler.InningNumber,
		&inningScore.ID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.InningNumber,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
		&inningScore.InningStatus,
		&wickets.ID,
		&wickets.MatchID,
		&wickets.TeamID,
		&wickets.BatsmanID,
		&wickets.BowlerID,
		&wickets.WicketsNumber,
		&wickets.WicketType,
		&wickets.BallNumber,
		&wickets.FielderID,
		&wickets.Score,
		&wickets.InningNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil, nil, nil
		}
		return nil, nil, nil, nil, nil, fmt.Errorf("Failed to scan query: ", err)
	}
	return &outBatsman, &notOutBatsman, &bowler, &inningScore, &wickets, nil
}

const addCricketWicketWithBowlType = `
WITH add_wicket AS (
    INSERT INTO wickets (
        match_id,
        team_id,
        batsman_id,
        bowler_id,
		inning_number,
        wickets_number,
        wicket_type,
        ball_number,
        fielder_id,
        score
    ) 
    SELECT 
        m.id,
        t.id,
        bp.id,
        bowler_p.id,
        $5,
        $6,
        fp.id,
        $8
        $9,
        $10
    FROM matches m
    CROSS JOIN teams t
    CROSS JOIN players bp  -- batsman player
    CROSS JOIN players bowler_p  -- bowler player
    LEFT JOIN players fp ON fp.public_id = $8  -- fielder player (optional)
    WHERE m.public_id = $1
      AND t.public_id = $2
      AND bp.public_id = $3
      AND bowler_p.public_id = $4
    RETURNING *
),
update_out_batsman AS (
    UPDATE batsman_score
    SET balls_faced = balls_faced + 1,
        runs_scored = runs_scored + (CASE WHEN is_striker THEN (CASE WHEN $9 > 0 THEN $9 ELSE 0 END) ELSE 0 END),
        is_currently_batting = false,
        is_striker = false
    FROM matches m, teams t, players p
    WHERE batsman_score.match_id = m.id 
      AND batsman_score.batsman_id = p.id 
      AND batsman_score.team_id = t.id
      AND m.public_id = $1
      AND t.public_id = $2
      AND p.public_id = $3
      AND batsman_score.inning_number = $10
    RETURNING batsman_score.*
),
update_not_out_batsman AS (
    UPDATE batsman_score
    SET balls_faced = balls_faced + 1,
        runs_scored = runs_scored + (CASE WHEN is_striker THEN (CASE WHEN $9 > 0 THEN $9 ELSE 0 END) ELSE 0 END)
    FROM matches m, teams t, players batsman_p
    WHERE batsman_score.match_id = m.id 
      AND batsman_score.team_id = t.id 
      AND batsman_score.batsman_id <> (SELECT id FROM players WHERE public_id = $3)
      AND batsman_score.is_currently_batting = true
      AND m.public_id = $1
      AND t.public_id = $2
      AND batsman_score.inning_number = $10
    RETURNING batsman_score.*
),
update_bowler AS (
    UPDATE bowler_score
    SET wickets = CASE
                    WHEN $6 != 'Run Out' THEN wickets + 1
                    ELSE wickets
                  END,
        runs = runs + (CASE WHEN $9 > 0 THEN $9 ELSE 0 END),
        ball_number = ball_number,
        wide = wide + (CASE WHEN $11 = 'wide' THEN 1 ELSE 0 END),
        no_ball = no_ball + (CASE WHEN $11 = 'no_ball' THEN 1 ELSE 0 END)
    FROM matches m, players p
    WHERE bowler_score.match_id = m.id 
      AND bowler_score.bowler_id = p.id 
      AND bowler_score.is_current_bowler = true
      AND m.public_id = $1
      AND p.public_id = $4
      AND bowler_score.inning_number = $10
    RETURNING bowler_score.*
),
update_inning_score AS (
    UPDATE cricket_score
    SET overs = overs,
        wickets = wickets + 1,
        score = score + (CASE WHEN $9 > 0 THEN $9 ELSE 0 END)
    FROM matches m, teams t
    WHERE cricket_score.match_id = m.id 
      AND cricket_score.team_id = t.id
      AND m.public_id = $1
      AND t.public_id = $2
      AND cricket_score.inning_number = $10
    RETURNING cricket_score.*
)
SELECT 
    o.*,
    n.*,
    b.*,
    sc.*,
    w.*
FROM add_wicket w
LEFT JOIN update_out_batsman o ON w.match_id = o.match_id
LEFT JOIN update_not_out_batsman n ON w.match_id = n.match_id
LEFT JOIN update_bowler b ON w.match_id = b.match_id
LEFT JOIN update_inning_score sc ON w.match_id = sc.match_id;
`

func (q *Queries) AddCricketWicketWithBowlType(ctx context.Context, matchPublicID, teamPublicID, batsmanPublicID, bowlerPublicID uuid.UUID, wicketNumber int, wicketType string, ballNumber int, fielderPublicID *uuid.UUID, score int, bowlType string, inningNumber int) (*models.BatsmanScore, *models.BatsmanScore, *models.BowlerScore, *models.CricketScore, *models.Wicket, error) {
	var outBatsman models.BatsmanScore
	var notOutBatsman models.BatsmanScore
	var bowler models.BowlerScore
	var inningScore models.CricketScore
	var wickets models.Wicket

	// Handle optional fielder parameter
	var fielderParam interface{}
	if fielderPublicID != nil {
		fielderParam = *fielderPublicID
	} else {
		fielderParam = nil
	}

	row := q.db.QueryRowContext(ctx, addCricketWicketWithBowlType,
		matchPublicID,
		teamPublicID,
		batsmanPublicID,
		bowlerPublicID,
		inningNumber,
		wicketNumber,
		wicketType,
		ballNumber,
		fielderParam,
		score,
		bowlType,
	)

	err := row.Scan(
		&outBatsman.ID,
		&outBatsman.PublicID,
		&outBatsman.MatchID,
		&outBatsman.TeamID,
		&outBatsman.BatsmanID,
		&outBatsman.Position,
		&outBatsman.RunsScored,
		&outBatsman.BallsFaced,
		&outBatsman.Fours,
		&outBatsman.Sixes,
		&outBatsman.BattingStatus,
		&outBatsman.IsStriker,
		&outBatsman.IsCurrentlyBatting,
		&outBatsman.InningNumber,
		&notOutBatsman.ID,
		&notOutBatsman.PublicID,
		&notOutBatsman.MatchID,
		&notOutBatsman.TeamID,
		&notOutBatsman.BatsmanID,
		&notOutBatsman.Position,
		&notOutBatsman.RunsScored,
		&notOutBatsman.BallsFaced,
		&notOutBatsman.Fours,
		&notOutBatsman.Sixes,
		&notOutBatsman.BattingStatus,
		&notOutBatsman.IsStriker,
		&notOutBatsman.IsCurrentlyBatting,
		&notOutBatsman.InningNumber,
		&bowler.ID,
		&bowler.PublicID,
		&bowler.MatchID,
		&bowler.TeamID,
		&bowler.BowlerID,
		&bowler.BallNumber,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
		&bowler.InningNumber,
		&inningScore.ID,
		&inningScore.PublicID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.InningNumber,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
		&inningScore.InningStatus,
		&wickets.ID,
		&wickets.PublicID,
		&wickets.MatchID,
		&wickets.TeamID,
		&wickets.BatsmanID,
		&wickets.BowlerID,
		&wickets.WicketsNumber,
		&wickets.WicketType,
		&wickets.BallNumber,
		&wickets.FielderID,
		&wickets.Score,
		&wickets.InningNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil, nil, nil
		}
		return nil, nil, nil, nil, nil, fmt.Errorf("Failed to scan query: %w", err)
	}

	return &outBatsman, &notOutBatsman, &bowler, &inningScore, &wickets, nil
}

const updateInningEndStatus = `
WITH update_inning_number AS (
	UPDATE cricket_score
	SET is_inning_completed = true
		inning_status = 'completed'
	WHERE match_id = $1 AND team_id = $2 AND inning_number= $3
	RETURNING *
),
update_batsman AS (
	UPDATE batsman_score
	SET is_striker = false,
		is_currently_batting = false
	WHERE match_id = $1 AND team_id = $2 AND is_striker = true AND inning_number= $3
	RETURNING *
),
update_bowler AS (
	UPDATE bowler_score
	SET is_current_bowler = false
	WHERE match_id = $1 AND is_current_bowler = true AND inning_number = $3
	RETURNING *
)
SELECT
	ui.*,
	ub.*,
	ubl.*
FROM update_batsman ub
LEFT JOIN update_bowler AS ubl ON ub.match_id = ubl.match_id AND ub.inning_number= ubl.inning_number
LEFT JOIN update_inning_number AS ui ON ub.match_id = ui.match_id AND ui.team_id = ub.team_id AND ui.inning_number= ub.inning_number
`

func (q *Queries) UpdateInningEndStatus(ctx context.Context, matchID, batsmanTeamID int32, inningNumber int) (*models.CricketScore, *models.BatsmanScore, *models.BowlerScore, error) {
	var inningScore models.CricketScore
	var batsmanScore models.BatsmanScore
	var bowler models.BowlerScore

	row := q.db.QueryRowContext(ctx, updateInningEndStatus, matchID, batsmanTeamID, inningNumber)

	err := row.Scan(
		&inningScore.ID,
		&inningScore.PublicID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.InningNumber,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
		&inningScore.InningStatus,
		&batsmanScore.ID,
		&batsmanScore.PublicID,
		&batsmanScore.MatchID,
		&batsmanScore.TeamID,
		&batsmanScore.BatsmanID,
		&batsmanScore.InningNumber,
		&batsmanScore.Position,
		&batsmanScore.RunsScored,
		&batsmanScore.BallsFaced,
		&batsmanScore.Fours,
		&batsmanScore.Sixes,
		&batsmanScore.BattingStatus,
		&batsmanScore.IsStriker,
		&batsmanScore.IsCurrentlyBatting,
		&bowler.ID,
		&bowler.PublicID,
		&bowler.MatchID,
		&bowler.TeamID,
		&bowler.BowlerID,
		&bowler.InningNumber,
		&bowler.BallNumber,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil
		}
		return nil, nil, nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return &inningScore, &batsmanScore, &bowler, nil
}

const updateInningEndStatusByPublicIDQuery = `
WITH resolve_ids AS (
    SELECT m.id AS match_id, t.id AS team_id
    FROM matches m, teams t
    WHERE m.public_id = $1
      AND t.public_id = $2
),
update_inning_number AS (
    UPDATE cricket_score cs
    SET is_inning_completed = true,
		inning_status = 'completed'
    FROM resolve_ids ri
    WHERE cs.match_id = ri.match_id 
      AND cs.team_id = ri.team_id 
      AND cs.inning_number = $3
    RETURNING 
        cs.id AS ui_id,
        cs.public_id AS ui_public_id,
        cs.match_id AS ui_match_id,
        cs.team_id AS ui_team_id,
        cs.inning_number AS ui_inning_number,
        cs.score AS ui_score,
        cs.wickets AS ui_wickets,
        cs.overs AS ui_overs,
        cs.run_rate AS ui_run_rate,
        cs.target_run_rate AS ui_target_run_rate,
        cs.follow_on AS ui_follow_on,
        cs.is_inning_completed AS ui_is_inning_completed,
        cs.declared AS ui_declared,
		cs.inning_status AS ui_inning_status
),
update_batsman AS (
    UPDATE batsman_score bs
    SET is_striker = false,
        is_currently_batting = false
    FROM resolve_ids ri
    WHERE bs.match_id = ri.match_id 
      AND bs.team_id = ri.team_id 
      AND bs.inning_number = $3
    RETURNING 
        bs.id AS ub_id,
        bs.public_id AS ub_public_id,
        bs.match_id AS ub_match_id,
        bs.team_id AS ub_team_id,
        bs.batsman_id AS ub_batsman_id,
        bs.position AS ub_position,
        bs.runs_scored AS ub_runs_scored,
        bs.balls_faced AS ub_balls_faced,
        bs.fours AS ub_fours,
        bs.sixes AS ub_sixes,
        bs.batting_status AS ub_batting_status,
        bs.is_striker AS ub_is_striker,
        bs.is_currently_batting AS ub_is_currently_batting,
        bs.inning_number AS ub_inning_number
),
update_bowler AS (
    UPDATE bowler_score b
    SET is_current_bowler = false
    FROM resolve_ids ri
    WHERE b.match_id = ri.match_id 
      AND b.is_current_bowler = true 
      AND b.inning_number = $3
    RETURNING 
        b.id AS ubl_id,
        b.public_id AS ubl_public_id,
        b.match_id AS ubl_match_id,
        b.team_id AS ubl_team_id,
        b.bowler_id AS ubl_bowler_id,
        b.ball_number AS ubl_ball_number,
        b.runs AS ubl_runs,
        b.wickets AS ubl_wickets,
        b.wide AS ubl_wide,
        b.no_ball AS ubl_no_ball,
        b.bowling_status AS ubl_bowling_status,
        b.is_current_bowler AS ubl_is_current_bowler,
        b.inning_number AS ubl_inning_number
)
SELECT
    ui.*,
    ub.*,
    ubl.*
FROM update_batsman ub
LEFT JOIN update_bowler ubl 
       ON ub.ub_match_id = ubl.ubl_match_id 
      AND ub.ub_inning_number = ubl.ubl_inning_number
LEFT JOIN update_inning_number ui 
       ON ub.ub_match_id = ui.ui_match_id 
      AND ub.ub_team_id = ui.ui_team_id 
      AND ub.ub_inning_number = ui.ui_inning_number;
`

func (q *Queries) UpdateInningEndStatusByPublicID(ctx context.Context, matchPublicID, batsmanTeamPublicID uuid.UUID, inningNumber int) (*models.CricketScore, *models.BatsmanScore, *models.BowlerScore, error) {
	var inningScore models.CricketScore
	var batsmanScore models.BatsmanScore
	var bowler models.BowlerScore
	fmt.Println("Batsman ID: ", batsmanTeamPublicID)
	fmt.Println("Match ID: ", matchPublicID)
	fmt.Println("Inning :NUmber; ", inningNumber)

	fmt.Println("Line no data: 1484: ")

	row := q.db.QueryRowContext(ctx, updateInningEndStatusByPublicIDQuery, matchPublicID, batsmanTeamPublicID, inningNumber)
	fmt.Println("Line no 1487: ", row)
	err := row.Scan(
		&inningScore.ID,
		&inningScore.PublicID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.InningNumber,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
		&inningScore.InningStatus,
		&batsmanScore.ID,
		&batsmanScore.PublicID,
		&batsmanScore.MatchID,
		&batsmanScore.TeamID,
		&batsmanScore.BatsmanID,
		&batsmanScore.Position,
		&batsmanScore.RunsScored,
		&batsmanScore.BallsFaced,
		&batsmanScore.Fours,
		&batsmanScore.Sixes,
		&batsmanScore.BattingStatus,
		&batsmanScore.IsStriker,
		&batsmanScore.IsCurrentlyBatting,
		&batsmanScore.InningNumber,
		&bowler.ID,
		&bowler.PublicID,
		&bowler.MatchID,
		&bowler.TeamID,
		&bowler.BowlerID,
		&bowler.BallNumber,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
		&bowler.InningNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No Rows in sql: ", err)
			return nil, nil, nil, err
		}
		return nil, nil, nil, fmt.Errorf("failed to exec query: %w", err)
	}

	fmt.Println("Rowls: ", err)

	return &inningScore, &batsmanScore, &bowler, nil
}

const updateInningScore = `
	WITH get_match AS (
    SELECT id AS match_id, home_team_id, away_team_id
    FROM matches
    WHERE public_id = $1
),
get_batting_team AS (
    SELECT id AS team_id
    FROM teams
    WHERE public_id = $2
),
get_batsman AS (
    SELECT id AS batsman_id
    FROM players
    WHERE public_id = $3
),
get_bowler AS (
    SELECT id AS bowler_id
    FROM players
    WHERE public_id = $4
),
get_bowling_team AS (
    SELECT CASE
        WHEN gm.home_team_id = (SELECT team_id FROM get_batting_team) THEN gm.away_team_id
        ELSE gm.home_team_id
    END AS team_id
    FROM get_match gm
),
update_batsman AS (
    UPDATE batsman_score
    SET 
        runs_scored = runs_scored + $5,
        balls_faced = balls_faced + 1,
        fours = fours + CASE WHEN $5 = 4 THEN 1 ELSE 0 END,
        sixes = sixes + CASE WHEN $5 = 6 THEN 1 ELSE 0 END
    WHERE match_id = (SELECT match_id FROM get_match)
      AND team_id = (SELECT team_id FROM get_batting_team)
      AND batsman_id = (SELECT batsman_id FROM get_batsman)
      AND inning_number = $6
    RETURNING *
),
update_bowler AS (
    UPDATE bowler_score
    SET 
        runs = runs + $5,
        ball_number = ball_number + 1
    WHERE match_id = (SELECT match_id FROM get_match)
      AND team_id = (SELECT team_id FROM get_bowling_team)
      AND bowler_id = (SELECT bowler_id FROM get_bowler)
      AND inning_number = $6
      AND $4 IS NOT NULL
    RETURNING *
),
update_inning_score AS (
    UPDATE cricket_score
    SET 
        score = score + $5,
        overs = overs + 1
    WHERE match_id = (SELECT match_id FROM get_match)
      AND team_id = (SELECT team_id FROM get_batting_team)
      AND inning_number = $6
    RETURNING *
)
SELECT 
    ub.*,
    ubl.*,
    uis.*
FROM update_batsman ub
LEFT JOIN update_bowler ubl ON TRUE
JOIN update_inning_score uis ON TRUE
`

func (q *Queries) UpdateInningScore(ctx context.Context, matchPublicID, batsmanTeamPublicID, batsmanPublicID, bowlerTeamID uuid.UUID, runsScored int32, inningNumber int) (*models.BatsmanScore, *models.BowlerScore, *models.CricketScore, error) {
	fmt.Println("Match Public Id: ", matchPublicID)
	fmt.Println("BatsmTeam ID: ", batsmanTeamPublicID)
	fmt.Println("Batsman ID: ", batsmanPublicID)
	fmt.Println("Bowler ID; ", bowlerTeamID)
	fmt.Println("Inning NUmber: ", inningNumber)
	var batsman models.BatsmanScore
	var bowler models.BowlerScore
	var inningScore models.CricketScore
	row := q.db.QueryRowContext(ctx, updateInningScore, matchPublicID, batsmanTeamPublicID, batsmanPublicID, bowlerTeamID, runsScored, inningNumber)

	err := row.Scan(
		&batsman.ID,
		&batsman.PublicID,
		&batsman.MatchID,
		&batsman.TeamID,
		&batsman.BatsmanID,
		&batsman.InningNumber,
		&batsman.Position,
		&batsman.RunsScored,
		&batsman.BallsFaced,
		&batsman.Fours,
		&batsman.Sixes,
		&batsman.BattingStatus,
		&batsman.IsStriker,
		&batsman.IsCurrentlyBatting,
		&bowler.ID,
		&bowler.PublicID,
		&bowler.MatchID,
		&bowler.TeamID,
		&bowler.BowlerID,
		&bowler.InningNumber,
		&bowler.BallNumber,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
		&inningScore.ID,
		&inningScore.PublicID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.InningNumber,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
		&inningScore.InningStatus,
	)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return &batsman, &bowler, &inningScore, nil
}

const updateSetBowlerStatus = `
	UPDATE bowler_score b
	SET is_current_bowler = NOT is_current_bowler
	FROM balls_status bs
	JOIN matches m ON m.id = b.match_id
	JOIN teams t ON t.id = b.team_id
	JOIN players p ON p.id = b.player_id
	WHERE m.public_id = $1 AND t.public_id = $2 AND p.public_id = $3 AND inning_number= $4
	RETURNING *
`

func (q *Queries) UpdateBowlingBowlerStatus(ctx context.Context, matchPublicID, teamPublicID, bowlerPublicID uuid.UUID, inningNumber int) (*models.BowlerScore, error) {
	var currentBowler models.BowlerScore

	row := q.db.QueryRowContext(ctx, updateSetBowlerStatus, matchPublicID, teamPublicID, bowlerPublicID, inningNumber)

	err := row.Scan(
		&currentBowler.ID,
		&currentBowler.PublicID,
		&currentBowler.MatchID,
		&currentBowler.TeamID,
		&currentBowler.BowlerID,
		&currentBowler.BallNumber,
		&currentBowler.Runs,
		&currentBowler.Wickets,
		&currentBowler.Wide,
		&currentBowler.NoBall,
		&currentBowler.BowlingStatus,
		&currentBowler.IsCurrentBowler,
		&currentBowler.InningNumber,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return &currentBowler, nil
}

const getCurrentBattingBatsmanQuery = `
	SELECT b.* FROM batsman_score b
	LEFT JOIN matches AS m ON m.id = b.match_id
	LEFT JOIN teams AS t ON t.id = b.team_id
	WHERE m.public_id=$1 AND t.public_id=$2 AND is_currently_batting=true AND inning_number= $3;
`

func (q *Queries) GetCurrentBattingBatsman(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, inningNumber int) ([]models.BatsmanScore, error) {
	rows, err := q.db.QueryContext(ctx, getCurrentBattingBatsmanQuery, matchPublicID, teamPublicID, inningNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batsmen []models.BatsmanScore
	for rows.Next() {
		var bat models.BatsmanScore
		err := rows.Scan(
			&bat.ID,
			&bat.PublicID,
			&bat.MatchID,
			&bat.TeamID,
			&bat.BatsmanID,
			&bat.InningNumber,
			&bat.Position,
			&bat.RunsScored,
			&bat.BallsFaced,
			&bat.Fours,
			&bat.Sixes,
			&bat.BattingStatus,
			&bat.IsStriker,
			&bat.IsCurrentlyBatting,
		)
		if err != nil {
			return nil, err
		}

		batsmen = append(batsmen, bat)
	}
	if len(batsmen) != 2 {
		return nil, fmt.Errorf("unexpected number of batsmen updated: expected 2, got %d", len(batsmen))
	}

	return batsmen, nil
}

const getCurrentBatsmanQuery = `
	SELECT 
    	JSON_BUILD_OBJECT(
			'team', JSON_BUILD_OBJECT(
				'id', tm.id, 
				'public_id', tm.public_id,
				'name', tm.name, 
				'slug', tm.slug, 
				'short_name', tm.shortname,
				'media_url', tm.media_url, 
				'gender', tm.gender, 
				'national', tm.national, 
				'country', tm.country, 
				'type', tm.type, 
				'player_count', tm.player_count, 
				'game_id', tm.game_id
			),
        	'batsman', JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', bt.id, 
					'public_id', bt.public_id,
					'batsman_id', bt.batsman_id,
					'player', JSON_BUILD_OBJECT('id',pl.id,'public_id',pl.public_id, 'name', pl.name, 'slug', pl.slug, 'short_name',pl.short_name, 'country', pl.country, 'positions', pl.positions, 'media_url', pl.media_url),
					'position', bt.position, 
					'runs_scored', bt.runs_scored, 
					'balls_faced', bt.balls_faced, 
					'fours', bt.fours, 
					'sixes', bt.sixes, 
					'batting_status', bt.batting_status, 
					'is_striker', bt.is_striker, 
					'is_currently_batting', bt.is_currently_batting,
					'inning_number', bt.inning_number
				)
        	)
    	)
	FROM batsman_score bt
	JOIN players AS pl ON pl.id = bt.batsman_id
	JOIN teams AS tm ON tm.id = bt.team_id
	JOIN matches AS m ON m.id = bt.match_id
	WHERE m.public_id = $1 AND tm.public_id = $2 AND bt.inning_number= $3 AND bt.is_currently_batting = true
	GROUP BY tm.id;
`

func (q *Queries) GetCurrentBatsman(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, inningNumber int) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getCurrentBatsmanQuery, matchPublicID, teamPublicID, inningNumber)

	var jsonBytes []byte
	if err := row.Scan(&jsonBytes); err != nil {
		return nil, fmt.Errorf("failed to scan json data: %w", err)
	}

	var currentBatsman interface{}

	err := json.Unmarshal(jsonBytes, &currentBatsman)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal bowler: ", err)
	}

	return currentBatsman, nil
}

const getCurrentBowlerQuery = `
	SELECT 
    	JSON_BUILD_OBJECT(
			'team', JSON_BUILD_OBJECT(
				'id', tm.id, 
				'public_id', tm.public_id,
				'name', tm.name, 
				'slug', tm.slug, 
				'short_name', tm.shortname,
				'media_url', tm.media_url, 
				'gender', tm.gender, 
				'national', tm.national, 
				'country', tm.country, 
				'type', tm.type, 
				'player_count', tm.player_count, 
				'game_id', tm.game_id
			),
        	'bowler', JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', bl.id, 
					'public_id', bl.public_id,
					'match_id', bl.match_id,
					'team_id', bl.team_id,
					'bowler_id', bl.bowler_id,
					'player', JSON_BUILD_OBJECT('id',pl.id,'public_id',pl.public_id, 'name', pl.name, 'slug', pl.slug, 'short_name',pl.short_name, 'country', pl.country, 'positions', pl.positions, 'media_url', pl.media_url),
					'runs', bl.runs, 
					'ball_number', bl.ball_number, 
					'wickets', bl.wickets, 
					'wide', bl.wide, 
					'no_ball', bl.no_ball,
					'bowling_status', bl.bowling_status,
					'is_current_bowler', bl.is_current_bowler,
					'inning_number', bl.inning_number
				)
			)
    	)
	FROM bowler_score bl
	JOIN matches AS m ON m.id = bl.match_id
	JOIN players AS pl ON pl.id = bl.bowler_id
	JOIN teams AS tm ON tm.id = bl.team_id
	WHERE m.public_id = $1 AND tm.public_id = $2 AND bl.inning_number= $3 AND bl.is_current_bowler = true
	GROUP BY tm.id;
`

func (q *Queries) GetCurrentBowler(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, inningNumber int) (interface{}, error) {
	var jsonBytes []byte
	row := q.db.QueryRowContext(ctx, getCurrentBowlerQuery, matchPublicID, teamPublicID, inningNumber)
	if err := row.Scan(&jsonBytes); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan json: %w", err)
	}

	var currentBowler interface{}

	err := json.Unmarshal(jsonBytes, &currentBowler)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal bowler: ", err)
	}

	return currentBowler, nil
}
