package types

import "github.com/dliluashvili/cowatchit/internal/models"

type RoomUser struct {
	ID string
}

type RoomParticipants struct {
	ID       string
	Username string
}

type RoomPageContext struct {
	User         RoomUser
	Room         models.Room
	IsHost       bool
	IsPlaying    bool
	Participants []RoomParticipants
	CurrentTime  int
	Duration     int
}
