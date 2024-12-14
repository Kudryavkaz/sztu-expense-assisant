from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from src.config import MYSQL_DATABASE, MYSQL_HOST, MYSQL_PASSWORD, MYSQL_PORT, MYSQL_USER
from src.resource.model.base import Base

sql_link = f"mysql+mysqlconnector://{MYSQL_USER}:{MYSQL_PASSWORD}@{MYSQL_HOST}:{MYSQL_PORT}/{MYSQL_DATABASE}"

engine = create_engine(
    sql_link,
    connect_args={
        "charset": "utf8mb4",
    },
    pool_size=10,
    max_overflow=20,
    pool_pre_ping=True,
    pool_use_lifo=True,
    pool_recycle=3600,
)

Session = sessionmaker(bind=engine)

Base.metadata.create_all(engine)
