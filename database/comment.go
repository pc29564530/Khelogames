package database

import (
	"context"
	"encoding/json"
	"khelogames/database/models"
	"log"

	"github.com/google/uuid"
)

const createComment = `
WITH threadID AS (
	SELECT * from threads 
	WHERE public_id = $1
),
userID AS (
	SELECT * from users
	WHERE public_id = $2
)
INSERT INTO comments (
    thread_id,
    user_id,
	comment_text,
    created_at
) 
SELECT
	threadID.id,
	userID.id,
	$3,
	CURRENT_TIMESTAMP
FROM threadID, userID
RETURNING *;
`

func (q *Queries) CreateComment(ctx context.Context, threadPublicID, commentPublicID uuid.UUID, commentText string) (models.Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment, threadPublicID, commentPublicID, commentText)
	var i models.Comment
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.ThreadID,
		&i.UserID,
		&i.ParentCommentID,
		&i.CommentText,
		&i.LikeCount,
		&i.ReplyCount,
		&i.IsDeleted,
		&i.IsEdited,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteComment = `
DELETE FROM comment
JOIN users AS u ON u.id = c.user_id
WHERE c.public_id=$1 AND u.public_id=$2
RETURNING *
`

func (q *Queries) DeleteComment(ctx context.Context, commentPublicID, userPublicID uuid.UUID) (models.Comment, error) {
	row := q.db.QueryRowContext(ctx, deleteComment, commentPublicID, userPublicID)
	var i models.Comment
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.ThreadID,
		&i.UserID,
		&i.ParentCommentID,
		&i.CommentText,
		&i.LikeCount,
		&i.ReplyCount,
		&i.IsDeleted,
		&i.IsEdited,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllComment = `
SELECT 
JSON_BUILD_OBJECT(
	'id', c.id 'public_id' c.public_id, 'thread_id', c.thread_id,'user_id',c.user_id, 'parent_comment_id', c.parent_comment_id, 'comment_text', c.comment_text,'like_count', c.like_count,'reply_count',c.reply_count, 'is_deleted',c.is_deleted, 'is_edited',c.is_edited, 'created_at',c.created_at, 'updated_at'c.updated_at,,
	'profile', JSON_BUILD_OBJECT('id', p.id, 'public_id',p.public_id, 'user_id',p.user_id,  'username',u.username,  'full_name',p.full_name,  'bio',p.bio,  'avatar_url',p.avatar_url,  'created_at',p.created_at )
) 
FROM comment c
JOIN threads AS t ON t.id = c.thread_id
JOIN profile AS p ON p.user_id = c.user_id
JOIN users AS u ON u.id = c.user_id
WHERE t.public_id=$1
ORDER BY c.id;
`

func (q *Queries) GetAllComment(ctx context.Context, publicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getAllComment, publicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []map[string]interface{}
	for rows.Next() {
		var jsonBytes []byte
		var data map[string]interface{}
		if err := rows.Scan(&jsonBytes); err != nil {
			return nil, err
		}
		err := json.Unmarshal(jsonBytes, &data)
		if err != nil {
			log.Fatal("Failed to unmarshal the ")
		}
		items = append(items, data)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
