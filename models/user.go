package models

import (
	"errors"

	"github.com/just-nibble/stellar-gin/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName      string `gorm:"size:255;" json:"full_name"`
	Email         string `gorm:"size:255;not null;unique" json:"email"`
	Phone         string `gorm:"size:255;not null;unique" json:"phone_number"`
	Role          string `gorm:"size:255;default:customer" json:"role"`
	Password      string `gorm:"type:varchar(255);not null;" json:"-"`
	EmailVerified bool   `gorm:"default:false" json:"email_verified"`
	Status        string `gorm:"size:255" json:"status"`
}

func (user *User) Save() (*User, error) {
	err := database.DB.Save(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &User{}, errors.New("email or phone number already used")
		}
		return &User{}, err
	}
	return user, nil
}

func (user *User) BeforeSave(*gorm.DB) error {
	if user.ID < 1 {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(passwordHash)
		return nil
	}
	return nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByEmail(email string) (User, error) {
	var user User
	err := database.DB.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func FindUserByPhone(phone string) (User, error) {
	var user User
	err := database.DB.Where("phone_number = ?", phone).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func FindUserById(id uint) (User, error) {
	var user User
	err := database.DB.Where("ID=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}
