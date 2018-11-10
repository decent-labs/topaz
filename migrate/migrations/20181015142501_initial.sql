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
    "interval" INTEGER NOT NULL,
    name CHARACTER varying(255) NOT NULL,
    last_flushed TIMESTAMP,
    user_id INTEGER REFERENCES users(id),
    eth_address CHARACTER varying(255) NOT NULL
);

CREATE TABLE flushes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    app_id INT REFERENCES apps(id),
    directory_hash CHARACTER varying(255) NOT NULL,
    "eth_transaction" CHARACTER varying(255) NOT NULL
);

CREATE TABLE objects (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    hash CHARACTER varying(255) NOT NULL,
    app_id INT REFERENCES apps(id),
    flush_id INT
);

-- +migrate Down

DROP TABLE objects, flushes, apps, users;
