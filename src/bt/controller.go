package bt

import (
	"algo/src/shared"
	"algo/src/util"
	"fmt"

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

func (controller *Controller) DoBacktest(gctx *gin.Context) {

	ctx := shared.GetReqCtx(gctx)

	// convert the strategyID to a uuid
	strategyID := util.Parse(gctx.Query("strategy_id"))

	runGroup := controller.engine.Start(ctx, strategyID)

	shared.SuccessResp(gctx, runGroup)
}

func (controller *Controller) Start(gctx *gin.Context) {
	// ctx := shared.GetReqCtx(gctx)
	// id := gctx.Param("id")
	// run, err := controller.service.GetRun(id)
	// if err != nil {
	// 	panic(err)
	// }
	// println(run)
	// indicators := controller.engine.Run(ctx, run)
	shared.SuccessResp(gctx, []string{})
}

func (controller *Controller) Pause(gctx *gin.Context) {}

func (controller *Controller) Resume(gctx *gin.Context) {}

func (controller *Controller) GetMetrics(gctx *gin.Context) {
	id := gctx.Param("id")

	run, err := controller.service.GetRun(id)
	if err != nil {
		panic(err)
	}
	metrics := controller.engine.CalculateRunMetrics(run)
	response := map[string]interface{}{}
	response["id"] = id
	response["metrics"] = metrics
	response["status"] = "work in progress"
	shared.SuccessResp(gctx, response)
}

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

func (controller *Controller) DeployLive(gctx *gin.Context) {
	ctx := shared.GetReqCtx(gctx)

	strategyID := util.Parse(gctx.Param("id"))
	strategy := controller.service.strategyService.Get(ctx, strategyID)

	if !strategy.IsActive {
		panic(shared.NewBadRequest("Strategy is not active", "INACTIVE_STRATEGY"))
	}

	runGroup := shared.NewRunGroup(strategyID)
	runs := make([]*shared.Run, 0)

	// Create runs for each instrument
	for _, instrumentID := range strategy.Instruments {
		run := shared.NewRun(
			runGroup.ID,
			strategyID,
			instrumentID,
		)
		run.RunType = shared.LiveTrade
		runs = append(runs, run)
	}

	runGroup.Runs = runs

	// Save runs to database first
	controller.service.runGroupRepo.Create(runGroup)
	controller.service.runRepo.Create(runs)

	// Start live trading for each run in separate goroutines
	errChan := make(chan error, len(runs))
	for _, run := range runs {
		go func(r *shared.Run) {
			errChan <- controller.engine.RunLiveTrading(ctx, r, strategy)
		}(run)
	}

	// Monitor for errors
	go func() {
		for err := range errChan {
			if err != nil {
				fmt.Printf("Live trading error: %v\n", err)
			}
		}
	}()

	shared.SuccessResp(gctx, runGroup)
}
