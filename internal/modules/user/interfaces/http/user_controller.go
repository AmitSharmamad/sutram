package http

import (
	"net/http"
	"time"

	"algo/internal/modules/user/application/dto"
	"algo/internal/modules/user/application/service"
	"algo/internal/modules/user/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController handles user-related HTTP requests
type UserController struct {
	authService *service.AuthService
	userService *service.UserService
}

// NewUserController creates a new user controller
func NewUserController(authService *service.AuthService, userService *service.UserService) *UserController {
	return &UserController{
		authService: authService,
		userService: userService,
	}
}

// Register handles user registration
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Register the user
	authResp, err := c.authService.Register(
		ctx.Request.Context(),
		req.Name,
		req.Username,
		req.Email,
		req.PhoneNumber,
		req.Password,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Set JWT token as a cookie
	setTokenCookie(ctx, authResp.Token, authResp.ExpiresAt)

	// Return success response with token and user data
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data: gin.H{
			"token":      authResp.Token,
			"expires_at": authResp.ExpiresAt,
			"user":       authResp.User,
		},
	})
}

// Login handles user login
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Login the user
	authResp, err := c.authService.Login(
		ctx.Request.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Authentication failed: " + err.Error(),
		})
		return
	}

	// Set JWT token as a cookie
	setTokenCookie(ctx, authResp.Token, authResp.ExpiresAt)

	// Check if this is an API request or a browser request
	if ctx.GetHeader("Content-Type") == "application/json" || ctx.GetHeader("Accept") == "application/json" {
		ctx.Redirect(http.StatusFound, "/home")
	} else {
		// Redirect to home page for browser clients
		ctx.Redirect(http.StatusFound, "/home")
	}
}

// RefreshToken refreshes a user's token
func (c *UserController) RefreshToken(ctx *gin.Context) {
	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Refresh the token
	authResp, err := c.authService.RefreshToken(ctx.Request.Context(), tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Token refresh failed: " + err.Error(),
		})
		return
	}

	// Set JWT token as a cookie
	setTokenCookie(ctx, authResp.Token, authResp.ExpiresAt)

	// Return success response with new token
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data: gin.H{
			"token":     authResp.Token,
			"expiresAt": authResp.ExpiresAt,
		},
	})
}

// GetProfile gets the current user's profile
func (c *UserController) GetProfile(ctx *gin.Context) {
	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Get the user from token
	user, err := c.authService.GetUserFromToken(ctx.Request.Context(), tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	// Convert user model to DTO
	userDTO := modelToDTO(user)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    userDTO,
	})
}

// UpdateProfile updates the current user's profile
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Get user ID from token
	userID, err := c.authService.GetUserIDFromToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	// Update the user profile
	user, err := c.userService.UpdateUser(
		ctx.Request.Context(),
		userID,
		req.Name,
		req.Username,
		req.Email,
		req.PhoneNumber,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Failed to update profile: " + err.Error(),
		})
		return
	}

	// Convert user model to DTO
	userDTO := modelToDTO(user)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    userDTO,
	})
}

// ChangePassword changes the current user's password
func (c *UserController) ChangePassword(ctx *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Get user ID from token
	userID, err := c.authService.GetUserIDFromToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	// Change the password
	authResp, err := c.authService.ChangePassword(
		ctx.Request.Context(),
		userID,
		req.CurrentPassword,
		req.NewPassword,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Failed to change password: " + err.Error(),
		})
		return
	}

	// Convert user model to DTO
	userDTO := modelToDTO(authResp.User)

	// Set JWT token as a cookie
	setTokenCookie(ctx, authResp.Token, authResp.ExpiresAt)

	// Create response
	response := dto.AuthResponse{
		Token:     authResp.Token,
		User:      userDTO,
		ExpiresAt: authResp.ExpiresAt.Unix(),
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Password changed successfully",
		Data:    response,
	})
}

// GetUserByID gets a user by ID (admin only)
func (c *UserController) GetUserByID(ctx *gin.Context) {
	// Get user ID from URL parameter
	userIDParam := ctx.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Check if user has admin role
	hasAdminRole, err := c.authService.HasRole(tokenString, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	if !hasAdminRole {
		ctx.JSON(http.StatusForbidden, dto.ApiResponse{
			Success: false,
			Error:   "Permission denied: Admin role required",
		})
		return
	}

	// Get the user
	user, err := c.userService.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ApiResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Convert user model to DTO
	userDTO := modelToDTO(user)

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    userDTO,
	})
}

// AssignRole assigns a role to a user (admin only)
func (c *UserController) AssignRole(ctx *gin.Context) {
	var req dto.AssignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from URL parameter
	userIDParam := ctx.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Check if user has admin role
	hasAdminRole, err := c.authService.HasRole(tokenString, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	if !hasAdminRole {
		ctx.JSON(http.StatusForbidden, dto.ApiResponse{
			Success: false,
			Error:   "Permission denied: Admin role required",
		})
		return
	}

	// Assign the role
	err = c.userService.AssignRoleToUser(ctx.Request.Context(), userID, req.RoleName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Failed to assign role: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Role assigned successfully",
	})
}

// RemoveRole removes a role from a user (admin only)
func (c *UserController) RemoveRole(ctx *gin.Context) {
	var req dto.AssignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from URL parameter
	userIDParam := ctx.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Check if user has admin role
	hasAdminRole, err := c.authService.HasRole(tokenString, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	if !hasAdminRole {
		ctx.JSON(http.StatusForbidden, dto.ApiResponse{
			Success: false,
			Error:   "Permission denied: Admin role required",
		})
		return
	}

	// Remove the role
	err = c.userService.RemoveRoleFromUser(ctx.Request.Context(), userID, req.RoleName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Failed to remove role: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "Role removed successfully",
	})
}

