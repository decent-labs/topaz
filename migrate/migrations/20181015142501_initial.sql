-- +migrate Up

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    name CHARACTER varying(255) NOT NULL
);

CREATE TABLE apps (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    "interval" int NOT NULL,
    name CHARACTER varying(255) NOT NULL,
    last_flushed TIMESTAMP
    user_id INT NOT NULL,
);

CREATE TABLE flushes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    app_id INT NOT NULL,
    directory_hash CHARACTER varying(255) NOT NULL,
    "transaction" CHARACTER varying(255) NOT NULL
);

CREATE TABLE objects (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    hash CHARACTER varying(255) NOT NULL,
    app_id INT NOT NULL,
    flush_id INT
);

-- +migrate Down

DROP TABLE objects, flushes, apps, users;
