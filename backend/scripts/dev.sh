#!/bin/bash
#
# 微服务开发模式管理脚本
# 用法:
#   ./scripts/dev.sh start    - 启动所有微服务
#   ./scripts/dev.sh stop     - 停止所有微服务
#   ./scripts/dev.sh restart  - 重启所有微服务
#   ./scripts/dev.sh status   - 查看服务状态
#   ./scripts/dev.sh logs     - 实时查看所有服务日志
#

# 不使用 set -e，让脚本在单个服务失败时继续处理其他服务

# ============ 配置 ============
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BIN_DIR="$ROOT_DIR/bin"
LOG_DIR="$ROOT_DIR/logs"
PID_DIR="$ROOT_DIR/.pids"

mkdir -p "$LOG_DIR" "$PID_DIR" "$BIN_DIR"

# 服务定义: 名称|二进制|配置文件|端口
SERVICES=(
    "auth|auth|auth/etc/auth.yaml|8081"
    "upload|upload|upload/etc/upload.yaml|8082"
    "diagnosis|diagnosis|diagnosis/etc/diagnosis.yaml|8083"
    "report|report|report/etc/report.yaml|8084"
    "gateway|gateway|gateway/etc/gateway.yaml|8080"
)

# 颜色定义（使用 $'...' 让 bash 在赋值时解析转义字符）
RED=$'\033[0;31m'
GREEN=$'\033[0;32m'
YELLOW=$'\033[1;33m'
CYAN=$'\033[0;36m'
GRAY=$'\033[0;90m'
NC=$'\033[0m' # No Color
BOLD=$'\033[1m'

# ============ 工具函数 ============

# 确保 go 命令可用
ensure_go() {
    if ! command -v go &>/dev/null; then
        # 尝试加载 Go 路径
        export PATH=$PATH:/usr/local/go/bin:~/go/bin
    fi
    if ! command -v go &>/dev/null; then
        echo -e "  ${RED}✗ 未找到 go 命令，请先运行: bash deploy/setup-wsl.sh${NC}"
        exit 1
    fi
}

get_pid() {
    local name=$1
    local pid_file="$PID_DIR/${name}.pid"
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            echo "$pid"
            return 0
        fi
    fi
    echo ""
    return 1
}

check_port() {
    local port=$1
    if command -v ss &>/dev/null; then
        ss -tlnp 2>/dev/null | grep -q ":${port} " && return 0
    elif command -v netstat &>/dev/null; then
        netstat -tlnp 2>/dev/null | grep -q ":${port} " && return 0
    fi
    return 1
}

# 检查基础设施端口是否可达（Docker容器：MySQL/etcd/Redis）
check_infra() {
    local missing=""
    local infra=(
        "MySQL|3307"
        "etcd|2379"
        "Redis|6380"
    )

    for item in "${infra[@]}"; do
        IFS='|' read -r name port <<< "$item"
        if ! check_port "$port"; then
            missing="${missing}  ${RED}✗ $name (${port})${NC}\n"
        fi
    done

    if [ -n "$missing" ]; then
        echo -e "  ${YELLOW}⚠ 以下基础设施未启动:${NC}"
        echo -e "$missing"
        echo -e "  ${YELLOW}请先在 deploy 目录启动 Docker 容器:${NC}"
        echo -e "  ${CYAN}  cd backend/deploy && docker-compose up -d mysql redis etcd rabbitmq milvus${NC}"
        echo -e "  ${CYAN}  docker-compose ps   # 等待所有服务 healthy${NC}"
        echo -e "  ${YELLOW}（注意: 先启动 Docker Desktop）${NC}"
        return 1
    fi
    return 0
}

print_header() {
    echo ""
    echo -e "${BOLD}========================================================${NC}"
    echo -e "${BOLD}  腺样体面容筛查系统 - 微服务管理${NC}"
    echo -e "${BOLD}========================================================${NC}"
    echo ""
}

print_status_table() {
    printf "\n"
    printf "  ${BOLD}%-12s %-8s %-10s %-12s %-8s${NC}\n" "服务" "端口" "进程状态" "端口监听" "PID"
    printf "  ${GRAY}------------------------------------------------------------${NC}\n"

    local all_running=true

    for svc in "${SERVICES[@]}"; do
        IFS='|' read -r name bin config port <<< "$svc"

        local pid=$(get_pid "$name")
        local pid_status process_status port_status

        if [ -n "$pid" ]; then
            pid_status="${GREEN}运行中${NC}"
            process_status="RUNNING"
        else
            pid_status="${RED}已停止${NC}"
            process_status="STOPPED"
            all_running=false
        fi

        if check_port "$port"; then
            port_status="${GREEN}✓ 监听${NC}"
        else
            port_status="${RED}✗ 未监听${NC}"
            all_running=false
        fi

        local pid_display=$([ -n "$pid" ] && echo "$pid" || echo "-")

        printf "  %-16s %-10s %-14s %-14s %-8s\n" \
            "$name" "$port" "$pid_status" "$port_status" "$pid_display"
    done

    printf "\n"

    if [ "$all_running" = true ]; then
        echo -e "  ${GREEN}${BOLD}✓ 所有服务运行正常${NC}"
    else
        echo -e "  ${YELLOW}${BOLD}⚠ 部分服务未正常运行${NC}"
    fi
    echo ""
}

