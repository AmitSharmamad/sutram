package model

import (
	"time"

	domainEntity "algo/internal/domain/entity"

	"github.com/google/uuid"
)

// User represents the user domain model specific to the user module
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

// Role represents a role in the user module
type Role struct {
	ID          int          `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	Permissions []Permission `json:"permissions,omitempty" db:"-"` // Permissions are loaded separately
}

// Permission represents a permission in the user module
type Permission struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	ResourceType string    `json:"resource_type" db:"resource_type"`
	Action       string    `json:"action" db:"action"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// FromDomainEntity converts a domain entity to a module user model
func FromDomainEntity(entity *domainEntity.User) *User {
	if entity == nil {
		return nil
	}

	user := &User{
		ID:          entity.ID,
		Name:        entity.Name,
		Username:    entity.Username,
		Email:       entity.Email,
		PhoneNumber: entity.PhoneNumber,
		Password:    entity.Password,
		Verified:    entity.Verified,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		IsActive:    entity.IsActive,
	}

	if len(entity.Roles) > 0 {
		user.Roles = make([]Role, len(entity.Roles))
		for i, role := range entity.Roles {
			user.Roles[i] = Role{
				ID:          role.ID,
				Name:        string(role.Name),
				Description: role.Description,
				CreatedAt:   role.CreatedAt,
				UpdatedAt:   role.UpdatedAt,
			}
		}
	}

	return user
}

// ToDomainEntity converts a module user model to a domain entity
func (u *User) ToDomainEntity() *domainEntity.User {
	if u == nil {
		return nil
	}

	user := &domainEntity.User{
		ID:          u.ID,
		Name:        u.Name,
		Username:    u.Username,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Password:    u.Password,
		Verified:    u.Verified,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		IsActive:    u.IsActive,
	}

	if len(u.Roles) > 0 {
		user.Roles = make([]domainEntity.Role, len(u.Roles))
		for i, role := range u.Roles {
			user.Roles[i] = domainEntity.Role{
				ID:          role.ID,
				Name:        domainEntity.RoleType(role.Name),
				Description: role.Description,
				CreatedAt:   role.CreatedAt,
				UpdatedAt:   role.UpdatedAt,
			}
		}
	}

	return user
}
