package models

import (
	"time"
)

type UserStatus string
type UserRole string

const (
	StatusActive  UserStatus = "active"
	StatusPassive UserStatus = "passive"
	StatusBanned  UserStatus = "banned"

	RoleUser       UserRole = "USER"
	RoleEditor     UserRole = "EDITOR"
	RoleAdmin      UserRole = "ADMIN"
	RoleSuperAdmin UserRole = "SUPER_ADMIN"
)

type User struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Username      string     `json:"username" gorm:"unique;not null"`
	Email         string     `json:"email" gorm:"unique;not null"`
	Password      string     `json:"-" gorm:"not null"` // "-" json tag'i bu alanın JSON dönüşlerinde gizlenmesini sağlar
	Avatar        string     `json:"avatar" gorm:"default:'/uploads/default/avatar.png'"`
	Status        UserStatus `json:"status" gorm:"type:user_status;default:active"`
	Role          UserRole   `json:"role" gorm:"type:user_role;default:USER"`
	IsRootAdmin   bool       `json:"-" gorm:"default:false"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	LastLoginDate time.Time  `json:"last_login_date"`
	BanReason     string     `json:"ban_reason,omitempty" gorm:"type:text"`
	BanEndDate    *time.Time `json:"ban_end_date,omitempty"`
}

// TableName GORM için tablo adını belirtir
func (User) TableName() string {
	return "users"
}

// ChangePasswordRequest şifre değiştirme isteği için model
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
} 