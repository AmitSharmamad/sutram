package auth

import "algo/src/core/db"

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

func (repo *Repo) Insert()    {}
func (repo *Repo) SelectOne() {}
func (repo *Repo) Select()    {}
func (repo *Repo) Update()    {}
func (repo *Repo) Delete()    {}
