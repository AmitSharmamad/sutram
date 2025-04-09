package auth

import (
	"algo/src/shared"
	"crypto/rand"
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	authService shared.AuthService
	userService shared.UserService
	google      *GoogleService
}

func NewControllerProv(
	authService shared.AuthService,
	userService shared.UserService,
	google *GoogleService,
) *Controller {
	return &Controller{
		authService: authService,
		userService: userService,
		google:      google,
	}
}

func (controller *Controller) Login(gctx *gin.Context) {
	req := LoginReq{}
	gctx.Bind(&req)

	user, err := controller.userService.GetByUsername(req.Username)
	if err != nil {
		shared.NewBadRequest("Invalid credentials", "INVALID_CREDENTIALS")
		return
	}

	shared.SuccessResp(gctx, &LoginRes{
		Token: controller.authService.GetToken(&user),
	})
}

func (controller *Controller) GoogleLogin(gctx *gin.Context) {
	// Generate random state
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		shared.BadReqResp(gctx, "Failed to generate state", "STATE_GENERATION_FAILED")
		return
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in session/cookie for validation
	gctx.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	// Redirect to Google's consent page
	url := controller.google.GetAuthURL(state)
	gctx.Redirect(302, url)
}

func (controller *Controller) GoogleCallback(gctx *gin.Context) {
	// Verify state
	state, _ := gctx.Cookie("oauth_state")
	if state != gctx.Query("state") {
		shared.BadReqResp(gctx, "Invalid state", "INVALID_STATE")
		return
	}

	code := gctx.Query("code")
	user, err := controller.google.HandleCallback(code)
	if err != nil {
		shared.BadReqResp(gctx, "Failed to handle Google callback", "GOOGLE_CALLBACK_FAILED")
		return
	}

	// Generate JWT token
	shared.SuccessResp(gctx, &LoginRes{
		Token: controller.authService.GetToken(user),
	})
}

func (controller *Controller) Logout(ctx *gin.Context) {
	// Clear any session cookies
	ctx.SetCookie("oauth_state", "", -1, "/", "", false, true)
}
