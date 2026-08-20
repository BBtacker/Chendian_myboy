#!/bin/bash
# ============================================
# WSL 环境搭建脚本 - 腺样体面容智能筛查系统
# 在 WSL 中运行: bash deploy/setup-wsl.sh
# ============================================

set -e

echo "=========================================="
echo "  腺样体面容智能筛查系统 - 环境搭建"
echo "=========================================="

# 检查是否在WSL中运行
if ! grep -qi microsoft /proc/version 2>/dev/null; then
    echo "[警告] 当前可能不在WSL环境中，但脚本仍可继续执行"
fi

# 1. 安装 Go
install_go() {
    GO_VERSION="1.22.12"
    GO_INSTALL_DIR="/usr/local/go"
    
    if command -v go &>/dev/null; then
        CURRENT_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        echo "[INFO] 已安装 Go $CURRENT_VERSION"
        if [[ "$CURRENT_VERSION" > "1.22" ]] || [[ "$CURRENT_VERSION" == "1.22"* ]]; then
            echo "[INFO] Go版本满足要求，跳过安装"
            return
        fi
    fi
    
    echo "[1/6] 安装 Go $GO_VERSION ..."
    cd /tmp
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O go.tar.gz
    sudo rm -rf $GO_INSTALL_DIR
    sudo tar -C /usr/local -xzf go.tar.gz
    rm go.tar.gz
    
    # 添加到PATH
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin:~/go/bin' >> ~/.bashrc
    fi
    export PATH=$PATH:/usr/local/go/bin:~/go/bin
    
    echo "[OK] Go 安装完成: $(go version)"
}

# 2. 安装 protoc
install_protoc() {
    PROTOC_VERSION="25.1"
    
    if command -v protoc &>/dev/null; then
        echo "[2/6] protoc 已安装: $(protoc --version)"
        return
    fi
    
    echo "[2/6] 安装 protoc ..."
    cd /tmp
    wget -q "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" -O protoc.zip
    sudo unzip -o protoc.zip -d /usr/local/protoc
    sudo ln -sf /usr/local/protoc/bin/protoc /usr/local/bin/protoc
    rm protoc.zip
    
    echo "[OK] protoc 安装完成: $(protoc --version)"
}

# 3. 安装 goctl 和 protoc 插件
install_tools() {
    export PATH=$PATH:/usr/local/go/bin:~/go/bin
    
    echo "[3/6] 安装 goctl 和 protoc 插件 ..."
    
    # 安装 goctl
    if ! command -v goctl &>/dev/null; then
        go install github.com/zeromicro/go-zero/tools/goctl@latest
    fi
    echo "[OK] goctl 安装完成"
    
    # 安装 protoc-gen-go 和 protoc-gen-go-grpc
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    echo "[OK] protoc 插件安装完成"
}

# 4. 生成 gRPC 代码
gen_proto() {
    export PATH=$PATH:/usr/local/go/bin:~/go/bin
    
    echo "[4/6] 生成 gRPC 代码 ..."
    cd /tmp/faceTest/backend || cd "$(dirname "$0")/.."
    
    # proto 文件已在子目录中，直接生成
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/auth/auth.proto
    
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/upload/upload.proto
    
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/diagnosis/diagnosis.proto
    
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/report/report.proto
    
    echo "[OK] gRPC 代码生成完成"
}

# 5. 下载依赖
install_deps() {
    export PATH=$PATH:/usr/local/go/bin:~/go/bin
    
    echo "[5/6] 下载 Go 依赖 ..."
    cd "$(dirname "$0")/.."
    go mod tidy
    echo "[OK] 依赖下载完成"
}

# 6. 构建所有服务
build_all() {
    export PATH=$PATH:/usr/local/go/bin:~/go/bin
    
    echo "[6/6] 构建所有微服务 ..."
    cd "$(dirname "$0")/.."
    mkdir -p bin
    
    for svc in gateway auth upload diagnosis report; do
        echo "  构建 $svc ..."
        cd $svc && go build -o ../bin/$svc . && cd ..
    done
    
    echo "[OK] 所有服务构建完成"
}

# 主流程
main() {
    install_go
    install_protoc
    install_tools
    gen_proto
    install_deps
    build_all
    
    echo ""
    echo "=========================================="
    echo "  环境搭建完成！"
    echo "=========================================="
    echo ""
    echo "后续步骤："
    echo "  1. 启动基础设施: cd deploy && docker-compose up -d"
    echo "  2. 启动微服务: make run-auth & make run-upload & make run-diagnosis & make run-report &"
    echo "  3. 启动网关: make run-gateway"
    echo "  4. 启动前端: cd ../front && npm run dev"
    echo ""
}

main
