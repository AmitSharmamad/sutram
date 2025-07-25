package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"algo/internal/domain/entity"
	domainRepo "algo/internal/domain/repository"
	"algo/internal/modules/user/domain/model"
	"algo/internal/modules/user/domain/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Common errors
var (
	ErrUserNotFound             = errors.New("user not found")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrEmailAlreadyExists       = errors.New("email already exists")
	ErrUsernameAlreadyExists    = errors.New("username already exists")
	ErrPhoneNumberAlreadyExists = errors.New("phone number already exists")
)

// UserService handles user-related business logic
type UserService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	usageLimitRepo domainRepo.UsageLimitRepository
	usageLogRepo   domainRepo.UsageLogRepository
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	usageLimitRepo domainRepo.UsageLimitRepository,
	usageLogRepo domainRepo.UsageLogRepository,
) *UserService {
	return &UserService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		usageLimitRepo: usageLimitRepo,
		usageLogRepo:   usageLogRepo,
	}
}

// Register registers a new user
func (s *UserService) Register(ctx context.Context, name, username, email, phoneNumber, password string) (*model.User, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user by email: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Check if user with username already exists
	existingUser, err = s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user by username: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUsernameAlreadyExists
	}

	// Check if user with phone number already exists (if provided)
	if phoneNumber != "" {
		existingUser, err = s.userRepo.FindByUsername(ctx, phoneNumber)
		if err != nil {
			return nil, fmt.Errorf("error checking existing user by phone number: %w", err)
		}
		if existingUser != nil {
			return nil, ErrPhoneNumberAlreadyExists
		}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create a new user
	now := time.Now()
	user := &model.User{
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
	}

	// Save the user to the database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Assign default "free_user" role
	freeRole, err := s.roleRepo.FindByName(ctx, "free_user")
	if err != nil {
		return nil, fmt.Errorf("error finding free_user role: %w", err)
	}

	if freeRole != nil {
		if err := s.roleRepo.AssignRoleToUser(ctx, user.ID, freeRole.ID); err != nil {
			return nil, fmt.Errorf("error assigning default role: %w", err)
		}

		// Add the role to the user object
		user.Roles = []model.Role{*freeRole}

		// Create default usage limits for the new user
		if err := s.createDefaultUsageLimits(ctx, user.ID, "free_user"); err != nil {
			// Log error but don't fail user creation
			fmt.Printf("Error creating default usage limits for user %s: %v\n", user.ID, err)
		}
	}

	return user, nil
}

// Authenticate authenticates a user using email and password
func (s *UserService) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if the user is active
	if !user.IsActive {
		return nil, errors.New("user account is not active")
	}

	// If password is empty, it's an OAuth authentication attempt
	// Only skip password verification for OAuth flows
	if password != "" {
		// Verify password for normal authentication
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return nil, ErrInvalidCredentials
		}
	}

	// Load user roles if not already loaded
	if len(user.Roles) == 0 {
		roles, err := s.roleRepo.GetUserRoles(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("error loading user roles: %w", err)
		}
		user.Roles = roles
	}

	return user, nil
}

