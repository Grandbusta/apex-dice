package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/Grandbusta/apex-dice/config"
	"gorm.io/gorm"
)

type User struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	Email           string           `gorm:"unique;not null" json:"email"`
	WalletBalance   int              `json:"wallet_balance"`
	WalletAsset     string           `gorm:"default:sats" json:"wallet_asset"`
	Games           []Game           `json:"games"`
	TransactionLogs []TransactionLog `json:"transactions"`
	CreatedAt       time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type PublicWallet struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	WalletBalance string `json:"wallet_balance"`
	WalletAsset   string `json:"wallet_asset"`
}

func (r *User) BeforeCreate(tx *gorm.DB) (err error) {
	r.WalletAsset = config.FUND_ASSET
	return
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	if err := db.Debug().Create(&u).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) GetWalletBalance(db *gorm.DB) (*PublicWallet, error) {
	if err := db.Debug().Where("email=?", u.Email).First(&u).Error; err != nil {
		return &PublicWallet{}, err
	}
	return &PublicWallet{
		UserID:        u.ID,
		Email:         u.Email,
		WalletBalance: strconv.Itoa(u.WalletBalance),
		WalletAsset:   u.WalletAsset,
	}, nil
}

func (u *User) AddToWallet(db *gorm.DB, amount int) error {
	if err := db.Debug().Where("email=?", u.Email).First(&u).Error; err != nil {
		return err
	}
	u.WalletBalance += amount
	return db.Debug().Where("email=?", u.Email).Updates(&u).Error
}

func (u *User) DeductWallet(db *gorm.DB, amount int) error {
	if err := db.Debug().Where("email=?", u.Email).First(&u).Error; err != nil {
		return err
	}
	if u.WalletBalance < amount {
		return errors.New("Balance not enough")
	}
	u.WalletBalance -= amount
	return db.Debug().Where("email=?", u.Email).Updates(&u).Error
}
