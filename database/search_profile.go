package database

import (
	"context"
	"khelogames/database/models"
)

const searchProfile = `
SELECT * from profile
WHERE full_name LIKE $1
`

func (q *Queries) SearchProfile(ctx context.Context, fullName string) ([]models.UsersProfile, error) {
	rows, err := q.db.QueryContext(ctx, searchProfile, fullName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.UsersProfile
	for rows.Next() {
		var i models.UsersProfile
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Username,
			&i.FullName,
			&i.Bio,
			&i.AvatarUrl,
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
