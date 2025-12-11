package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/helpers"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/dliluashvili/cowatchit/internal/services"
	"github.com/dliluashvili/cowatchit/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type WebSocketHandler struct {
	validate                *validator.Validate
	websocketManagerService *services.WebSocketManagerService
	sessionService          *services.SessionService
	roomService             *services.RoomService
	roomMessageService      *services.RoomMessageService
}

func NewWebSocketHandler(
	validate *validator.Validate,
	websocketManagerService *services.WebSocketManagerService,
	sessionService *services.SessionService,
	roomService *services.RoomService,
	roomMessageService *services.RoomMessageService,
) *WebSocketHandler {
	return &WebSocketHandler{
		validate:                validate,
		websocketManagerService: websocketManagerService,
		sessionService:          sessionService,
		roomService:             roomService,
		roomMessageService:      roomMessageService,
	}
}

func (h *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})

	if err != nil {
		log.Println("WebSocket accept error:", err)
		return
	}

	defer conn.Close(websocket.StatusInternalError, "internal error")

	ctx := r.Context()

	// Authenticate session
	sessionModel, err := h.authenticateSession(r)
	if err != nil {
		log.Println("Authentication error:", err)
		conn.Close(websocket.StatusPolicyViolation, "authentication failed")
		return
	}

	log.Println("User authenticated:", sessionModel.User.Username)

	// Create WebSocket context
	socketID := helpers.GenerateSocketID()

	wsCtx := &types.WebSocketContext{
		ID:     socketID,
		User:   sessionModel.User,
		RoomID: uuid.Nil, // Will be set when user joins a room
		Conn:   conn,
	}

	// Send socketId to client immediately (no need to wait for IDENTIFY request)
	identifyData := struct {
		SocketID     string    `json:"socket_id"`
		AuthID       uuid.UUID `json:"auth_id"`
		AuthUsername string    `json:"auth_username"`
	}{
		SocketID:     socketID,
		AuthID:       sessionModel.User.ID,
		AuthUsername: sessionModel.User.Username,
	}

	rawData, _ := json.Marshal(identifyData)

	identifyMsg := types.WSMessage{
		Type:  types.TypeEvent,
		Event: types.EventIdentify,
		Data:  rawData,
	}

	data, _ := json.Marshal(identifyMsg)

	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		log.Println("Failed to send socket ID:", err)
		return
	}

	log.Printf("Socket ID sent to user %s: %s", sessionModel.User.Username, socketID)

	// Message handling loop
	h.handleMessageLoop(ctx, conn, sessionModel, wsCtx)
}

func (h *WebSocketHandler) authenticateSession(r *http.Request) (*models.Session, error) {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		return nil, fmt.Errorf("missing session cookie: %w", err)
	}

	return h.sessionService.GetUserBySession(r.Context(), cookie.Value)
}

func (h *WebSocketHandler) handleMessageLoop(
	ctx context.Context,
	conn *websocket.Conn,
	sessionModel *models.Session,
	wsCtx *types.WebSocketContext,
) {
	defer func() {
		h.websocketManagerService.Unregister(wsCtx.ID)

		participants, err := h.websocketManagerService.GetRoomParticipants(wsCtx.RoomID)

		if err != nil {
			log.Printf("error particpant: %v", err)
			return
		}

		userLeftData := struct {
			IsHost              bool   `json:"is_host"`
			UserID              string `json:"user_id"`
			Username            string `json:"username"`
			SocketID            string `json:"socket_id"`
			CountedParticipants int    `json:"counted_participants"`
		}{
			IsHost:              wsCtx.IsHost,
			UserID:              sessionModel.User.ID.String(),
			Username:            sessionModel.User.Username,
			SocketID:            wsCtx.ID,
			CountedParticipants: len(*participants),
		}

		broadcastRawData, _ := json.Marshal(userLeftData)

		broadcastMsg := types.WSMessage{
			Type:  types.TypeEvent,
			Event: types.EventUserLeft,
			Data:  broadcastRawData,
		}

		messageDataJSON, _ := json.Marshal(broadcastMsg)
		if err := h.websocketManagerService.BroadcastToRoom(ctx, wsCtx.RoomID, messageDataJSON, wsCtx.ID); err != nil {
			log.Printf("Broadcast message error: %v", err)
			return
		}

	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgType, data, err := conn.Read(ctx)
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}

			// Only handle text messages
			if msgType != websocket.MessageText {
				continue
			}

			var msg types.WSMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				h.sendWSError(ctx, conn, "invalid message format")
				continue
			}

			// Process room actions
			if err := h.handleRoomAction(ctx, conn, sessionModel, wsCtx, msg); err != nil {
				log.Printf("Room action error: %v", err)
				h.sendWSError(ctx, conn, "failed to process room action")
			}
		}
	}
}

