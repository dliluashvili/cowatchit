package repositories

import (
	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomMessageRepository struct {
	db *gorm.DB
}

func NewRoomMessageRepository(db *gorm.DB) *RoomMessageRepository {
	return &RoomMessageRepository{
		db: db,
	}
}

func (rmp *RoomMessageRepository) Create(createRoomMessageDto *dtos.CreateRoomMessageDto) (*models.RoomMessage, error) {
	roomMessage := &models.RoomMessage{
		ID:             uuid.New(),
		Content:        createRoomMessageDto.Content,
		RoomID:         createRoomMessageDto.RoomID,
		SenderID:       createRoomMessageDto.SenderID,
		SenderUsername: createRoomMessageDto.SenderUsername,
		IsHost:         createRoomMessageDto.IsHost,
	}

	result := rmp.db.Create(roomMessage)

	if result.Error != nil {
		return nil, result.Error
	}

	return roomMessage, nil
}

func (rmp *RoomMessageRepository) FindByRoom(roomID uuid.UUID) ([]*models.RoomMessage, error) {
	var roomMessages []*models.RoomMessage

	result := rmp.db.Where("room_id = ?", roomID).Find(&roomMessages)

	if result.Error != nil {
		return nil, result.Error
	}

	return roomMessages, nil
}
