clickhouse-client --multiquery <sutram.sql

clickhouse-client --multiquery <user.sql
clickhouse-client --multiquery <stock.sql
clickhouse-client --multiquery <strategy.sql
clickhouse-client --multiquery <bt.sql
clickhouse-client --multiquery <pt.sql
clickhouse-client --multiquery <lt.sql
clickhouse-client --multiquery <broker.sql

clickhouse-client --query="INSERT INTO sutram.stock_TESLA FORMAT CSV" <ohlcv_data.csv

clickhouse-client --query="SELECT COUNT(*) FROM sutram.stock_TESLA;"

clickhouse-client --multiquery <default.sql
