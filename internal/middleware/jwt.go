package middleware

import (
	_ "fmt"
	"log"
	_ "log"
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

// GetUserRoles retrieves roles from context (set by RoleMiddleware or Gateway headers)
func GetUserRoles(c *gin.Context) []string {
	rolesHeader := c.GetHeader("X-User-Roles")
	if rolesHeader == "" {
		return []string{}
	}
	rawList := strings.Split(rolesHeader, ",")
	roles := make([]string, 0, len(rawList))
	for _, r := range rawList {
		roles = append(roles, strings.ToUpper(strings.TrimPrefix(strings.TrimSpace(r), "ROLE_")))
	}
	return roles
}

func HasRole(roles []string, target string) bool {
	target = strings.ToUpper(target)
	for _, r := range roles {
		if r == target {
			return true
		}
	}
	return false
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
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := c.GetHeader("X-User-Roles")
		if roles == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions (no roles found)"})
			return
		}

		roleList := strings.Split(roles, ",")
		hasRole := false

		for _, r := range roleList {
			cleanRole := strings.ToUpper(strings.TrimSpace(r))
			// Debug log for roles
			log.Printf("[AUTH DEBUG] Raw role from header: %s, Cleaned role: %s", r, cleanRole)

			// Remove ROLE_ prefix if it exists to compare pure roles
			cleanRole = strings.TrimPrefix(cleanRole, "ROLE_")

			for _, reqRole := range requiredRoles {
				if cleanRole == strings.ToUpper(reqRole) {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			log.Printf("[AUTH DEBUG] Access denied. Required one of: %v, but user has: %s", requiredRoles, roles)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}
