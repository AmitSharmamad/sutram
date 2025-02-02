USE sutram;

-- Broker Management Tables
CREATE TABLE IF NOT EXISTS broker (
    id UUID DEFAULT generateUUIDv4(),
    name String,
    client_id String,
    client_secret String,
    auth_type String,
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS broker_user (
    user_id UUID,
    broker_id UUID,
    token String,
    expiry DateTime,
    PRIMARY KEY (user_id, broker_id)
) ENGINE = MergeTree;