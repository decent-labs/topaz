-- +migrate Up

ALTER TABLE users ALTER COLUMN name DROP NOT NULL;

-- +migrate Down

ALTER TABLE users ALTER COLUMN name SET NOT NULL;