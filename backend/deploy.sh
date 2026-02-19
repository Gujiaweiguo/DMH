#!/bin/bash
# DMH 后端部署脚本
# 用途：自动化 DMH 后端服务的部署流程

set -e  # 遇到错误即退出
set -u  # 使用未定义变量时报错

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
PROJECT_NAME="dmh"
BACKEND_DIR="/opt/code/dmh/backend/api"
BINARY_NAME="dmh"
CONFIG_FILE="etc/dmh-api.yaml"
SERVICE_NAME="dmh-api"
PID_FILE="/tmp/${SERVICE_NAME}.pid"
LOG_FILE="/tmp/log/${SERVICE_NAME}/deploy_$(date +%Y%m%d_%H%M%S).log"
BACKUP_DIR="/opt/code/dmh/backups/$(date +%Y%m%d_%H%M%S)"
BUILD_DIR="/tmp/${SERVICE_NAME}_build"

# 日志函数
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 帮助信息
show_help() {
    cat << EOF
DMH 后端部署脚本

用法: $0 [选项] [环境]

选项:
    -h, --help          显示帮助信息
    -b, --build         仅构建，不部署
    -d, --deploy         部署到生产环境
    -s, --start          仅启动服务
    -t, --stop           停止服务
    -r, --restart        重启服务
    -c, --check          检查服务状态
    --skip-tests        跳过构建和测试
    --skip-backup      跳过备份

环境:
    dev               开发环境
    prod              生产环境（默认）

示例:
    $0 deploy prod          # 部署到生产环境
    $0 restart              # 重启服务
    $0 stop                 # 停止服务
    $0 -h                   # 显示帮助信息

EOF
}

# 检查环境依赖
check_dependencies() {
    log "检查环境依赖..."

    # 检查 Go
    if ! command -v go &> /dev/null; then
        error "Go 未安装，请先安装 Go 1.18+"
    fi
    GO_VERSION=$(go version | awk '{print $3}')
    log "Go 版本: $GO_VERSION"

    # 检查 MySQL
    if ! command -v mysql &> /dev/null; then
        warning "MySQL 客户端未安装，但数据库连接由应用直接配置"
    else
        MYSQL_VERSION=$(mysql --version | awk '{print $5}')
        log "MySQL 客户端版本: $MYSQL_VERSION"
    fi
}

# 创建备份
create_backup() {
    log "创建备份..."

    mkdir -p "$BACKUP_DIR"

    # 备份当前运行的二进制文件
    if [ -f "/tmp/$BINARY_NAME" ]; then
        cp "/tmp/$BINARY_NAME" "$BACKUP_DIR/"
        log "已备份当前运行的二进制文件"
    fi

    # 备份配置文件
    if [ -f "$BACKEND_DIR/$CONFIG_FILE" ]; then
        cp "$BACKEND_DIR/$CONFIG_FILE" "$BACKUP_DIR/"
        log "已备份配置文件"
    fi

    # 备份数据库
    BACKUP_SQL="$BACKUP_DIR/database_backup.sql"
    if [ -f "$BACKEND_DIR/../migrations" ]; then
        log "数据库迁移脚本已存在，跳过备份"
    else
        warning "未找到迁移脚本，建议手动备份数据库"
    fi

    log "备份完成: $BACKUP_DIR"
}

# 构建应用
build_app() {
    log "开始构建应用..."

    # 创建构建目录
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"

    # 进入项目目录
    cd "$BACKEND_DIR"

    # 清理构建缓存
    go clean -modcache

    # 下载依赖
    log "下载依赖..."
    go mod tidy

    # 运行测试
    if [ "$SKIP_TESTS" != "true" ]; then
        log "运行测试..."
        go test -v ./... || warning "部分测试失败，继续构建"
    fi

    # 构建二进制文件
    log "构建二进制文件..."
    go build -ldflags="-s -w" -o "$BUILD_DIR/$BINARY_NAME" dmh.go

    # 验证构建结果
    if [ ! -f "$BUILD_DIR/$BINARY_NAME" ]; then
        error "构建失败，未找到二进制文件"
    fi

    # 检查二进制文件大小
    SIZE=$(du -h "$BUILD_DIR/$BINARY_NAME" | cut -f1)
    log "构建成功: $BUILD_DIR/$BINARY_NAME ($SIZE)"

    # 复制到临时目录
    cp "$BUILD_DIR/$BINARY_NAME" "/tmp/$BINARY_NAME"
    log "二进制文件已复制到 /tmp/$BINARY_NAME"
}

# 停止服务
stop_service() {
    log "停止服务..."

    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            log "发送 SIGTERM 到进程 $PID"
            kill -TERM "$PID"

            # 等待进程优雅退出
            local count=0
            local max_wait=30
            while ps -p "$PID" > /dev/null 2>&1 && [ $count -lt $max_wait ]; do
                sleep 1
                count=$((count + 1))
            done

            # 如果进程仍在运行，强制杀掉
            if ps -p "$PID" > /dev/null 2>&1; then
                warning "进程未正常退出，强制终止"
                kill -9 "$PID"
                sleep 2
            fi

            rm -f "$PID_FILE"
            log "服务已停止"
        else
            log "PID 文件存在但进程不在运行"
            rm -f "$PID_FILE"
        fi
    else
        log "服务未在运行"
    fi

    # 检查并清理残留进程
    sleep 1
    for pid in $(ps aux | grep "[d]mh" | grep -v grep | awk '{print $2}'); do
        warning "发现残留进程: $pid，正在清理"
        kill -9 "$pid" 2>/dev/null || true
    done
}