# ============ 命令实现 ============

cmd_start() {
    print_header
    ensure_go

    # 检查基础设施（MySQL/etcd/Redis）是否已启动
    if ! check_infra; then
        return 1
    fi

    # 先检查是否已有服务在运行
    local running_count=0
    for svc in "${SERVICES[@]}"; do
        IFS='|' read -r name bin config port <<< "$svc"
        if [ -n "$(get_pid "$name")" ]; then
            ((running_count++))
        fi
    done

    if [ "$running_count" -gt 0 ]; then
        echo -e "  ${YELLOW}⚠ 检测到 $running_count 个服务已在运行${NC}"
        echo -e "  ${GRAY}如需重启请先执行: make dev-stop${NC}"
        print_status_table
        return 0
    fi

    echo -e "  ${CYAN}正在启动微服务...${NC}\n"

    # 检查 go.sum 是否存在，不存在则先 go mod tidy
    if [ ! -f "$ROOT_DIR/go.sum" ]; then
        echo -e "  ${YELLOW}⚠ go.sum 不存在，正在下载依赖 (go mod tidy)...${NC}"
        local tidy_output
        tidy_output=$(cd "$ROOT_DIR" && go mod tidy 2>&1)
        local tidy_status=$?
        if [ $tidy_status -ne 0 ] || [ ! -f "$ROOT_DIR/go.sum" ]; then
            echo -e "  ${RED}✗ go mod tidy 失败${NC}"
            echo -e "  ${GRAY}--- 错误信息 ---${NC}"
            echo "$tidy_output" | head -30 | sed 's/^/  /'
            echo -e "  ${GRAY}----------------${NC}"
            echo -e "  ${YELLOW}请手动修复后运行: cd backend && go mod tidy${NC}"
            return 1
        fi
        echo -e "  ${GREEN}✓ 依赖下载完成${NC}\n"
    fi

    local success_count=0
    local fail_count=0

    for svc in "${SERVICES[@]}"; do
        IFS='|' read -r name bin config port <<< "$svc"

        local binary="$BIN_DIR/$bin"
        local config_file="$ROOT_DIR/$config"
        local log_file="$LOG_DIR/${name}.log"
        local pid_file="$PID_DIR/${name}.pid"

        # 检查二进制是否存在，或源码是否有更新（有更新则重新编译）
        local need_build=false
        if [ ! -f "$binary" ]; then
            need_build=true
            echo -e "  ${GRAY}[$name] 二进制不存在，正在编译...${NC}"
        elif [ -n "$(find "$ROOT_DIR/$name" "$ROOT_DIR/common" "$ROOT_DIR/proto" -name '*.go' -newer "$binary" 2>/dev/null | head -1)" ]; then
            need_build=true
            echo -e "  ${GRAY}[$name] 源码有更新，正在重新编译...${NC}"
        fi

        if [ "$need_build" = true ]; then
            local build_output
            build_output=$(cd "$ROOT_DIR/$name" && go build -o "$binary" "$name.go" 2>&1)
            local build_status=$?
            if [ $build_status -ne 0 ]; then
                echo "$build_output" >> "$log_file"
                echo -e "  ${RED}✗ [$name] 编译失败${NC}"
                echo -e "  ${GRAY}--- 编译错误 ---${NC}"
                echo "$build_output" | head -20 | sed 's/^/  /'
                echo -e "  ${GRAY}----------------${NC}"
                ((fail_count++))
                continue
            fi
            echo -e "  ${GREEN}✓ [$name] 编译成功${NC}"
        fi

        # 检查配置文件
        if [ ! -f "$config_file" ]; then
            echo -e "  ${RED}✗ [$name] 配置文件不存在: $config${NC}"
            ((fail_count++))
            continue
        fi

        # 启动服务
        echo -e "  ${CYAN}→ [$name] 启动中 (端口 $port)...${NC}"

        # 清空旧日志
        > "$log_file"

        nohup "$binary" -f "$config_file" >> "$log_file" 2>&1 &
        local pid=$!
        echo "$pid" > "$pid_file"

        # 等待进程启动 (最多3秒)
        sleep 1

        if kill -0 "$pid" 2>/dev/null; then
            # 检查是否立即退出（比如配置错误）
            wait_time=0
            while [ $wait_time -lt 2 ]; do
                if ! kill -0 "$pid" 2>/dev/null; then
                    break
                fi
                sleep 1
                ((wait_time++))
            done

            if kill -0 "$pid" 2>/dev/null; then
                echo -e "  ${GREEN}✓ [$name] 启动成功 (PID: $pid, 端口: $port)${NC}"
                ((success_count++))
            else
                echo -e "  ${RED}✗ [$name] 启动后立即退出，请检查日志: logs/${name}.log${NC}"
                rm -f "$pid_file"
                # 显示最后几行错误日志
                echo -e "  ${GRAY}--- 日志摘要 ---${NC}"
                tail -5 "$log_file" 2>/dev/null | sed 's/^/  /'
                echo -e "  ${GRAY}----------------${NC}"
                ((fail_count++))
            fi
        else
            echo -e "  ${RED}✗ [$name] 启动失败，请检查日志: logs/${name}.log${NC}"
            rm -f "$pid_file"
            tail -5 "$log_file" 2>/dev/null | sed 's/^/  /'
            ((fail_count++))
        fi
    done

    echo ""
    echo -e "  ${BOLD}启动结果: ${GREEN}$success_count 成功${NC}${BOLD}, ${RED}$fail_count 失败${NC}"

    # 最终状态表
    print_status_table

    if [ "$fail_count" -gt 0 ]; then
        echo -e "  ${YELLOW}提示: 使用 ${BOLD}make dev-logs${NC}${YELLOW} 查看实时日志${NC}"
        echo -e "  ${YELLOW}提示: 使用 ${BOLD}make dev-stop${NC}${YELLOW} 停止所有服务${NC}"
    fi
}

