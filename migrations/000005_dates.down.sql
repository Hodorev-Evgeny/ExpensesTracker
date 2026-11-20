ALTER TABLE trackerapp.transactions
    ALTER COLUMN date TYPE TIMESTAMP WITHOUT TIME ZONE USING date AT TIME ZONE 'Europe/Moscow';

DROP TABLE IF EXISTS trackerapp.outbox;