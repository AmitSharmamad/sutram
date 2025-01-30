package strategy

import (
	"algo/src/shared"

	"github.com/gin-gonic/gin"
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

	id := gctx.Param("id")
	println(id)

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

	controller.service.Create(ctx, strategy)

	s := controller.service.Get(ctx, "")

	shared.SuccessResp(gctx, s)
}
