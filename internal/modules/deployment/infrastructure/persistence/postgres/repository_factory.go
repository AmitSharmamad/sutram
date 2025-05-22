package postgres

import (
	"database/sql"

	"algo/internal/modules/deployment/domain/repository"
)

// Factory for creating deployment repositories
func NewRepositoryFactory(db *sql.DB) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

type RepositoryFactory struct {
	db *sql.DB
}

func (f *RepositoryFactory) CreateDeploymentRepository() repository.DeploymentRepository {
	return NewDeploymentRepository(f.db)
}
