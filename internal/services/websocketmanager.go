package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/coder/websocket"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/dliluashvili/cowatchit/internal/types"
	"github.com/google/uuid"
)

const MaxConnectionsPerUser = 5

type WebSocketManagerService struct {
	mu sync.RWMutex

	// socketId -> userId
	socketIdToUserId map[string]uuid.UUID

	// userId -> {socketIds}
	userIDSocketIDs map[uuid.UUID]map[string]bool

	// userId_socketId -> conn
	userSocketConnections map[string]*websocket.Conn

	// roomId -> RoomMetadata
	roomMetadata map[uuid.UUID]*types.RoomMetadata
}

func NewWebSocketManagerService() *WebSocketManagerService {
	return &WebSocketManagerService{
		socketIdToUserId:      make(map[string]uuid.UUID),
		userIDSocketIDs:       make(map[uuid.UUID]map[string]bool),
		userSocketConnections: make(map[string]*websocket.Conn),
		roomMetadata:          make(map[uuid.UUID]*types.RoomMetadata),
	}
}

// Register a new socket
func (sm *WebSocketManagerService) Register(ctx *types.WebSocketContext, room *models.Room) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Check if user already has max connections
	userSocketCount := len(sm.userIDSocketIDs[ctx.User.ID])

	if userSocketCount >= MaxConnectionsPerUser {
		return fmt.Errorf("user has reached maximum connections (%d)", MaxConnectionsPerUser)
	}

	// If user is joining a different room, kick them from previous room
	if socketIDs, exists := sm.userIDSocketIDs[ctx.User.ID]; exists && len(socketIDs) > 0 {
		// Get first socket ID for room lookup
		var firstSocketID string
		for sockID := range socketIDs {
			firstSocketID = sockID
			break
		}
		oldRoomID := sm.getRoomIDFromSocket(firstSocketID)
		if oldRoomID != ctx.RoomID {
			sm.kickUserFromRoom(ctx.User.ID)
		}
	}

	// Create or get room metadata
	roomMeta, roomExists := sm.roomMetadata[ctx.RoomID]

	if !roomExists {
		if ctx.IsHost {
			// Init
			roomMeta = &types.RoomMetadata{
				RoomID:             ctx.RoomID,
				Capacity:           room.Capacity,
				CurrentTimeSeconds: 0,
				State:              types.StateStop,
				HostID:             room.HostID,
				SocketIDs:          make(map[string]bool),
				Users:              make(types.UserIDInfo),
			}
			sm.roomMetadata[ctx.RoomID] = roomMeta
		} else {
			return fmt.Errorf("room doesnt exist")
		}
	}

	// Check room capacity
	if len(roomMeta.Users) >= roomMeta.Capacity {
		return fmt.Errorf("room has reached maximum capacity (%d)", roomMeta.Capacity)
	}

	// Add socket to room
	roomMeta.SocketIDs[ctx.ID] = true

	roomMeta.Users[ctx.User.ID] = &types.UserInfo{
		Username: ctx.User.Username,
		IsHost:   ctx.IsHost,
	}

	// Add mappings
	sm.socketIdToUserId[ctx.ID] = ctx.User.ID

	// Initialize user's socket map if needed
	if _, exists := sm.userIDSocketIDs[ctx.User.ID]; !exists {
		sm.userIDSocketIDs[ctx.User.ID] = make(map[string]bool)
	}
	sm.userIDSocketIDs[ctx.User.ID][ctx.ID] = true

	compositeKey := sm.getCompositeKey(ctx.User.ID, ctx.ID)
	sm.userSocketConnections[compositeKey] = ctx.Conn

	return nil
}

// Unregister a socket
func (sm *WebSocketManagerService) Unregister(socketID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Get userID from socketID
	userID, userExists := sm.socketIdToUserId[socketID]

	if !userExists {
		return
	}

	// Get roomID from room metadata
	roomID := sm.getRoomIDFromSocket(socketID)

	// Remove socket from room
	if roomMeta, roomMetaExists := sm.roomMetadata[roomID]; roomMetaExists {
		// Delete socket from set
		delete(roomMeta.SocketIDs, socketID)

		// Check if user has any other sockets in room
		userHasOtherSockets := false
		for sockID := range roomMeta.SocketIDs {
			if sm.socketIdToUserId[sockID] == userID {
				userHasOtherSockets = true
				break
			}
		}

		// Remove user from room if no other sockets
		if !userHasOtherSockets {
			delete(roomMeta.Users, userID)
		}

		// Remove room if empty
		if len(roomMeta.SocketIDs) == 0 {
			delete(sm.roomMetadata, roomID)
		}
	}

	// Remove from userIDSocketIDs
	if socketIDs, exists := sm.userIDSocketIDs[userID]; exists {
		delete(socketIDs, socketID)
		if len(socketIDs) == 0 {
			delete(sm.userIDSocketIDs, userID)
		}
	}

	// Remove mappings
	delete(sm.socketIdToUserId, socketID)
	compositeKey := sm.getCompositeKey(userID, socketID)
	delete(sm.userSocketConnections, compositeKey)
}

