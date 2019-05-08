-- +migrate Up

CREATE TABLE blockchain_networks (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE,
  name CHARACTER varying(255) NOT NULL
);

CREATE TABLE blockchain_transactions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE,
  proof_id uuid REFERENCES proofs(id) NOT NULL,
  blockchain_network_id uuid REFERENCES blockchain_networks(id) NOT NULL,
  transaction_hash CHARACTER varying(255) NOT NULL
);

INSERT INTO blockchain_networks (created_at, updated_at, name)
VALUES (now(), now(), 'ethereum goerli');

INSERT INTO blockchain_transactions (created_at, updated_at, proof_id, blockchain_network_id, transaction_hash)
SELECT now(), now(), p.id, bn.id, p.eth_transaction
FROM proofs p
  JOIN blockchain_networks bn ON bn.name = 'ethereum goerli';

ALTER TABLE proofs DROP COLUMN eth_transaction;

-- +migrate Down

DROP TABLE blockchain_networks, blockchain_transactions;
