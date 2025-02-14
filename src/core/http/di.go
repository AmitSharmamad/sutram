package http

import (
	"github.com/google/wire"
)

type Root struct{}

func NewRoot(
	*AuthMiddleware,
) *Root {
	return new(Root)
}

var Set = wire.NewSet(
	NewRoot,
	NewAuthMiddleware,
)
