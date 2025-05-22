package http

import (
	"net/http"
	"os"
	"time"

	"algo/internal/modules/oauth/application/dto"
	"algo/internal/modules/oauth/application/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OAuthController handles OAuth-related HTTP requests
type OAuthController struct {
	oauthService *service.OAuthService
}

// NewOAuthController creates a new OAuth controller
func NewOAuthController(oauthService *service.OAuthService) *OAuthController {
	return &OAuthController{
		oauthService: oauthService,
	}
}

// InitiateLogin initiates the OAuth login flow
// @Summary Initiate OAuth login
// @Description Start the OAuth login flow for a specific provider
// @Tags OAuth
// @Accept json
// @Produce json
// @Param provider path string true "OAuth provider (google, github, etc.)"
// @Success 200 {object} dto.ApiResponse{data=map[string]string} "Login URL"
// @Success 302 {string} string "Redirect to provider login page"
// @Failure 400 {object} dto.ApiResponse "Invalid provider"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /oauth/{provider}/login [get]
func (c *OAuthController) InitiateLogin(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider == "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Provider is required",
		})
		return
	}

	// Get login URL
	loginURL, err := c.oauthService.GetLoginURL(ctx.Request.Context(), provider)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get login URL: " + err.Error(),
		})
		return
	}

	// Check if this is an API request or a browser request
	if ctx.GetHeader("Content-Type") == "application/json" || ctx.GetHeader("Accept") == "application/json" {
		// API request, return JSON
		ctx.JSON(http.StatusOK, dto.ApiResponse{
			Success: true,
			Data:    loginURL,
		})
	} else {
		// Browser request, redirect
		ctx.Redirect(http.StatusFound, loginURL.URL)
	}
}

// HandleCallback handles the OAuth callback
// HandleCallback handles the OAuth callback
// @Summary Handle OAuth callback
// @Description Process the OAuth callback from the provider
// @Tags OAuth
// @Accept json
// @Produce json
// @Param provider path string true "OAuth provider (google, github, etc.)"
// @Param code query string true "Authorization code"
// @Param state query string true "State parameter for CSRF protection"
// @Param error query string false "Error message from provider"
// @Success 200 {object} dto.ApiResponse{data=map[string]interface{}} "User and token"
// @Success 302 {string} string "Redirect to home page with cookie"
// @Failure 400 {object} dto.ApiResponse "Invalid parameters"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /oauth/{provider}/callback [get]
func (c *OAuthController) HandleCallback(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider == "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Provider is required",
		})
		return
	}

	// Get callback params
	code := ctx.Query("code")
	state := ctx.Query("state")
	errorMsg := ctx.Query("error")

	if errorMsg != "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "OAuth error: " + errorMsg,
		})
		return
	}

	if code == "" || state == "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Code and state are required",
		})
		return
	}

	// Handle callback
	user, token, err := c.oauthService.HandleCallback(ctx.Request.Context(), provider, code, state)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to handle callback: " + err.Error(),
		})
		return
	}

	// Set JWT token as a cookie
	SetTokenCookie(ctx, token, user.ID.String())

	// Check if this is an API request or a browser request
	if ctx.GetHeader("Content-Type") == "application/json" || ctx.GetHeader("Accept") == "application/json" {
		// API request, return JSON
		ctx.JSON(http.StatusOK, dto.ApiResponse{
			Success: true,
			Data: gin.H{
				"token": token,
				"user":  user,
			},
		})
	} else {
		// Browser request, redirect to home page
		ctx.Redirect(http.StatusFound, "/")
	}
}

// GetAvailableProviders gets all available OAuth providers
// GetAvailableProviders gets all available OAuth providers
// @Summary Get available OAuth providers
// @Description Retrieve a list of all available OAuth providers
// @Tags OAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=[]string} "List of available providers"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /oauth/providers [get]
func (c *OAuthController) GetAvailableProviders(ctx *gin.Context) {
	providers, err := c.oauthService.GetAvailableProviders(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get providers: " + err.Error(),
		})
		return
	}

	// Add login URLs to providers
	for _, provider := range providers {
		loginURL, err := c.oauthService.GetLoginURL(ctx.Request.Context(), provider.Name)
		if err == nil {
			provider.LoginURL = loginURL.URL
		}
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    providers,
	})
}

// GetUserConnections gets all OAuth connections for a user
// GetUserConnections gets all OAuth connections for a user
// @Summary Get user OAuth connections
// @Description Retrieve all OAuth connections for the authenticated user
// @Tags OAuth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=[]dto.OAuthConnectionResponse} "List of user's OAuth connections"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /oauth/connections [get]
func (c *OAuthController) GetUserConnections(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	// Convert userID to UUID
	id, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	connections, err := c.oauthService.GetUserConnections(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get connections: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    connections,
	})
}

// DisconnectProvider disconnects a provider from a user
func (c *OAuthController) DisconnectProvider(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	// Convert userID to UUID
	id, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	provider := ctx.Param("provider")
	if provider == "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Provider is required",
		})
		return
	}

	err = c.oauthService.DisconnectProvider(ctx.Request.Context(), id, provider)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to disconnect provider: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
	})
}

// SetTokenCookie sets the JWT token as a cookie
func SetTokenCookie(ctx *gin.Context, token string, userID string) {
	expiresAt := time.Now().Add(24 * time.Hour) // 24-hour expiry

	// Get the current environment
	appEnv := getAppEnvironment()

	// Set secure flag based on environment (only in production)
	isSecure := appEnv == "production"

	// For localhost and development, use lax SameSite policy
	sameSiteMode := http.SameSiteLaxMode
	if appEnv == "production" {
		sameSiteMode = http.SameSiteStrictMode
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: sameSiteMode,
	}
	http.SetCookie(ctx.Writer, cookie)
}

// getAppEnvironment returns the current application environment
func getAppEnvironment() string {
	// Default to development if not set
	appEnv := "development"

	// You can add environment variable check here if needed
	if env := os.Getenv("APP_ENV"); env != "" {
		appEnv = env
	}

	return appEnv
}
