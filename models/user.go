package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Email         string    `gorm:"unique;not null" json:"email"`
	WalletBalance int       `json:"wallet_balance"`
	WalletAsset   string    `gorm:"default:sats" json:"wallet_asset"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type PublicWallet struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	WalletBalance int    `json:"wallet_balance"`
	WalletAsset   string `json:"wallet_asset"`
}

func (u *User) GetWalletBalance(db *gorm.DB) (*PublicWallet, error) {
	if err := db.Debug().Where("email=?", u.Email).First(&u).Error; err != nil {
		return &PublicWallet{}, err
	}
	return &PublicWallet{
		UserID:        u.ID,
		Email:         u.Email,
		WalletBalance: u.WalletBalance,
		WalletAsset:   u.WalletAsset,
	}, nil
}

func (u *User) AddToWallet(db *gorm.DB, amount int) error {
	return db.Model(&u).UpdateColumn("wallet_balance", gorm.Expr("wallet_balance + ?", amount)).Error
}
