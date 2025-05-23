package deployment

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	backtestEngine "algo/internal/modules/backtest/application/service"
	"algo/internal/modules/deployment/application/service"
	"algo/internal/modules/deployment/domain/repository"
	"algo/internal/modules/deployment/infrastructure/persistence/postgres"
	"algo/internal/modules/deployment/interfaces/http"
	instrumentRepository "algo/internal/modules/instrument/domain/repository"
	strategyRepository "algo/internal/modules/strategy/domain/repository"
)

// ProvideDeploymentRepository provides a deployment repository
func ProvideDeploymentRepository(db *sql.DB) repository.DeploymentRepository {
	return postgres.NewDeploymentRepository(db)
}

// ProvideDeploymentLogRepository provides a deployment log repository
func ProvideDeploymentLogRepository(db *sqlx.DB) repository.DeploymentLogRepository {
	return postgres.NewPgDeploymentLogRepository(db)
}

// ProvideDeploymentPositionRepository provides a deployment position repository
func ProvideDeploymentPositionRepository(db *sqlx.DB) repository.DeploymentPositionRepository {
	return postgres.NewPgDeploymentPositionRepository(db)
}

// ProvideDeploymentMetricsRepository provides a deployment metrics repository
func ProvideDeploymentMetricsRepository(db *sqlx.DB) repository.DeploymentMetricsRepository {
	return postgres.NewPgDeploymentMetricsRepository(db)
}

// ProvideDeploymentPositionLogRepository provides a deployment position log repository
func ProvideDeploymentPositionLogRepository(db *sqlx.DB) repository.DeploymentPositionLogRepository {
	return postgres.NewPgDeploymentPositionLogRepository(db)
}

// ProvideDeploymentService provides a deployment service
func ProvideDeploymentService(
	deploymentRepo repository.DeploymentRepository,
	strategyRepo strategyRepository.StrategyRepository,
	logRepo repository.DeploymentLogRepository,
	positionRepo repository.DeploymentPositionRepository,
	metricsRepo repository.DeploymentMetricsRepository,
	instrumentRepo instrumentRepository.InstrumentRepository,
	engine *backtestEngine.BacktestEngine,
) service.DeploymentService {
	return service.NewDeploymentService(
		deploymentRepo,
		strategyRepo,
		logRepo,
		positionRepo,
		metricsRepo,
		instrumentRepo,
		engine,
	)
}

// ProvideDeploymentPositionLogService provides a deployment position log service
func ProvideDeploymentPositionLogService(
	repo repository.DeploymentPositionLogRepository,
) *service.DeploymentPositionLogService {
	return service.NewDeploymentPositionLogService(repo)
}

// ProvideDeploymentController provides a deployment controller
func ProvideDeploymentController(svc service.DeploymentService) *http.DeploymentController {
	return http.NewDeploymentController(svc)
}

// ProvideDeploymentLogController provides a deployment log controller
func ProvideDeploymentLogController(svc service.DeploymentService) *http.DeploymentLogController {
	return http.NewDeploymentLogController(svc)
}

// ProvideDeploymentPositionController provides a deployment position controller
func ProvideDeploymentPositionController(svc service.DeploymentService) *http.DeploymentPositionController {
	return http.NewDeploymentPositionController(svc)
}

// ProvideDeploymentMetricsController provides a deployment metrics controller
func ProvideDeploymentMetricsController(svc service.DeploymentService) *http.DeploymentMetricsController {
	return http.NewDeploymentMetricsController(svc)
}

// ProvideDeploymentPositionLogController provides a deployment position log controller
func ProvideDeploymentPositionLogController(svc *service.DeploymentPositionLogService) *http.DeploymentPositionLogController {
	return http.NewDeploymentPositionLogController(svc)
}
