package deployment

import (
	"database/sql"

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

// ProvideDeploymentService provides a deployment service
func ProvideDeploymentService(
	deploymentRepo repository.DeploymentRepository,
	strategyRepo strategyRepository.StrategyRepository,
) service.DeploymentService {
	return service.NewDeploymentService(deploymentRepo, strategyRepo)
}

// ProvideDeploymentController provides a deployment controller
func ProvideDeploymentController(svc service.DeploymentService) *http.DeploymentController {
	return http.NewDeploymentController(svc)
}
