// GetAvailableProviders gets all available OAuth providers
// @Summary Get available OAuth providers
// @Description Retrieve a list of all available OAuth providers
// @Tags OAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=[]string} "List of available providers"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /oauth/providers [get]
func (c *OAuthController) GetAvailableProviders(ctx *gin.Context) {
	providers, err := c.oauthService.GetAvailableProviders(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Success: false,
			Error:   "Failed to get providers: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.ApiResponse{
		Success: true,
		Data:    providers,
	})
}