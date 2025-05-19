# OAuth Configuration

The application supports OAuth authentication for user registration and login. Currently, the following providers are supported:

- Google

## Setting up Google OAuth

1. Create a Google OAuth 2.0 client ID in the [Google Cloud Console](https://console.developers.google.com/).
2. Configure the OAuth consent screen with your application information.
3. Add the authorized redirect URI: `http://localhost:9090/api/auth/oauth/google/callback` (adjust the hostname and port as needed for your environment).
4. Copy the client ID and client secret.
5. Update your `.env` file with the following variables:

```bash
GOOGLE_OAUTH_CLIENT_ID=your_client_id
GOOGLE_OAUTH_CLIENT_SECRET=your_client_secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:9090/api/auth/oauth/google/callback
```

## OAuth Endpoints

- **Initiate Google Login**: GET `/api/auth/google/login`
- **Google OAuth Callback**: GET `/api/auth/oauth/google/callback?code={code}&state={state}`
- **Available OAuth Providers**: GET `/api/auth/oauth/providers`
- **User OAuth Connections**: GET `/api/user/oauth/connections` (requires authentication)
- **Disconnect OAuth Provider**: DELETE `/api/user/oauth/:provider` (requires authentication)

## OAuth Flow

1. User initiates login by navigating to `/api/auth/google/login`
2. User is redirected to Google's authentication screen
3. After successful authentication, Google redirects back to the callback URL
4. The application verifies the authentication and creates or retrieves the user account
5. User is logged in and JWT token is provided
