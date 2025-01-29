package services

import (
	"context"
	"errors"
	"log"
	"reflect"
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
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	user.Password = string(hashedPassword)

	// Check if username exists
	var existingUser models.User
	if err := s.db.Where("username = ? AND deleted_at IS NULL", user.Username).First(&existingUser).Error; err == nil {
		log.Printf("Registration failed: Username '%s' is already taken", user.Username)
		return errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error checking username existence: %v", err)
		return err
	}
	log.Printf("Username '%s' is available for registration", user.Username)

	// Check if email exists
	if err := s.db.Where("email = ? AND deleted_at IS NULL", user.Email).First(&existingUser).Error; err == nil {
		log.Printf("Registration failed: Email '%s' is already registered", user.Email)
		return errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error checking email existence: %v", err)
		return err
	}
	log.Printf("Email '%s' is available for registration", user.Email)

	// Set default values
	user.CreatedAt = time.Now()
	user.LastLoginDate = time.Now()

	// Create user
	if err := s.db.Create(user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	log.Printf("User registered successfully: %s (%s)", user.Username, user.Email)
	return nil
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

func (s *AuthService) AdminLogin(identifier, password string) (*models.User, error) {
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

	// Check if user is admin
	if user.Role != models.RoleAdmin && user.Role != models.RoleSuperAdmin {
		return nil, ErrUnauthorized
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

func (s *AuthService) GetAdminProfile(userID uint, profile interface{}) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Check if user is admin
	if user.Role != models.RoleAdmin && user.Role != models.RoleSuperAdmin {
		return errors.New("admin privileges required")
	}

	// Map user data to profile
	val := reflect.ValueOf(profile).Elem()
	userVal := reflect.ValueOf(user)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if userField := userVal.FieldByName(field.Name); userField.IsValid() {
			val.Field(i).Set(userField)
		}
	}

	return nil
} 