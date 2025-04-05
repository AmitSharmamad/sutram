package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterMiddleware(engine *gin.Engine) {
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add middleware to allow multiple reads of request body
	engine.Use(func(c *gin.Context) {
		buf, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
		c.Next()
	})

	engine.Use(GlobalErrorHandler)
}
