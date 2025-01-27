package shared

type AuthService interface {
	GetToken(user *User) string
	ValidateToken(token string) (*User, error)
}
