package community

import (
	"github.com/google/wire"
)

// ProvideSet provides the community dependencies.
var ProvideSet = wire.NewSet(
	NewRepository,
	NewService,
	NewController,
	NewRouter,
)
