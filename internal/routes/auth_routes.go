package routes

import (
	"github.com/anilsoylu/answer-backend/internal/handlers"
	"github.com/anilsoylu/answer-backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		users := protected.Group("/users")
		{
			users.POST("/freeze", authHandler.FreezeAccount)
			users.DELETE("/:id", authHandler.DeleteAccount)
			users.PUT("/status", authHandler.UpdateUserStatus)
			users.PUT("/role", authHandler.UpdateUserRole)
			users.PUT("/profile", authHandler.UpdateProfile)
			users.PUT("/password", authHandler.UpdatePassword)
		}
	}
} 