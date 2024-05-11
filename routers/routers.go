package routers

import (
	"io"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/just-nibble/stellar-gin/controllers"
)

var f, _ = os.Create("gin.log")

func NewRouter() *gin.Engine {
	config := cors.DefaultConfig()
	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := gin.Default()
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:8000", "http://localhost:5173", "http://localhost:5174"},
	// 	AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{
		"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "authorization",
		"Referer", "User-Agent",
	}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	v1Routes := r.Group("/api/v1")
	authRoutes := v1Routes.Group("/auth")
	authRoutes.POST("/register", controllers.Register)
	authRoutes.POST("/login", controllers.Login)

	cardRoutes := v1Routes.Group("/gift-card")
	cardRoutes.POST("/create", controllers.CreateGiftCard)
	cardRoutes.POST("/redeem", controllers.RedeemGiftCard)

	return r
}
