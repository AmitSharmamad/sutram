package bt

import (
	"algo/src/shared"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
	engine  *Engine
}

func NewControllerProv(
	service *Service,
	engine *Engine,
) *Controller {
	return &Controller{
		service: service,
		engine:  engine,
	}
}

func (controller *Controller) CreateRun(gctx *gin.Context) {
	run := shared.NewRun()

	gctx.Bind(run)

	// validate the run input

	run = controller.service.CreateRun(run)

	shared.SuccessResp(gctx, run)
}

func (controller *Controller) Start(gctx *gin.Context) {
	id := gctx.Param("id")
	run, err := controller.service.GetRun(id)
	if err != nil {
		panic("")
	}
	controller.engine.Run(run)
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
