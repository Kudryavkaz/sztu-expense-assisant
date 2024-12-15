FROM golang:1.23.2-bookworm as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# 设置工作目录
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 下载依赖项和编译成二进制可执行文件
RUN go build -o ./api main.go

FROM debian:bookworm-slim

RUN set -ex && \
    echo "deb http://deb.debian.org/debian bookworm main" > /etc/apt/sources.list && \
    echo "deb http://security.debian.org/debian-security bookworm-security/updates main" >> /etc/apt/sources.list && \
    sed -i "s@http://\(deb\|security\).debian.org@[trusted=yes] http://mirrors.tencent.com@g" /etc/apt/sources.list && \
    apt update && \
    apt install apt-utils &&\
    apt-get install -y --no-install-recommends ca-certificates tzdata curl libc6-dev &&\
    ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\
    dpkg-reconfigure -f noninteractive tzdata

WORKDIR /app

COPY --from=builder /build/api ./

RUN chmod +x ./api
