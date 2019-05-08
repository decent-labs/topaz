-- +migrate Up

CREATE TABLE api_tokens (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    user_id uuid REFERENCES users(id) NOT NULL,
    token CHARACTER varying(255) UNIQUE NOT NULL
);

-- +migrate Down

DROP TABLE api_tokens;
