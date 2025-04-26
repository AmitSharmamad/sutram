

-- Broker Management Tables
CREATE TABLE IF NOT EXISTS broker (
    id UUID,
    name VARCHAR(255),
    client_id VARCHAR(255),
    client_secret VARCHAR(255),
    auth_type VARCHAR(50),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS broker_user (
    id UUID,
    user_id UUID,
    broker_id UUID,
    token VARCHAR(255),
    expiry TIMESTAMP,

    FOREIGN KEY (broker_id) REFERENCES broker(id),
    FOREIGN KEY (user_id) REFERENCES "user"(id)
);
