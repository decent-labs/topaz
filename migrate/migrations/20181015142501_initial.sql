-- +migrate Up

CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name character varying(255),
    "interval" interval,
    created_at timestamp without time zone DEFAULT now(),
    flushed_at timestamp without time zone
);

CREATE TABLE flushes (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid,
    hash character varying(255),
    created_at timestamp without time zone DEFAULT now()
);

CREATE TABLE objects (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid,
    flush_id uuid,
    hash character varying(255),
    created_at timestamp without time zone DEFAULT now()
);

ALTER TABLE ONLY users ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY flushes ADD CONSTRAINT flushes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY flushes ADD CONSTRAINT flushes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE ONLY objects ADD CONSTRAINT objects_pkey PRIMARY KEY (id);
ALTER TABLE ONLY objects ADD CONSTRAINT objects_flush_id_fkey FOREIGN KEY (flush_id) REFERENCES flushes(id);
ALTER TABLE ONLY objects ADD CONSTRAINT objects_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);

-- +migrate Down

DROP TABLE objects, flushes, users;
