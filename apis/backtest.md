# Backtest API Endpoints

This document outlines the API endpoints for the backtest module of the algo trading platform.

## Authentication

All authenticated endpoints require authentication via Bearer token.

```
Authorization: Bearer <token>
```

## Endpoints

### 1. Health Check

Check if the backtest module is running.

```bash
curl -X GET \
  'https://api.example.com/backtest/health'
```

**Response:**
```json
{
  "status": "Backtest module is running"
}
```

### 2. Run Backtest

Executes a backtest for a given strategy using Server-Sent Events for real-time updates.

```bash
curl -X POST \
  'https://api.example.com/backtest/run?strategy_id=123e4567-e89b-12d3-a456-426614174000' \
  -H 'Authorization: Bearer <token>' \
  -H 'Accept: text/event-stream'
```

**Query Parameters:**
- `strategy_id`: UUID of the strategy to backtest (required)

**Response:**
Server-Sent Events stream with the following event types:

- `results`: Contains the backtest results
- `done`: Indicates completion with success/error status

**Event Examples:**
```
event: results
data: {"results": [...], "timestamp": "2025-06-02T10:30:00Z"}

event: done
data: {"message": "Backtest completed successfully", "error": false, "timestamp": "2025-06-02T10:30:05Z"}
```

### 3. Export Backtest to Excel

Exports backtest results to an Excel file for download.

```bash
curl -X GET \
  'https://api.example.com/backtest/export/excel?strategy_id=123e4567-e89b-12d3-a456-426614174000' \
  -H 'Authorization: Bearer <token>' \
  --output backtest_results.xlsx
```

**Query Parameters:**
- `strategy_id`: UUID of the strategy to export backtest data for (required)

**Response:**
- **Content-Type**: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- **Content-Disposition**: `attachment; filename=backtest_results_{strategy_id}_{timestamp}.xlsx`

**Excel File Structure:**

The exported Excel file contains separate sheets for each instrument backtested, with the following sections:

#### Instrument Information
- Symbol
- Exchange  
- Instrument ID

#### Positions Table
| Column | Description |
|--------|-------------|
| Position ID | Unique identifier for the position |
| State | Position state (NOT_ENTERED, ENTERED, EXITED) |
| Direction | Position direction (LONG, SHORT) |
| Entry Signal | Signal that triggered entry |
| Entry Time Signal | Time when entry signal was generated |
| Entry Price Signal | Price at entry signal |
| Exit Signal | Signal that triggered exit |
| Exit Time Signal | Time when exit signal was generated |
| Exit Price Signal | Price at exit signal |
| P&L Absolute | Absolute profit/loss amount |
| P&L Percentage | Percentage profit/loss |
| Position Size | Size of the position |
| Transaction Cost | Total transaction costs |
| Entry Fees | Fees paid on entry |
| Exit Fees | Fees paid on exit |
| Net P&L | Net profit/loss after fees |

#### Performance Metrics
| Metric | Description |
|--------|-------------|
| Total Return % | Overall return percentage |
| Annualized Return | Annualized return rate |
| Risk Adjusted Return | Risk-adjusted return metric |
| Sharpe Ratio | Risk-adjusted return ratio |
| Sortino Ratio | Downside deviation adjusted return |
| Max Drawdown | Maximum peak-to-trough decline |
| Max Drawdown Duration | Duration of maximum drawdown |
| Value at Risk | Potential loss at given confidence level |
| Standard Deviation | Volatility measure |
| Downside Deviation | Downside volatility measure |
| Total Trades | Total number of trades executed |
| Win Rate | Percentage of profitable trades |
| Average Win Size | Average size of winning trades |
| Average Loss Size | Average size of losing trades |
| Average Trade Duration | Average duration of trades |
| Profit Factor | Ratio of gross profit to gross loss |
| Average Position Size | Average size of positions |
| Max Position Size | Maximum position size |
| Position Turnover | Portfolio turnover rate |
| Average Exposure | Average market exposure |
| Total Commission | Total commission costs |
| Total Slippage | Total slippage costs |
| Transaction Costs % | Transaction costs as percentage |

**Error Responses:**

```json
{
  "error": "strategy_id query parameter is required"
}
```

```json
{
  "error": "Invalid strategy_id format"
}
```

```json
{
  "error": "User not authenticated"
}
```

```json
{
  "error": "Failed to retrieve backtest data: [error details]"
}
```

```json
{
  "error": "Failed to generate Excel file: [error details]"
}
```

## Usage Examples

### JavaScript (Browser)

```javascript
// Export backtest to Excel
async function exportBacktestToExcel(strategyId) {
  try {
    const response = await fetch(`/backtest/export/excel?strategy_id=${strategyId}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    });

    if (!response.ok) {
      throw new Error('Export failed');
    }

    // Create download link
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `backtest_results_${strategyId}_${new Date().getTime()}.xlsx`;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  } catch (error) {
    console.error('Export failed:', error);
  }
}
```

### Python

```python
import requests

def export_backtest_to_excel(strategy_id, auth_token):
    """Export backtest results to Excel file"""
    url = f"https://api.example.com/backtest/export/excel"
    headers = {
        'Authorization': f'Bearer {auth_token}'
    }
    params = {
        'strategy_id': strategy_id
    }
    
    response = requests.get(url, headers=headers, params=params)
    
    if response.status_code == 200:
        filename = f"backtest_results_{strategy_id}.xlsx"
        with open(filename, 'wb') as f:
            f.write(response.content)
        print(f"Excel file saved as {filename}")
    else:
        print(f"Export failed: {response.json().get('error', 'Unknown error')}")
```

## Notes

- The Excel export endpoint runs the backtest engine to generate fresh results
- Large datasets may take time to process and generate
- The Excel file includes formatted headers and proper column widths
- Multiple instruments will result in multiple sheets in the Excel file
- All timestamp fields are formatted as "YYYY-MM-DD HH:MM:SS"
- Numeric values are preserved with appropriate precision
