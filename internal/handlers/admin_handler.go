package handlers

import (
	"net/http"

	"github.com/anilsoylu/answer-backend/internal/services"
	"github.com/anilsoylu/answer-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	authService *services.AuthService
}

func NewAdminHandler(authService *services.AuthService) *AdminHandler {
	return &AdminHandler{
		authService: authService,
	}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": "Invalid request body",
			},
		})
		return
	}

	user, err := h.authService.AdminLogin(req.Identifier, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "unauthorized",
				"message": err.Error(),
			},
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Email, user.Role, user.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "token_error",
				"message": "Failed to generate token",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"status":     user.Status,
				"role":       user.Role,
				"avatar":     user.Avatar,
				"created_at": user.CreatedAt,
			},
		},
	})
}

func (h *AdminHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var user struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Status    string `json:"status"`
		Role      string `json:"role"`
		Avatar    string `json:"avatar"`
		CreatedAt string `json:"created_at"`
	}

	if err := h.authService.GetAdminProfile(userID, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "internal_error",
				"message": "Failed to get admin profile",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": user,
		},
	})
} 