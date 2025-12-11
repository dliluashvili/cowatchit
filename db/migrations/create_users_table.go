package migrations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username    string    `gorm:"type:varchar(255);uniqueIndex"`
	Email       string    `gorm:"type:varchar(255);uniqueIndex"`
	Password    string    `gorm:"type:varchar(255)"`
	Gender      string    `gorm:"type:varchar(1);index"`
	DateOfBirth time.Time
	Age         uint8
	DeletedAt   time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func CreateUserTable(dbconnection *gorm.DB) {
	dbconnection.AutoMigrate(&User{})
}
