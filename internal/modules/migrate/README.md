# Data Migration Module

This module fetches historical financial data from the Dhan API and stores it as CSV files.

## Features

- Fetches data for predefined securities from January 2020 until the current date
- Processes data in 2-month chunks to manage API requests efficiently
- Organizes data by year into separate CSV files with the format `{securityId}_{year}.csv`
- Includes proper formatting for dates and numerical values

## Usage

### Environment Setup

You can set your Dhan API access token as an environment variable:

```bash
export DHAN_ACCESS_TOKEN="your-access-token"
```

Alternatively, create a `.env` file in the project root with:

```
DHAN_ACCESS_TOKEN=your-access-token
```

### Running the Migration

```bash
# Build the migration tool
go build -o s-migrate cmd/migrate/main.go

# Run with environment variable
./s-migrate

# Or provide the token directly
./s-migrate -token="your-access-token"
```

## Output

The migration process will create CSV files in the `data/historical` directory with the following structure:

- Each file is named according to the pattern `{securityId}_{year}.csv`
- CSV headers: exchange_segment_id, open, high, low, close, volume, timestamp, created_at, updated_at
- Dates are formatted as "YYYY-MM-DD HH:MM:SS"

## Customizing Securities

To add or modify the securities for which data is fetched, edit the `NewDataMigrator` function in `internal/modules/migrate/migrate.go`.
