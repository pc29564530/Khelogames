package database

import (
	"context"
	"database/sql"
	"encoding/json"
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
	is_current_bowler,
	inning
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *
`

type AddCricketBallParams struct {
	MatchID         int64  `json:"match_id"`
	TeamID          int64  `json:"team_id"`
	BowlerID        int64  `json:"bowler_id"`
	Ball            int32  `json:"ball"`
	Runs            int32  `json:"runs"`
	Wickets         int32  `json:"wickets"`
	Wide            int32  `json:"wide"`
	NoBall          int32  `json:"no_ball"`
	BowlingStatus   bool   `json:"bowling_status"`
	IsCurrentBowler bool   `json:"is_current_bowler"`
	Inning          string `json:"inning"` // inning1 or inning2
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
		arg.Inning,
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
		&i.Inning)
	return i, err
}

const getCricketStricker = `
	SELECT * FROM bats
	WHERE match_id=$1 AND team_id=$2 AND is_currently_batting=true AND is_striker=true AND inning = $3;
`

func (q *Queries) GetCricketStricker(ctx context.Context, matchID, teamID int64, inning string) (*models.Bat, error) {
	row := q.db.QueryRowContext(ctx, getCricketStricker, matchID, teamID, inning)

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
		&i.Inning,
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
	is_currently_batting,
	inning
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;
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
	Inning             string `json:"inning"`
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
		arg.Inning,
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
		&i.Inning,
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
	fielder_id,
	score,
	inning
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
	Score         *int32 `json:"score"`
	Inning        string `json:"inning"`
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
		arg.Score,
		arg.Inning,
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
		&i.Score,
		&i.Inning,
	)
	return i, err
}

const getCricketBall = `
SELECT * FROM balls
WHERE match_id=$1 AND bowler_id=$2 AND inning = $3 LIMIT 1
`

// filter according to the inning also in this case
// not used function
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
		&i.Inning,
	)
	return i, err
}

const getCricketBalls = `
SELECT * FROM balls
WHERE match_id=$1 AND team_id = $2
ORDER BY id, inning
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
			&i.Inning,
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

// Not used function
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
		&i.Inning,
	)
	return i, err
}

const getCricketPlayersScore = `
SELECT * FROM bats
WHERE match_id = $1 AND team_id = $2
ORDER BY id, inning
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
			&i.BattingStatus,
			&i.IsStriker,
			&i.IsCurrentlyBatting,
			&i.Inning,
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
		&i.MatchID,
		&i.TeamID,
		&i.BatsmanID,
		&i.BowlerID,
		&i.WicketsNumber,
		&i.WicketType,
		&i.BallNumber,
		&i.FielderID,
		&i.Score,
		&i.Inning,
	)
	return i, err
}

const getCricketWickets = `
SELECT * FROM wickets
WHERE match_id=$1 AND team_id=$2
ORDER BY id, inning
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
			&i.Score,
			&i.Inning,
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
WHERE match_id = $6 AND bowler_id = $7 AND team_id=$8 AND inning = $9
RETURNING *;
`

type UpdateCricketBowlerParams struct {
	Ball     int32  `json:"ball"`
	Runs     int32  `json:"runs"`
	Wickets  int32  `json:"wickets"`
	Wide     int32  `json:"wide"`
	NoBall   int32  `json:"no_ball"`
	MatchID  int64  `json:"match_id"`
	BowlerID int64  `json:"bowler_id"`
	TeamID   int64  `json:"team_id"`
	Inning   string `json:"inning"`
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
		arg.Inning,
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
		&i.Inning,
	)
	return i, err
}

const updateCricketRunsScored = `
UPDATE bats
SET runs_scored = $1,
    balls_faced = $2,
    fours = $3,
    sixes = $4
WHERE match_id = $5 AND batsman_id = $6 AND team_id=$7 AND inning = $8
RETURNING *
`

type UpdateCricketRunsScoredParams struct {
	RunsScored int32  `json:"runs_scored"`
	BallsFaced int32  `json:"balls_faced"`
	Fours      int32  `json:"fours"`
	Sixes      int32  `json:"sixes"`
	MatchID    int64  `json:"match_id"`
	BatsmanID  int64  `json:"batsman_id"`
	TeamID     int64  `json:"team_id"`
	Inning     string `json:"inning"`
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
		arg.Inning,
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
		&i.Inning,
	)
	return i, err
}

// const updateInningRunsScored = `
// UPDATE bats
// SET runs_scored = runs_scored + $1,
//     balls_faced = balls_faced + 1,
//     fours = fours + CASE WHEN $1 = 4 THEN 1 ELSE 0 END,
//     sixes = sixes + CASE WHEN $1 = 6 THEN 1 ELSE 0 END
// WHERE match_id = $5 AND batsman_id = $6 AND inning = $7
// RETURNING *;
// `