func (h *WebSocketHandler) handleRoomAction(
	ctx context.Context,
	conn *websocket.Conn,
	sessionModel *models.Session,
	wsCtx *types.WebSocketContext,
	msg types.WSMessage,
) error {
	switch msg.Event {
	case types.EventUserJoinRequest:
		return h.handleRoomJoin(ctx, conn, sessionModel, wsCtx, msg.Data)
	case types.EventUserLeft:
		return h.handleRoomLeave(ctx, conn, sessionModel, wsCtx, msg.Data)
	case types.EventChatMessageSend:
		return h.handleRoomMessage(ctx, conn, sessionModel, wsCtx, msg.Data)
	case types.EventRoomMessagesRequest:
		return h.handleRoomMessages(ctx, conn, sessionModel, wsCtx, msg.Data)
	case types.EventHostStateSend:
		return h.handleHostStateChange(ctx, conn, sessionModel, wsCtx, msg.Data)
	default:
		return fmt.Errorf("unknown room action: %s", msg.Event)
	}
}

func (h *WebSocketHandler) handleRoomJoin(
	ctx context.Context,
	conn *websocket.Conn,
	sessionModel *models.Session,
	wsCtx *types.WebSocketContext,
	payload json.RawMessage,
) error {
	log.Println("User is joining room")

	type JoinPayload struct {
		RoomID string `json:"room_id" validate:"required,uuid"`
	}

	var joinPayload JoinPayload
	if err := json.Unmarshal(payload, &joinPayload); err != nil {
		return h.sendWSError(ctx, conn, "invalid join payload")
	}

	if err := h.validate.Struct(joinPayload); err != nil {
		return h.sendWSError(ctx, conn, "validation failed")
	}

	// Parse room ID
	roomID, err := uuid.Parse(joinPayload.RoomID)

	if err != nil {
		return h.sendWSError(ctx, conn, "invalid room ID")
	}

	// Get room details
	room, err := h.roomService.FindOne(roomID)

	if err != nil {
		return h.sendWSError(ctx, conn, "room not found")
	}

	// Update WebSocket context with room

	isHost := room.HostID == sessionModel.User.ID

	wsCtx.IsHost = isHost
	wsCtx.RoomID = roomID

	// Register connection in WebSocket manager

	if err := h.websocketManagerService.Register(wsCtx, room); err != nil {
		return h.sendWSError(ctx, conn, err.Error())
	}

	participants, err := h.websocketManagerService.GetRoomParticipants(roomID)

	if err != nil {
		return h.sendWSError(ctx, conn, "Bad error")
	}

	roomMetadata := h.websocketManagerService.GetRoomMetadata(roomID)

	if roomMetadata == nil {
		fmt.Println("room doesnt exist ;)")
		return h.sendWSError(ctx, conn, "Bad error")
	}

	// Send join confirmation to user
	joinData := struct {
		Title              string            `json:"title"`
		Host               string            `json:"host"`
		IsHost             bool              `json:"is_host"`
		Src                string            `json:"src"`
		State              string            `json:"state"`
		CurrentTimeSeconds float64           `json:"current_time_seconds"`
		Participants       *types.UserIDInfo `json:"participants"`
	}{
		Title:              room.Title,
		Host:               room.HostUsername,
		IsHost:             isHost,
		Src:                room.Src,
		State:              roomMetadata.State,
		CurrentTimeSeconds: roomMetadata.CurrentTimeSeconds,
		Participants:       participants,
	}

	rawData, _ := json.Marshal(joinData)

	joinMsg := types.WSMessage{
		Type:  types.TypeEvent,
		Event: types.EventUserJoinAnswer,
		Data:  rawData,
	}
	confirmationData, _ := json.Marshal(joinMsg)

	if err := conn.Write(ctx, websocket.MessageText, confirmationData); err != nil {
		h.websocketManagerService.Unregister(wsCtx.ID)
		return err
	}

	log.Printf("User %s joined room %s", sessionModel.User.Username, roomID)

	userJointData := struct {
		IsHost              bool   `json:"is_host"`
		UserID              string `json:"user_id"`
		Username            string `json:"username"`
		CountedParticipants int    `json:"counted_participants"`
	}{
		IsHost:              wsCtx.IsHost,
		UserID:              sessionModel.User.ID.String(),
		Username:            sessionModel.User.Username,
		CountedParticipants: len(*participants),
	}

	broadcastRawData, _ := json.Marshal(userJointData)

	broadcastMsg := types.WSMessage{
		Type:  types.TypeEvent,
		Event: types.EventUserJoint,
		Data:  broadcastRawData,
	}

	messageDataJSON, _ := json.Marshal(broadcastMsg)
	if err := h.websocketManagerService.BroadcastToRoom(ctx, wsCtx.RoomID, messageDataJSON, wsCtx.ID); err != nil {
		log.Printf("Broadcast message error: %v", err)
		return h.sendWSError(ctx, conn, "failed to send message")
	}

	return nil
}

