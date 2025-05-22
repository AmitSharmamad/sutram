// SetTokenCookie sets the JWT token as a cookie
func SetTokenCookie(ctx *gin.Context, token, userID string) {
	cookieMaxAge := 3600 * 24 * 30 // 30 days
	domain := ""                   // Use the domain from which the request was made
	
	// Get domain from env var if available
	if envDomain := os.Getenv("COOKIE_DOMAIN"); envDomain != "" {
		domain = envDomain
	}

	secure := false
	if os.Getenv("ENVIRONMENT") == "production" {
		secure = true
	}

	// Set the token cookie
	ctx.SetCookie(
		"token",                      // name
		token,                        // value
		cookieMaxAge,                 // max age in seconds
		"/",                          // path
		domain,                       // domain
		secure,                       // secure
		true,                         // http only
	)

	// Set the user ID cookie (not HTTP only so it's accessible to JS)
	ctx.SetCookie(
		"user_id",                    // name
		userID,                       // value
		cookieMaxAge,                 // max age in seconds
		"/",                          // path
		domain,                       // domain
		secure,                       // secure
		false,                        // http only
	)
}