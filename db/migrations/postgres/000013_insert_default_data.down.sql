-- Remove default system user and trading strategies

-- Delete default strategies (created by system user)
DELETE FROM strategy WHERE user_id = 'e65c67f1-2b43-4581-b33b-ec26295e741b';

-- Delete system user
DELETE FROM "user" WHERE id = 'e65c67f1-2b43-4581-b33b-ec26295e741b';
