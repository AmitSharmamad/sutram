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

func (*Controller) GetStrategies(ctx *gin.Context) {
	ctx.SecureJSON(200, true)
}
func (this *Controller) GetStrategy(ctx *gin.Context) {
	results, err := this.service.repo.db.Select("SELECT COUNT(*) FROM sutram.stock_TESLA;")
	if err != nil {
		println("from inside the error ", err.Error())
		ctx.JSON(200, gin.H{
			"message": "error",
		})
		return
	}
	println(err)
	ctx.JSON(200, gin.H{
		"message": results,
	})
}
func (*Controller) GetStrategyRuns(ctx *gin.Context) {}
func (*Controller) GetStrategyRun(ctx *gin.Context)  {}
func (*Controller) Run(ctx *gin.Context)             {}
