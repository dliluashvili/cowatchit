package models

import (
	"time"

	"github.com/google/uuid"
)

type RoomRepository struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	RoomID    uuid.UUID `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
