USE sutram;

clickhouse-client --query="INSERT INTO instrument_master FORMAT CSV" < instrument_master.csv

clickhouse-client --query="INSERT INTO instrument FORMAT CSV" < instrument.csv