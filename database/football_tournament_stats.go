package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type FootballPlayerStat struct {
	PlayerID   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	StatValue  int    `json:"stat_value"`
}

func (q *Queries) GetFootballFootballPlayerStat(ctx context.Context, query string, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {

	rows, err := q.db.QueryContext(ctx, query, tournamentPublicID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query : ", err)
	}
	defer rows.Close()

	var FootballPlayerStats []map[string]interface{}
	for rows.Next() {
		var stats FootballPlayerStat
		err := rows.Scan(&stats.PlayerID, &stats.PlayerName, &stats.TeamName, &stats.StatValue)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Error: ", err)
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan: ", err)
		}

		FootballPlayerStats = append(FootballPlayerStats, map[string]interface{}{
			"player_id":   stats.PlayerID,
			"player_name": stats.PlayerName,
			"team_name":   stats.TeamName,
			"stat_value":  stats.StatValue,
		})
	}

	return FootballPlayerStats, nil
}

// player goals
const getFootballTournamentGoals = `
SELECT 
	p.id AS player_id,
	p.name,
	tm.name AS team_name,
	COUNT(*) FILTER(WHERE fi.incident_type = 'goal') AS goals
FROM football_incidents fi
LEFT JOIN matches AS m ON fi.match_id = m.id
INNER JOIN football_incident_player AS fip ON fip.incident_id = fi.id
INNER JOIN players p ON p.id = fip.player_id
LEFT JOIN teams tm ON tm.id = fi.team_id
LEFT JOIN tournaments t ON t.id = m.tournament_id
WHERE t.public_id = $1 AND p.id IS NOT NULL
GROUP BY p.id, p.name, tm.name
HAVING COUNT(*) FILTER(WHERE fi.incident_type = 'goal') > 0
ORDER BY goals DESC;
`

func (q *Queries) GetFootballTournamentPlayersGoals(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetFootballFootballPlayerStat(ctx, getFootballTournamentGoals, tournamentPublicID)
}

// player yellow cards
const getFootballTournamentYellowCards = `
	SELECT 
		p.id AS player_id,
		p.name,
		tm.name AS team_name,
		COUNT(*) FILTER(WHERE fi.incident_type = 'yellow_card') AS yellow_cards
	FROM football_incidents fi
	LEFT JOIN matches AS m ON fi.match_id = m.id
	JOIN football_incident_player AS fip ON fip.incident_id = fi.id
	JOIN players p ON p.id = fip.player_id
	LEFT JOIN teams tm ON tm.id = fi.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1 AND p.id IS NOT NULL
	GROUP BY p.id, p.name, tm.name
	HAVING COUNT(*) FILTER(WHERE fi.incident_type = 'yellow_card') > 0
	ORDER BY yellow_cards DESC;
`

func (q *Queries) GetFootballTournamentPlayersYellowCard(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetFootballFootballPlayerStat(ctx, getFootballTournamentYellowCards, tournamentPublicID)
}

// player yellow cards
const getFootballTournamentRedCards = `
	SELECT 
		p.id AS player_id,
		p.name,
		tm.name AS team_name,
		COUNT(*) FILTER(WHERE fi.incident_type = 'red_cards') AS red_cards
	FROM matches m
	LEFT JOIN football_incidents AS fi ON fi.match_id = m.id
	JOIN football_incident_player AS fip ON fip.incident_id = fi.id
	JOIN players p ON p.id = fip.player_id
	LEFT JOIN teams tm ON tm.id = fi.team_id
	LEFT JOIN tournaments t ON t.id = m.tournament_id
	WHERE t.public_id = $1 AND p.id IS NOT NULL
	GROUP BY p.id, p.name, tm.name
	HAVING COUNT(*) FILTER(WHERE fi.incident_type = 'red_cards') > 0
	ORDER BY red_cards DESC;
`

func (q *Queries) GetFootballTournamentPlayersRedCard(ctx context.Context, tournamentPublicID uuid.UUID) ([]map[string]interface{}, error) {
	return q.GetFootballFootballPlayerStat(ctx, getFootballTournamentRedCards, tournamentPublicID)
}
