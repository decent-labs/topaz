-- +migrate Up

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    name CHARACTER varying(255) NOT NULL,
    email CHARACTER varying(255) UNIQUE NOT NULL,
    password CHARACTER varying(255) NOT NULL
);

CREATE TABLE apps (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    interval INTEGER NOT NULL,
    name CHARACTER varying(255) NOT NULL,
    last_batched INTEGER,
    user_id INTEGER REFERENCES users(id),
    eth_address CHARACTER varying(255) NOT NULL
);

CREATE TABLE batches (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    app_id INTEGER REFERENCES apps(id),
    unix_timestamp INTEGER NOT NULL
);

CREATE TABLE proofs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    batch_id INTEGER REFERENCES batches(id),
    merkle_root CHARACTER varying(255) NOT NULL,
    eth_transaction CHARACTER varying(255) NOT NULL
);

CREATE TABLE objects (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    data_blob BYTEA NOT NULL,
    hash CHARACTER varying(255) NOT NULL,
    unix_timestamp INTEGER NOT NULL,
    app_id INTEGER REFERENCES apps(id),
    proof_id INTEGER
);

-- +migrate Down

DROP TABLE objects, batches, apps, users, proofs;
