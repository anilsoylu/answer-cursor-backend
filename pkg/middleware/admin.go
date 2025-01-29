package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		role := c.GetString("role")
		status := c.GetString("status")

		// Check if user has admin privileges
		if role != "ADMIN" && role != "SUPER_ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "Admin privileges required",
				},
			})
			c.Abort()
			return
		}

		// Check if user status is active
		if status != "active" {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "Account is not active",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
} 