// func (q *Queries) UpdateBatsmanScored(ctx context.Context, runsScored, ballsFaced, fours, sixes int32, matchID, batsmanID int64, inning string) (models.Bat, error) {
// 	row := q.db.QueryRowContext(ctx, updateInningRunsScored,
// 		runsScored,
// 		ballsFaced,
// 		fours,
// 		sixes,
// 		matchID,
// 		batsmanID,
// 		inning,
// 	)
// 	var i models.Bat
// 	err := row.Scan(
// 		&i.ID,
// 		&i.BatsmanID,
// 		&i.TeamID,
// 		&i.MatchID,
// 		&i.Position,
// 		&i.RunsScored,
// 		&i.BallsFaced,
// 		&i.Fours,
// 		&i.Sixes,
// 		&i.BattingStatus,
// 		&i.IsStriker,
// 		&i.IsCurrentlyBatting,
// 		&i.Inning,
// 	)
// 	return i, err
// }

// const updateRegularRunsScored = `
// UPDATE bats
// SET runs_scored = runs_scored + $1,
//     balls_faced = balls_faced + 1,
//     fours = fours + CASE WHEN $1 = 4 THEN 1 ELSE 0 END,
//     sixes = sixes + CASE WHEN $1 = 6 THEN 1 ELSE 0 END
// WHERE match_id = $2 AND batsman_id = $3 AND is_striker=true AND inning = $4
// RETURNING *;
// `

// func (q *Queries) UpdateCricketBatsmanScore(ctx context.Context, runsScored int32, matchID, batsmanID int64, inning string) (models.Bat, error) {
// 	row := q.db.QueryRowContext(ctx, updateRegularRunsScored,
// 		runsScored,
// 		matchID,
// 		batsmanID,
// 		inning,
// 	)
// 	var i models.Bat
// 	err := row.Scan(
// 		&i.ID,
// 		&i.BatsmanID,
// 		&i.TeamID,
// 		&i.MatchID,
// 		&i.Position,
// 		&i.RunsScored,
// 		&i.BallsFaced,
// 		&i.Fours,
// 		&i.Sixes,
// 		&i.BattingStatus,
// 		&i.IsStriker,
// 		&i.IsCurrentlyBatting,
// 		&i.Inning,
// 	)
// 	return i, err
// }

const updateBowlingStats = `
UPDATE balls
SET runs = runs + $1,
    ball = ball + 1
WHERE match_id = $2 AND bowler_id = $3 AND is_current_bowler=true AND inning = $4
RETURNING *;
`

func (q *Queries) UpdateBowlerStats(ctx context.Context, runs int32, matchID, bowlerID int64, inning string) (models.Ball, error) {
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
		&i.Inning,
	)
	return i, err
}

const getCurrentPlayingBatsman = `
	SELECT * FROM bats bs
	WHERE bs.match_id = $1 AND bs.batting_status = true AND inning = $2;
`

