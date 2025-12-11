package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/dliluashvili/cowatchit/internal/repositories"
	"github.com/google/uuid"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrRoomFull     = errors.New("room is full")
)

type RoomService struct {
	roomRepository   *repositories.RoomRepository
	userService      *UserService
	roomRedisService *RoomRedisService
}

func NewRoomService(rp *repositories.RoomRepository, us *UserService, rrs *RoomRedisService) *RoomService {
	return &RoomService{
		roomRepository:   rp,
		userService:      us,
		roomRedisService: rrs,
	}
}

func (rs *RoomService) Create(ctx context.Context, dto *dtos.CreateRoomServiceDto) (*models.Room, error) {
	host, err := rs.userService.FindByID(dto.HostID.String())

	if err != nil {
		return nil, err
	}

	repoDto := &dtos.CreateRoomRepoDto{
		HostUsername:         host.Username,
		CreateRoomServiceDto: dto,
	}

	// Create room in Postgres
	room, err := rs.roomRepository.Create(repoDto)

	if err != nil {
		return nil, fmt.Errorf("error creating room in db: %w", err)
	}

	// Create room in Redis
	err = rs.roomRedisService.CreateRoom(
		ctx,
		room.ID,
		room.HostID,
		dto.Src,
		dto.Capacity,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating room in redis: %w", err)
	}

	// Add host to the room
	err = rs.joinRoom(ctx, room.ID, room.HostID)
	if err != nil {
		return nil, fmt.Errorf("error adding host to room: %w", err)
	}

	return room, nil
}

func (rs *RoomService) Find(dto *dtos.FindRoomDto) ([]models.Room, error) {
	return rs.roomRepository.Find(dto)
}

func (rs *RoomService) FindOne(ID uuid.UUID) (*models.Room, error) {
	return rs.roomRepository.FindOne(ID)
}

func (rs *RoomService) Exists(ID uuid.UUID) (bool, error) {
	return rs.roomRepository.Exists(ID)
}

func (rs *RoomService) Join(ctx context.Context, roomID, userID *uuid.UUID) error {
	return rs.joinRoom(ctx, *roomID, *userID)
}

// Internal method to handle room joining logic
func (rs *RoomService) joinRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	// Check if room exists
	findDto := &dtos.FindRoomDto{ID: &roomID}

	rooms, err := rs.roomRepository.Find(findDto)

	if err != nil {
		return fmt.Errorf("error finding room: %w", err)
	}
	if len(rooms) == 0 {
		return ErrRoomNotFound
	}

	// Check room capacity in Redis
	capacity, err := rs.roomRedisService.GetRoomCapacity(ctx, roomID)
	if err != nil {
		return fmt.Errorf("error getting room capacity: %w", err)
	}

	// Count current users
	currentUsers, err := rs.roomRedisService.CountRoomUsers(ctx, roomID)
	if err != nil {
		return fmt.Errorf("error counting room users: %w", err)
	}

	// Check if room is full
	if currentUsers >= int64(capacity) {
		return ErrRoomFull
	}

	// Add user to Redis
	err = rs.roomRedisService.AddUser(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("error adding user to room in Redis: %w", err)
	}

	// Add user to room in Postgres
	// err = rs.roomRepository.AddUserToRoom(roomID, userID)
	// if err != nil {
	// Rollback Redis addition if Postgres fails
	// 	_ = rs.roomRedisService.RemoveUser(ctx, roomID, userID)
	// 	return fmt.Errorf("error adding user to room in Postgres: %w", err)
	// }

	return nil
}
