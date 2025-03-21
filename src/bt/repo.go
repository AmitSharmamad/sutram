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
