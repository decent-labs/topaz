-- +migrate Up

CREATE TABLE blockchain_explorers (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE,
  blockchain_network_id uuid REFERENCES blockchain_networks(id) NOT NULL,
  url_template CHARACTER varying(255) NOT NULL
);

INSERT INTO blockchain_explorers (created_at, updated_at, blockchain_network_id, url_template)
SELECT now(), now(), id, 'https://goerli.etherscan.io/tx/{transaction_hash}'
FROM blockchain_networks bn
WHERE bn.name = 'ethereum goerli';

INSERT INTO blockchain_explorers (created_at, updated_at, blockchain_network_id, url_template)
SELECT now(), now(), id, 'https://blockscout.com/eth/goerli/tx/{transaction_hash}/internal_transactions'
FROM blockchain_networks bn
WHERE bn.name = 'ethereum goerli';

-- +migrate Down

DROP TABLE blockchain_explorers;
