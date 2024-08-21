package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokha04/go-e-commerce/tokens"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ClientToken := ctx.Request.Header.Get("token")
		if ClientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no authorization header provided"})
			ctx.Abort()
			return
		}

		claims, err := tokens.ValidateToken(ClientToken)
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("uid", claims.Uid)
		ctx.Next()
	}
}
