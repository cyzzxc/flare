# 多阶段构建 Dockerfile for Flare
# 用于快速构建和推送到 Docker Hub

# 阶段 1: 构建阶段
FROM golang:1.22.2-alpine3.19 AS builder

# 设置工作目录
WORKDIR /app

# 安装编译依赖
RUN apk add --no-cache git bash gcc musl-dev upx

# 配置 Go 环境
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn
ENV TZ=Asia/Shanghai

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 运行构建脚本（生成内嵌资源）
RUN go run build/build.go

# 编译二进制文件
RUN export VERSION=$(git describe --tags --always 2>/dev/null || echo "dev") && \
    export COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    export BUILD_DATE=$(date +%FT%T%z) && \
    echo "Building version: $VERSION / $COMMIT / $BUILD_DATE" && \
    go build -ldflags "-w -s \
        -X 'github.com/soulteary/flare/internal/version.Version=$VERSION' \
        -X 'github.com/soulteary/flare/internal/version.Commit=$COMMIT' \
        -X 'github.com/soulteary/flare/internal/version.BuildDate=$BUILD_DATE'" \
        -o flare main.go

# 压缩二进制文件（可选，减小镜像大小）
RUN upx -9 -o flare.minify flare && mv flare.minify flare

# 阶段 2: 运行阶段
FROM alpine:3.19

# 安装运行时依赖
RUN apk add --no-cache tzdata ca-certificates

# 设置时区
ENV TZ=Asia/Shanghai
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    rm -rf /var/cache/apk/*

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/flare /bin/flare

# 设置工作目录（配置文件存放位置）
WORKDIR /app

# 暴露默认端口
EXPOSE 5005

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:5005/health || exit 1

# 启动应用
CMD ["flare"]
