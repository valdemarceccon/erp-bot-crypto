-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (role_name, created_at, updated_at) VALUES ('admin', now(), now()) returning id;
INSERT INTO permission (permission_name, created_at, updated_at) VALUES ('ListUsers', now(), now()), ('WriteApiKeys', now(), now()), ('ReadApiKeys', now(), now());
INSERT INTO role_permission (role_id, permission_id, created_at, updated_at) SELECT r.id, p.id, now(), now() FROM permission p, roles r WHERE role_name = 'admin';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permission WHERE role_id = (SELECT id FROM roles WHERE role_name = 'admin');
DELETE FROM permission;
DELETE FROM user_roles where exists (select * from roles where role_name = 'admin');
DELETE FROM roles;

-- +goose StatementEnd
