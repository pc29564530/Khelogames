package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"

	"github.com/gin-gonic/gin"
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
    no_ball,
	bowling_status,
	is_current_bowler
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *
`

type AddCricketBallParams struct {
	MatchID         int64 `json:"match_id"`
	TeamID          int64 `json:"team_id"`
	BowlerID        int64 `json:"bowler_id"`
	Ball            int32 `json:"ball"`
	Runs            int32 `json:"runs"`
	Wickets         int32 `json:"wickets"`
	Wide            int32 `json:"wide"`
	NoBall          int32 `json:"no_ball"`
	BowlingStatus   bool  `json:"bowling_status"`
	IsCurrentBowler bool  `json:"is_current_bowler"`
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
		arg.BowlingStatus,
		arg.IsCurrentBowler,
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
		&i.BowlingStatus,
		&i.IsCurrentBowler,
	)
	return i, err
}

const getCricketStricker = `
	SELECT * FROM bats
	WHERE match_id=$1 AND team_id=$2 AND is_currently_batting=true AND is_striker=true;
`

func (q *Queries) GetCricketStricker(ctx context.Context, matchID, teamID int64) (*models.Bat, error) {
	row := q.db.QueryRowContext(ctx, getCricketStricker, matchID, teamID)

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
INSERT INTO bats (
    batsman_id,
    match_id,
    team_id,
    position,
    runs_scored,
    balls_faced,
    fours,
    sixes,
	batting_status,
	is_striker,
	is_currently_batting
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, batsman_id, team_id, match_id, position, runs_scored, balls_faced, fours, sixes, batting_status, is_striker, is_currently_batting
`

type AddCricketBatsScoreParams struct {
	BatsmanID          int64  `json:"batsman_id"`
	MatchID            int64  `json:"match_id"`
	TeamID             int64  `json:"team_id"`
	Position           string `json:"position"`
	RunsScored         int32  `json:"runs_scored"`
	BallsFaced         int32  `json:"balls_faced"`
	Fours              int32  `json:"fours"`
	Sixes              int32  `json:"sixes"`
	BattingStatus      bool   `json:"batting_status"`
	IsStriker          bool   `json:"is_striker"`
	IsCurrentlyBatting bool   `json:"is_currently_batting"`
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
		arg.BattingStatus,
		arg.IsStriker,
		arg.IsCurrentlyBatting,
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
		&i.BattingStatus,
		&i.IsStriker,
		&i.IsCurrentlyBatting,
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
    ball_number,
	fielder_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *
`

type AddCricketWicketsParams struct {
	MatchID       int64  `json:"match_id"`
	TeamID        int64  `json:"team_id"`
	BatsmanID     int64  `json:"batsman_id"`
	BowlerID      int64  `json:"bowler_id"`
	WicketsNumber int32  `json:"wickets_number"`
	WicketType    string `json:"wicket_type"`
	BallNumber    int32  `json:"ball_number"`
	FielderID     *int32 `json:"fielder_id"`
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
		arg.FielderID,
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
		&i.FielderID,
	)
	return i, err
}

const getCricketBall = `
SELECT * FROM balls
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
		&i.BowlingStatus,
		&i.IsCurrentBowler,
	)
	return i, err
}

const getCricketBalls = `
SELECT * FROM balls
WHERE match_id=$1 AND team_id=$2
ORDER BY id
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
SELECT * FROM bats
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
		&i.BattingStatus,
		&i.IsStriker,
		&i.IsCurrentlyBatting,
	)
	return i, err
}

const getCricketPlayersScore = `
SELECT * FROM bats
WHERE match_id = $1 AND team_id = $2
ORDER BY id
`

type GetCricketPlayersScoreParams struct {
	TeamID  int64 `json:"team_id"`
	MatchID int64 `json:"match_id"`
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

const getCricketWickets = `
SELECT * FROM wickets
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
			&i.FielderID,
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
RETURNING *
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
		&i.BattingStatus,
		&i.IsStriker,
		&i.IsCurrentlyBatting,
	)
	return i, err
}