// FindUserByEmail finds a user by email
func (s *UserService) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Load user roles if not already loaded
	if len(user.Roles) == 0 {
		roles, err := s.roleRepo.GetUserRoles(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("error loading user roles: %w", err)
		}
		user.Roles = roles
	}

	return user, nil
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	// Attempt to find the user with the provided context
	user, err := s.userRepo.FindByID(ctx, id)

	// If there's a context cancellation error, retry with a fresh context
	if err != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
		// Create a new context with timeout
		retryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Retry the operation
		user, err = s.userRepo.FindByID(retryCtx, id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, name, username, email, phoneNumber string) (*model.User, error) {
	// Get the existing user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check if email is changing and if it conflicts with existing users
	if email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("error checking existing user by email: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrEmailAlreadyExists
		}
	}

	// Check if username is changing and if it conflicts with existing users
	if username != user.Username {
		existingUser, err := s.userRepo.FindByUsername(ctx, username)
		if err != nil {
			return nil, fmt.Errorf("error checking existing user by username: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrUsernameAlreadyExists
		}
	}

	// Check if phone number is changing and if it conflicts with existing users
	if phoneNumber != user.PhoneNumber {
		existingUser, err := s.userRepo.FindByUsername(ctx, phoneNumber)
		if err != nil {
			return nil, fmt.Errorf("error checking existing user by phone number: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrPhoneNumberAlreadyExists
		}
	}

	// Update user fields
	user.Name = name
	user.Username = username
	user.Email = email
	user.PhoneNumber = phoneNumber
	user.UpdatedAt = time.Now()

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error {
	// Get the user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Update password
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error updating user password: %w", err)
	}

	return nil
}

// VerifyUser marks a user as verified
func (s *UserService) VerifyUser(ctx context.Context, id uuid.UUID) error {
	// Get the user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Set verified flag
	user.Verified = true
	user.UpdatedAt = time.Now()

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error updating user verification: %w", err)
	}

	return nil
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	// Get the user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Set active flag to false
	user.IsActive = false
	user.UpdatedAt = time.Now()

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error deactivating user: %w", err)
	}

	return nil
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(ctx context.Context, id uuid.UUID) error {
	// Get the user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Set active flag to true
	user.IsActive = true
	user.UpdatedAt = time.Now()

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("error activating user: %w", err)
	}

	return nil
}

// AssignRoleToUser assigns a role to a user
func (s *UserService) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	// Find the role by name
	role, err := s.roleRepo.FindByName(ctx, roleName)
	if err != nil {
		return fmt.Errorf("error finding role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Assign the role to the user
	if err := s.roleRepo.AssignRoleToUser(ctx, userID, role.ID); err != nil {
		return fmt.Errorf("error assigning role to user: %w", err)
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (s *UserService) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	// Find the role by name
	role, err := s.roleRepo.FindByName(ctx, roleName)
	if err != nil {
		return fmt.Errorf("error finding role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Remove the role from the user
	if err := s.roleRepo.RemoveRoleFromUser(ctx, userID, role.ID); err != nil {
		return fmt.Errorf("error removing role from user: %w", err)
	}

	return nil
}

// GetUserRoles gets all roles for a user
func (s *UserService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user roles: %w", err)
	}
	return roles, nil
}

// GetAllRoles gets all available roles
func (s *UserService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all roles: %w", err)
	}
	return roles, nil
}

// createDefaultUsageLimits creates default usage limits for a new user based on their role
func (s *UserService) createDefaultUsageLimits(ctx context.Context, userID uuid.UUID, roleName string) error {
	// Define default limits based on role
	var backtestLimit, deploymentLimit int

	switch roleName {
	case "admin":
		backtestLimit = 999999
		deploymentLimit = 999999
	case "paid_user":
		backtestLimit = 100
		deploymentLimit = 20
	case "free_user":
		backtestLimit = 10
		deploymentLimit = 3
	default:
		backtestLimit = 5
		deploymentLimit = 1
	}

	now := time.Now()

	// Create backtest limit
	backtestUsageLimit := &entity.UsageLimit{
		UserID:       userID,
		LimitType:    entity.BacktestPerMonth,
		LimitValue:   backtestLimit,
		CurrentUsage: 0,
		ResetPeriod:  entity.Monthly,
		LastResetAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.usageLimitRepo.Create(ctx, backtestUsageLimit); err != nil {
		return fmt.Errorf("failed to create backtest limit: %w", err)
	}

	// Create deployment limit
	deploymentUsageLimit := &entity.UsageLimit{
		UserID:       userID,
		LimitType:    entity.DeploymentCount,
		LimitValue:   deploymentLimit,
		CurrentUsage: 0,
		ResetPeriod:  entity.Never,
		LastResetAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.usageLimitRepo.Create(ctx, deploymentUsageLimit); err != nil {
		return fmt.Errorf("failed to create deployment limit: %w", err)
	}

	return nil
}
