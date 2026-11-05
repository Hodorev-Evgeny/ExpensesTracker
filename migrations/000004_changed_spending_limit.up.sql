UPDATE trackerapp.spending_limit AS l
SET user_id = c.user_id
    FROM trackerapp.categories AS c
WHERE c.limit_id = l.id;

DELETE FROM trackerapp.spending_limit
WHERE user_id IS NULL;

ALTER TABLE trackerapp.spending_limit
    ALTER COLUMN user_id SET NOT NULL;