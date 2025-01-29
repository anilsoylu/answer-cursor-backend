package models

// RegisterRequest represents the model for user registration request
type RegisterRequest struct {
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
}

// LoginRequest represents the model for user login request
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // Email or Username
	Password   string `json:"password" validate:"required"`
}

// AuthResponse represents the model for token response
type AuthResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
}

// UserResponse represents the model for user information response
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Status    UserStatus `json:"status"`
	Role      UserRole   `json:"role"`
}

// UpdateUserRoleRequest represents the model for updating user role request
type UpdateUserRoleRequest struct {
	UserID uint     `json:"user_id" validate:"required"`
	Role   UserRole `json:"role" validate:"required,oneof=USER EDITOR ADMIN SUPER_ADMIN"`
}

// UpdateUserStatusRequest represents the model for updating user status request
type UpdateUserStatusRequest struct {
	Status UserStatus `json:"status" validate:"required,oneof=active passive banned"`
}

// UpdateProfileRequest represents the model for updating user profile request
type UpdateProfileRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Avatar   string `json:"avatar,omitempty"`
}

// UpdatePasswordRequest represents the model for updating password request
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// BanUserRequest represents the model for banning user request
type BanUserRequest struct {
	UserID      uint   `json:"user_id" validate:"required"`
	BanReason   string `json:"ban_reason" validate:"required,min=10,max=500"`
	BanDuration string `json:"ban_duration" validate:"required,oneof=1_day 1_week 1_month permanent"`
} 