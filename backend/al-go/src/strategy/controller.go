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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid expression(s)",
		})
		return
	}

	user, err := controller.userService.GetByUsername("amit")
	if err != nil {
		// println("coming insided error uuid", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid user",
		})
		return
	}

	userID, _ := uuid.Parse(user.ID)

	strategy.UserID = userID

	new, err := controller.service.Create(strategy)

	if err != nil || !new {
		println("coming insided error create", err, new)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": true,
			"error":  err,
		})
		return
	}

	s := controller.service.Get("")

	println("coming here")

	ctx.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   s,
	})
}

func (*Controller) Run(ctx *gin.Context) {}
