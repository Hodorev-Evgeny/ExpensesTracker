CREATE SCHEMA trackerapp;


CREATE TABLE trackerapp.users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(15) CHECK (
        phone_number ~ '^\+[0-9]+$'
        AND
        char_length(phone_number) BETWEEN 10 AND 15
        ),
    password VARCHAR NOT NULL,
    time_add TIMESTAMPTZ NOT NULL
);

CREATE TABLE trackerapp.categories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(125) NOT NULL,
    user_id INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES trackerapp.users(id) ON DELETE CASCADE,
    CONSTRAINT uq_categories_user_title UNIQUE (user_id, title)
);

CREATE TABLE trackerapp.transactions (
    id SERIAL PRIMARY KEY,
    sum INTEGER NOT NULL,
    type_transaction VARCHAR(125) NOT NULL CHECK (
        type_transaction = 'Income' OR type_transaction = 'Expenditure'
        ),
    date TIMESTAMP NOT NULL,
    category_id INTEGER NOT NULL,
    comments VARCHAR(1000),
    user_id INTEGER NOT NULL,
    time_create TIMESTAMPTZ NOT NULL,
    time_changes TIMESTAMPTZ,

    FOREIGN KEY (user_id) REFERENCES trackerapp.users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES trackerapp.categories(id) ON DELETE CASCADE
);