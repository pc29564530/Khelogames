package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addPlayerBattingStats = `
	WITH resolve_ids AS (
		SELECT p.id AS player_id FROM players p
		WHERE p.public_id = $1
	)
	INSERT INTO player_batting_stats (
		player_id,
		match_type,
		matches,
		innings,
		runs,
		balls
		sixes,
		fours,
		fifties,
		hundreds,
		best_score,
		average,
		strike_rate,
		created_at,
		updated_at
	)
	SELECT 
		ri.player_id, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	FROM resolve_ids ri
	RETURNING *;
`

func (q *Queries) AddPlayerCricketStats(ctx context.Context, playerPublicID uuid.UUID, matchType string, matches, battingInnings int32, runs, ballsFaced, fours, sixes, fifties, hundreds, bestScore int, bowlingInnings int32, wickets, runsConceded, ballsBowled, fourWickets, fiveWickets int) (*models.CricketPlayerStats, error) {

	rows := q.db.QueryRowContext(ctx, addPlayerBattingStats, playerPublicID, matchType, matches, battingInnings, runs, ballsFaced, fours, sixes, fifties, hundreds, bestScore, bowlingInnings, wickets, runsConceded, ballsBowled, fourWickets, fiveWickets)

	var i models.CricketPlayerStats
	err := rows.Scan(
		&i.ID,
		&i.PublicID,
		&i.PlayerID,
		&i.MatchType,
		&i.Matches,
		&i.BattingInnings,
		&i.Runs,
		&i.Balls,
		&i.Sixes,
		&i.Fours,
		&i.Fifties,
		&i.Hundreds,
		&i.BowlingInnings,
		&i.Wickets,
		&i.BallsBowled,
		&i.RunsConceded,
		&i.FourWickets,
		&i.FiveWickets,
		&i.CreatedAT,
		&i.UpdatedAT,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan the query data: ", err)
	}
	return &i, nil
}

func (q *Queries) AddPlayerBattingStats(ctx context.Context, playerPublicID uuid.UUID, matchType string, totalMatches, totalInnings, runs, balls, fours, sixes, fifties, hundreds, bestScore int, average, strikeRate string) (*models.PlayerBattingStats, error) {

	rows := q.db.QueryRowContext(ctx, addPlayerBattingStats, playerPublicID, matchType, totalMatches, totalInnings, runs, balls, sixes, fours, fifties, hundreds, bestScore, average, strikeRate)

	var i models.PlayerBattingStats
	err := rows.Scan(
		&i.ID,
		&i.PublicID,
		&i.PlayerID,
		&i.MatchType,
		&i.Matches,
		&i.Innings,
		&i.Runs,
		&i.Balls,
		&i.Sixes,
		&i.Fours,
		&i.Fifties,
		&i.Hundreds,
		&i.BestScore,
		&i.Average,
		&i.StrikeRate,
		&i.CreatedAT,
		&i.UpdatedAT,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan the query data: ", err)
	}
	return &i, nil
}

const getPlayerCricketStatsByMatchType = `
	SELECT * FROM player_cricket_stats pcs
	JOIN players p ON p.id = pcs.player_id
	WHERE p.public_id=$1
`

func (q *Queries) GetPlayerCricketStatsByMatchType(context context.Context, playerPublicID uuid.UUID) (*[]models.CricketPlayerStats, error) {
	rows, err := q.db.QueryContext(context, getCricketPlayerBattingStats, playerPublicID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query database: ", err)
	}
	defer rows.Close()

	var playerStats []models.CricketPlayerStats

	for rows.Next() {
		var i models.CricketPlayerStats
		err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.PlayerID,
			&i.MatchType,
			&i.Matches,
			&i.BattingInnings,
			&i.Runs,
			&i.Balls,
			&i.Sixes,
			&i.Fours,
			&i.Fifties,
			&i.Hundreds,
			&i.BowlingInnings,
			&i.Wickets,
			&i.BallsBowled,
			&i.RunsConceded,
			&i.FourWickets,
			&i.FiveWickets,
			&i.CreatedAT,
			&i.UpdatedAT,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("Failed to scan the query: ", err)
		}

		playerStats = append(playerStats, i)

	}
	return &playerStats, nil
}

const getCricketPlayerBowlingStatsByMatchType = `
	SELECT * FROM player_batting_stats pbs
	JOIN players p ON p.id = pbs.player_id
	WHERE p.public_id=$1 AND pbs.match_type=$2
`

func (q *Queries) GetCricketPlayerBowlingStatsByMatchType(context context.Context, playerPublicID uuid.UUID, matchType string) (*models.PlayerBowlingStats, error) {
	var i models.PlayerBowlingStats
	rows := q.db.QueryRowContext(context, getCricketPlayerBowlingStats, playerPublicID, matchType)
	err := rows.Scan(
		&i.ID,
		&i.PublicID,
		&i.PlayerID,
		&i.MatchType,
		&i.Matches,
		&i.Innings,
		&i.Wickets,
		&i.Runs,
		&i.Balls,
		&i.Average,
		&i.StrikeRate,
		&i.FourWickets,
		&i.FiveWickets,
		&i.CreatedAT,
		&i.UpdatedAT,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan the query: ", err)
	}
	return &i, nil
}

const getCricketPlayerBattingStats = `
	SELECT * FROM player_batting_stats
	WHERE player_id=$1
`

