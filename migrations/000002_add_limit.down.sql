ALTER TABLE trackerapp.categories
DROP CONSTRAINT IF EXISTS fk_transactions_spending_limit;

ALTER TABLE trackerapp.categories
DROP COLUMN IF EXISTS limit_id;

DROP TABLE IF EXISTS trackerapp.spending_limit;