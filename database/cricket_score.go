package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const getCricketBatsmanScoreByTeamID = `
	SELECT * FROM bats b
	JOIN teams t ON t.id = b.team_id
	WHERE t.public_id=$1
`

func (q *Queries) GetCricketBatsmanScoreByTeamID(ctx context.Context, teamPublicID uuid.UUID) (*[]models.Bat, error) {
	rows, err := q.db.QueryContext(ctx, getCricketBatsmanScoreByTeamID, teamPublicID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query: ", err)
	}
	var batsmanScore []models.Bat

	for rows.Next() {
		var i models.Bat
		err := rows.Scan(
			&i.ID,
			&i.PublicID,
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
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("Failed to scan the query: ", err)
		}
	}

	return &batsmanScore, nil
}

const getCricketBowlerScoreByTeamID = `
	SELECT * FROM balls b
	JOIN teams t ON t.id = b.team_id
	WHERE t.public_id=$1
`

func (q *Queries) GetCricketBowlerScoreByTeamID(ctx context.Context, teamPublicID uuid.UUID) (*[]models.Ball, error) {
	rows, err := q.db.QueryContext(ctx, getCricketBowlerScoreByTeamID, teamPublicID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query: ", err)
	}
	var bowlerScore []models.Ball

	for rows.Next() {
		var i models.Ball
		err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.MatchID,
			&i.TeamID,
			&i.BowlerID,
			&i.Ball,
			&i.Runs,
			&i.Wickets,
			&i.Wide,
			&i.NoBall,
			&i.InningNumber,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("Failed to scan the query: ", err)
		}
	}

	return &bowlerScore, nil
}

const getCricketScore = `
SELECT * FROM cricket_score cs
WHERE cs.match_id=$1 AND cs.team_id=$2
`

func (q *Queries) GetCricketScore(ctx context.Context, matchID, teamID int32) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, getCricketScore, matchID, teamID)
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
	)
	return i, err
}

const getCricketScoreByPublicID = `
SELECT * FROM cricket_score cs
JOIN matches m ON m.id = cs.match_id
JOIN teams t ON t.id = cs.team_id
WHERE m.public_id=$1 AND t.public_id=$2
`

func (q *Queries) GetCricketScoreByPublicID(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, getCricketScoreByPublicID, matchPublicID, teamPublicID)
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
	)
	return i, err
}

const getCricketScores = `
SELECT * FROM cricket_score
WHERE m.match_id = $1
ORDER BY inning_number;
`

func (q *Queries) GetCricketScores(ctx context.Context, matchID int32) ([]models.CricketScore, error) {
	row, err := q.db.QueryContext(ctx, getCricketScores, matchID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get all cricket scores inning: %w", err)
	}
	defer row.Close()

	var cricketScores []models.CricketScore
	for row.Next() {
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
		)
		if err != nil {
			return nil, fmt.Errorf("Not able to scan row: %w", err)
		}
		cricketScores = append(cricketScores, i)
	}

	// Check for row iteration errors
	if err = row.Err(); err != nil {
		return nil, fmt.Errorf("Row iteration error: %w", err)
	}

	return cricketScores, nil
}

const newCricketScore = `
WITH matchID AS (
	SELECT * FROM matches WHERE public_id = $1
),
teamID AS (
	SELECT * FROM teams WHERE public_id = $2
)
INSERT INTO cricket_score (
    match_id,
    team_id,
    inning_number,
    score,
    wickets,
    overs,
    run_rate,
    target_run_rate,
	follow_on,
	is_inning_completed,
	declared
)
SELECT 
	matchID.id,
	teamID.id,
	$3,
	$4,
	$5,
	$6,
	CAST($7 AS numeric(5,2)),
	CAST($8 AS numeric(5,2)),
	$9,
	$10,
	$11
FROM matchID, teamID
RETURNING *;
`

type NewCricketScoreParams struct {
	MatchID           int32  `json:"match_id"`
	TeamID            int32  `json:"team_id"`
	InningNumber      int    `json:"inning_number"`
	Score             int32  `json:"score"`
	Wickets           int32  `json:"wickets"`
	Overs             int32  `json:"overs"`
	RunRate           string `json:"run_rate"`
	TargetRunRate     string `json:"target_run_rate"`
	FollowOn          bool   `json:"follow_on"`
	IsInningCompleted bool   `json:"is_inning_completed"`
	Declared          bool   `json:"declared"`
}

func (q *Queries) NewCricketScore(ctx context.Context, arg NewCricketScoreParams) (models.CricketScore, error) {

	row := q.db.QueryRowContext(ctx, newCricketScore,
		arg.MatchID,
		arg.TeamID,
		arg.InningNumber,
		arg.Score,
		arg.Wickets,
		arg.Overs,
		arg.RunRate,
		arg.TargetRunRate,
		arg.FollowOn,
		arg.IsInningCompleted,
		arg.Declared,
	)
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
	)
	return i, err
}

