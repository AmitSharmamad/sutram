package shared

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	GetToken(user *User) string
}

func TokenValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		println("Validating the token")

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
