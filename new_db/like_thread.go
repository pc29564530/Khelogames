package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const checkUserCount = `
SELECT COUNT(*) AS user_count
FROM like_thread l1
JOIN users u ON l1.username = u.username
WHERE l1.thread_id = $1
AND u.username = $2
`

type CheckUserCountParams struct {
	ThreadID int64  `json:"thread_id"`
	Username string `json:"username"`
}

func (q *Queries) CheckUserCount(ctx context.Context, arg CheckUserCountParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkUserCount, arg.ThreadID, arg.Username)
	var user_count int64
	err := row.Scan(&user_count)
	return user_count, err
}

const countLikeUser = `
SELECT COUNT(*) FROM like_thread
WHERE thread_id = $1
`

func (q *Queries) CountLikeUser(ctx context.Context, threadID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, countLikeUser, threadID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createLike = `
INSERT INTO like_thread (
    thread_id,
    username
) VALUES (
  $1, $2
) RETURNING id, thread_id, username
`

type CreateLikeParams struct {
	ThreadID int64  `json:"thread_id"`
	Username string `json:"username"`
}

func (q *Queries) CreateLike(ctx context.Context, arg CreateLikeParams) (models.LikeThread, error) {
	row := q.db.QueryRowContext(ctx, createLike, arg.ThreadID, arg.Username)
	var i models.LikeThread
	err := row.Scan(&i.ID, &i.ThreadID, &i.Username)
	return i, err
}

const getLike = `-- name: GetLike :one
SELECT id, thread_id, username FROM like_thread
WHERE username = $1 LIMIT $1
`

func (q *Queries) GetLike(ctx context.Context, limit int32) (models.LikeThread, error) {
	row := q.db.QueryRowContext(ctx, getLike, limit)
	var i models.LikeThread
	err := row.Scan(&i.ID, &i.ThreadID, &i.Username)
	return i, err
}

const userListLike = `-- name: UserListLike :many
SELECT id, thread_id, username FROM like_thread
ORDER BY username
`

func (q *Queries) UserListLike(ctx context.Context) ([]models.LikeThread, error) {
	rows, err := q.db.QueryContext(ctx, userListLike)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.LikeThread
	for rows.Next() {
		var i models.LikeThread
		if err := rows.Scan(&i.ID, &i.ThreadID, &i.Username); err != nil {
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
