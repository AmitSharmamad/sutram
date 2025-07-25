-- Add unique constraint to ensure only one open position per deployment
-- This prevents multiple open positions for the same deployment at the database level

-- First, identify and handle any existing violations
DO $$
DECLARE
    violation_count INT;
BEGIN
    -- Count deployments with multiple open positions
    SELECT COUNT(*)
    INTO violation_count
    FROM (
        SELECT deployment_id, COUNT(*) as open_count
        FROM deployment_position
        WHERE status = 'OPEN'
        GROUP BY deployment_id
        HAVING COUNT(*) > 1
    ) violations;

    -- If violations exist, close older positions keeping only the latest one
    IF violation_count > 0 THEN
        RAISE NOTICE 'Found % deployment(s) with multiple open positions. Closing older positions...', violation_count;
        
        -- Close older open positions, keeping only the most recent one per deployment
        UPDATE deployment_position
        SET 
            status = 'CLOSED',
            exit_reason = 'SYSTEM_CLEANUP',
            exit_time = COALESCE(exit_time, NOW()),
            exit_price = COALESCE(exit_price, entry_price),
            updated_at = NOW()
        WHERE id IN (
            SELECT dp.id
            FROM deployment_position dp
            INNER JOIN (
                SELECT 
                    deployment_id,
                    MAX(entry_time) as latest_entry_time
                FROM deployment_position
                WHERE status = 'OPEN'
                GROUP BY deployment_id
                HAVING COUNT(*) > 1
            ) latest ON dp.deployment_id = latest.deployment_id
            WHERE dp.status = 'OPEN'
            AND dp.entry_time < latest.latest_entry_time
        );
        
        RAISE NOTICE 'Closed older positions. Each deployment now has at most one open position.';
    ELSE
        RAISE NOTICE 'No deployments with multiple open positions found.';
    END IF;
END $$;

-- Create the unique constraint to prevent future violations
-- This ensures only one open position per deployment
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_deployment_position_one_open_per_deployment
ON deployment_position (deployment_id)
WHERE status = 'OPEN';

-- Add a comment to document the constraint
COMMENT ON INDEX idx_deployment_position_one_open_per_deployment IS 
'Ensures only one open position per deployment at any time. Prevents multiple concurrent open positions.';
