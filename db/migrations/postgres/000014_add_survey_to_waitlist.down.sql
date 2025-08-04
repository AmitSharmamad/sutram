-- Remove survey reference from waitlist
DROP INDEX IF EXISTS idx_waitlist_survey_id;
ALTER TABLE waitlist DROP COLUMN IF EXISTS survey_id;

-- Drop survey table and related indexes
DROP INDEX IF EXISTS idx_survey_experience_level;
DROP INDEX IF EXISTS idx_survey_trading_type;
DROP INDEX IF EXISTS idx_survey_nature_of_work;
DROP INDEX IF EXISTS idx_survey_created_at;
DROP INDEX IF EXISTS idx_survey_email;
DROP TABLE IF EXISTS survey;

-- Drop enum types
DROP TYPE IF EXISTS work_type;
DROP TYPE IF EXISTS experience_level;
DROP TYPE IF EXISTS trading_type;
DROP TYPE IF EXISTS nature_of_work;
