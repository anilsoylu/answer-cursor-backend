package models

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus string
type UserRole string

const (
	StatusActive  UserStatus = "active"
	StatusPassive UserStatus = "passive"
	StatusBanned  UserStatus = "banned"
	StatusFrozen  UserStatus = "frozen"

	RoleUser       UserRole = "USER"
	RoleEditor     UserRole = "EDITOR"
	RoleAdmin      UserRole = "ADMIN"
	RoleSuperAdmin UserRole = "SUPER_ADMIN"
)

// User represents the user model in the database
type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Username      string         `json:"username" gorm:"not null"`
	Email         string         `json:"email" gorm:"not null"`
	Password      string         `json:"-" gorm:"not null"`
	Avatar        string         `json:"avatar"`
	Status        UserStatus     `json:"status" gorm:"type:user_status"`
	Role          UserRole       `json:"role" gorm:"type:user_role"`
	IsRootAdmin   bool          `json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	LastLoginDate time.Time      `json:"last_login_date"`
	BanReason     string         `json:"ban_reason,omitempty"`
	BanEndDate    *time.Time     `json:"ban_end_date,omitempty"`
	FrozenReason  string         `json:"frozen_reason,omitempty"`
	FrozenDate    *time.Time     `json:"frozen_date,omitempty"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate is a GORM hook that runs before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Status == "" {
		u.Status = StatusActive
	}
	if u.Role == "" && !u.IsRootAdmin {
		u.Role = RoleUser
	}
	if u.Avatar == "" {
		u.Avatar = "/uploads/default/avatar.png"
	}
	return nil
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// ChangePasswordRequest represents the model for password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
} 