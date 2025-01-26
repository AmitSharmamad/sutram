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

func (controller *Controller) GetAll(ctx *gin.Context) {
	println("coming here")
}

func (controller *Controller) Get(ctx *gin.Context) {
	data, err := controller.service.Get("")
	if err != nil {
		panic(err)
	}

	ctx.JSON(200, gin.H{
		"data": data,
	})
}
