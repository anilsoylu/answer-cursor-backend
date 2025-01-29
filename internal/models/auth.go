package models

// RegisterRequest kullanıcı kayıt isteği için model
type RegisterRequest struct {
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
}

// LoginRequest kullanıcı giriş isteği için model
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // Email veya Username
	Password   string `json:"password" validate:"required"`
}

// AuthResponse token yanıtı için model
type AuthResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
}

// UserResponse kullanıcı bilgileri yanıtı için model
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Status    UserStatus `json:"status"`
	Role      UserRole   `json:"role"`
}

// UpdateUserRoleRequest kullanıcı rolü güncelleme isteği için model
type UpdateUserRoleRequest struct {
	UserID uint     `json:"user_id" validate:"required"`
	Role   UserRole `json:"role" validate:"required,oneof=USER EDITOR ADMIN SUPER_ADMIN"`
}

// UpdateUserStatusRequest kullanıcı durumu güncelleme isteği için model
type UpdateUserStatusRequest struct {
	Status UserStatus `json:"status" validate:"required,oneof=active passive banned"`
}

// UpdateProfileRequest kullanıcı profili güncelleme isteği için model
type UpdateProfileRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Avatar   string `json:"avatar,omitempty"`
} 