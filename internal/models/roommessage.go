package models

import (
	"time"

	"github.com/google/uuid"
)

type RoomMessage struct {
	ID             uuid.UUID `json:"id"`
	SenderID       uuid.UUID `json:"sender_id"`
	SenderUsername string    `json:"sender_username"`
	RoomID         uuid.UUID `json:"room_id"`
	Content        string    `json:"content"`
	IsHost         bool      `json:"is_host"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
