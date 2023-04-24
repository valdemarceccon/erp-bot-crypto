"""change api key fields

Revision ID: 92fa26954ea9
Revises: 49b2685f3964
Create Date: 2023-04-24 00:33:29.992117

"""
import sqlalchemy as sa
from alembic import op


# revision identifiers, used by Alembic.
revision = "92fa26954ea9"
down_revision = "49b2685f3964"
branch_labels = None
depends_on = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.add_column("apikeys", sa.Column("secret", sa.String(length=255), nullable=False))
    op.drop_column("apikeys", "api_secret")
    op.drop_column("apikeys", "apikey")
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.add_column(
        "apikeys",
        sa.Column(
            "apikey", sa.VARCHAR(length=255), autoincrement=False, nullable=False
        ),
    )
    op.add_column(
        "apikeys",
        sa.Column(
            "api_secret", sa.VARCHAR(length=255), autoincrement=False, nullable=False
        ),
    )
    op.drop_column("apikeys", "secret")
    # ### end Alembic commands ###