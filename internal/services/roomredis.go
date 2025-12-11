package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	roomPrefix         = "room:"
	roomUsersSetPrefix = "room:users:"
	activeRoomsKey     = "active:rooms"

	defaultRoomTTL = 24 * time.Hour
)

type RoomRedisService struct {
	client *redis.Client
}

func NewRoomRedisService(client *redis.Client) *RoomRedisService {
	return &RoomRedisService{
		client: client,
	}
}

// CreateRoom stores basic room information in Redis
func (r *RoomRedisService) CreateRoom(
	ctx context.Context,
	roomID, hostID uuid.UUID,
	src string,
	capacity int,
) error {
	// Prepare Redis keys
	roomKey := fmt.Sprintf("%s%s", roomPrefix, roomID.String())
	usersSetKey := fmt.Sprintf("%s%s", roomUsersSetPrefix, roomID.String())

	pipe := r.client.TxPipeline()

	// Store room metadata
	pipe.HSet(ctx, roomKey, map[string]any{
		"host_id":  hostID.String(),
		"src":      src,
		"capacity": capacity,
	})

	// Set room expiration
	pipe.Expire(ctx, roomKey, defaultRoomTTL)

	// Create an empty users set
	pipe.Del(ctx, usersSetKey)

	// Execute transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create room in Redis: %w", err)
	}

	return nil
}

// AddUser adds a user to a room's user set
func (r *RoomRedisService) AddUser(
	ctx context.Context,
	roomID, userID uuid.UUID,
) error {
	// Prepare Redis keys
	usersSetKey := fmt.Sprintf("%s%s", roomUsersSetPrefix, roomID.String())

	// Add user to room's user set
	err := r.client.SAdd(ctx, usersSetKey, userID.String()).Err()
	if err != nil {
		return fmt.Errorf("failed to add user to room: %w", err)
	}

	return nil
}

// RemoveUser removes a user from a room's user set
func (r *RoomRedisService) RemoveUser(
	ctx context.Context,
	roomID, userID uuid.UUID,
) error {
	// Prepare Redis keys
	usersSetKey := fmt.Sprintf("%s%s", roomUsersSetPrefix, roomID.String())

	// Remove user from room's user set
	err := r.client.SRem(ctx, usersSetKey, userID.String()).Err()
	if err != nil {
		return fmt.Errorf("failed to remove user from room: %w", err)
	}

	return nil
}

// GetRoomUsers retrieves all users in a room
func (r *RoomRedisService) GetRoomUsers(
	ctx context.Context,
	roomID uuid.UUID,
) ([]string, error) {
	// Prepare Redis key
	usersSetKey := fmt.Sprintf("%s%s", roomUsersSetPrefix, roomID.String())

	// Get all users in the room
	users, err := r.client.SMembers(ctx, usersSetKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get room users: %w", err)
	}

	return users, nil
}

// GetRoomCapacity retrieves the room's capacity
func (r *RoomRedisService) GetRoomCapacity(
	ctx context.Context,
	roomID uuid.UUID,
) (int, error) {
	// Prepare Redis key
	roomKey := fmt.Sprintf("%s%s", roomPrefix, roomID.String())

	// Get capacity from room metadata
	capacity, err := r.client.HGet(ctx, roomKey, "capacity").Int()
	if err != nil {
		return 0, fmt.Errorf("failed to get room capacity: %w", err)
	}

	return capacity, nil
}

// CountRoomUsers counts the number of users in a room
func (r *RoomRedisService) CountRoomUsers(
	ctx context.Context,
	roomID uuid.UUID,
) (int64, error) {
	// Prepare Redis key
	usersSetKey := fmt.Sprintf("%s%s", roomUsersSetPrefix, roomID.String())

	// Count users in the room
	count, err := r.client.SCard(ctx, usersSetKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to count room users: %w", err)
	}

	return count, nil
}

// DeleteRoom removes all room-related data from Redis
func (r *RoomRedisService) DeleteRoom(
	ctx context.Context,
	roomID uuid.UUID,
) error {
	// Prepare Redis keys
	roomKey := fmt.Sprintf("%s%s", roomPrefix, roomID.String())
	usersSetKey := fmt.Sprintf("%s%s", roomUsersSetPrefix, roomID.String())

	// Use Redis transaction for atomic operations
	pipe := r.client.TxPipeline()

	// Delete room metadata
	pipe.Del(ctx, roomKey)
	// Delete users set
	pipe.Del(ctx, usersSetKey)

	// Execute transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	return nil
}
