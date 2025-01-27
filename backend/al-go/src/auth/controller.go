package auth

import (
	"net/http"

	"algo/src/shared"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	authService shared.AuthService
	userService shared.UserService
}

func NewControllerProv(
	authService shared.AuthService,
	userService shared.UserService,
) *Controller {
	return &Controller{
		authService: authService,
		userService: userService,
	}
}

func (controller *Controller) Login(ctx *gin.Context) {
	req := LoginReq{}
	ctx.Bind(&req)

	user, err := controller.userService.GetByUsername(req.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"data":   err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": true,
		"data": &LoginRes{
			Token: controller.authService.GetToken(&user),
		},
	})
}

func (controller *Controller) Logout(ctx *gin.Context) {}
