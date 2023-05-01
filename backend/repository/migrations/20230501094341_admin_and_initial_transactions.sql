-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (role_name) VALUES ('admin') returning id;
INSERT INTO permission (permission_name) VALUES ('ListUsers'), ('WriteApiKeys'), ('ReadApiKeys');
INSERT INTO role_permission (role_id, permission_id) SELECT r.id, p.id FROM permission p, roles r WHERE role_name = 'admin';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permission WHERE role_id = (SELECT id FROM roles WHERE role_name = 'admin');
DELETE FROM permission;
DELETE FROM roles;
-- +goose StatementEnd
