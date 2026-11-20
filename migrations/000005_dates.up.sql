ALTER TABLE trackerapp.transactions
    ALTER COLUMN date TYPE TIMESTAMPTZ USING date AT TIME ZONE 'Europe/Moscow';

CREATE TABLE trackerapp.outbox(
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER NOT NULL UNIQUE ,
    status VARCHAR(6) NOT NULL CHECK (
        status = 'unsent' OR status = 'sent'
        ) DEFAULT 'unsent',

    FOREIGN KEY (transaction_id) REFERENCES trackerapp.transactions ON DELETE CASCADE
);