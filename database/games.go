package database

import (
	"context"
	"khelogames/database/models"
)

const getGame = `
SELECT id, name, min_players FROM games
WHERE id=$1
`

func (q *Queries) GetGame(ctx context.Context, id int64) (models.Game, error) {
	row := q.db.QueryRowContext(ctx, getGame, id)
	var i models.Game
	err := row.Scan(&i.ID, &i.Name, &i.MinPlayers)
	return i, err
}

const getGameByName = `
	SELECT * FROM games
	WHERE name=$1
`

func (q *Queries) GetGamebyName(ctx context.Context, name string) (models.Game, error) {
	row := q.db.QueryRowContext(ctx, getGameByName, name)
	var i models.Game
	err := row.Scan(&i.ID, &i.Name, &i.MinPlayers)
	return i, err
}

const getGames = `
SELECT id, name, min_players FROM games
`

func (q *Queries) GetGames(ctx context.Context) ([]models.Game, error) {
	rows, err := q.db.QueryContext(ctx, getGames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Game
	for rows.Next() {
		var i models.Game
		if err := rows.Scan(&i.ID, &i.Name, &i.MinPlayers); err != nil {
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
