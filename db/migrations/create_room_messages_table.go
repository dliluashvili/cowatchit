package migrations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomMessage struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	SenderID       uuid.UUID `gorm:"type:uuid;not null"`
	User           User      `gorm:"foreignKey:SenderID"`
	SenderUsername string    `gorm:"type:varchar(255)"`
	RoomID         uuid.UUID `gorm:"type:uuid;not null;index"`
	Room           Room      `gorm:"foreignKey:RoomID"`
	IsHost         bool      `gorm:"type:bool;not null;default:false"`
	Content        string    `gorm:"type:text"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func CreateRoomMessageTable(db *gorm.DB) {
	db.AutoMigrate(&RoomMessage{})
}
