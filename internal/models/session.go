package models

import (
	"time"
)

type Session struct {
	User      *User     `json:"user"`
	SessionID string    `json:"session_id"`
	SocketID  *string   `json:"socket_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
