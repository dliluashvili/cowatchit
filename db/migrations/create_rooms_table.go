package migrations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	HostID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Host         User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:HostID;references:ID"`
	HostUsername string    `gorm:"type:varchar(255);not null"`
	Title        string    `gorm:"type:varchar(255);"`
	Capacity     int
	Description  string `gorm:"type:varchar(800)"`
	Src          string `gorm:"type:varchar(500);not null"`
	Poster       string `gorm:"type:varchar(500)"`
	Private      bool
	Hidden       bool
	Password     string    `gorm:"type:varchar(30)"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func CreateRoomTable(dbconnection *gorm.DB) {
	dbconnection.AutoMigrate(&Room{})
}
