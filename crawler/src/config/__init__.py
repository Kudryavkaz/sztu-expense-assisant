import os

GRPC_THREAD_POOL_WORKER = int(os.environ.get("GRPC_THREAD_POOL_WORKER", "10"))
GRPC_SERVER_PORT = os.environ.get("GRPC_SERVER_PORT")
JWT_SECRET_KEY = os.environ.get("JWT_SECRET_KEY")

# log
LOG_LEVEL = os.environ.get("LOG_LEVEL", "INFO")

# db
MYSQL_HOST = os.environ.get("MYSQL_HOST", "localhost")
MYSQL_PORT = os.environ.get("MYSQL_PORT", "3306")
MYSQL_USER = os.environ.get("MYSQL_USER", "")
MYSQL_PASSWORD = os.environ.get("MYSQL_PASSWORD", "")
MYSQL_DATABASE = os.environ.get("MYSQL_DATABASE", "tree_sitter_parser")

TENCENTCLOUD_SECRET_ID = os.environ.get("TENCENTCLOUD_SECRET_ID", "")
TENCENTCLOUD_SECRET_KEY = os.environ.get("TENCENTCLOUD_SECRET_KEY", "")
