package middleware

import (
	"net/http"
	"strconv"
	"strings"

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

// RoleMiddleware checks the X-User-Roles header
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := c.GetHeader("X-User-Roles")
		if roles == "" {
			// Fallback or legacy check
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions (no roles found)"})
			return
		}

		// Simplified check: search for the role in the comma-separated list
		// In a real scenario, you'd split and check properly
		
		roleList := strings.Split(roles, ",")
		hasRole := false
		for _, r := range roleList {
			if strings.TrimSpace(r) == requiredRole || strings.TrimSpace(r) == "ROLE_" + requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}
