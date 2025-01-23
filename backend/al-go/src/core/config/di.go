package env

import "github.com/google/wire"

type Root struct{}

func NewRoot(*EnvProv) *Root {
	return &Root{}
}

var Set = wire.NewSet(
	NewRoot,
	NewEnvProv,
)
