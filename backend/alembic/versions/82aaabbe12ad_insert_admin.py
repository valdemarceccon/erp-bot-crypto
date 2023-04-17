"""Insert admin

Revision ID: 82aaabbe12ad
Revises: c1915f90ffa7
Create Date: 2023-04-17 15:00:56.753669

"""
import sqlalchemy as sa
from alembic import op


# revision identifiers, used by Alembic.
revision = "82aaabbe12ad"
down_revision = "c1915f90ffa7"
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.execute("INSERT INTO users_role (name) VALUES ('admin')")


def downgrade() -> None:
    op.execute("DELETE FROM users_role WHERE name = 'admin'")
