"""Create user admin

Revision ID: 7ce2a5703dad
Revises: dae99ac4f75d
Create Date: 2023-04-17 14:24:21.083407

"""
import sqlalchemy as sa
from alembic import op
from src.models.user import Role
from src.models.user import User


# revision identifiers, used by Alembic.
revision = "7ce2a5703dad"
down_revision = "dae99ac4f75d"
branch_labels = None
depends_on = None


def upgrade() -> None:
    ...


def downgrade() -> None:
    ...
