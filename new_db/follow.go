package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const createFollowing = `
INSERT INTO follow (
    follower_owner,
    following_owner,
    created_at
) VALUES (
             $1, $2, CURRENT_TIMESTAMP
) RETURNING id, follower_owner, following_owner, created_at
`

type CreateFollowingParams struct {
	FollowerOwner  string `json:"follower_owner"`
	FollowingOwner string `json:"following_owner"`
}

func (q *Queries) CreateFollowing(ctx context.Context, arg CreateFollowingParams) (models.Follow, error) {
	row := q.db.QueryRowContext(ctx, createFollowing, arg.FollowerOwner, arg.FollowingOwner)
	var i models.Follow
	err := row.Scan(
		&i.ID,
		&i.FollowerOwner,
		&i.FollowingOwner,
		&i.CreatedAt,
	)
	return i, err
}

const deleteFollowing = `
DELETE FROM follow
WHERE following_owner = $1 RETURNING id, follower_owner, following_owner, created_at
`

func (q *Queries) DeleteFollowing(ctx context.Context, followingOwner string) (models.Follow, error) {
	row := q.db.QueryRowContext(ctx, deleteFollowing, followingOwner)
	var i models.Follow
	err := row.Scan(
		&i.ID,
		&i.FollowerOwner,
		&i.FollowingOwner,
		&i.CreatedAt,
	)
	return i, err
}

const getAllFollower = `
SELECT DISTINCT follower_owner FROM follow
WHERE following_owner = $1
`

func (q *Queries) GetAllFollower(ctx context.Context, followingOwner string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getAllFollower, followingOwner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var follower_owner string
		if err := rows.Scan(&follower_owner); err != nil {
			return nil, err
		}
		items = append(items, follower_owner)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllFollowing = `
SELECT DISTINCT following_owner FROM follow
WHERE follower_owner =  $1
`

func (q *Queries) GetAllFollowing(ctx context.Context, followerOwner string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getAllFollowing, followerOwner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var following_owner string
		if err := rows.Scan(&following_owner); err != nil {
			return nil, err
		}
		items = append(items, following_owner)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
