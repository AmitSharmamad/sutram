// InitiateLogin initiates the OAuth login flow
// @Summary Initiate OAuth login
// @Description Start the OAuth login flow for a specific provider
// @Tags OAuth
// @Accept json
// @Produce json
// @Param provider path string true "OAuth provider (google, github, etc.)"
// @Success 200 {object} dto.ApiResponse{data=map[string]string} "Login URL"
// @Success 302 {string} string "Redirect to provider login page"
// @Failure 400 {object} dto.ApiResponse "Invalid provider"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /oauth/{provider}/login [get]
func (c *OAuthController) InitiateLogin(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider == "" {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Provider is required",
		})
		return
	}

	// Get login URL
	loginURL, err := c.oauthService.GetLoginURL(ctx.Request.Context(), provider)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get login URL: " + err.Error(),
		})
		return
	}

	// Check if this is an API request or a browser request
	if ctx.GetHeader("Content-Type") == "application/json" || ctx.GetHeader("Accept") == "application/json" {
		// API request, return JSON
		ctx.JSON(http.StatusOK, dto.ApiResponse{
			Success: true,
			Data: gin.H{
				"url": loginURL,
			},
		})
	} else {
		// Browser request, redirect
		ctx.Redirect(http.StatusFound, loginURL)
	}
}