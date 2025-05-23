package deployment

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"algo/internal/modules/deployment/application/service"
	"algo/internal/modules/deployment/domain/repository"
	"algo/internal/modules/deployment/infrastructure/persistence/postgres"
	"algo/internal/modules/deployment/interfaces/http"
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

// ProvideDeploymentService provides a deployment service
func ProvideDeploymentService(
	deploymentRepo repository.DeploymentRepository,
	strategyRepo strategyRepository.StrategyRepository,
	logRepo repository.DeploymentLogRepository,
	positionRepo repository.DeploymentPositionRepository,
	metricsRepo repository.DeploymentMetricsRepository,
) service.DeploymentService {
	return service.NewDeploymentService(deploymentRepo, strategyRepo, logRepo, positionRepo, metricsRepo)
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
