package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const addFootballSubstitution = `
INSERT INTO football_substitution (
    team_id,
    player_id,
    match_id,
    position
) VALUES ( $1, $2, $3, $4 )
RETURNING id, team_id, player_id, match_id, position
`

type AddFootballSubstitutionParams struct {
	TeamID   int64  `json:"team_id"`
	PlayerID int64  `json:"player_id"`
	MatchID  int64  `json:"match_id"`
	Position string `json:"position"`
}

func (q *Queries) AddFootballSubstitution(ctx context.Context, arg AddFootballSubstitutionParams) (models.FootballSubstitution, error) {
	row := q.db.QueryRowContext(ctx, addFootballSubstitution,
		arg.TeamID,
		arg.PlayerID,
		arg.MatchID,
		arg.Position,
	)
	var i models.FootballSubstitution
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.PlayerID,
		&i.MatchID,
		&i.Position,
	)
	return i, err
}

const getFootballSubstitution = `
SELECT id, team_id, player_id, match_id, position FROM football_substitution
WHERE match_id=$1 AND team_id=$2
`

type GetFootballSubstitutionParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) GetFootballSubstitution(ctx context.Context, arg GetFootballSubstitutionParams) ([]models.FootballSubstitution, error) {
	rows, err := q.db.QueryContext(ctx, getFootballSubstitution, arg.MatchID, arg.TeamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.FootballSubstitution
	for rows.Next() {
		var i models.FootballSubstitution
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.PlayerID,
			&i.MatchID,
			&i.Position,
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
