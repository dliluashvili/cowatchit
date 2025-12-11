package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dliluashvili/cowatchit/internal/helpers"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/dliluashvili/cowatchit/internal/shared/constants"

	"github.com/redis/go-redis/v9"
)

type SessionService struct {
	redisClient *redis.Client
}

func NewSessionService(rc *redis.Client) *SessionService {
	return &SessionService{
		redisClient: rc,
	}
}

func (s *SessionService) CreateAndSave(ctx context.Context, user *models.User) (*models.Session, error) {
	sessionId, err := helpers.GenerateSessionID()

	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(constants.SessionDuration)

	session := &models.Session{
		User:      user,
		SessionID: sessionId,
		ExpiresAt: expiresAt,
	}

	data, err := json.Marshal(session)

	if err != nil {
		return nil, err
	}

	// Redis keys
	sessionKey := fmt.Sprintf("session:%s", sessionId)
	userSessionKey := fmt.Sprintf("user_session:%s", user.ID.String())

	ttl := time.Until(expiresAt)

	// Step 1: Check if user already has a session
	oldSessionId, err := s.redisClient.Get(ctx, userSessionKey).Result()

	if err == nil {
		// If old session exists, delete old session key
		oldSessionKey := fmt.Sprintf("session:%s", oldSessionId)
		_ = s.redisClient.Del(ctx, oldSessionKey).Err() // Ignore error
	}

	// Step 2: Save new session
	if err := s.redisClient.Set(ctx, sessionKey, data, ttl).Err(); err != nil {
		return nil, err
	}

	// Step 3: Update user → session ID mapping
	if err := s.redisClient.Set(ctx, userSessionKey, sessionId, ttl).Err(); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) GetUserBySession(ctx context.Context, sessionId string) (*models.Session, error) {
	sessionKey := fmt.Sprintf("session:%s", sessionId)

	val, err := s.redisClient.Get(ctx, sessionKey).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}

	var session models.Session

	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, fmt.Errorf("failed to parse session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	if session.User == nil {
		return nil, fmt.Errorf("user not found in session")
	}

	return &session, nil
}

func (s *SessionService) SetSocketId(ctx context.Context, sessionId string, socketId *string) error {
	key := fmt.Sprintf("session:%s", sessionId)

	val, err := s.redisClient.Get(ctx, key).Result()

	if err != nil {
		return fmt.Errorf("failed to fetch session: %w", err)
	}

	// Unmarshal existing session
	var session models.Session
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Update socket ID
	session.SocketID = socketId

	// Marshal and overwrite session
	updated, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal updated session: %w", err)
	}

	// Re-save with remaining TTL
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	if err := s.redisClient.Set(ctx, key, updated, ttl).Err(); err != nil {
		return fmt.Errorf("failed to update session in Redis: %w", err)
	}

	return nil
}
