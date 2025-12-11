package repositories

import (
	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (rp *RoomRepository) Create(dto *dtos.CreateRoomRepoDto) (*models.Room, error) {

	hidden := false

	if dto.Private {
		hidden = true
	}

	room := &models.Room{
		ID:           uuid.New(),
		HostID:       dto.HostID,
		HostUsername: dto.HostUsername,
		Title:        dto.Title,
		Description:  dto.Description,
		Src:          dto.Src,
		Capacity:     dto.Capacity,
		Private:      dto.Private,
		Password:     dto.Password,
		Hidden:       hidden,
	}

	result := rp.db.Create(room)

	if result.Error != nil {
		return nil, result.Error
	}

	return room, nil
}

func (rp *RoomRepository) Find(dto *dtos.FindRoomDto) ([]models.Room, error) {
	var rooms []models.Room

	query := rp.db

	if dto.Filter != nil && *dto.Filter != "" {
		switch *dto.Filter {
		case "public":
			query = query.Where("private = ?", false)
		case "private":
			query = query.Where("private = ?", true)
		}
	}

	if dto.My != nil {
		query.Where("host_id = ?", dto.AuthUserID)
	}

	if dto.Keyword != nil && *dto.Keyword != "" {
		query.Where("title ILIKE ?", "%"+*dto.Keyword+"%")
	}

	result := query.Find(&rooms)

	if result.Error != nil {
		return nil, result.Error
	}

	return rooms, nil
}

func (rp *RoomRepository) FindOne(ID uuid.UUID) (*models.Room, error) {
	var room models.Room

	result := rp.db.Where("id = ?", ID).First(&room)

	return &room, result.Error
}

func (rp *RoomRepository) Exists(ID uuid.UUID) (bool, error) {
	var count int64
	result := rp.db.Where("id = ?", ID).Model(&models.Room{}).Count(&count)

	return count > 0, result.Error
}

func (rp *RoomRepository) Join() {}
