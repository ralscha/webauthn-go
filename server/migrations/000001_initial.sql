-- +goose Up
CREATE TABLE app_user
(
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(255) UNIQUE NOT NULL,
    registration_start TIMESTAMPTZ         NULL
);

CREATE TABLE app_credentials
(
    id          BYTEA  NOT NULL,
    app_user_id BIGINT NOT NULL,
    public_key  BYTEA  NOT NULL,
    attestation_type VARCHAR(50) NULL,
    transports VARCHAR(255) NULL,
    aa_guid BYTEA NOT NULL,
    sign_count INTEGER NOT NULL,
    PRIMARY KEY (id, app_user_id),
    FOREIGN KEY (app_user_id) REFERENCES app_user (id) ON DELETE CASCADE
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
DROP TABLE app_credentials;
DROP TABLE app_user;
