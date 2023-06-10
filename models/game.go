package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Active     bool      `gorm:"default:true" json:"active"`
	Target     int       `json:"target"`
	FirstRoll  int       `json:"first_roll"`
	SecondRoll int       `json:"second_roll"`
	Cost       int       `json:"cost"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (g *Game) SaveGame(db *gorm.DB) (*Game, error) {
	if err := db.Debug().Create(&g).Error; err != nil {
		return &Game{}, err
	}
	return g, nil
}
