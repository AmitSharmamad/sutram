CREATE TABLE IF NOT EXISTS deployment_position_log (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    deployment_position_id UUID NOT NULL REFERENCES deployment_position(id) ON DELETE CASCADE,
    results JSONB,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_deployment_position_log_position_id ON deployment_position_log(deployment_position_id);
CREATE INDEX IF NOT EXISTS idx_deployment_position_log_timestamp ON deployment_position_log(timestamp);
