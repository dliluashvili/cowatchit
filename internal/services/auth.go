package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/helpers"
	"github.com/dliluashvili/cowatchit/internal/models"
)

type AuthService struct {
	sessionService *SessionService
	userService    *UserService
}

func NewAuthService(sessionService *SessionService, userService *UserService) *AuthService {
	return &AuthService{
		sessionService: sessionService,
		userService:    userService,
	}
}

func (s *AuthService) SignIn(ctx context.Context, dto *dtos.SignInDto) (*models.Session, int, error) {
	if dto.Password == nil || *dto.Password == "" {
		return nil, 500, errors.New("password is required")
	}

	user, err := s.checkUser(*dto.Username)
	if err != nil {
		return nil, 500, errors.New("invalid credentials")
	}

	if user == nil {
		return nil, 500, errors.New("invalid credentials")
	}

	password := *dto.Password

	err = helpers.ComparePassword(user.Password, password)
	if err != nil {
		return nil, 500, errors.New("invalid credentials")
	}

	if helpers.NeedsRehash(user.Password, helpers.DefaultCost) {
		newHash, _ := helpers.HashPassword(password)
		err := s.userService.UpdatePasswordHash(user.ID, newHash)

		if err != nil {
			return nil, 500, errors.New("failed on password rehash")
		}
	}

	// Generate session
	session, err := s.generateSessionModel(ctx, user)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to create session: %w", err)
	}

	return session, 200, nil
}

func (s *AuthService) SignUp(ctx context.Context, dto *dtos.SignUpDto) (*models.Session, int, error) {
	found, err := s.checkUser(*dto.Username)

	if err != nil {
		return nil, 500, err
	}

	if found != nil {
		return nil, 409, fmt.Errorf("already exists")
	}

	if *dto.Password != *dto.PasswordConfirmation {
		return nil, 500, fmt.Errorf("bad request")
	}

	hashed, err := helpers.HashPassword(*dto.Password)

	if err != nil {
		return nil, 500, err
	}

	newUser, err := s.userService.Create(&dtos.CreateUserDto{
		Username:    dto.Username,
		Gender:      dto.Gender,
		Password:    &hashed,
		Email:       dto.Email,
		DateOfBirth: helpers.StringToDate(*dto.DateOfBirth),
	})

	if err != nil || newUser == nil {
		return nil, 500, err
	}

	session, err := s.generateSessionModel(ctx, newUser)

	if err != nil {
		return nil, 500, err
	}

	return session, 201, nil

}

func (s *AuthService) generateSessionModel(ctx context.Context, user *models.User) (*models.Session, error) {
	session, err := s.sessionService.CreateAndSave(ctx, &models.User{
		ID:       user.ID,
		Username: user.Username,
		Age:      user.Age,
		Gender:   user.Gender,
	})

	return session, err
}

func (s *AuthService) checkUser(username string) (*models.User, error) {
	return s.userService.FindByUsername(username)
}
