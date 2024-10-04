package database

import (
	"context"
	"khelogames/database/models"
)

const createComment = `
INSERT INTO comment (
    thread_id,
    owner,
	comment_text,
    created_at
) VALUES (
    $1, $2, $3, CURRENT_TIMESTAMP
)
RETURNING id, thread_id, owner, comment_text, created_at
`

type CreateCommentParams struct {
	ThreadID    int64  `json:"thread_id"`
	Owner       string `json:"owner"`
	CommentText string `json:"comment_text"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (models.Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment, arg.ThreadID, arg.Owner, arg.CommentText)
	var i models.Comment
	err := row.Scan(
		&i.ID,
		&i.ThreadID,
		&i.Owner,
		&i.CommentText,
		&i.CreatedAt,
	)
	return i, err
}

const deleteComment = `
DELETE FROM comment
WHERE id=$1 AND owner=$2
RETURNING id, thread_id, owner, comment_text, created_at
`

type DeleteCommentParams struct {
	ID    int64  `json:"id"`
	Owner string `json:"owner"`
}

func (q *Queries) DeleteComment(ctx context.Context, arg DeleteCommentParams) (models.Comment, error) {
	row := q.db.QueryRowContext(ctx, deleteComment, arg.ID, arg.Owner)
	var i models.Comment
	err := row.Scan(
		&i.ID,
		&i.ThreadID,
		&i.Owner,
		&i.CommentText,
		&i.CreatedAt,
	)
	return i, err
}

const getAllComment = `
SELECT id, thread_id, owner, comment_text, created_at FROM comment
WHERE thread_id=$1
`

func (q *Queries) GetAllComment(ctx context.Context, threadID int64) ([]models.Comment, error) {
	rows, err := q.db.QueryContext(ctx, getAllComment, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Comment
	for rows.Next() {
		var i models.Comment
		if err := rows.Scan(
			&i.ID,
			&i.ThreadID,
			&i.Owner,
			&i.CommentText,
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

const getCommentByUser = `
SELECT id, thread_id, owner, comment_text, created_at FROM comment
WHERE owner=$1
`

func (q *Queries) GetCommentByUser(ctx context.Context, owner string) ([]models.Comment, error) {
	rows, err := q.db.QueryContext(ctx, getCommentByUser, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Comment
	for rows.Next() {
		var i models.Comment
		if err := rows.Scan(
			&i.ID,
			&i.ThreadID,
			&i.Owner,
			&i.CommentText,
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
