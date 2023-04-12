from sqlalchemy import Column, Integer, String, DateTime, ForeignKey, create_engine, UniqueConstraint
from sqlalchemy.orm import declarative_base, relationship
from sqlalchemy.sql import func

Base = declarative_base()

# Update the User class
class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True)
    email = Column(String, nullable=False, unique=True)
    hashed_password = Column(String, nullable=False)  # Rename this field
    jwt_token = Column(String, nullable=False)
    created_at = Column(DateTime(timezone=True), default=func.now())

    # Add a unique constraint for the email field
    __table_args__ = (UniqueConstraint("email"),)


    # Add a unique constraint for the email field
    __table_args__ = (UniqueConstraint("email"),)

class BotStatus(Base):
    __tablename__ = "bot_statuses"

    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, ForeignKey("users.id"))
    status = Column(String, nullable=False)

    user = relationship("User", back_populates="bot_status")

User.bot_status = relationship("BotStatus", uselist=False, back_populates="user")

class PaymentProof(Base):
    __tablename__ = "payment_proofs"

    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, ForeignKey("users.id"))
    transaction_hash = Column(String, nullable=False)
    created_at = Column(DateTime(timezone=True), default=func.now())

    user = relationship("User", back_populates="payment_proofs")

User.payment_proofs = relationship("PaymentProof", back_populates="user")
