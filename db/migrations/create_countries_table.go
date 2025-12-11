package migrations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Country struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title      string    `gorm:"type:text"`
	PhoneCode  string    `gorm:"type:text"`
	EmojiU     string    `gorm:"type:text"`
	Native     string    `gorm:"type:text"`
	OriginalID string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func CreateCountriesTable(db *gorm.DB) {
	db.AutoMigrate(&Country{})
}
