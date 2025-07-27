-- Create instrument data table

CREATE TABLE IF NOT EXISTS sutram.instrument (
    exchange_security_id Int32,
    open Float32,
    high Float32,
    low Float32,
    close Float32,
    volume Float32,
    `timestamp` DateTime64(6, 'Asia/Kolkata'),
    `created_at` DateTime64(6, 'Asia/Kolkata'),
    `updated_at` DateTime64(6, 'Asia/Kolkata')
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(`timestamp`)
ORDER BY (exchange_security_id, `timestamp`);
