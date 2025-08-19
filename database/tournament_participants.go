package database

import (
	"context"
	"encoding/json"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addTournamentParticipantsQuery = `
WITH tournament_resolved AS (
    SELECT id AS tournament_id
    FROM tournaments
    WHERE public_id = $1
),
entity_resolved AS (
    SELECT id AS entity_id
    FROM teams
    WHERE public_id = $3 AND $4 = 'team'
    UNION
    SELECT id AS entity_id
    FROM players
    WHERE public_id = $3 AND $4 = 'player'
)
INSERT INTO tournament_participants (
    tournament_id,
    group_id,
    entity_id,
    entity_type,
    seed_number,
    status
)
SELECT t.tournament_id, $2, e.entity_id, $4, $5, $6
FROM tournament_resolved t
JOIN entity_resolved e ON TRUE
RETURNING *;
`

func (q *Queries) AddTournamentParticipants(ctx context.Context, tournamentPublicID uuid.UUID, groupID int32, entityPublicID uuid.UUID, entityType string, seedNumber int, status string) (*models.TournamentParticipants, error) {
	row := q.db.QueryRowContext(ctx, addTournamentParticipantsQuery, tournamentPublicID, groupID, entityPublicID, entityType, seedNumber, status)
	var i models.TournamentParticipants
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.TournamentID,
		&i.GroupID,
		&i.EntityID,
		&i.EntityType,
		&i.SeedNumber,
		&i.Status,
	)

	if err != nil {

	}
	return &i, err
}

const getTournamentParticipantsQuery = `
	SELECT JSON_BUILD_OBJECT(
		'id', tp.id, 'public_id', tp.public_id 'tournament_id', tp.tournament_id, 'group_id', tp.group_id, 'entity_id', tp.entity_id, 'entity_type', tp.entity_type, 'seed_number', tp.seed_number, 'status', tp.status,
		'entity', CASE 
			WHEN tp.entity_type = 'team'
			THEN JSON_BUILD_OBJECT(
				'id', tm.id, 'public_id', tm.public_id, 'user_id', tm.user_id, 'name', tm.name, 'slug', tm.slug, 'short_name', tm.shortname, 'admin', tm.admin, 'media_url', tm.media_url, 'gender', tm.gender, 'national', tm.national, 'country', tm.country, 'type', tm.type, 'player_count', tm.player_count, 'game_id', tm.game_id
			) ELSE NULL
			WHEN tp.entity_type = 'player'
			THEN JSON_BUILD_OBJECT (
				p.id, p.public_id, p.user_id, p.game_id, p.name, p.slug, p.short_name, p.media_url, p.positions, p.country, p.created_at, p.updated_at
			) ELSE NULL
		END
	) FROM tournament_participants tp
	JOIN tournaments AS t ON t.id = tp.tournament_id
	JOIN teams AS tm ON tm.id = tp.entity_id
	JOIN players AS p ON p.id = tp.entity_id
	WHERE t.public_id = $1
`

func (q *Queries) GetTournamentParticipants(ctx context.Context, tournamentPublicID uuid.UUID, entityType string) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentParticipantsQuery, tournamentPublicID, entityType)
	if err != nil {
		return nil, err
	}
	var tournamentParticipants []map[string]interface{}
	for rows.Next() {
		var participants map[string]interface{}
		var jsonByte []byte
		err := rows.Scan(&jsonByte)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonByte, &participants)
		if err != nil {
			return nil, err
		}

		tournamentParticipants = append(tournamentParticipants, participants)
	}
	return tournamentParticipants, nil
}
