package strategy

import (
	"database/sql"

	"algo/src/core/db"
)

type Strategies []Strategy

type Repo struct {
	Name string
	db   db.Database
}

func NewRepoProv(
	factory *db.DBFactory,
) *Repo {
	return &Repo{
		Name: "strategy",
		db:   factory.Get("clickhouse"),
	}
}

func (repo *Repo) Create(strategy *Strategy) (sql.Result, error) {
	return repo.db.Insert(
		"INSERT INTO ? VALUES (?,?,?,?)",
		repo.Name,
		strategy.Name,
	)
}

func (repo *Repo) Get(id string) (*Strategy, error) {
	// results, err := repo.db.Select(
	// 	`
	// 	SELECT *
	// 	FROM ?
	// 	WHERE id = '?';
	// 	`,
	// 	repo.Name,
	// 	id,
	// )
	// if err != nil {
	// 	return nil, err
	// }
	// return results[0].(*Strategy), nil
	return nil, nil
}

func (repo *Repo) GetLatest() (any, error) {
	return nil, nil
	// return repo.db.Select(
	// 	`
	// 		SELECT *
	// 		FROM ?
	// 		ORDER BY DESC
	// 		LIMIT 1;
	// 	`,
	// 	repo.Name,
	// )
}

func (this *Repo) GetAll() (Strategies, error) {
	var strategies Strategies
	this.db.Select(&strategies, "SELECT * FROM strategy LIMIT 10;")
	return strategies, nil
}

func (repo *Repo) Update(id string, strategy *Strategy) (sql.Result, error) {
	return repo.db.Update("UPDATE strategy SET ... WHERE id = {id}")
}

func (repo *Repo) Delete(query string) (sql.Result, error) {
	return repo.db.Delete("DELETE FROM strategy")
}
