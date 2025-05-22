// HandleCallback handles the OAuth callback
// @Summary Handle OAuth callback
// @Description Process the OAuth callback from the provider
// @Tags OAuth
// @Accept json
// @Produce json
// @Param provider path string true "OAuth provider (google, github, etc.)"
// @Param code query string true "Authorization code"
// @Param state query string true "State parameter for CSRF protection"
// @Param error query string false "Error message from provider"
// @Success 200 {object} dto.ApiResponse{data=map[string]interface{}} "User and token"
// @Success 302 {string} string "Redirect to home page with cookie"
// @Failure 400 {object} dto.ApiResponse "Invalid parameters"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /oauth/{provider}/callback [get]
func (c *OAuthController) HandleCallback(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider == "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Provider is required",
		})
		return
	}

	// Get callback params
	code := ctx.Query("code")
	state := ctx.Query("state")
	errorMsg := ctx.Query("error")

	if errorMsg != "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "OAuth error: " + errorMsg,
		})
		return
	}

	if code == "" || state == "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Code and state are required",
		})
		return
	}

	// Handle callback
	user, token, err := c.oauthService.HandleCallback(ctx.Request.Context(), provider, code, state)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to handle callback: " + err.Error(),
		})
		return
	}

	// Set JWT token as a cookie
	SetTokenCookie(ctx, token, user.ID.String())

	// Check if this is an API request or a browser request
	if ctx.GetHeader("Content-Type") == "application/json" || ctx.GetHeader("Accept") == "application/json" {
		// API request, return JSON
		ctx.JSON(http.StatusOK, dto.ApiResponse{
			Success: true,
			Data: gin.H{
				"token": token,
				"user":  user,
			},
		})
	} else {
		// Browser request, redirect to home page
		ctx.Redirect(http.StatusFound, "/")
	}
}