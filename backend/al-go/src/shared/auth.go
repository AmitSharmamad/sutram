package shared

type AuthService interface {
	GetToken(user *User) string
}