func (q *Queries) GetCurrentPlayingBatsmen(ctx context.Context, matchID int64, inning string) ([]models.Bat, error) {
	rows, err := q.db.QueryContext(ctx, getCurrentPlayingBatsman, matchID, inning)
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
			&i.Inning,
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

func (q *Queries) ToggleCricketStricker(ctx context.Context, matchID int64, inning string) ([]models.Bat, error) {
	const query = `
		UPDATE bats
		SET is_striker = NOT is_striker
		WHERE match_id = $1 AND is_currently_batting = true AND inning = $2
		RETURNING *;
	`

	rows, err := q.db.QueryContext(ctx, query, matchID, inning)
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
			&bat.Inning,
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
	WITH update_bowler AS (
		UPDATE balls
		SET 
			wide = wide + 1, 
			runs = runs + $4
		WHERE 
			match_id = $1 
			AND bowler_id = $2 
			AND team_id = (
				SELECT CASE 
					WHEN home_team_id = $3 THEN away_team_id 
					ELSE home_team_id 
				END AS bowler_team_id
				FROM matches
				WHERE id = $1
			) 
			AND is_current_bowler = true 
			AND inning = $5
		RETURNING *
	),
	update_inning_score AS (
		UPDATE cricket_score
		SET 
			score = score + $4 + 1
		WHERE 
			match_id = $1
			AND team_id = $3
			AND inning = $5
		RETURNING *
	),
	update_batsman AS (
		UPDATE bats
		SET 
			runs_scored = runs_scored + $4
		WHERE 
			match_id = $1 
			AND is_striker = true 
			AND inning = $5
		RETURNING *
	)
	SELECT 
		ub.*, 
		ubl.*, 
		uis.*
	FROM update_batsman ub
	JOIN update_bowler ubl 
		ON ub.match_id = ubl.match_id 
		AND ub.inning = ubl.inning
	JOIN update_inning_score uis 
		ON ub.match_id = uis.match_id 
		AND ub.inning = uis.inning
	WHERE ubl.team_id = (
		SELECT CASE 
			WHEN home_team_id = $3 THEN away_team_id 
			ELSE home_team_id 
		END AS bowler_team_id
		FROM matches
		WHERE id = $1
	);
`

func (q *Queries) UpdateWideRuns(ctx context.Context, matchID, bowlerID, battingTeamID int64, runsScored int32, inning string) (*models.Bat, *models.Ball, *models.CricketScore, error) {
	var bowler models.Ball
	var batsman models.Bat
	var inningScore models.CricketScore
	row := q.db.QueryRowContext(ctx, updateWideRun, matchID, bowlerID, battingTeamID, runsScored, inning)
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
		&batsman.Inning,
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
		&bowler.Inning,
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
	WHERE match_id = $1 AND bowler_id = $2 AND is_current_bowler = true  AND inning = $5
	RETURNING *
),
update_inning_score AS (
	UPDATE cricket_score
	SET score = score + 1
	WHERE match_id = $1 AND team_id = $3 AND inning = $5
	RETURNING *
),
update_batsman AS (
	UPDATE bats
	SET runs_scored = runs_scored + $4
	WHERE match_id = $1 AND is_striker = true AND inning = $5
	RETURNING *
)
SELECT 
	ub.*, 
	ubl.*, 
	uis.*
FROM update_batsman ub
JOIN update_bowler ubl ON ub.match_id = ubl.match_id AND ub.inning = ubl.inning
JOIN update_inning_score uis ON ub.match_id = uis.match_id AND ub.inning = uis.inning
`

func (q *Queries) UpdateNoBallsRuns(ctx *gin.Context, matchID, bowlerID, battingTeamID int64, runsScored int32, inning string) (*models.Bat, *models.Ball, *models.CricketScore, error) {
	var bowler models.Ball
	var batsman models.Bat
	var inningScore models.CricketScore
	row := q.db.QueryRowContext(ctx, updateNoBallRun, matchID, bowlerID, battingTeamID, runsScored, inning)
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
		&batsman.Inning,
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
		&bowler.Inning,
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

// Enhance about the no ball
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
        fielder_id,
        score,
		inning
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING *
),
update_out_batsman AS (
    UPDATE bats
    SET balls_faced = balls_faced + 1,
        runs_scored = runs_scored + (CASE WHEN is_striker THEN (CASE WHEN $10 > 0 THEN $10 ELSE 0 END) ELSE 0 END),
        is_currently_batting = false,
        is_striker = false
    WHERE match_id = $1 
      AND batsman_id = $3 
      AND team_id = $2
	  AND inning = $10
    RETURNING *
),
update_not_out_batsman AS (
    UPDATE bats
    SET balls_faced = balls_faced + 1,
        runs_scored = runs_scored + (CASE WHEN is_striker THEN (CASE WHEN $10 > 0 THEN $10 ELSE 0 END) ELSE 0 END)
    WHERE match_id = $1 
      AND team_id = $2 
      AND batsman_id <> $3 
      AND is_currently_batting = true
	  AND inning = $10
    RETURNING *
),
update_bowler AS (
    UPDATE balls
    SET wickets = CASE
                    WHEN $6 != 'Run Out' THEN wickets + 1
                    ELSE wickets
                  END,
        runs = runs + (CASE WHEN $10 > 0 THEN $10 ELSE 0 END),
        ball = ball + 1
    WHERE match_id = $1 
      AND bowler_id = $4 
      AND is_current_bowler = true
	  AND inning = $10
    RETURNING *
),
update_inning_score AS (
    UPDATE cricket_score
    SET overs = overs + 1,
        wickets = wickets + 1,
        score = score + (CASE WHEN $10 > 0 THEN $10 ELSE 0 END)
    WHERE match_id = $1 
      AND team_id = $2
	  AND inning = $10
    RETURNING *
)
SELECT 
	o.*,
	n.*,
	b.*,
	sc.*,
    w.*
FROM add_wicket w
JOIN update_out_batsman o ON w.match_id = o.match_id AND w.team_id = o.team_id AND o.inning = w.inning
JOIN update_not_out_batsman n ON w.match_id = n.match_id AND w.team_id = n.team_id AND n.inning = w.inning
JOIN update_bowler b ON w.match_id = b.match_id AND w.bowler_id = b.bowler_id AND b.inning = w.inning
JOIN update_inning_score sc ON w.match_id = sc.match_id AND w.team_id = sc.team_id AND sc.inning = w.inning;
`

func (q *Queries) AddCricketWicket(ctx context.Context, matchID, teamID, batsmanID, bowlerID int64, wicketNumber int, wicketType string, ballNumber int, fielderID int64, score int32, runsScored int32, inning string) (*models.Bat, *models.Bat, *models.Ball, *models.CricketScore, *models.Wicket, error) {
	var outBatsman models.Bat
	var notOutBatsman models.Bat
	var bowler models.Ball
	var inningScore models.CricketScore
	var wickets models.Wicket

	row := q.db.QueryRowContext(ctx, addCricketWicket, matchID, teamID, batsmanID, bowlerID, wicketNumber, wicketType, ballNumber, fielderID, score, runsScored, inning)
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
        wickets_number,
        wicket_type,
        ball_number,
        fielder_id,
        score,
		inning
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING *
),
update_out_batsman AS (
    UPDATE bats
    SET balls_faced = balls_faced + 1,
        runs_scored = runs_scored + (CASE WHEN is_striker THEN (CASE WHEN $10 > 0 THEN $10 ELSE 0 END) ELSE 0 END),
        is_currently_batting = false,
        is_striker = false
    WHERE match_id = $1 
      AND batsman_id = $3 
      AND team_id = $2
	  AND inning = $10
    RETURNING *
),
update_not_out_batsman AS (
    UPDATE bats
    SET balls_faced = balls_faced + 1,
        runs_scored = runs_scored + (CASE WHEN is_striker THEN (CASE WHEN $10 > 0 THEN $10 ELSE 0 END) ELSE 0 END)
    WHERE match_id = $1 
      AND team_id = $2 
      AND batsman_id <> $3 
      AND is_currently_batting = true
	  AND inning = $10
    RETURNING *
),
update_bowler AS (
    UPDATE balls
    SET wickets = CASE
                    WHEN $6 != 'Run Out' THEN wickets + 1
                    ELSE wickets
                  END,
        runs = runs + (CASE WHEN $10 > 0 THEN $10 ELSE 0 END),
        ball = ball,
		wide = wide + (CASE WHEN $11 = 'wide' THEN 1 ELSE 0 END),
        no_ball = no_ball + (CASE WHEN $11 = 'no_ball' THEN 1 ELSE 0 END)
    WHERE match_id = $1 
      AND bowler_id = $4 
      AND is_current_bowler = true
	  AND inning = $10
    RETURNING *
),
update_inning_score AS (
    UPDATE cricket_score
    SET overs = overs,
        wickets = wickets + 1,
        score = score + (CASE WHEN $10 > 0 THEN $10 ELSE 0 END)
    WHERE match_id = $1 
      AND team_id = $2
	  AND inning = $10
    RETURNING *
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

func (q *Queries) AddCricketWicketWithBowlType(ctx context.Context, matchID, teamID, batsmanID, bowlerID int64, wicketNumber int, wicketType string, ballNumber int, fielderID int64, score int32, runsScored int32, bowlType string, inning string) (*models.Bat, *models.Bat, *models.Ball, *models.CricketScore, *models.Wicket, error) {
	var outBatsman models.Bat
	var notOutBatsman models.Bat
	var bowler models.Ball
	var inningScore models.CricketScore
	var wickets models.Wicket

	row := q.db.QueryRowContext(ctx, addCricketWicketWithBowlType, matchID, teamID, batsmanID, bowlerID, wicketNumber, wicketType, ballNumber, fielderID, score, runsScored, bowlType, inning)
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil, nil, nil
		}
		return nil, nil, nil, nil, nil, fmt.Errorf("Failed to scan query: ", err)
	}
	return &outBatsman, &notOutBatsman, &bowler, &inningScore, &wickets, nil
}

const updateInningEndStatus = `
WITH update_inning AS (
	UPDATE cricket_score
	SET is_inning_completed = true
	WHERE match_id = $1 AND team_id = $2 AND inning = $3
	RETURNING *
),
update_batsman AS (
	UPDATE bats
	SET is_striker = false
	WHERE match_id = $1 AND team_id = $2 AND is_striker = true AND inning = $3
	RETURNING *
),
update_bowler AS (
	UPDATE balls
	SET is_current_bowler = false
	WHERE match_id = $1 AND is_current_bowler = true AND inning = $3
	RETURNING *
)
SELECT 
	ui.*,
	ub.*,
	ubl.*
FROM update_batsman ub
JOIN update_bowler AS ubl ON ub.match_id = ubl.match_id AND ub.team_id = ubl.team_id AND ub.inning = ubl.inning
JOIN update_inning AS ui ON ub.match_id = ui.match_id AND ui.team_id = ub.team_id AND ui.inning = ub.inning
`

func (q *Queries) UpdateInningEndStatus(ctx context.Context, matchID, batsmanTeamID int64, inning string) (*models.CricketScore, *models.Bat, *models.Ball, error) {
	var inningScore models.CricketScore
	var batsmanScore models.Bat
	var bowler models.Ball

	row := q.db.QueryRowContext(ctx, updateInningEndStatus, matchID, batsmanTeamID)

	err := row.Scan(
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
		&batsmanScore.ID,
		&batsmanScore.BatsmanID,
		&batsmanScore.TeamID,
		&batsmanScore.MatchID,
		&batsmanScore.Position,
		&batsmanScore.RunsScored,
		&batsmanScore.BallsFaced,
		&batsmanScore.Fours,
		&batsmanScore.Sixes,
		&batsmanScore.BattingStatus,
		&batsmanScore.IsStriker,
		&batsmanScore.IsCurrentlyBatting,
		&batsmanScore.Inning,
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
		&bowler.Inning,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil
		}
		return nil, nil, nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return &inningScore, &batsmanScore, &bowler, nil
}

const updateInningScore = `
	WITH update_batsman AS (
		UPDATE bats
		SET runs_scored = runs_scored + $5,
			balls_faced = balls_faced + 1,
			fours = fours + CASE WHEN $5 = 4 THEN 1 ELSE 0 END,
			sixes = sixes + CASE WHEN $5 = 6 THEN 1 ELSE 0 END
		WHERE match_id = $1 AND team_id = $2 AND batsman_id = $3 AND inning = $6
		RETURNING *
	),
	get_bowling_team AS (
		SELECT CASE
			WHEN m.home_team_id = $2 THEN m.away_team_id 
				ELSE m.home_team_id 
			END AS bowler_team_id
		FROM matches AS m WHERE id = $1
	),
	update_bowler AS (
		UPDATE balls
		SET runs = runs + $5,
			ball = ball + 1
		WHERE match_id = $1 AND team_id = (SELECT bowler_team_id FROM get_bowling_team)  AND bowler_id = $4 AND inning = $6
		RETURNING *
	),
	update_inning_score AS (
		UPDATE cricket_score
		SET score = score + $5,
			overs = overs + 1
		WHERE match_id = $1 AND team_id = $2 AND inning = $6
		RETURNING *
	)
	SELECT 
		ub.*, 
		ubl.*, 
		uis.*
	FROM update_batsman ub
	JOIN update_bowler ubl ON ub.match_id = ubl.match_id AND ubl.team_id =  (SELECT bowler_team_id FROM get_bowling_team)
		AND ub.inning = ubl.inning
	JOIN update_inning_score uis ON ub.match_id = uis.match_id AND ub.team_id = uis.team_id AND ub.inning = ubl.inning;
`

func (q *Queries) UpdateInningScore(ctx context.Context, matchID, batsmanTeamID, batsmanID, bowlerID int64, runsScored int32, inning string) (*models.Bat, *models.Ball, *models.CricketScore, error) {
	var batsman models.Bat
	var bowler models.Ball
	var inningScore models.CricketScore
	row := q.db.QueryRowContext(ctx, updateInningScore, matchID, batsmanTeamID, batsmanID, bowlerID, runsScored, inning)

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
		&batsman.Inning,
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
		&bowler.Inning,
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
	WHERE match_id = $1 AND team_id = $2 AND bowler_id = $3 AND inning = $4
	RETURNING *
`

func (q *Queries) UpdateBowlingBowlerStatus(ctx context.Context, matchID, teamID, bowlerID int64, inning string) (*models.Ball, error) {
	var currentBowler models.Ball

	row := q.db.QueryRowContext(ctx, updateSetBowlerStatus, matchID, bowlerID, inning)

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
		&currentBowler.Inning,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return &currentBowler, nil
}

const getCurrentBattingBatsmanQuery = `
	SELECT * FROM bats
	WHERE match_id=$1 AND team_id=$2 AND is_currently_batting=true AND inning = $3;
`

func (q *Queries) GetCurrentBattingBatsman(ctx context.Context, matchID, teamID int64, inning string) ([]models.Bat, error) {
	rows, err := q.db.QueryContext(ctx, getCurrentBattingBatsmanQuery, matchID, teamID, inning)
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
			&bat.Inning,
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
				'name', tm.name, 
				'slug', tm.slug, 
				'short_name', tm.shortname, 
				'admin', tm.admin, 
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
					'batsman_id', bt.batsman_id,
					'player', JSON_BUILD_OBJECT('id',pl.id,'username',pl.username, 'name', pl.player_name, 'slug', pl.slug, 'short_name',pl.short_name, 'country', pl.country, 'positions', pl.positions, 'media_url', pl.media_url),
					'position', bt.position, 
					'runs_scored', bt.runs_scored, 
					'balls_faced', bt.balls_faced, 
					'fours', bt.fours, 
					'sixes', bt.sixes, 
					'batting_status', bt.batting_status, 
					'is_striker', bt.is_striker, 
					'is_currently_batting', bt.is_currently_batting,
					'inning', bt.inning
				)
        	)
    	) AS team_data
	FROM bats bt
	JOIN players AS pl ON pl.id = bt.batsman_id
	JOIN teams AS tm ON tm.id = bt.team_id
	WHERE bt.match_id = $1 AND bt.team_id = $2 AND bt.inning = $3 AND bt.is_currently_batting = true
	GROUP BY tm.id;
