from datetime import datetime

from sqlalchemy import BigInteger, Column, DateTime, String
from sqlalchemy.orm import relationship

from src.resource.model.base import Base


class User(Base):
    __tablename__ = "users"

    created_at = Column(DateTime, default=datetime.now)
    updated_at = Column(DateTime, default=datetime.now, onupdate=datetime.now)
    deleted_at = Column(DateTime)
    id = Column(BigInteger, primary_key=True)
    name = Column(String(255), nullable=False)
    email = Column(String(255), unique=True)
    avatar = Column(String)
    points = Column(BigInteger)
    tof4_bind_id = Column(String(255), nullable=True, unique=True)
    uin_bind_id = Column(String(255), nullable=True, unique=True)
    credentials = relationship("Credential", back_populates="User")
    user_groups = relationship("UserGroup", secondary="user_groups_users", back_populates="users")