func (h *WebSocketHandler) handleRoomMessage(
	ctx context.Context,
	conn *websocket.Conn,
	sessionModel *models.Session,
	wsCtx *types.WebSocketContext,
	payload json.RawMessage,
) error {
	type MessagePayload struct {
		RoomId   string `json:"room_id" validate:"required,uuid"`
		SocketID string `json:"socket_id" `
		Content  string `json:"content" validate:"required,max=250"`
	}

	var msgPayload MessagePayload

	if err := json.Unmarshal(payload, &msgPayload); err != nil {
		return h.sendWSError(ctx, conn, "invalid message payload")
	}

	if err := h.validate.Struct(msgPayload); err != nil {
		return h.sendWSError(ctx, conn, "validation failed")
	}

	roomID, err := uuid.Parse(msgPayload.RoomId)

	if err != nil {
		fmt.Println("err", err)
		return h.sendWSError(ctx, conn, "validation failed")
	}

	createMessageRoomDto := &dtos.CreateRoomMessageDto{
		Content:        msgPayload.Content,
		RoomID:         roomID,
		SenderID:       sessionModel.User.ID,
		IsHost:         wsCtx.IsHost,
		SenderUsername: sessionModel.User.Username,
	}

	isUserAllowed := h.websocketManagerService.IsUserAllowed(roomID, createMessageRoomDto.SenderID, wsCtx.ID)

	if !isUserAllowed {
		fmt.Println("not allowed room socket")

		return h.sendWSError(ctx, conn, "bad request")
	}

	roomMessageService, err := h.roomMessageService.Create(createMessageRoomDto)

	if err != nil {
		fmt.Println("err", err)
		return h.sendWSError(ctx, conn, "very bad error :X")
	}

	messageData := struct {
		SenderID       uuid.UUID `json:"sender_id"`
		SenderUsername string    `json:"sender_username"`
		Content        string    `json:"content"`
		IsHost         bool      `json:"is_host"`
		CreatedAt      string    `json:"created_at"`
	}{
		SenderID:       roomMessageService.SenderID,
		SenderUsername: roomMessageService.SenderUsername,
		Content:        roomMessageService.Content,
		IsHost:         roomMessageService.IsHost,
		CreatedAt:      roomMessageService.CreatedAt.String(),
	}

	rawData, _ := json.Marshal(messageData)

	messageMsg := types.WSMessage{
		Type:  types.TypeEvent,
		Event: types.EventChatMessageReceived,
		Data:  rawData,
	}
	messageDataJSON, _ := json.Marshal(messageMsg)

	if err := h.websocketManagerService.BroadcastToRoom(ctx, wsCtx.RoomID, messageDataJSON, wsCtx.ID); err != nil {
		log.Printf("Broadcast message error: %v", err)
		return h.sendWSError(ctx, conn, "failed to send message")
	}

	return nil
}