cmd_stop() {
    print_header

    echo -e "  ${CYAN}正在停止微服务...${NC}\n"

    local stopped=0
    local not_running=0

    for svc in "${SERVICES[@]}"; do
        IFS='|' read -r name bin config port <<< "$svc"

        local pid=$(get_pid "$name")

        if [ -n "$pid" ]; then
            # 先尝试优雅关闭
            kill "$pid" 2>/dev/null
            echo -e "  ${CYAN}→ [$name] 发送 SIGTERM (PID: $pid)...${NC}"

            # 等待3秒
            local wait_time=0
            while [ $wait_time -lt 3 ]; do
                if ! kill -0 "$pid" 2>/dev/null; then
                    break
                fi
                sleep 1
                ((wait_time++))
            done

            # 如果还在运行，强制杀
            if kill -0 "$pid" 2>/dev/null; then
                kill -9 "$pid" 2>/dev/null
                echo -e "  ${YELLOW}⚠ [$name] 强制终止 (SIGKILL)${NC}"
            else
                echo -e "  ${GREEN}✓ [$name] 已停止${NC}"
            fi

            rm -f "$PID_DIR/${name}.pid"
            ((stopped++))
        else
            echo -e "  ${GRAY}○ [$name] 未在运行${NC}"
            ((not_running++))
        fi
    done

    echo ""
    echo -e "  ${BOLD}停止结果: ${GREEN}$stopped 个已停止${NC}, ${GRAY}$not_running 个未运行${NC}"
    echo ""
}

cmd_status() {
    print_header
    print_status_table
}

cmd_logs() {
    print_header
    echo -e "  ${CYAN}实时日志 (Ctrl+C 退出):${NC}\n"

    # 检查日志文件是否存在
    local log_files=()
    for svc in "${SERVICES[@]}"; do
        IFS='|' read -r name bin config port <<< "$svc"
        local log_file="$LOG_DIR/${name}.log"
        if [ -f "$log_file" ]; then
            log_files+=("$log_file")
        fi
    done

    if [ ${#log_files[@]} -eq 0 ]; then
        echo -e "  ${YELLOW}没有找到日志文件，服务可能未启动过${NC}"
        echo ""
        return 0
    fi

    # 使用 tail -f 跟踪所有日志，加上前缀
    local tail_args=()
    for log_file in "${log_files[@]}"; do
        local name=$(basename "$log_file" .log)
        tail_args+=("-f" "$log_file")
    done

    # 用 awk 加上服务名前缀
    tail -f "${log_files[@]}" 2>/dev/null | awk \
        -v LOG_DIR="$LOG_DIR" \
        '
        {
            # 根据文件名确定服务
            fname = FILENAME
            sub(LOG_DIR "/", "", fname)
            sub(".log", "", fname)
            printf "[%s] %s\n", fname, $0
            fflush()
        }
        ' || {
            # 如果上面的方式不支持 FILENAME，退回简单模式
            echo -e "  ${YELLOW}多文件 tail 不支持，切换为合并模式...${NC}"
            tail -f "${log_files[@]}"
        }
}

cmd_restart() {
    cmd_stop
    echo -e "  ${GRAY}等待 2 秒后重启...${NC}"
    sleep 2
    cmd_start
}

# ============ 主入口 ============

case "${1:-}" in
    start)
        cmd_start
        ;;
    stop)
        cmd_stop
        ;;
    restart)
        cmd_restart
        ;;
    status)
        cmd_status
        ;;
    logs)
        cmd_logs
        ;;
    *)
        echo "用法: $0 {start|stop|restart|status|logs}"
        echo ""
        echo "  start   - 启动所有微服务"
        echo "  stop    - 停止所有微服务"
        echo "  restart - 重启所有微服务"
        echo "  status  - 查看服务运行状态"
        echo "  logs    - 实时查看所有服务日志"
        exit 1
        ;;
esac
