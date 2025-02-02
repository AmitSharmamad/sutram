package http

import (
	"log"
	"strings"
	"time"

	"algo/src/core/logger"
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
	return func(gctx *gin.Context) {
		t := time.Now()

		logger.Log("validating the token")

		authorization := gctx.GetHeader("Authorization")

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

		gctx.Set("user_id", user.ID)

		gctx.Next()

		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := gctx.Writer.Status()
		log.Println(status)
	}
}
