package token

import (
	"context"
	"fmt"
	"time"

	"khelogames/database"

	"github.com/google/uuid"
)

func CreateNewToken(
	ctx context.Context,
	store *database.Store,
	maker Maker,
	userID int32,
	publicID uuid.UUID,
	accessDuration, refreshDuration time.Duration,
	userAgent, clientIP string,
) (map[string]interface{}, error) {
	fmt.Println("Maker: ", maker)
	accessToken, accessPayload, err := maker.CreateToken(publicID, userID, accessDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshPayload, err := maker.CreateToken(publicID, userID, refreshDuration)
	if err != nil {
		return nil, err
	}

	session, err := store.CreateSessions(ctx, database.CreateSessionsParams{
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"accessToken":    accessToken,
		"accessPayload":  accessPayload,
		"refreshToken":   refreshToken,
		"refreshPayload": refreshPayload,
		"session":        session,
	}, nil
}

func CreateNewTokenTx(
	ctx context.Context,
	q *database.Queries,
	maker Maker,
	userID int32,
	publicID uuid.UUID,
	accessDuration, refreshDuration time.Duration,
	userAgent, clientIP string,
) (map[string]interface{}, error) {

	accessToken, accessPayload, err := maker.CreateToken(publicID, userID, accessDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshPayload, err := maker.CreateToken(publicID, userID, refreshDuration)
	if err != nil {
		return nil, err
	}

	session, err := q.CreateSessions(ctx, database.CreateSessionsParams{
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
	})
	if err != nil {
		return nil, fmt.Errorf("session insert failed: %w", err)
	}

	return map[string]interface{}{
		"accessToken":    accessToken,
		"accessPayload":  accessPayload,
		"refreshToken":   refreshToken,
		"refreshPayload": refreshPayload,
		"session":        session,
	}, nil
}
