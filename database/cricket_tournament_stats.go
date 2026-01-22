package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type PlayerStat struct {
	PlayerID   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	StatValue  string `json:"stat_value"`
}

func (q *Queries) GetPlayerStat(ctx context.Context, query string, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, query, tournamentPublicID)
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
		p.name AS player_name,
		tm.name AS team_name,
		SUM(b.runs_scored) AS total_runs
	FROM batsman_score b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	HAVING SUM(b.runs_scored) > 0
	ORDER BY total_runs DESC;
`

func (q *Queries) GetCricketTournamentMostRuns(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, cricketTournamentsMostRuns, tournamentPublicID)
}

const getCricketTournamentHighestRuns = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		MAX(b.runs_scored) AS highest_runs
	FROM batsman_score b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	HAVING SUM (b.runs_scored) > 0
	ORDER BY highest_runs DESC;
`

func (q *Queries) GetCricketTournamentHighestRuns(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentHighestRuns, tournamentPublicID)
}

const getCricketTournamentMostSixes = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		sum(b.sixes) AS sixes
	FROM batsman_score b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	HAVING SUM(b.sixes) > 0
	ORDER BY sixes DESC;
`

func (q *Queries) GetCricketTournamentMostSixes(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentMostSixes, tournamentPublicID)
}

const getCricketTournamentMostFours = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		sum(b.fours) AS fours
	FROM batsman_score b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	HAVING SUM(b.fours) > 0
	ORDER BY fours DESC;
`

func (q *Queries) GetCricketTournamentMostFours(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentMostFours, tournamentPublicID)
}

const getCricketTournamentMostFifties = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		COUNT(*) FILTER(WHERE b.runs_scored >= 50 AND b.runs_scored < 100) AS fifties
	FROM batsman_score b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	HAVING COUNT(*) FILTER(WHERE b.runs_scored >= 50 AND b.runs_scored < 100) > 0
	ORDER BY fifties DESC;
`

func (q *Queries) GetCricketTournamentMostFifties(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentMostFifties, tournamentPublicID)
}

const getCricketTournamentMostHundreds = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		COUNT(*) FILTER(WHERE b.runs_scored >= 100) AS hundreds
	FROM batsman_score b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	HAVING COUNT(*) FILTER(WHERE b.runs_scored >= 100) > 0
	ORDER BY hundreds DESC;
`

func (q *Queries) GetCricketTournamentMostHundreds(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentMostHundreds, tournamentPublicID)
}

//Bowling Parts
//Most Wickets

const getCricketTournamentMostWickets = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		SUM(b.wickets) AS wickets
	FROM balls b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.bowler_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	HAVING SUM(b.wickets) > 0
	ORDER BY wickets DESC;
`

func (q *Queries) GetCricketTournamentMostWickets(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentMostWickets, tournamentPublicID)
}

// EconomyRate
const getCricketTournamentEconomyRate = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		ROUND(sum(b.runs)::numeric/sum(ball-wide-no_ball) * 6, 2) AS economy_rate
	FROM balls b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.bowler_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	ORDER BY economy_rate DESC;
`

func (q *Queries) GetCricketTournamentBowlingEconomyRate(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentEconomyRate, tournamentPublicID)
}

// Bowling Average
const getCricketTournamentBowlingAverage = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		ROUND(SUM(b.runs)::numeric / NULLIF(SUM(b.wickets), 0), 2) AS bowling_average
	FROM balls b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.bowler_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	HAVING SUM(b.wickets) > 0
	ORDER BY bowling_average DESC;
`

func (q *Queries) GetCricketTournamentBowlingAverage(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentBowlingAverage, tournamentPublicID)
}

// Bowling Strike Rate
const getCricketTournamentBowlingStrikeRate = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		ROUND(sum(b.runs)::numeric / sum(ball) - sum(wide) - sum(no_ball), 2) AS bowling_strike_rate
	FROM balls b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.bowler_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	ORDER BY bowling_strike_rate DESC;
`

func (q *Queries) GetCricketTournamentBowlingStrikeRate(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentBowlingStrikeRate, tournamentPublicID)
}

// Bowling Strike Rate
const getCricketTournamentBowlingFiveWicketHaul = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		COUNT(*) FILTER(WHERE b.wickets >= 0) AS five_wickets_haul
	FROM balls b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.bowler_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	ORDER BY five_wickets_haul DESC;
`

func (q *Queries) GetCricketTournamentBowlingFiveWicketHaul(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentBowlingFiveWicketHaul, tournamentPublicID)
}

// batting average
const getCricketTournamentBattingAverage = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		ROUND(
			SUM(b.runs_scored)::numeric / NULLIF(
				COUNT(*) FILTER (
					WHERE cs.is_inning_completed = true AND b.is_currently_batting = false
				), 0
			), 2
		) AS batting_average
	FROM batsman_score b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	JOIN cricket_score cs ON cs.match_id = m.id AND cs.team_id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	HAVING
		SUM(b.runs_scored)::numeric / NULLIF(
			COUNT(*) FILTER (
				WHERE cs.is_inning_completed = true AND b.is_currently_batting = false
			), 0
		) > 0
	ORDER BY batting_average DESC;
`

func (q *Queries) GetCricketTournamentBattingAverage(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentBowlingAverage, tournamentPublicID)
}

// Batting Strike Rate
const getCricketTournamentBattingStrikeRate = `
	SELECT 
		p.id AS player_id,
		p.name AS player_name,
		tm.name AS team_name,
		ROUND(sum(b.runs_scored)::numeric / sum(balls_faced), 2) AS batting_strike_rate
	FROM batsman_score b
	LEFT JOIN matches m ON m.id = b.match_id
	LEFT JOIN players p ON p.id = b.batsman_id
	LEFT JOIN teams tm ON tm.id = b.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1
	GROUP BY p.id, p.name, tm.name
	ORDER BY batting_strike_rate DESC;
`

func (q *Queries) GetCricketTournamentBattingStrikeRate(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetPlayerStat(ctx, getCricketTournamentBattingStrikeRate, tournamentPublicID)
}
