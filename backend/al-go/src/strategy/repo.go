package strategy

import (
	"context"

	"algo/src/core/db"
	"algo/src/shared"

	"github.com/google/uuid"
)

type Repo struct {
	Name string
	db   db.Database
}

func NewRepoProv(
	factory *db.DBFactory,
) *Repo {
	return &Repo{
		Name: "sutram.strategy",
		db:   factory.Get("clickhouse"),
	}
}

func (repo *Repo) Create(ctx context.Context, strategy *shared.Strategy) error {
	userID := shared.FromReqCtx[uuid.UUID](ctx, shared.UserID)

	strategy.UserID = userID

	_, err := repo.db.Insert(
		`
		INSERT INTO sutram.strategy (
			id, user_id, name, description, entry_expression, 
      exit_expression
    ) VALUES (
    	:id, :user_id, :name, :description, :entry_expression,
    	:exit_expression
    );
		`,
		strategy,
	)
	return err
}

func (repo *Repo) Get(
	ctx context.Context,
	id uuid.UUID,
) (*shared.Strategy, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	var strategies shared.Strategies

	err := repo.db.Select(
		&strategies,
		`
		SELECT *
		FROM $1
		WHERE user_id = $2 
		AND id = $3;
		`,
		repo.Name,
		userID,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &strategies[0], nil
}

func (repo *Repo) GetLatest(
	ctx context.Context,
) (*shared.Strategy, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	var strategies shared.Strategies

	err := repo.db.Select(
		&strategies,
		`
			SELECT *
			FROM $1
			WHERE user_id = $2
			ORDER BY timestamp DESC
			LIMIT 1;
		`,
		repo.Name,
		userID,
	)

	return &strategies[0], err
}

func (repo *Repo) GetAll(
	ctx context.Context,
) (*shared.Strategies, error) {
	userID := shared.FromReqCtx[string](ctx, shared.UserID)

	var strategies shared.Strategies

	repo.db.Select(
		&strategies,
		"SELECT * FROM $1 WHERE user_id = $2;",
		repo.Name,
		userID,
	)

	return &strategies, nil
}
