package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserID retrieves the UserID safely from the context
func GetUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	userID, ok := val.(uint)
	return userID, ok
}

// GatewayAuthMiddleware reads the X-User-Id header injected by Spring Cloud Gateway
func GatewayAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-Id")
		
		if userIDStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "x-user-id header is missing"})
			return
		}

		userIDFloat, err := strconv.ParseFloat(userIDStr, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "x-user-id header is invalid"})
			return
		}

		c.Set("user_id", uint(userIDFloat))
		c.Next()
	}
}
