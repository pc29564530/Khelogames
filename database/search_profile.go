package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const searchProfile = `
SELECT * from profile p
JOIN users u ON p.user_id = u.id
WHERE full_name LIKE $1
`

type searchUserProfile struct {
	ID        int64     `json:"id"`
	PublicID  uuid.UUID `json:"public_id"`
	UserID    int32     `json:"user_id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	AvatarUrl string    `json:"avatar_url"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

func (q *Queries) SearchProfile(ctx context.Context, fullName string) ([]searchUserProfile, error) {
	rows, err := q.db.QueryContext(ctx, searchProfile, fullName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []searchUserProfile
	for rows.Next() {
		var i searchUserProfile
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.Username,
			&i.FullName,
			&i.Bio,
			&i.AvatarUrl,
			&i.CreatedAT,
			&i.UpdatedAT,
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
