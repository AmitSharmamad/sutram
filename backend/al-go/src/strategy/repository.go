package strategy

import (
	"algo/src/core/db"
	"database/sql"
)

type Repo struct {
	Name string
	DB   db.Database
}

func NewRepo() *Repo {
	return &Repo{
		Name: "strategy",
		// DB:   db,
	}
}

func (repo *Repo) Insert(strategy *Strategy) (sql.Result, error) {
	return repo.DB.Insert(
		"INSERT INTO ? VALUES (?,?,?,?)",
		repo.Name,
		strategy.Name,
	)
}

func (repo *Repo) Get(id string) (any, error) {
	return repo.DB.Select(
		`
		SELECT * 
		FROM ? 
		WHERE id = '?';
		`,
		repo.Name,
		id,
	)
}

func (repo *Repo) GetLatest() (any, error) {
	return repo.DB.Select(
		`
			SELECT * 
			FROM ? 
			ORDER BY DESC
			LIMIT 1;
		`,
		repo.Name,
	)
}

func (repo *Repo) GetAll(id string) ([]any, error) {
	return repo.DB.Select("SELECT * FROM strategy;")
}

func (repo *Repo) Update(id string, strategy *Strategy) (sql.Result, error) {
	return repo.DB.Update("UPDATE strategy SET ... WHERE id = {id}")
}

func (repo *Repo) Delete(query string) (sql.Result, error) {
	return repo.DB.Delete("DELETE FROM strategy")
}
