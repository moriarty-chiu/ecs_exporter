# -------- Stage 1: Build --------
FROM golang:1.22 AS builder

# 设置 Go Module 和代理
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# 拷贝 go.mod 和 go.sum 先执行 mod 下载
COPY go.mod go.sum ./
RUN go mod download

# 拷贝全部源代码
COPY . .

# 构建二进制
RUN go build -o ecs_exporter .

# -------- Stage 2: Final Image --------
FROM alpine:latest

# 安装 curl 用于健康检查（可选）
RUN apk add --no-cache curl ca-certificates

WORKDIR /app

# 拷贝构建好的二进制
COPY --from=builder /app/ecs_exporter .

# 创建非 root 用户运行程序（可选）
RUN adduser -D -g '' ecsuser
USER ecsuser

# 设置暴露端口
EXPOSE 9100

# 启动程序
ENTRYPOINT ["./ecs_exporter"]

# 健康检查（每 30 秒一次）
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:9100/healthz || exit 1
