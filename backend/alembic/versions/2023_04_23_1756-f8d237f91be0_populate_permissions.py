"""populate permissions

Revision ID: f8d237f91be0
Revises: 34d120bbc603
Create Date: 2023-04-23 17:56:14.397959

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
revision = "f8d237f91be0"
down_revision = "34d120bbc603"
branch_labels = None
depends_on = None


# This migration adds pre-defined roles to the roles table
def upgrade():
    session = Session(bind=op.get_bind())

    # Create a set of existing role names
    existing_roles = set([r.name for r in session.query(Permission)])

    # Loop over each value in the PermissionEnum
    for permission in PermissionEnum:
        permission_name = permission.value

        # Only add the role if it doesn't already exist in the table
        if permission_name not in existing_roles:
            new_permission = Permission(name=permission_name)
            session.add(new_permission)

    session.commit()


# This migration removes the pre-defined roles from the roles table
def downgrade():
    session = Session(bind=op.get_bind())

    # Loop over each value in the PermissionEnum
    for permission in PermissionEnum:
        permission_name = permission.value

        # Delete the role if it exists in the table
        role = session.query(Permission).filter_by(name=permission_name).first()
        if role is not None:
            session.delete(role)

    session.commit()
