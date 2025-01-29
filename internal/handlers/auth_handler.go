package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anilsoylu/answer-backend/internal/models"
	"github.com/anilsoylu/answer-backend/internal/services"
	"github.com/anilsoylu/answer-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Username string         `json:"username" binding:"required,min=3,max=50"`
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6"`
	Role     models.UserRole `json:"role" binding:"omitempty,oneof=USER EDITOR ADMIN SUPER_ADMIN"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"user"`
}

type FreezeAccountRequest struct {
	UserID       uint   `json:"user_id" binding:"required"`
	FreezeReason string `json:"freeze_reason" binding:"required"`
}

type AuthHandler struct {
	authService *services.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	user := &models.User{
		Username:      req.Username,
		Email:         req.Email,
		Password:      req.Password,
		Status:        models.StatusActive,
		Role:          models.RoleUser, // Always starts with USER role
		CreatedAt:     time.Now(),
		LastLoginDate: time.Now(),
	}

	if err := h.authService.Register(user); err != nil {
		switch err {
		case services.ErrUsernameTaken:
			c.JSON(http.StatusConflict, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "username_taken",
					"message": "Username is already taken",
				},
			})
		case services.ErrEmailTaken:
			c.JSON(http.StatusConflict, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "email_taken",
					"message": "Email is already taken",
				},
			})
		case services.ErrUserFrozen:
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "account_frozen",
					"message": "This account is frozen, please contact support",
				},
			})
		case services.ErrUserBanned:
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "account_banned",
					"message": "This account is banned",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "internal_error",
					"message": "An error occurred",
				},
			})
		}
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

	c.JSON(http.StatusCreated, gin.H{
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

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	user, err := h.authService.Login(req.Identifier, req.Password)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "invalid_credentials",
					"message": "Invalid username/email or password",
				},
			})
		case services.ErrUserNotActive:
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "user_not_active",
					"message": "User account is not active",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "internal_error",
					"message": err.Error(),
				},
			})
		}
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

func (h *AuthHandler) UpdateUserRole(c *gin.Context) {
	var req models.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	// Get requester's information
	requesterID := c.GetUint("user_id")
	var requester models.User
	if err := h.authService.GetUserByID(requesterID, &requester); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "internal_error",
				"message": "An error occurred while fetching user information",
			},
		})
		return
	}

	if err := h.authService.UpdateUserRole(req.UserID, req.Role, &requester); err != nil {
		switch err.Error() {
		case "root admin's role cannot be changed":
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "Root admin's role cannot be changed",
				},
			})
		case "only root admin can assign SUPER_ADMIN role":
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "Only root admin can assign SUPER_ADMIN role",
				},
			})
		case "only root admin can change SUPER_ADMIN's role":
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "Only root admin can change SUPER_ADMIN's role",
				},
			})
		case "admin can only modify USER and EDITOR roles":
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "Admin can only modify USER and EDITOR roles",
				},
			})
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "not_found",
					"message": "User not found",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "internal_error",
					"message": "An error occurred while updating user role",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "User role updated successfully",
		},
	})
}

func (h *AuthHandler) UpdateUserStatus(c *gin.Context) {
	var req models.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	// Get user ID from token
	userID := c.GetUint("user_id")
	
	// Get requester's role
	requesterRole := models.UserRole(c.GetString("role"))

	// Get target user ID from URL parameter
	targetUserID := userID // Varsayılan olarak kendi ID'si

	if targetID := c.Query("user_id"); targetID != "" {
		// Eğer query parameter varsa ve admin/super_admin ise başka kullanıcıyı güncelle
		if requesterRole == models.RoleAdmin || requesterRole == models.RoleSuperAdmin {
			if id, err := strconv.ParseUint(targetID, 10, 32); err == nil {
				targetUserID = uint(id)
			}
		}
	}

	if err := h.authService.UpdateUserStatus(targetUserID, req.Status, requesterRole); err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "not_found",
					"message": "User not found",
				},
			})
		case services.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "unauthorized",
					"message": "You are not authorized to perform this action",
				},
			})
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "You cannot change this user's status",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "internal_error",
					"message": "Failed to update user status",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "User status updated successfully",
		},
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	// Validasyon kontrolü
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": err.Error(),
			},
		})
		return
	}

	userID := c.GetUint("user_id")
	updatedUser, err := h.authService.UpdateProfile(userID, req.Username, req.Email, req.Avatar)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "not_found",
					"message": "User not found",
				},
			})
		case services.ErrDuplicateEntry:
			c.JSON(http.StatusConflict, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "conflict",
					"message": "Username or email already in use",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "internal_error",
					"message": "Failed to update profile",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": gin.H{
				"id":       updatedUser.ID,
				"username": updatedUser.Username,
				"email":    updatedUser.Email,
				"avatar":   updatedUser.Avatar,
				"status":   updatedUser.Status,
				"role":     updatedUser.Role,
			},
			"message": "Profile updated successfully",
		},
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	// Validasyon kontrolü
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": err.Error(),
			},
		})
		return
	}

	userID := c.GetUint("user_id")
	if err := h.authService.ChangePassword(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "invalid_password",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "Password updated successfully",
		},
	})
}

