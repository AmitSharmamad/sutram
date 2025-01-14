
-- Template for stock tables (replace SYMBOL with actual stock symbol)
CREATE TABLE IF NOT EXISTS stock_{SYMBOL} (
    timestamp DateTime,
    open Decimal64(2),
    high Decimal64(2),
    close Decimal64(2),
    low Decimal64(2),
    volume UInt64,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (timestamp)
) ENGINE = MergeTree;
