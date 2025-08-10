package database

import (
	"context"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const checkUserCount = `
SELECT COUNT(*) AS user_count
FROM user_like_thread l1
JOIN users u ON l1.user_id = u.id
JOIN threads t ON t.id = l1.thread_id
WHERE u.public_id = $1
AND t.public_id = $2
`

func (q *Queries) CheckUserCount(ctx context.Context, userPublicID, threadPublicID uuid.UUID) (int, error) {
	row := q.db.QueryRowContext(ctx, checkUserCount, userPublicID, threadPublicID)
	var user_count int
	err := row.Scan(&user_count)
	return user_count, err
}

const countUserLike = `
SELECT COUNT(*) FROM user_like_thread ut
JOIN threads t ON t.id = ut.thread_id
WHERE t.public_id = $1
`

func (q *Queries) CountLikeUser(ctx context.Context, threadID uuid.UUID) (int, error) {
	row := q.db.QueryRowContext(ctx, countUserLike, threadID)
	var count int
	err := row.Scan(&count)
	return count, err
}

const createLike = `
WITH userID AS (
	SELECT id FROM users WHERE public_id = $1
),
threadID AS (
	SELECT id FROM threads WHERE public_id = $2
)
INSERT INTO user_like_thread (
    thread_id,
    user_id
)
SELECT
	threadID.id,
	userID.id
FROM userID, threadID
RETURNING *;
`

func (q *Queries) CreateLike(ctx context.Context, userPublicID, threadPublicID uuid.UUID) (models.UserLikeThread, error) {
	row := q.db.QueryRowContext(ctx, createLike, userPublicID, threadPublicID)
	var i models.UserLikeThread
	err := row.Scan(&i.ID, &i.ThreadID, &i.UserID)
	return i, err
}

const getLike = `
SELECT id, thread_id, user_id FROM user_like_thread
WHERE user_id = $1 LIMIT $1
`

func (q *Queries) GetLike(ctx context.Context, limit int32) (models.UserLikeThread, error) {
	row := q.db.QueryRowContext(ctx, getLike, limit)
	var i models.UserLikeThread
	err := row.Scan(&i.ID, &i.ThreadID, &i.UserID)
	return i, err
}

const userListLike = `
SELECT id, thread_id, user_id FROM user_like_thread
ORDER BY id
`

func (q *Queries) UserListLike(ctx context.Context) ([]models.UserLikeThread, error) {
	rows, err := q.db.QueryContext(ctx, userListLike)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.UserLikeThread
	for rows.Next() {
		var i models.UserLikeThread
		if err := rows.Scan(&i.ID, &i.ThreadID, &i.UserID); err != nil {
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
