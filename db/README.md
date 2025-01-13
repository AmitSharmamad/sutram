# tables

---

## user management

- user

  - id [ PK ]
  - name
  - username
  - email
  - phone_number
  - password
  - created_at
  - updated_at

- portfolio

  - id [ PK ]
  - amount
  - profit
  - loss
  - created_at
  - updated_at

### A portfolio can have multiple logs

- portfolio_log

  - id [ PK ]
  - portfolio_id
  - action_type [ FUND_ADDED, FUND_WITHDRAWN, FUND_REFUNDED ]
  - created_at
  - updated_at

- notification

  - id
  - user_id [ user(id) ]

- setting

  - id
  - user_id [ user(id) ]
  - config [ mode (LIGHT, DARK), etc... ]
  - created_at
  - updated_at

---

## strategy management

- strategy

  - id [ PK ]
  - name
  - description
  - symbol
  - instrument_type [ EQUITY, FUTURES, OPTIONS ]
  - position_type [ BUY, SELL ]
  - entry_expression
  - exit_expression
  - stop_loss
  - trailing_stop_loss
  - max_trades
  - created_at
  - updated_at

- strategy_history
  - id [ PK ]
  - strategy_id [ strategy(id) ]
  - change [ json(strategy) ]
  - created_at
  - updated_at

---

## stock management - data integration

### Every stock has its own table

- stock\_{SYMBOL} [ stock_AAPL, stock_TSLA ]

  - timestamp [ PK ]
  - open
  - high
  - close
  - low
  - volume
  - created_at
  - updated_at

---

## backtesting engine

### A strategy can have multiple backtesting runs

- bt_run

  - id [ PK ]
  - strategy_id
  - state [ not_picked, started, in_progress, ended ]
  - qty
  - period_start_time
  - period_end_time
  - run_start_time
  - run_stop_time
  - created_at
  - updated_at

### A run can have multiple positions (trades)

- bt_position

  - id [ PK ]
  - run_id
  - state [ enter, entered, hold, exited ]
  - entry_price
  - entry_time
  - exit_price
  - exit_time
  - exit_reason
  - created_at
  - updated_at

- bt_position_log

  - id [ PK ]
  - position_id [ position(id) ]
  - state [ enter, entered, hold, exited ]
  - created_at

---

## paper trading

### A strategy can have multiple paper trading runs

- pt_run

  - id [ PK ]
  - strategy_id
  - state [ not_picked, started, in_progress, ended ]
  - start_time
  - stop_time
  - created_at
  - updated_at

### A run can have multiple positions (trades)

- pt_position

  - id [ PK ]
  - run_id
  - state [ enter, entered, hold, exited ]
  - entry_price
  - entry_time
  - exit_price
  - exit_time
  - exit_reason
  - created_at
  - updated_at

- pt_position_log

  - id [ PK ]
  - position_id [ position(id) ]
  - state [ enter, entered, hold, exited ]
  - created_at

---

## live trading

### A strategy can have multiple paper trading runs

- lt_run

  - id [ PK ]
  - strategy_id
  - state [ not_picked, started, in_progress, ended ]
  - start_time
  - stop_time
  - created_at
  - updated_at

### A run can have multiple positions (trades)

- lt_position

  - id [ PK ]
  - run_id
  - state [ enter, entered, hold, exited ]
  - entry_price
  - entry_time
  - exit_price
  - exit_time
  - exit_reason
  - created_at
  - updated_at

- lt_position_log

  - id [ PK ]
  - position_id [ position(id) ]
  - state [ enter, entered, hold, exited ]
  - created_at

---

## broker management - broker integration

- broker

  - id [ PK ]
  - name
  - client_id
  - client_secret
  - auth_type [ oauth2.0 ]

- broker_user

  - user_id [ user (id) ]
  - broker_id [ broker (id) ]
  - token
  - expiry
  - is_active
