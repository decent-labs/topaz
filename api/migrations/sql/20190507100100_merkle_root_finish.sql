-- +migrate Up

ALTER TABLE proofs
DROP COLUMN merkle_root;

ALTER TABLE proofs
RENAME COLUMN merkle_root_bytes TO merkle_root;

ALTER TABLE proofs
ALTER COLUMN merkle_root SET NOT NULL;

-- +migrate Down

ALTER TABLE proofs
ALTER COLUMN merkle_root DROP NOT NULL;

ALTER TABLE proofs
RENAME COLUMN merkle_root TO merkle_root_bytes;

ALTER TABLE proofs
ADD COLUMN merkle_root CHARACTER varying(255) NOT NULL;
