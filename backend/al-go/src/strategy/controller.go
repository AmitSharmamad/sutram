package strategy

import (
	"algo/src/data"

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

func (*Controller) GetStrategies(ctx *gin.Context) {
	ctx.SecureJSON(200, true)
}

func (this *Controller) GetStrategy(ctx *gin.Context) {
	var data []data.Data
	_, err := this.service.repo.db.SelectOne("SELECT * FROM sutram.stock_TESLA LIMIT 10;", &data)
	if err != nil {
		println("from inside the error ", err.Error())
		ctx.JSON(200, gin.H{
			"message": "error",
		})
		return
	}

	// r, _ := json.Marshal(data)

	ctx.JSON(200, gin.H{
		"message": data,
	})
}
func (*Controller) GetStrategyRuns(ctx *gin.Context) {}
func (*Controller) GetStrategyRun(ctx *gin.Context)  {}
func (*Controller) Run(ctx *gin.Context)             {}
