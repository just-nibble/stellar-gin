package models

import (
	"github.com/just-nibble/stellar-gin/database"
	"gorm.io/gorm"
)

type GiftCard struct {
	gorm.Model
	Address       string `json:"address"`
	Amount        int    `json:"amount"`
	Code          string `json:"code"`
	Net           string `json:"net"`
	Hash          string `json:"hash"`
	Status        string `json:"status"`
	ReceiverEmail string `json:"receiver_email"`
	UserID        uint   `json:"user_id"`
	User          *User  `gorm:"foreignKey:UserID;references:ID" json:"user"`
}

func (g *GiftCard) Save() (*GiftCard, error) {
	err := database.DB.Save(&g).Error
	if err != nil {
		return &GiftCard{}, err
	}

	return g, nil
}
