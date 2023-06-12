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

func GetTrasactionLogs(db *gorm.DB) ([]TransactionLog, error) {
	var logs []TransactionLog
	err := db.Debug().Model(&TransactionLog{}).Find(&logs).Error
	if err != nil {
		return logs, err
	}
	return logs, nil

}
