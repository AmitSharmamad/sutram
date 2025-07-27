-- Options Trading Leg Table
-- Captures all details for individual legs in multi-leg options strategies

CREATE TABLE option_leg (
    id UUID DEFAULT gen_random_uuid(),

    -- Strategy/Trade identification
    strategy_id UUID NOT NULL,
    leg_sequence INTEGER NOT NULL, -- Order of legs in the strategy (1, 2, 3, etc.)

    -- Index/Underlying details
    select_index VARCHAR(50) NOT NULL, -- e.g., 'Nifty', 'BankNifty', 'FinNifty'
    segment segment NOT NULL, -- Using existing segment enum

    -- Option type (only applicable when segment is 'F&O')
    option_type option_type, -- Using proper enum for option type

    -- Action type
    action_type transaction_type NOT NULL, -- Using existing transaction_type enum

    -- Strike price details
    strike_type strike_type NOT NULL,
    strike_price DECIMAL(18,8), -- Actual strike price (populated based on strike_price_type)

    -- Quantity
    total_lot INTEGER NOT NULL,
    lot_size INTEGER, -- Lot size for the instrument

    -- Order details
    order_type order_type NOT NULL, -- Order type for the leg
    tick_offset DECIMAL(18,8), -- Tick size for the instrument

    -- Expiry details
    expiry_type expiry_type NOT NULL,
    expiry_cycle expiry_cycle NOT NULL,
    expiry_date DATE,

    -- Exchange and instrument details
    exchange exchange NOT NULL,
    asset asset NOT NULL,
    exchange_security_id INTEGER,
    instrument_type instrument_type NOT NULL,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    FOREIGN KEY (strategy_id) REFERENCES strategy(id) ON DELETE CASCADE,

    -- Additional constraints
    CONSTRAINT valid_option_details CHECK (
        (segment = 'F&O' AND option_type IS NOT NULL) OR
        (segment = 'EQ' AND option_type IS NULL)
    ),

    CONSTRAINT valid_strike_price CHECK (
        (segment = 'F&O' AND strike_price IS NOT NULL) OR
        (segment = 'EQ' AND strike_price IS NULL)
    )
);

-- Indexes for better performance
CREATE INDEX idx_option_leg_strategy_id ON option_leg(strategy_id);
CREATE INDEX idx_option_leg_select_index ON option_leg(select_index);
CREATE INDEX idx_option_leg_expiry_date ON option_leg(expiry_date);

-- Trigger to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_option_leg_updated_at
    BEFORE UPDATE ON option_leg
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
