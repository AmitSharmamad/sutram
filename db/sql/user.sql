USE sutram;

-- User Management Tables
CREATE TABLE IF NOT EXISTS user (
    id UUID DEFAULT generateUUIDv4(),
    name String,
    username String,
    email String,
    phone_number String,
    password String,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    is_active Boolean DEFAULT true,
    PRIMARY KEY (id)
) ENGINE = MergeTree;