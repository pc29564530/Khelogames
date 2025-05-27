package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"
)

const addPlayerBattingStats = `
	INSERT INTO player_batting_stats (
		player_id,
		match_type,
		total_matches,
		total_innings,
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
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING *;
`

func (q *Queries) AddPlayerBattingStats(ctx context.Context, playerID int32, matchType string, totalMatches, totalInnings, runs, balls, fours, sixes, fifties, hundreds, bestScore int, average, strikeRate string) (*models.PlayerBattingStats, error) {

	rows := q.db.QueryRowContext(ctx, addPlayerBattingStats, playerID, matchType, totalMatches, totalInnings, runs, balls, sixes, fours, fifties, hundreds, bestScore, average, strikeRate)

	var i models.PlayerBattingStats
	err := rows.Scan(
		&i.ID,
		&i.PlayerID,
		&i.MatchType,
		&i.TotalMatches,
		&i.TotalInnings,
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

const getCricketPlayerBattingStatsByMatchType = `
	SELECT * FROM player_batting_stats
	WHERE player_id=$1 AND match_type=$2
`

func (q *Queries) GetCricketPlayerBattingStatsByMatchType(context context.Context, playerID int64, matchType string) (*models.PlayerBattingStats, error) {
	var i models.PlayerBattingStats
	rows := q.db.QueryRowContext(context, getCricketPlayerBattingStats, playerID, matchType)
	err := rows.Scan(
		&i.ID,
		&i.PlayerID,
		&i.MatchType,
		&i.TotalMatches,
		&i.TotalInnings,
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
	return &i, nil
}

const getCricketPlayerBowlingStatsByMatchType = `
	SELECT * FROM player_batting_stats
	WHERE player_id=$1 AND match_type=$2
`

func (q *Queries) GetCricketPlayerBowlingStatsByMatchType(context context.Context, playerID int64, matchType string) (*models.PlayerBowlingStats, error) {
	var i models.PlayerBowlingStats
	rows := q.db.QueryRowContext(context, getCricketPlayerBowlingStats, playerID, matchType)
	err := rows.Scan(
		&i.ID,
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
			&i.TotalMatches,
			&i.TotalInnings,
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

const addPlayerBowlingStats = `
	INSERT INTO player_bowling_stats (
		player_id,
		match_type,
		matches,
		innings,
		wickets,
		runs,
		balls,
		four_wickets,
		five_wickets,
		average,
		strike_rate,
		economy_rate,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, $12)
	RETURNING *;
`

func (q *Queries) AddPlayerBowlingStats(ctx context.Context, playerID int32, matchType string, matches, innings, wickets, runs, balls int, average, strikeRate, economy_rate string, four_wickets, five_wickets int) (*models.PlayerBowlingStats, error) {

	rows := q.db.QueryRowContext(ctx, addPlayerBowlingStats, playerID, matchType, matches, innings, wickets, runs, balls, average, strikeRate, economy_rate, four_wickets, five_wickets)

	var i models.PlayerBowlingStats
	err := rows.Scan(
		&i.ID,
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
		return nil, fmt.Errorf("Failed to scan the query data: ", err)
	}
	return &i, nil
}

const getCricketPlayerBowlingStats = `
	SELECT * FROM player_bowling_stats
	WHERE player_id=$1
`

func (q *Queries) GetCricketPlayerBowlingStats(context context.Context, playerID int64) (*[]models.PlayerBowlingStats, error) {
	var playerStats []models.PlayerBowlingStats
	rows, err := q.db.QueryContext(context, getCricketPlayerBowlingStats, playerID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query: ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var i models.PlayerBowlingStats
		err := rows.Scan(
			&i.ID,
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
	UPDATE player_batting_stats pbs
	SET (
		total_matches = total_matches + 1,
		total_innings = total_innings + 1,
		runs = runs + $3,
		balls = balls + $4
		sixes = sixes + $5,
		fours = fours + $6,
		fifties = fifties + (CASE $3 < 100 AND $3 >= 50 THEN 1 ELSE 0 END),
		hundreds = hudreds + (CASE $3 >= 100 THEN 1 ELSE 0 END),
		best_score = CASE $3 > best_score THEN $3 ELSE best_score END,
		average = CASE
			WHEN (total_innings + 1) > 0
			THEN TO_CHAR((runs + $3)::DECIMAL / (total_innings + 1), 'FM999999.00')
			ELSE '0.00'
		END,
		strike_rate = CASE 
			WHEN (balls + $3) > 0 
			THEN TO_CHAR((runs + $3)::DECIMAL / (balls + $4), 'FM999999.00')
			ELSE '0.00'
		END,
		Update_now
	)
	WHERE player_id = $1 AND match_type = $2
	RETURNING *;
`

func (q *Queries) UpdatePlayerBattingStats(ctx context.Context, playerID int32, matchType string, totalMatches, totalInnings, runs, balls, fours, sixes, fifties, hundreds, bestScore int, average, strikeRate string) (*models.PlayerBattingStats, error) {
	var i models.PlayerBattingStats
	rows := q.db.QueryRowContext(ctx, updatePlayerBattingStatsQuery, playerID, matchType)
	err := rows.Scan(
		&i.ID,
		&i.PlayerID,
		&i.MatchType,
		&i.TotalMatches,
		&i.TotalInnings,
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
		return nil, fmt.Errorf("Failed to scan the query data: ", err)
	}
	return &i, nil
}

// update player bowler
const updatePlayerBowlingStatsQuery = `
	UPDATE player_bowling_stats
	SET
		matches = matches + 1,
		innings = innings + 1,
		wickets = wickets + $3,
		runs = runs + $4,
		balls = balls + $5,
		average = CASE 
			WHEN (wickets + $3) > 0 
			THEN TO_CHAR((runs + $4)::DECIMAL / (wickets + $3), 'FM999999.00')
			ELSE '0.00'
		END,
		strike_rate = CASE 
			WHEN (wickets + $3) > 0 
			THEN TO_CHAR((balls + $5)::DECIMAL / (wickets + $3), 'FM999999.00')
			ELSE '0.00'
		END,
		economy_rate = CASE 
			WHEN (balls + $5) > 0 
			THEN TO_CHAR((runs + $4)::DECIMAL / ((balls + $5)::DECIMAL / 6), 'FM999999.00')
			ELSE '0.00'
		END,
		four_wickets = four_wickets + (CASE WHEN $3 = 4 THEN 1 ELSE 0 END),
		five_wickets = five_wickets + (CASE WHEN $3 >= 5 THEN 1 ELSE 0 END),
		updated_at = NOW()
	WHERE player_id = $1 AND match_type = $2;

`

func (q *Queries) UpdatePlayerBowlingStats(ctx context.Context, playerID int32, matchType string, matches, innings, wickets, runs, balls int, average, strikeRate, economy_rate string, four_wickets, five_wickets int) (*models.PlayerBowlingStats, error) {
	var i models.PlayerBowlingStats
	rows := q.db.QueryRowContext(ctx, updatePlayerBowlingStatsQuery, playerID, matchType)
	err := rows.Scan(
		&i.ID,
		&i.PlayerID,
		&i.MatchType,
		&i.Matches,
		&i.Innings,
		&i.Wickets,
		&i.Runs,
		&i.Balls,
		&i.Average,
		&i.StrikeRate,
		&i.EconomyRate,
		&i.FourWickets,
		&i.FiveWickets,
		&i.CreatedAT,
		&i.UpdatedAT,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan the query data: ", err)
	}
	return &i, nil
}
