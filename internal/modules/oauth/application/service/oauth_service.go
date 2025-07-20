package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	userEntity "algo/internal/domain/entity"
	"algo/internal/modules/oauth/application/dto"
	"algo/internal/modules/oauth/application/service/provider"
	"algo/internal/modules/oauth/domain/model"
	"algo/internal/modules/oauth/domain/repository"
	userService "algo/internal/modules/user/application/service"

	"github.com/google/uuid"
)

// OAuthService handles OAuth-related functionality
type OAuthService struct {
	providerRepo  repository.OAuthProviderRepository
	userOAuthRepo repository.UserOAuthRepository
	userService   *userService.UserService
	authService   *userService.AuthService
	providers     map[string]provider.OAuthProvider
	stateStore    map[string]time.Time // Simple in-memory store for state tokens
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(
	providerRepo repository.OAuthProviderRepository,
	userOAuthRepo repository.UserOAuthRepository,
	userService *userService.UserService,
	authService *userService.AuthService,
) *OAuthService {
	service := &OAuthService{
		providerRepo:  providerRepo,
		userOAuthRepo: userOAuthRepo,
		userService:   userService,
		authService:   authService,
		providers:     make(map[string]provider.OAuthProvider),
		stateStore:    make(map[string]time.Time),
	}

	// Register providers
	service.providers["google"] = provider.NewGoogleOAuthProvider()

	// Start a goroutine to clean up expired state tokens
	go service.cleanupExpiredStates()

	return service
}

// isEmailWhitelisted checks if an email is in the whitelist
func (s *OAuthService) isEmailWhitelisted(email string) bool {
	// Whitelist of allowed emails
	allowedEmails := []string{
		"admin@sutram.co",
		"amitsharmamad@gmail.com",
		"Narendra.polisetty09@gmail.com",
		"kkmarella@gmail.com",
		"ajitsharma6329@gmail.com",
		"koushikmarella@gmail.com",
	}

	// Check if the email is in the whitelist
	for _, allowedEmail := range allowedEmails {
		if strings.EqualFold(email, allowedEmail) {
			return true
		}
	}

	return false
}

// cleanupExpiredStates periodically removes expired state tokens
func (s *OAuthService) cleanupExpiredStates() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for state, expiresAt := range s.stateStore {
			if now.After(expiresAt) {
				delete(s.stateStore, state)
			}
		}
	}
}

// GenerateState generates a random state token for OAuth flow
func (s *OAuthService) GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)
	s.stateStore[state] = time.Now().Add(10 * time.Minute) // State valid for 10 minutes
	return state, nil
}

// ValidateState validates the state token
func (s *OAuthService) ValidateState(state string) bool {
	expiresAt, ok := s.stateStore[state]
	if !ok {
		return false
	}

	if time.Now().After(expiresAt) {
		delete(s.stateStore, state)
		return false
	}

	// State is valid, remove it from store to prevent reuse
	delete(s.stateStore, state)
	return true
}

// GetLoginURL returns the URL to redirect the user to for login
func (s *OAuthService) GetLoginURL(ctx context.Context, providerName string) (*dto.OAuthLoginURL, error) {
	provider := s.providers[providerName]
	if provider == nil {
		return nil, fmt.Errorf("provider %s not supported", providerName)
	}

	state, err := s.GenerateState()
	if err != nil {
		return nil, err
	}

	url := provider.GetLoginURL(state)
	return &dto.OAuthLoginURL{
		URL:        url,
		Provider:   providerName,
		StateToken: state,
	}, nil
}

