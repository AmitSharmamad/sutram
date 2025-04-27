//go:build wireinject
// +build wireinject

package main

import (
	"algo/internal/modules/user"
	
	"github.com/google/wire"
)

// InitializeApp sets up the application with all its dependencies
func InitializeApp() (*App, error) {
	wire.Build(
		// Core components
		NewApp,

		// User module
		user.New,
	)
	return &App{}, nil
}

// App represents the application container
type App struct {
	UserModule *user.Module
}

// NewApp creates a new application instance
func NewApp(userModule *user.Module) *App {
	return &App{
		UserModule: userModule,
	}
}
