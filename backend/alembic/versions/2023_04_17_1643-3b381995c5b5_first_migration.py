"""First migration

Revision ID: 3b381995c5b5
Revises:
Create Date: 2023-04-17 16:43:47.479613

"""
import sqlalchemy as sa
from alembic import op


# revision identifiers, used by Alembic.
revision = "3b381995c5b5"
down_revision = None
branch_labels = None
depends_on = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.create_table(
        "users_role",
        sa.Column("id", sa.Integer(), autoincrement=True, nullable=False),
        sa.Column("name", sa.String(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(op.f("ix_users_role_id"), "users_role", ["id"], unique=False)
    op.create_index(op.f("ix_users_role_name"), "users_role", ["name"], unique=False)
    op.create_table(
        "users",
        sa.Column("email", sa.String(), nullable=False),
        sa.Column("name", sa.String(), nullable=True),
        sa.Column("hashed_password", sa.String(), nullable=True),
        sa.Column("role_id", sa.Integer(), nullable=True),
        sa.ForeignKeyConstraint(
            ["role_id"],
            ["users_role.id"],
        ),
        sa.PrimaryKeyConstraint("email"),
    )
    op.create_index(op.f("ix_users_email"), "users", ["email"], unique=False)
    op.create_index(op.f("ix_users_name"), "users", ["name"], unique=False)
    op.create_table(
        "apikeys",
        sa.Column("id", sa.Integer(), autoincrement=True, nullable=False),
        sa.Column("user_email", sa.String(), nullable=False),
        sa.Column("apikey", sa.String(), nullable=False),
        sa.ForeignKeyConstraint(
            ["user_email"],
            ["users.email"],
        ),
        sa.PrimaryKeyConstraint("id", "user_email"),
    )
    op.create_index(op.f("ix_apikeys_id"), "apikeys", ["id"], unique=False)
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_index(op.f("ix_apikeys_id"), table_name="apikeys")
    op.drop_table("apikeys")
    op.drop_index(op.f("ix_users_name"), table_name="users")
    op.drop_index(op.f("ix_users_email"), table_name="users")
    op.drop_table("users")
    op.drop_index(op.f("ix_users_role_name"), table_name="users_role")
    op.drop_index(op.f("ix_users_role_id"), table_name="users_role")
    op.drop_table("users_role")
    # ### end Alembic commands ###