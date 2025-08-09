-- Migration: Drop option_leg table
-- Description: Drops the option_leg table and related objects

-- Drop trigger
DROP TRIGGER IF EXISTS update_option_leg_updated_at ON option_leg;

-- Drop indexes
DROP INDEX IF EXISTS idx_option_leg_expiry_date;
DROP INDEX IF EXISTS idx_option_leg_select_index;
DROP INDEX IF EXISTS idx_option_leg_strategy_id;

-- Drop table
DROP TABLE IF EXISTS option_leg;
