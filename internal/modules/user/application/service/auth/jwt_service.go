package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService provides JWT authentication functionality
type JWTService struct {
	secretKey     string
	tokenDuration time.Duration
}

// TokenClaims represents the JWT token claims
type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT authentication service
func NewJWTService(secretKey string, tokenDuration time.Duration) *JWTService {
	return &JWTService{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// GenerateToken generates a JWT token for a user
func (s *JWTService) GenerateToken(userID uuid.UUID, email string, roles []string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "sutram",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GetUserIDFromToken extracts the user ID from a token
func (s *JWTService) GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

// GetRolesFromToken extracts the roles from a token
func (s *JWTService) GetRolesFromToken(tokenString string) ([]string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	return claims.Roles, nil
}

// HasRole checks if the token contains a specific role
func (s *JWTService) HasRole(tokenString string, role string) (bool, error) {
	roles, err := s.GetRolesFromToken(tokenString)
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		if r == role {
			return true, nil
		}
	}

	return false, nil
}
