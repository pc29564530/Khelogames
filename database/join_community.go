package database

import (
	"context"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addJoinCommunity = `
WITH communityID AS (
	SELECT id FROM communities WHERE public_id = $1
),
userID AS (
	SELECT id FROM users WHERE public_id = $2
)
INSERT INTO join_community (
    communityID.id,
    userID.id
) 
SELECT
	$1,
	$2
FROM communityID, userID
RETURNING *;
`

func (q *Queries) AddJoinCommunity(ctx context.Context, communityID, userID uuid.UUID) (models.JoinCommunity, error) {
	row := q.db.QueryRowContext(ctx, addJoinCommunity, communityID, userID)
	var i models.JoinCommunity
	err := row.Scan(&i.ID, &i.PublicID, &i.CommunityID, &i.UserID)
	return i, err
}

const getCommunityByUser = `
SELECT c.* FROM join_community jc
JOIN communities AS c ON jc.communityID = c.id
JOIN users AS u ON jc.userID = u.id
WHERE u.public_id=$1
ORDER BY id
`

func (q *Queries) GetCommunityByUser(ctx context.Context, userPublicID uuid.UUID) ([]models.Communities, error) {
	rows, err := q.db.QueryContext(ctx, getCommunityByUser, userPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Communities
	for rows.Next() {
		var i models.Communities
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.Name,
			&i.Description,
			&i.IsActive,
			&i.MemberCount,
			&i.AvatarUrl,
			&i.CoverImageUrl,
			&i.CreatedAt,
			&i.UpdatedAt); err != nil {
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
