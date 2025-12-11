package migrations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	SessionID string    `gorm:"type:varchar(255);uniqueIndex"`
	ExpiresAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateSessionTable(dbconnection *gorm.DB) {
	dbconnection.AutoMigrate(&Session{})
}
