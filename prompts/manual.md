✅ ✅ create an api `deployment/:deploymentId/positions/:positionId/close` to close the open position

**Implementation Notes**:
1. Added a new route at `deployment/:deployment_id/positions/:positionID/close` to the existing ClosePosition controller method
2. Updated the controller to handle both route formats:
   - Original route: `/deployments/positions/:positionID/close`
   - New route: `/deployment/:deployment_id/positions/:positionID/close`
3. Added validation to ensure the position belongs to the specified deployment
4. Updated the API documentation to include both route formats
5. Fixed an issue with route parameter naming conventions (using `deployment_id` instead of `deploymentID`)

API implementation details:
- Added a new route in `/internal/modules/deployment/module.go`
- Updated the `ClosePosition` controller method to handle the new URL format 
- Added validation to ensure the position belongs to the specified deployment
- Updated the API documentation with the new endpoint format