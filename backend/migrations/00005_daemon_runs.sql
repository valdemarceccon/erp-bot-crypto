-- +goose Up
-- +goose StatementBegin
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
DROP TABLE damon_process;
-- +goose StatementEnd