// Get room ID from socket (helper function)
func (sm *WebSocketManagerService) getRoomIDFromSocket(socketID string) uuid.UUID {
	for roomID, roomMeta := range sm.roomMetadata {
		if roomMeta.SocketIDs[socketID] {
			return roomID
		}
	}
	return uuid.UUID{}
}

func (sm *WebSocketManagerService) GetRoomMetadata(roomID uuid.UUID) *types.RoomMetadata {
	roomMeta, exists := sm.roomMetadata[roomID]

	if !exists {
		return nil
	}

	return roomMeta
}

func (sm *WebSocketManagerService) IsUserAllowed(roomID, userID uuid.UUID, socketId string) bool {

	roomMetadata := sm.GetRoomMetadata(roomID)

	if roomMetadata == nil {
		fmt.Println("room doesnt exist")

		return false
	}

	_, userExistsInRoom := roomMetadata.Users[userID]

	if !userExistsInRoom {
		fmt.Println("not allowed room userid")

		return false
	}

	_, socketExistsInRoom := roomMetadata.SocketIDs[socketId]

	if !socketExistsInRoom {
		fmt.Println("not allowed room socket")

		return false
	}

	return true
}

// Kick user from room (close all their connections)
func (sm *WebSocketManagerService) kickUserFromRoom(userID uuid.UUID) {
	socketIDs, exists := sm.userIDSocketIDs[userID]
	if !exists {
		return
	}

	// Close all connections for this user
	for socketID := range socketIDs {
		if conn, exists := sm.userSocketConnections[sm.getCompositeKey(userID, socketID)]; exists {
			conn.Close(websocket.StatusGoingAway, "user joined another room")
		}
	}
}

// Get composite key
func (sm *WebSocketManagerService) getCompositeKey(userID uuid.UUID, socketID string) string {
	return fmt.Sprintf("%s_%s", userID.String(), socketID)
}

// Get connection from socketID
func (sm *WebSocketManagerService) GetConnFromSocket(socketID string) (*websocket.Conn, error) {
	sm.mu.RLock()
	userID, exists := sm.socketIdToUserId[socketID]
	if !exists {
		sm.mu.RUnlock()
		return nil, fmt.Errorf("socket not found: %s", socketID)
	}

	compositeKey := sm.getCompositeKey(userID, socketID)
	conn, connExists := sm.userSocketConnections[compositeKey]
	sm.mu.RUnlock()

	if !connExists {
		return nil, fmt.Errorf("connection not found for socket: %s", socketID)
	}

	return conn, nil
}

// Count room participants
func (sm *WebSocketManagerService) GetRoomParticipants(roomID uuid.UUID) (*types.UserIDInfo, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	roomMeta, exists := sm.roomMetadata[roomID]

	if !exists {
		return nil, fmt.Errorf("room not found: %s", roomID)
	}

	return &roomMeta.Users, nil
}

