-- +goose Up
-- +goose StatementBegin
ALTER TABLE closed_pnl
ALTER COLUMN createdTime TYPE BIGINT USING createdTime::BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE closed_pnl
ALTER COLUMN createdTime TYPE VARCHAR(100) USING createdTime::VARCHAR(100);
-- +goose StatementEnd
