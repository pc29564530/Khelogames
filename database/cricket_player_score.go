package database

import (
	"context"
	"khelogames/database/models"
)

const addCricketBall = `
INSERT INTO balls (
    match_id,
    team_id,
    bowler_id,
    ball,
    runs,
    wickets,
    wide,
    no_ball
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, team_id, match_id, bowler_id, ball, runs, wickets, wide, no_ball
`

type AddCricketBallParams struct {
	MatchID  int64 `json:"match_id"`
	TeamID   int64 `json:"team_id"`
	BowlerID int64 `json:"bowler_id"`
	Ball     int32 `json:"ball"`
	Runs     int32 `json:"runs"`
	Wickets  int32 `json:"wickets"`
	Wide     int32 `json:"wide"`
	NoBall   int32 `json:"no_ball"`
}

func (q *Queries) AddCricketBall(ctx context.Context, arg AddCricketBallParams) (models.Ball, error) {
	row := q.db.QueryRowContext(ctx, addCricketBall,
		arg.MatchID,
		arg.TeamID,
		arg.BowlerID,
		arg.Ball,
		arg.Runs,
		arg.Wickets,
		arg.Wide,
		arg.NoBall,
	)
	var i models.Ball
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.MatchID,
		&i.BowlerID,
		&i.Ball,
		&i.Runs,
		&i.Wickets,
		&i.Wide,
		&i.NoBall,
	)
	return i, err
}

const addCricketBatsScore = `
INSERT INTO bats (
    batsman_id,
    match_id,
    team_id,
    position,
    runs_scored,
    balls_faced,
    fours,
    sixes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, batsman_id, team_id, match_id, position, runs_scored, balls_faced, fours, sixes
`

type AddCricketBatsScoreParams struct {
	BatsmanID  int64 `json:"batsman_id"`
	MatchID    int64 `json:"match_id"`
	TeamID     int64 `json:"team_id"`
	Position   int32 `json:"position"`
	RunsScored int32 `json:"runs_scored"`
	BallsFaced int32 `json:"balls_faced"`
	Fours      int32 `json:"fours"`
	Sixes      int32 `json:"sixes"`
}

func (q *Queries) AddCricketBatsScore(ctx context.Context, arg AddCricketBatsScoreParams) (models.Bat, error) {
	row := q.db.QueryRowContext(ctx, addCricketBatsScore,
		arg.BatsmanID,
		arg.MatchID,
		arg.TeamID,
		arg.Position,
		arg.RunsScored,
		arg.BallsFaced,
		arg.Fours,
		arg.Sixes,
	)
	var i models.Bat
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
	)
	return i, err
}

const addCricketWickets = `
INSERT INTO wickets (
    match_id,
    team_id,
    batsman_id,
    bowler_id,
    wickets_number,
    wicket_type,
    ball_number
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, match_id, team_id, batsman_id, bowler_id, wickets_number, wicket_type, ball_number
`

type AddCricketWicketsParams struct {
	MatchID       int64  `json:"match_id"`
	TeamID        int64  `json:"team_id"`
	BatsmanID     int64  `json:"batsman_id"`
	BowlerID      int64  `json:"bowler_id"`
	WicketsNumber int32  `json:"wickets_number"`
	WicketType    string `json:"wicket_type"`
	BallNumber    int32  `json:"ball_number"`
}

func (q *Queries) AddCricketWickets(ctx context.Context, arg AddCricketWicketsParams) (models.Wicket, error) {
	row := q.db.QueryRowContext(ctx, addCricketWickets,
		arg.MatchID,
		arg.TeamID,
		arg.BatsmanID,
		arg.BowlerID,
		arg.WicketsNumber,
		arg.WicketType,
		arg.BallNumber,
	)
	var i models.Wicket
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.BatsmanID,
		&i.BowlerID,
		&i.WicketsNumber,
		&i.WicketType,
		&i.BallNumber,
	)
	return i, err
}

const getCricketBall = `
SELECT id, team_id, match_id, bowler_id, ball, runs, wickets, wide, no_ball FROM balls
WHERE match_id=$1 AND bowler_id=$2 LIMIT 1
`

type GetCricketBallParams struct {
	MatchID  int64 `json:"match_id"`
	BowlerID int64 `json:"bowler_id"`
}

func (q *Queries) GetCricketBall(ctx context.Context, arg GetCricketBallParams) (models.Ball, error) {
	row := q.db.QueryRowContext(ctx, getCricketBall, arg.MatchID, arg.BowlerID)
	var i models.Ball
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.MatchID,
		&i.BowlerID,
		&i.Ball,
		&i.Runs,
		&i.Wickets,
		&i.Wide,
		&i.NoBall,
	)
	return i, err
}

const getCricketBalls = `
SELECT id, team_id, match_id, bowler_id, ball, runs, wickets, wide, no_ball FROM balls
WHERE match_id=$1 AND team_id=$2
`

type GetCricketBallsParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) GetCricketBalls(ctx context.Context, arg GetCricketBallsParams) ([]models.Ball, error) {
	rows, err := q.db.QueryContext(ctx, getCricketBalls, arg.MatchID, arg.TeamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Ball
	for rows.Next() {
		var i models.Ball
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.MatchID,
			&i.BowlerID,
			&i.Ball,
			&i.Runs,
			&i.Wickets,
			&i.Wide,
			&i.NoBall,
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
SELECT id, batsman_id, team_id, match_id, position, runs_scored, balls_faced, fours, sixes FROM bats
WHERE match_id=$1 AND batsman_id=$2 LIMIT 1
`

type GetCricketPlayerScoreParams struct {
	MatchID   int64 `json:"match_id"`
	BatsmanID int64 `json:"batsman_id"`
}

func (q *Queries) GetCricketPlayerScore(ctx context.Context, arg GetCricketPlayerScoreParams) (models.Bat, error) {
	row := q.db.QueryRowContext(ctx, getCricketPlayerScore, arg.MatchID, arg.BatsmanID)
	var i models.Bat
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
	)
	return i, err
}

const getCricketPlayersScore = `
SELECT id, batsman_id, team_id, match_id, position, runs_scored, balls_faced, fours, sixes FROM bats
WHERE match_id=$1 AND team_id=$2
ORDER BY position
`

type GetCricketPlayersScoreParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) GetCricketPlayersScore(ctx context.Context, arg GetCricketPlayersScoreParams) ([]models.Bat, error) {
	rows, err := q.db.QueryContext(ctx, getCricketPlayersScore, arg.MatchID, arg.TeamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Bat
	for rows.Next() {
		var i models.Bat
		if err := rows.Scan(
			&i.ID,
			&i.BatsmanID,
			&i.TeamID,
			&i.MatchID,
			&i.Position,
			&i.RunsScored,
			&i.BallsFaced,
			&i.Fours,
			&i.Sixes,
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

const getCricketWicket = `
SELECT id, match_id, team_id, batsman_id, bowler_id, wickets_number, wicket_type, ball_number FROM wickets
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
		&i.MatchID,
		&i.TeamID,
		&i.BatsmanID,
		&i.BowlerID,
		&i.WicketsNumber,
		&i.WicketType,
		&i.BallNumber,
	)
	return i, err
}

const getCricketWickets = `-- name: GetCricketWickets :many
SELECT id, match_id, team_id, batsman_id, bowler_id, wickets_number, wicket_type, ball_number FROM wickets
WHERE match_id=$1 AND team_id=$2
`

type GetCricketWicketsParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) GetCricketWickets(ctx context.Context, arg GetCricketWicketsParams) ([]models.Wicket, error) {
	rows, err := q.db.QueryContext(ctx, getCricketWickets, arg.MatchID, arg.TeamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Wicket
	for rows.Next() {
		var i models.Wicket
		if err := rows.Scan(
			&i.ID,
			&i.MatchID,
			&i.TeamID,
			&i.BatsmanID,
			&i.BowlerID,
			&i.WicketsNumber,
			&i.WicketType,
			&i.BallNumber,
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
UPDATE balls
SET 
    ball = $1,
    runs = $2,
    wickets = $3,
    wide = $4,
    no_ball = $5
WHERE match_id = $6 AND bowler_id = $7 AND team_id=$8
RETURNING id, team_id, match_id, bowler_id, ball, runs, wickets, wide, no_ball
`

type UpdateCricketBowlerParams struct {
	Ball     int32 `json:"ball"`
	Runs     int32 `json:"runs"`
	Wickets  int32 `json:"wickets"`
	Wide     int32 `json:"wide"`
	NoBall   int32 `json:"no_ball"`
	MatchID  int64 `json:"match_id"`
	BowlerID int64 `json:"bowler_id"`
	TeamID   int64 `json:"team_id"`
}

func (q *Queries) UpdateCricketBowler(ctx context.Context, arg UpdateCricketBowlerParams) (models.Ball, error) {
	row := q.db.QueryRowContext(ctx, updateCricketBowler,
		arg.Ball,
		arg.Runs,
		arg.Wickets,
		arg.Wide,
		arg.NoBall,
		arg.MatchID,
		arg.BowlerID,
		arg.TeamID,
	)
	var i models.Ball
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.MatchID,
		&i.BowlerID,
		&i.Ball,
		&i.Runs,
		&i.Wickets,
		&i.Wide,
		&i.NoBall,
	)
	return i, err
}

const updateCricketRunsScored = `
UPDATE bats
SET runs_scored = $1,
    balls_faced = $2,
    fours = $3,
    sixes = $4
WHERE match_id = $5 AND batsman_id = $6 AND team_id=$7
RETURNING id, batsman_id, team_id, match_id, position, runs_scored, balls_faced, fours, sixes
`

type UpdateCricketRunsScoredParams struct {
	RunsScored int32 `json:"runs_scored"`
	BallsFaced int32 `json:"balls_faced"`
	Fours      int32 `json:"fours"`
	Sixes      int32 `json:"sixes"`
	MatchID    int64 `json:"match_id"`
	BatsmanID  int64 `json:"batsman_id"`
	TeamID     int64 `json:"team_id"`
}

func (q *Queries) UpdateCricketRunsScored(ctx context.Context, arg UpdateCricketRunsScoredParams) (models.Bat, error) {
	row := q.db.QueryRowContext(ctx, updateCricketRunsScored,
		arg.RunsScored,
		arg.BallsFaced,
		arg.Fours,
		arg.Sixes,
		arg.MatchID,
		arg.BatsmanID,
		arg.TeamID,
	)
	var i models.Bat
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
	)
	return i, err
}
