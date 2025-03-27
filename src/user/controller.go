package user

import (
	"algo/src/shared"
	"algo/src/util"

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

func (controller *Controller) GetAll(gctx *gin.Context) {
	limit := gctx.DefaultQuery("limit", "10")
	offset := gctx.DefaultQuery("offset", "0")

	// Convert to int
	limitInt := util.ParseIntWithDefault(limit, 10)
	offsetInt := util.ParseIntWithDefault(offset, 0)

	users, err := controller.service.GetAll(offsetInt, limitInt)
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

	users, err := controller.service.GetAll(0, 1)
	if err != nil {
		panic(err)
	}

	shared.SuccessResp(gctx, users)
}
