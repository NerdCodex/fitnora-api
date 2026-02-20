package middleware

import (
	"backend/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
		c.Abort()
		return
	}

	// Check if the Authorization header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
		c.Abort()
		return
	}

	// Get the token part from the Authorization header
	tokenString := authHeader[7:]
	claims, err := services.ValidateToken(tokenString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	c.Set("claims", claims)
	c.Next()
}