// UpdatePassword şifre güncelleme handler'ı
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": err.Error(),
			},
		})
		return
	}

	err := h.authService.UpdatePassword(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "internal_error",
				"message": "Failed to update password",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "Password updated successfully",
		},
	})
}

// BanUser kullanıcı banlama handler'ı
func (h *AuthHandler) BanUser(c *gin.Context) {
	var req models.BanUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	// Validasyon kontrolü
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": err.Error(),
			},
		})
		return
	}

	// İşlemi yapan kullanıcının rolünü al
	requesterRole := models.UserRole(c.GetString("role"))

	if err := h.authService.BanUser(c.Request.Context(), req.UserID, req.BanReason, req.BanDuration, requesterRole); err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "not_found",
					"message": "User not found",
				},
			})
		case services.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "unauthorized",
					"message": "You are not authorized to ban users",
				},
			})
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "You cannot ban this user",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "internal_error",
					"message": "Failed to ban user",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "User banned successfully",
		},
	})
}

func (h *AuthHandler) FreezeAccount(c *gin.Context) {
	var req FreezeAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	// Get requester's ID and role
	userID := c.GetUint("user_id")
	requesterRole := c.GetString("role")

	// Kullanıcı sadece kendi hesabını dondurabilir
	if req.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "forbidden",
				"message": "You can only freeze your own account",
			},
		})
		return
	}

	if err := h.authService.FreezeAccount(c.Request.Context(), userID, req.FreezeReason, models.UserRole(requesterRole)); err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "user_not_found",
					"message": "User not found",
				},
			})
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "You don't have permission to freeze this account",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "internal_error",
					"message": err.Error(),
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "Account has been frozen successfully",
		},
	})
}

// DeleteAccount kullanıcıyı siler (soft delete)
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	// Get user ID from URL parameter
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "invalid_id",
				"message": "Invalid user ID",
			},
		})
		return
	}

	// Get requester's ID and role
	requesterID := c.GetUint("user_id")
	requesterRole := models.UserRole(c.GetString("role"))

	if err := h.authService.DeleteAccount(c.Request.Context(), uint(userID), requesterID, requesterRole); err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "user_not_found",
					"message": "User not found",
				},
			})
		case services.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "unauthorized",
					"message": "You can only delete your own account or be a SUPER_ADMIN to delete others",
				},
			})
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "forbidden",
					"message": "SUPER_ADMIN accounts cannot be deleted",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "internal_error",
					"message": "Failed to delete account",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"message": "Account has been soft deleted successfully",
		},
	})
}

func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "validation_error",
				"message": utils.GetValidationError(err),
			},
		})
		return
	}

	user, err := h.authService.AdminLogin(req.Identifier, req.Password)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "invalid_credentials",
					"message": "Invalid username/email or password",
				},
			})
		case services.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "unauthorized",
					"message": "Admin privileges required",
				},
			})
		case services.ErrUserNotActive:
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "user_not_active",
					"message": "User account is not active",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "internal_error",
					"message": err.Error(),
				},
			})
		}
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

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := h.authService.GetUserByID(userID, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "internal_error",
				"message": "Failed to get user profile",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
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