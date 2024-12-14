from pathlib import Path
from sys import stdout
from typing import TYPE_CHECKING

from loguru import logger as _logger

from src.config import LOG_LEVEL

LOG_FORMAT = (
    "\033[0m<green>{time:YYYY-MM-DD HH:mm:ss.SSS}</green> | <level>{level}</level> | <cyan>[{function}]</cyan>: <level>{message}</level>\033[0m"
)

if TYPE_CHECKING:
    from loguru import Logger


def init_logger() -> "Logger":
    log_root = Path("/data")  # type:ignore

    info = "info.log"
    success = "success.log"
    error = "error.log"

    log_format = "<green>{time:YYYY-MM-DD HH:mm:ss.SSS}</green> | <level>{level}</level> | <cyan>[{function}]</cyan>: <level>{message}</level>"

    if not log_root.exists():
        log_root.mkdir(parents=True)
    if not log_root.is_dir():
        raise ValueError("LOG_ROOT is not a directory")

    _logger.remove()  # remove origin handler
    _logger.add(stdout, colorize=True, enqueue=True, level=LOG_LEVEL, format=log_format)
    _logger.add(log_root.joinpath(info), encoding="utf-8", rotation="100MB", enqueue=True, level="INFO", format=log_format)
    _logger.add(log_root.joinpath(success), encoding="utf-8", rotation="100MB", enqueue=True, level="SUCCESS", format=log_format)
    _logger.add(log_root.joinpath(error), encoding="utf-8", rotation="100MB", enqueue=True, level="ERROR", format=log_format)

    _logger.info("[Init Logger] Init logger successfully")

    return _logger


logger = init_logger()
