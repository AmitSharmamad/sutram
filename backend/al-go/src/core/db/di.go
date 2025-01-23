package db

import "github.com/google/wire"

type Root struct{}

func NewRoot(*DBFactory) *Root {
	return &Root{}
}

var Set = wire.NewSet(
	NewRoot,
	NewFactoryProv,
	NewClickHouseProv,
	NewPostgresProv,
)
