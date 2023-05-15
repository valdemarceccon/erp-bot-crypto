-- +goose Up
-- +goose StatementBegin
CREATE TABLE permission (
    id SERIAL PRIMARY KEY,
    permission_name VARCHAR(255) UNIQUE NOT NULL,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(255) UNIQUE NOT NULL,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    fullname VARCHAR(255) NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

CREATE INDEX ix_users_username ON users(username);
CREATE INDEX ix_users_email ON users(email);

CREATE TABLE api_key (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    api_key_name VARCHAR(255) NOT NULL,
    exchange VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    api_secret VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

CREATE INDEX ix_api_key_user_id ON api_key(user_id);

CREATE TABLE role_permission (
    role_id INTEGER NOT NULL REFERENCES roles(id),
    permission_id INTEGER NOT NULL REFERENCES permission(id),
    PRIMARY KEY (role_id, permission_id),
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

CREATE TABLE user_roles (
    user_id INTEGER NOT NULL REFERENCES users(id),
    role_id INTEGER NOT NULL REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id),
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

CREATE TABLE closed_pnl (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    api_key_id INTEGER NOT NULL REFERENCES api_key(id),
    symbol VARCHAR(100) NOT NULL,
    orderId VARCHAR(100) NOT NULL,
    execType VARCHAR(100) NOT NULL,
    closedPnl DECIMAL(25,15) NOT NULL,
    createdTime BIGINT NOT NULL,
    updatedTime BIGINT NOT NULL,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

CREATE INDEX ix_closed_pnl_api_key_id ON closed_pnl(api_key_id);
CREATE INDEX ix_closed_pnl_user_id ON closed_pnl(user_id);

CREATE TABLE bot_start (
  id SERIAL,
  api_key_id INTEGER NOT NULL,
  start_time TIMESTAMP NOT NULL,
  wallet_balance NUMERIC(23,8) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (api_key_id) REFERENCES api_key(id)
);

CREATE TABLE bot_stop (
  id SERIAL,
  stop_time TIMESTAMP NOT NULL,
  start_time_id INTEGER NOT NULL,
  wallet_balance NUMERIC(23,8) NOT NULL,
  PRIMARY key (id),
  FOREIGN KEY (start_time_id) REFERENCES bot_start(id)
);

CREATE TABLE comission (
  user_id INTEGER NOT NULL,
  comission_date DATE NOT NULL,
  balance NUMERIC(23,8) NOT NULL,
  high_watermark NUMERIC(23,8) NOT NULL,
  tpnl NUMERIC(23,8) NOT NULL,
  net_profit NUMERIC(23,8) NOT NULL,
  fee NUMERIC(23,8) NOT NULL,
  PRIMARY KEY (user_id, comission_date),
  FOREIGN KEY (user_id) REFERENCES users(id)

);

CREATE TABLE damon_process (
  id SERIAL,
  api_key_id INTEGER,
  start_time TIMESTAMP,
  end_time TIMESTAMP,
  status INTEGER, -- 0 ok 1 not ok
  message TEXT,
  PRIMARY KEY(id),
  FOREIGN KEY(api_key_id) REFERENCES api_key(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS damon_process;
DROP TABLE IF EXISTS comission;
DROP TABLE IF EXISTS bot_stop;
DROP TABLE IF EXISTS bot_start;
DROP TABLE IF EXISTS closed_pnl;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permission;
DROP TABLE IF EXISTS api_key;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permission;
-- +goose StatementEnd
