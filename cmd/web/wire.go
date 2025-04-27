//go:build wireinject
// +build wireinject

package main

import (
	"algo/internal/modules/instrument"
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

		// Instrument module
		instrument.New,
	)
	return &App{}, nil
}

// App represents the application container
type App struct {
	UserModule       *user.Module
	InstrumentModule *instrument.Module
}

// NewApp creates a new application instance
func NewApp(userModule *user.Module, instrumentModule *instrument.Module) *App {
	return &App{
		UserModule:       userModule,
		InstrumentModule: instrumentModule,
	}
}
