
# Instrument Master Migration Guide

This guide explains how to migrate instrument master data from Dhan's CSV format to your ClickHouse database.

## Prerequisites

1. ClickHouse database should be running and accessible
2. The `instrument_master` table should be created using the provided migration
3. CSV file with Dhan's instrument data should be available

## Usage

### Command Line Interface

The migrate command supports three types of migrations:

1. **Historical Data Migration** (existing functionality)
2. **Instrument Master Migration** (new functionality)
3. **Generate Sample Excel** (for testing)

#### Instrument Master Migration

```bash
# Basic usage - migrate from CSV file
./bin/migrate -type=instrument-master -csv=path/to/dhan_instruments.csv

# Truncate existing data before migration (recommended for fresh imports)
./bin/migrate -type=instrument-master -csv=path/to/dhan_instruments.csv -truncate=true

# Test with a small sample first
head -100 your_large_file.csv > test_sample.csv
./bin/migrate -type=instrument-master -csv=test_sample.csv -truncate=true
```

#### Historical Data Migration

```bash
# Requires Dhan access token
./bin/migrate -type=historical -token=your_dhan_access_token
```

#### Generate Sample Excel (for testing)

```bash
./bin/migrate -type=generate-sample
```

### CSV File Format

Your CSV file should contain columns with Dhan's naming convention. The following columns are supported:

**Required Columns:**
- `SECURITY_ID` → `exchange_security_id`
- `EXCH_ID` → `exchange` (NSE, BSE, MCX)
- `SEGMENT` → `segment` (E=EQUITY, D=DERIVATIVE, M=COMMODITY, C=CURRENCY)
- `SYMBOL_NAME` → `symbol`
- `INSTRUMENT` → `instrument`
- `DISPLAY_NAME` → `instrument_name`
- `LOT_SIZE` → `lot_units` (handles decimal values like "1.0")
- `TICK_SIZE` → `tick_size`

**Optional Columns:**
- `ISIN` → `isin` (handles "NA" values as null)
- `UNDERLYING_SECURITY_ID` → `underlying_security_id`
- `UNDERLYING_SYMBOL` → `underlying_symbol`
- `SM_EXPIRY_DATE` → `expiry`
- `EXPIRY_FLAG` → `expiry_flag`
- `STRIKE_PRICE` → `strike_price` (handles negative values)
- `OPTION_TYPE` → `option_type`
- `INSTRUMENT_TYPE` → `instrument_type`
- Various margin and bracket order parameters

### Data Transformations

The migration automatically handles:

1. **Enum Mappings:**
   - Exchange: NSE, BSE, MCX
   - Segment: E→EQUITY, D→DERIVATIVE, M→COMMODITY, C→CURRENCY
   - Instrument types: EQUITY, FUTSTK→STOCK_FUTURES, FUTCUR→CURRENCY_FUTURES, etc.

2. **Data Type Conversions:**
   - String to integer for IDs
   - Decimal string to float for prices and percentages (handles "1.0" format)
   - Date parsing for expiry dates (multiple formats supported)
   - "NA" values converted to database null
   - Null handling for optional fields

3. **Default Values:**
   - `created_at` and `updated_at` are set to current timestamp
   - `expiry` defaults to 2099-12-31 for non-derivative instruments
   - Flag fields default to appropriate values if not provided

### Error Handling

- Invalid rows are logged as warnings and skipped
- Migration continues even if some rows fail
- Detailed error messages help identify data issues
- Database transactions ensure data consistency
- Progress reporting shows processed vs skipped rows

### Performance

- Data is processed in batches of 1000 records for optimal performance
- Bulk insert operations minimize database round trips
- Progress is reported during migration with batch timing information
- Large files (180K+ records) are handled efficiently
- **Timing Metrics**: Shows total migration time, average time per row, and batch processing times

## Environment Variables

Set the ClickHouse connection string:
```bash
export CH_URL="clickhouse://user:password@host:port/database"
```

## Examples

```bash
# Example 1: Fresh migration with data cleanup
./bin/migrate -type=instrument-master -csv=api-scrip-master-detailed_AUG11.csv -truncate=true

# Example 2: Test with small sample first
head -1000 api-scrip-master-detailed_AUG11.csv > test_1000.csv
./bin/migrate -type=instrument-master -csv=test_1000.csv -truncate=true

# Example 3: Append new instruments to existing data
./bin/migrate -type=instrument-master -csv=new_instruments.csv

# Example 4: Historical price data migration
export DHAN_ACCESS_TOKEN="your_token_here"
./bin/migrate -type=historical
```

## Real-World Usage

Successfully tested with Dhan's actual CSV file containing 180K+ instrument records:
- Processed large CSV files efficiently
- Handled decimal lot sizes (1.0, 25.0, etc.)
- Managed "NA" values in ISIN fields
- Processed negative strike prices for currency futures
- Skipped invalid records while continuing migration

## Troubleshooting

1. **Database Connection Issues:**
   - Verify CH_URL environment variable
   - Check ClickHouse server status
   - Ensure database and table exist

2. **CSV File Issues:**
   - Verify file path and permissions
   - Check CSV format and column headers
   - Ensure required columns are present
   - Test with small sample first

3. **Data Validation Issues:**
   - Check logs for specific row errors
   - Verify enum values match expected mappings
   - Ensure required fields are not empty
   - Use test samples to identify data issues

4. **Performance Issues:**
   - For very large files, consider splitting into smaller chunks
   - Monitor ClickHouse server resources during migration
   - Use truncate=true for fresh imports to avoid duplicate key issues
