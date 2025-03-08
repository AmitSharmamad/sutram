### Slippage Cost

Slippage cost is the monetary difference between your expected/intended execution price and the actual price at which your order was filled, multiplied by the position size, representing the additional cost incurred due to market movement during order execution.

```
slippage_cost = position_size * (actual_execution_price - intended_price)
```


### Transaction Cost

Transaction cost represents the total cost of executing a trade, including all fees, commissions, and price differences between expected and actual execution prices.

```
transaction_cost = entry_fees + exit_fees + slippage_cost
```