const updateCricketInnings = `
UPDATE cricket_score cs
SET inning_number=$3
FROM update_inning ui
JOIN matches m ON m.id = ui.match_id
JOIN teams t ON t.id = ui.team_id
WHERE m.public_id=$1 AND t.public_id=$2 AND ui.match_id = m.id
RETURNING cs.*
`

func (q *Queries) UpdateCricketInnings(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, inningNumber int) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketInnings, matchPublicID, teamPublicID, inningNumber)
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
	)
	return i, err
}

const updateCricketOvers = `
UPDATE cricket_score cs
SET overs = (
        SELECT COALESCE(SUM(bt.balls_faced), 0)
        FROM bats bt
        WHERE bt.match_id = cs.match_id AND bt.inning_number= cs.inning_number AND bt.team_id = cs.team_id
    )
FROM update_over uo
JOIN matches m ON m.id = uo.match_id
JOIN teams t ON t.id = uo.team_id
WHERE m.public_id = $1 AND t.public_id= $2 AND cs.inning_number = $3
RETURNING *;
`

func (q *Queries) UpdateCricketOvers(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, inningNumber int) (*models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketOvers, matchPublicID, teamPublicID, inningNumber)
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

const updateCricketScore = `
UPDATE cricket_score cs
SET score = (
        SELECT SUM(bt.runs_scored) + SUM(bl.wide + bl.no_ball)
        FROM bats bt
		LEFT JOIN balls AS bl ON bt.match_id=bl.match_id AND bl.team_id =  $3 AND bl.bowling_status=true
        WHERE bt.match_id = cs.match_id AND bt.inning_number = cs.inning_number AND bt.team_id=cs.team_id
        GROUP BY (bt.match_id, bt.inning_number, bt.team_id)
    )
WHERE cs.match_id=$1 AND cs.team_id=$2 AND cs.inning_number= $4
RETURNING *;
`

type UpdateCricketScoreParams struct {
	MatchID       int64 `json:"match_id"`
	BattingTeamID int64 `json:"batting_team_id"`
	BowlingTeamID int64 `json:"bowling_team_id"`
	InningNumber  int   `json:"inning_number"`
}

func (q *Queries) UpdateCricketScore(ctx context.Context, arg UpdateCricketScoreParams) (*models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketScore, arg.MatchID, arg.BattingTeamID, arg.BowlingTeamID, arg.InningNumber)
	var i models.CricketScore
	err := row.Scan(
		&i.ID,
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, err
}

// const updateCricketWickets = `
// UPDATE cricket_score cs
// SET wickets = (
//         SELECT COUNT(*)
//         FROM wickets w
//         WHERE w.match_id = cs.match_id AND w.inning_number= cs.inning_numberAND w.team_id = cs.team_id
//     )
// WHERE cs.match_id = $1 AND cs.inning_number= $2 cs.team_id = $3
// RETURNING *;
// `

// type UpdateCricketWicketsParams struct {
// 	MatchID int64  `json:"match_id"`
// 	inning_number string `json:"inning"`
// 	TeamID  int64  `json:"team_id"`
// }

// func (q *Queries) UpdateCricketWickets(ctx context.Context, arg UpdateCricketWicketsParams) (models.CricketScore, error) {
// 	row := q.db.QueryRowContext(ctx, updateCricketWickets, arg.MatchID, arg.InningNumber, arg.TeamID)
// 	var i models.CricketScore
// 	err := row.Scan(
// 		&i.ID,
// 		&i.MatchID,
// 		&i.TeamID,
// 		&i.InningNumber,
// 		&i.Score,
// 		&i.Wickets,
// 		&i.Overs,
// 		&i.RunRate,
// 		&i.TargetRunRate,
// 		&i.FollowOn,
// 		&i.IsInningCompleted,
// 		&i.Declared,
// 	)
// 	return i, err
// }

// const updateCricketEndInnings = `
// UPDATE cricket_score
// SET is_inning_completed=true
// WHERE match_id=$1 AND team_id=$2 AND inning=$3
// RETURNING *
// `

// func (q *Queries) UpdateCricketEndInnings(ctx context.Context, matchID, teamID int64, inning_numberstring) (models.CricketScore, error) {
// 	row := q.db.QueryRowContext(ctx, updateCricketEndInnings, matchID, teamID, inning)
// 	var i models.CricketScore
// 	err := row.Scan(
// 		&i.ID,
// 		&i.MatchID,
// 		&i.TeamID,
// 		&i.InningNumber,
// 		&i.Score,
// 		&i.Wickets,
// 		&i.Overs,
// 		&i.RunRate,
// 		&i.TargetRunRate,
// 		&i.FollowOn,
// 		&i.IsInningCompleted,
// 		&i.Declared,
// 	)
// 	return i, err
// }
