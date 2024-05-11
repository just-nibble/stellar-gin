package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bitgifty.com/stellar/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GenerateJWT(user models.User) (string, error) {
	loadEnv()
	var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":             user.ID,
		"role":           user.Role,
		"fullname":       user.FullName,
		"email":          user.Email,
		"phone":          user.Phone,
		"email_verified": user.EmailVerified,
		"status":         user.Status,
		"iat":            time.Now().Unix(),
		"eat":            time.Now().Add(time.Hour * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

func ValidateJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func GetRole(context *gin.Context) (string, error) {
	err := ValidateJWT(context)
	if err != nil {
		return "", err
	}
	token, err := getToken(context)
	if err != nil {
		return "", nil
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	role := string(claims["role"].(string))

	return role, nil
}

func RoleCheck(context *gin.Context, roleCheck string) (bool, error) {
	role, err := GetRole(context)
	if err != nil {
		return false, nil
	}
	if role != roleCheck {
		return false, errors.New("role does not match")
	}
	return true, nil
}

func CurrentUser(context *gin.Context) (models.User, error) {
	user := models.User{}

	err := ValidateJWT(context)
	if err != nil {
		return user, err
	}

	token, err := getToken(context)
	if err != nil {
		return user, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("claims")
		return user, errors.New("invalid token provided")
	}

	Id, ok := claims["id"].(float64)
	userId := uint(Id)
	if !ok {
		return user, errors.New("invalid token provided")
	}

	role, ok := claims["role"].(string)
	if !ok {
		log.Println("role")
		return user, errors.New("invalid token provided")
	}

	FullName, ok := claims["fullname"].(string)
	if !ok {
		log.Println("fullname")
		return user, errors.New("invalid token provided")
	}

	Email, ok := claims["email"].(string)
	if !ok {
		log.Println("email")
		return user, errors.New("invalid token provided")
	}

	Phone, ok := claims["phone"].(string)
	if !ok {
		log.Println("phone")
		return user, errors.New("invalid token provided")
	}

	EmailVerified, ok := claims["email_verified"].(bool)
	if !ok {
		log.Println("email verified")
		return user, errors.New("invalid token provided")
	}

	Status, ok := claims["status"].(string)
	if !ok {
		log.Println("sta")
		return user, errors.New("invalid token provided")
	}
	user.ID = userId
	user.FullName = FullName
	user.Email = Email
	user.Phone = Phone
	user.EmailVerified = EmailVerified
	user.Status = Status
	user.Role = role

	return user, nil
}

func getToken(context *gin.Context) (*jwt.Token, error) {
	loadEnv()
	var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
