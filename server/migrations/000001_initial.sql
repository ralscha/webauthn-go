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
    user_id          INTEGER      NOT NULL,
    webauthn_user_id BYTEA UNIQUE NOT NULL,
    created_at       TIMESTAMP    NOT NULL DEFAULT now(),
    last_used        TIMESTAMP    NULL,
    aaguid           BYTEA     NULL,
    attestation_type VARCHAR(32),
    attachment       VARCHAR(64)  NOT NULL,
    transport        VARCHAR(64)  NOT NULL DEFAULT '',
    sign_count       INTEGER      NOT NULL DEFAULT 0,
    clone_warning    BOOLEAN      NOT NULL DEFAULT FALSE,
    present          BOOLEAN      NOT NULL DEFAULT FALSE,
    verified         BOOLEAN      NOT NULL DEFAULT FALSE,
    backup_eligible  BOOLEAN      NOT NULL DEFAULT FALSE,
    backup_state     BOOLEAN      NOT NULL DEFAULT FALSE,
    public_key       BYTEA        NOT NULL,

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