// HandleCallback handles the OAuth callback
func (s *OAuthService) HandleCallback(ctx context.Context, providerName, code, state string) (*userEntity.User, string, error) {
	// Create a timeout context to ensure the operation doesn't hang indefinitely
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate state
	if !s.ValidateState(state) {
		return nil, "", errors.New("invalid state token")
	}

	// Get provider
	provider := s.providers[providerName]
	if provider == nil {
		return nil, "", fmt.Errorf("provider %s not supported", providerName)
	}

	// Get user info from provider
	oauthResp, err := provider.GetUserInfo(timeoutCtx, code)
	if err != nil {
		return nil, "", err
	}

	// Check if the email is whitelisted
	if !s.isEmailWhitelisted(oauthResp.Email) {
		return nil, "", fmt.Errorf("email %s is not authorized for this application", oauthResp.Email)
	}

	// Get provider ID
	dbProvider, err := s.providerRepo.GetByName(timeoutCtx, providerName)
	if err != nil {
		return nil, "", err
	}
	if dbProvider == nil {
		return nil, "", fmt.Errorf("provider %s not found in database", providerName)
	}

	// Check if user is already connected
	userOAuth, err := s.userOAuthRepo.GetByProviderUserID(timeoutCtx, dbProvider.ID, oauthResp.ProviderID)
	if err != nil {
		return nil, "", err
	}

	var entityUser *userEntity.User
	var jwtToken string

	if userOAuth != nil {
		// User already exists, update tokens
		err = s.userOAuthRepo.UpdateTokens(
			timeoutCtx,
			userOAuth.ID,
			oauthResp.AccessToken,
			oauthResp.RefreshToken,
			oauthResp.ExpiresAt.Unix(),
		)
		if err != nil {
			return nil, "", err
		}

		// Get user with a fresh context to avoid cancellation
		getUserCtx, getUserCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer getUserCancel()

		user, err := s.userService.GetUserByID(getUserCtx, userOAuth.UserID)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				// If the context was canceled, retry with a fresh context
				retryCtx, retryCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer retryCancel()
				user, err = s.userService.GetUserByID(retryCtx, userOAuth.UserID)
				if err != nil {
					return nil, "", fmt.Errorf("error finding user after retry: %w", err)
				}
			} else {
				return nil, "", fmt.Errorf("error finding user: %w", err)
			}
		}

		// Convert model.User to entity.User
		entityUser = &userEntity.User{
			ID:          user.ID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Verified:    user.Verified,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			IsActive:    user.IsActive,
		}

		// Get JWT token with a fresh context
		authCtx, authCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer authCancel()

		// For OAuth users, we pass empty password to indicate OAuth authentication
		authResp, err := s.authService.Login(authCtx, user.Email, "")
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate token for returning OAuth user: %w", err)
		}
		jwtToken = authResp.Token
	} else {
		// New user, try to find if the user exists by email
		findCtx, findCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer findCancel()

		// Try to find the user by email directly
		existingUser, err := s.userService.FindUserByEmail(findCtx, oauthResp.Email)

		// Check if user exists
		var userExists bool
		if err == nil && existingUser != nil {
			userExists = true
		}

		if userExists && existingUser != nil {
			// Email exists, link account
			entityUser = &userEntity.User{
				ID:          existingUser.ID,
				Name:        existingUser.Name,
				Username:    existingUser.Username,
				Email:       existingUser.Email,
				PhoneNumber: existingUser.PhoneNumber,
				Verified:    existingUser.Verified,
				CreatedAt:   existingUser.CreatedAt,
				UpdatedAt:   existingUser.UpdatedAt,
				IsActive:    existingUser.IsActive,
			}

			// Create OAuth connection for existing user
			userDataJSON, err := json.Marshal(map[string]interface{}{
				"name":      oauthResp.Name,
				"email":     oauthResp.Email,
				"avatarURL": oauthResp.AvatarURL,
			})
			if err != nil {
				return nil, "", err
			}

			userOAuth = &model.UserOAuth{
				UserID:         existingUser.ID,
				ProviderID:     dbProvider.ID,
				ProviderUserID: oauthResp.ProviderID,
				AccessToken:    oauthResp.AccessToken,
				RefreshToken:   oauthResp.RefreshToken,
				TokenExpiresAt: oauthResp.ExpiresAt,
				UserData:       userDataJSON,
			}

			createCtx, createCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer createCancel()

			err = s.userOAuthRepo.Create(createCtx, userOAuth)
			if err != nil {
				return nil, "", err
			}

			// Get JWT token with a fresh context
			jwtCtx, jwtCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer jwtCancel()

			// For OAuth users, we pass empty password to indicate OAuth authentication
			authResp, err := s.authService.Login(jwtCtx, existingUser.Email, "")
			if err != nil {
				return nil, "", fmt.Errorf("failed to generate token for existing OAuth user: %w", err)
			}
			jwtToken = authResp.Token
		} else {
			// Create new user
			// Generate a random username based on name and random string
			randomBytes := make([]byte, 4)
			_, err = rand.Read(randomBytes)
			if err != nil {
				return nil, "", err
			}
			randomSuffix := base64.URLEncoding.EncodeToString(randomBytes)[:6]

			// Create a username from the name (convert to lowercase, replace spaces with underscores)
			baseName := oauthResp.Name
			if baseName == "" {
				// If name is empty, use part of the email
				parts := strings.Split(oauthResp.Email, "@")
				baseName = parts[0]
			}

			username := strings.ToLower(strings.ReplaceAll(baseName, " ", "_")) + "_" + randomSuffix

			// Generate a random secure password for OAuth users
			passwordBytes := make([]byte, 32)
			_, err = rand.Read(passwordBytes)
			if err != nil {
				return nil, "", err
			}
			randomPassword := base64.StdEncoding.EncodeToString(passwordBytes)

			// Register new user with a fresh context
			registerCtx, registerCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer registerCancel()

			newUser, err := s.userService.Register(
				registerCtx,
				oauthResp.Name,
				username,
				oauthResp.Email,
				"", // No phone number from OAuth
				randomPassword,
			)
			if err != nil {
				return nil, "", err
			}

			// Convert model.User to entity.User
			entityUser = &userEntity.User{
				ID:          newUser.ID,
				Name:        newUser.Name,
				Username:    newUser.Username,
				Email:       newUser.Email,
				PhoneNumber: newUser.PhoneNumber,
				Verified:    true, // Mark as verified
				CreatedAt:   newUser.CreatedAt,
				UpdatedAt:   newUser.UpdatedAt,
				IsActive:    newUser.IsActive,
			}

			// Create OAuth connection for new user
			userDataJSON, err := json.Marshal(map[string]interface{}{
				"name":      oauthResp.Name,
				"email":     oauthResp.Email,
				"avatarURL": oauthResp.AvatarURL,
			})
			if err != nil {
				return nil, "", err
			}

			userOAuth = &model.UserOAuth{
				UserID:         newUser.ID,
				ProviderID:     dbProvider.ID,
				ProviderUserID: oauthResp.ProviderID,
				AccessToken:    oauthResp.AccessToken,
				RefreshToken:   oauthResp.RefreshToken,
				TokenExpiresAt: oauthResp.ExpiresAt,
				UserData:       userDataJSON,
			}

			// Create OAuth connection with a fresh context
			createCtx, createCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer createCancel()

			err = s.userOAuthRepo.Create(createCtx, userOAuth)
			if err != nil {
				return nil, "", err
			}

			// Get JWT token with a fresh context
			jwtCtx, jwtCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer jwtCancel()

			// For OAuth users, we pass empty password to indicate OAuth authentication
			authResp, err := s.authService.Login(jwtCtx, newUser.Email, "")
			if err != nil {
				return nil, "", fmt.Errorf("failed to generate token for OAuth user: %w", err)
			}
			jwtToken = authResp.Token
		}
	}

	return entityUser, jwtToken, nil
}

