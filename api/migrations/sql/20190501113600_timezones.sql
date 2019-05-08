-- +migrate Up

ALTER TABLE users ALTER COLUMN created_at type timestamp with time zone;
ALTER TABLE users ALTER COLUMN updated_at type timestamp with time zone;
ALTER TABLE users ALTER COLUMN deleted_at type timestamp with time zone;

ALTER TABLE api_tokens ALTER COLUMN created_at type timestamp with time zone;
ALTER TABLE api_tokens ALTER COLUMN updated_at type timestamp with time zone;
ALTER TABLE api_tokens ALTER COLUMN deleted_at type timestamp with time zone;

ALTER TABLE apps ALTER COLUMN created_at type timestamp with time zone;
ALTER TABLE apps ALTER COLUMN updated_at type timestamp with time zone;
ALTER TABLE apps ALTER COLUMN deleted_at type timestamp with time zone;

ALTER TABLE objects ALTER COLUMN created_at type timestamp with time zone;
ALTER TABLE objects ALTER COLUMN updated_at type timestamp with time zone;
ALTER TABLE objects ALTER COLUMN deleted_at type timestamp with time zone;

ALTER TABLE hashes ALTER COLUMN created_at type timestamp with time zone;
ALTER TABLE hashes ALTER COLUMN updated_at type timestamp with time zone;
ALTER TABLE hashes ALTER COLUMN deleted_at type timestamp with time zone;

ALTER TABLE proofs ALTER COLUMN created_at type timestamp with time zone;
ALTER TABLE proofs ALTER COLUMN updated_at type timestamp with time zone;
ALTER TABLE proofs ALTER COLUMN deleted_at type timestamp with time zone;

-- +migrate Down

ALTER TABLE users ALTER COLUMN created_at type timestamp;
ALTER TABLE users ALTER COLUMN updated_at type timestamp;
ALTER TABLE users ALTER COLUMN deleted_at type timestamp;

ALTER TABLE api_tokens ALTER COLUMN created_at type timestamp;
ALTER TABLE api_tokens ALTER COLUMN updated_at type timestamp;
ALTER TABLE api_tokens ALTER COLUMN deleted_at type timestamp;

ALTER TABLE apps ALTER COLUMN created_at type timestamp;
ALTER TABLE apps ALTER COLUMN updated_at type timestamp;
ALTER TABLE apps ALTER COLUMN deleted_at type timestamp;

ALTER TABLE objects ALTER COLUMN created_at type timestamp;
ALTER TABLE objects ALTER COLUMN updated_at type timestamp;
ALTER TABLE objects ALTER COLUMN deleted_at type timestamp;

ALTER TABLE hashes ALTER COLUMN created_at type timestamp;
ALTER TABLE hashes ALTER COLUMN updated_at type timestamp;
ALTER TABLE hashes ALTER COLUMN deleted_at type timestamp;

ALTER TABLE proofs ALTER COLUMN created_at type timestamp;
ALTER TABLE proofs ALTER COLUMN updated_at type timestamp;
ALTER TABLE proofs ALTER COLUMN deleted_at type timestamp;
