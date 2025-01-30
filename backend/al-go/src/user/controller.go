package user

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

func (controller *Controller) Get(gctx *gin.Context) {
	users, err := controller.service.Get()
	if err != nil {
		panic(err)
	}

	shared.SuccessResp(gctx, users)
}

func (controller *Controller) Create(gctx *gin.Context) {
	user := User{}
	gctx.Bind(&user)

	_, err := controller.service.Create(&user)
	if err != nil {
		panic(err)
	}

	users, err := controller.service.Get()
	if err != nil {
		panic(err)
	}

	shared.SuccessResp(gctx, users)
}