// Get all connections for user
func (sm *WebSocketManagerService) GetUserConnections(userID uuid.UUID) ([]*websocket.Conn, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	socketIDs, exists := sm.userIDSocketIDs[userID]
	if !exists || len(socketIDs) == 0 {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	connections := make([]*websocket.Conn, 0, len(socketIDs))
	for socketID := range socketIDs {
		compositeKey := sm.getCompositeKey(userID, socketID)
		if conn, exists := sm.userSocketConnections[compositeKey]; exists {
			connections = append(connections, conn)
		}
	}

	if len(connections) == 0 {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	return connections, nil
}

// Get user's room
func (sm *WebSocketManagerService) GetUserRoom(userID uuid.UUID) (uuid.UUID, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	socketIDs, exists := sm.userIDSocketIDs[userID]
	if !exists || len(socketIDs) == 0 {
		return uuid.UUID{}, fmt.Errorf("user not found: %s", userID)
	}

	// Get first socket ID
	var firstSocketID string
	for sockID := range socketIDs {
		firstSocketID = sockID
		break
	}

	// Get room ID from first socket
	for roomID, roomMeta := range sm.roomMetadata {
		if roomMeta.SocketIDs[firstSocketID] {
			return roomID, nil
		}
	}

	return uuid.UUID{}, fmt.Errorf("user room not found: %s", userID)
}

// Broadcast to room (exclude sender)
func (sm *WebSocketManagerService) BroadcastToRoom(ctx context.Context, roomID uuid.UUID, message []byte, excludeSocketID string) error {
	sm.mu.RLock()

	roomMetadata := sm.GetRoomMetadata(roomID)

	if roomMetadata == nil {
		sm.mu.RUnlock()
		return fmt.Errorf("room not found: %s", roomID)
	}

	var conns []*websocket.Conn

	for socketID := range roomMetadata.SocketIDs {
		if socketID == excludeSocketID {
			continue
		}

		userID := sm.socketIdToUserId[socketID]
		compositeKey := sm.getCompositeKey(userID, socketID)
		if conn, exists := sm.userSocketConnections[compositeKey]; exists {
			conns = append(conns, conn)
		}
	}

	sm.mu.RUnlock()

	var errs []error
	for _, conn := range conns {
		if err := conn.Write(ctx, websocket.MessageText, message); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("broadcast errors: %v", errs)
	}

	return nil
}

// Broadcast to room (include sender)
func (sm *WebSocketManagerService) BroadcastToRoomIncludeSender(ctx context.Context, roomID uuid.UUID, message []byte) error {
	sm.mu.RLock()

	roomMeta, exists := sm.roomMetadata[roomID]

	if !exists {
		sm.mu.RUnlock()
		return fmt.Errorf("room not found: %s", roomID)
	}

	var conns []*websocket.Conn

	for socketID := range roomMeta.SocketIDs {
		userID := sm.socketIdToUserId[socketID]
		compositeKey := sm.getCompositeKey(userID, socketID)
		if conn, exists := sm.userSocketConnections[compositeKey]; exists {
			conns = append(conns, conn)
		}
	}

	sm.mu.RUnlock()

	var errs []error
	for _, conn := range conns {
		if err := conn.Write(ctx, websocket.MessageText, message); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("broadcast errors: %v", errs)
	}

	return nil
}

// Send to specific socket
func (sm *WebSocketManagerService) SendToSocket(ctx context.Context, socketID string, message []byte) error {
	conn, err := sm.GetConnFromSocket(socketID)
	if err != nil {
		return err
	}

	return conn.Write(ctx, websocket.MessageText, message)
}

// Send to all sockets of a user
func (sm *WebSocketManagerService) SendToUser(ctx context.Context, userID uuid.UUID, message []byte) error {
	conns, err := sm.GetUserConnections(userID)
	if err != nil {
		return err
	}

	var errs []error
	for _, conn := range conns {
		if err := conn.Write(ctx, websocket.MessageText, message); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("send errors: %v", errs)
	}

	return nil
}

// Get all users in a room
func (sm *WebSocketManagerService) GetUsersInRoom(roomID uuid.UUID) []uuid.UUID {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	roomMeta, exists := sm.roomMetadata[roomID]
	if !exists {
		return []uuid.UUID{}
	}

	users := make([]uuid.UUID, 0, len(roomMeta.Users))
	for userID := range roomMeta.Users {
		users = append(users, userID)
	}
	return users
}

// Get room user count
func (sm *WebSocketManagerService) GetRoomUserCount(roomID uuid.UUID) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	roomMeta, exists := sm.roomMetadata[roomID]
	if !exists {
		return 0
	}

	return len(roomMeta.Users)
}

// Get room socket count
func (sm *WebSocketManagerService) GetRoomSocketCount(roomID uuid.UUID) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	roomMeta, exists := sm.roomMetadata[roomID]
	if !exists {
		return 0
	}

	return len(roomMeta.SocketIDs)
}

// Check if socket exists
func (sm *WebSocketManagerService) SocketExists(socketID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, exists := sm.socketIdToUserId[socketID]
	return exists
}

// Get room host
func (sm *WebSocketManagerService) GetRoomHost(roomID uuid.UUID) (uuid.UUID, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	roomMeta, exists := sm.roomMetadata[roomID]
	if !exists {
		return uuid.UUID{}, fmt.Errorf("room not found: %s", roomID)
	}

	return roomMeta.HostID, nil
}

// Get user socket count
func (sm *WebSocketManagerService) GetUserSocketCount(userID uuid.UUID) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	socketIDs, exists := sm.userIDSocketIDs[userID]
	if !exists {
		return 0
	}

	return len(socketIDs)
}
