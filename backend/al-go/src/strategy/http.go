package strategy

// REST API

type Handler struct {
	Service *Service
}

func NewHandlerProv(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}
