package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"
)

const addPlayerCricketStats = `
	player_id,
	match_type,
	matches,
	batting_innings,
	batting_runs,
	balls_faced
	sixes,
	fours,
	fifties,
	hundreds,
	bowling_innings,
	wickets,
	runs_conceded,
	balls_bowled,
	four_wickets,
	five_wickets,
	created_at,
	updated_at
`

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
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING *;
`

func (q *Queries) AddPlayerCricketStats(ctx context.Context, playerID int32, matchType string, matches, battingInnings int32, runs, ballsFaced, fours, sixes, fifties, hundreds, bestScore int, bowlingInnings int32, wickets, runsConceded, ballsBowled, fourWickets, fiveWickets int) (*models.PlayerCricketStats, error) {

	rows := q.db.QueryRowContext(ctx, addPlayerBattingStats, playerID, matchType, matches, battingInnings, runs, ballsFaced, fours, sixes, fifties, hundreds, bestScore, bowlingInnings, wickets, runsConceded, ballsBowled, fourWickets, fiveWickets)

	var i models.PlayerCricketStats
	err := rows.Scan(
		&i.ID,
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

const getPlayerCricketStatsByMatchType = `
	SELECT * FROM player_cricket_stats
	WHERE player_id=$1
`

func (q *Queries) GetPlayerCricketStatsByMatchType(context context.Context, playerID int64) (*[]models.PlayerCricketStats, error) {
	rows, err := q.db.QueryContext(context, getCricketPlayerBattingStats, playerID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query database: ", err)
	}
	defer rows.Close()

	var playerStats []models.PlayerCricketStats

	for rows.Next() {
		var i models.PlayerCricketStats
		err := rows.Scan(
			&i.ID,
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
	WHERE player_id = $1 AND match_type = $2
	RETURNING *;
`

func (q *Queries) UpdatePlayerBattingStats(ctx context.Context, playerID int32, matchType string, runs, balls, fours, sixes, fifties, hundreds, bestScore int) (*models.PlayerCricketStats, error) {
	var i models.PlayerCricketStats
	rows := q.db.QueryRowContext(ctx, updatePlayerBattingStatsQuery, playerID, matchType)
	err := rows.Scan(
		&i.ID,
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
	WHERE player_id = $1 AND match_type = $2;

`

func (q *Queries) UpdatePlayerBowlingStats(ctx context.Context, playerID int32, matchType string, wickets, runsConceded, ballsBowled int, fourWickets, fiveWickets int) (*models.PlayerCricketStats, error) {
	var i models.PlayerCricketStats
	rows := q.db.QueryRowContext(ctx, updatePlayerBowlingStatsQuery, playerID, matchType)
	err := rows.Scan(
		&i.ID,
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
		return nil, fmt.Errorf("Failed to scan the query data: ", err)
	}
	return &i, nil
}