const updateInningRunsScored = `
UPDATE bats
SET runs_scored = runs_scored + $1,
    balls_faced = balls_faced + 1,
    fours = fours + CASE WHEN $1 = 4 THEN 1 ELSE 0 END,
    sixes = sixes + CASE WHEN $1 = 6 THEN 1 ELSE 0 END
WHERE match_id = $5 AND batsman_id = $6
RETURNING id, batsman_id, team_id, match_id, position, runs_scored, balls_faced, fours, sixes;
`

func (q *Queries) UpdateBatsmanScored(ctx context.Context, runsScored, ballsFaced, fours, sixes int32, matchID, batsmanID int64) (models.Bat, error) {
	row := q.db.QueryRowContext(ctx, updateInningRunsScored,
		runsScored,
		ballsFaced,
		fours,
		sixes,
		matchID,
		batsmanID,
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
		&i.BattingStatus,
		&i.IsStriker,
		&i.IsCurrentlyBatting,
	)
	return i, err
}

const updateRegularRunsScored = `
UPDATE bats
SET runs_scored = runs_scored + $1,
    balls_faced = balls_faced + 1,
    fours = fours + CASE WHEN $1 = 4 THEN 1 ELSE 0 END,
    sixes = sixes + CASE WHEN $1 = 6 THEN 1 ELSE 0 END
WHERE match_id = $2 AND batsman_id = $3 AND is_striker=true
RETURNING *;
`

func (q *Queries) UpdateCricketBatsmanScore(ctx context.Context, runsScored int32, matchID, batsmanID int64) (models.Bat, error) {
	row := q.db.QueryRowContext(ctx, updateRegularRunsScored,
		runsScored,
		matchID,
		batsmanID,
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
		&i.BattingStatus,
		&i.IsStriker,
		&i.IsCurrentlyBatting,
	)
	return i, err
}

const updateBowlingStats = `
UPDATE balls
SET runs = runs + $1,
    ball = ball + 1
WHERE match_id = $2 AND bowler_id = $3 AND is_current_bowler=true
RETURNING *;
`

