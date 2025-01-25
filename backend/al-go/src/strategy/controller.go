package strategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func (this *Controller) GetStrategies(ctx *gin.Context) {
	strategies, err := this.service.GetAll()
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   strategies,
	})
}

func (this *Controller) GetStrategy(ctx *gin.Context) {
}

func (*Controller) GetStrategyRuns(ctx *gin.Context) {}

func (*Controller) GetStrategyRun(ctx *gin.Context) {}

func (*Controller) Run(ctx *gin.Context) {}
