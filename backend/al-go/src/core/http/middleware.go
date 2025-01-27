package http

import (
	"log"
	"strings"
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
	return func(ctx *gin.Context) {
		t := time.Now()

		println("Validating the token")

		authorization := ctx.GetHeader("Authorization")

		if len(strings.TrimSpace(authorization)) == 0 {
			println("No auth token")
			panic(shared.NewUnauthorizedError("Unauthorized", "Unauthorized"))
		}

		values := strings.Split(authorization, " ")

		if len(values) != 2 || values[0] != "Bearer" {
			println("Incorrect auth token format")
			panic(shared.NewUnauthorizedError("Unauthorized", "Unauthorized"))
		}

		token := values[1]

		user, err := middleware.authService.ValidateToken(token)
		if err != nil {
			println("Error in decoding the token")
			panic(shared.NewUnauthorizedError("Unauthorized", "Unauthorized"))
		}

		ctx.Set("user_id", user.ID)

		ctx.Next()

		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := ctx.Writer.Status()
		log.Println(status)
	}
}
