package bt

import "github.com/google/wire"

type Root struct {
}

func NewRoot(*Controller) *Root {
	return &Root{}
}

var Set = wire.NewSet(
	NewRoot,
	NewRepo,
	NewService,
	NewController,
)
