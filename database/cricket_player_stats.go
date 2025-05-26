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
		sixes,
		fours,
		fifties,
		hundreds,
		best_score,
		average,
		strike_rate,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING *;
`

func (q *Queries) AddPlayerBattingStats(ctx context.Context, playerID int32, matchType string, totalMatches, totalInnings, runs, sixes, fours, fifties, hundreds, bestScore int, average, strikeRate string) (*models.PlayerBattingStats, error) {

	rows := q.db.QueryRowContext(ctx, addPlayerBattingStats, playerID, matchType, totalMatches, totalInnings, runs, sixes, fours, fifties, hundreds, bestScore, average, strikeRate)

	var i models.PlayerBattingStats
	err := rows.Scan(
		&i.ID,
		&i.PlayerID,
		&i.MatchType,
		&i.TotalMatches,
		&i.TotalInnings,
		&i.Runs,
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
		four_wickets,
		five_wickets,
		economy_rate,
		average,
		strike_rate,
		runs,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, $12)
	RETURNING *;
`

func (q *Queries) AddPlayerBowlingStats(ctx context.Context, playerID int32, matchType string, matches, innings, wickets int, average, strikeRate, economy_rate string, four_wickets, five_wickets, runs int) (*models.PlayerBowlingStats, error) {

	rows := q.db.QueryRowContext(ctx, addPlayerBowlingStats, playerID, matchType, matches, innings, wickets, average, strikeRate, economy_rate, four_wickets, five_wickets, runs)

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
