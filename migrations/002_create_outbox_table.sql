-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

       CREATE TABLE outbox_events(
           id UUID primary key default gen_random_uuid(),
           aggregate_type text not null,
           aggregate_id text not null,
           event_type varchar(50) not null,
           payload jsonb not null,
           created_at timestamp DEFAULT CURRENT_TIMESTAMP,
           processed boolean default FALSE
       );

-- +goose Down
drop table outbox_events;