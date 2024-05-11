package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/just-nibble/stellar-gin/helpers"
	"github.com/just-nibble/stellar-gin/models"
)

type AuthenticationInput struct {
	FullName    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input AuthenticationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "invalid input", "error": err.Error()})
		return
	}
	user := models.User{
		FullName: input.FullName,
		Email:    input.Email,
		Phone:    input.PhoneNumber,
		Password: input.Password,
		Role:     "user",
	}

	savedUser, err := user.Save()

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "user already exists", "error": err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "unexpected error, please try again later", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "user created successfully", "data": savedUser})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "invalid input", "error": err.Error()})
		return
	}

	user, err := models.FindUserByEmail(input.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "user not found", "error": err.Error()})
		return
	}

	err = user.ValidatePassword(input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "wrong email or password", "error": err.Error()})
		return
	}

	jwt, err := helpers.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "login success", "data": user, "jwt": jwt})
}
