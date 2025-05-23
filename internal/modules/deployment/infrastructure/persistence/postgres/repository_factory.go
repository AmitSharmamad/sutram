package postgres

import (
	"database/sql"

	"algo/internal/modules/deployment/domain/repository"

	"github.com/jmoiron/sqlx"
)

// Factory for creating deployment repositories
func NewRepositoryFactory(db *sql.DB) *RepositoryFactory {
	return &RepositoryFactory{
		db:     db,
		sqlxDB: sqlx.NewDb(db, "postgres"),
	}
}

type RepositoryFactory struct {
	db     *sql.DB
	sqlxDB *sqlx.DB
}

func (f *RepositoryFactory) CreateDeploymentRepository() repository.DeploymentRepository {
	return NewDeploymentRepository(f.db)
}

func (f *RepositoryFactory) CreateDeploymentPositionLogRepository() repository.DeploymentPositionLogRepository {
	return NewPgDeploymentPositionLogRepository(f.sqlxDB)
}
