-- +goose Up
CREATE TABLE app_user
(
    id                 BIGSERIAL PRIMARY KEY,
    username           VARCHAR(255) UNIQUE NOT NULL,
    registration_start TIMESTAMPTZ         NULL
);

CREATE TABLE credentials
(
    id              BIGSERIAL PRIMARY KEY,
    app_user_id     BIGINT NOT NULL,
    count           BIGINT NOT NULL,
    public_key_cose BYTEA  NOT NULL,
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
DROP TABLE credentials;
DROP TABLE app_user;
