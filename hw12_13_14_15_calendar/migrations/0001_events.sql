-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE "events" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "title" TEXT NOT NULL,
    "description" TEXT,
    "begin" TIMESTAMP NOT NULL,
    "end" TIMESTAMP NOT NULL,
    "user_id" INT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE "events";
