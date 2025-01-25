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

func (this *Repo) Insert()    {}
func (this *Repo) SelectOne() {}
func (this *Repo) Select()    {}
func (this *Repo) Update()    {}
func (this *Repo) Delete()    {}
