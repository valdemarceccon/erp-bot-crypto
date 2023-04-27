"""Api admin view

Revision ID: 92e862656f74
Revises: f4314becde3c
Create Date: 2023-04-27 17:44:03.746474

"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.orm import sessionmaker
from src.models.roles import PermissionEnum
from src.models.user import Permission
from src.models.user import Role


Session = sessionmaker()


# revision identifiers, used by Alembic.
revision = "92e862656f74"
down_revision = "f4314becde3c"
branch_labels = None
depends_on = None


def upgrade():
    session = Session(bind=op.get_bind())
    admin_role = session.query(Role).filter(Role.name == "Admin").first()

    if not admin_role:
        raise LookupError("No admin role")

    # Create a set of existing role names
    existing_permissions = set([r.name for r in session.query(Permission)])

    # Loop over each value in the PermissionEnum
    for permission in PermissionEnum:
        permission_name = permission.value

        # Only add the role if it doesn't already exist in the table
        if permission_name not in existing_permissions:
            new_permission = Permission(name=permission_name)
            admin_role.permissions.append(new_permission)
            session.add(new_permission)

    session.commit()


# This migration removes the pre-defined roles from the roles table
def downgrade():
    session = Session(bind=op.get_bind())

    admin_role = session.query(Role).filter(Role.name == "Admin").first()

    if not admin_role:
        raise LookupError("No admin role")

    admin_role.permissions.clear()

    # Loop over each value in the PermissionEnum
    for permission in PermissionEnum:
        permission_name = permission.value

        # Delete the role if it exists in the table
        role = session.query(Permission).filter_by(name=permission_name).first()
        if role is not None:
            session.delete(role)

    session.commit()
