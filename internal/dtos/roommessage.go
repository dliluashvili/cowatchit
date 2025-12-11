package dtos

import "github.com/google/uuid"

type CreateRoomMessageDto struct {
	Content        string    `json:"content"`
	RoomID         uuid.UUID `json:"room_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	SenderUsername string    `json:"sender_username"`
	IsHost         bool      `json:"is_host"`
}
