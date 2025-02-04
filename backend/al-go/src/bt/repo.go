package bt

import (
	"algo/src/core/db"
)

type RunPositionRepo struct {
	RunName      string
	PositionName string
	db           db.Database
}

func NewRunPositionRepoProv(
	factory *db.DBFactory,
) *RunPositionRepo {
	return &RunPositionRepo{
		RunName:      `sutram.bt_run`,
		PositionName: `sutram.bt_position`,
	}
}

// func (repo *RunPositionRepo) Get(id string) (Runs, error) {
// 	var runs RunPosition
//
// 	err := repo.db.Select(&runs, `
// 		SELECT *
// 		FROM $1
// 		JOIN $2
// 		ON $1.id = $2.run_id
// 		ORDER BY timestamp DESC;
// 		`,
// 		repo.RunName,
// 		repo.PositionName,
// 	)
// 	return runs, err
// }
