package services

import (
	"context"
	"errors"
	"time"

	"github.com/anilsoylu/answer-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUsernameTaken      = errors.New("username is already taken")
	ErrEmailTaken         = errors.New("email is already taken")
	ErrUserFrozen         = errors.New("this account is frozen, please contact support")
	ErrUserBanned         = errors.New("this account is banned")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotActive       = errors.New("user is not active")
	ErrUserNotFound        = errors.New("user not found")
	ErrDuplicateEntry     = errors.New("duplicate entry")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(user *models.User) error {
	// Check if username is taken by an active or banned user
	var existingUser models.User
	if err := s.db.Where("username = ? AND deleted_at IS NULL", user.Username).First(&existingUser).Error; err == nil {
		if existingUser.Status == models.StatusBanned {
			return ErrUserBanned
		}
		return ErrUsernameTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check if email is taken by an active or banned user
	if err := s.db.Where("email = ? AND deleted_at IS NULL", user.Email).First(&existingUser).Error; err == nil {
		if existingUser.Status == models.StatusBanned {
			return ErrUserBanned
		}
		return ErrEmailTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Create user with transaction
	tx := s.db.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Reload user to get the correct ID
	if err := tx.First(user, user.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AuthService) Login(identifier, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if user.Status != models.StatusActive {
		return nil, ErrUserNotActive
	}

	// Update last login date
	user.LastLoginDate = time.Now()
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) UpdateUserRole(userID uint, newRole models.UserRole, requester *models.User) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check permissions and role hierarchy
	if user.Role == models.RoleSuperAdmin && !requester.IsRootAdmin {
		return errors.New("only root admin can change SUPER_ADMIN's role")
	}

	if newRole == models.RoleSuperAdmin && !requester.IsRootAdmin {
		return errors.New("only root admin can assign SUPER_ADMIN role")
	}

	if requester.Role == models.RoleAdmin {
		if user.Role != models.RoleUser && user.Role != models.RoleEditor {
			return errors.New("admin can only modify USER and EDITOR roles")
		}
	}

	user.Role = newRole
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdateUserStatus(userID uint, newStatus models.UserStatus, requesterRole models.UserRole) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check permissions
	if user.Role == models.RoleSuperAdmin {
		return ErrForbidden
	}

	// Kullanıcı kendi durumunu güncelleyebilir
	if requesterRole == models.RoleUser || requesterRole == models.RoleEditor {
		// Normal kullanıcılar sadece kendi durumlarını active ve passive yapabilir
		if newStatus == models.StatusBanned {
			return ErrForbidden
		}
	} else if requesterRole == models.RoleAdmin {
		// Admin başka bir admin'in durumunu değiştiremez
		if user.Role == models.RoleAdmin {
			return ErrForbidden
		}
	} else if requesterRole != models.RoleSuperAdmin {
		return ErrUnauthorized
	}

	user.Status = newStatus
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdateProfile(userID uint, username, email, avatar string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Check for duplicate username/email
	if username != "" && username != user.Username {
		var existingUser models.User
		if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
			return nil, ErrDuplicateEntry
		}
		user.Username = username
	}

	if email != "" && email != user.Email {
		var existingUser models.User
		if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
			return nil, ErrDuplicateEntry
		}
		user.Email = email
	}

	if avatar != "" {
		user.Avatar = avatar
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash and save new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetUserByID(userID uint, user *models.User) error {
	if err := s.db.First(user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

// UpdatePassword kullanıcı şifresini günceller
func (s *AuthService) UpdatePassword(ctx context.Context, userID uint, req models.UpdatePasswordRequest) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Mevcut şifrenin doğruluğunu kontrol et
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Yeni şifreyi hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Şifreyi güncelle
	user.Password = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// BanUser kullanıcıyı banlar
func (s *AuthService) BanUser(ctx context.Context, userID uint, banReason, banDuration string, requesterRole models.UserRole) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Yetki kontrolü
	if requesterRole != models.RoleAdmin && requesterRole != models.RoleSuperAdmin {
		return ErrUnauthorized
	}

	// Admin, diğer adminleri ve super adminleri banlayamaz
	if requesterRole == models.RoleAdmin && (user.Role == models.RoleAdmin || user.Role == models.RoleSuperAdmin) {
		return ErrForbidden
	}

	// Super admin'i kimse banlayamaz
	if user.Role == models.RoleSuperAdmin {
		return ErrForbidden
	}

	// Ban süresini hesapla
	var banEndDate *time.Time
	now := time.Now()

	switch banDuration {
	case "1_day":
		end := now.Add(24 * time.Hour)
		banEndDate = &end
	case "1_week":
		end := now.Add(7 * 24 * time.Hour)
		banEndDate = &end
	case "1_month":
		end := now.AddDate(0, 1, 0)
		banEndDate = &end
	case "permanent":
		banEndDate = nil
	}

	// Kullanıcıyı banla
	user.Status = models.StatusBanned
	user.BanReason = banReason
	user.BanEndDate = banEndDate

	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// FreezeAccount hesabı dondurur (soft delete)
func (s *AuthService) FreezeAccount(ctx context.Context, userID uint, freezeReason string, requesterRole models.UserRole) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Super admin hesabı dondurulamaz
	if user.Role == models.RoleSuperAdmin {
		return ErrForbidden
	}

	// Hesap dondurma sebebi zorunlu
	if freezeReason == "" {
		return errors.New("freeze reason is required")
	}

	// Soft delete işlemi
	if err := s.db.Delete(&user).Error; err != nil {
		return err
	}

	// Dondurma bilgilerini güncelle
	user.Status = models.StatusFrozen
	user.FrozenReason = freezeReason
	now := time.Now()
	user.FrozenDate = &now

	// Soft delete edilmiş kaydı güncelle
	if err := s.db.Unscoped().Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// DeleteAccount hesabı siler (soft delete)
func (s *AuthService) DeleteAccount(ctx context.Context, userID uint, requesterID uint, requesterRole models.UserRole) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Kullanıcı kendi hesabını silebilir veya SUPER_ADMIN başka hesapları silebilir
	if userID != requesterID && requesterRole != models.RoleSuperAdmin {
		return ErrUnauthorized
	}

	// Super admin hesabı silinemez
	if user.Role == models.RoleSuperAdmin {
		return ErrForbidden
	}

	// Kullanıcıyı soft delete yap
	if err := s.db.Delete(&user).Error; err != nil {
		return err
	}

	// Status'u güncelle
	user.Status = models.StatusFrozen
	if err := s.db.Unscoped().Save(&user).Error; err != nil {
		return err
	}

	return nil
} 