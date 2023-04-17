"""Insert admin

Revision ID: 6c63057e83e1
Revises: 3b381995c5b5
Create Date: 2023-04-17 16:44:25.179388

"""
import sqlalchemy as sa
from alembic import op


# revision identifiers, used by Alembic.
revision = "6c63057e83e1"
down_revision = "3b381995c5b5"
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.execute("INSERT INTO users_role (name) VALUES ('admin')")


def downgrade() -> None:
    op.execute("DELETE FROM users_role WHERE name = 'admin'")
