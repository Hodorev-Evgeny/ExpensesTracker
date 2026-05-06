CREATE TABLE trackerapp.spending_limit (
    id SERIAL PRIMARY KEY,
    duration TIMESTAMP,
    amount_limit INTEGER
);


ALTER TABLE trackerapp.categories
ADD COLUMN limit_id  INTEGER;

ALTER TABLE trackerapp.categories
ADD CONSTRAINT fk_transactions_spending_limit
FOREIGN KEY (limit_id)
REFERENCES trackerapp.spending_limit(id);
