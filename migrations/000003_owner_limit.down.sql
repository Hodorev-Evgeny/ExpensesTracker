ALTER TABLE trackerapp.transactions
ALTER COLUMN date TYPE TIMESTAMP WITHOUT TIME ZONE;

ALTER TABLE trackerapp.spending_limit
DROP COLUMN IF EXISTS user_id;