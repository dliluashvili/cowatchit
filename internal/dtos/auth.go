package dtos

import (
	"time"
)

type SignUpDto struct {
	Username             *string `json:"username" validate:"required,unique,username"`
	Email                *string `json:"email" validate:"required,email,unique"`
	Gender               *string `json:"gender" validate:"required,gender"`
	Password             *string `json:"password" validate:"required,min=6,max=30"`
	PasswordConfirmation *string `json:"password_confirmation" validate:"required,min=6,max=30"`
	DateOfBirth          *string `json:"date_of_birth" validate:"required,dob"`
}

type SignInDto struct {
	Username *string `json:"username" validate:"required,username"`
	Password *string `json:"password" validate:"required,min=6,max=30"`
}

type CreateUserDto struct {
	Username             *string   `json:"username"`
	Email                *string   `json:"email"`
	Gender               *string   `json:"gender"`
	Password             *string   `json:"password"`
	Age                  *uint8    `json:"age"`
	PasswordConfirmation *string   `json:"password_confirmation"`
	DateOfBirth          time.Time `json:"date_of_birth"`
}
