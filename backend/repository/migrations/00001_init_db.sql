-- +goose Up
-- +goose StatementBegin
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    UNIQUE (name)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    hashed_password VARCHAR(255) NOT NULL
);

CREATE TABLE api_key (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    exchange VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    secret VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL
);

CREATE INDEX ix_api_key_user_id ON api_key(user_id);

CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(id),
    permission_id INTEGER NOT NULL REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE user_roles (
    user_id INTEGER NOT NULL REFERENCES users(id),
    role_id INTEGER NOT NULL REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE closed_pnl (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    api_key_id INTEGER NOT NULL REFERENCES api_key(id),
    symbol VARCHAR(100) NOT NULL,
    orderId VARCHAR(100) NOT NULL,
    side VARCHAR(100) NOT NULL,
    qty VARCHAR(100) NOT NULL,
    orderPrice VARCHAR(100) NOT NULL,
    orderType VARCHAR(100) NOT NULL,
    execType VARCHAR(100) NOT NULL,
    closedSize VARCHAR(100) NOT NULL,
    cumEntryValue VARCHAR(100) NOT NULL,
    avgEntryPrice VARCHAR(100) NOT NULL,
    cumExitValue VARCHAR(100) NOT NULL,
    avgExitPrice VARCHAR(100) NOT NULL,
    closedPnl VARCHAR(100) NOT NULL,
    fillCount VARCHAR(100) NOT NULL,
    leverage VARCHAR(100) NOT NULL,
    createdTime VARCHAR(100) NOT NULL,
    updatedTime VARCHAR(100) NOT NULL
);

CREATE INDEX ix_closed_pnl_api_key_id ON closed_pnl(api_key_id);
CREATE INDEX ix_closed_pnl_user_id ON closed_pnl(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS closed_pnl;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS api_key;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;
-- +goose StatementEnd
