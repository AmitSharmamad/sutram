// Package main provides a simple Swagger docs utility
package main

import (
	_ "github.com/swaggo/swag"
)

// @title Sutram API
// @version 1.0
// @description API server for the Sutram algorithmic trading platform
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@sutram.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9090
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token

// @tag.name Strategy
// @tag.description Strategy management endpoints

// @tag.name User
// @tag.description User management endpoints

// @tag.name Instrument
// @tag.description Instrument management endpoints

// @tag.name Deployment
// @tag.description Strategy deployment endpoints

// @tag.name Backtest
// @tag.description Strategy backtesting endpoints

// @tag.name Marketplace
// @tag.description Marketplace endpoints

// @tag.name Community
// @tag.description Community feature endpoints

// @tag.name OAuth
// @tag.description OAuth authentication endpoints

func main() {
	// This is just a placeholder file for Swagger documentation
	// It's not meant to be executed
}