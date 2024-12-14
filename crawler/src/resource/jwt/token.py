from typing import Any, Dict

import jwt

from src.config import JWT_SECRET_KEY
from src.log import logger
from src.resource.database.mysql import Session
from src.resource.model.user import User

if not JWT_SECRET_KEY:
    raise ValueError("JWT_SECRET_KEY must be set")


def parse_token(jwt_token: str | bytes) -> Dict[str, Any]:
    try:
        payload = jwt.decode(jwt_token, bytes(
            JWT_SECRET_KEY, encoding="utf-8"), algorithms=["HS256"])
        return payload
    except jwt.ExpiredSignatureError:
        logger.error("[ParseToken] Token has expired")
        raise ValueError("Token has expired")
    except (jwt.InvalidSignatureError, jwt.InvalidTokenError):
        logger.error("[ParseToken] Token is invalid")
        raise ValueError("Token is invalid")
    except Exception as e:
        logger.error("[ParseToken] Token parse error: {e}")
        raise ValueError(f"Token parse error: {e}")


def validate_payload(payload: Dict[str, Any]) -> int:
    user_id = payload.get("user_id")
    if not user_id:
        logger.error("[ValidateToken]")
        raise ValueError("Token is invalid")
    with Session() as session:
        if not session.is_active:
            session.rollback()
            session.close()
        else:
            session.commit()

        query = session.query(User.id).filter(
            User.id == user_id).filter(User.deleted_at.is_(None))
        result = query.first()

    if not result:
        logger.error("[ValidateToken]")
        raise ValueError("User is not exists")

    return user_id
