-- Options Trading Leg Table
-- Captures all details for individual legs in multi-leg options strategies

CREATE TABLE option_legs (
    leg_id UUID DEFAULT gen_random_uuid(),

    -- Strategy/Trade identification
    strategy_id UUID NOT NULL,
    leg_sequence INTEGER NOT NULL, -- Order of legs in the strategy (1, 2, 3, etc.)

    -- Index/Underlying details
    select_index VARCHAR(50) NOT NULL, -- e.g., 'Nifty', 'BankNifty', 'FinNifty'
    segment segment NOT NULL DEFAULT 'F&O', -- Using existing segment enum

    -- Option type (only applicable when segment is 'F&O')
    option_type option_type_enum, -- Using proper enum for option type

    -- Action type
    action_type transaction_type NOT NULL, -- Using existing transaction_type enum

    -- Strike price details
    strike_price_type strike_price_type NOT NULL,
    strike_price DECIMAL(18,8), -- Actual strike price (populated based on strike_price_type)
    strike_offset INTEGER, -- Offset from ATM (e.g., +100, -50 for OTM/ITM)

    -- Quantity
    total_lot INTEGER NOT NULL DEFAULT 1,
    lot_size INTEGER, -- Lot size for the instrument
    total_quantity INTEGER, -- total_lot * lot_size

    -- Order details
    order_type order_type NOT NULL DEFAULT 'MARKET', -- Order type for the leg
    tick_offset DECIMAL(18,8), -- Tick size for the instrument

    -- Expiry details
    expiry_type expiry_type NOT NULL,
    expiry_date DATE,
    days_to_expiry INTEGER,

    -- Status and timestamps
    leg_status leg_status DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (leg_id),
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
CREATE INDEX idx_option_legs_strategy_id ON option_legs(strategy_id);
CREATE INDEX idx_option_legs_select_index ON option_legs(select_index);
CREATE INDEX idx_option_legs_expiry_date ON option_legs(expiry_date);
CREATE INDEX idx_option_legs_leg_status ON option_legs(leg_status);

-- Trigger to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_option_legs_updated_at
    BEFORE UPDATE ON option_legs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
