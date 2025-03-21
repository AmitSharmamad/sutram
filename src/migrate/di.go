package migrate

import (
	"github.com/google/wire"
)

type Root struct{}

func NewRoot() *Root {
	return new(Root)
}

var Set = wire.NewSet(
	NewRoot,
)
