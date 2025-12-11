package services

import (
	"fmt"

	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/helpers"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/dliluashvili/cowatchit/internal/repositories"
	"github.com/google/uuid"
)

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService(r *repositories.UserRepository) *UserService {
	return &UserService{
		repository: r,
	}
}

func (s *UserService) Create(dto *dtos.CreateUserDto) (*models.User, error) {
	age := helpers.CalculateAge(dto.DateOfBirth)

	if age < 18 {
		return nil, fmt.Errorf("incorrect age")
	}

	if age > 255 {
		return nil, fmt.Errorf("age too large for storage")
	}

	ageUint8 := uint8(age)

	dto.Age = &ageUint8

	return s.repository.Create(dto)
}

func (s *UserService) Me(userID uuid.UUID) (*models.User, error) {
	user, err := s.repository.FindByField(map[string]any{
		"id": userID,
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindByUsername(username string) (*models.User, error) {
	return s.repository.FindByField(map[string]any{"username": username})
}

func (s *UserService) FindByID(idStr string) (*models.User, error) {
	id, err := uuid.Parse(idStr)

	if err != nil {
		return nil, fmt.Errorf("invalid userid")
	}

	return s.repository.FindByField(map[string]any{"id": id})
}

func (s *UserService) UpdatePasswordHash(ID uuid.UUID, passwordHash string) error {
	return s.repository.Update(ID, map[string]any{
		"password": passwordHash,
	})
}

func (s *UserService) Delete(ID uuid.UUID) (*bool, error) {
	return s.repository.Delete(ID)
}