func (q *Queries) GetCricketPlayerBattingStats(context context.Context, playerID int64) (*[]models.PlayerBattingStats, error) {
	var playerStats []models.PlayerBattingStats
	rows, err := q.db.QueryContext(context, getCricketPlayerBattingStats, playerID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query: ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var i models.PlayerBattingStats
		err := rows.Scan(
			&i.ID,
			&i.PlayerID,
			&i.MatchType,
			&i.Matches,
			&i.Innings,
			&i.Runs,
			&i.Balls,
			&i.Sixes,
			&i.Fours,
			&i.Fifties,
			&i.Hundreds,
			&i.BestScore,
			&i.Average,
			&i.StrikeRate,
			&i.CreatedAT,
			&i.UpdatedAT,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("Failed to scan the query: ", err)
		}

		playerStats = append(playerStats, i)
	}
	return &playerStats, nil
}

const getCricketPlayerBowlingStats = `
	SELECT * FROM player_bowling_stats
	JOIN players p ON p.id = pbs.player_id
	WHERE p.public_id=$1
`

func (q *Queries) GetCricketPlayerBowlingStats(context context.Context, playerPublicID uuid.UUID) (*[]models.PlayerBowlingStats, error) {
	var playerStats []models.PlayerBowlingStats
	rows, err := q.db.QueryContext(context, getCricketPlayerBowlingStats, playerPublicID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query: ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var i models.PlayerBowlingStats
		err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.PlayerID,
			&i.MatchType,
			&i.Matches,
			&i.Innings,
			&i.Wickets,
			&i.Average,
			&i.StrikeRate,
			&i.FourWickets,
			&i.FiveWickets,
			&i.Runs,
			&i.CreatedAT,
			&i.UpdatedAT,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("Failed to scan the query: ", err)
		}

		playerStats = append(playerStats, i)
	}
	return &playerStats, nil
}

const updatePlayerBattingStatsQuery = `
	UPDATE player_cricket_stats pbs
	SET (
		matches = matches + 1,
		batting_innings = batting_innings + 1,
		runs = runs + $3,
		balls = balls + $4,
		sixes = sixes + $5,
		fours = fours + $6,
		fifties = fifties + (CASE $3 < 100 AND $3 >= 50 THEN 1 ELSE 0 END),
		hundreds = hudreds + (CASE $3 >= 100 THEN 1 ELSE 0 END),
		best_score = CASE $3 > best_score THEN $3 ELSE best_score END,
		Update_now
	)
	FROM cricket_stats cs
	JOIN players p ON p.id = cs.player_id
	WHERE p.public_id = $1 AND match_type = $2
	RETURNING *;
`

func (q *Queries) UpdatePlayerBattingStats(ctx context.Context, playerPublicID uuid.UUID, matchType string, runs, balls, fours, sixes, fifties, hundreds, bestScore int) (*models.CricketPlayerStats, error) {
	var i models.CricketPlayerStats
	rows := q.db.QueryRowContext(ctx, updatePlayerBattingStatsQuery, playerPublicID, matchType)
	err := rows.Scan(
		&i.ID,
		&i.PublicID,
		&i.PlayerID,
		&i.MatchType,
		&i.Matches,
		&i.BattingInnings,
		&i.Runs,
		&i.Balls,
		&i.Sixes,
		&i.Fours,
		&i.Fifties,
		&i.Hundreds,
		&i.BestScore,
		&i.BowlingInnings,
		&i.Wickets,
		&i.BallsBowled,
		&i.RunsConceded,
		&i.FourWickets,
		&i.FiveWickets,
		&i.CreatedAT,
		&i.UpdatedAT,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan the query data: ", err)
	}
	return &i, nil
}

const updatePlayerBowlingStatsQuery = `
	UPDATE player_cricket_stats
	SET
		bowling_innings = bowling_innings + 1,
		wickets = wickets + $3,
		runs_conceded = runs_conceded + $4,
		balls_bowled = balls_bowled + $5,
		four_wickets = four_wickets + (CASE WHEN $3 = 4 THEN 1 ELSE 0 END),
		five_wickets = five_wickets + (CASE WHEN $3 >= 5 THEN 1 ELSE 0 END),
		updated_at = NOW()
	FROM cricket_stats cs
	JOIN players p ON p.id = cs.player_id
	WHERE p.public_id = $1 AND match_type = $2;

`

func (q *Queries) UpdatePlayerBowlingStats(ctx context.Context, playerPublicID uuid.UUID, matchType string, wickets, runsConceded, ballsBowled int, fourWickets, fiveWickets int) (*models.CricketPlayerStats, error) {
	var i models.CricketPlayerStats
	rows := q.db.QueryRowContext(ctx, updatePlayerBowlingStatsQuery, playerPublicID, matchType)
	err := rows.Scan(
		&i.ID,
		&i.PublicID,
		&i.PlayerID,
		&i.MatchType,
		&i.Matches,
		&i.BattingInnings,
		&i.Runs,
		&i.Balls,
		&i.Sixes,
		&i.Fours,
		&i.Fifties,
		&i.Hundreds,
		&i.BestScore,
		&i.BowlingInnings,
		&i.Wickets,
		&i.BallsBowled,
		&i.RunsConceded,
		&i.FourWickets,
		&i.FiveWickets,
		&i.CreatedAT,
		&i.UpdatedAT,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan the query data: ", err)
	}
	return &i, nil
}
