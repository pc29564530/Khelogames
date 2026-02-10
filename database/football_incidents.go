package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addFootballSubsPlayer = `
WITH incidentID AS (
	SELECT * FROM football_incidents WHERE public_id = $1
),
playerInID AS (
	SELECT * FROM players WHERE public_id = $2
),
playerOutID AS (
	SELECT * FROM players WHERE public_id = $3
)
INSERT INTO football_substitutions_player (
    incident_id,
    player_in_id,
    player_out_id
)
SELECT
	JSON_BUILD_OBJECT (
		'incident_id': incidentID.id,
		'player_in', JSON_BUILD_OBJECT(
			'id', playerInID.id,
			'public_id', playerInID.public_id,
			'user_id', playerInID.user_id,
			'game_id', playerInID.game_id,
			'name', playerInID.name,
			'slug', playerInID.slug,
			'short_name', playerInID.ShortName,
			'media_url', playerInID.media_url,
			'positions', playerInID.positions,
			'country', playerInID.country,
			'created_at', playerInID.created_at,
			'updated_at', playerInID.updated_at
		),
		'player_out', JSON_BUILD_OBJECT(
				'id', playerOutID.id,
				'public_id', playerOutID.public_id,
				'user_id', playerOutID.user_id,
				'game_id', playerOutID.game_id,
				'name', playerOutID.name,
				'slug', playerOutID.slug,
				'short_name', playerOutID.ShortName,
				'media_url', playerOutID.media_url,
				'positions', playerOutID.positions,
				'country', playerOutID.country,
				'created_at', playerOutID.created_at,
				'updated_at', playerOutID.updated_at
		
		)
    )
FROM incidentID, playerInID, playerOutID	
RETURNING *;
`

func (q *Queries) ADDFootballSubsPlayer(ctx context.Context, incidentPublicID, playerInPublicID, playerOutPublicID uuid.UUID) (*map[string]interface{}, error) {
	row := q.db.QueryRowContext(ctx, addFootballSubsPlayer, incidentPublicID, playerInPublicID, playerOutPublicID)
	var i map[string]interface{}
	var jsonByte []byte

	err := row.Scan(&jsonByte)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: ", err)
	}

	err = json.Unmarshal(jsonByte, &i)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal: ", err)
	}
	return &i, err
}

const addFootballIncidentPlayer = `
WITH playerID AS (
	SELECT id FROM players WHERE public_id = $2
)
INSERT INTO football_incident_player (
    incident_id,
    player_id
)
SELECT $1, playerID.id FROM playerID
RETURNING *;
`

func (q *Queries) AddFootballIncidentPlayer(ctx context.Context, incidentID int64, playerPublicID uuid.UUID) (models.FootballIncidentPlayer, error) {
	row := q.db.QueryRowContext(ctx, addFootballIncidentPlayer, incidentID, playerPublicID)

	var i models.FootballIncidentPlayer
	err := row.Scan(&i.ID, &i.IncidentID, &i.PlayerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.FootballIncidentPlayer{}, fmt.Errorf("Failed to add the incident player no row found: ", err)
		}
	}
	return i, err
}

const createFootballIncidents = `
WITH tournamentID AS (
	SELECT id FROM tournaments WHERE public_id = $1
),
matchID AS (
	SELECT id FROM matches WHERE public_id = $2
),
teamID AS (
	SELECT id FROM teams WHERE public_id = $3
)
INSERT INTO football_incidents (
	tournament_id,
    match_id,
    team_id,
    periods,
    incident_type,
    incident_time,
    description,
    penalty_shootout_scored
)
SELECT 
	tournamentID.id,
	matchID.id,
	CASE 
        WHEN $5 = 'period' THEN NULL 
        ELSE teamID.id 
    END,
	$4,
	$5,
	$6,
	$7,
	$8
FROM tournamentID
JOIN matchID ON TRUE
LEFT JOIN teamID ON TRUE
RETURNING *;

`

type CreateFootballIncidentsParams struct {
	TournamentPublicID    uuid.UUID  `json:"tournament_public_id"`
	MatchPublicID         uuid.UUID  `json:"match_public_id"`
	TeamPublicID          *uuid.UUID `json:"team_public_id"`
	Periods               string     `json:"periods"`
	IncidentType          string     `json:"incident_type"`
	IncidentTime          int        `json:"incident_time"`
	Description           string     `json:"description"`
	PenaltyShootoutScored bool       `json:"penalty_shootout_scored"`
}

func (q *Queries) CreateFootballIncidents(ctx context.Context, arg CreateFootballIncidentsParams) (*models.FootballIncident, error) {
	row := q.db.QueryRowContext(ctx, createFootballIncidents,
		arg.TournamentPublicID,
		arg.MatchPublicID,
		arg.TeamPublicID,
		arg.Periods,
		arg.IncidentType,
		arg.IncidentTime,
		arg.Description,
		arg.PenaltyShootoutScored,
	)
	var i models.FootballIncident
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.TournamentID,
		&i.MatchID,
		&i.TeamID,
		&i.Periods,
		&i.IncidentType,
		&i.IncidentTime,
		&i.Description,
		&i.PenaltyShootoutScored,
		&i.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return &i, err
}

