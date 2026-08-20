# 图像上传服务 Dockerfile
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

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/upload ./upload

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/bin/upload /app/upload
COPY --from=builder /app/deploy/configs/upload-docker.yaml /app/etc/upload.yaml

RUN mkdir -p /app/uploads

EXPOSE 8082

CMD ["/app/upload", "-f", "/app/etc/upload.yaml"]
