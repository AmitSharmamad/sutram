# WebSocket Connection Troubleshooting Guide

## Issue: "websocket: bad handshake" Error

### Root Cause Analysis
The "bad handshake" error typically occurs due to authentication or connection parameter issues when connecting to the Dhan WebSocket API.

## Implemented Fixes

### 1. Enhanced Error Logging
Added comprehensive error logging to help diagnose connection issues:
- HTTP response status and headers
- Response body content for error details
- JWT token validation with expiry checking
- Parameter validation before connection attempt

### 2. JWT Token Validation
Added `validateJWTToken()` function that:
- Validates JWT token format (3 parts separated by dots)
- Decodes and parses the payload
- Checks token expiry against current time
- Validates issuer ("dhan")
- Logs token expiry time and client ID

### 3. Parameter Validation
Enhanced validation for required environment variables:
- `WS_TOKEN`: Dhan access token (JWT format)
- `WS_CLIENT_ID`: Dhan client ID
- Proper error handling with retry logic

## Common Causes and Solutions

### 1. Expired JWT Token
**Symptoms:**
- "bad handshake" error
- HTTP 401 Unauthorized response

**Solution:**
- Generate a new access token from Dhan API
- Update `WS_TOKEN` in your `.env` file
- The token in your current `.env` has expiry timestamp `1750134166`

**Token Generation:**
```bash
# Use Dhan API to generate a new token
curl -X POST "https://api.dhan.co/v2/login" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "your_user_id",
    "password": "your_password",
    "twoFA": "your_2fa_code"
  }'
```

### 2. Invalid Client ID
**Symptoms:**
- "bad handshake" error
- Parameters mismatch

**Solution:**
- Verify `WS_CLIENT_ID` matches your Dhan client ID
- Should be numeric (e.g., "1105711358")

### 3. Network/Firewall Issues
**Symptoms:**
- Connection timeout
- DNS resolution failure

**Solution:**
- Check internet connectivity
- Verify firewall allows WebSocket connections to `api-feed.dhan.co`
- Try connecting from a different network

### 4. Incorrect URL Parameters
**Current URL Format:**
```
wss://api-feed.dhan.co?version=2&token=<JWT_TOKEN>&clientId=<CLIENT_ID>&authType=2
```

**Validation:**
- `version=2` (required)
- `authType=2` (required)
- Token and clientId must be URL-encoded if containing special characters

## Debugging Steps

### 1. Check Current Token Status
```bash
# Run the application and check logs for JWT validation
./bin/live
```

Look for log messages:
- "JWT token is valid. Expires at: ..."
- "JWT Token validation failed: ..."

### 2. Test Token Manually
You can decode your JWT token at https://jwt.io to check:
- Expiry time (`exp` field)
- Client ID (`dhanClientId` field)
- Issuer (`iss` should be "dhan")

### 3. Network Connectivity Test
```bash
# Test basic connectivity to Dhan API
curl -I https://api-feed.dhan.co

# Test WebSocket endpoint (should return HTTP 400 without proper parameters)
curl -I "wss://api-feed.dhan.co"
```

## Updated Environment Variables

Ensure your `.env` file has:
```properties
WS_UPSTREAM_URL=wss://api-feed.dhan.co
WS_TOKEN=<your_valid_jwt_token>
WS_CLIENT_ID=<your_dhan_client_id>
```

## Next Steps

1. **Generate New Token**: If token is expired, generate a new one via Dhan API
2. **Update .env**: Replace `WS_TOKEN` with the new token
3. **Restart Application**: The enhanced logging will provide detailed error information
4. **Monitor Logs**: Check for specific error messages in the enhanced logging

## Enhanced Logging Output

With the new changes, you'll see detailed logs like:
```
Connecting to upstream WebSocket: wss://api-feed.dhan.co
Using Client ID: 1105711358
JWT token is valid. Expires at: 2025-12-15T10:30:00Z
Full WebSocket URL: wss://api-feed.dhan.co?version=2&token=eyJ0eXAiOiJKV1QiLCJh...&clientId=1105711358&authType=2
```

Or in case of errors:
```
JWT Token validation failed: JWT token has expired. Expiry: 1750134166, Current: 1750134200
HTTP Response Status: 401 Unauthorized
Response Body: {"error": "Invalid or expired token"}
```

This will help pinpoint the exact cause of the connection failure.
