#!/bin/bash
# DMH 回滚脚本
# 用途：在部署失败时快速回滚到上一个稳定版本

set -e
set -u

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 配置
PROJECT_NAME="dmh"
BACKUP_ROOT="/opt/code/dmh/backups"
SERVICE_NAME="dmh-api"
CURRENT_BACKUP="$BACKUP_ROOT/current"
LOG_FILE="/var/log/${SERVICE_NAME}/rollback_$(date +%Y%m%d_%H%M%S).log"

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
DMH 回滚脚本

用法: $0 [选项] [备份目录]

选项:
    -h, --help          显示帮助信息
    -l, --list          列出所有可用的备份
    -s, --status         显示当前部署状态
    -f, --force         强制回滚，不进行确认
    --dry-run            模拟回滚，不实际执行

示例:
    $0 list                         # 列出所有备份
    $0 rollback 20250130_143000      # 回滚到指定备份
    $0 rollback                    # 回滚到上一个备份
    $0 -f rollback                 # 强制回滚
    $0 -s                          # 检查当前状态

EOF
}

# 列出所有备份
list_backups() {
    log "可用的备份列表："
    echo ""

    if [ ! -d "$BACKUP_ROOT" ]; then
        warning "未找到备份根目录"
        return
    fi

    local count=0
    for backup_dir in $(ls -t "$BACKUP_ROOT"); do
        count=$((count + 1))
        backup_time=$(echo "$backup_dir" | grep -oP '\d{8}')
        backup_date=$(date -d "${backup_time:0:0:0}" +'%Y-%m-%d %H:%M:%S')
        
        # 检查备份完整性
        if [ -f "$backup_dir/$PROJECT_NAME" ]; then
            size=$(du -h "$backup_dir/$PROJECT_NAME" | cut -f1)
            status="${GREEN}可恢复${NC}"
        else
            size="-"
            status="${RED}不完整${NC}"
        fi

        # 标记当前备份
        if [ "$backup_dir" = "$CURRENT_BACKUP" ]; then
            echo "  [$count] $backup_date  $status - $size - ${YELLOW}当前版本${NC}"
        else
            echo "  [$count] $backup_date  $status - $size"
        fi
    done

    echo ""
    echo "总计: $count 个备份"
}

