package database

import (
	"context"
	"encoding/json"
	"khelogames/database/models"
	"log"

	"github.com/google/uuid"
)

const createThread = `
INSERT INTO threads (
    user_id,
    community_id,
    title,
    content,
	media_url,
    media_type,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6 CURRENT_TIMESTAMP
) RETURNING *
`

type CreateThreadParams struct {
	UserID      int32  `json:"user_id"`
	CommunityID int32  `json:"community_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MediaUrl    string `json:"media_url"`
	MediaType   string `json:"media_type"`
}

func (q *Queries) CreateThread(ctx context.Context, arg CreateThreadParams) (models.Thread, error) {
	row := q.db.QueryRowContext(ctx, createThread,
		arg.UserID,
		arg.CommunityID,
		arg.Title,
		arg.Content,
		arg.MediaUrl,
		arg.MediaType,
	)
	var i models.Thread
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.CommunityID,
		&i.Title,
		&i.Content,
		&i.MediaUrl,
		&i.MediaType,
		&i.LikeCount,
		&i.CommentCount,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

const deleteThread = `
DELETE FROM threads
WHERE public_id = $1
RETURNING *
`

func (q *Queries) DeleteThread(ctx context.Context, publicID uuid.UUID) (models.Thread, error) {
	row := q.db.QueryRowContext(ctx, deleteThread, publicID)
	var i models.Thread
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.CommunityID,
		&i.Title,
		&i.Content,
		&i.MediaUrl,
		&i.MediaType,
		&i.LikeCount,
		&i.CommentCount,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

const getAllThreads = `
	SELECT 
	JSON_BUILD_OBJECT(
		'id', t.id 'public_id' t.public_id,'user_id',c.user_id, 'community_id', t.community_id, 'title', t.title, 'content', t.content,'media_url', t.media_url, 'media_type', 'like_count', t.like_count, 'comment_count',t.comment_count, 'is_deleted',t.is_deleted, 'created_at',c.created_at,
		'profile', JSON_BUILD_OBJECT('id', p.id, 'public_id',p.public_id, 'user_id',p.user_id,  'username',u.username,  'full_name',p.full_name,  'bio',p.bio,  'avatar_url',p.avatar_url,  'created_at',p.created_at )
	) 
	FROM threads t
	JOIN profile AS p ON p.user_id = c.user_id
	JOIN users AS u ON u.id = c.user_id
`

func (q *Queries) GetAllThreads(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getAllThreads)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []map[string]interface{}
	for rows.Next() {
		var jsonByte []byte
		var data map[string]interface{}
		if err := rows.Scan(&jsonByte); err != nil {
			return nil, err
		}
		err := json.Unmarshal(jsonByte, &data)
		if err != nil {
			log.Fatal("Failed to unmarshal: ", err)
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

const getThreadsByCommunities = `
	SELECT 
	JSON_BUILD_OBJECT(
		'id', t.id 'public_id' t.public_id,'user_id',c.user_id, 'community_id', t.community_id, 'title', t.title, 'content', t.content,'media_url', t.media_url, 'media_type', 'like_count', t.like_count, 'comment_count',t.comment_count, 'is_deleted',t.is_deleted, 'created_at',c.created_at,
		'profile', JSON_BUILD_OBJECT('id', p.id, 'public_id',p.public_id, 'user_id',p.user_id,  'username',u.username,  'full_name',p.full_name,  'bio',p.bio,  'avatar_url',p.avatar_url,  'created_at',p.created_at )
	) 
	FROM threads t
	JOIN profile AS p ON p.user_id = c.user_id
	JOIN users AS u ON u.id = c.user_id
	JOIN communities AS c ON c.id = t.community_id
	WHERE c.public_id = $1
`

func (q *Queries) GetAllThreadsByCommunities(ctx context.Context, communityPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getThreadsByCommunities, communityPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []map[string]interface{}
	for rows.Next() {
		var jsonByte []byte
		var item map[string]interface{}
		if err := rows.Scan(&jsonByte); err != nil {
			return nil, err
		}

		err := json.Unmarshal(jsonByte, &item)
		if err != nil {
			log.Fatal("Failed to unmarshal: ", err)
			return nil, err
		}

		items = append(items, item)
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
SELECT 
	JSON_BUILD_OBJECT(
		'id', t.id 'public_id' t.public_id,'user_id',c.user_id, 'community_id', t.community_id, 'title', t.title, 'content', t.content,'media_url', t.media_url, 'media_type', 'like_count', t.like_count, 'comment_count',t.comment_count, 'is_deleted',t.is_deleted, 'created_at',c.created_at,
		'profile', JSON_BUILD_OBJECT('id', p.id, 'public_id',p.public_id, 'user_id',p.user_id,  'username',u.username,  'full_name',p.full_name,  'bio',p.bio,  'avatar_url',p.avatar_url,  'created_at',p.created_at )
	) 
FROM threads t
JOIN profile AS p ON p.user_id = t.user_id
JOIN users AS u ON u.id = t.user_id
WHERE t.public_id = $1
`

func (q *Queries) GetThread(ctx context.Context, publicID uuid.UUID) (map[string]interface{}, error) {
	row := q.db.QueryRowContext(ctx, getThread, publicID)
	var jsonByte []byte
	var item map[string]interface{}
	if err := row.Scan(&jsonByte); err != nil {
		return nil, err
	}

	err := json.Unmarshal(jsonByte, &item)
	if err != nil {
		log.Fatal("Failed to unmarshal: ", err)
		return nil, err
	}
	return item, err
}

const getThreadByUser = `
SELECT *t. FROM threads
JOIN users AS u ON u.id = t.user_id
WHERE u.publid_id=$1
`

func (q *Queries) GetThreadUser(ctx context.Context, publicID uuid.UUID) ([]models.Thread, error) {
	rows, err := q.db.QueryContext(ctx, getThreadByUser, publicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Thread
	for rows.Next() {
		var i models.Thread
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.CommunityID,
			&i.Title,
			&i.Content,
			&i.MediaUrl,
			&i.MediaType,
			&i.LikeCount,
			&i.CommentCount,
			&i.IsDeleted,
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
	UPDATE threads t
	SET like_count= like_count + 1
	WHERE t.public_id=$1
	RETURNING *
`

type UpdateThreadLikeParams struct {
	PublicID uuid.UUID `json:"public_id"`
}

func (q *Queries) UpdateThreadLike(ctx context.Context, publicID uuid.UUID) (models.Thread, error) {
	row := q.db.QueryRowContext(ctx, updateThreadLike, publicID)
	var i models.Thread
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.CommunityID,
		&i.Title,
		&i.Content,
		&i.MediaUrl,
		&i.MediaType,
		&i.LikeCount,
		&i.CommentCount,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

const updateThreadComment = `
	UPDATE threads t
	SET comment_count = comment_count + 1
	WHERE t.public_id=$1
	RETURNING *
`

func (q *Queries) UpdateThreadCommentCount(ctx context.Context, publicID uuid.UUID) (models.Thread, error) {
	row := q.db.QueryRowContext(ctx, updateThreadLike, publicID)
	var i models.Thread
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.CommunityID,
		&i.Title,
		&i.Content,
		&i.MediaUrl,
		&i.MediaType,
		&i.LikeCount,
		&i.CommentCount,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}
