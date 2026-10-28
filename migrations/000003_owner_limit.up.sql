ALTER TABLE trackerapp.spending_limit
    ADD COLUMN user_id INTEGER NOT NULL REFERENCES trackerapp.users(id) ON DELETE CASCADE;

ALTER TABLE trackerapp.transactions
    ALTER COLUMN date TYPE TIMESTAMPTZ USING date AT TIME ZONE 'UTC';