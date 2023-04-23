"""populate permissions

Revision ID: 1038a80f8498
Revises: b75f3890c9b5
Create Date: 2023-04-23 21:56:23.298703

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
revision = "1038a80f8498"
down_revision = "b75f3890c9b5"
branch_labels = None
depends_on = None


def upgrade():
    session = Session(bind=op.get_bind())

    # Create a set of existing role names
    existing_permissions = set([r.name for r in session.query(Permission)])

    # Loop over each value in the PermissionEnum
    for permission in PermissionEnum:
        permission_name = permission.value

        # Only add the role if it doesn't already exist in the table
        if permission_name not in existing_permissions:
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
