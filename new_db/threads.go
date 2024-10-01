package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const createThread = `
INSERT INTO threads (
    username,
    communities_name,
    title,
    content,
    media_type,
    media_url,
    like_count,
    created_at
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP
) RETURNING id, username, communities_name, title, content, media_type, media_url, like_count, created_at
`

type CreateThreadParams struct {
	Username        string `json:"username"`
	CommunitiesName string `json:"communities_name"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	MediaType       string `json:"media_type"`
	MediaUrl        string `json:"media_url"`
	LikeCount       int64  `json:"like_count"`
}

func (q *Queries) CreateThread(ctx context.Context, arg CreateThreadParams) (models.Thread, error) {
	row := q.db.QueryRowContext(ctx, createThread,
		arg.Username,
		arg.CommunitiesName,
		arg.Title,
		arg.Content,
		arg.MediaType,
		arg.MediaUrl,
		arg.LikeCount,
	)
	var i models.Thread
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CommunitiesName,
		&i.Title,
		&i.Content,
		&i.MediaType,
		&i.MediaUrl,
		&i.LikeCount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteThread = `
DELETE FROM threads
WHERE id = $1
RETURNING id, username, communities_name, title, content, media_type, media_url, like_count, created_at
`

func (q *Queries) DeleteThread(ctx context.Context, id int64) (models.Thread, error) {
	row := q.db.QueryRowContext(ctx, deleteThread, id)
	var i models.Thread
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CommunitiesName,
		&i.Title,
		&i.Content,
		&i.MediaType,
		&i.MediaUrl,
		&i.LikeCount,
		&i.CreatedAt,
	)
	return i, err
}

const getAllThreads = `
SELECT id, username, communities_name, title, content, media_type, media_url, like_count, created_at FROM threads
`

func (q *Queries) GetAllThreads(ctx context.Context) ([]models.Thread, error) {
	rows, err := q.db.QueryContext(ctx, getAllThreads)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Thread
	for rows.Next() {
		var i models.Thread
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.CommunitiesName,
			&i.Title,
			&i.Content,
			&i.MediaType,
			&i.MediaUrl,
			&i.LikeCount,
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

const getAllThreadsByCommunities = `
SELECT id, username, communities_name, title, content, media_type, media_url, like_count, created_at FROM threads
WHERE communities_name = $1
`

func (q *Queries) GetAllThreadsByCommunities(ctx context.Context, communitiesName string) ([]models.Thread, error) {
	rows, err := q.db.QueryContext(ctx, getAllThreadsByCommunities, communitiesName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Thread
	for rows.Next() {
		var i models.Thread
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.CommunitiesName,
			&i.Title,
			&i.Content,
			&i.MediaType,
			&i.MediaUrl,
			&i.LikeCount,
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

const getThread = `
SELECT id, username, communities_name, title, content, media_type, media_url, like_count, created_at FROM threads
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetThread(ctx context.Context, id int64) (models.Thread, error) {
	row := q.db.QueryRowContext(ctx, getThread, id)
	var i models.Thread
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CommunitiesName,
		&i.Title,
		&i.Content,
		&i.MediaType,
		&i.MediaUrl,
		&i.LikeCount,
		&i.CreatedAt,
	)
	return i, err
}

const getThreadByThreadID = `
SELECT id, username, communities_name, title, content, media_type, media_url, like_count, created_at FROM threads
WHERE id = $1
`

func (q *Queries) GetThreadByThreadID(ctx context.Context, id int64) (models.Thread, error) {
	row := q.db.QueryRowContext(ctx, getThreadByThreadID, id)
	var i models.Thread
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CommunitiesName,
		&i.Title,
		&i.Content,
		&i.MediaType,
		&i.MediaUrl,
		&i.LikeCount,
		&i.CreatedAt,
	)
	return i, err
}

const getThreadUser = `
SELECT id, username, communities_name, title, content, media_type, media_url, like_count, created_at FROM threads
WHERE username=$1
`

func (q *Queries) GetThreadUser(ctx context.Context, username string) ([]models.Thread, error) {
	rows, err := q.db.QueryContext(ctx, getThreadUser, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Thread
	for rows.Next() {
		var i models.Thread
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.CommunitiesName,
			&i.Title,
			&i.Content,
			&i.MediaType,
			&i.MediaUrl,
			&i.LikeCount,
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

const updateThreadLike = `
UPDATE threads
SET like_count=$1
WHERE id=$2
RETURNING id, username, communities_name, title, content, media_type, media_url, like_count, created_at
`

type UpdateThreadLikeParams struct {
	LikeCount int64 `json:"like_count"`
	ID        int64 `json:"id"`
}

func (q *Queries) UpdateThreadLike(ctx context.Context, arg UpdateThreadLikeParams) (models.Thread, error) {
	row := q.db.QueryRowContext(ctx, updateThreadLike, arg.LikeCount, arg.ID)
	var i models.Thread
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CommunitiesName,
		&i.Title,
		&i.Content,
		&i.MediaType,
		&i.MediaUrl,
		&i.LikeCount,
		&i.CreatedAt,
	)
	return i, err
}
