#!/bin/bash
# DMH 容器化回滚脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }
info() { echo -e "${BLUE}[INFO]${NC} $1"; }

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENT_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$DEPLOYMENT_DIR")"
BACKEND_DIR="$PROJECT_ROOT/backend"
COMPOSE_FILE="$DEPLOYMENT_DIR/docker-compose.yml"

# 步骤1: 确认操作
echo "========================================="
warn "即将停止容器化部署，回滚到独立进程"
warn "此操作将："
warn "  1. 停止并删除Docker容器"
warn "  2. 恢复使用deploy.sh部署脚本"
warn "  3. 后端服务以独立进程方式运行"
warn "  4. 前端静态文件需要单独部署（如nginx）"
echo "========================================="
read -p "确认继续？(yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    info "操作已取消"
    exit 0
fi

# 步骤2: 停止Docker Compose
log "步骤 1: 停止 Docker Compose..."
cd "$DEPLOYMENT_DIR"
if docker compose -f "$COMPOSE_FILE" ps 2>/dev/null | grep -q "dmh-nginx\|dmh-api"; then
    info "停止容器..."
    docker compose -f "$COMPOSE_FILE" down
    log "容器已停止"
else
    info "没有运行中的容器"
fi

# 步骤3: 部署后端服务
log "步骤 2: 部署后端服务..."
cd "$BACKEND_DIR"
if [ -f "deploy.sh" ]; then
    info "构建并部署服务..."
    ./deploy.sh -d prod || error "部署失败"
else
    warn "deploy.sh 不存在，跳过后端部署"
fi

# 步骤4: 验证服务
log "步骤 3: 验证服务..."
sleep 3

# 测试API登录
info "测试 API 登录..."
LOGIN_RESULT=$(curl -s -X POST http://localhost:8889/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"123456"}')

if echo "$LOGIN_RESULT" | grep -q "token"; then
    log "✓ API 登录测试通过"
else
    warn "API 登录测试失败"
fi

# 步骤5: 显示管理命令
log "========================================="
log "回滚完成！"
log "========================================="
info ""
info "当前部署方式: 独立进程（后端）"
info "服务端口: 8889"
info ""
info "后端管理命令："
info "  查看服务状态:  cd $BACKEND_DIR && ./deploy.sh -c"
info "  重启服务:      ./deploy.sh -r"
info "  停止服务:      ./deploy.sh -t"
info "  查看日志:      tail -f /tmp/log/dmh-api/deploy_*.log"
info ""
warn "注意: 前端静态文件需要单独配置nginx部署！"
warn "如需切换回容器部署："
warn "  cd $DEPLOYMENT_DIR/scripts && ./deploy-containers.sh"
warn "========================================="
