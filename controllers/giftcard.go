package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/just-nibble/stellar-gin/database"
	"github.com/just-nibble/stellar-gin/helpers"
	"github.com/just-nibble/stellar-gin/models"
	"github.com/just-nibble/stellar-gin/pkg"
)

const DEST = ""
const ADMINKEY = ""

type CreateGiftCardInput struct {
	Address       string `json:"address" binding:"required"`
	Amount        int    `json:"amount" binding:"required"`
	Net           string `json:"net" binding:"required"`
	ReceiverEmail string `json:"receiver_email" binding:"required"`
	Key           string `json:"secret_key" binding:"required"`
}

type RedeemGiftCardInput struct {
	Code    string `json:"code" binding:"required"`
	Address string `json:"address" binding:"required"`
}

func CreateGiftCard(c *gin.Context) {
	currentUser, err := helpers.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "unathorized user", "error": "unathorized user"})
		return
	}

	var input CreateGiftCardInput

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "invalid input", "error": err.Error()})
		return
	}

	giftCard := models.GiftCard{
		Address:       input.Address,
		Amount:        input.Amount,
		Net:           input.Net,
		ReceiverEmail: input.ReceiverEmail,
		UserID:        currentUser.ID,
	}

	stellarClient := pkg.StellarClient{Net: input.Net}

	amount := strconv.Itoa(input.Amount)

	hash, err := stellarClient.BuildTransaction(DEST, amount, input.Key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "unexpected error, please try again later", "error": err.Error()})
		return
	}

	giftCard.Code = helpers.RandString(12)
	giftCard.Hash = *hash

	savedCard, err := giftCard.Save()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "giftcard already exists", "error": err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "unexpected error, please try again later", "error": err.Error()})
		return
	}

	go func() {
		helpers.SendGiftCardEmail("info@bitgifty.com", []string{input.ReceiverEmail}, giftCard.Code)
	}()

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "giftcard created successfully", "data": savedCard})
}

func RedeemGiftCard(c *gin.Context) {
	currentUser, err := helpers.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "unathorized user", "error": "unathorized user"})
		return
	}

	var input RedeemGiftCardInput

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "invalid input", "error": err.Error()})
		return
	}

	giftCard := models.GiftCard{}

	err = database.DB.Where("giftcards.code = ?", input.Code).First(&giftCard).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "giftcard not found", "error": err.Error()})
		return
	}

	if giftCard.ReceiverEmail != currentUser.Email || giftCard.Status == "used" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "shamless comrade", "error": "unathorized user"})
	}

	stellarClient := pkg.StellarClient{}
	amount := strconv.Itoa(giftCard.Amount)

	hash, err := stellarClient.BuildTransaction(input.Address, amount, ADMINKEY)

	redeem := models.Redeem{}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "unexpected error, please try again later", "error": err.Error()})
		return
	}

	redeem.Hash = *hash
	savedEntry, err := redeem.Save()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "unexpected error, please try again later", "error": err.Error()})
		return
	}

	giftCard.Status = "used"
	_, err = giftCard.Save()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "unexpected error, please try again later", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "giftcard redeemed successfully", "data": savedEntry})
}
