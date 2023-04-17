from passlib.hash import bcrypt
from sqlalchemy.orm import Session
from src.models.user import User
from src.schemas.user import UserCreate


def get_password_hash(password: str) -> str:
    return bcrypt.hash(password)


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return bcrypt.verify(plain_password, hashed_password)


def create_user(db: Session, user: UserCreate) -> User:
    hashed_password = get_password_hash(user.password)
    db_user = User(email=user.email, hashed_password=hashed_password, name=user.name)
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    return db_user


def authenticate_user(db: Session, email: str, password: str) -> User | None:
    user = get_user(db, email)

    if not user:
        return None

    if not verify_password(password, str(user.hashed_password)):
        return None
    return user


def get_user(db: Session, email: str) -> User | None:
    user: User = db.query(User).filter(User.email == email).first()
    if not user:
        return None

    return user
