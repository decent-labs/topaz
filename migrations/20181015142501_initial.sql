-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    name CHARACTER varying(255) NOT NULL,
    email CHARACTER varying(255) UNIQUE NOT NULL,
    password CHARACTER varying(255) NOT NULL
);

CREATE TABLE apps (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    interval INTEGER NOT NULL,
    name CHARACTER varying(255) NOT NULL,
    last_batched INTEGER,
    user_id uuid REFERENCES users(id),
    eth_address CHARACTER varying(255) NOT NULL
);

CREATE TABLE batches (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    app_id uuid REFERENCES apps(id),
    unix_timestamp INTEGER NOT NULL
);

CREATE TABLE proofs (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    batch_id uuid REFERENCES batches(id),
    merkle_root CHARACTER varying(255) NOT NULL,
    eth_transaction CHARACTER varying(255) NOT NULL
);

CREATE TABLE objects (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    app_id uuid REFERENCES apps(id),
    uuid uuid NOT NULL
);

CREATE TABLE hashes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    object_id uuid REFERENCES objects(id),
    proof_id uuid,
    hash BYTEA NOT NULL,
    unix_timestamp INTEGER NOT NULL
);

-- +migrate Down

DROP TABLE objects, batches, apps, users, proofs, hashes;
