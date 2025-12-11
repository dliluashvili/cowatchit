package types

import (
	"encoding/json"

	"github.com/coder/websocket"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/google/uuid"
)

const (
	TypeEvent = "EVENT"
	TypeError = "ERROR"
)

const (
	StateStop    = "STOP"
	StatePaused  = "PAUSED"
	StatePlaying = "PLAYING"
	StateEnd     = "END"
)

const (
	EventHostStateSend       = "HOST_STATE_SEND"
	EventHostStateReceived   = "HOST_STATE_RECEIVED"
	EventUserStateSend       = "USER_STATE_SEND"
	EventUserStateReceived   = "USER_STATE_RECEIVED"
	EventVideoPaused         = "VIDEO_PAUSED"
	EventVideoPlaying        = "VIDEO_PLAYING"
	EventChatMessageSend     = "CHAT_MESSAGE_SEND"
	EventChatMessageReceived = "CHAT_MESSAGE_RECEIVED"
	EventUserJoinRequest     = "USER_JOIN_REQUEST"
	EventUserJoinAnswer      = "USER_JOIN_ANSWER"
	EventRoomMessagesRequest = "ROOM_MESSAGES_REQUEST"
	EventRoomMessagesAnswer  = "ROOM_MESSAGES_ANSWER"
	EventUserJoint           = "USER_JOINT"
	EventUserLeft            = "USER_LEFT"
	EventIdentify            = "IDENTIFY"
)

type WSMessage struct {
	Type  string          `json:"type"`
	Event string          `json:"event,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type WebSocketContext struct {
	ID     string // socketId
	User   *models.User
	RoomID uuid.UUID // roomID
	IsHost bool
	Conn   *websocket.Conn // WebSocket connection
}

type UserInfo struct {
	Username string `json:"username"`
	IsHost   bool   `json:"is_host"`
}

type UserIDInfo = map[uuid.UUID]*UserInfo

type RoomMetadata struct {
	RoomID             uuid.UUID
	Capacity           int
	State              string
	CurrentTimeSeconds float64
	HostUsername       string
	HostID             uuid.UUID
	SocketIDs          map[string]bool // set of socket IDs
	Users              UserIDInfo      // set of user IDs with their info
}