// GetUserConnections gets all OAuth connections for a user
func (s *OAuthService) GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*dto.UserOAuthInfo, error) {
	// Get all providers
	providers, err := s.providerRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var connections []*dto.UserOAuthInfo

	// Check each provider
	for _, provider := range providers {
		userOAuth, err := s.userOAuthRepo.GetByUserIDAndProvider(ctx, userID, provider.ID)
		if err != nil {
			return nil, err
		}

		connection := &dto.UserOAuthInfo{
			UserID:    userID,
			Provider:  provider.Name,
			Connected: false,
		}

		if userOAuth != nil {
			connection.Connected = true
			connection.ProviderUserID = userOAuth.ProviderUserID
			connection.ConnectedAt = userOAuth.CreatedAt
		}

		connections = append(connections, connection)
	}

	return connections, nil
}

// DisconnectProvider disconnects a provider from a user
func (s *OAuthService) DisconnectProvider(ctx context.Context, userID uuid.UUID, providerName string) error {
	// Get provider
	dbProvider, err := s.providerRepo.GetByName(ctx, providerName)
	if err != nil {
		return err
	}
	if dbProvider == nil {
		return fmt.Errorf("provider %s not found", providerName)
	}

	// Get user OAuth connection
	userOAuth, err := s.userOAuthRepo.GetByUserIDAndProvider(ctx, userID, dbProvider.ID)
	if err != nil {
		return err
	}
	if userOAuth == nil {
		return fmt.Errorf("user is not connected to %s", providerName)
	}

	// Delete connection
	return s.userOAuthRepo.DeleteByID(ctx, userOAuth.ID)
}

// GetAvailableProviders gets all available OAuth providers
func (s *OAuthService) GetAvailableProviders(ctx context.Context) ([]*dto.OAuthProviderInfo, error) {
	// Get all providers
	providers, err := s.providerRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []*dto.OAuthProviderInfo

	for _, provider := range providers {
		result = append(result, &dto.OAuthProviderInfo{
			ID:          provider.ID,
			Name:        provider.Name,
			DisplayName: provider.DisplayName,
		})
	}

	return result, nil
}
