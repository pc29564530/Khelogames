package database

import (
	"context"
	"khelogames/database/models"
)

const getGroups = `
	SELECT * FROM groups;
`

func (q *Queries) GetGroups(ctx context.Context) ([]models.Group, error) {
	rows, err := q.db.QueryContext(ctx, getGroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Group
	for rows.Next() {
		var i models.Group
		if err := rows.Scan(
			&i.ID,
			&i.Name,
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

const getGroup = `
	SELECT * FROM groups
	WHERE id=$1;
`

func (q *Queries) GetGroup(ctx context.Context) ([]models.Group, error) {
	rows, err := q.db.QueryContext(ctx, getGroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Group
	for rows.Next() {
		var i models.Group
		if err := rows.Scan(
			&i.ID,
			&i.Name,
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
