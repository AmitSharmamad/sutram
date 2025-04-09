package user

import (
	"algo/src/core/db"
	"algo/src/shared"
	"fmt"

	"github.com/lib/pq"
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
		db:   factory.Get("postgres"),
	}
}

func (repo *Repo) Create(user *shared.User) (bool, error) {
	res, err := repo.db.Insert(
		`INSERT INTO "user" (name, username, email, password, phone_number) VALUES (
		:name, :username, :email, :password, :phone_number
		)`,
		user,
	)
	if err != nil {
		// Check if this is a unique violation error
		// NOTE: Check if this is a how it's supposed to be
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			errorMessage := "Username or email already exists"
			return false, shared.NewBadRequest(errorMessage, "DUPLICATE_RECORD")
		}
		return false, err
	}
	rows, e := res.RowsAffected()
	if rows != 1 {
		return false, e
	}
	return true, nil
}

func (repo *Repo) GetAll(offset int, limit int) (shared.Users, error) {
	var users shared.Users

	query := fmt.Sprintf(`
        SELECT *
        FROM "user"
		ORDER BY created_at DESC
		LIMIT %d OFFSET %d
    `, limit, offset)

	err := repo.db.Select(
		&users,
		query,
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
