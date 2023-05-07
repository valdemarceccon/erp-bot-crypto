-- +goose Up
-- +goose StatementBegin
INSERT INTO permission (permission_name, created_at, updated_at) VALUES ('RunDataCollector', now(), now());
INSERT INTO role_permission (role_id, permission_id, created_at, updated_at) SELECT r.id, p.id, now(), now() FROM permission p, roles r WHERE role_name = 'admin' and permission_name = 'RunDataCollector';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permission WHERE role_id = (SELECT id FROM roles WHERE role_name = 'admin' and deleted_at is null);
DELETE FROM permission where permission_name = 'RunDataCollector';
-- +goose StatementEnd
