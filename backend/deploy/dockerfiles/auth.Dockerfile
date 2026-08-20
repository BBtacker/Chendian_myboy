# 认证服务 Dockerfile
FROM golang:1.22-alpine AS builder

# 国内 Go 模块镜像，避免 proxy.golang.org 连不上
ENV GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct
ENV GOSUMDB=off
ENV GOFLAGS=-mod=mod

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# .pb.go 已提交进仓库，无需在容器内重新生成 protobuf 代码

# 编译认证服务
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/auth ./auth

# 运行阶段
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/bin/auth /app/auth
COPY --from=builder /app/deploy/configs/auth-docker.yaml /app/etc/auth.yaml

EXPOSE 8081

CMD ["/app/auth", "-f", "/app/etc/auth.yaml"]
