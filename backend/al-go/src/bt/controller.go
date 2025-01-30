package bt

import (
	"algo/src/shared"

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

func (controller *Controller) Run(gctx *gin.Context) {
	strategyID := gctx.Query("strategy")

	run := controller.service.CreateRun(strategyID)

	shared.SuccessResp(gctx, run)
}

func (controller *Controller) Pause(gctx *gin.Context) {}

func (controller *Controller) Resume(gctx *gin.Context) {}

func (controller *Controller) GetAllRuns(gctx *gin.Context) {
	strategyID := gctx.Query("strategy")

	runs := controller.service.GetRuns(strategyID)

	shared.SuccessResp(gctx, runs)
}

func (controller *Controller) GetRun(gctx *gin.Context) {
	id := gctx.Param("id")

	run, positions := controller.service.GetRun(id)

	shared.SuccessResp(gctx, &shared.RunPosition{
		Run:       run,
		Positions: positions,
	})
}

func (controller *Controller) GetRunLogs(gctx *gin.Context) {
	id := gctx.Param("id")

	logs := controller.service.GetRunLogs(id)

	shared.SuccessResp(gctx, logs)
}
