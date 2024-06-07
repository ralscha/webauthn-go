-- +goose Up
CREATE TABLE users
(
    id                 SERIAL PRIMARY KEY,
    username           VARCHAR   NOT NULL,
    registration_start TIMESTAMP NULL,
    created_at         TIMESTAMP DEFAULT now()
);

CREATE TABLE credentials
(
    cred_id          BYTEA PRIMARY KEY,
    cred_public_key  BYTEA        NOT NULL,
    user_id          INTEGER      NOT NULL,
    webauthn_user_id BYTEA UNIQUE NOT NULL,
    counter          INTEGER      NOT NULL,
    created_at       TIMESTAMP DEFAULT now(),
    last_used        TIMESTAMP,

    UNIQUE (webauthn_user_id, user_id),
    CONSTRAINT fk_internal_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);


CREATE TABLE sessions
(
    token  TEXT PRIMARY KEY,
    data   BYTEA       NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);


-- +goose Down
DROP TABLE sessions;
DROP TABLE credentials;
DROP TABLE users;
