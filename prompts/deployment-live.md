do the above actions in the following order

```
1. abstract out all the core logic into a module `core`
2. backtesting should still work with the same flow without any functional changes or flow changes (make double sure of this)
3. make use of all the core methods to process the events from live, example `isStopLossTriggered`, `isTargetProfitTriggered` etc.
4. use the existing CalculateAll method in an indicator for backtesting to calculate all the indicator values at once
5. use the existing Calculate method in an indicator for live trading
6. make sure that there is no discrepency in the logic between Calculate and CalculateAll
7. upon deploying a deployment, make sure that all the metadata is stored in the deployment_position_log table which are as follows
    a. calculate the values for all the indicators in the entry and exit expressions and store it in the metadata field of the deployment_position_log
    b. when a data event occurs and Run method is called on `deployment` object do the following
        i. run the checks on timeframe, entry start time and entry end time to see if the execution should go ahead, otherwise short circuit
        ii. look for the exact static previous point of the timestamp and see if there's indicator calculation available
        iii. if yes, pick the indicators and calculate the current indicator value with only that
        iv. if not, calculate the current indicator value from scratch like brute force
```

any database migrations needed, don't write migration files
instead update the existing schema as requried and make changes to the respective repos, services and apis