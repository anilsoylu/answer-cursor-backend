package routes

import (
	"github.com/anilsoylu/answer-backend/internal/handlers"
	"github.com/anilsoylu/answer-backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	admin := router.Group("/api/v1/admin")
	{
		// Public admin routes
		admin.POST("/login", authHandler.AdminLogin)

		// Protected admin routes
		protected := admin.Group("")
		protected.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			protected.GET("/me", authHandler.Me)
		}
	}
} 