func (q *Queries) UpdateBowlerStats(ctx context.Context, runs int32, matchID, bowlerID int64) (models.Ball, error) {
	row := q.db.QueryRowContext(ctx, updateBowlingStats,
		runs,
		matchID,
		bowlerID,
	)
	var i models.Ball
	err := row.Scan(
		&i.ID,
		&i.BowlerID,
		&i.TeamID,
		&i.MatchID,
		&i.Ball,
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
	SELECT * FROM bats bs
	WHERE bs.match_id = $1 AND bs.batting_status = true;
`

func (q *Queries) GetCurrentPlayingBatsmen(ctx context.Context, matchID int64) ([]models.Bat, error) {
	rows, err := q.db.QueryContext(ctx, getCurrentPlayingBatsman, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batsmen []models.Bat
	for rows.Next() {
		var i models.Bat
		err := rows.Scan(
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

const updateStricketSwapQuery = `
	UPDATE bats
	SET is_striker = NOT is_striker;
	WHERE match_id=$2 AND is_currently_batting=true
	RETURNING *;
`

func (q *Queries) ToggleCricketStricker(ctx context.Context, matchID int64) ([]models.Bat, error) {
	const query = `
		UPDATE bats
		SET is_striker = NOT is_striker
		WHERE match_id = $1 AND is_currently_batting = true
		RETURNING *;
	`

	rows, err := q.db.QueryContext(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batsmen []models.Bat
	for rows.Next() {
		var bat models.Bat
		err := rows.Scan(
			&bat.ID,
			&bat.BatsmanID,
			&bat.TeamID,
			&bat.MatchID,
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

func (q *Queries) UpdateWideball(ctx context.Context, matchID, battingTeamID, bowlerID int64) error {
	_, err := q.db.ExecContext(ctx, `
		UPDATE balls
		SET wide = wide + 1;
		WHERE match_id=$1 AND bowler_id=$2
		RETURNING *;
	`, matchID)
	if err != nil {
		return fmt.Errorf("failed to toggle striker: %w", err)
	}

	return nil
}

const updateWideRun = `
	WITH update_bowler AS (
		UPDATE balls
		SET wide = wide + 1, 
		    runs = runs + 1
		WHERE match_id = $1 AND bowler_id = $2 AND is_current_bowler = true
		RETURNING *
	),
	update_inning_score AS (
		UPDATE cricket_score
		SET score = score + 1
		WHERE match_id = $1 AND team_id = $3
		RETURNING *
	),
	update_batsman AS (
		UPDATE bats
		SET runs_scored = runs_scored + $4
		WHERE match_id = $1 AND is_striker = true
		RETURNING *
	)
	SELECT 
		ub.*, 
		ubl.*, 
		uis.*
	FROM update_batsman ub
	JOIN update_bowler ubl ON ub.match_id = ubl.match_id
	JOIN update_inning_score uis ON ub.match_id = uis.match_id
`

func (q *Queries) UpdateWideRuns(ctx context.Context, matchID, bowlerID, battingTeamID int64, runsScored int32) (*models.Bat, *models.Ball, *models.CricketScore, error) {
	var bowler models.Ball
	var batsman models.Bat
	var inningScore models.CricketScore
	row := q.db.QueryRowContext(ctx, updateWideRun, matchID, bowlerID, battingTeamID, runsScored)
	err := row.Scan(
		&batsman.ID,
		&batsman.BatsmanID,
		&batsman.TeamID,
		&batsman.MatchID,
		&batsman.Position,
		&batsman.RunsScored,
		&batsman.BallsFaced,
		&batsman.Fours,
		&batsman.Sixes,
		&batsman.BattingStatus,
		&batsman.IsStriker,
		&batsman.IsCurrentlyBatting,
		&bowler.ID,
		&bowler.BowlerID,
		&bowler.TeamID,
		&bowler.MatchID,
		&bowler.Ball,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
		&inningScore.ID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.Inning,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return &batsman, &bowler, &inningScore, nil
}

const updateNoBallRun = `
WITH update_bowler AS (
	UPDATE balls
	SET no_ball = no_ball + 1, 
		runs = runs + 1 + $4
	WHERE match_id = $1 AND bowler_id = $2 AND is_current_bowler = true
	RETURNING *
),
update_inning_score AS (
	UPDATE cricket_score
	SET score = score + 1
	WHERE match_id = $1 AND team_id = $3
	RETURNING *
),
update_batsman AS (
	UPDATE bats
	SET runs_scored = runs_scored + $4
	WHERE match_id = $1 AND is_striker = true
	RETURNING *
)
SELECT 
	ub.*, 
	ubl.*, 
	uis.*
FROM update_batsman ub
JOIN update_bowler ubl ON ub.match_id = ubl.match_id
JOIN update_inning_score uis ON ub.match_id = uis.match_id
`

func (q *Queries) UpdateNoBallsRuns(ctx *gin.Context, matchID, bowlerID, battingTeamID int64, runsScored int32) (*models.Bat, *models.Ball, *models.CricketScore, error) {
	var bowler models.Ball
	var batsman models.Bat
	var inningScore models.CricketScore
	row := q.db.QueryRowContext(ctx, updateNoBallRun, matchID, bowlerID, battingTeamID, runsScored)
	err := row.Scan(
		&batsman.ID,
		&batsman.BatsmanID,
		&batsman.TeamID,
		&batsman.MatchID,
		&batsman.Position,
		&batsman.RunsScored,
		&batsman.BallsFaced,
		&batsman.Fours,
		&batsman.Sixes,
		&batsman.BattingStatus,
		&batsman.IsStriker,
		&batsman.IsCurrentlyBatting,
		&bowler.ID,
		&bowler.BowlerID,
		&bowler.TeamID,
		&bowler.MatchID,
		&bowler.Ball,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
		&inningScore.ID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.Inning,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return &batsman, &bowler, &inningScore, nil
}

const addCricketWicket = `
WITH add_wicket AS (
    INSERT INTO wickets (
        match_id,
        team_id,
        batsman_id,
        bowler_id,
        wickets_number,
        wicket_type,
        ball_number,
        fielder_id
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING *
),
update_bowler AS (
    UPDATE balls
    SET 
        wickets = CASE
				WHEN $6 != 'Run Out' THEN wickets + 1
				ELSE wickets
			END,
        ball = ball + 1
    WHERE match_id = $1 AND bowler_id = $4 AND is_current_bowler = true
    RETURNING *
),
update_batsman AS (
    UPDATE bats
    SET 
        balls_faced = balls_faced + 1,
        runs_scored = runs_scored + CASE 
            WHEN $9 > 0 THEN $9
            ELSE 0
        END,
		is_currently_batting = NOT is_currently_batting,
		is_striker = false
    WHERE match_id = $1 AND batsman_id = $3 AND is_currently_batting = true AND is_striker = true
    RETURNING *
)
SELECT 
    w.*,
    b.*,
    ba.*
FROM add_wicket w
JOIN update_bowler b ON w.match_id = b.match_id
JOIN update_batsman ba ON w.match_id = ba.match_id;
`

func (q *Queries) AddCricketWicket(ctx context.Context, matchID, teamID, batsmanID, bowlerID int64, wicketNumber int, wicketType string, ballNumber int, fielderID int64, runsScored int32) error {
	_, err := q.db.ExecContext(ctx, addCricketWicket, matchID, teamID, batsmanID, bowlerID, wicketNumber, wicketType, ballNumber, fielderID, runsScored)
	if err != nil {
		return fmt.Errorf("Failed to exec querys: ", err)
	}
	return nil
}

type BattingScore struct {
	ID                 int64  `json:"id"`
	BatsmanID          int64  `json:"batsman_id"`
	MatchID            int64  `json:"match_id"`
	TeamID             int64  `json:"team_id"`
	Position           string `json:"position"`
	RunsScored         int32  `json:"runs_scored"`
	BallsFaced         int32  `json:"balls_faced"`
	Fours              int32  `json:"fours"`
	Sixes              int32  `json:"sixes"`
	BattingStatus      bool   `json:"batting_status"`
	IsStriker          bool   `json:"is_striker"`
	IsCurrentlyBatting bool   `json:"is_currently_batting"`
}

type BowlingScore struct {
	ID              int64 `json:"id"`
	MatchID         int64 `json:"match_id"`
	TeamID          int64 `json:"team_id"`
	BowlerID        int64 `json:"bowler_id"`
	Ball            int32 `json:"ball"`
	Runs            int32 `json:"runs"`
	Wickets         int32 `json:"wickets"`
	Wide            int32 `json:"wide"`
	NoBall          int32 `json:"no_ball"`
	BowlingStatus   bool  `json:"bowling_status"`
	IsCurrentBowler bool  `json:"is_current_bowler"`
}

type InningScore struct {
	ID                int64  `json:"id"`
	MatchID           int64  `json:"match_id"`
	TeamID            int64  `json:"team_id"`
	Inning            string `json:"inning"`
	Score             int32  `json:"score"`
	Wickets           int32  `json:"wickets"`
	Overs             int32  `json:"overs"`
	RunRate           string `json:"run_rate"`
	TargetRunRate     string `json:"target_run_rate"`
	FollowOn          bool   `json:"follow_on"`
	IsInningCompleted bool   `json:"is_inning_completed"`
	Declared          bool   `json:"declared"`
}

const updateInningScore = `
	WITH update_batsman AS (
		UPDATE bats
		SET runs_scored = runs_scored + $5,
			balls_faced = balls_faced + 1,
			fours = fours + CASE WHEN $5 = 4 THEN 1 ELSE 0 END,
			sixes = sixes + CASE WHEN $5 = 6 THEN 1 ELSE 0 END
		WHERE match_id = $1 AND batsman_id = $3
		RETURNING *
	),
	update_bowler AS (
		UPDATE balls
		SET runs = runs + $5,
			ball = ball + 1
		WHERE match_id = $1 AND bowler_id = $4
		RETURNING *
	),
	update_inning_score AS (
		UPDATE cricket_score
		SET score = score + $5,
			overs = overs + 1
		WHERE match_id = $1 AND team_id = $2
		RETURNING *
	)
	SELECT 
		ub.*, 
		ubl.*, 
		uis.*
	FROM update_batsman ub
	JOIN update_bowler ubl ON ub.match_id = ubl.match_id
	JOIN update_inning_score uis ON ub.match_id = uis.match_id;
`

func (q *Queries) UpdateInningScore(ctx context.Context, matchID, batsmanTeamID, batsmanID, bowlerID int64, runsScored int32) (*models.Bat, *models.Ball, *models.CricketScore, error) {
	var batsman models.Bat
	var bowler models.Ball
	var inningScore models.CricketScore

	row := q.db.QueryRowContext(ctx, updateInningScore, matchID, batsmanTeamID, batsmanID, bowlerID, runsScored)

	err := row.Scan(
		&batsman.ID,
		&batsman.BatsmanID,
		&batsman.TeamID,
		&batsman.MatchID,
		&batsman.Position,
		&batsman.RunsScored,
		&batsman.BallsFaced,
		&batsman.Fours,
		&batsman.Sixes,
		&batsman.BattingStatus,
		&batsman.IsStriker,
		&batsman.IsCurrentlyBatting,
		&bowler.ID,
		&bowler.TeamID,
		&bowler.MatchID,
		&bowler.BowlerID,
		&bowler.Ball,
		&bowler.Runs,
		&bowler.Wickets,
		&bowler.Wide,
		&bowler.NoBall,
		&bowler.BowlingStatus,
		&bowler.IsCurrentBowler,
		&inningScore.ID,
		&inningScore.MatchID,
		&inningScore.TeamID,
		&inningScore.Inning,
		&inningScore.Score,
		&inningScore.Wickets,
		&inningScore.Overs,
		&inningScore.RunRate,
		&inningScore.TargetRunRate,
		&inningScore.FollowOn,
		&inningScore.IsInningCompleted,
		&inningScore.Declared,
	)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return &batsman, &bowler, &inningScore, nil
}

const updateSetBowlerStatus = `
	UPDATE balls
	SET is_current_bowler = NOT is_current_bowler
	WHERE match_id = $1 AND bowler_id = $2
	RETURNING *
`

func (q *Queries) UpdateBowlingBowlerStatus(ctx context.Context, matchID, bowlerID int64) (*models.Ball, error) {
	var currentBowler models.Ball

	row := q.db.QueryRowContext(ctx, updateSetBowlerStatus, matchID, bowlerID)

	err := row.Scan(
		&currentBowler.ID,
		&currentBowler.TeamID,
		&currentBowler.MatchID,
		&currentBowler.BowlerID,
		&currentBowler.Ball,
		&currentBowler.Runs,
		&currentBowler.Wickets,
		&currentBowler.Wide,
		&currentBowler.NoBall,
		&currentBowler.BowlingStatus,
		&currentBowler.IsCurrentBowler,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return &currentBowler, nil
}

const getCurrentBattingBatsmanQuery = `
	SELECT * FROM bats
	WHERE match_id=$1 AND team_id=$2 AND is_currently_batting=true;
`

func (q *Queries) GetCurrentBattingBatsman(ctx context.Context, matchID, batsmanID int64) ([]models.Bat, error) {
	rows, err := q.db.QueryContext(ctx, getCurrentBattingBatsmanQuery, matchID, batsmanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batsmen []models.Bat
	for rows.Next() {
		var bat models.Bat
		err := rows.Scan(
			&bat.ID,
			&bat.BatsmanID,
			&bat.TeamID,
			&bat.MatchID,
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
