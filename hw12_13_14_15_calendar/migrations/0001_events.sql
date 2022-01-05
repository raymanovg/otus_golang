-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE "events" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "title" TEXT NOT NULL,
    "desc" TEXT,
    "begin" TIMESTAMP NOT NULL,
    "end" TIMESTAMP NOT NULL,
    "userID" INT NOT NULL,
    "createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE "events";
