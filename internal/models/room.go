package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID           uuid.UUID `json:"id"`
	HostID       uuid.UUID `json:"user_id"`
	HostUsername string    `json:"host_username"`
	Title        string    `json:"title"`
	Capacity     int       `json:"capacity"`
	Description  string    `json:"description"`
	Src          string    `json:"src"`
	Poster       string    `json:"poster"`
	Private      bool      `json:"private"`
	Password     string    `json:"password"`
	Hidden       bool      `json:"hidden"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
