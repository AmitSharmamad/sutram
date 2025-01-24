package strategy

import "github.com/google/wire"

type Root struct {
	Controller *Controller
}

func NewRootProv(controller *Controller) *Root {
	return &Root{
		Controller: controller,
	}
}

var Set = wire.NewSet(
	NewRootProv,
	NewRepoProv,
	NewServiceProv,
	NewControllerProv,
	RegisterRoutesProv,
)
