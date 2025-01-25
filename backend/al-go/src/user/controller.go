package user

import (
	"net/http"

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

func (this *Controller) Get(ctx *gin.Context) {
	users, err := this.service.Get()
	if err != nil {
		panic(err)
	}

	ctx.JSON(200, gin.H{
		"result": true,
		"data":   users,
	})
}

func (this *Controller) Create(ctx *gin.Context) {
	user := User{}
	ctx.Bind(&user)

	_, err := this.service.Create(&user)
	if err != nil {
		panic(err)
	}

	users, err := this.service.Get()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   users,
	})
}

func (this *Controller) Login(ctx *gin.Context) {
	req := LoginReq{}
	ctx.Bind(&req)

	user, err := this.service.GetByUsername(req.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"data":   err,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   user,
	})
}

func (this *Controller) Logout(ctx *gin.Context) {}