const getFootballIncidentWithPlayer = `
SELECT 
    fi.id,
	fi.public_id,
	fi.tournament_id,
    fi.match_id, 
    NULL AS team_id, 
    fi.periods, 
    fi.incident_type, 
    fi.incident_time, 
    fi.description, 
    fi.penalty_shootout_scored,
    NULL AS players
FROM 
    football_incidents fi
LEFT JOIN matches AS m ON  m.id = fi.match_id
WHERE 
    m.public_id = $1 AND 
    (fi.periods = 'half_time' OR fi.periods = 'full_time' OR fi.periods = 'extra_time')
UNION ALL
SELECT 
    fi.id,
	fi.public_id,
	fi.tournament_id,
    fi.match_id, 
    fi.team_id, 
    fi.periods, 
    fi.incident_type, 
    fi.incident_time, 
    fi.description, 
    fi.penalty_shootout_scored,
    CASE
        WHEN fi.incident_type = 'substitutions' THEN 
            JSON_BUILD_OBJECT(
                'player_in', JSON_BUILD_OBJECT('id',player_in.id, 'public_id', player_in.public_id, 'user_id', player_in.user_id, 'name', player_in.name, 'slug', player_in.slug, 'short_name',player_in.short_name, 'country', player_in.country, 'positions', player_in.positions, 'media_url', player_in.media_url ),
                'player_out', JSON_BUILD_OBJECT('id',player_out.id, 'public_id', player_out.public_id, 'user_id', player_out.user_id, 'name', player_out.name, 'slug', player_out.slug, 'short_name',player_out.short_name, 'country', player_out.country, 'positions', player_out.positions, 'media_url', player_out.media_url)
            )
        ELSE
            JSON_BUILD_OBJECT(
                'player', JSON_BUILD_OBJECT('id',player_incident.id,'public_id', player_incident.public_id, 'user_id', player_incident.user_id, 'name', player_incident.name, 'slug', player_incident.slug, 'short_name',player_incident.short_name, 'country', player_incident.country, 'positions', player_incident.positions, 'media_url', player_incident.media_url)
            )
    END AS players
FROM 
    football_incidents fi
LEFT JOIN
	matches AS m ON m.id = fi.match_id 
LEFT JOIN 
    football_incident_player AS fip ON fip.incident_id=fi.id
LEFT JOIN 
    players AS player_incident ON player_incident.id = fip.player_id
LEFT JOIN 
    football_substitutions_player AS fis ON fis.incident_id=fi.id
LEFT JOIN 
    players AS player_in ON player_in.id = fis.player_in_id
LEFT JOIN 
    players AS player_out ON player_out.id = fis.player_out_id
WHERE
    m.public_id = $1 AND
    (fi.periods IS NULL OR fi.periods NOT IN ('half_time', 'full_time', 'extra_time'))
ORDER BY
    id ASC;
`

type GetFootballIncidentWithPlayerRow struct {
	ID                    int64       `json:"id"`
	PublicID              uuid.UUID   `json:"public_id"`
	TournamentID          int32       `json:"tournament_id"`
	MatchID               int32       `json:"match_id"`
	TeamID                *int32      `json:"team_id"`
	Periods               string      `json:"periods"`
	IncidentType          string      `json:"incident_type"`
	IncidentTime          int64       `json:"incident_time"`
	Description           string      `json:"description"`
	PenaltyShootoutScored bool        `json:"penalty_shootout_scored"`
	Players               interface{} `json:"players"`
}

func (q *Queries) GetFootballIncidentWithPlayer(ctx context.Context, matchPublicID uuid.UUID) ([]GetFootballIncidentWithPlayerRow, error) {
	rows, err := q.db.QueryContext(ctx, getFootballIncidentWithPlayer, matchPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFootballIncidentWithPlayerRow
	for rows.Next() {
		var i GetFootballIncidentWithPlayerRow
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.TournamentID,
			&i.MatchID,
			&i.TeamID,
			&i.Periods,
			&i.IncidentType,
			&i.IncidentTime,
			&i.Description,
			&i.PenaltyShootoutScored,
			&i.Players,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFootballScoreByIncidentTime = `
SELECT SUM ( CASE WHEN team_id=$3 AND incident_type='goal' THEN 1 ELSE 0 END )
FROM football_incidents fi
WHERE fi.id = $1 AND fi.match_id = $2 AND fi.team_id = $3
`

func (q *Queries) GetFootballScoreByIncidentTime(ctx context.Context, incidentID, matchID, teamID int32) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getFootballScoreByIncidentTime, incidentID, matchID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var sum int64
		if err := rows.Scan(&sum); err != nil {
			return nil, err
		}
		items = append(items, sum)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFootballShootoutScoreByTeam = `
SELECT SUM ( CASE WHEN team_id=$1 AND incident_type='penalty_shootout' AND penalty_shootout_scored='t' THEN 1 ELSE 0 END )
FROM football_incidents fi
JOIN matches m ON  m.id = fi.match_id
WHERE fi.public_id=$1 AND m.public_id = $2 AND fi.team_id = $3
`

func (q *Queries) GetFootballShootoutScoreByTeam(ctx context.Context, incidentPublicID, matchPublicID uuid.UUID, teamID int32) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getFootballShootoutScoreByTeam, incidentPublicID, matchPublicID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var sum int64
		if err := rows.Scan(&sum); err != nil {
			return nil, err
		}
		items = append(items, sum)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFootballIncidentByTeam = `
	SELECT *
	FROM football_incidents fi
	WHERE fi.match_id = $1 AND fi.team_id = $2;
`

func (q *Queries) GetFootballIncidentByTeam(ctx context.Context, matchID int64, teamID int32) (*[]models.FootballIncident, error) {
	row, err := q.db.QueryContext(ctx, getFootballIncidentByTeam, matchID, teamID)
	if err != nil {
		return nil, err
	}
	var stats []models.FootballIncident
	for row.Next() {
		var stat models.FootballIncident
		err := row.Scan(
			&stat.ID,
			&stat.PublicID,
			&stat.MatchID,
			&stat.TournamentID,
			&stat.TeamID,
			&stat.Periods,
			&stat.IncidentType,
			&stat.IncidentTime,
			&stat.Description,
			&stat.PenaltyShootoutScored,
			&stat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		stats = append(stats, stat)
	}
	return &stats, nil
}
