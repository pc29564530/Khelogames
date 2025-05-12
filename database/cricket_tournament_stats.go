package database

import (
	"context"
	"database/sql"
	"fmt"
)

const cricketTournamentsStats = `
	SELECT
		p.id AS player_id,
		p.player_name AS player_name,
		tm.name AS team_name,
		SUM(b.runs_scored) AS total_runs
	FROM bats b
	JOIN matches m ON m.id = b.match_id
	JOIN players p ON p.id = b.batsman_id
	JOIN teams tm ON tm.id = b.team_id
	WHERE m.tournament_id = $1
	GROUP BY p.id, p.player_name, tm.name
	ORDER BY total_runs DESC;
`

type mostRun struct {
	PlayerID   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	TotalRuns  string `json:"total_runs"`
}

func (q *Queries) GetCricketTournamentMostRuns(ctx context.Context, tournamentID int64) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, cricketTournamentsStats, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query : ", err)
	}

	defer rows.Close()

	var mostRuns []map[string]interface{}
	for rows.Next() {
		var stats mostRun
		err := rows.Scan(&stats.PlayerID, &stats.PlayerName, &stats.TeamName, &stats.TotalRuns)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan: ", err)
		}
		fmt.Println("stats Line no 49: ", stats)
		mostRuns = append(mostRuns, map[string]interface{}{
			"player_id":   stats.PlayerID,
			"player_name": stats.PlayerName,
			"team_name":   stats.TeamName,
			"total_runs":  stats.TotalRuns,
		})
	}

	return mostRuns, nil

}

const getCricketTournamentHighestRuns = `
	SELECT 
		p.id AS player_id,
		p.player_name,
		tm.name AS team_name,
		MAX(b.runs_scored) AS highest_runs
	FROM bats b
	JOIN matches m ON m.id = b.match_id
	JOIN players p ON p.id = b.batsman_id
	JOIN teams tm ON tm.id = b.team_id
	WHERE m.tournament_id = $1
	GROUP BY p.id, p.player_name, tm.name
	ORDER BY highest_runs DESC;
`

type HighestRunsPlayer struct {
	PlayerID    int64  `json:"player_id"`
	PlayerName  string `json:"player_name"`
	TeamName    string `json:"team_name"`
	HighestRuns string `json:"highest_runs"`
}

func (q *Queries) GetCricketTournamentHighestRuns(ctx context.Context, tournamentID int64) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getCricketTournamentHighestRuns, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query : ", err)
	}

	defer rows.Close()

	var HighestRunsPlayers []map[string]interface{}
	for rows.Next() {
		var stats HighestRunsPlayer
		err := rows.Scan(&stats.PlayerID, &stats.PlayerName, &stats.TeamName, &stats.HighestRuns)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan: ", err)
		}
		fmt.Println("stats Line no 49: ", stats)
		HighestRunsPlayers = append(HighestRunsPlayers, map[string]interface{}{
			"player_id":    stats.PlayerID,
			"player_name":  stats.PlayerName,
			"team_name":    stats.TeamName,
			"highest_runs": stats.HighestRuns,
		})
	}

	return HighestRunsPlayers, nil

}

const getCricketTournamentMostSixes = `
	SELECT 
		p.id AS player_id,
		p.player_name,
		tm.name AS team_name,
		sum(b.sixes) AS sixes
	FROM bats b
	JOIN matches m ON m.id = b.match_id
	JOIN players p ON p.id = b.batsman_id
	JOIN teams tm ON tm.id = b.team_id
	WHERE m.tournament_id = $1
	GROUP BY p.id, p.player_name, tm.name
	ORDER BY sixes DESC;
`

type MostSixesByPlayer struct {
	PlayerID   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	MostSixes  string `json:"most_sixes"`
}

func (q *Queries) GetCricketTournamentMostSixes(ctx context.Context, tournamentID int64) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getCricketTournamentMostSixes, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query : ", err)
	}

	defer rows.Close()

	var MostSixesByPlayers []map[string]interface{}
	for rows.Next() {
		var stats MostSixesByPlayer
		err := rows.Scan(&stats.PlayerID, &stats.PlayerName, &stats.TeamName, &stats.MostSixes)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan: ", err)
		}
		fmt.Println("stats Line no 49: ", stats)
		MostSixesByPlayers = append(MostSixesByPlayers, map[string]interface{}{
			"player_id":   stats.PlayerID,
			"player_name": stats.PlayerName,
			"team_name":   stats.TeamName,
			"most_sixes":  stats.MostSixes,
		})
	}

	return MostSixesByPlayers, nil

}
