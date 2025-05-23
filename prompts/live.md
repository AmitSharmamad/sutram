create a `Initialize` method in deployment_service which should do the following

1. call this method from inside cmd/live.go
2. read all the active deployments for all the users
3. group the deployments by instrument id in a hashmap
4. once there's an event received for an instrument, do the following
    a. persist the data in the database
    b. find the respective instrument from the above group
    c. call a function `Run` for all the deployments against it


Add a method `Run` on the deployment object and 
add all the logic of executing the strategy against the data
reuse all the logic that is available in `backtest/engine.go`
refactor the logic in engine.go carefully and reuse all the core logic