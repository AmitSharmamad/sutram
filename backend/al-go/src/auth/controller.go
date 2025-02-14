package auth

import (
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

func (controller *Controller) Login(gctx *gin.Context) {
	req := LoginReq{}
	gctx.Bind(&req)

	user, err := controller.userService.GetByUsername(req.Username)
	if err != nil {
		shared.NewBadRequest("", "")
		return
	}

	shared.SuccessResp(gctx, &LoginRes{
		Token: controller.authService.GetToken(&user),
	})
}

func (controller *Controller) Logout(ctx *gin.Context) {}
