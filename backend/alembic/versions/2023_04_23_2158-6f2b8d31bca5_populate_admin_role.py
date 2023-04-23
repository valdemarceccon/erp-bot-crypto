"""populate admin role

Revision ID: 6f2b8d31bca5
Revises: 1038a80f8498
Create Date: 2023-04-23 21:58:01.714615

"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.orm import sessionmaker
from src.models.roles import PermissionEnum
from src.models.user import Permission
from src.models.user import Role
from src.models.user import RolePermission
from src.models.user import UserRole


Session = sessionmaker()


# revision identifiers, used by Alembic.
revision = "6f2b8d31bca5"
down_revision = "1038a80f8498"
branch_labels = None
depends_on = None


def upgrade() -> None:
    session = Session(bind=op.get_bind())

    # Create a set of existing role names
    permissions = session.query(Permission).all()
    admin_role = Role(name="Admin")
    admin_role.permissions = permissions
    session.add(admin_role)
    session.commit()


def downgrade() -> None:
    session = Session(bind=op.get_bind())

    admin_role = session.query(Role).filter(Role.name == "Admin").first()
    session.delete(admin_role)
    session.commit()
