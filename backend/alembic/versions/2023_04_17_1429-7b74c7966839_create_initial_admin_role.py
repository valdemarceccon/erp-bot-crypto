"""Create initial admin role

Revision ID: 7b74c7966839
Revises: 6fefa0b968eb
Create Date: 2023-04-17 14:29:37.414844

"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy import text


# revision identifiers, used by Alembic.
revision = "7b74c7966839"
down_revision = "6fefa0b968eb"
branch_labels = None
depends_on = None


def upgrade():
    # ...
    op.execute("INSERT INTO roles (name) VALUES ('admin');")


def downgrade():
    # ...

    op.execute(
        "DELETE FROM user_roles WHERE role_id = (SELECT role_id FROM roles where name = 'admin')"
    )
    op.execute("DELETE FROM roles WHERE name = 'admin';")
