from datetime import datetime

from sqlalchemy import BigInteger, Column, DateTime, String

from src.resource.model.base import Base


class User(Base):
    __tablename__ = "users"

    created_at = Column(DateTime, default=datetime.now)
    updated_at = Column(DateTime, default=datetime.now, onupdate=datetime.now)
    deleted_at = Column(DateTime)
    id = Column(BigInteger, primary_key=True)
    account = Column(String(127), unique=True)
    password = Column(String(255), nullable=False)
    sztu_password = Column(String(63))
    sztu_accountn = Column(String(63))
    cookie = Column(String(255))
