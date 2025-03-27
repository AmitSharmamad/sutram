package stock

type Controller struct {
	service *Service
}

func NewControllerProv(
	service *Service,
) *Controller {
	return &Controller{
		service: service,
	}
}
