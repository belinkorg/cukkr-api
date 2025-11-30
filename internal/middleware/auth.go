package middleware

import (
	"bLink-app/pkg/errors"
	"bLink-app/pkg/jwt"
	"bLink-app/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Should use authorization in header", errors.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "There should split with space between bearer and token", errors.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenParts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid token", err)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
