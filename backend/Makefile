.PHONY: build run clean proto docker-up docker-down setup dev dev-stop dev-restart dev-status dev-logs

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod

# Service names
SERVICES=gateway auth upload diagnosis report

# Dev script
DEV_SCRIPT=bash scripts/dev.sh

build: $(SERVICES)

gateway:
	cd gateway && $(GOBUILD) -o ../bin/gateway gateway.go

auth:
	cd auth && $(GOBUILD) -o ../bin/auth auth.go

upload:
	cd upload && $(GOBUILD) -o ../bin/upload upload.go

diagnosis:
	cd diagnosis && $(GOBUILD) -o ../bin/diagnosis diagnosis.go

report:
	cd report && $(GOBUILD) -o ../bin/report report.go

run-gateway: gateway
	./bin/gateway -f gateway/etc/gateway.yaml

run-auth: auth
	./bin/auth -f auth/etc/auth.yaml

run-upload: upload
	./bin/upload -f upload/etc/upload.yaml

run-diagnosis: diagnosis
	./bin/diagnosis -f diagnosis/etc/diagnosis.yaml

run-report: report
	./bin/report -f report/etc/report.yaml

# Generate gRPC code from proto files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/auth.proto proto/upload/upload.proto proto/diagnosis/diagnosis.proto proto/report/report.proto

# ============ 开发模式命令 ============

# 启动所有微服务（后台运行，带状态面板）
dev:
	$(DEV_SCRIPT) start

# 停止所有微服务
dev-stop:
	$(DEV_SCRIPT) stop

# 重启所有微服务
dev-restart:
	$(DEV_SCRIPT) restart

# 查看服务运行状态
dev-status:
	$(DEV_SCRIPT) status

# 实时查看所有服务日志
dev-logs:
	$(DEV_SCRIPT) logs

# ============ Docker 命令 ============
docker-up:
	cd deploy && docker-compose up -d

docker-down:
	cd deploy && docker-compose down

docker-build:
	docker-compose -f deploy/docker-compose.yml build

# Setup WSL environment
setup:
	bash deploy/setup-wsl.sh

clean:
	rm -rf bin/ logs/ .pids/

tidy:
	$(GOMOD) tidy

test:
	$(GOCMD) test ./... -v
