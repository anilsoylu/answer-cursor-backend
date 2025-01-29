package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "unauthorized",
					"message": "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "unauthorized",
					"message": "Invalid authorization header format",
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // TODO: Use environment variable
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "unauthorized",
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "unauthorized",
					"message": "Invalid token claims",
				},
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Set("email", claims["email"].(string))
		c.Set("role", claims["role"].(string))
		c.Set("status", claims["status"].(string))

		c.Next()
	}
} 