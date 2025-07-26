# Option Leg Module

This module provides comprehensive functionality for managing option legs in multi-leg options trading strategies.

## Features

- **Domain Model**: Complete option leg model with all required fields and validation
- **Repository Interface**: Abstract data access layer for option legs
- **PostgreSQL Implementation**: Concrete implementation using PostgreSQL
- **Service Layer**: Business logic for option leg operations
- **HTTP API**: RESTful endpoints for option leg management
- **Validation**: Comprehensive validation for different segment types (F&O vs EQ)

## Database Schema

The module uses the `option_leg` table with the following key fields:

- **id**: UUID primary key
- **strategy_id**: Reference to the parent strategy
- **leg_sequence**: Order of legs in the strategy (1, 2, 3, etc.)
- **select_index**: Index/Underlying (e.g., 'Nifty', 'BankNifty')
- **segment**: Trading segment (F&O or EQ)
- **option_type**: CALL or PUT (only for F&O segment)
- **action_type**: BUY or SELL
- **strike_price_type**: ATM, ITM, OTM, or Custom
- **strike_price**: Actual strike price
- **total_lot**: Number of lots
- **order_type**: MARKET or LIMIT
- **expiry_type**: Weekly, Monthly, or Quarterly
- **leg_status**: ACTIVE, INACTIVE, EXECUTED, or CANCELLED

## API Endpoints

### Core Operations
- `POST /api/v1/option-legs` - Create option leg
- `GET /api/v1/option-legs/:id` - Get option leg by ID
- `PUT /api/v1/option-legs/:id` - Update option leg
- `DELETE /api/v1/option-legs/:id` - Delete option leg

### Status Management
- `PATCH /api/v1/option-legs/:id/status` - Update option leg status

### Search and Filter
- `GET /api/v1/option-legs/search` - Search with filters
- `GET /api/v1/option-legs/index/:index` - Get by index/underlying
- `GET /api/v1/option-legs/status/:status` - Get by status

### Strategy Integration
- `GET /api/v1/strategies/:strategy_id/option-legs` - Get all legs for strategy
- `DELETE /api/v1/strategies/:strategy_id/option-legs` - Delete all legs for strategy

## Usage

### Creating an Option Leg

```json
POST /api/v1/option-legs
{
  "strategy_id": "uuid",
  "leg_sequence": 1,
  "select_index": "Nifty",
  "segment": "F&O",
  "option_type": "CALL",
  "action_type": "BUY",
  "strike_price_type": "ATM",
  "total_lot": 1,
  "order_type": "MARKET",
  "expiry_type": "Weekly"
}
```

### Updating Option Leg Status

```json
PATCH /api/v1/option-legs/:id/status
{
  "status": "EXECUTED"
}
```

### Searching Option Legs

```
GET /api/v1/option-legs/search?strategy_id=uuid&segment=F&O&option_type=CALL&page=1&per_page=20
```

## Integration

To integrate the option leg module into your application:

1. Import the module: `import backtestModule "algo/internal/modules/backtest"`
2. Create and initialize: `optionLegMod := optionLegModule.New()`
3. Set database connection: `optionLegMod.SetDB(db)`
4. Register routes: `optionLegMod.RegisterRoutes(router)`
5. Access service: `service := optionLegMod.GetOptionLegService()`

## Validation Rules

### F&O Segment
- Must have `option_type` (CALL or PUT)
- Must have `strike_price` for custom strike price type
- Can have strike offset for ATM/ITM/OTM calculations

### EQ Segment  
- Must NOT have `option_type`
- Must NOT have `strike_price`
- Focused on equity trading

### General Rules
- `leg_sequence` must be > 0 and unique per strategy
- `total_lot` must be > 0
- Total quantity automatically calculated as `total_lot * lot_size`

## Dependencies

- PostgreSQL database with option_leg table
- Go UUID package for ID generation
- Gin framework for HTTP endpoints
- Standard Go SQL package for database operations
