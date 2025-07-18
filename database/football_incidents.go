package database

import (
	"context"
	"khelogames/database/models"
)

const aDDFootballSubsPlayer = `
INSERT INTO football_substitutions_player (
    incident_id,
    player_in_id,
    player_out_id
) VALUES ($1, $2, $3)
RETURNING id, incident_id, player_in_id, player_out_id
`

type ADDFootballSubsPlayerParams struct {
	IncidentID  int64 `json:"incident_id"`
	PlayerInID  int64 `json:"player_in_id"`
	PlayerOutID int64 `json:"player_out_id"`
}

func (q *Queries) ADDFootballSubsPlayer(ctx context.Context, arg ADDFootballSubsPlayerParams) (models.FootballSubstitutionsPlayer, error) {
	row := q.db.QueryRowContext(ctx, aDDFootballSubsPlayer, arg.IncidentID, arg.PlayerInID, arg.PlayerOutID)
	var i models.FootballSubstitutionsPlayer
	err := row.Scan(
		&i.ID,
		&i.IncidentID,
		&i.PlayerInID,
		&i.PlayerOutID,
	)
	return i, err
}

const addFootballIncidentPlayer = `
INSERT INTO football_incident_player (
    incident_id,
    player_id
) VALUES ($1, $2)
RETURNING id, incident_id, player_id
`

type AddFootballIncidentPlayerParams struct {
	IncidentID int64 `json:"incident_id"`
	PlayerID   int64 `json:"player_id"`
}

func (q *Queries) AddFootballIncidentPlayer(ctx context.Context, arg AddFootballIncidentPlayerParams) (models.FootballIncidentPlayer, error) {
	row := q.db.QueryRowContext(ctx, addFootballIncidentPlayer, arg.IncidentID, arg.PlayerID)
	var i models.FootballIncidentPlayer
	err := row.Scan(&i.ID, &i.IncidentID, &i.PlayerID)
	return i, err
}

const createFootballIncidents = `
INSERT INTO football_incidents (
    match_id,
    team_id,
    periods,
    incident_type,
    incident_time,
    description,
    penalty_shootout_scored,
	tournament_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, match_id, team_id, periods, incident_type, incident_time, description, created_at, penalty_shootout_scored, tournament_id
`

type CreateFootballIncidentsParams struct {
	MatchID               int64  `json:"match_id"`
	TeamID                *int64 `json:"team_id"`
	Periods               string `json:"periods"`
	IncidentType          string `json:"incident_type"`
	IncidentTime          int64  `json:"incident_time"`
	Description           string `json:"description"`
	PenaltyShootoutScored bool   `json:"penalty_shootout_scored"`
	TournamentID          int32  `json:"tournament_id"`
}

func (q *Queries) CreateFootballIncidents(ctx context.Context, arg CreateFootballIncidentsParams) (models.FootballIncident, error) {
	row := q.db.QueryRowContext(ctx, createFootballIncidents,
		arg.MatchID,
		arg.TeamID,
		arg.Periods,
		arg.IncidentType,
		arg.IncidentTime,
		arg.Description,
		arg.PenaltyShootoutScored,
		arg.TournamentID,
	)
	var i models.FootballIncident
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.Periods,
		&i.IncidentType,
		&i.IncidentTime,
		&i.Description,
		&i.CreatedAt,
		&i.PenaltyShootoutScored,
		&i.TournamentID,
	)
	return i, err
}

const getFootballIncidentByGoal = `
SELECT (fi.id, fi.match_id, fi.team_id, fi.incident_type, fi.incident_time, fi.player_id, fi.assist_player_id, fi.description, fi.created_at) FROM football_incidents AS fi
WHERE match_id=$1 AND incident_type="goal"
ORDER BY incident_time DESC
`

