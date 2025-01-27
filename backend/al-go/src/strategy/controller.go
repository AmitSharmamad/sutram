package strategy

import (
	"net/http"

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

func (controller *Controller) GetAll(ctx *gin.Context) {
	strategies, err := controller.service.GetAll()
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   strategies,
	})
}

func (controller *Controller) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	println(id)

	strategy := controller.service.Get(id)

	ctx.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   strategy,
	})
}

func (*Controller) GetRuns(ctx *gin.Context) {}

func (*Controller) GetRun(ctx *gin.Context) {}

func (controller *Controller) Create(ctx *gin.Context) {
	strategy := &Strategy{}
	ctx.Bind(strategy)

	isValid := controller.service.Validate(strategy)

	if !isValid {
		panic(shared.NewBadRequest("Invalid Expression(s)", "INVALID_EXPRESSION"))
	}

	user, err := controller.userService.GetByUsername("amit")
	if err != nil {
		panic(shared.NewBadRequest("Invalid User", "INVALID_USER"))
	}

	userID, _ := uuid.Parse(user.ID)

	strategy.UserID = userID

	controller.service.Create(strategy)

	s := controller.service.Get("")

	ctx.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   s,
	})
}

func (*Controller) Run(ctx *gin.Context) {}