# 检查当前状态
check_status() {
    log "检查当前部署状态..."
    echo ""

    # 检查服务状态
    if [ -f "/tmp/${SERVICE_NAME}.pid" ]; then
        PID=$(cat "/tmp/${SERVICE_NAME}.pid")
        if ps -p "$PID" > /dev/null 2>&1; then
            echo "服务状态: ${GREEN}运行中${NC} (PID: $PID)"
            
            # 检查运行时间
            runtime=$(ps -p "$PID" -o etime= | awk '{print int($9)}')
            runtime_min=$((runtime / 60))
            echo "运行时间: ${runtime_min} 分钟"
            
            # 检查当前备份
            if [ -d "$CURRENT_BACKUP" ]; then
                current_backup_time=$(basename "$CURRENT_BACKUP" | grep -oP '\d{8}')
                current_backup_date=$(date -d "${current_backup_time:0:0:0}" +'%Y-%m-%d %H:%M:%S')
                echo "当前备份: ${current_backup_date}"
            else
                echo "当前备份: ${RED}无${NC}"
            fi
        else
            echo "服务状态: ${RED}未运行${NC}"
            rm -f "/tmp/${SERVICE_NAME}.pid"
        fi
    else
        echo "服务状态: ${YELLOW}未知${NC}"
    fi

    # 检查备份目录
    if [ ! -d "$BACKUP_ROOT" ]; then
        echo "备份目录: ${RED}不存在${NC}"
    else
        backup_count=$(ls "$BACKUP_ROOT" 2>/dev/null | wc -l)
        echo "备份目录: ${GREEN}存在${NC} ($backup_count 个备份)"
    fi

    # 检查当前二进制文件
    if [ -f "/tmp/$PROJECT_NAME" ]; then
        bin_time=$(stat -c %Y /tmp/$PROJECT_NAME | awk '{print $6" " "$7" | cut -d'.' -f1)
        bin_date=$(date -d "$bin_time" +'%Y-%m-%d %H:%M:%S')
        bin_size=$(du -h /tmp/$PROJECT_NAME | cut -f1)
        echo "二进制文件: ${GREEN}存在${NC} ($bin_date, $bin_size)"
    else
        echo "二进制文件: ${RED}不存在${NC}"
    fi
}

# 回滚到指定备份
rollback_to_backup() {
    local backup_dir="$1"
    local force="$2"

    log "准备回滚到: $backup_dir"
    echo ""

    # 检查备份目录
    if [ ! -d "$backup_dir" ]; then
        error "备份目录不存在: $backup_dir"
    fi

    # 检查备份完整性
    if [ ! -f "$backup_dir/$PROJECT_NAME" ]; then
        error "备份不完整，缺少二进制文件"
    fi

    if [ ! -f "$backup_dir/${CONFIG_FILE:-etc/dmh-api.yaml}" ]; then
        warning "备份中缺少配置文件"
    fi

    # 如果不是强制回滚，进行确认
    if [ "$force" != "--force" ]; then
        echo "即将回滚到: $(basename $backup_dir)"
        echo "  - 二进制文件: $backup_dir/$PROJECT_NAME"
        if [ -f "$backup_dir/${CONFIG_FILE:-etc/dmh-api.yaml}" ]; then
            echo "  - 配置文件: $backup_dir/${CONFIG_FILE:-etc/dmh-api.yaml}"
        fi
        echo ""
        read -p "确认回滚？(yes/no): " confirm
        if [ "$confirm" != "yes" ] && [ "$confirm" != "y" ]; then
            echo "已取消回滚"
            exit 0
        fi
    fi

    log "开始回滚..."
    echo ""

    # 停止服务
    if [ -f "/tmp/${SERVICE_NAME}.pid" ]; then
        local pid=$(cat "/tmp/${SERVICE_NAME}.pid")
        log "停止服务 (PID: $pid)..."
        kill -TERM "$pid" 2>/dev/null || true
        sleep 3
        
        # 强制停止如果仍在运行
        if ps -p "$pid" > /dev/null 2>&1; then
            log "强制停止服务..."
            kill -9 "$pid" 2>/dev/null || true
            sleep 2
        fi
    fi

    # 恢复二进制文件
    log "恢复二进制文件..."
    cp "$backup_dir/$PROJECT_NAME" "/tmp/$PROJECT_NAME"
    chmod +x "/tmp/$PROJECT_NAME"

    # 恢复配置文件
    if [ -f "$backup_dir/${CONFIG_FILE:-etc/dmh-api.yaml}" ]; then
        log "恢复配置文件..."
        cp "$backup_dir/${CONFIG_FILE:-etc/dmh-api.yaml}" "/opt/code/dmh/backend/api/${CONFIG_FILE:-etc/dmh-api.yaml}"
    fi

    # 更新当前备份链接
    if [ "$backup_dir" != "$CURRENT_BACKUP" ]; then
        log "更新当前备份链接..."
        rm -f "$CURRENT_BACKUP"
        ln -s "$backup_dir" "$CURRENT_BACKUP"
    fi

    # 重启服务
    log "重启服务..."
    cd /opt/code/dmh/backend/api
    nohup "/tmp/$PROJECT_NAME" -f "${CONFIG_FILE:-etc/dmh-api.yaml}" >> "$LOG_FILE" 2>&1 &
    local new_pid=$!
    echo "$new_pid" > "/tmp/${SERVICE_NAME}.pid"

    # 等待服务启动
    log "等待服务启动..."
    local max_wait=10
    local count=0
    while [ $count -lt $max_wait ]; do
        sleep 1
        if netstat -tlnp 2>/dev/null | grep -q ":8889.*LISTEN"; then
            log "${GREEN}回滚完成！${NC}"
            log "服务启动成功 (PID: $new_pid)"
            
            # 健康检查
            sleep 2
            local health_code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8889/health")
            if [ "$health_code" = "200" ]; then
                log "健康检查通过"
            else
                warning "健康检查失败 (HTTP $health_code)，请检查日志: $LOG_FILE"
            fi
            return 0
        fi
        count=$((count + 1))
    done

    error "服务启动超时，请检查日志: $LOG_FILE"
}

# 回滚到上一个备份
rollback_to_previous() {
    log "查找上一个备份..."

    if [ ! -d "$BACKUP_ROOT" ]; then
        error "未找到备份目录: $BACKUP_ROOT"
    fi

    # 查找上一个备份（排除 current）
    local previous_backup=$(ls -t "$BACKUP_ROOT" | grep -v "/current$" | head -1)

    if [ -z "$previous_backup" ]; then
        error "未找到可用的备份"
    fi

    rollback_to_backup "$previous_backup" "$1"
}

# 解析命令行参数
ACTION="status"
FORCE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -l|--list)
            ACTION="list"
            shift
            ;;
        -s|--status)
            ACTION="status"
            shift
            ;;
        -f|--force)
            FORCE="--force"
            shift
            ;;
        --dry-run)
            DRY_RUN="--dry-run"
            shift
            ;;
        rollback)
            ACTION="rollback"
            shift
            if [ -n "$1" ] && [ -d "$BACKUP_ROOT/$1" ]; then
                BACKUP_ARG="$1"
                shift
            fi
            ;;
        *)
            if [ -d "$BACKUP_ROOT/$1" ] || [ "$1" == "rollback" ]; then
                BACKUP_ARG="$1"
                shift
            else
                warning "未知选项: $1"
                show_help
                exit 1
            fi
            ;;
    esac
done

# 执行操作
case $ACTION in
    list)
        list_backups
        ;;
    status)
        check_status
        ;;
    rollback)
        if [ -n "$BACKUP_ARG" ]; then
            rollback_to_backup "$BACKUP_ARG" "$FORCE"
        else
            rollback_to_previous ""
        fi
        ;;
    *)
        error "未知操作: $ACTION"
        ;;
esac
