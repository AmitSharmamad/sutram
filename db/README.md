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
  - is_active

- portfolio

  - id [ PK ]
  - user_id [ user(id) ]
  - capital
  - returns
  - created_at
  - updated_at

- setting

  - id
  - user_id [ user(id) ]
  - config [ json(mode (LIGHT, DARK), etc...) ]
  - created_at
  - updated_at

- notification

  - id
  - user_id [ user(id) ]
  - type [ BT_DONE, LT_TRADE_ENTERED, LT_TRADE_EXITED, etc ]
  - status [ DELIVER, DELIVERED ]
  - receiver [ MOBILE, EMAIL ]
  - content

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

- stock_master

  - id [ PK ]
  - exchange [ NSE, BSE, etc... ]
  - name
  - symbol
  - is_active
  - created_at
  - updated_at

---

## strategy management

### A User can have multiple strategies

- strategy

  - id [ PK ]
  - user_id [ user(id) ]
  - name
  - description
  - symbol
  - instrument_type [ EQUITY, FUTURES, OPTIONS ]
  - stock [ stock_master(id) ]
  - position_type [ BUY, SELL ]
  - entry_expression
  - exit_expression
  - stop_loss
  - trailing_stop_loss
  - max_trades
  - created_at
  - updated_at
  - is_active

- strategy_history
  - id [ PK ]
  - strategy_id [ strategy(id) ]
  - change [ json(strategy) ]
  - created_at
  - updated_at

---

## backtest engine

### A strategy can have multiple backtesting runs

- bt_run

  - id [ PK ]
  - strategy_id [ strategy(id) ]
  - state [ start, in_progress, ended ]
  - qty
  - period_start_time
  - period_end_time
  - started_time
  - ended_time
  - created_at
  - updated_at

### A run can have multiple positions (trades)

- bt_position

  - id [ PK ]
  - run_id [ run(id) ]
  - state [ enter, entered, hold, exited ]
  - entry_price
  - entry_time
  - exit_price
  - exit_time
  - exit_reason
  - created_at
  - updated_at

### A Position can have multiple logs

- bt_position_log

  - id [ PK ]
  - position_id [ position(id) ]
  - state [ enter, entered, hold, exited ]
  - created_at

---

## paper trading

### A strategy can have multiple backtesting runs

- pt_run

  - id [ PK ]
  - strategy_id [ strategy(id) ]
  - state [ start, in_progress, ended ]
  - qty
  - period_start_time
  - period_end_time
  - started_time
  - ended_time
  - created_at
  - updated_at

### A run can have multiple positions (trades)

- pt_position

  - id [ PK ]
  - run_id [ run(id) ]
  - state [ enter, entered, hold, exited ]
  - entry_price
  - entry_time
  - exit_price
  - exit_time
  - exit_reason
  - created_at
  - updated_at

### A Position can have multiple logs

- pt_position_log

  - id [ PK ]
  - position_id [ position(id) ]
  - state [ enter, entered, hold, exited ]
  - created_at

---

## live trading

### A strategy can have multiple backtesting runs

- lt_run

  - id [ PK ]
  - strategy_id [ strategy(id) ]
  - state [ start, in_progress, ended ]
  - qty
  - period_start_time
  - period_end_time
  - started_time
  - ended_time
  - created_at
  - updated_at

### A run can have multiple positions (trades)

- lt_position

  - id [ PK ]
  - run_id [ run(id) ]
  - state [ enter, entered, hold, exited ]
  - entry_price
  - entry_time
  - exit_price
  - exit_time
  - exit_reason
  - created_at
  - updated_at

### A Position can have multiple logs

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
