-- Broker Management Tables
CREATE TABLE IF NOT EXISTS broker (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    client_id VARCHAR(255),
    client_secret VARCHAR(255),
    auth_type VARCHAR(50),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS broker_user (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID,
    broker_id UUID,
    token VARCHAR(255),
    expiry TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (broker_id) REFERENCES broker(id),
    FOREIGN KEY (user_id) REFERENCES "user"(id)
);

CREATE TABLE "order" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    position_id UUID NOT NULL,
    exchange exchange NOT NULL,
    segment segment NOT NULL,
    instrument VARCHAR(50) NOT NULL,
    transaction_type transaction_type NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    executed_at TIMESTAMP NOT NULL,
    state order_state NOT NULL,
    product_type product_type NOT NULL,
    order_type order_type NOT NULL,
    broker_order_id VARCHAR(100),
    broker_order_response JSONB,
    algo_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES "user"(id)
);

-- Indexes
CREATE INDEX idx_order_user_id ON "order"(user_id);
CREATE INDEX idx_order_instrument ON "order"(instrument);
CREATE INDEX idx_order_state ON "order"(state);
