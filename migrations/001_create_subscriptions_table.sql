-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE subscriptions (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               service_name VARCHAR(100) NOT NULL,
                               price INTEGER NOT NULL,
                               user_id VARCHAR(100) NOT NULL,
                               start_date VARCHAR(7) NOT NULL
);

-- +goose Down
DROP TABLE subscriptions;
DROP EXTENSION IF EXISTS "pgcrypto";