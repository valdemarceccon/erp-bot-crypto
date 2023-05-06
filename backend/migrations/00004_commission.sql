-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table comission;
drop table bot_stop;
drop table bot_start;
-- +goose StatementEnd
