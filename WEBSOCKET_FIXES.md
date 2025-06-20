# WebSocket Contract Fixes - Dhan API v2

## Overview
Fixed the WebSocket contracts in `cmd/live/main.go` to comply with the official Dhan API v2 Live Market Feed documentation at: https://dhanhq.co/docs/v2/live-market-feed

## Changes Made

### 1. Binary Packet Parsing Corrections

#### Ticker Packet (Feed Code 2)
- **Fixed**: Correct byte positions for Last Traded Price (LTP) and Last Trade Time (LTT)
- **Structure**: 
  - Bytes 8-12: LTP (float32)
  - Bytes 12-16: LTT (int32)

#### Quote Packet (Feed Code 4)
- **Fixed**: Proper parsing of all quote data fields according to documentation
- **Structure**:
  - Bytes 8-12: Latest Traded Price (float32)
  - Bytes 12-14: Last Traded Quantity (int16)
  - Bytes 14-18: Last Trade Time (int32)
  - Bytes 18-22: Average Trade Price (float32)
  - Bytes 22-26: Volume (int32)
  - Bytes 26-30: Total Sell Quantity (int32)
  - Bytes 30-34: Total Buy Quantity (int32)
  - Bytes 34-38: Day Open Value (float32)
  - Bytes 38-42: Day Close Value (float32)
  - Bytes 42-46: Day High Value (float32)
  - Bytes 46-50: Day Low Value (float32)

#### Full Packet (Feed Code 8)
- **Fixed**: Updated to handle the complete 163-byte full packet structure
- **Changes**:
  - Updated minimum packet size check from 63 to 163 bytes
  - Correct OHLC field positions
  - Added comments for Open Interest fields
  - Market depth data location documented (bytes 62-162)

#### Open Interest Packet (Feed Code 5)
- **Enhanced**: Added proper logging with security ID and exchange information
- **Structure**: Bytes 8-12 contain Open Interest value

#### Previous Close Packet (Feed Code 6)
- **Enhanced**: Added proper logging with security ID and exchange information
- **Structure**:
  - Bytes 8-12: Previous day closing price (float32)
  - Bytes 12-16: Previous day Open Interest (int32)

### 2. Response Header Parsing
All packet types now properly extract:
- Feed Response Code (byte 0)
- Message Length (bytes 1-3)
- Exchange Segment (byte 3)
- Security ID (bytes 4-8)

### 3. Connection URL Format
Added documentation comment confirming the correct URL format:
```
wss://api-feed.dhan.co?version=2&token=eyxxxxx&clientId=100xxxxxxx&authType=2
```

### 4. JSON Request Structure
The subscription request structure already matches the documentation:
```json
{
    "RequestCode": 15,
    "InstrumentCount": 2,
    "InstrumentList": [
        {
            "ExchangeSegment": "NSE_EQ",
            "SecurityId": "1333"
        }
    ]
}
```

## Key Compliance Features

### 1. Subscription Limits
- ✅ Handles up to 100 instruments per subscription message
- ✅ Supports chunking for larger subscription requests
- ✅ Can handle up to 5000 instruments per connection

### 2. Connection Management
- ✅ Proper ping-pong handling for connection keep-alive
- ✅ Automatic reconnection logic
- ✅ Resubscription after reconnection

### 3. Data Processing
- ✅ Binary message parsing for high-speed data
- ✅ Proper fallback handling for missing OHLC values
- ✅ Correct timestamp conversion from Unix format

## Testing Recommendations

1. **Environment Variables**: Ensure these are set correctly:
   - `WS_UPSTREAM_URL`: Default "wss://api-feed.dhan.co"
   - `WS_TOKEN`: Your Dhan access token
   - `WS_CLIENT_ID`: Your Dhan client ID

2. **Subscription Testing**: 
   - Test with various exchange segments (NSE_EQ, BSE_EQ, NSE_FNO, etc.)
   - Verify subscription chunking with >100 instruments
   - Test reconnection scenarios

3. **Data Validation**:
   - Compare received data with Dhan's trading platform
   - Verify timestamp accuracy
   - Check OHLC data consistency

## Performance Optimizations

- Binary packet processing for minimal latency
- Efficient memory usage with proper buffer sizes
- Concurrent client handling with goroutines
- Optimized data broadcasting to multiple clients

## Error Handling

- Connection failure recovery with automatic retry
- Invalid packet size validation
- Graceful handling of malformed data
- Proper logging for debugging

The implementation now fully complies with Dhan's Live Market Feed WebSocket API v2 specifications.
