package strategy

import (
	"fmt"

	"algo/src/shared"

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

	gctx.Bind(strategy)

	isValid := controller.service.Validate(strategy)

	if !isValid {
		panic(shared.NewBadRequest("Invalid Expression(s)", "INVALID_EXPRESSION"))
	}

	fmt.Printf("%+v\n", strategy)

	strategy, err := controller.service.Create(ctx, strategy)
	if err != nil {
		panic(err)
	}

	s := controller.service.Get(ctx, strategy.ID)

	shared.SuccessResp(gctx, s)
}