`

func (q *Queries) GetCurrentBatsman(ctx context.Context, matchID, teamID int64, inning string) (interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getCurrentBatsmanQuery, matchID, teamID, inning)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	var jsonBytes []byte
	if rows.Next() {
		if err := rows.Scan(&jsonBytes); err != nil {
			return nil, fmt.Errorf("failed to scan json data: %w", err)
		}
	}

	var currentBatsman interface{}

	err = json.Unmarshal(jsonBytes, &currentBatsman)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal batsman: ", err)
	}

	return currentBatsman, nil
}

const getCurrentBowlerQuery = `
	SELECT 
    	JSON_BUILD_OBJECT(
			'team', JSON_BUILD_OBJECT(
				'id', tm.id, 
				'name', tm.name, 
				'slug', tm.slug, 
				'short_name', tm.shortname, 
				'admin', tm.admin, 
				'media_url', tm.media_url, 
				'gender', tm.gender, 
				'national', tm.national, 
				'country', tm.country, 
				'type', tm.type, 
				'player_count', tm.player_count, 
				'game_id', tm.game_id
			),
        	'bowler',
				JSON_BUILD_OBJECT(
					'id', bl.id, 
					'team_id', bl.team_id,
					'match_id', bl.match_id,
					'bowler_id', bl.bowler_id,
					'player', JSON_BUILD_OBJECT('id',pl.id,'username',pl.username, 'name', pl.player_name, 'slug', pl.slug, 'short_name',pl.short_name, 'country', pl.country, 'positions', pl.positions, 'media_url', pl.media_url),
					'runs', bl.runs, 
					'ball', bl.ball, 
					'wickets', bl.wickets, 
					'wide', bl.wide, 
					'no_ball', bl.no_ball,
					'bowling_status', bl.bowling_status,
					'is_current_bowler', bl.is_current_bowler,
					'inning', bl.inning,
				)
    	) AS team_data
	FROM balls bl
	JOIN players AS pl ON pl.id = bl.bowler_id
	JOIN teams AS tm ON tm.id = bl.team_id
	WHERE bl.match_id = $1 AND bl.team_id = $2 AND AND bl.inning = $3 AND bl.is_current_bowler = true
`

func (q *Queries) GetCurrentBowler(ctx context.Context, matchID, teamID int64, inning string) (interface{}, error) {
	var jsonBytes []byte
	row := q.db.QueryRowContext(ctx, getCurrentBowlerQuery, matchID, teamID, inning)
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
