package migrations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomUser struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoomID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Room      Room      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:RoomID;references:ID"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateRoomUserTable(dbconnection *gorm.DB) {
	dbconnection.AutoMigrate(&RoomUser{})
}
