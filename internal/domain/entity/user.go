package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity in the system
type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Password    string    `json:"-" db:"password"` // Never expose password in JSON
	Verified    bool      `json:"verified" db:"verified"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Roles       []Role    `json:"roles,omitempty" db:"-"` // Roles are loaded separately
}

// Role represents a role entity in the system
type Role struct {
	ID          int       `json:"id" db:"id"`
	Name        RoleType  `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RoleType is an enum representing the possible role types
type RoleType string

// Role type constants
const (
	RoleAdmin     RoleType = "admin"
	RolePaidUser  RoleType = "paid_user"
	RoleFreeUser  RoleType = "free_user"
	RoleModerator RoleType = "moderator"
	RoleAnalyst   RoleType = "analyst"
)

// Permission represents a permission entity in the system
type Permission struct {
	ID           int            `json:"id" db:"id"`
	Name         PermissionType `json:"name" db:"name"`
	Description  string         `json:"description" db:"description"`
	ResourceType string         `json:"resource_type" db:"resource_type"`
	Action       string         `json:"action" db:"action"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
}

// PermissionType is an enum representing the possible permission types
type PermissionType string

// Permission type constants
const (
	// Strategy permissions
	PermCreateStrategy    PermissionType = "create_strategy"
	PermViewAnyStrategy   PermissionType = "view_any_strategy"
	PermEditOwnStrategy   PermissionType = "edit_own_strategy"
	PermDeleteOwnStrategy PermissionType = "delete_own_strategy"

	// Backtest permissions
	PermRunBacktest     PermissionType = "run_backtest"
	PermViewOwnBacktest PermissionType = "view_own_backtest"
	PermViewAnyBacktest PermissionType = "view_any_backtest"

	// Marketplace permissions
	PermPublishStrategy  PermissionType = "publish_strategy"
	PermPurchaseStrategy PermissionType = "purchase_strategy"

	// Paper trading
	PermRunPaperTrade PermissionType = "run_paper_trade"

	// Live trading
	PermRunLiveTrade PermissionType = "run_live_trade"

	// Community permissions
	PermCreateTopic     PermissionType = "create_topic"
	PermModerateContent PermissionType = "moderate_content"
)

// NewUser creates a new user with default values
func NewUser(name, username, email, phoneNumber, password string) (*User, error) {
	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:          uuid.New(),
		Name:        name,
		Username:    username,
		Email:       email,
		PhoneNumber: phoneNumber,
		Password:    string(hashedPassword),
		Verified:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsActive:    true,
	}, nil
}

// ComparePassword compares the provided password with the stored hash
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(roleName RoleType) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// SetVerified marks the user as verified
func (u *User) SetVerified() {
	u.Verified = true
	u.UpdatedAt = time.Now()
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}
