package strategy

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

func (*Controller) GetStrategies(ctx *gin.Context)   {}
func (*Controller) GetStrategy(ctx *gin.Context)     {}
func (*Controller) GetStrategyRuns(ctx *gin.Context) {}
func (*Controller) GetStrategyRun(ctx *gin.Context)  {}
func (*Controller) Run(ctx *gin.Context)             {}
