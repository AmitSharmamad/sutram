-- Remove initial data
DELETE FROM user_usage_limits;
DELETE FROM oauth_provider WHERE name IN ('google', 'github');