func (q *Queries) GetFootballIncidentByGoal(ctx context.Context, matchID int64) ([]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getFootballIncidentByGoal, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var column_1 interface{}
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFootballIncidentBySubstitution = `
SELECT (fi.id, fi.match_id, fi.team_id, fi.incident_type, fi.incident_time, fi.substitution_in_player_id, fi.substitution_out_player_id, fi.description, fi.created_at) FROM football_incidents AS fi
WHERE match_id=$1 AND incident_type="substitutions"
ORDER BY incident_time DESC
`

func (q *Queries) GetFootballIncidentBySubstitution(ctx context.Context, matchID int64) ([]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getFootballIncidentBySubstitution, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var column_1 interface{}
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFootballIncidentPlayer = `
SELECT id, incident_id, player_id FROM football_incident_player
WHERE incident_id=$1
`

func (q *Queries) GetFootballIncidentPlayer(ctx context.Context, incidentID int64) (models.FootballIncidentPlayer, error) {
	row := q.db.QueryRowContext(ctx, getFootballIncidentPlayer, incidentID)
	var i models.FootballIncidentPlayer
	err := row.Scan(&i.ID, &i.IncidentID, &i.PlayerID)
	return i, err
}

const getFootballIncidentSubsPlayer = `
SELECT id, incident_id, player_in_id, player_out_id FROM football_substitutions_player
WHERE incident_id=$1
`

func (q *Queries) GetFootballIncidentSubsPlayer(ctx context.Context, incidentID int64) (models.FootballSubstitutionsPlayer, error) {
	row := q.db.QueryRowContext(ctx, getFootballIncidentSubsPlayer, incidentID)
	var i models.FootballSubstitutionsPlayer
	err := row.Scan(
		&i.ID,
		&i.IncidentID,
		&i.PlayerInID,
		&i.PlayerOutID,
	)
	return i, err
}

const getFootballIncidentWithPlayer = `
SELECT 
    fi.id, 
    fi.match_id, 
    NULL AS team_id, 
    fi.periods, 
    fi.incident_type, 
    fi.incident_time, 
    fi.description, 
    fi.penalty_shootout_scored,
	fi.tournament_id,
    NULL AS players
FROM 
    football_incidents fi
WHERE 
    fi.match_id = $1 AND 
    (fi.periods = 'half_time' OR fi.periods = 'full_time' OR fi.periods = 'extra_time')
UNION ALL
SELECT 
    fi.id, 
    fi.match_id, 
    fi.team_id, 
    fi.periods, 
    fi.incident_type, 
    fi.incident_time, 
    fi.description, 
    fi.penalty_shootout_scored,
	fi.tournament_id,
    CASE
        WHEN fi.incident_type = 'substitutions' THEN 
            JSON_BUILD_OBJECT(
                'player_in', JSON_BUILD_OBJECT('id',player_in.id,'username',player_in.username, 'name', player_in.player_name, 'slug', player_in.slug, 'short_name',player_in.short_name, 'country', player_in.country, 'positions', player_in.positions, 'media_url', player_in.media_url ),
                'player_out', JSON_BUILD_OBJECT('id',player_out.id,'username',player_out.username, 'name', player_out.player_name, 'slug', player_out.slug, 'short_name',player_out.short_name, 'country', player_out.country, 'positions', player_out.positions, 'media_url', player_out.media_url)
            )
        ELSE
            JSON_BUILD_OBJECT(
                'player', JSON_BUILD_OBJECT('id',player_incident.id,'username',player_incident.username, 'name', player_incident.player_name, 'slug', player_incident.slug, 'short_name',player_incident.short_name, 'country', player_incident.country, 'positions', player_incident.positions, 'media_url', player_incident.media_url)
            )
    END AS players
FROM 
    football_incidents fi
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
    fi.match_id = $1 AND 
    (fi.periods IS NULL OR fi.periods NOT IN ('half_time', 'full_time', 'extra_time'))
ORDER BY 
    incident_time DESC;
`

type GetFootballIncidentWithPlayerRow struct {
	ID                    int64       `json:"id"`
	MatchID               int64       `json:"match_id"`
	TeamID                *int64      `json:"team_id"`
	Periods               string      `json:"periods"`
	IncidentType          string      `json:"incident_type"`
	IncidentTime          int64       `json:"incident_time"`
	Description           string      `json:"description"`
	PenaltyShootoutScored bool        `json:"penalty_shootout_scored"`
	TournamentID          int64       `json:"tournament_id"`
	Players               interface{} `json:"players"`
}

func (q *Queries) GetFootballIncidentWithPlayer(ctx context.Context, matchID int64) ([]GetFootballIncidentWithPlayerRow, error) {
	rows, err := q.db.QueryContext(ctx, getFootballIncidentWithPlayer, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFootballIncidentWithPlayerRow
	for rows.Next() {
		var i GetFootballIncidentWithPlayerRow
		if err := rows.Scan(
			&i.ID,
			&i.MatchID,
			&i.TeamID,
			&i.Periods,
			&i.IncidentType,
			&i.IncidentTime,
			&i.Description,
			&i.PenaltyShootoutScored,
			&i.TournamentID,
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

const getFootballIncidents = `
SELECT id, match_id, team_id, periods, incident_type, incident_time, description, created_at, penalty_shootout_scored FROM football_incidents
WHERE match_id=$1
ORDER BY created_at DESC
`

func (q *Queries) GetFootballIncidents(ctx context.Context, matchID int64) ([]models.FootballIncident, error) {
	rows, err := q.db.QueryContext(ctx, getFootballIncidents, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.FootballIncident
	for rows.Next() {
		var i models.FootballIncident
		if err := rows.Scan(
			&i.ID,
			&i.MatchID,
			&i.TeamID,
			&i.Periods,
			&i.IncidentType,
			&i.IncidentTime,
			&i.Description,
			&i.CreatedAt,
			&i.PenaltyShootoutScored,
			&i.TournamentID,
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
SELECT SUM ( CASE WHEN team_id=$1 AND incident_type='goal' THEN 1 ELSE 0 END )
FROM football_incidents
WHERE match_id = $2 AND id <= $3
`

type GetFootballScoreByIncidentTimeParams struct {
	TeamID  int64 `json:"team_id"`
	MatchID int64 `json:"match_id"`
	ID      int64 `json:"id"`
}

func (q *Queries) GetFootballScoreByIncidentTime(ctx context.Context, arg GetFootballScoreByIncidentTimeParams) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getFootballScoreByIncidentTime, arg.TeamID, arg.MatchID, arg.ID)
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
FROM football_incidents
WHERE match_id=$2 AND id <= $3
`

type GetFootballShootoutScoreByTeamParams struct {
	TeamID  int64 `json:"team_id"`
	MatchID int64 `json:"match_id"`
	ID      int64 `json:"id"`
}

func (q *Queries) GetFootballShootoutScoreByTeam(ctx context.Context, arg GetFootballShootoutScoreByTeamParams) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getFootballShootoutScoreByTeam, arg.TeamID, arg.MatchID, arg.ID)
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
