package models

import (
	"bitgifty.com/stellar/database"
	"gorm.io/gorm"
)

type Redeem struct {
	gorm.Model
	Address    string    `json:"address"`
	GiftCardID uint      `json:"giftcard_id"`
	GiftCard   *GiftCard `json:"giftcard"`
	Hash       string    `json:"hash"`
	UserID     uint      `json:"user_id"`
	User       *User     `gorm:"foreignKey:UserID;references:ID" json:"user"`
}

func (r *Redeem) Save() (*Redeem, error) {
	err := database.DB.Save(&r).Error
	if err != nil {
		return &Redeem{}, err
	}

	return r, nil
}