// GetAllRoles gets all available roles (admin only)
func (c *UserController) GetAllRoles(ctx *gin.Context) {
	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Check if user has admin role
	hasAdminRole, err := c.authService.HasRole(tokenString, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	if !hasAdminRole {
		ctx.JSON(http.StatusForbidden, dto.ApiResponse{
			Success: false,
			Error:   "Permission denied: Admin role required",
		})
		return
	}

	// Get all roles
	roles, err := c.userService.GetAllRoles(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get roles: " + err.Error(),
		})
		return
	}

	// Convert to simple format
	roleNames := make([]map[string]interface{}, len(roles))
	for i, role := range roles {
		roleNames[i] = map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"description": role.Description,
		}
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    roleNames,
	})
}

// GetUserRoles gets all roles for a user
func (c *UserController) GetUserRoles(ctx *gin.Context) {
	// Get user ID from URL parameter
	userIDParam := ctx.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Get authenticated user ID
	authUserID, err := c.authService.GetUserIDFromToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	// Check if user is requesting their own roles or has admin role
	isOwnRoles := authUserID == userID
	hasAdminRole, err := c.authService.HasRole(tokenString, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	if !isOwnRoles && !hasAdminRole {
		ctx.JSON(http.StatusForbidden, dto.ApiResponse{
			Success: false,
			Error:   "Permission denied",
		})
		return
	}

	// Get user roles
	roles, err := c.userService.GetUserRoles(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get user roles: " + err.Error(),
		})
		return
	}

	// Convert to simple format
	roleNames := make([]map[string]interface{}, len(roles))
	for i, role := range roles {
		roleNames[i] = map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"description": role.Description,
		}
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    roleNames,
	})
}

// ActivateUser activates a user account (admin only)
func (c *UserController) ActivateUser(ctx *gin.Context) {
	// Get user ID from URL parameter
	userIDParam := ctx.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Check if user has admin role
	hasAdminRole, err := c.authService.HasRole(tokenString, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	if !hasAdminRole {
		ctx.JSON(http.StatusForbidden, dto.ApiResponse{
			Success: false,
			Error:   "Permission denied: Admin role required",
		})
		return
	}

	// Activate the user
	err = c.userService.ActivateUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Failed to activate user: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "User activated successfully",
	})
}

// DeactivateUser deactivates a user account (admin only)
func (c *UserController) DeactivateUser(ctx *gin.Context) {
	// Get user ID from URL parameter
	userIDParam := ctx.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get the token from header
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "No token provided",
		})
		return
	}

	// Check if user has admin role
	hasAdminRole, err := c.authService.HasRole(tokenString, "admin")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid token: " + err.Error(),
		})
		return
	}

	if !hasAdminRole {
		ctx.JSON(http.StatusForbidden, dto.ApiResponse{
			Success: false,
			Error:   "Permission denied: Admin role required",
		})
		return
	}

	// Deactivate the user
	err = c.userService.DeactivateUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Failed to deactivate user: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "User deactivated successfully",
	})
}

// CheckMe checks if the user is logged in by verifying the token cookie
func (c *UserController) CheckMe(ctx *gin.Context) {
	// Get the token from header or cookie
	tokenString := extractToken(ctx)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Not logged in",
		})
		return
	}

	// Validate the token
	_, err := c.authService.ValidateToken(ctx.Request.Context(), tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "Invalid session: " + err.Error(),
		})
		return
	}

	// If token is valid, return success
	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Message: "User is logged in",
	})
}

// ExtractToken extracts the JWT token from the Authorization header or cookie
// Returns the token without the "Bearer " prefix or empty string if not found
// This function is exported so it can be used by other modules
func ExtractToken(ctx *gin.Context) string {
	// First, try to get token from Authorization header
	// bearerToken := ctx.GetHeader("Authorization")
	// if len(bearerToken) > 7 && strings.HasPrefix(strings.ToUpper(bearerToken), "BEARER ") {
	// 	return bearerToken[7:]
	// }

	// If not in header, try to get from cookie
	cookie, err := ctx.Request.Cookie("token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

// For backward compatibility, keep the unexported version
func extractToken(ctx *gin.Context) string {
	return ExtractToken(ctx)
}

// Helper function to get user ID from token
func (c *UserController) GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	return c.authService.GetUserIDFromToken(tokenString)
}

// Helper function to set token cookie
func setTokenCookie(ctx *gin.Context, token string, expiresAt time.Time) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(ctx.Writer, cookie)
}

// Logout handles user logout by clearing the token cookie
func (c *UserController) Logout(ctx *gin.Context) {
	// Create an expired cookie with the same attributes to clear it
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Set expiry in the past
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Instruct browser to delete the cookie immediately
	}
	http.SetCookie(ctx.Writer, cookie)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// Helper function to convert user model to DTO
func modelToDTO(user *model.User) *dto.UserDTO {
	if user == nil {
		return nil
	}

	// Extract role names
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	return &dto.UserDTO{
		ID:          user.ID.String(),
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Verified:    user.Verified,
		IsActive:    user.IsActive,
		Roles:       roleNames,
	}
}

// GetAuthService returns the auth service instance
func (c *UserController) GetAuthService() *service.AuthService {
	return c.authService
}
