"""populate_roles

Revision ID: 3217b6975fb5
Revises: cf90aaa8b2b7
Create Date: 2023-04-23 16:58:24.428680

"""
import sqlalchemy as sa
from alembic import op


# revision identifiers, used by Alembic.
revision = "3217b6975fb5"
down_revision = "cf90aaa8b2b7"
branch_labels = None
depends_on = None


from alembic import op
import sqlalchemy as sa
from sqlalchemy.orm import sessionmaker
from src.models.roles import PermissionEnum
from src.models.user import Permission, UserRole, RolePermission, Role

Session = sessionmaker()


# This migration adds pre-defined roles to the roles table
def upgrade():
    session = Session(bind=op.get_bind())

    # Create a set of existing role names
    existing_roles = set([r.name for r in session.query(Role)])

    # Loop over each value in the PermissionEnum
    for permission in PermissionEnum:
        role_name = permission.value

        # Only add the role if it doesn't already exist in the table
        if role_name not in existing_roles:
            role = Role(name=role_name)
            session.add(role)

    session.commit()


# This migration removes the pre-defined roles from the roles table
def downgrade():
    session = Session(bind=op.get_bind())

    # Loop over each value in the PermissionEnum
    for permission in PermissionEnum:
        role_name = permission.value

        # Delete the role if it exists in the table
        role = session.query(Role).filter_by(name=role_name).first()
        if role is not None:
            session.delete(role)

    session.commit()
