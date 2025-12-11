package services

import (
	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/dliluashvili/cowatchit/internal/repositories"
	"github.com/google/uuid"
)

type RoomMessageService struct {
	roomMessageRepository *repositories.RoomMessageRepository
}

func NewRoomMessageService(rmp *repositories.RoomMessageRepository) *RoomMessageService {
	return &RoomMessageService{
		roomMessageRepository: rmp,
	}
}

func (roomMessageService *RoomMessageService) Create(createRoomMessageDto *dtos.CreateRoomMessageDto) (*models.RoomMessage, error) {
	return roomMessageService.roomMessageRepository.Create(createRoomMessageDto)
}

func (roomMessageService *RoomMessageService) GetRoomMessages(roomID uuid.UUID) ([]*models.RoomMessage, error) {
	return roomMessageService.roomMessageRepository.FindByRoom(roomID)
}
