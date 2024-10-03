package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const searchProfile = `
SELECT id, owner, full_name, bio, avatar_url, created_at from profile
WHERE full_name LIKE $1
`

func (q *Queries) SearchProfile(ctx context.Context, fullName string) ([]models.Profile, error) {
	rows, err := q.db.QueryContext(ctx, searchProfile, fullName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Profile
	for rows.Next() {
		var i models.Profile
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.FullName,
			&i.Bio,
			&i.AvatarUrl,
			&i.CreatedAt,
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
