-- Create instrument index

CREATE INDEX IF NOT EXISTS idx_instrument_exchange_security_id
ON sutram.instrument (exchange_security_id) TYPE bloom_filter GRANULARITY 4;
