package user

import (
	"algo/src/core/db"
	"algo/src/shared"
)

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

func (repo *Repo) Create(user *User) (bool, error) {
	res, err := repo.db.Insert(
		`INSERT INTO sutram.user (name, username, email, password) VALUES (
		:name, :username, :email, :password
		)`,
		user,
	)
	if err != nil {
		return false, err
	}
	rows, e := res.RowsAffected()
	if rows != 1 {
		return false, e
	}
	return true, nil
}

func (repo *Repo) GetAll() {}

func (repo *Repo) Get() (shared.Users, error) {
	var users shared.Users
	err := repo.db.Select(
		&users,
		"SELECT * FROM sutram.user LIMIT 10;",
	)
	return users, err
}

func (repo *Repo) GetUsername(username string) (shared.Users, error) {
	var users shared.Users
	err := repo.db.Select(
		&users,
		"SELECT * FROM sutram.user WHERE username = $1;",
		username,
	)
	return users, err
}

func (repo *Repo) Update() {}

func (repo *Repo) Delete() {}
