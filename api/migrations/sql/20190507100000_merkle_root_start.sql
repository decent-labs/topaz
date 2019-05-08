-- +migrate Up

ALTER TABLE proofs
ADD COLUMN merkle_root_bytes BYTEA;

-- +migrate Down

ALTER TABLE proofs
DROP COLUMN merkle_root_bytes;
