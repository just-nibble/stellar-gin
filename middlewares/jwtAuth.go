package middlewares

import (
	"log"
	"net/http"

	"bitgifty.com/stellar/helpers"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := helpers.ValidateJWT(context)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Authentication required", "error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}
