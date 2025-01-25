package stock

import "github.com/gin-gonic/gin"

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

func (this *Controller) Get(ctx *gin.Context) {}
