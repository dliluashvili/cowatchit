package dtos

import "github.com/google/uuid"

type CreateRoomDto struct {
	Title       string `json:"title" validate:"required,roomtitle"`
	Capacity    int    `json:"capacity" validate:"required,min=2,max=10"`
	Description string `json:"description" validate:"required"`
	Src         string `json:"src" validate:"required,url,max=500"`
	Private     bool   `json:"private"`
	Password    string `json:"password" validate:"required_if=Private true,omitempty,min=3,max=20"`
}

type CreateRoomServiceDto struct {
	HostID uuid.UUID `json:"host_id"`
	*CreateRoomDto
}

type CreateRoomRepoDto struct {
	HostUsername string `json:"host_username"`
	*CreateRoomServiceDto
}

type FindRoomDto struct {
	Keyword    *string    `json:"keyword"`
	Filter     *string    `json:"filter"`
	ID         *uuid.UUID `json:"id"`
	My         *bool      `json:"my"`
	AuthUserID *uuid.UUID `json:"auth_user_id"`
}
