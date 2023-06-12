package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    int       `json:"user_id"`
	Action    string    `json:"action"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (t *TransactionLog) Add(db *gorm.DB) error {
	if err := db.Debug().Create(&t).Error; err != nil {
		return err
	}
	return nil
}
