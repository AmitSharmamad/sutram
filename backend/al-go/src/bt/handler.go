package bt

import "github.com/gin-gonic/gin"

type IController interface {
	Run(c *gin.Context)
}

type Controller struct {
	service *Service
}

func NewController(
	service *Service,
) *Controller {
	return &Controller{
		service: service,
	}
}

func (*Controller) Run(c *gin.Context) {

}