# 检查服务状态
check_service() {
    log "检查服务状态..."

    if [ ! -f "$PID_FILE" ]; then
        echo "状态: 未运行"
        return 1
    fi

    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        # 检查端口是否在监听
        if netstat -tlnp 2>/dev/null | grep -q ":8889.*LISTEN.*$PID"; then
            echo "状态: 运行中 (PID: $PID, 端口: 8889)"
            return 0
        else
            warning "进程存在但端口未在监听，可能启动中或异常"
            echo "状态: 异常 (PID: $PID)"
            return 1
        fi
    else
        echo "状态: 未运行 (PID 文件残留)"
        rm -f "$PID_FILE"
        return 1
    fi
}

# 启动服务
start_service() {
    log "启动服务..."

    # 检查是否已运行
    if check_service &> /dev/null; then
        warning "服务已在运行，跳过启动"
        return
    fi

    # 创建日志目录
    mkdir -p "$(dirname "$LOG_FILE")"

    # 启动服务
    log "启动 $SERVICE_NAME 服务..."
    nohup "/tmp/$BINARY_NAME" -f "$BACKEND_DIR/$CONFIG_FILE" >> "$LOG_FILE" 2>&1 &
    SERVICE_PID=$!

    # 保存 PID
    echo "$SERVICE_PID" > "$PID_FILE"

    # 等待服务启动
    log "等待服务启动..."
    local max_wait=10
    local count=0
    while [ $count -lt $max_wait ]; do
        sleep 1
        if netstat -tlnp 2>/dev/null | grep -q ":8889.*LISTEN"; then
            log "服务启动成功 (PID: $SERVICE_PID, 端口: 8889)"
            return 0
        fi
        count=$((count + 1))
    done

    error "服务启动超时，请检查日志: $LOG_FILE"
}

# 部署服务
deploy_service() {
    log "开始部署服务..."

    # 检查环境
    if [ "$ENV" = "prod" ]; then
        log "生产环境部署"
    else
        log "开发环境部署"
    fi

    # 执行备份
    if [ "$SKIP_BACKUP" != "true" ]; then
        create_backup
    fi

    # 构建应用
    build_app

    # 停止旧服务
    stop_service

    # 启动新服务
    start_service

    # 健康检查
    log "执行健康检查..."
    sleep 3
    HEALTH_URL="http://localhost:8889/health"
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$HEALTH_URL" || echo "000")

    if [ "$HTTP_CODE" = "200" ]; then
        log "${GREEN}部署成功！${NC}"
        log "健康检查通过: $HEALTH_URL"
    else
        error "健康检查失败 (HTTP $HTTP_CODE)，请检查日志: $LOG_FILE"
    fi
}

# 回滚操作
rollback() {
    log "开始回滚操作..."

    # 停止当前服务
    stop_service

    # 检查备份目录
    if [ ! -d "$BACKUP_DIR" ]; then
        error "未找到备份目录: $BACKUP_DIR"
    fi

    # 恢复二进制文件
    if [ -f "$BACKUP_DIR/$BINARY_NAME" ]; then
        cp "$BACKUP_DIR/$BINARY_NAME" "/tmp/$BINARY_NAME"
        log "已恢复二进制文件"
    fi

    # 恢复配置文件
    if [ -f "$BACKUP_DIR/$CONFIG_FILE" ]; then
        cp "$BACKUP_DIR/$CONFIG_FILE" "$BACKEND_DIR/$CONFIG_FILE"
        log "已恢复配置文件"
    fi

    # 重启服务
    start_service

    log "${GREEN}回滚完成！${NC}"
}

# 解析命令行参数
ENV="prod"
BUILD_ONLY="false"
SKIP_TESTS="false"
SKIP_BACKUP="false"
ACTION="deploy"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -b|--build)
            BUILD_ONLY="true"
            ACTION="build"
            shift
            ;;
        -d|--deploy)
            ACTION="deploy"
            shift
            ;;
        -s|--start)
            ACTION="start"
            shift
            ;;
        -t|--stop)
            ACTION="stop"
            shift
            ;;
        -r|--restart)
            ACTION="restart"
            shift
            ;;
        -c|--check)
            ACTION="check"
            shift
            ;;
        --skip-tests)
            SKIP_TESTS="true"
            shift
            ;;
        --skip-backup)
            SKIP_BACKUP="true"
            shift
            ;;
        dev|prod)
            ENV="$1"
            shift
            ;;
        *)
            warning "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 根据操作执行
case $ACTION in
    deploy)
        check_dependencies
        deploy_service
        ;;
    build)
        check_dependencies
        build_app
        ;;
    start)
        check_dependencies
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        check_dependencies
        stop_service
        start_service
        ;;
    check)
        check_service
        ;;
    rollback)
        rollback
        ;;
    *)
        error "未知操作: $ACTION"
        ;;
esac
