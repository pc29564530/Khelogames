package database

import (
	"context"
	"encoding/json"
	"log"
)

//Get Top Performer

const getFootballTopPerformer = `
SELECT 
	json_build_object(
		'player', json_build_object(
			'public_id', p.public_id,
			'name', p.name
		),
		'goals', COUNT(fi.id)
	)
FROM matches m
JOIN football_incidents fi ON fi.match_id = m.id  AND fi.incident_type = 'goals'
JOIN football_incident_player fip ON fip.incident_id = fi.id
JOIN players p 
	ON p.id = fip.player_id
WHERE m.game_id = 1 AND m.status_code = 'finished'
	AND m.start_timestamp BETWEEN 
		(EXTRACT(EPOCH FROM NOW()) * 1000) - (7 * 24 * 60 * 60 * 1000)
		AND (EXTRACT(EPOCH FROM NOW()) * 1000)
GROUP BY p.id, p.public_id, p.name
ORDER BY COUNT(fi.id) DESC
LIMIT 5;
`

func (q *Queries) GetFootballTopPerformer(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getFootballTopPerformer)
	if err != nil {
		log.Printf("Failed to query: %v", err)
		return nil, err
	}

	defer rows.Close()

	// initialize slice so empty result returns []
	topPerformers := make([]map[string]interface{}, 0)
	for rows.Next() {
		var topPerformer map[string]interface{}
		var jsonByte []byte
		err := rows.Scan(&jsonByte)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		err = json.Unmarshal(jsonByte, &topPerformer)
		if err != nil {
			log.Printf("Failed to unmarshal: %v", err)
			return nil, err
		}
		topPerformers = append(topPerformers, topPerformer)
	}
	//check row iteration error
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}
	return topPerformers, nil
}

const getCricketBatsmanTopPerformer = `
	SELECT 
		json_build_object(
			'player', json_build_object(
				'public_id', p.public_id,
				'name', p.name
			),
			'total_runs', SUM(bs.runs_scored)
		)
	FROM matches m
	JOIN batsman_score bs
		ON bs.match_id = m.id
	JOIN players p 
		ON p.id = bs.batsman_id
	WHERE m.game_id = 2 AND m.status_code = 'finished'
		AND m.start_timestamp BETWEEN 
			(EXTRACT(EPOCH FROM NOW()) * 1000) - (7 * 24 * 60 * 60 * 1000)
			AND (EXTRACT(EPOCH FROM NOW()) * 1000)
	GROUP BY p.id, p.public_id, p.name
	ORDER BY COUNT(bs.id) DESC
	LIMIT 5;
`

func (q *Queries) GetCricketTopBattingPerformer(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getCricketBatsmanTopPerformer)
	if err != nil {
		log.Printf("Failed to query: %v", err)
		return nil, err
	}

	defer rows.Close()

	var topPerformers []map[string]interface{}
	for rows.Next() {
		var topPerformer map[string]interface{}
		var jsonByte []byte
		err := rows.Scan(&jsonByte)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		err = json.Unmarshal(jsonByte, &topPerformer)
		if err != nil {
			log.Printf("Failed to unmarshal: %v", err)
			return nil, err
		}
		topPerformers = append(topPerformers, topPerformer)
	}
	//check row iteration error
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}
	return topPerformers, nil
}

const getCricketBowlingTopPerformer = `
	SELECT 
		json_build_object(
			'player', json_build_object(
				'public_id', p.public_id,
				'name', p.name
			),
			'wickets', SUM(bs.wickets)
		)
	FROM matches m
	JOIN bolwer_score bs
		ON bs.match_id = m.id
	JOIN players p 
		ON p.id = bs.bowler_id
	WHERE m.game_id = 2 AND m.status_code = 'finished'
		AND m.start_timestamp BETWEEN 
			(EXTRACT(EPOCH FROM NOW()) * 1000) - (7 * 24 * 60 * 60 * 1000)
			AND (EXTRACT(EPOCH FROM NOW()) * 1000)
	GROUP BY p.id, p.public_id, p.name
	ORDER BY COUNT(bs.id) DESC
	LIMIT 5;
`

func (q *Queries) GetCricketTopBowlingPerformer(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getCricketBatsmanTopPerformer)
	if err != nil {
		log.Printf("Failed to query: %v", err)
		return nil, err
	}

	defer rows.Close()

	var topPerformers []map[string]interface{}
	for rows.Next() {
		var topPerformer map[string]interface{}
		var jsonByte []byte
		err := rows.Scan(&jsonByte)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		err = json.Unmarshal(jsonByte, &topPerformer)
		if err != nil {
			log.Printf("Failed to unmarshal: %v", err)
			return nil, err
		}
		topPerformers = append(topPerformers, topPerformer)
	}
	//check row iteration error
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}
	return topPerformers, nil
}