func (h *WebSocketHandler) handleHostStateChange(
	ctx context.Context,
	conn *websocket.Conn,
	sessionModel *models.Session,
	wsCtx *types.WebSocketContext,
	payload json.RawMessage,
) error {

	if !wsCtx.IsHost {
		return h.sendWSError(ctx, conn, "invalid request broooo")
	}

	type MessagePayload struct {
		State              string  `json:"state"`
		CurrentTimeSeconds float64 `json:"current_time_seconds"`
		RoomId             string  `json:"room_id" validate:"required,uuid"`
		SocketID           string  `json:"socket_id" `
	}

	var msgPayload MessagePayload

	if err := json.Unmarshal(payload, &msgPayload); err != nil {
		return h.sendWSError(ctx, conn, "invalid message payload")
	}

	if err := h.validate.Struct(msgPayload); err != nil {
		return h.sendWSError(ctx, conn, "validation failed")
	}

	roomID, err := uuid.Parse(msgPayload.RoomId)

	if err != nil {
		fmt.Println("err", err)
		return h.sendWSError(ctx, conn, "validation failed")
	}

	userID := sessionModel.User.ID

	isUserAllowed := h.websocketManagerService.IsUserAllowed(roomID, userID, wsCtx.ID)

	if !isUserAllowed {
		fmt.Println("not allowed room socket")

		return h.sendWSError(ctx, conn, "bad request")
	}

	messageData := struct {
		State              string  `json:"state"`
		CurrentTimeSeconds float64 `json:"current_time_seconds"`
	}{
		State:              msgPayload.State,
		CurrentTimeSeconds: msgPayload.CurrentTimeSeconds,
	}

	rawData, _ := json.Marshal(messageData)

	messageMsg := types.WSMessage{
		Type:  types.TypeEvent,
		Event: types.EventHostStateReceived,
		Data:  rawData,
	}

	broadcastData, _ := json.Marshal(messageMsg)

	if err := h.websocketManagerService.BroadcastToRoom(ctx, wsCtx.RoomID, broadcastData, wsCtx.ID); err != nil {
		log.Printf("Broadcast leave error: %v", err)
	}

	return nil
}

