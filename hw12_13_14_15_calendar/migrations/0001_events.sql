-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE "events" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "title" TEXT,
    "description" TEXT,
    "begin" TIMESTAMP,
    "end" TIMESTAMP,
    "user_id" INT,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE "events";
