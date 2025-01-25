clickhouse-client --multiquery <sutram.sql

clickhouse-client --query="INSERT INTO sutram.stock_TESLA FORMAT CSV" <ohlcv_data.csv

clickhouse-client --query="SELECT COUNT(*) FROM sutram.stock_TESLA;"
