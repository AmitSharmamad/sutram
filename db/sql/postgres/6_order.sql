
CREATE TABLE "order" (
    id UUID PRIMARY KEY,
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

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Indexes
CREATE INDEX idx_user_id ON "order"(user_id);
CREATE INDEX idx_instrument ON "order"(instrument);
CREATE INDEX idx_state ON "order"(state);
