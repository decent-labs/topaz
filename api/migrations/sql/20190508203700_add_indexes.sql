-- +migrate Up

CREATE INDEX api_tokens_user_id ON api_tokens(user_id);
CREATE INDEX apps_user_id ON apps(user_id);
CREATE INDEX apps_last_proofed ON apps(last_proofed);
CREATE INDEX apps_interval ON apps(interval);
CREATE INDEX objects_app_id ON objects(app_id);
CREATE INDEX hashes_object_id ON hashes(object_id);
CREATE INDEX hashes_proof_id ON hashes(proof_id);
CREATE INDEX proofs_app_id ON proofs(app_id);
CREATE INDEX blockchain_transactions_proof_id ON blockchain_transactions(proof_id);
CREATE INDEX blockchain_transactions_blockchain_network_id ON blockchain_transactions(blockchain_network_id);

-- +migrate Down

DROP INDEX api_tokens_user_id;
DROP INDEX apps_user_id;
DROP INDEX apps_last_proofed;
DROP INDEX apps_interval;
DROP INDEX objects_app_id;
DROP INDEX hashes_object_id;
DROP INDEX hashes_proof_id;
DROP INDEX proofs_app_id;
DROP INDEX blockchain_transactions_proof_id;
DROP INDEX blockchain_transactions_blockchain_network_id;