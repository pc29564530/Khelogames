package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khelogames/database/models"
	"time"

	"github.com/google/uuid"
)

const createNewMessage = `
WITH senderID AS (
    SELECT id FROM users WHERE public_id = $1
),
receiverID AS (
    SELECT user_id FROM user_profiles WHERE public_id = $2
)
INSERT INTO message (
    sender_id,
    receiver_id,
    content,
    media_url,
    media_type,
    created_at,
	sent_at
)
SELECT 
    senderID.id,
    receiverID.user_id,
    $3,
    $4,
    $5,
    CURRENT_TIMESTAMP,
	$6
FROM senderID, receiverID
RETURNING *;
`

type CreateNewMessageParams struct {
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	MediaUrl   string    `json:"media_url"`
	MediaType  string    `json:"media_type"`
	SentAt     time.Time `json:"sent_at"`
}

func (q *Queries) CreateNewMessage(ctx context.Context, arg CreateNewMessageParams) (models.Message, error) {
	row := q.db.QueryRowContext(ctx, createNewMessage,
		arg.SenderID,
		arg.ReceiverID,
		arg.Content,
		arg.MediaUrl,
		arg.MediaType,
		arg.SentAt,
	)
	var i models.Message
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.SenderID,
		&i.ReceiverID,
		&i.Content,
		&i.MediaUrl,
		&i.MediaType,
		&i.IsSeen,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.SentAt,
		&i.IsDelivered,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return i, err
}

const deleteMessage = `
DELETE FROM message
WHERE sender_id=$1 and receiver_id=$2
RETURNING *
`

type DeleteMessageParams struct {
	SenderID string `json:"sender_id"`
	ID       int64  `json:"id"`
}

func (q *Queries) DeleteMessage(ctx context.Context, arg DeleteMessageParams) (models.Message, error) {
	row := q.db.QueryRowContext(ctx, deleteMessage, arg.SenderID, arg.ID)
	var i models.Message
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.SenderID,
		&i.ReceiverID,
		&i.Content,
		&i.MediaUrl,
		&i.MediaType,
		&i.IsSeen,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.SentAt,
		&i.IsDelivered,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return i, err
}

// fetch all the message by the receiver id
const getMessageByReceiver = `
SELECT JSON_BUILD_OBJECT(
	'public_id', m.public_id,
	'sender', JSON_BUILD_OBJECT(
		'public_id', sp.public_id,
		'username', sender.username,
		'full_name', sender.full_name,
		'bio', sp.bio,
		'avatar_url', sp.avatar_url,
		'created_at', sp.created_at
	),
	'receiver', JSON_BUILD_OBJECT(
		'public_id', rp.public_id,
		'username', receiver.username,
		'full_name', receiver.full_name,
		'bio', rp.bio,
		'avatar_url', rp.avatar_url,
		'created_at', rp.created_at
	),
	'content', m.content,
	'media_url', m.media_url,
	'media_type', m.media_type,
	'is_seen', m.is_seen,
	'is_deleted', m.is_deleted,
	'created_at', m.created_at,
	'sent_at', m.sent_at,
	'is_delivered', m.is_delivered
)
FROM message m
JOIN users sender   ON sender.id = m.sender_id
JOIN user_profiles sp ON sp.user_id = sender.id
JOIN users receiver ON receiver.id = m.receiver_id
JOIN user_profiles rp ON rp.user_id = receiver.id
WHERE 
   (sender.public_id = $1 AND receiver.public_id = (
       SELECT u.public_id 
       FROM user_profiles up
       JOIN users u ON up.user_id = u.id
       WHERE up.public_id = $2
   ))
   OR
   (sender.public_id = (
       SELECT u.public_id 
       FROM user_profiles up
       JOIN users u ON up.user_id = u.id
       WHERE up.public_id = $2
   ) AND receiver.public_id = $1)
ORDER BY m.id ASC;
`

type GetMessageByReceiverParams struct {
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
}

func (q *Queries) GetMessageByReceiver(ctx context.Context, arg GetMessageByReceiverParams) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getMessageByReceiver, arg.SenderID, arg.ReceiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []map[string]interface{}
	for rows.Next() {
		var item map[string]interface{}
		var jsonByte []byte
		if err := rows.Scan(&jsonByte); err != nil {
			return nil, err
		}
		err := json.Unmarshal(jsonByte, &item)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal: ", err)
		}
		items = append(items, item)
	}
	// fmt.Println("Message Items; ", items)
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// get users by message sender
const getUserByMessageSender = `
	SELECT DISTINCT ON (up.id)
		JSON_BUILD_OBJECT(
			'public_id', up.public_id,
			'username', u.username,
			'full_name', u.full_name,
			'avatar_url', up.avatar_url
		) AS sender_profile
	FROM message m
	JOIN users u ON u.id = m.sender_id
	JOIN user_profiles up ON up.user_id = u.id
	WHERE m.receiver_id = (SELECT id FROM users WHERE public_id = $1)
	ORDER BY up.id, m.created_at DESC;
`

func (q *Queries) GetUserByMessageSend(ctx context.Context, senderID uuid.UUID) ([]map[string]interface{}, error) {
	fmt.Println("Sender Public ID: ", senderID)

	rows, err := q.db.QueryContext(ctx, getUserByMessageSender, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var jsonByte []byte
		if err := rows.Scan(&jsonByte); err != nil {
			return nil, err
		}

		item := make(map[string]interface{})
		if err := json.Unmarshal(jsonByte, &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// schedule delete message
const scheduledDeleteMessage = `
DELETE FROM message
WHERE is_deleted = TRUE AND deleted_at < NOW() - INTERVAL '30 days'
RETURNING *
`

func (q *Queries) ScheduledDeleteMessage(ctx context.Context) ([]models.Message, error) {
	rows, err := q.db.QueryContext(ctx, scheduledDeleteMessage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Message
	for rows.Next() {
		var i models.Message
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.SenderID,
			&i.ReceiverID,
			&i.Content,
			&i.MediaUrl,
			&i.MediaType,
			&i.IsSeen,
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
