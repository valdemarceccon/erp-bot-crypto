"""Create initial admin permission in the admin role

Revision ID: 8c13cfbad613
Revises: 7b74c7966839
Create Date: 2023-04-17 14:36:29.309071

"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy import text


# revision identifiers, used by Alembic.
revision = "8c13cfbad613"
down_revision = "7b74c7966839"
branch_labels = None
depends_on = None


def upgrade():
    ...
    # connection = op.get_bind()
    # connection.execute(text("INSERT INTO permissions (name) VALUES ('admin');"))
    # connection.execute(
    #     text(
    #         "INSERT INTO role_permissions (role_id, permission_id) SELECT roles.id, permissions.id FROM roles, permissions WHERE roles.name = 'admin' AND permissions.name = 'admin';"
    #     )
    # )


def downgrade():
    ...
    # connection = op.get_bind()
    # connection.execute(
    #     text(
    #         "DELETE FROM role_permissions USING roles, permissions WHERE roles.id = role_permissions.role_id AND permissions.id = role_permissions.permission_id AND roles.name = 'admin' AND permissions.name = 'admin';"
    #     )
    # )
    # connection.execute(text("DELETE FROM permissions WHERE name = 'admin';"))
