package repositories

import (
	"errors"
	"time"

	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Update(userID uuid.UUID, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}

	safeUpdates := make(map[string]any, len(updates))

	for k, v := range updates {
		if k == "images" {
			if images, ok := v.([]string); ok {
				safeUpdates[k] = pq.StringArray(images)
				continue
			}
		}
		safeUpdates[k] = v
	}

	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(safeUpdates).Error
}

func (r *UserRepository) Create(dto *dtos.CreateUserDto) (*models.User, error) {
	user := models.User{
		ID:          uuid.New(),
		Username:    *dto.Username,
		Email:       *dto.Email,
		Gender:      *dto.Gender,
		Password:    *dto.Password,
		Age:         *dto.Age,
		DateOfBirth: dto.DateOfBirth,
		DeletedAt:   nil,
	}

	result := r.db.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserRepository) FindByField(params map[string]any) (*models.User, error) {
	query := r.db

	for k, v := range params {
		if k == "" {
			continue // Skip invalid keys
		}
		query = query.Where(k+" = ?", v)
	}

	var u models.User

	if err := query.First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}

	}

	return &u, nil
}

func (r *UserRepository) Delete(ID uuid.UUID) (*bool, error) {
	result := r.db.Model(&models.User{}).
		Where("id = ?", ID).
		Updates(map[string]any{
			"deletedAt": time.Now().UTC(),
		})

	if result.Error != nil {
		return nil, result.Error
	}

	success := result.RowsAffected > 0

	return &success, nil
}
