//go:build wireinject
// +build wireinject

package main

import (
	"algo/internal/modules/backtest"
	"algo/internal/modules/deployment"
	"algo/internal/modules/instrument"
	"algo/internal/modules/oauth"
	"algo/internal/modules/strategy"
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

		// Strategy module
		strategy.New,

		// Deployment module
		deployment.New,

		// Backtest module
		backtest.New,

		// OAuth module
		oauth.New,
	)
	return &App{}, nil
}

// App represents the application container
type App struct {
	UserModule       *user.Module
	InstrumentModule *instrument.Module
	StrategyModule   *strategy.Module
	DeploymentModule *deployment.Module
	BacktestModule   *backtest.Module
	OauthModule      *oauth.Module
}

// NewApp creates a new application instance
func NewApp(userModule *user.Module, instrumentModule *instrument.Module, strategyModule *strategy.Module, deploymentModule *deployment.Module, backtestModule *backtest.Module, oauthModule *oauth.Module) *App {
	return &App{
		UserModule:       userModule,
		InstrumentModule: instrumentModule,
		StrategyModule:   strategyModule,
		DeploymentModule: deploymentModule,
		BacktestModule:   backtestModule,
		OauthModule:      oauthModule,
	}
}
