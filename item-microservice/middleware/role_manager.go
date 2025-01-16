package middleware

import (
	"github.com/gin-gonic/gin"
)

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("Roles")
		print(userRoles)
		print(exists)
		if !exists {
			c.JSON(403, gin.H{"error": "No roles found"})
			c.Abort()
			return
		}

		roles, ok := userRoles.([]string)
		if !ok {
			c.JSON(403, gin.H{"error": "Invalid roles format"})
			c.Abort()
			return
		}

		hasRole := false
		for _, userRole := range roles {
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(403, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
