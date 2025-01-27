package http

import (
	"log"
	"time"

	"algo/src/shared"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService shared.AuthService
}

func NewAuthMiddleware(
	authService shared.AuthService,
) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (middleware *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		println("Validating the token")
		middleware.authService.ValidateToken("")

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}