func (h *WebSocketHandler) handleRoomMessages(
	ctx context.Context,
	conn *websocket.Conn,
	sessionModel *models.Session,
	wsCtx *types.WebSocketContext,
	payload json.RawMessage,
) error {
	type MessagePayload struct {
		RoomId   string `json:"room_id" validate:"required,uuid"`
		SocketID string `json:"socket_id" `
	}

	var msgPayload MessagePayload

	if err := json.Unmarshal(payload, &msgPayload); err != nil {
		return h.sendWSError(ctx, conn, "invalid message payload")
	}

	if err := h.validate.Struct(msgPayload); err != nil {
		return h.sendWSError(ctx, conn, "validation failed")
	}

	roomID, err := uuid.Parse(msgPayload.RoomId)

	if err != nil {
		fmt.Println("err", err)
		return h.sendWSError(ctx, conn, "validation failed")
	}

	userID := sessionModel.User.ID

	isUserAllowed := h.websocketManagerService.IsUserAllowed(roomID, userID, wsCtx.ID)

	if !isUserAllowed {
		fmt.Println("not allowed room socket")

		return h.sendWSError(ctx, conn, "bad request")
	}

	roomMessages, err := h.roomMessageService.GetRoomMessages(roomID)

	if err != nil {
		fmt.Println("err", err)
		return h.sendWSError(ctx, conn, "very bad error :X")
	}

	messageData := struct {
		Messages []*models.RoomMessage `json:"messages"`
		AuthID   uuid.UUID             `json:"auth_id"`
	}{
		Messages: roomMessages,
		AuthID:   userID,
	}

	rawData, _ := json.Marshal(messageData)

	messageMsg := types.WSMessage{
		Type:  types.TypeEvent,
		Event: types.EventRoomMessagesAnswer,
		Data:  rawData,
	}
	messageDataJSON, _ := json.Marshal(messageMsg)

	if err := conn.Write(ctx, websocket.MessageText, messageDataJSON); err != nil {
		log.Println("Failed to send socket ID:", err)
		return nil
	}

	return nil
}

func (h *WebSocketHandler) handleRoomLeave(
	ctx context.Context,
	conn *websocket.Conn,
	sessionModel *models.Session,
	wsCtx *types.WebSocketContext,
	payload json.RawMessage,
) error {
	fmt.Println("user left")

	if wsCtx.RoomID == uuid.Nil {
		return h.sendWSError(ctx, conn, "not in a room")
	}

	roomID := wsCtx.RoomID

	// Send leave confirmation to user
	leaveData := struct {
		RoomID   string `json:"room_id"`
		SocketID string `json:"socket_id"`
		Message  string `json:"message"`
	}{
		RoomID:   roomID.String(),
		SocketID: wsCtx.ID,
		Message:  "Successfully left the room",
	}

	rawData, _ := json.Marshal(leaveData)
	leaveMsg := types.WSMessage{
		Type:  types.TypeEvent,
		Event: types.EventUserLeft,
		Data:  rawData,
	}
	confirmationData, _ := json.Marshal(leaveMsg)
	if err := conn.Write(ctx, websocket.MessageText, confirmationData); err != nil {
		return err
	}

	log.Printf("User %s left room %s", sessionModel.User.Username, roomID)

	// Broadcast user left to room
	userLeftData := struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		SocketID string `json:"socket_id"`
	}{
		UserID:   sessionModel.User.ID.String(),
		Username: sessionModel.User.Username,
		SocketID: wsCtx.ID,
	}

	broadcastRawData, _ := json.Marshal(userLeftData)
	broadcastMsg := types.WSMessage{
		Type:  types.TypeEvent,
		Event: types.EventUserLeft,
		Data:  broadcastRawData,
	}
	broadcastData, _ := json.Marshal(broadcastMsg)
	if err := h.websocketManagerService.BroadcastToRoom(ctx, roomID, broadcastData, wsCtx.ID); err != nil {
		log.Printf("Broadcast leave error: %v", err)
	}

	// Unregister from WebSocket manager
	h.websocketManagerService.Unregister(wsCtx.ID)

	return nil
}

func (h *WebSocketHandler) sendWSError(ctx context.Context, conn *websocket.Conn, errMsg string) error {
	errorData := struct {
		Message string `json:"message"`
	}{
		Message: errMsg,
	}

	rawData, _ := json.Marshal(errorData)
	errorMsg := types.WSMessage{
		Type:  types.TypeError,
		Event: types.TypeError,
		Data:  rawData,
	}

	data, _ := json.Marshal(errorMsg)
	return conn.Write(ctx, websocket.MessageText, data)
}
