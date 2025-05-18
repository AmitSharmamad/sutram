package dto

// UserDTO represents a user data transfer object
type UserDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	Verified    bool     `json:"verified"`
	IsActive    bool     `json:"is_active"`
	Roles       []string `json:"roles,omitempty"`
}

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Password    string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest represents the request payload for updating user profile
type UpdateUserRequest struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// ChangePasswordRequest represents the request payload for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token     string   `json:"token"`
	User      *UserDTO `json:"user"`
	ExpiresAt int64    `json:"expires_at"` // Unix timestamp
}

// AssignRoleRequest represents the request payload for assigning a role to a user
type AssignRoleRequest struct {
	RoleName string `json:"role_name" binding:"required"`
}

// ApiResponse is a generic wrapper for API responses
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
