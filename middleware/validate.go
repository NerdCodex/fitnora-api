package middleware

import (
	"backend/models"
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
	claimsData, err := services.ValidateToken(tokenString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	// Safe extraction
	userID, ok := claimsData["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token payload"})
		c.Abort()
		return
	}

	userEmail, ok := claimsData["user_email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token payload"})
		c.Abort()
		return
	}

	claims := &models.AccessTokenClaims{
		UserID:    uint64(userID),
		UserEmail: userEmail,
	}

	c.Set("claims", claims)
	c.Next()
}
