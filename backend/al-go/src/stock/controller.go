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

func (this *Controller) GetAll(ctx *gin.Context) {
	println("coming here")
}

func (this *Controller) Get(ctx *gin.Context) {
	data, err := this.service.Get("")

	if err != nil {
		panic(err)
	}

	ctx.JSON(200, gin.H{
		"data": data,
	})

}
