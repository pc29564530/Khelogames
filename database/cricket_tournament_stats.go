package database

import (
	"context"
	"database/sql"
	"fmt"
)

type PlayerStat struct {
	PlayerID   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	StatValue  string `json:"stat_value"`
}

func (q *Queries) getPlayerStat(ctx context.Context, query string, tournamentID int64) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, query, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query : ", err)
	}

	defer rows.Close()

	var playerStats []map[string]interface{}
	for rows.Next() {
		var stats PlayerStat
		err := rows.Scan(&stats.PlayerID, &stats.PlayerName, &stats.TeamName, &stats.StatValue)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan: ", err)
		}
		playerStats = append(playerStats, map[string]interface{}{
			"player_id":   stats.PlayerID,
			"player_name": stats.PlayerName,
			"team_name":   stats.TeamName,
			"stat_value":  stats.StatValue,
		})
	}

	return playerStats, nil
}

const cricketTournamentsMostRuns = `
	SELECT
		p.id AS player_id,
		p.player_name AS player_name,
		tm.name AS team_name,
		SUM(b.runs_scored) AS total_runs
	FROM bats b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	WHERE m.tournament_id = $1
	GROUP BY p.id, p.player_name, tm.name
	HAVING SUM(b.runs_scored) > 0
	ORDER BY total_runs DESC;
`

func (q *Queries) GetCricketTournamentMostRuns(ctx context.Context, tournamentID int64) ([]map[string]interface{}, error) {
	return q.getPlayerStat(ctx, cricketTournamentsMostRuns, tournamentID)
}

const getCricketTournamentHighestRuns = `
	SELECT 
		p.id AS player_id,
		p.player_name,
		tm.name AS team_name,
		MAX(b.runs_scored) AS highest_runs
	FROM bats b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	WHERE m.tournament_id = $1
	GROUP BY p.id, p.player_name, tm.name
	HAVING SUM (b.runs_scored) > 0
	ORDER BY highest_runs DESC;
`

func (q *Queries) GetCricketTournamentHighestRuns(ctx context.Context, tournamentID int64) ([]map[string]interface{}, error) {
	return q.getPlayerStat(ctx, getCricketTournamentHighestRuns, tournamentID)
}

const getCricketTournamentMostSixes = `
	SELECT 
		p.id AS player_id,
		p.player_name,
		tm.name AS team_name,
		sum(b.sixes) AS sixes
	FROM bats b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	WHERE m.tournament_id = $1
	GROUP BY p.id, p.player_name, tm.name
	HAVING SUM(b.sixes) > 0
	ORDER BY sixes DESC;
`

func (q *Queries) GetCricketTournamentMostSixes(ctx context.Context, tournamentID int64) ([]map[string]interface{}, error) {
	return q.getPlayerStat(ctx, getCricketTournamentMostSixes, tournamentID)
}

const getCricketTournamentMostFours = `
	SELECT 
		p.id AS player_id,
		p.player_name,
		tm.name AS team_name,
		sum(b.fours) AS fours
	FROM bats b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	WHERE m.tournament_id = $1
	GROUP BY p.id, p.player_name, tm.name
	HAVING SUM(b.fours) > 0
	ORDER BY fours DESC;
`

func (q *Queries) GetCricketTournamentMostFours(ctx context.Context, tournamentID int64) ([]map[string]interface{}, error) {
	return q.getPlayerStat(ctx, getCricketTournamentMostFours, tournamentID)
}

const getCricketTournamentMostFifties = `
	SELECT 
		p.id AS player_id,
		p.player_name,
		tm.name AS team_name,
		COUNT(*) FILTER(WHERE b.runs_scored >= 50 AND b.runs_scored < 100) AS fifties
	FROM bats b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	WHERE m.tournament_id = $1
	GROUP BY p.id, p.player_name, tm.name
	HAVING COUNT(*) FILTER(WHERE b.runs_scored >= 50 AND b.runs_scored < 100) > 0
	ORDER BY fifties DESC;
`

func (q *Queries) GetCricketTournamentMostFifties(ctx context.Context, tournamentID int64) ([]map[string]interface{}, error) {
	return q.getPlayerStat(ctx, getCricketTournamentMostFifties, tournamentID)
}

const getCricketTournamentMostHundreds = `
	SELECT 
		p.id AS player_id,
		p.player_name,
		tm.name AS team_name,
		COUNT(*) FILTER(WHERE b.runs_scored >= 100) AS hundreds
	FROM bats b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	WHERE m.tournament_id = $1
	GROUP BY p.id, p.player_name, tm.name
	HAVING COUNT(*) FILTER(WHERE b.runs_scored >= 100) > 0
	ORDER BY hundreds DESC;
`

func (q *Queries) GetCricketTournamentMostHundreds(ctx context.Context, tournamentID int64) ([]map[string]interface{}, error) {
	return q.getPlayerStat(ctx, getCricketTournamentMostHundreds, tournamentID)
}
