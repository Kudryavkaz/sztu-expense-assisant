FROM python:3.12-slim

RUN pip config set global.index-url https://mirrors.cloud.tencent.com/pypi/simple

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

RUN pip install poetry

WORKDIR /app

COPY . /app

RUN poetry install --no-root
