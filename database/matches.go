package database

import (
	"context"
	"encoding/json"
	"log"
)

const getAllMatches = `
SELECT 
    json_build_object(
        'id', m.id,
        'tournament_id', m.tournament_id,
        'away_team_id', m.away_team_id,
        'home_team_id', m.home_team_id,
        'start_timestamp', m.start_timestamp,
        'end_timestamp', m.end_timestamp,
        'type', m.type,
        'status_code', m.status_code,
        'result', m.result,
        'stage', m.stage,
        'homeTeam', json_build_object(
            'id', ht.id,
            'name', ht.name,
            'slug', ht.slug,
            'short_name', ht.shortname,
            'admin', ht.admin,
            'media_url', ht.media_url,
            'gender', ht.gender,
            'national', ht.national,
            'country', ht.country,
            'type', ht.type,
            'player_count', ht.player_count,
            'game_id', ht.game_id
        ),
        'homeScore', CASE 
            WHEN g.name = 'football' AND fs.id IS NOT NULL AND fs.team_id=m.home_team_id THEN json_build_object(
                'id', fs.id,
                'match_id', fs.match_id,
                'team_id', fs.team_id,
                'first_half', fs.first_half,
                'second_half', fs.second_half,
                'goals', fs.goals
            )
            WHEN g.name = 'cricket' AND cs.id IS NOT NULL AND fs.team_id=m.home_team_id THEN json_build_object(
                'id', cs.id,
                'match_id', cs.match_id,
                'team_id', cs.team_id,
                'inning', cs.inning,
                'score', cs.score,
                'wickets', cs.wickets,
                'overs', cs.overs,
                'run_rate', cs.run_rate,
                'target_run_rate', cs.target_run_rate
            )
            ELSE NULL
        END,
        'awayTeam', json_build_object(
            'id', at.id,
            'name', at.name,
            'slug', at.slug,
            'short_name', at.shortname,
            'admin', at.admin,
            'media_url', at.media_url,
            'gender', at.gender,
            'national', at.national,
            'country', at.country,
            'type', at.type,
            'player_count', at.player_count,
            'game_id', at.game_id
        ),
        'awayScore', CASE 
            WHEN g.name = 'football' AND fs.id IS NOT NULL AND cs.team_id=m.away_team_id THEN json_build_object(
                'id', fs.id,
                'match_id', fs.match_id,
                'team_id', fs.team_id,
                'first_half', fs.first_half,
                'second_half', fs.second_half,
                'goals', fs.goals
            )
            WHEN g.name = 'cricket' AND cs.id IS NOT NULL AND cs.team_id=m.away_team_id THEN json_build_object(
                'id', cs.id,
                'match_id', cs.match_id,
                'team_id', cs.team_id,
                'inning', cs.inning,
                'score', cs.score,
                'wickets', cs.wickets,
                'overs', cs.overs,
                'run_rate', cs.run_rate,
                'target_run_rate', cs.target_run_rate
            )
            ELSE NULL
        END,
        'tournament', json_build_object(
            'id', t.id,
            'name', t.name,
            'slug', t.slug,
            'country', t.country,
			'status_code', t.status_code,
            'level', t.level,
            'start_timestamp', t.start_timestamp,
            'game_id', t.game_id,
            'group_count', t.group_count,
            'max_group_team', t.max_group_team,
            'stage', t.stage
        )
    ) AS response
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
LEFT JOIN tournaments AS t ON m.tournament_id = t.id
JOIN games g ON t.game_id = g.id
LEFT JOIN football_score fs ON m.id = fs.match_id AND g.name = 'football'
LEFT JOIN cricket_score cs ON m.id = cs.match_id AND g.name = 'cricket'
WHERE m.start_timestamp >= $1
  AND t.game_id = $2
ORDER BY m.start_timestamp;

`

type MatchesResponse struct {
	Response interface{} `json:"response"`
}

func (q *Queries) GetAllMatches(ctx context.Context, date int32, gameID int64) ([]map[string]interface{}, error) {
	var matches []map[string]interface{}

	rows, err := q.db.QueryContext(ctx, getAllMatches, date, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var matchResponse MatchesResponse
		if err := rows.Scan(&matchResponse.Response); err != nil {
			log.Printf("Failed to scan match: %v", err)
			continue
		}

		var matchResult map[string]interface{}
		data := matchResponse.Response.([]byte)
		if err := json.Unmarshal(data, &matchResult); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		matches = append(matches, matchResult)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}