-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE events (
    id bigserial NOT NULL PRIMARY KEY,
    title text,
    description text,
    event_time timestamptz(0),
    user_id text,
    created_at timestamptz(0) NOT NULL DEFAULT NOW(),
    updated_at timestamptz(0)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE events;
