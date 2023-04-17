"""Permissions enum

Revision ID: a425e90fdfa0
Revises: 9454644cac93
Create Date: 2023-04-17 18:37:43.674315

"""
import sqlalchemy as sa
from alembic import op


# revision identifiers, used by Alembic.
revision = "a425e90fdfa0"
down_revision = "9454644cac93"
branch_labels = None
depends_on = None


from alembic import op
import sqlalchemy as sa
from src.models.roles import PermissionEnum


def upgrade():
    op.execute(
        "DELETE FROM role_permissions USING roles, permissions WHERE roles.id = role_permissions.role_id AND permissions.id = role_permissions.permission_id AND roles.name = 'admin' AND permissions.name = 'admin';"
    )
    op.execute("DELETE FROM permissions WHERE name = 'admin';")

    permissions_table = sa.table(
        "permissions",
        sa.column("name", sa.String),
    )

    for permission in PermissionEnum:
        op.execute(
            f"""
            INSERT INTO permissions (name)
            SELECT '{permission.value}'
            WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE name='{permission.value}');
            """
        )

    op.execute(
        f"INSERT INTO role_permissions (role_id, permission_id) SELECT roles.id, permissions.id FROM roles, permissions WHERE roles.name = 'admin' AND permissions.name = '{PermissionEnum.ADMIN.value}';"
    )


def downgrade():
    op.execute(
        f"""
        DELETE FROM role_permissions
        WHERE role_id IN (SELECT id FROM roles WHERE name = 'admin')
        AND permission_id IN (SELECT id FROM permissions WHERE name = '{PermissionEnum.ADMIN.value}');
        """
    )

    for permission in PermissionEnum:
        op.execute(f"DELETE FROM permissions WHERE name='{permission.value}';")

    op.execute("INSERT INTO permissions (name) VALUES ('admin');")
    op.execute(
        "INSERT INTO role_permissions (role_id, permission_id) SELECT roles.id, permissions.id FROM roles, permissions WHERE roles.name = 'admin' AND permissions.name = 'admin';"
    )
