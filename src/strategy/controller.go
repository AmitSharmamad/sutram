package strategy

import (
	"algo/src/shared"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	service     *Service
	userService shared.UserService
}

func NewControllerProv(
	service *Service,
	userService shared.UserService,
) *Controller {
	return &Controller{
		service:     service,
		userService: userService,
	}
}

func (controller *Controller) GetAll(gctx *gin.Context) {
	ctx := shared.GetReqCtx(gctx)

	strategies, err := controller.service.GetAll(ctx)
	if err != nil {
		panic(err)
	}

	shared.SuccessResp(gctx, strategies)
}

func (controller *Controller) Get(gctx *gin.Context) {
	ctx := shared.GetReqCtx(gctx)

	id, _ := uuid.Parse(gctx.Param("id"))

	strategy := controller.service.Get(ctx, id)

	shared.SuccessResp(gctx, strategy)
}

func (controller *Controller) Create(gctx *gin.Context) {
	ctx := shared.GetReqCtx(gctx)

	strategy := new(shared.Strategy)

	if err := gctx.ShouldBind(strategy); err != nil {
		fmt.Printf("Binding error: %v\n", err)
		panic(shared.NewBadRequest("Invalid request body: "+err.Error(), "INVALID_REQUEST"))
	}

	// // Validate dates
	// now := time.Now()
	// if strategy.StartPeriod.After(strategy.EndPeriod.Time) {
	// 	panic(shared.NewBadRequest("Start date must be before end date", "INVALID_DATE_RANGE"))
	// }
	// if strategy.EndPeriod.Before(now) {
	// 	panic(shared.NewBadRequest("End date cannot be in the past", "INVALID_DATE"))
	// }

	// TODO: Work on the validation logic
	// isValid := controller.service.Validate(strategy)
	// if !isValid {
	// 	panic(shared.NewBadRequest("Invalid Expression(s)", "INVALID_EXPRESSION"))
	// }

	fmt.Printf("Received strategy: %+v\n", strategy)

	strategy, err := controller.service.Create(ctx, strategy)
	if err != nil {
		panic(err)
	}

	s := controller.service.Get(ctx, strategy.ID)

	shared.SuccessResp(gctx, s)
}

func (controller *Controller) Update(gctx *gin.Context) {
	ctx := shared.GetReqCtx(gctx)

	id, err := uuid.Parse(gctx.Param("id"))
	if err != nil {
		panic(shared.NewBadRequest("Invalid strategy ID", "INVALID_ID"))
	}

	strategy := new(shared.Strategy)
	if err := gctx.ShouldBind(strategy); err != nil {
		fmt.Printf("Binding error: %v\n", err)
		panic(shared.NewBadRequest("Invalid request body: "+err.Error(), "INVALID_REQUEST"))
	}

	strategy, err = controller.service.Update(ctx, id, strategy)
	if err != nil {
		panic(err)
	}

	shared.SuccessResp(gctx, strategy)
}
